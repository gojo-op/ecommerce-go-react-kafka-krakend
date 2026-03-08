package main

import (
    "log"
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "order-service/internal/events"
    "order-service/internal/config"
    "order-service/internal/controllers"
    "order-service/internal/models"
    "order-service/internal/services"
    odb "order-service/internal/db"
)

func main() {
    cfg, err := config.Load()
    if err != nil { log.Fatalf("config error: %v", err) }
    gdb, err := odb.Open()
    if err != nil { log.Fatalf("db error: %v", err) }
    pub := events.New()
    _ = gdb.AutoMigrate(&models.Order{}, &models.OrderItem{})
    svc := services.New(gdb, pub)
    ctl := controllers.New(svc)
    r := gin.New()
    r.Use(gin.Logger(), gin.Recovery())
    r.GET("/health", func(c *gin.Context){ c.JSON(http.StatusOK, gin.H{"status":"healthy","service":"order-service","timestamp": time.Now().Unix()}) })
    r.HEAD("/health", func(c *gin.Context){ c.Status(http.StatusOK) })
    api := r.Group("/api/v1")
    api.POST("/orders", ctl.Create)
    api.POST("/orders/checkout", ctl.Checkout)
    api.GET("/orders", ctl.ListByUser)
    api.GET("/orders/:id", ctl.Get)
    api.PATCH("/orders/:id/status", ctl.UpdateStatus)
    srv := &http.Server{ Addr: ":" + cfg.ServicePort, Handler: r }
    // No cross-service seeding; orders are created by user actions
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { log.Fatalf("server error: %v", err) }
}
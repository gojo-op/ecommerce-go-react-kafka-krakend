package main

import (
    "log"
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "cart-service/internal/events"
    "cart-service/internal/config"
    "cart-service/internal/controllers"
    "cart-service/internal/services"
    cdb "cart-service/internal/db"
    "cart-service/internal/models"
)

func main() {
    cfg, err := config.Load()
    if err != nil { log.Fatalf("config error: %v", err) }
    gdb, err := cdb.Open()
    if err != nil { log.Fatalf("db error: %v", err) }
    _ = gdb.AutoMigrate(&models.CartEntity{}, &models.CartItemEntity{})
    pub := events.New()

    svc := services.New(gdb, pub)
    ctl := controllers.New(svc)

    r := gin.New()
    r.Use(gin.Logger(), gin.Recovery())
    r.GET("/health", func(c *gin.Context){ c.JSON(http.StatusOK, gin.H{"status":"healthy","service":"cart-service","timestamp": time.Now().Unix()}) })
    r.HEAD("/health", func(c *gin.Context){ c.Status(http.StatusOK) })
    api := r.Group("/api/v1")
    api.GET("/cart/:user_id", ctl.Get)
    api.POST("/cart/:user_id/items", ctl.AddItem)
    api.DELETE("/cart/:user_id/items/:sku", ctl.RemoveItem)
    api.DELETE("/cart/:user_id", ctl.Clear)
    srv := &http.Server{ Addr: ":" + cfg.ServicePort, Handler: r }
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { log.Fatalf("server error: %v", err) }
}
package main

import (
    "context"
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "notification-service/internal/config"
    "notification-service/internal/controllers"
    "notification-service/internal/services"
)

func main() {
    cfg, err := config.Load()
    if err != nil { panic(err) }
    svc := services.New()
    handler := services.NewHandler(svc)
    topics := []string{ "order.created", "order.status_changed", "payment.processed", "payment.failed" }
    consumer, err := services.StartConsumer(context.Background(), topics, handler)
    if err != nil { panic(err) }
    defer consumer.Close()
    r := gin.New()
    r.Use(gin.Logger(), gin.Recovery())
    r.GET("/health", func(c *gin.Context){ c.JSON(http.StatusOK, gin.H{"status":"healthy","service":"notification-service","timestamp": time.Now().Unix()}) })
    r.HEAD("/health", func(c *gin.Context){ c.Status(http.StatusOK) })
    api := r.Group("/api/v1")
    ctl := controllers.New(svc)
    api.GET("/notifications", ctl.List)
    srv := &http.Server{ Addr: ":" + cfg.ServicePort, Handler: r }
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { panic(err) }
}
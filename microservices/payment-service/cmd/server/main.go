package main

import (
    "log"
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "payment-service/internal/config"
    "payment-service/internal/controllers"
    "payment-service/internal/services"
    "payment-service/internal/events"
)

func main() {
    cfg, err := config.Load()
    if err != nil { log.Fatalf("config error: %v", err) }
    pub := events.New()
    svc := services.New(pub)
    ctl := controllers.New(svc)
    r := gin.New()
    r.Use(gin.Logger(), gin.Recovery())
    r.GET("/health", func(c *gin.Context){ c.JSON(http.StatusOK, gin.H{"status":"healthy","service":"payment-service","timestamp": time.Now().Unix()}) })
    r.HEAD("/health", func(c *gin.Context){ c.Status(http.StatusOK) })
    api := r.Group("/api/v1")
    api.POST("/payments/intent", ctl.CreateIntent)
    api.POST("/payments/webhooks/stripe", ctl.StripeWebhook)
    api.POST("/payments/webhooks/razorpay", ctl.RazorpayWebhook)
    srv := &http.Server{ Addr: ":" + cfg.ServicePort, Handler: r }
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { log.Fatalf("server error: %v", err) }
}
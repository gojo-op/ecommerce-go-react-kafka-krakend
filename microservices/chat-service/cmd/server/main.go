package main

import (
    "log"
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "chat-service/internal/config"
    "chat-service/internal/controllers"
    "chat-service/internal/events"
)

func main() {
    cfg, err := config.Load()
    if err != nil { log.Fatalf("config error: %v", err) }
    pub := events.New()
    ctl := controllers.New(pub)
    r := gin.New()
    r.Use(gin.Logger(), gin.Recovery())
    r.GET("/health", func(c *gin.Context){ c.JSON(http.StatusOK, gin.H{"status":"healthy","service":"chat-service","timestamp": time.Now().Unix()}) })
    r.HEAD("/health", func(c *gin.Context){ c.Status(http.StatusOK) })
    r.GET("/ws", ctl.WS)
    srv := &http.Server{ Addr: ":" + cfg.ServicePort, Handler: r }
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { log.Fatalf("server error: %v", err) }
}
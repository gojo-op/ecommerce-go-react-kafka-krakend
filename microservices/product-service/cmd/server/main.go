package main

import (
    "log"
    "net/http"
    "time"
    "context"
    "github.com/gin-gonic/gin"
    "product-service/internal/config"
    "product-service/internal/controllers"
    "product-service/internal/models"
    "product-service/internal/services"
    pdb "product-service/internal/db"
    "product-service/internal/events"
)

func main() {
    cfg, err := config.Load()
    if err != nil { log.Fatalf("config error: %v", err) }

    gdb, err := pdb.Open()
    if err != nil { log.Fatalf("db error: %v", err) }

    pub := events.New()

    _ = gdb.AutoMigrate(&models.Product{})

    svc := services.New(gdb, pub)
    ctl := controllers.New(svc)

    r := gin.New()
    r.Use(gin.Logger(), gin.Recovery())

    r.GET("/health", func(c *gin.Context){ c.JSON(http.StatusOK, gin.H{"status":"healthy","service":"product-service","timestamp": time.Now().Unix()}) })
    r.HEAD("/health", func(c *gin.Context){ c.Status(http.StatusOK) })

    api := r.Group("/api/v1")
    api.GET("/products", ctl.List)
    api.GET("/products/:id", ctl.Get)
    api.GET("/products/sku/:sku", ctl.GetBySKU)
    api.POST("/products", ctl.Create)
    api.PUT("/products/:id", ctl.Update)
    api.DELETE("/products/:id", ctl.Delete)
    api.PATCH("/products/:id/stock", ctl.UpdateStock)

    srv := &http.Server{ Addr: ":" + cfg.ServicePort, Handler: r }
    var prodCount int64
    gdb.Raw("SELECT COUNT(*) FROM products").Scan(&prodCount)
    if prodCount == 0 {
        ctx := context.Background()
        _, _ = svc.Create(ctx, &models.CreateProductRequest{ Name: "Wireless Headphones", SKU: "SKU-HEAD-001", Description: "Bluetooth over-ear headphones", Price: 8999, Currency: "USD", Stock: 50, Category: "Audio", ImageURL: "https://picsum.photos/seed/headphones/600/400" })
        _, _ = svc.Create(ctx, &models.CreateProductRequest{ Name: "Gaming Mouse", SKU: "SKU-MOUS-002", Description: "Ergonomic gaming mouse", Price: 4999, Currency: "USD", Stock: 100, Category: "Accessories", ImageURL: "https://picsum.photos/seed/mouse/600/400" })
        _, _ = svc.Create(ctx, &models.CreateProductRequest{ Name: "Mechanical Keyboard", SKU: "SKU-KEYB-003", Description: "RGB mechanical keyboard", Price: 9999, Currency: "USD", Stock: 75, Category: "Accessories", ImageURL: "https://picsum.photos/seed/keyboard/600/400" })
        _, _ = svc.Create(ctx, &models.CreateProductRequest{ Name: "4K Monitor", SKU: "SKU-MONI-004", Description: "27-inch 4K UHD monitor", Price: 24999, Currency: "USD", Stock: 30, Category: "Displays", ImageURL: "https://picsum.photos/seed/monitor/600/400" })
        _, _ = svc.Create(ctx, &models.CreateProductRequest{ Name: "USB-C Dock", SKU: "SKU-DOCK-005", Description: "Multiport USB-C docking station", Price: 6999, Currency: "USD", Stock: 80, Category: "Accessories", ImageURL: "https://picsum.photos/seed/dock/600/400" })
    }
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { log.Fatalf("server error: %v", err) }
}
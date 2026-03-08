package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
    "auth-service/internal/models"

    "auth-service/internal/config"
    "auth-service/internal/controllers"
    "auth-service/internal/repositories"
    "auth-service/internal/services"
    amw "auth-service/internal/middleware"
    adb "auth-service/internal/db"
)

func main() {
	// Load configuration
	cfg, err := config.LoadAuthConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

    // Set Gin mode
    if os.Getenv("ENVIRONMENT") == "production" {
        gin.SetMode(gin.ReleaseMode)
    }

    // Initialize database (SQLite per service)
    db, err := adb.Open()
    if err != nil { log.Fatalf("Failed to initialize database: %v", err) }

    _ = db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.UserProfile{}, &models.Address{})

    var cacheClient interface{}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	permissionRepo := repositories.NewPermissionRepository(db)

	// Initialize services
    authService := services.NewAuthService(
        userRepo,
        roleRepo,
        permissionRepo,
        &CacheAdapter{},
        &NopKafka{},
        &services.AuthConfig{
            JWTSecret:        []byte(cfg.JWTSecret),
            JWTAccessExpiry:  cfg.JWTAccessExpiry,
            JWTRefreshExpiry: cfg.JWTRefreshExpiry,
            KafkaTopics: struct{ UserRegistered, UserUpdated, UserDeleted, RoleAssigned, RoleRevoked string }{
                UserRegistered: cfg.Kafka.Topics.UserRegistered,
                UserUpdated:    cfg.Kafka.Topics.UserUpdated,
                UserDeleted:    cfg.Kafka.Topics.UserDeleted,
                RoleAssigned:   cfg.Kafka.Topics.RoleAssigned,
                RoleRevoked:    cfg.Kafka.Topics.RoleRevoked,
            },
        },
    )

    var userCount int64
    db.Raw("SELECT COUNT(*) FROM users").Scan(&userCount)
    if userCount == 0 {
        // Seed users via service to ensure proper password hashing
        _, _ = authService.Register(context.Background(), &models.RegisterRequest{
            Email:     "admin@example.com",
            Username:  "admin",
            Password:  "password123",
            FirstName: "Admin",
            LastName:  "User",
        })
        _, _ = authService.Register(context.Background(), &models.RegisterRequest{
            Email:     "customer@example.com",
            Username:  "customer",
            Password:  "password123",
            FirstName: "John",
            LastName:  "Doe",
        })
        var adminID string
        db.Raw("SELECT id FROM users WHERE email = ?", "admin@example.com").Scan(&adminID)
        if adminID != "" {
            _ = authService.AssignRole(context.Background(), adminID, "admin")
        }
    }

    // Initialize controllers
    addrRepo := repositories.NewAddressRepository(db)
    authController := controllers.NewAuthController(authService, addrRepo)

    // JWT middleware
    authMw := amw.NewAuthMiddleware(string(cfg.JWTSecret))

	// Create Gin router
	router := gin.New()

	// Apply global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

    // Health check endpoint
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{ "status": "healthy", "service": cfg.ServiceName, "timestamp": time.Now().Unix() })
    })
    router.HEAD("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	// API routes
	api := router.Group("/api/v1")
	{
		// Public auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/refresh", authController.RefreshToken)
		}

    // Protected auth routes
    protected := api.Group("/auth")
    protected.Use(authMw.RequireAuth())
    {
            protected.POST("/logout", authController.Logout)
            protected.GET("/profile", authController.GetProfile)
            protected.PUT("/profile", authController.UpdateProfile)
            protected.POST("/change-password", authController.ChangePassword)
            protected.GET("/addresses", authController.ListAddresses)
            protected.POST("/addresses", authController.CreateAddress)
            protected.PUT("/addresses/:id", authController.UpdateAddress)
            protected.DELETE("/addresses/:id", authController.DeleteAddress)
    }

    // Admin routes
    admin := api.Group("/auth")
    admin.Use(authMw.RequireRole("admin"))
    {
            admin.POST("/users/:user_id/roles", authController.AssignRole)
            admin.DELETE("/users/:user_id/roles", authController.RevokeRole)
    }
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.ServicePort,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting %s on port %s", cfg.ServiceName, cfg.ServicePort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

    // Cleanup

    _ = cacheClient

    // SQLite gorm has no Close; noop

	log.Println("Server exited")
}
type CacheAdapter struct{}
func (a *CacheAdapter) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error { return nil }
func (a *CacheAdapter) Get(ctx context.Context, key string, dest interface{}) error { return fmt.Errorf("no cache") }
func (a *CacheAdapter) Delete(ctx context.Context, key string) error { return nil }

type NopKafka struct{}
func (n *NopKafka) Publish(ctx context.Context, topic string, event services.Event) error { return nil }
package main

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/product"
	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/infrastructure/config"
	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/infrastructure/delivery"
	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/infrastructure/persistence"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Database connection established")

	// Initialize repositories
	productRepo := persistence.NewProductRepository(db)

	// Initialize use cases
	createProductUseCase := product.NewCreateProductUseCase(productRepo)
	getProductUseCase := product.NewGetProductUseCase(productRepo)

	// Initialize handlers
	productHandler := delivery.NewProductHandler(createProductUseCase, getProductUseCase)

	// Set Gin mode based on environment
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.New()

	// Apply global middleware
	router.Use(gin.Recovery())
	router.Use(delivery.LoggerMiddleware())
	router.Use(delivery.ErrorHandlerMiddleware())
	router.Use(delivery.CORSMiddleware())

	// Register routes
	registerRoutes(router, productHandler)

	// Start server in a goroutine
	serverAddr := cfg.GetServerAddress()
	go func() {
		log.Printf("Starting server on %s", serverAddr)
		if err := router.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}

// initDatabase initializes and returns a database connection
func initDatabase(cfg *config.Config) (*sql.DB, error) {
	dsn := cfg.GetDatabaseDSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}

// registerRoutes registers all API routes
func registerRoutes(router *gin.Engine, productHandler *delivery.ProductHandler) {
	// Health check endpoint
	router.GET("/health", delivery.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Product routes
		products := v1.Group("/products")
		{
			products.POST("", productHandler.Create)
			products.GET("/:id", productHandler.Get)
		}
	}
}

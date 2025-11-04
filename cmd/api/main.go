package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/inventory/command"
	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/inventory/query"
	productcommand "github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/product/command"
	productquery "github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/product/query"
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

	// Initialize repositories (CQRS: separate command and query repositories)
	productCmdRepo := persistence.NewProductCommandRepository(db)
	productQueryRepo := persistence.NewProductQueryRepository(db)
	inventoryCmdRepo := persistence.NewInventoryCommandRepository(db)
	inventoryQueryRepo := persistence.NewInventoryQueryRepository(db)

	// STEP 1: Initialize product queries (without inventory integration first)
	getProductQueryBasic := productquery.NewGetProductQuery(productQueryRepo)

	// STEP 2: Create adapter for Inventory → Product communication
	productQueryAdapter := query.NewProductQueryAdapter(getProductQueryBasic)

	// STEP 3: Initialize inventory commands and queries with product query adapter injection
	// This demonstrates Inventory → Product module communication
	createInventoryCommand := command.NewCreateInventoryCommand(
		inventoryCmdRepo,
		inventoryQueryRepo,
		productQueryAdapter,
	)
	getInventoryQuery := query.NewGetInventoryQuery(
		inventoryQueryRepo,
		productQueryAdapter,
	)
	adjustInventoryCommand := command.NewAdjustInventoryCommand(
		inventoryCmdRepo,
		inventoryQueryRepo,
		productQueryAdapter,
	)

	// STEP 4: Create adapter for Product → Inventory communication
	// Wrap GetInventoryQuery.Execute to match the function signature expected by ProductInventoryAdapter
	inventoryAdapterFunc := func(ctx context.Context, productID string) (*productquery.InventoryOutput, error) {
		output, err := getInventoryQuery.Execute(ctx, productID)
		if err != nil {
			return nil, err
		}
		return &productquery.InventoryOutput{
			Quantity:          output.Quantity,
			AvailableQuantity: output.AvailableQuantity,
		}, nil
	}
	inventoryAdapter := productquery.NewProductInventoryAdapter(inventoryAdapterFunc)

	// STEP 5: Re-initialize product query WITH inventory integration
	// This demonstrates Product → Inventory bidirectional module communication
	getProductQuery := productquery.NewGetProductQueryWithInventory(productQueryRepo, inventoryAdapter)

	// Initialize product command
	createProductCommand := productcommand.NewCreateProductCommand(productCmdRepo)

	// Initialize handlers
	productHandler := delivery.NewProductHandler(createProductCommand, getProductQuery)
	inventoryHandler := delivery.NewInventoryHandler(createInventoryCommand, getInventoryQuery, adjustInventoryCommand)

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
	registerRoutes(router, productHandler, inventoryHandler)

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
	// db.SetConnMaxLifetime(5 * time.Minute) // Uncomment if needed

	return db, nil
}

// registerRoutes registers all API routes
func registerRoutes(router *gin.Engine, productHandler *delivery.ProductHandler, inventoryHandler *delivery.InventoryHandler) {
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

		// Inventory routes
		inventoryGroup := v1.Group("/inventory")
		{
			inventoryGroup.POST("", inventoryHandler.Create)
			inventoryGroup.GET("/:productId", inventoryHandler.Get)
			inventoryGroup.PATCH("/adjust", inventoryHandler.Adjust)
		}
	}
}

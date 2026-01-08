package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vhvplatform/go-shared/config"
	"github.com/vhvplatform/go-shared/logger"
	"github.com/vhvplatform/go-shared/mongodb"
	"github.com/vhvplatform/go-tenant-service/internal/grpc"
	"github.com/vhvplatform/go-tenant-service/internal/handler"
	"github.com/vhvplatform/go-tenant-service/internal/repository"
	"github.com/vhvplatform/go-tenant-service/internal/service"

	pb "github.com/vhvplatform/go-tenant-service/proto"
	"go.uber.org/zap"
	grpcServer "google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Initialize logger
	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer log.Sync()

	log.Info("Starting Tenant Service", zap.String("environment", cfg.Environment))

	// Initialize MongoDB
	mongoClient, err := mongodb.NewClient(context.Background(), mongodb.Config{
		URI:         cfg.MongoDB.URI,
		Database:    cfg.MongoDB.Database,
		MaxPoolSize: cfg.MongoDB.MaxPoolSize,
		MinPoolSize: cfg.MongoDB.MinPoolSize,
	})
	if err != nil {
		log.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer mongoClient.Close(context.Background())

	// Initialize repositories
	tenantRepo := repository.NewTenantRepository(mongoClient.Database())
	tenantUserRepo := repository.NewTenantUserRepository(mongoClient.Database())

	// Initialize services
	tenantService := service.NewTenantService(tenantRepo, tenantUserRepo, log)

	// Start gRPC server
	grpcPort := os.Getenv("TENANT_SERVICE_PORT")
	if grpcPort == "" {
		grpcPort = "50053"
	}
	go startGRPCServer(tenantService, log, grpcPort)

	// Start HTTP server
	httpPort := os.Getenv("TENANT_SERVICE_HTTP_PORT")
	if httpPort == "" {
		httpPort = "8083"
	}
	startHTTPServer(tenantService, log, httpPort)
}

func startGRPCServer(tenantService *service.TenantService, log *logger.Logger, port string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal("Failed to listen", zap.Error(err))
	}

	grpcSrv := grpcServer.NewServer()
	tenantGrpcServer := grpc.NewTenantServiceServer(tenantService, log)
	pb.RegisterTenantServiceServer(grpcSrv, tenantGrpcServer)

	// Register health check service
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcSrv, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	log.Info("gRPC server listening", zap.String("port", port))
	if err := grpcSrv.Serve(lis); err != nil {
		log.Fatal("Failed to serve gRPC", zap.Error(err))
	}
}

func startHTTPServer(tenantService *service.TenantService, log *logger.Logger, port string) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Initialize handlers
	tenantHandler := handler.NewTenantHandler(tenantService, log)

	// Health check endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	router.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	// API routes
	v1 := router.Group("/api/v1")
	{
		tenants := v1.Group("/tenants")
		{
			tenants.POST("", tenantHandler.CreateTenant)
			tenants.GET("", tenantHandler.ListTenants)
			tenants.GET("/:id", tenantHandler.GetTenant)
			tenants.PUT("/:id", tenantHandler.UpdateTenant)
			tenants.DELETE("/:id", tenantHandler.DeleteTenant)
			tenants.POST("/:id/users", tenantHandler.AddUserToTenant)
			tenants.DELETE("/:id/users/:user_id", tenantHandler.RemoveUserFromTenant)
		}
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		log.Info("HTTP server listening", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exited")
}

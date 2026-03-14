package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"millog_backend/internal/database"
	"millog_backend/internal/handlers"
	"millog_backend/internal/middleware"
	"millog_backend/internal/models"
	"millog_backend/internal/repositories"
	"millog_backend/internal/services"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, relying on system environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		slog.Error("DATABASE_URL is not set")
		os.Exit(1)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-in-production"
		slog.Warn("JWT_SECRET not set, using default")
	}

	ctx := context.Background()

	var dbPool *pgxpool.Pool
	for attempt := 1; attempt <= 10; attempt++ {
		connCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		pool, err := services.NewPostgresDB(connCtx, dbURL)
		cancel()
		if err == nil {
			dbPool = pool
			break
		}
		slog.Warn("Database connection failed, retrying...", "attempt", attempt, "error", err)
		if attempt < 10 {
			time.Sleep(2 * time.Second)
		} else {
			slog.Error("Failed to connect to database after 10 attempts", "error", err)
			os.Exit(1)
		}
	}
	defer dbPool.Close()

	migrateCtx, cancelMigrate := context.WithTimeout(ctx, 30*time.Second)
	defer cancelMigrate()
	if err := database.Migrate(migrateCtx, dbPool); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}
	userRepo := repositories.NewUserRepository()
	tokenRepo := repositories.NewInviteTokenRepository()
	refreshTokenRepo := repositories.NewRefreshTokenRepository()

	emailService, err := services.NewEmailService(
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_EMAIL"),
		os.Getenv("SMTP_PASSWORD"),
	)
	if err != nil {
		slog.Error("Failed to init email service", "error", err)
		os.Exit(1)
	}

	authService := services.NewAuthService(userRepo, tokenRepo, refreshTokenRepo, dbPool, emailService, jwtSecret)
	authHandler := handlers.NewAuthHandler(authService)

	catRepo := repositories.NewCategoryRepository()
	resRepo := repositories.NewResourceRepository()
	reqRepo := repositories.NewSupplyRequestRepository()
	unitRepo := repositories.NewUnitRepository()
	volReqRepo := repositories.NewVolunteerRequestRepository()
	invService := services.NewInventoryService(catRepo, resRepo, dbPool)
	reqService := services.NewRequestService(reqRepo, resRepo, dbPool)
	unitService := services.NewUnitService(unitRepo, dbPool)
	volReqService := services.NewVolunteerRequestService(volReqRepo, dbPool)
	invHandler := handlers.NewInventoryHandler(invService)
	reqHandler := handlers.NewRequestHandler(reqService)
	unitHandler := handlers.NewUnitHandler(unitService)
	volReqHandler := handlers.NewVolunteerRequestHandler(volReqService)

	r := gin.Default()

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/register", authHandler.RegisterVolunteer)
			auth.POST("/setup-password", authHandler.SetupPassword)
			auth.GET("/me", middleware.AuthMiddleware(jwtSecret), authHandler.Me)
		}

		api.POST("/bootstrap", authHandler.BootstrapAdmin)

		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(jwtSecret), middleware.RequireAnyRole(models.UserCreatorRoles))
		{
			admin.POST("/users", authHandler.RegisterUser)
		}

		// Units: Admin + commanders + logists + storekeepers
		units := api.Group("/units")
		units.Use(middleware.AuthMiddleware(jwtSecret))
		{
			units.GET("", unitHandler.List)
			units.POST("", middleware.RequireAnyRole(models.UnitManagerRoles), unitHandler.Create)
		}

		// Inventory: storekeepers + company sergeant
		inv := api.Group("/inventory")
		inv.Use(middleware.AuthMiddleware(jwtSecret))
		{
			inv.GET("/categories", invHandler.ListCategories)
			inv.GET("/resources", invHandler.ListResources)
			inv.POST("/categories", middleware.RequireAnyRole(models.InventoryManagerRoles), invHandler.CreateCategory)
			inv.POST("/resources", middleware.RequireAnyRole(models.InventoryManagerRoles), invHandler.CreateResource)
		}

		// Supply requests: commanders + logists + sergeant create; commanders + logists approve
		requests := api.Group("/requests")
		requests.Use(middleware.AuthMiddleware(jwtSecret))
		{
			requests.POST("", middleware.RequireAnyRole(models.SupplyRequestCreatorRoles), reqHandler.Create)
			requests.GET("", reqHandler.List)
			requests.POST("/:id/approve", middleware.RequireAnyRole(models.SupplyRequestApproverRoles), reqHandler.Approve)
		}

		// Volunteer requests: military creates, volunteers take and complete
		volRequests := api.Group("/volunteer-requests")
		volRequests.Use(middleware.AuthMiddleware(jwtSecret))
		{
			volRequests.GET("", volReqHandler.List)
			volRequests.POST("", middleware.RequireAnyRole(models.VolunteerRequestCreatorRoles), volReqHandler.Create)
			volRequests.POST("/:id/take", middleware.RequireRoles(models.RoleVolunteer), volReqHandler.Take)
			volRequests.POST("/:id/complete", middleware.RequireRoles(models.RoleVolunteer), volReqHandler.Complete)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		slog.Info("Server is starting", "port", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exited properly")
}

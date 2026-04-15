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

	catRepo := repositories.NewCategoryRepository()
	resRepo := repositories.NewResourceRepository()
	reqRepo := repositories.NewSupplyRequestRepository()
	unitRepo := repositories.NewUnitRepository()
	volReqRepo := repositories.NewVolunteerRequestRepository()
	vehicleRepo := repositories.NewVehicleRepository()
	fuelRepo := repositories.NewFuelRepository()
	warehouseRepo := repositories.NewWarehouseRepository()
	analyticsRepo := repositories.NewAnalyticsRepository()

	invService := services.NewInventoryService(catRepo, resRepo, userRepo, dbPool)
	reqService := services.NewRequestService(reqRepo, resRepo, userRepo, dbPool)
	unitService := services.NewUnitService(unitRepo, userRepo, dbPool)
	volReqService := services.NewVolunteerRequestService(volReqRepo, dbPool)
	warehouseService := services.NewWarehouseService(warehouseRepo, dbPool)
	analyticsService := services.NewAnalyticsService(analyticsRepo, dbPool) // dbPool - твоє з'єднання

	invHandler := handlers.NewInventoryHandler(invService)
	reqHandler := handlers.NewRequestHandler(reqService)
	unitHandler := handlers.NewUnitHandler(unitService)
	volReqHandler := handlers.NewVolunteerRequestHandler(volReqService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	warehouseHandler := handlers.NewWarehouseHandler(warehouseService)

	authService := services.NewAuthService(userRepo, unitRepo, tokenRepo, refreshTokenRepo, dbPool, emailService, jwtSecret)
	authHandler := handlers.NewAuthHandler(authService)
	fuelService := services.NewFuelService(fuelRepo, dbPool)
	fuelHandler := handlers.NewFuelHandler(fuelService)

	vehicleService := services.NewVehicleService(vehicleRepo, dbPool)
	vehicleHandler := handlers.NewVehicleHandler(vehicleService)

	r := gin.Default()

	os.MkdirAll("uploads/maintenance", os.ModePerm)

	r.Static("/uploads", "./uploads")

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/register", authHandler.RegisterVolunteer)
			auth.POST("/setup-password", authHandler.SetupPassword)
			auth.POST("/forgot-password", authHandler.RequestPasswordReset)
			auth.GET("/me", middleware.AuthMiddleware(jwtSecret, dbPool), authHandler.Me)
		}

		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(jwtSecret, dbPool))
		{
			users.GET("/commanders", authHandler.ListCommanders)
			users.GET("/visible", authHandler.GetVisibleUsers)
			users.PUT("/:id/role", authHandler.UpdateRoleAndUnit)
			users.PUT("/:id/block", authHandler.BlockUser)
			users.PUT("/:id/unblock", authHandler.UnblockUser)
			users.PATCH("/profile", authHandler.UpdateMyProfile)
			users.PATCH("/password", authHandler.UpdateMyPassword)
		}

		api.POST("/bootstrap", authHandler.BootstrapAdmin)

		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(jwtSecret, dbPool), middleware.RequireAnyRole(models.UserCreatorRoles))
		{
			admin.POST("/users", authHandler.RegisterUser)
		}

		// Units: Admin + commanders + logists + storekeepers
		units := api.Group("/units")
		units.Use(middleware.AuthMiddleware(jwtSecret, dbPool))
		{
			units.GET("", unitHandler.List)
			units.GET("/available", unitHandler.GetAvailableForRole)
			units.POST("", middleware.RequireAnyRole(models.UnitManagerRoles), unitHandler.Create)
			units.POST("/:id/change-commander", middleware.RequireAnyRole(models.UnitManagerRoles), unitHandler.ChangeCommander)
			units.GET("/my-hierarchy", unitHandler.GetMyHierarchyForRole)
		}

		// Inventory: storekeepers + company sergeant
		inv := api.Group("/inventory")
		inv.Use(middleware.AuthMiddleware(jwtSecret, dbPool))
		{
			inv.GET("/categories", invHandler.ListCategories)
			inv.GET("/resources", invHandler.ListResources)
			inv.POST("/categories", middleware.RequireAnyRole(models.InventoryManagerRoles), invHandler.CreateCategory)
			inv.POST("/resources", middleware.RequireAnyRole(models.InventoryManagerRoles), invHandler.CreateResource)
			inv.POST("/resources/:id/write-off", middleware.RequireAnyRole(models.InventoryManagerRoles), invHandler.WriteOff)
			inv.PATCH("/resources/:id", middleware.RequireAnyRole(models.InventoryManagerRoles), invHandler.UpdateResource)
			inv.DELETE("/resources/:id", middleware.RequireAnyRole(models.InventoryManagerRoles), invHandler.Delete)
			inv.POST("/resources/:id/assign", middleware.RequireAnyRole(models.InventoryManagerRoles), invHandler.Assign)
			inv.GET("/my-equipment", middleware.AuthMiddleware(jwtSecret, dbPool), invHandler.GetMyEquipment)
			inv.POST("/issue", middleware.AuthMiddleware(jwtSecret, dbPool), invHandler.IssueResource)
			inv.POST("/shipments", middleware.AuthMiddleware(jwtSecret, dbPool), invHandler.CreateShipment)
			inv.GET("/warehouse/:id", middleware.AuthMiddleware(jwtSecret, dbPool), invHandler.GetByWarehouse)
			inv.POST("/shipments/:id/receive", middleware.AuthMiddleware(jwtSecret, dbPool), invHandler.ReceiveShipment)
			inv.GET("/shipments", middleware.AuthMiddleware(jwtSecret, dbPool), invHandler.ListShipments)
			inv.GET("/shipments/:id/pdf", middleware.AuthMiddleware(jwtSecret, dbPool), invHandler.DownloadShipmentPDF)
			inv.GET("/resources/:id/qr", invHandler.DownloadResourceQR)
		}

		// Supply requests: commanders + logists + sergeant create; commanders + logists approve
		requests := api.Group("/requests")
		requests.Use(middleware.AuthMiddleware(jwtSecret, dbPool))
		{
			requests.POST("", middleware.RequireAnyRole(models.SupplyRequestCreatorRoles), reqHandler.Create)
			requests.GET("", reqHandler.List)
			requests.POST("/:id/approve", middleware.RequireAnyRole(models.SupplyRequestApproverRoles), reqHandler.Approve)
		}

		// Volunteer requests: military creates, volunteers take and complete
		volRequests := api.Group("/volunteer-requests")
		volRequests.Use(middleware.AuthMiddleware(jwtSecret, dbPool))
		{
			// Перегляд списку (доступно всім)
			volRequests.GET("", volReqHandler.List)

			// Створення заявки (тільки для військових)
			volRequests.POST("", middleware.RequireAnyRole(models.MilitaryInventoryRoles), volReqHandler.Create)

			// Дії ВОЛОНТЕРА (Взяти в роботу, Доставити)
			volRequests.POST("/:id/take", middleware.RequireAnyRole([]models.UserRole{models.RoleVolunteer}), volReqHandler.Take)
			volRequests.POST("/:id/deliver", middleware.RequireAnyRole([]models.UserRole{models.RoleVolunteer}), volReqHandler.Deliver)

			// Дії ВІЙСЬКОВИХ (Прийняти на баланс, Відхилити, Скасувати)
			volRequests.POST("/:id/accept", middleware.RequireAnyRole(models.MilitaryInventoryRoles), volReqHandler.Accept)
			volRequests.POST("/:id/reject", middleware.RequireAnyRole(models.MilitaryInventoryRoles), volReqHandler.Reject)
			volRequests.POST("/:id/cancel", middleware.RequireAnyRole(models.MilitaryInventoryRoles), volReqHandler.Cancel)
		}

		// Fuel records: logists + commanders create; logists + commanders view
		vehicleGroup := api.Group("/vehicles")
		vehicleGroup.Use(middleware.AuthMiddleware(jwtSecret, dbPool))
		vehicleGroup.Use(middleware.RequireAnyRole(models.FuelRecordCreatorRoles))
		{
			vehicleGroup.POST("", vehicleHandler.Create)
			vehicleGroup.GET("", vehicleHandler.GetAll)
			vehicleGroup.GET("/:id", vehicleHandler.GetByID)
			vehicleGroup.PATCH("/:id/status", vehicleHandler.UpdateStatus)

			vehicleGroup.POST("/:id/fuel", fuelHandler.CreateRecord)
			vehicleGroup.GET("/:id/fuel", fuelHandler.GetHistory)
			vehicleGroup.POST("/:id/maintenance", vehicleHandler.PerformMaintenance)
			vehicleGroup.GET("/:id/maintenance", vehicleHandler.GetMaintenanceHistory)
			vehicleGroup.PATCH("/:id/driver", vehicleHandler.AssignDriver)
			vehicleGroup.GET("/:id/drivers", vehicleHandler.GetDriverHistory)
		}

		warehouseGroup := api.Group("/warehouses")
		warehouseGroup.Use(middleware.AuthMiddleware(jwtSecret, dbPool), middleware.RequireAnyRole(models.WarehouseManagerRoles))
		{
			warehouseGroup.GET("", warehouseHandler.List)
			warehouseGroup.POST("", warehouseHandler.Create)
			warehouseGroup.PATCH("/:id/location", warehouseHandler.UpdateLocation)
		}

		analyticsGroup := api.Group("/analytics")
		analyticsGroup.Use(middleware.AuthMiddleware(jwtSecret, dbPool))
		{
			analyticsGroup.GET("/dashboard", analyticsHandler.GetDashboard)
			analyticsGroup.POST("/auto-replenish", analyticsHandler.AutoReplenish)
			analyticsGroup.GET("/export/inventory", analyticsHandler.ExportInventory)
			analyticsGroup.GET("/export/fuel", analyticsHandler.ExportFuel)
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

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

	"Omnilog_backend/internal/database"
	"Omnilog_backend/internal/handlers"
	"Omnilog_backend/internal/middleware"
	"Omnilog_backend/internal/models"
	"Omnilog_backend/internal/repositories"
	"Omnilog_backend/internal/services"
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
	volReqRepo := repositories.NewContractorRequestRepository()
	vehicleRepo := repositories.NewVehicleRepository()
	fuelRepo := repositories.NewFuelRepository()
	warehouseRepo := repositories.NewWarehouseRepository()
	analyticsRepo := repositories.NewAnalyticsRepository()
	auditRepo := repositories.NewAuditLogRepository()
	tenantRepo := repositories.NewTenantRepository()
	notificationRepo := repositories.NewNotificationRepository()

	// Спочатку створюємо базові сервіси
	notificationService := services.NewNotificationService(notificationRepo, dbPool)
	auditService := services.NewAuditService(auditRepo, dbPool)
	unitService := services.NewUnitService(unitRepo, userRepo, dbPool)
	warehouseService := services.NewWarehouseService(warehouseRepo, dbPool)
	analyticsService := services.NewAnalyticsService(analyticsRepo, dbPool)
	limitationService := services.NewLimitationService(dbPool)
	slaMonitor := services.NewSLAMonitor(dbPool, reqRepo, auditRepo, emailService)

	// Тепер створюємо сервіси, які залежать від notificationService
	invService := services.NewInventoryService(catRepo, resRepo, userRepo, dbPool, notificationService)
	reqService := services.NewRequestService(reqRepo, resRepo, userRepo, dbPool, notificationService)
	volReqService := services.NewContractorRequestService(volReqRepo, dbPool)

	authService := services.NewAuthService(userRepo, unitRepo, tokenRepo, refreshTokenRepo, dbPool, emailService, jwtSecret)
	invHandler := handlers.NewInventoryHandler(invService, auditService, limitationService, authService)
	reqHandler := handlers.NewRequestHandler(reqService, auditService, slaMonitor)
	unitHandler := handlers.NewUnitHandler(unitService, auditService)
	volReqHandler := handlers.NewContractorRequestHandler(volReqService, auditService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService, auditService)
	warehouseHandler := handlers.NewWarehouseHandler(warehouseService, auditService, limitationService, authService)

	authHandler := handlers.NewAuthHandler(authService, auditService)
	fuelService := services.NewFuelService(fuelRepo, dbPool)
	fuelHandler := handlers.NewFuelHandler(fuelService, auditService)

	vehicleService := services.NewVehicleService(vehicleRepo, dbPool)
	vehicleHandler := handlers.NewVehicleHandler(vehicleService, auditService)
	auditHandler := handlers.NewAuditHandler(auditService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	slaMonitor.Start(context.Background())

	r := gin.Default()

	// === SECURITY MIDDLEWARE ===
	allowedOrigins := []string{
		os.Getenv("FRONTEND_URL"),
		"http://localhost:5173",
		"http://localhost:3000",
		"https://localhost",
	}
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORSMiddleware(allowedOrigins))

	// Rate limiter для auth endpoints: 10 спроб за 5 хвилин з однієї IP
	authRateLimiter := middleware.NewRateLimiter(10, 5*time.Minute)

	os.MkdirAll("uploads/maintenance", os.ModePerm)

	// /uploads доступний лише через автентифікований proxy-endpoint
	// НЕ виставляємо r.Static("/uploads", "./uploads") публічно

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", middleware.RateLimitMiddleware(authRateLimiter), authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/register", middleware.RateLimitMiddleware(authRateLimiter), authHandler.RegisterContractor)
			auth.POST("/setup-password", authHandler.SetupPassword)
			auth.POST("/forgot-password", middleware.RateLimitMiddleware(authRateLimiter), authHandler.RequestPasswordReset)
			auth.POST("/tenants/signup", middleware.RateLimitMiddleware(authRateLimiter), authHandler.SignupTenant)
			auth.GET("/me", middleware.AuthMiddleware(jwtSecret, dbPool), authHandler.Me)
		}

		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(jwtSecret, dbPool))
		{
			users.GET("/commanders", authHandler.ListCommanders)
			users.GET("/visible", authHandler.GetVisibleUsers)
			users.GET("/limits", authHandler.GetUserLimits)
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
			admin.GET("/audit-logs", auditHandler.GetLogs)
			admin.POST("/sla/trigger", reqHandler.TriggerCheck)
		}

		// === PLATFORM ADMIN API (SYSTEM_ADMIN only, cross-tenant) ===
		platformHandler := handlers.NewPlatformHandler(tenantRepo, userRepo, dbPool)
		platform := api.Group("/platform")
		platform.Use(middleware.AuthMiddleware(jwtSecret, dbPool), middleware.RequireSystemAdmin())
		{
			platform.GET("/stats", platformHandler.Stats)
			platform.GET("/tenants", platformHandler.ListTenants)
			platform.GET("/tenants/:id", platformHandler.GetTenant)
			platform.PATCH("/tenants/:id/tier", platformHandler.UpdateTier)
			platform.PATCH("/tenants/:id/active", platformHandler.SetActive)
			platform.DELETE("/tenants/:id", platformHandler.DeleteTenant)
		}

		// Units: Admin + commanders + logists + storekeepers
		units := api.Group("/units")
		units.Use(middleware.AuthMiddleware(jwtSecret, dbPool))
		{
			units.GET("", unitHandler.List)
			units.GET("/available", unitHandler.GetAvailableForRole)
			units.POST("", middleware.RequireAnyRole(models.UnitManagerRoles), unitHandler.Create)
			units.POST("/:id/change-commander", middleware.RequireAnyRole(models.UnitManagerRoles), unitHandler.ChangeManager)
			units.GET("/my-hierarchy", unitHandler.GetMyHierarchyForRole)
			//units.GET("/:id", unitHandler.GetByID) // Якщо в тебе ще немає отримання одного підрозділу
			units.PATCH("/:id", middleware.RequireAnyRole(models.UnitManagerRoles), unitHandler.UpdateUnit)
			units.DELETE("/:id", middleware.RequireAnyRole(models.UnitManagerRoles), unitHandler.DeleteUnit)
		}

		// Inventory: storekeepers + company sergeant
		inv := api.Group("/inventory")
		inv.Use(middleware.AuthMiddleware(jwtSecret, dbPool))
		{
			inv.GET("/categories", invHandler.ListCategories)
			inv.GET("/resources", invHandler.ListResources)
			inv.GET("/resources/unique-names", invHandler.GetUniqueResourceNames) // Унікальні назви для форми заявки
			inv.GET("/resources/:id", invHandler.GetResource)
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
			inv.POST("/shipments/:id/start", middleware.AuthMiddleware(jwtSecret, dbPool), invHandler.StartShipment)
			inv.POST("/shipments/:id/receive", middleware.AuthMiddleware(jwtSecret, dbPool), invHandler.ReceiveShipment)
			inv.GET("/shipments", middleware.AuthMiddleware(jwtSecret, dbPool), invHandler.ListShipments)
			inv.GET("/shipments/my", middleware.AuthMiddleware(jwtSecret, dbPool), invHandler.ListMyShipments)
			inv.GET("/shipments/:id/pdf", middleware.AuthMiddleware(jwtSecret, dbPool), invHandler.DownloadShipmentPDF)
			inv.GET("/resources/:id/qr", invHandler.DownloadResourceQR)
			inv.PATCH("/categories/:id", middleware.RequireAnyRole(models.InventoryManagerRoles), invHandler.UpdateCategory)
			inv.DELETE("/categories/:id", middleware.RequireAnyRole(models.InventoryManagerRoles), invHandler.DeleteCategory)
			inv.POST("/audit", invHandler.SubmitAudit)
			inv.GET("/resources/import/template", invHandler.DownloadImportTemplate) // Завантаження шаблону
			// 🚀 PRO FEATURE: Excel import з захистом
			inv.POST("/resources/import", middleware.RequireAnyRole(models.InventoryManagerRoles), middleware.RequireSubscriptionTier("PRO", dbPool), invHandler.ImportExcel) // Завантаження заповненого Excel
		}

		// Supply requests: commanders + logists + sergeant create; commanders + logists approve
		requests := api.Group("/requests")
		requests.Use(middleware.AuthMiddleware(jwtSecret, dbPool))
		{
			requests.POST("", middleware.RequireAnyRole(models.SupplyRequestCreatorRoles), reqHandler.Create)
			requests.GET("", reqHandler.List)
			requests.POST("/:id/approve", middleware.RequireAnyRole(models.SupplyRequestApproverRoles), reqHandler.Approve)
			requests.POST("/:id/reject", middleware.RequireAnyRole(models.SupplyRequestApproverRoles), reqHandler.Reject) // 👈 НОВЕ (відмова логіста)
			requests.POST("/:id/cancel", reqHandler.Cancel)
			requests.GET("/:id", reqHandler.GetByID)
			// 🚀 PRO FEATURE: Smart Dispatch з захистом
			requests.POST("/smart-dispatch-preview", middleware.RequireSubscriptionTier("PRO", dbPool), reqHandler.SmartDispatchPreview)
			requests.POST("/smart-dispatch-confirm", middleware.RequireSubscriptionTier("PRO", dbPool), reqHandler.SmartDispatchConfirm)

		}

		// CONTRACTOR requests: military creates, CONTRACTORs take and complete
		contractorReqs := api.Group("/contractor-requests")
		contractorReqs.Use(middleware.AuthMiddleware(jwtSecret, dbPool))
		{
			// Перегляд списку (доступно всім)
			contractorReqs.GET("", volReqHandler.List)

			// Створення заявки (тільки для військових)
			contractorReqs.POST("", middleware.RequireAnyRole(models.ContractorRequestCreatorRoles), volReqHandler.Create)

			// Дії ВОЛОНТЕРА (Взяти в роботу, Доставити)
			contractorReqs.POST("/:id/take", middleware.RequireAnyRole([]models.UserRole{models.RoleContractor}), volReqHandler.Take)
			contractorReqs.POST("/:id/deliver", middleware.RequireAnyRole([]models.UserRole{models.RoleContractor}), volReqHandler.Deliver)

			// Дії ВІЙСЬКОВИХ (Прийняти на баланс, Відхилити, Скасувати)
			contractorReqs.POST("/:id/accept", middleware.RequireAnyRole(models.ContractorRequestCreatorRoles), volReqHandler.Accept)
			contractorReqs.POST("/:id/reject", middleware.RequireAnyRole(models.ContractorRequestCreatorRoles), volReqHandler.Reject)
			contractorReqs.POST("/:id/cancel", middleware.RequireAnyRole(models.ContractorRequestCreatorRoles), volReqHandler.Cancel)
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

			vehicleGroup.PATCH("/:id", vehicleHandler.Update)
			vehicleGroup.DELETE("/:id", vehicleHandler.Delete)
			vehicleGroup.GET("/available-for-route", vehicleHandler.GetAvailableForShipment)
		}

		// Захищений endpoint для скачування файлів ТО (лише для авторизованих)
		r.GET("/uploads/maintenance/:filename",
			middleware.AuthMiddleware(jwtSecret, dbPool),
			func(c *gin.Context) {
				filename := c.Param("filename")
				// Захист від path traversal — лише ім'я файлу без слешів
				for _, ch := range filename {
					if ch == '/' || ch == '\\' || ch == '.' && filename[0] == '.' {
						c.AbortWithStatus(http.StatusBadRequest)
						return
					}
				}
				c.File("./uploads/maintenance/" + filename)
			},
		)

		warehouseGroup := api.Group("/warehouses")
		warehouseGroup.Use(middleware.AuthMiddleware(jwtSecret, dbPool), middleware.RequireAnyRole(models.WarehouseManagerRoles))
		{
			warehouseGroup.GET("", warehouseHandler.List)
			warehouseGroup.POST("", warehouseHandler.Create)
			warehouseGroup.PATCH("/:id/location", warehouseHandler.UpdateLocation)
			warehouseGroup.PATCH("/:id", warehouseHandler.UpdateWarehouse)
			warehouseGroup.DELETE("/:id", warehouseHandler.Delete)
		}

		analyticsGroup := api.Group("/analytics")
		analyticsGroup.Use(middleware.AuthMiddleware(jwtSecret, dbPool))
		{
			// Базовий дашборд (SLA/TCO/ризики) доступний на всіх тарифах —
			// це перегляд власних даних, а не платна аналітика.
			analyticsGroup.GET("/dashboard", analyticsHandler.GetDashboard)
			// 🚀 PRO FEATURE: Smart-поповнення (дія, а не перегляд)
			analyticsGroup.POST("/auto-replenish", middleware.RequireSubscriptionTier("PRO", dbPool), analyticsHandler.AutoReplenish)
			// 🚀 НОВІ PRO FEATURES: KPI та Forecast
			analyticsGroup.GET("/kpi", middleware.RequireSubscriptionTier("PRO", dbPool), analyticsHandler.GetAdvancedKPIs)
			analyticsGroup.GET("/forecast", middleware.RequireSubscriptionTier("PRO", dbPool), analyticsHandler.GetDemandForecast)
			// 🚀 ДІЇ PRO FEATURES: Maintenance та Fuel Anomalies
			analyticsGroup.GET("/maintenance", middleware.RequireSubscriptionTier("PRO", dbPool), analyticsHandler.GetPredictiveMaintenanceSchedule)
			analyticsGroup.GET("/fuel-anomalies", middleware.RequireSubscriptionTier("PRO", dbPool), analyticsHandler.GetFuelAnomalyDetection)
			analyticsGroup.GET("/export/inventory", analyticsHandler.ExportInventory)
			analyticsGroup.GET("/export/fuel", analyticsHandler.ExportFuel)
		}

		// Notifications
		notificationsGroup := api.Group("/notifications")
		notificationsGroup.Use(middleware.AuthMiddleware(jwtSecret, dbPool))
		{
			notificationsGroup.GET("", notificationHandler.ListNotifications)
			notificationsGroup.GET("/unread-count", notificationHandler.GetUnreadCount)
			notificationsGroup.PATCH("/:id/read", notificationHandler.MarkAsRead)
			notificationsGroup.POST("/mark-all-read", notificationHandler.MarkAllAsRead)
			notificationsGroup.DELETE("/:id", notificationHandler.DeleteNotification)
		}

		// 🚀 GPS TRACKING & GEOFENCING (PRO FEATURE #5)
		gpsService := services.NewGPSTrackingService(dbPool)
		gpsHandler := handlers.NewGPSTrackingHandler(gpsService, auditService)

		{
			gpsGroup := r.Group("/api/gps")
			gpsGroup.Use(middleware.AuthMiddleware(jwtSecret, dbPool))

			// PRO endpoints
			gpsGroup.POST("/locations", middleware.RequireSubscriptionTier("PRO", dbPool), gpsHandler.RecordVehicleLocation)
			gpsGroup.GET("/fleet-map", middleware.RequireSubscriptionTier("PRO", dbPool), gpsHandler.GetFleetMap)
			gpsGroup.GET("/trajectory", middleware.RequireSubscriptionTier("PRO", dbPool), gpsHandler.GetVehicleTrajectory)
			gpsGroup.POST("/geofences", middleware.RequireSubscriptionTier("PRO", dbPool), gpsHandler.CreateGeofence)
			gpsGroup.GET("/geofences", middleware.RequireSubscriptionTier("PRO", dbPool), gpsHandler.GetGeofences)
			gpsGroup.GET("/geofence-alerts", middleware.RequireSubscriptionTier("PRO", dbPool), gpsHandler.GetGeofenceAlerts)
			gpsGroup.GET("/fleet-status", middleware.RequireSubscriptionTier("PRO", dbPool), gpsHandler.GetFleetStatus)
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

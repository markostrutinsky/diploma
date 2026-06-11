package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"Omnilog_backend/internal/models"
	"Omnilog_backend/internal/services"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	service      *services.AnalyticsService
	auditService *services.AuditService
}

func NewAnalyticsHandler(s *services.AnalyticsService, auditService *services.AuditService) *AnalyticsHandler {
	return &AnalyticsHandler{
		service:      s,
		auditService: auditService,
	}
}

func (h *AnalyticsHandler) GetDashboard(c *gin.Context) {
	ctx := c.Request.Context()

	// Читаємо дати (за замовчуванням - останні 30 днів)
	defaultEnd := time.Now().Format("2006-01-02")
	defaultStart := time.Now().AddDate(0, 0, -30).Format("2006-01-02")

	startDateStr := c.DefaultQuery("start", defaultStart)
	endDateStr := c.DefaultQuery("end", defaultEnd)

	unitID := c.Query("unit_id")
	stats, err := h.service.GetDashboardAnalytics(ctx, startDateStr, endDateStr, unitID)
	if err != nil {
		// 🛡️ БЕЗПЕКА: Логуємо повну помилку на сервері, відправляємо безпечне повідомлення на клієнт
		log.Printf("ERROR: GetDashboard failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка завантаження аналітики"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// НОВА ФУНКЦІЯ: Smart Поповнення з модального вікна
func (h *AnalyticsHandler) AutoReplenish(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. Отримуємо ID користувача з middleware (зазвичай це "user_id" або "userID")
	userID := c.GetString("user_id")
	if userID == "" {
		// Fallback, якщо раптом middleware назвав змінну інакше
		userID = c.GetString("userID")
	}
	if userID == "" {
		userID = "system_smart_replenish" // Якщо зовсім немає ID, ставимо системний маркер
	}

	// 2. Читаємо JSON з тіла запиту
	var req models.SmartReplenishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 3. Відправляємо в сервіс
	count, err := h.service.RunSmartReplenish(ctx, req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process smart replenish"})
		return
	}

	go func(uID string, requestsCount int) {
		_ = h.auditService.LogAction(context.Background(), uID, "CREATE", "SMART_REPLENISH", "", fmt.Sprintf("Сформовано автоматичні заявки на поповнення (%d шт)", requestsCount))
	}(userID, count)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully created requests", "count": count})
}

// ExportInventory віддає XLSX файл із залишками
func (h *AnalyticsHandler) ExportInventory(c *gin.Context) {
	userID := c.GetString("user_id")

	// Дозволяємо фільтрувати по філії (опціонально)
	var unitID *int
	if idStr := c.Query("unit_id"); idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err == nil {
			unitID = &id
		}
	}

	fileBytes, err := h.service.GenerateInventoryExcel(c.Request.Context(), unitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не вдалося згенерувати звіт"})
		return
	}

	go func(uID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "EXPORT", "INVENTORY", "", "Експорт звіту залишків на складах (Excel)")
	}(userID)

	fileName := fmt.Sprintf("Inventory_Report_%s.xlsx", time.Now().Format("2006-01-02"))

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)
}

// ExportFuel віддає XLSX файл з витратами пального
func (h *AnalyticsHandler) ExportFuel(c *gin.Context) {
	userID := c.GetString("user_id")
	startStr := c.Query("start")
	endStr := c.Query("end")
	var unitID *int
	if idStr := c.Query("unit_id"); idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err == nil {
			unitID = &id
		}
	}

	// Дефолт: за останні 30 днів
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	if startStr != "" {
		if parsed, err := time.Parse("2006-01-02", startStr); err == nil {
			startDate = parsed
		}
	}
	if endStr != "" {
		if parsed, err := time.Parse("2006-01-02", endStr); err == nil {
			endDate = parsed.Add(23*time.Hour + 59*time.Minute) // Кінець дня
		}
	}

	fileBytes, err := h.service.GenerateFuelExcel(c.Request.Context(), startDate, endDate, unitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не вдалося згенерувати звіт"})
		return
	}

	go func(uID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "EXPORT", "FUEL", "", "Експорт звіту витрат пального (Excel)")
	}(userID)

	fileName := fmt.Sprintf("Fuel_Report_%s.xlsx", time.Now().Format("2006-01-02"))

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)
}

// 🚀 PRO FEATURE #1: GetAdvancedKPIs повертає розширену аналітику
// GET /api/analytics/kpi?start=2026-01-01&end=2026-04-22&unit_id=1
func (h *AnalyticsHandler) GetAdvancedKPIs(c *gin.Context) {
	ctx := c.Request.Context()

	startDate := c.DefaultQuery("start", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDate := c.DefaultQuery("end", time.Now().Format("2006-01-02"))
	unitIDStr := c.Query("unit_id")

	var unitID int64
	if unitIDStr != "" {
		if id, err := strconv.ParseInt(unitIDStr, 10, 64); err == nil {
			unitID = id
		}
	}

	kpis, err := h.service.GetAdvancedKPIs(ctx, startDate, endDate, unitID)
	if err != nil {
		log.Printf("ERROR: GetAdvancedKPIs failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка завантаження KPI"})
		return
	}

	c.JSON(http.StatusOK, kpis)
}

// 🚀 PRO FEATURE #2: GetDemandForecast прогнозує попит на 3 місяці вперед
// GET /api/analytics/forecast?unit_id=1
func (h *AnalyticsHandler) GetDemandForecast(c *gin.Context) {
	ctx := c.Request.Context()

	unitIDStr := c.Query("unit_id")
	var unitID int64
	if unitIDStr != "" {
		if id, err := strconv.ParseInt(unitIDStr, 10, 64); err == nil {
			unitID = id
		}
	}

	forecast, err := h.service.GetDemandForecast(ctx, unitID)
	if err != nil {
		log.Printf("ERROR: GetDemandForecast failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка завантаження прогнозу"})
		return
	}

	c.JSON(http.StatusOK, forecast)
}

// GetPredictiveMaintenanceSchedule returns predicted maintenance schedule for vehicles
func (h *AnalyticsHandler) GetPredictiveMaintenanceSchedule(c *gin.Context) {
	ctx := c.Request.Context()

	unitIDStr := c.Query("unit_id")
	var unitID int64
	if unitIDStr != "" {
		if id, err := strconv.ParseInt(unitIDStr, 10, 64); err == nil {
			unitID = id
		}
	}

	schedule, err := h.service.GetPredictiveMaintenanceSchedule(ctx, unitID)
	if err != nil {
		log.Printf("ERROR: GetPredictiveMaintenanceSchedule failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка завантаження прогнозу обслуговування"})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// GetFuelAnomalyDetection detects fuel consumption anomalies
func (h *AnalyticsHandler) GetFuelAnomalyDetection(c *gin.Context) {
	ctx := c.Request.Context()

	unitIDStr := c.Query("unit_id")
	var unitID int64
	if unitIDStr != "" {
		if id, err := strconv.ParseInt(unitIDStr, 10, 64); err == nil {
			unitID = id
		}
	}

	anomalies, err := h.service.GetFuelAnomalyDetection(ctx, unitID)
	if err != nil {
		log.Printf("ERROR: GetFuelAnomalyDetection failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка аналізу витрат палива"})
		return
	}

	c.JSON(http.StatusOK, anomalies)
}

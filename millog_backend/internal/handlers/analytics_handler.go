package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"millog_backend/internal/models"
	"millog_backend/internal/services"

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load dashboard analytics"})
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
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)
}

// ExportFuel віддає XLSX файл з витратами пального
func (h *AnalyticsHandler) ExportFuel(c *gin.Context) {
	userID := c.GetString("user_id")
	startStr := c.Query("start")
	endStr := c.Query("end")

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

	fileBytes, err := h.service.GenerateFuelExcel(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не вдалося згенерувати звіт"})
		return
	}

	go func(uID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "EXPORT", "FUEL", "", "Експорт звіту витрат пального (Excel)")
	}(userID)

	fileName := fmt.Sprintf("Fuel_Report_%s.xlsx", time.Now().Format("2006-01-02"))

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)
}

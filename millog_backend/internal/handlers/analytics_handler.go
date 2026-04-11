package handlers

import (
	"net/http"
	"time"

	"millog_backend/internal/models"
	"millog_backend/internal/services"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	service *services.AnalyticsService
}

func NewAnalyticsHandler(s *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: s}
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

	c.JSON(http.StatusOK, gin.H{"message": "Successfully created requests", "count": count})
}

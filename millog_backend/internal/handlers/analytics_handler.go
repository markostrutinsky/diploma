package handlers

import (
	"log"
	"net/http"

	"millog_backend/internal/services"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	service *services.AnalyticsService
}

func NewAnalyticsHandler(s *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: s}
}

// GetDashboard обробляє запит на отримання статистики
func (h *AnalyticsHandler) GetDashboard(c *gin.Context) {
	ctx := c.Request.Context()

	// ...
	stats, err := h.service.GetDashboardAnalytics(ctx)
	if err != nil {
		// ДОДАЄМО ВИВІД ПОМИЛКИ В КОНСОЛЬ (імпортуй пакет "log", якщо його немає)
		log.Printf("🔥 ПОМИЛКА АНАЛІТИКИ: %v", err)

		// ВІДПРАВЛЯЄМО ТЕКСТ ПОМИЛКИ НА ФРОНТ (замість стандартного тексту)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// ...

	c.JSON(http.StatusOK, stats)
}

package handlers

import (
	"net/http"
	"time"

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

	stats, err := h.service.GetDashboardAnalytics(ctx, startDateStr, endDateStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load dashboard analytics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *AnalyticsHandler) AutoReplenish(c *gin.Context) {
	ctx := c.Request.Context()
	count, err := h.service.RunAutoReplenish(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create auto-requests"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully created requests", "count": count})
}

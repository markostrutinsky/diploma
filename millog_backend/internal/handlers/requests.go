package handlers

import (
	"fmt"
	"millog_backend/internal/middleware"
	"millog_backend/internal/models"
	"millog_backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	reqService *services.RequestService
}

func NewRequestHandler(svc *services.RequestService) *RequestHandler {
	return &RequestHandler{reqService: svc}
}

func (h *RequestHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	var req models.CreateSupplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sr, err := h.reqService.Create(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sr)
}

func (h *RequestHandler) List(c *gin.Context) {
	// 1. Дістаємо дані з токена
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не знайдено токен авторизації"})
		return
	}
	claims := claimsVal.(*middleware.Claims)

	userRole := models.UserRole(claims.Role)

	list, err := h.reqService.List(c.Request.Context(), string(userRole), &claims.UnitID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if list == nil {
		list = []models.SupplyRequest{}
	}

	c.JSON(http.StatusOK, list)
}

func (h *RequestHandler) Approve(c *gin.Context) {
	id := c.Param("id")

	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не знайдено токен авторизації"})
		return
	}
	claims := claimsVal.(*middleware.Claims)

	userID := fmt.Sprint(claims.ID)
	userRole := models.UserRole(claims.Role)

	var req models.ApproveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.reqService.Approve(c.Request.Context(), id, userID, userRole, req.Approved, req.Comment); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Заявку успішно оброблено"})
}

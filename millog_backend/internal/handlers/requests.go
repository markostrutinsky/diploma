package handlers

import (
	"context"
	"fmt"
	"millog_backend/internal/middleware"
	"millog_backend/internal/models"
	"millog_backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	reqService   *services.RequestService
	auditService *services.AuditService
}

func NewRequestHandler(svc *services.RequestService, auditService *services.AuditService) *RequestHandler {
	return &RequestHandler{reqService: svc, auditService: auditService}
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

	go func(uID string, entityID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "CREATE", "SUPPLY_REQUEST", entityID, "Створено нову заявку на забезпечення")
	}(userID, fmt.Sprintf("%v", sr.ID))

	c.JSON(http.StatusCreated, sr)
}

func (h *RequestHandler) List(c *gin.Context) {
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

	go func(uID string, entityID string, isApproved bool) {
		actionDesc := "Відхилено заявку (на етапі погодження)"
		if isApproved {
			actionDesc = "Погоджено заявку"
		}
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "SUPPLY_REQUEST", entityID, actionDesc)
	}(userID, id, req.Approved)

	c.JSON(http.StatusOK, gin.H{"message": "Заявку успішно оброблено"})
}

func (h *RequestHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	req, err := h.reqService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заявку не знайдено"})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *RequestHandler) Reject(c *gin.Context) {
	reqID := c.Param("id")
	userID := c.GetString("user_id")

	var req struct {
		Comment string `json:"comment" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Вкажіть причину відмови"})
		return
	}

	if err := h.reqService.Reject(c.Request.Context(), reqID, req.Comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка відхилення"})
		return
	}

	go func(u, r, comment string) {
		_ = h.auditService.LogAction(context.Background(), u, "REJECT", "SUPPLY_REQUEST", r, "Відхилено заявку: "+comment)
	}(userID, reqID, req.Comment)

	c.JSON(http.StatusOK, gin.H{"message": "Заявку відхилено"})
}

func (h *RequestHandler) Cancel(c *gin.Context) {
	reqID := c.Param("id")
	userID := c.GetString("user_id")

	if err := h.reqService.Cancel(c.Request.Context(), reqID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func(u, r string) {
		_ = h.auditService.LogAction(context.Background(), u, "CANCEL", "SUPPLY_REQUEST", r, "Скасовано власну заявку")
	}(userID, reqID)

	c.JSON(http.StatusOK, gin.H{"message": "Заявку скасовано"})
}

// SmartDispatchPreview відповідає за генерацію прев'ю розумного пакування
func (h *RequestHandler) SmartDispatchPreview(c *gin.Context) {
	var req models.SmartDispatchReq

	// Валідація вхідного JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некоректний запит: необхідно передати масив request_ids"})
		return
	}

	// Виклик сервісу
	result, err := h.reqService.GetSmartDispatchPreview(c.Request.Context(), req.RequestIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Повертаємо готовий результат
	c.JSON(http.StatusOK, result)
}

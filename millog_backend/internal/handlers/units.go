package handlers

import (
	"Omnilog_backend/internal/middleware"
	"Omnilog_backend/internal/models"
	"Omnilog_backend/internal/services"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UnitHandler struct {
	svc          *services.UnitService
	auditService *services.AuditService
}

func NewUnitHandler(svc *services.UnitService, auditService *services.AuditService) *UnitHandler {
	return &UnitHandler{svc: svc, auditService: auditService}
}

func (h *UnitHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")

	var req models.CreateUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roleVal, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "не вдалося ідентифікувати користувача"})
		return
	}

	creatorRole, ok := roleVal.(models.UserRole)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "помилка обробки ролі користувача"})
		return
	}

	u, err := h.svc.Create(c.Request.Context(), &req, creatorRole)
	if err != nil {
		if strings.Contains(err.Error(), "недостатньо прав") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func(uID string, entityID string, name string) {
		_ = h.auditService.LogAction(context.Background(), uID, "CREATE", "UNIT", entityID, "Створено новий відділ: "+name)
	}(userID, fmt.Sprintf("%v", u.ID), req.Name)

	c.JSON(http.StatusCreated, u)
}

func (h *UnitHandler) List(c *gin.Context) {
	userIDVal, existsId := c.Get("user_id")
	roleVal, existsRole := c.Get("user_role")

	if !existsId || !existsRole {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "неавторизований доступ"})
		return
	}

	userID := userIDVal.(string)
	role := roleVal.(models.UserRole)

	units, err := h.svc.GetVisibleUnits(c.Request.Context(), userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, units)
}

func (h *UnitHandler) GetAvailableForRole(c *gin.Context) {
	role := c.Query("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "параметр role є обов'язковим"})
		return
	}

	units, err := h.svc.GetAvailableForRole(c.Request.Context(), role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, units)
}

func (h *UnitHandler) ChangeManager(c *gin.Context) {
	userID := c.GetString("user_id")
	unitIDStr := c.Param("id")
	targetUnitID, err := strconv.ParseInt(unitIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Невалідний ID відділу"})
		return
	}

	var req models.ChangeManagerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не вказано нового керівника"})
		return
	}

	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не знайдено токен авторизації"})
		return
	}
	claims := claimsVal.(*middleware.Claims)

	err = h.svc.ChangeCommander(c.Request.Context(), targetUnitID, req.NewManagerID, string(claims.Role), claims.UnitID)
	if err != nil {
		if err.Error() == "відмовлено в доступі: ви не можете змінити керівника цього відділу" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func(uID string, entityID string, newMgr string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "UNIT", entityID, "Змінено керівника відділу на користувача: "+newMgr)
	}(userID, unitIDStr, req.NewManagerID)

	c.JSON(http.StatusOK, gin.H{"message": "Керівника успішно змінено"})
}

func (h *UnitHandler) GetMyHierarchyForRole(c *gin.Context) {
	role := c.Query("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "параметр role є обов'язковим"})
		return
	}

	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "неавторизований доступ"})
		return
	}
	claims := claimsVal.(*middleware.Claims)

	if claims.UnitID == 0 {
		c.JSON(http.StatusOK, []models.Unit{})
		return
	}

	units, err := h.svc.GetMyHierarchyForRole(c.Request.Context(), role, claims.UnitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, units)
}

func (h *UnitHandler) UpdateUnit(c *gin.Context) {
	unitID := c.Param("id")
	userID := c.GetString("user_id")

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильні дані запиту"})
		return
	}

	err := h.svc.UpdateUnit(c.Request.Context(), unitID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка оновлення відділу"})
		return
	}

	go func(uID, unID, name string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "UNIT", unID, "Оновлено назву відділу на: "+name)
	}(userID, unitID, req.Name)

	c.JSON(http.StatusOK, gin.H{"message": "Відділ оновлено"})
}

func (h *UnitHandler) DeleteUnit(c *gin.Context) {
	unitID := c.Param("id")
	userID := c.GetString("user_id")

	err := h.svc.DeleteUnit(c.Request.Context(), unitID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неможливо видалити: до відділу прив'язані люди або майно"})
		return
	}

	go func(uID, unID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "DELETE", "UNIT", unID, "Видалено відділ")
	}(userID, unitID)

	c.JSON(http.StatusOK, gin.H{"message": "Відділ видалено"})
}

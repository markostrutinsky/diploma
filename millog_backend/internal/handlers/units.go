package handlers

import (
	"context"
	"millog_backend/internal/middleware"
	"millog_backend/internal/models"
	"millog_backend/internal/services"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "role parameter is required"})
		return
	}

	units, err := h.svc.GetAvailableForRole(c.Request.Context(), role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, units)
}

func (h *UnitHandler) ChangeCommander(c *gin.Context) {
	unitIDStr := c.Param("id")
	targetUnitID, err := strconv.ParseInt(unitIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Невалідний ID підрозділу"})
		return
	}

	var req models.ChangeCommanderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не вказано нового командира"})
		return
	}

	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не знайдено токен авторизації"})
		return
	}
	claims := claimsVal.(*middleware.Claims)

	err = h.svc.ChangeCommander(c.Request.Context(), targetUnitID, req.NewCommanderID, string(claims.Role), claims.UnitID)
	if err != nil {
		if err.Error() == "відмовлено в доступі: ви не можете змінити командира цього підрозділу" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Командира успішно змінено"})
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

// UpdateUnit оновлює дані підрозділу
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка оновлення підрозділу"})
		return
	}

	go func(uID, unID, name string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "UNIT", unID, "Оновлено назву підрозділу на: "+name)
	}(userID, unitID, req.Name)

	c.JSON(http.StatusOK, gin.H{"message": "Підрозділ оновлено"})
}

// DeleteUnit видаляє підрозділ
func (h *UnitHandler) DeleteUnit(c *gin.Context) {
	unitID := c.Param("id")
	userID := c.GetString("user_id")

	err := h.svc.DeleteUnit(c.Request.Context(), unitID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неможливо видалити: до підрозділу прив'язані люди або майно"})
		return
	}

	go func(uID, unID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "DELETE", "UNIT", unID, "Видалено підрозділ (організаційно-штатна структура)")
	}(userID, unitID)

	c.JSON(http.StatusOK, gin.H{"message": "Підрозділ ліквідовано"})
}

package handlers

import (
	"context"
	"millog_backend/internal/middleware"
	"millog_backend/internal/models"
	"millog_backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WarehouseHandler struct {
	service      *services.WarehouseService
	auditService *services.AuditService
}

func NewWarehouseHandler(service *services.WarehouseService, audit *services.AuditService) *WarehouseHandler {
	return &WarehouseHandler{service: service, auditService: audit}
}

func (h *WarehouseHandler) Create(c *gin.Context) {
	var req models.CreateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "не знайдено токен авторизації"})
		return
	}
	claims := claimsVal.(*middleware.Claims)

	// Логіка доступу:
	// Якщо це НЕ Адмін, примусово ставимо його власний UnitID.
	// Якщо Адмін — дозволяємо зберегти той UnitID, який прийшов з фронтенду.
	if claims.Role != models.RoleAdmin {
		if claims.UnitID == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "відмовлено в доступі: ваш профіль не прив'язаний до підрозділу"})
			return
		}
		req.UnitID = claims.UnitID
	}

	w, err := h.service.CreateWarehouse(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, w)
}

func (h *WarehouseHandler) List(c *gin.Context) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "не знайдено токен авторизації"})
		return
	}
	claims := claimsVal.(*middleware.Claims)

	var targetUnitID int64

	// ГОЛОВНА ЛОГІКА ВІДОБРАЖЕННЯ:
	// Адмін бачить усі склади (передаємо 0). Інші — тільки свої.
	if claims.Role == models.RoleAdmin {
		targetUnitID = 0
	} else {
		if claims.UnitID == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "відмовлено в доступі: ваш профіль не прив'язаний до підрозділу"})
			return
		}
		targetUnitID = claims.UnitID
	}

	list, err := h.service.ListWarehouses(c.Request.Context(), targetUnitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

// Додай цей метод до WarehouseHandler
func (h *WarehouseHandler) UpdateLocation(c *gin.Context) {
	warehouseID := c.Param("id")

	// Структура для парсингу JSON
	var req struct {
		Latitude  float64 `json:"latitude" binding:"required"`
		Longitude float64 `json:"longitude" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Невірний формат координат"})
		return
	}

	err := h.service.UpdateLocation(c.Request.Context(), warehouseID, req.Latitude, req.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка оновлення дислокації: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Дислокацію успішно оновлено"})
}

func (h *WarehouseHandler) UpdateWarehouse(c *gin.Context) {
	warehouseID := c.Param("id")
	userID := c.GetString("user_id")

	var req struct {
		Name          string `json:"name" binding:"required"`
		CapacityLevel string `json:"capacity_level"`
		ZoneType      string `json:"zone_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильні дані запиту"})
		return
	}

	if err := h.service.UpdateWarehouse(c.Request.Context(), warehouseID, req.Name, req.CapacityLevel, req.ZoneType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка оновлення складу"})
		return
	}

	// 🛡️ Аудит
	go func(u, w, n string) {
		_ = h.auditService.LogAction(context.Background(), u, "UPDATE", "WAREHOUSE", w, "Оновлено параметри складу: "+n)
	}(userID, warehouseID, req.Name)

	c.JSON(http.StatusOK, gin.H{"message": "Склад оновлено"})
}

// 3. Додаємо метод видалення
func (h *WarehouseHandler) Delete(c *gin.Context) {
	warehouseID := c.Param("id")
	userID := c.GetString("user_id")

	if err := h.service.DeleteWarehouse(c.Request.Context(), warehouseID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Поверне нашу помилку про непустий склад
		return
	}

	// 🛡️ Аудит
	go func(u, w string) {
		_ = h.auditService.LogAction(context.Background(), u, "DELETE", "WAREHOUSE", w, "Видалено порожній склад")
	}(userID, warehouseID)

	c.JSON(http.StatusOK, gin.H{"message": "Склад видалено"})
}

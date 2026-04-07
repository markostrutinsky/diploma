package handlers

import (
	"fmt"
	"millog_backend/internal/middleware"
	"millog_backend/internal/models"
	"millog_backend/internal/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	invService *services.InventoryService
}

func NewInventoryHandler(inv *services.InventoryService) *InventoryHandler {
	return &InventoryHandler{invService: inv}
}

func (h *InventoryHandler) CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat, err := h.invService.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func (h *InventoryHandler) ListCategories(c *gin.Context) {
	list, err := h.invService.ListCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *InventoryHandler) CreateResource(c *gin.Context) {
	var req models.CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.invService.CreateResource(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *InventoryHandler) ListResources(c *gin.Context) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "claims not found"})
		return
	}
	claims := claimsVal.(*middleware.Claims)
	userRole := string(claims.Role)
	userUnitID := claims.UnitID

	var finalUnitID *int64

	var requestedUnitID *int64
	if s := c.Query("unit_id"); s != "" {
		if id, err := strconv.ParseInt(s, 10, 64); err == nil {
			requestedUnitID = &id
		}
	}

	if userRole == "ADMIN" {
		finalUnitID = requestedUnitID
	} else {
		if requestedUnitID != nil {
			finalUnitID = requestedUnitID
		} else {
			finalUnitID = &userUnitID
		}
	}

	list, err := h.invService.ListResources(c.Request.Context(), finalUnitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *InventoryHandler) WriteOff(c *gin.Context) {
	resourceID := c.Param("id")
	var req models.WriteOffResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невірний формат: вкажіть кількість (quantity)"})
		return
	}

	err := h.invService.WriteOff(c.Request.Context(), resourceID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Успішно списано"})
}

func (h *InventoryHandler) UpdateResource(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resource id is required"})
		return
	}

	var req models.UpdateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format: " + err.Error()})
		return
	}

	err := h.invService.UpdateResource(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "repository update failed: resource not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
			return
		}

		fmt.Println("🚨 UPDATE ERROR:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "resource successfully updated"})
}

func (h *InventoryHandler) Transfer(c *gin.Context) {
	// 1. Беремо ID майна з URL (наприклад: /api/resources/:id/transfer)
	resourceID := c.Param("id")
	if resourceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не вказано ID ресурсу"})
		return
	}

	// 2. Парсимо тіло запиту (кількість, цільовий склад)
	var req models.TransferResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невірний формат даних: " + err.Error()})
		return
	}

	// 3. Викликаємо сервіс
	err := h.invService.Transfer(c.Request.Context(), resourceID, req)
	if err != nil {
		// Якщо помилка пов'язана з нестачею майна, краще віддавати 400 Bad Request
		if strings.Contains(err.Error(), "недостатньо майна") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Успішна відповідь
	c.JSON(http.StatusOK, gin.H{"message": "Майно успішно переміщено"})
}

func (h *InventoryHandler) Delete(c *gin.Context) {
	resourceID := c.Param("id")
	if resourceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не вказано ID ресурсу"})
		return
	}

	err := h.invService.Delete(c.Request.Context(), resourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ресурс успішно видалено"})
}

func (h *InventoryHandler) Assign(c *gin.Context) {
	resourceID := c.Param("id")
	var req models.AssignResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невірний формат запиту"})
		return
	}

	err := h.invService.Assign(c.Request.Context(), resourceID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Майно успішно видано"})
}

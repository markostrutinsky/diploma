package handlers

import (
	"millog_backend/internal/models"
	"millog_backend/internal/services"
	"net/http"
	"strconv"

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
	var unitID *int64
	if s := c.Query("unit_id"); s != "" {
		if id, err := strconv.ParseInt(s, 10, 64); err == nil {
			unitID = &id
		}
	}
	list, err := h.invService.ListResources(c.Request.Context(), unitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

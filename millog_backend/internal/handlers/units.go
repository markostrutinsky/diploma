package handlers

import (
	"millog_backend/internal/models"
	"millog_backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UnitHandler struct {
	svc *services.UnitService
}

func NewUnitHandler(svc *services.UnitService) *UnitHandler {
	return &UnitHandler{svc: svc}
}

func (h *UnitHandler) Create(c *gin.Context) {
	var req models.CreateUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, u)
}

func (h *UnitHandler) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

package handlers

import (
	"millog_backend/internal/models"
	"millog_backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VolunteerRequestHandler struct {
	svc *services.VolunteerRequestService
}

func NewVolunteerRequestHandler(svc *services.VolunteerRequestService) *VolunteerRequestHandler {
	return &VolunteerRequestHandler{svc: svc}
}

func (h *VolunteerRequestHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	var req models.CreateVolunteerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vr, err := h.svc.Create(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, vr)
}

func (h *VolunteerRequestHandler) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *VolunteerRequestHandler) Take(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")
	if err := h.svc.Take(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Request taken"})
}

func (h *VolunteerRequestHandler) Complete(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")
	if err := h.svc.Complete(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Request completed"})
}

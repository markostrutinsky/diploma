package handlers

import (
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
	list, err := h.reqService.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *RequestHandler) Approve(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")
	var req models.ApproveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.reqService.Approve(c.Request.Context(), id, userID, req.Approved, req.Comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Request processed"})
}

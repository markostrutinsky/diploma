package handlers

import (
	"millog_backend/internal/models"
	"millog_backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CONTRACTORRequestHandler struct {
	svc *services.CONTRACTORRequestService
}

func NewCONTRACTORRequestHandler(svc *services.CONTRACTORRequestService) *CONTRACTORRequestHandler {
	return &CONTRACTORRequestHandler{svc: svc}
}

func (h *CONTRACTORRequestHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	tokenUnitID := c.GetInt64("unit_id")

	var req models.CreateCONTRACTORRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	finalUnitID := tokenUnitID
	if req.UnitID != nil {
		finalUnitID = *req.UnitID
	}

	vr, err := h.svc.Create(c.Request.Context(), userID, finalUnitID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, vr)
}

// Список заявок
func (h *CONTRACTORRequestHandler) List(c *gin.Context) {
	status := c.Query("status")

	list, err := h.svc.List(c.Request.Context(), models.CONTRACTORRequestStatus(status))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if list == nil {
		list = []models.CONTRACTORRequest{}
	}

	c.JSON(http.StatusOK, list)
}

// Волонтер бере в роботу
func (h *CONTRACTORRequestHandler) Take(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")
	if err := h.svc.UpdateStatus(c.Request.Context(), id, userID, models.CONTRACTORTaken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Заявку взято в роботу"})
}

// Волонтер відправляє/доставляє (Замість старого Complete)
func (h *CONTRACTORRequestHandler) Deliver(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")
	if err := h.svc.UpdateStatus(c.Request.Context(), id, userID, models.CONTRACTORDelivered); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Статус змінено на 'Доставлено'. Очікує прийомки."})
}

// ВІЙСЬКОВИЙ: Скасування заявки (поки вона ще OPEN)
func (h *CONTRACTORRequestHandler) Cancel(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")
	if err := h.svc.UpdateStatus(c.Request.Context(), id, userID, models.CONTRACTORCanceled); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Заявку скасовано"})
}

// ВІЙСЬКОВИЙ: Відхилення доставленого майна (брак)
func (h *CONTRACTORRequestHandler) Reject(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")
	if err := h.svc.UpdateStatus(c.Request.Context(), id, userID, models.CONTRACTORRejected); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Прийомку відхилено"})
}

// ВІЙСЬКОВИЙ (Комірник/Старшина): Прийомка на баланс (НАЙГОЛОВНІШЕ)
func (h *CONTRACTORRequestHandler) Accept(c *gin.Context) {
	commanderID := c.GetString("user_id") // ID того, хто приймає
	unitID := c.GetInt64("unit_id")       // Склад підрозділу
	requestID := c.Param("id")

	var payload models.AcceptCONTRACTORPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильні дані для складу: " + err.Error()})
		return
	}

	// Передаємо в сервіс: ID заявки, хто приймає, на який склад (unitID) і самі дані про майно
	if err := h.svc.AcceptAndStore(c.Request.Context(), requestID, commanderID, unitID, payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Майно успішно прийнято на баланс підрозділу!"})
}

package handlers

import (
	"Omnilog_backend/internal/models"
	"Omnilog_backend/internal/services"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContractorRequestHandler struct {
	svc          *services.ContractorRequestService
	auditService *services.AuditService
}

func NewContractorRequestHandler(svc *services.ContractorRequestService, auditService *services.AuditService) *ContractorRequestHandler {
	return &ContractorRequestHandler{svc: svc, auditService: auditService}
}

func (h *ContractorRequestHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	tokenUnitID := c.GetInt64("unit_id")

	var req models.CreateContractorRequest
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

	// 🛡️ Аудит: Створення
	go func(uID string, entityID string) {
		_ = h.auditService.LogAction(
			context.Background(),
			uID,
			"CREATE",
			"CONTRACTOR_REQUEST",
			entityID,
			"Створено новий зовнішній запит",
		)
	}(userID, vr.ID)

	c.JSON(http.StatusCreated, vr)
}

// Список заявок
func (h *ContractorRequestHandler) List(c *gin.Context) {
	status := c.Query("status")

	list, err := h.svc.List(c.Request.Context(), models.ContractorRequestStatus(status))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if list == nil {
		list = []models.ContractorRequest{}
	}

	c.JSON(http.StatusOK, list)
}

// Волонтер бере в роботу
func (h *ContractorRequestHandler) Take(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	if err := h.svc.UpdateStatus(c.Request.Context(), id, userID, models.ContractorTaken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 🛡️ Аудит: Оновлення статусу
	go func(uID string, entityID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "CONTRACTOR_REQUEST", entityID, "Заявку взято в роботу")
	}(userID, id)

	c.JSON(http.StatusOK, gin.H{"message": "Заявку взято в роботу"})
}

// Волонтер відправляє/доставляє
func (h *ContractorRequestHandler) Deliver(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	if err := h.svc.UpdateStatus(c.Request.Context(), id, userID, models.ContractorDelivered); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 🛡️ Аудит
	go func(uID string, entityID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "CONTRACTOR_REQUEST", entityID, "Статус змінено на 'Доставлено'")
	}(userID, id)

	c.JSON(http.StatusOK, gin.H{"message": "Статус змінено на 'Доставлено'. Очікує прийомки."})
}

// ВІЙСЬКОВИЙ: Скасування заявки
func (h *ContractorRequestHandler) Cancel(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	if err := h.svc.UpdateStatus(c.Request.Context(), id, userID, models.ContractorCanceled); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 🛡️ Аудит
	go func(uID string, entityID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "CONTRACTOR_REQUEST", entityID, "Заявку скасовано")
	}(userID, id)

	c.JSON(http.StatusOK, gin.H{"message": "Заявку скасовано"})
}

// ВІЙСЬКОВИЙ: Відхилення доставленого майна
func (h *ContractorRequestHandler) Reject(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	if err := h.svc.UpdateStatus(c.Request.Context(), id, userID, models.ContractorRejected); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 🛡️ Аудит
	go func(uID string, entityID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "CONTRACTOR_REQUEST", entityID, "Прийомку відхилено (брак)")
	}(userID, id)

	c.JSON(http.StatusOK, gin.H{"message": "Прийомку відхилено"})
}

// ВІЙСЬКОВИЙ (Комірник/Старшина): Прийомка на баланс
func (h *ContractorRequestHandler) Accept(c *gin.Context) {
	commanderID := c.GetString("user_id") // ID того, хто приймає
	unitID := c.GetInt64("unit_id")       // Склад підрозділу
	requestID := c.Param("id")

	var payload models.AcceptContractorPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильні дані для складу: " + err.Error()})
		return
	}

	if err := h.svc.AcceptAndStore(c.Request.Context(), requestID, commanderID, unitID, payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 🛡️ Аудит: Прийомка
	go func(uID string, entityID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "CONTRACTOR_REQUEST", entityID, "Майно успішно прийнято на баланс")
	}(commanderID, requestID) // Тут передаємо commanderID, бо дію робить військовий

	c.JSON(http.StatusOK, gin.H{"message": "Майно успішно прийнято на баланс підрозділу!"})
}

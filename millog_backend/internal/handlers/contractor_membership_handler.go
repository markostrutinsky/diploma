package handlers

import (
	"context"
	"net/http"

	"Omnilog_backend/internal/models"
	"Omnilog_backend/internal/services"

	"github.com/gin-gonic/gin"
)

type ContractorMembershipHandler struct {
	svc          *services.ContractorMembershipService
	auditService *services.AuditService
}

func NewContractorMembershipHandler(svc *services.ContractorMembershipService, auditService *services.AuditService) *ContractorMembershipHandler {
	return &ContractorMembershipHandler{svc: svc, auditService: auditService}
}

// ListForTenant — адмін організації бачить підрядників, які подалися на співпрацю.
// Опційний фільтр ?status=PENDING|APPROVED|REJECTED.
func (h *ContractorMembershipHandler) ListForTenant(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Ця операція доступна лише в контексті організації"})
		return
	}
	status := models.ContractorMembershipStatus(c.Query("status"))

	list, err := h.svc.ListByTenant(c.Request.Context(), tenantID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if list == nil {
		list = []models.ContractorMembership{}
	}
	c.JSON(http.StatusOK, list)
}

// Approve — адмін організації підтверджує підрядника.
func (h *ContractorMembershipHandler) Approve(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	deciderID := c.GetString("user_id")
	id := c.Param("id")

	if tenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Ця операція доступна лише в контексті організації"})
		return
	}

	if err := h.svc.Approve(c.Request.Context(), id, tenantID, deciderID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func(uID, entityID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "CONTRACTOR_MEMBERSHIP", entityID, "Підрядника схвалено для співпраці")
	}(deciderID, id)

	c.JSON(http.StatusOK, gin.H{"message": "Підрядника схвалено"})
}

// Reject — адмін організації відхиляє підрядника.
func (h *ContractorMembershipHandler) Reject(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	deciderID := c.GetString("user_id")
	id := c.Param("id")

	if tenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Ця операція доступна лише в контексті організації"})
		return
	}

	if err := h.svc.Reject(c.Request.Context(), id, tenantID, deciderID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func(uID, entityID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "CONTRACTOR_MEMBERSHIP", entityID, "Підрядника відхилено")
	}(deciderID, id)

	c.JSON(http.StatusOK, gin.H{"message": "Підрядника відхилено"})
}

// ListMine — self-view підрядника: з якими організаціями він співпрацює та в якому статусі.
func (h *ContractorMembershipHandler) ListMine(c *gin.Context) {
	userID := c.GetString("user_id")

	list, err := h.svc.ListByContractor(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if list == nil {
		list = []models.ContractorMembership{}
	}
	c.JSON(http.StatusOK, list)
}

// Apply — підрядник самостійно надсилає заявку на співпрацю з організацією (explicit,
// без спроби взяти завдання). Повертає поточний статус членства, щоб фронтенд одразу
// показав коректний стан (надіслано / вже очікує / вже схвалено / відхилено).
func (h *ContractorMembershipHandler) Apply(c *gin.Context) {
	contractorID := c.GetString("user_id")

	var req struct {
		TenantID string `json:"tenant_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Вкажіть організацію для співпраці"})
		return
	}

	status, err := h.svc.Apply(c.Request.Context(), contractorID, req.TenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func(uID, entityID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "CREATE", "CONTRACTOR_MEMBERSHIP", entityID, "Підрядник надіслав заявку на співпрацю")
	}(contractorID, req.TenantID)

	c.JSON(http.StatusOK, gin.H{"status": status, "message": "Заявку на співпрацю надіслано"})
}

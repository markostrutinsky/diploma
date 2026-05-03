package handlers

import (
	"Omnilog_backend/internal/middleware"
	"Omnilog_backend/internal/models"
	"Omnilog_backend/internal/services"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	reqService   *services.RequestService
	auditService *services.AuditService
	slaMonitor   *services.SLAMonitor
}

func NewRequestHandler(svc *services.RequestService, auditService *services.AuditService, slaMonitor *services.SLAMonitor) *RequestHandler {
	return &RequestHandler{reqService: svc, auditService: auditService, slaMonitor: slaMonitor}
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

	go func(uID string, entityID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "CREATE", "SUPPLY_REQUEST", entityID, "Створено нову заявку на забезпечення")
	}(userID, fmt.Sprintf("%v", sr.ID))

	c.JSON(http.StatusCreated, sr)
}

func (h *RequestHandler) List(c *gin.Context) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не знайдено токен авторизації"})
		return
	}
	claims := claimsVal.(*middleware.Claims)

	userRole := models.UserRole(claims.Role)

	list, err := h.reqService.List(c.Request.Context(), string(userRole), &claims.UnitID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if list == nil {
		list = []models.SupplyRequest{}
	}

	c.JSON(http.StatusOK, list)
}

func (h *RequestHandler) Approve(c *gin.Context) {
	id := c.Param("id")

	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не знайдено токен авторизації"})
		return
	}
	claims := claimsVal.(*middleware.Claims)

	userID := fmt.Sprint(claims.ID)
	userRole := models.UserRole(claims.Role)

	var req models.ApproveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.reqService.Approve(c.Request.Context(), id, userID, userRole, req.Approved, req.Comment); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	go func(uID string, entityID string, isApproved bool) {
		actionDesc := "Відхилено заявку (на етапі погодження)"
		if isApproved {
			actionDesc = "Погоджено заявку"
		}
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "SUPPLY_REQUEST", entityID, actionDesc)
	}(userID, id, req.Approved)

	c.JSON(http.StatusOK, gin.H{"message": "Заявку успішно оброблено"})
}

func (h *RequestHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	req, err := h.reqService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заявку не знайдено"})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *RequestHandler) Reject(c *gin.Context) {
	reqID := c.Param("id")
	userID := c.GetString("user_id")

	var req struct {
		Comment string `json:"comment" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Вкажіть причину відмови"})
		return
	}

	if err := h.reqService.Reject(c.Request.Context(), reqID, req.Comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка відхилення"})
		return
	}

	go func(u, r, comment string) {
		_ = h.auditService.LogAction(context.Background(), u, "REJECT", "SUPPLY_REQUEST", r, "Відхилено заявку: "+comment)
	}(userID, reqID, req.Comment)

	c.JSON(http.StatusOK, gin.H{"message": "Заявку відхилено"})
}

func (h *RequestHandler) Cancel(c *gin.Context) {
	reqID := c.Param("id")
	userID := c.GetString("user_id")

	if err := h.reqService.Cancel(c.Request.Context(), reqID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func(u, r string) {
		_ = h.auditService.LogAction(context.Background(), u, "CANCEL", "SUPPLY_REQUEST", r, "Скасовано власну заявку")
	}(userID, reqID)

	c.JSON(http.StatusOK, gin.H{"message": "Заявку скасовано"})
}

// SmartDispatchPreview відповідає за генерацію прев'ю розумного пакування
func (h *RequestHandler) SmartDispatchPreview(c *gin.Context) {
	var req models.SmartDispatchReq

	// Валідація вхідного JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некоректний запит: необхідно передати масив request_ids"})
		return
	}

	// Виклик сервісу
	result, err := h.reqService.GetSmartDispatchPreview(c.Request.Context(), req.RequestIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Повертаємо готовий результат
	c.JSON(http.StatusOK, result)
}

// SmartDispatchConfirm — затверджує результат Smart Розподілу: для кожної
// машини створюється окремий shipment (та сама логіка, що й у ручного рейсу).
func (h *RequestHandler) SmartDispatchConfirm(c *gin.Context) {
	userID := c.GetString("user_id")

	var req models.SmartDispatchConfirmReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := h.reqService.ConfirmSmartDispatch(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             err.Error(),
			"successful_routes": count,
		})
		return
	}

	go func(u string, n int) {
		_ = h.auditService.LogAction(
			context.Background(),
			u,
			"CREATE",
			"SHIPMENT",
			"smart-batch",
			fmt.Sprintf("Smart Розподіл: створено %d рейсів", n),
		)
	}(userID, count)

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Успішно створено %d рейсів", count),
		"count":   count,
	})
}

func (h *RequestHandler) TriggerCheck(c *gin.Context) {
	ctx := c.Request.Context()

	newlyEscalated, err := h.slaMonitor.CheckPendingRequests(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка під час перевірки SLA"})
		return
	}

	// Додатковий контекст для UX: скільки вже ескальованих заявок висить у системі
	// та скільки ще PENDING (і скільки з них наближаються до дедлайну).
	existingEscalated, _ := h.slaMonitor.GetEscalatedCount(ctx)
	pendingTotal, pendingSoon, _ := h.slaMonitor.GetPendingStats(ctx)

	var message string
	if newlyEscalated > 0 {
		message = fmt.Sprintf("Ескальовано %d нових заявок. Всього ESCALATED у системі: %d.",
			newlyEscalated, existingEscalated)
	} else if existingEscalated > 0 {
		message = fmt.Sprintf("Нових порушень SLA не знайдено. У системі вже є %d ескальованих заявок, які потребують уваги менеджера.",
			existingEscalated)
	} else {
		message = "Порушень SLA не знайдено. Всі заявки у межах нормативу."
	}

	c.JSON(http.StatusOK, gin.H{
		"message":            message,
		"escalated_count":    newlyEscalated,    // скільки НОВИХ ескалацій за цей запуск
		"existing_escalated": existingEscalated, // скільки вже ESCALATED у системі
		"pending_total":      pendingTotal,      // скільки зараз у PENDING
		"pending_near_sla":   pendingSoon,       // з них наближаються до дедлайну
	})
}

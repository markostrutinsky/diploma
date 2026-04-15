package handlers

import (
	"context"
	"fmt"
	"millog_backend/internal/middleware"
	"millog_backend/internal/models"
	"millog_backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	invService   *services.InventoryService
	auditService *services.AuditService
}

func NewInventoryHandler(inv *services.InventoryService, audit *services.AuditService) *InventoryHandler {
	return &InventoryHandler{invService: inv, auditService: audit}
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

func (h *InventoryHandler) GetResource(c *gin.Context) {
	resourceID := c.Param("id")
	if resourceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не вказано ID ресурсу"})
		return
	}

	res, err := h.invService.GetResource(c.Request.Context(), resourceID)
	if err != nil {
		if err.Error() == "repository get failed: resource not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "ресурс не знайдено"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
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
	userID := c.GetString("user_id")
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

	go func(uID string, rID string) {
		// Використовуємо context.Background(), бо контекст запиту 'c' вмирає після відповіді
		_ = h.auditService.LogAction(
			context.Background(),
			uID,
			"WRITE_OFF",
			"RESOURCE",
			rID,
			"Списання майна зі складу", // Можна додати кількість, якщо вона є в змінній
		)
	}(userID, resourceID)

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

func (h *InventoryHandler) Delete(c *gin.Context) {
	resourceID := c.Param("id")
	userID := c.GetString("user_id")
	if resourceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не вказано ID ресурсу"})
		return
	}

	err := h.invService.Delete(c.Request.Context(), resourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	go func(uID string, rID string) {
		_ = h.auditService.LogAction(
			context.Background(),
			uID,
			"DELETE",
			"RESOURCE",
			rID,
			"Безповоротне видалення картки майна",
		)
	}(userID, resourceID)
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

func (h *InventoryHandler) GetMyEquipment(c *gin.Context) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не знайдено токен авторизації"})
		return
	}
	claims := claimsVal.(*middleware.Claims)

	// ФІКС: Беремо правильне поле з токена
	userID := claims.UserID

	items, err := h.invService.GetMyEquipment(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка отримання майна: " + err.Error()})
		return
	}

	if items == nil {
		items = []models.MyEquipmentItem{}
	}

	c.JSON(http.StatusOK, items)
}

// POST /api/inventory/issue
func (h *InventoryHandler) IssueResource(c *gin.Context) {
	var req models.IssueResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Невірний формат даних: " + err.Error()})
		return
	}

	// Дістаємо дані командира з JWT токена
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизовано"})
		return
	}
	claims := claimsVal.(*middleware.Claims)

	var unitID *int64
	if claims.UnitID != 0 {
		val := claims.UnitID
		unitID = &val
	}

	// Виконуємо видачу
	err := h.invService.IssueResource(c.Request.Context(), unitID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Майно успішно видано військовослужбовцю"})
}

// POST /api/shipments
func (h *InventoryHandler) CreateShipment(c *gin.Context) {
	var req models.CreateShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Невірний формат даних: " + err.Error()})
		return
	}

	err := h.invService.CreateShipment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Рейс успішно сформовано, транспорт відправлено"})
}

func (h *InventoryHandler) GetByWarehouse(c *gin.Context) {
	warehouseID := c.Param("id")

	items, err := h.invService.GetByWarehouse(c.Request.Context(), warehouseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Якщо товарів немає, віддаємо порожній масив замість null
	if items == nil {
		items = []models.InventoryItem{}
	}

	c.JSON(http.StatusOK, items)
}

func (h *InventoryHandler) ListShipments(c *gin.Context) {
	list, err := h.invService.ListShipments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *InventoryHandler) ReceiveShipment(c *gin.Context) {
	shipmentID := c.Param("id")
	err := h.invService.ReceiveShipment(c.Request.Context(), shipmentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Вантаж успішно прийнято на склад!"})
}

// DownloadShipmentPDF віддає згенеровану ТТН у форматі PDF
func (h *InventoryHandler) DownloadShipmentPDF(c *gin.Context) {
	shipmentID := c.Param("id")

	// Викликаємо сервіс
	pdfBytes, err := h.invService.GenerateShipmentPDF(c.Request.Context(), shipmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Кажемо браузеру завантажити цей файл з конкретною назвою
	filename := fmt.Sprintf("Waybill_%s.pdf", shipmentID)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/pdf")

	// Віддаємо байти (HTTP статус 200)
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// DownloadResourceQR віддає згенерований QR-код у форматі PNG
func (h *InventoryHandler) DownloadResourceQR(c *gin.Context) {
	resourceID := c.Param("id")

	pngBytes, err := h.invService.GenerateResourceQR(resourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fileName := fmt.Sprintf("QR_%s.png", resourceID)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Header("Content-Type", "image/png")

	c.Data(http.StatusOK, "image/png", pngBytes)
}

// UpdateCategory оновлює назву або опис категорії
func (h *InventoryHandler) UpdateCategory(c *gin.Context) {
	categoryID := c.Param("id")
	userID := c.GetString("user_id")

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильні дані запиту"})
		return
	}

	err := h.invService.UpdateCategory(c.Request.Context(), categoryID, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не вдалося оновити категорію"})
		return
	}

	go func(uID, cID, name string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "CATEGORY", cID, "Оновлено назву/опис категорії: "+name)
	}(userID, categoryID, req.Name)

	c.JSON(http.StatusOK, gin.H{"message": "Категорію успішно оновлено"})
}

// DeleteCategory видаляє категорію (якщо вона порожня)
func (h *InventoryHandler) DeleteCategory(c *gin.Context) {
	categoryID := c.Param("id")
	userID := c.GetString("user_id")

	err := h.invService.DeleteCategory(c.Request.Context(), categoryID)
	if err != nil {
		if err.Error() == "категорія не порожня" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неможливо видалити: у цій категорії ще є майно"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка видалення категорії"})
		return
	}

	go func(uID, cID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "DELETE", "CATEGORY", cID, "Видалено категорію майна")
	}(userID, categoryID)

	c.JSON(http.StatusOK, gin.H{"message": "Категорію видалено"})
}

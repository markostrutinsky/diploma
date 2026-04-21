package handlers

import (
	"context"
	"fmt"
	"millog_backend/internal/middleware"
	"millog_backend/internal/models"
	"millog_backend/internal/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type InventoryHandler struct {
	invService        *services.InventoryService
	auditService      *services.AuditService
	limitationService *services.LimitationService
	authService       *services.AuthService
}

func NewInventoryHandler(inv *services.InventoryService, audit *services.AuditService, limitation *services.LimitationService, auth *services.AuthService) *InventoryHandler {
	return &InventoryHandler{
		invService:        inv,
		auditService:      audit,
		limitationService: limitation,
		authService:       auth,
	}
}

func (h *InventoryHandler) CreateCategory(c *gin.Context) {
	userID := c.GetString("user_id")
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

	go func(uID string, entityID string, name string) {
		_ = h.auditService.LogAction(context.Background(), uID, "CREATE", "CATEGORY", entityID, "Створено нову категорію: "+name)
	}(userID, fmt.Sprintf("%v", cat.ID), req.Name)

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
	userID := c.GetString("user_id")
	var req models.CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 🛡️ Перевіримо ліміт на ресурсах
	if err := h.limitationService.CheckResourceLimit(c.Request.Context(), req.UnitID); err != nil {
		c.JSON(http.StatusPaymentRequired, gin.H{
			"error":       err.Error(),
			"upgrade_url": "/billing?plan=pro",
		})
		return
	}

	res, err := h.invService.CreateResource(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func(uID string, entityID string, name string) {
		_ = h.auditService.LogAction(context.Background(), uID, "CREATE", "RESOURCE", entityID, "Створено нову картку майна: "+name)
	}(userID, fmt.Sprintf("%v", res.ID), req.Name)

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
		_ = h.auditService.LogAction(
			context.Background(),
			uID,
			"WRITE_OFF",
			"RESOURCE",
			rID,
			"Списання майна зі складу",
		)
	}(userID, resourceID)

	c.JSON(http.StatusOK, gin.H{"message": "Успішно списано"})
}

func (h *InventoryHandler) UpdateResource(c *gin.Context) {
	userID := c.GetString("user_id")
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

	go func(uID string, rID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "RESOURCE", rID, "Оновлено дані картки майна")
	}(userID, id)

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
	userID := c.GetString("user_id")
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

	go func(uID string, rID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "RESOURCE", rID, "Майно видано користувачу")
	}(userID, resourceID)

	c.JSON(http.StatusOK, gin.H{"message": "Майно успішно видано"})
}

func (h *InventoryHandler) GetMyEquipment(c *gin.Context) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не знайдено токен авторизації"})
		return
	}
	claims := claimsVal.(*middleware.Claims)
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

func (h *InventoryHandler) IssueResource(c *gin.Context) {
	userID := c.GetString("user_id")
	var req models.IssueResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Невірний формат даних: " + err.Error()})
		return
	}

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

	err := h.invService.IssueResource(c.Request.Context(), unitID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func(uID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "RESOURCE", "", "Майно видано користувачу")
	}(userID)

	c.JSON(http.StatusOK, gin.H{"message": "Майно успішно видано користувачу"})
}

func (h *InventoryHandler) CreateShipment(c *gin.Context) {
	userID := c.GetString("user_id")
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

	go func(uID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "CREATE", "SHIPMENT", "", "Сформовано новий рейс")
	}(userID)

	c.JSON(http.StatusOK, gin.H{"message": "Рейс успішно сформовано, транспорт відправлено"})
}

func (h *InventoryHandler) GetByWarehouse(c *gin.Context) {
	warehouseID := c.Param("id")

	items, err := h.invService.GetByWarehouse(c.Request.Context(), warehouseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
	userID := c.GetString("user_id")
	shipmentID := c.Param("id")
	err := h.invService.ReceiveShipment(c.Request.Context(), shipmentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func(uID string, sID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "SHIPMENT", sID, "Вантаж прийнято на склад")
	}(userID, shipmentID)

	c.JSON(http.StatusOK, gin.H{"message": "Вантаж успішно прийнято на склад!"})
}

func (h *InventoryHandler) DownloadShipmentPDF(c *gin.Context) {
	userID := c.GetString("user_id")
	shipmentID := c.Param("id")

	pdfBytes, err := h.invService.GenerateShipmentPDF(c.Request.Context(), shipmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func(uID string, sID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "READ", "SHIPMENT", sID, "Завантажено ТТН рейсу")
	}(userID, shipmentID)

	filename := fmt.Sprintf("Waybill_%s.pdf", shipmentID)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

func (h *InventoryHandler) DownloadResourceQR(c *gin.Context) {
	userID := c.GetString("user_id")
	resourceID := c.Param("id")

	pngBytes, err := h.invService.GenerateResourceQR(resourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func(uID string, rID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "READ", "RESOURCE", rID, "Завантажено QR-код майна")
	}(userID, resourceID)

	fileName := fmt.Sprintf("QR_%s.png", resourceID)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Header("Content-Type", "image/png")

	c.Data(http.StatusOK, "image/png", pngBytes)
}

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

func (h *InventoryHandler) SubmitAudit(c *gin.Context) {
	var req models.SubmitAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Невірний формат даних: " + err.Error()})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не знайдено ID користувача"})
		return
	}

	err := h.invService.SubmitInventoryAudit(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка збереження результатів переобліку: " + err.Error()})
		return
	}

	go func(uID, wID string) {
		_ = h.auditService.LogAction(
			context.Background(),
			uID,
			"INVENTORY_AUDIT",
			"WAREHOUSE",
			wID,
			"Проведено переоблік складу. Зафіксовано акт розбіжностей.",
		)
	}(userID, req.WarehouseID)

	c.JSON(http.StatusOK, gin.H{"message": "Акт переобліку успішно сформовано, залишки оновлено!"})
}

// Допоміжна функція для очищення тексту перед порівнянням
func normalizeText(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// 1. Простий і красивий шаблон (без організацій і складів)
func (h *InventoryHandler) DownloadImportTemplate(c *gin.Context) {
	f := excelize.NewFile()
	defer f.Close()
	sheet := "Шаблон"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{
		"Назва майна (Обов'язково)",
		"Назва категорії (Обов'язково)",
		"Кількість",
		"Одиниці (PCS/KG/L/KIT)",
		"Вага 1 од. (кг)",
		"Заводський штрих-код",
		"Мін. залишок",
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, header)
	}

	f.SetColWidth(sheet, "A", "B", 35)
	f.SetColWidth(sheet, "C", "G", 15)

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=OmniLog_Import_Template.xlsx")
	f.Write(c.Writer)
}

// 1. Оновлений Імпорт (з виводом помилок у консоль)
func (h *InventoryHandler) ImportExcel(c *gin.Context) {
	unitIDStr := c.PostForm("unit_id")
	warehouseID := c.PostForm("warehouse_id")
	unitID, _ := strconv.ParseInt(unitIDStr, 10, 64)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл не знайдено"})
		return
	}

	file, _ := fileHeader.Open()
	defer file.Close()
	f, _ := excelize.OpenReader(file)
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil || len(rows) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл порожній"})
		return
	}

	ctx := c.Request.Context()
	categories, err := h.invService.GetAllCategories(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Критична помилка завантаження категорій з бази: %v", err),
		})
		return
	}

	successCount := 0
	var importErrors []string

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		getCol := func(idx int) string {
			if idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}

		name := getCol(0)
		categoryName := getCol(1)
		if name == "" || categoryName == "" {
			continue
		}

		// 1. Шукаємо або створюємо категорію (безпечним методом)
		catID, catErr := h.smartCategoryLookup(ctx, categoryName, &categories)
		if catErr != nil {
			importErrors = append(importErrors, fmt.Sprintf("Рядок %d (%s): %v", i+1, name, catErr))
			continue // Блокуємо і йдемо до наступного
		}

		qty, _ := strconv.Atoi(getCol(2))
		unitType := getCol(3)
		if unitType == "" {
			unitType = "PCS"
		}

		// Захист від ком замість крапок (напр. 0,6 замість 0.6)
		weightStr := strings.ReplaceAll(getCol(4), ",", ".")
		weight, _ := strconv.ParseFloat(weightStr, 64)
		if weight <= 0 {
			weight = 1.0
		}
		minQty, _ := strconv.Atoi(getCol(6))

		// 2. Перевіряємо, чи існує вже такий ресурс
		existingRes, err := h.invService.GetResourceByNameAndWarehouse(ctx, name, warehouseID)

		if err == nil && existingRes != nil {
			// РЕСУРС ІСНУЄ -> Плюсуємо кількість
			newQty := existingRes.Quantity + qty
			updateErr := h.invService.UpdateResourceQuantity(ctx, existingRes.ID, newQty)
			if updateErr != nil {
				importErrors = append(importErrors, fmt.Sprintf("Рядок %d: Помилка оновлення кількості: %v", i+1, updateErr))
			} else {
				successCount++
			}
		} else {
			// РЕСУРС НОВИЙ -> Створюємо
			req := models.CreateResourceRequest{
				UnitID:      unitID,
				WarehouseID: &warehouseID,
				CategoryID:  catID,
				Name:        name,
				Quantity:    qty,
				UnitType:    models.UnitMeasurement(unitType),
				WeightKg:    weight,
				Barcode:     getCol(5),
				MinQuantity: minQty,
				Condition:   models.ConditionNew,
			}

			_, createErr := h.invService.CreateResource(ctx, &req)
			if createErr != nil {
				importErrors = append(importErrors, fmt.Sprintf("Рядок %d (%s): Помилка створення БД - %v", i+1, name, createErr))
			} else {
				successCount++
			}
		}
	}

	// 🔥 ДРУКУЄМО ПОМИЛКИ В КОНСОЛЬ GO ДЛЯ ДЕБАГУ
	if len(importErrors) > 0 {
		fmt.Println("=== ПОМИЛКИ ІМПОРТУ ===")
		for _, e := range importErrors {
			fmt.Println(e)
		}
		fmt.Println("=======================")
	}

	c.JSON(http.StatusOK, gin.H{
		"success_count":   successCount,
		"total_processed": len(rows) - 1,
		"errors":          importErrors,
	})
}

// 2. Виправлений, безпечний розумний пошук
func (h *InventoryHandler) smartCategoryLookup(ctx context.Context, name string, cats *[]models.ResourceCategory) (string, error) {
	search := strings.ToLower(strings.TrimSpace(name))
	if search == "" {
		return "", fmt.Errorf("порожня назва категорії")
	}

	// Етап 1: Шукаємо ТОЧНИЙ збіг (пріоритет)
	for _, c := range *cats {
		if strings.ToLower(c.Name) == search {
			return c.ID, nil
		}
	}

	// Етап 2: Безпечний частковий збіг
	// (Ігноруємо слова менше 4 символів, щоб не ловити "та", "іт", "ка")
	for _, c := range *cats {
		dbName := strings.ToLower(c.Name)
		if len(dbName) > 3 && strings.Contains(search, dbName) {
			return c.ID, nil
		}
		if len(search) > 3 && strings.Contains(dbName, search) {
			return c.ID, nil
		}
	}

	// Етап 3: Якщо нічого не підійшло - створюємо нову
	newCat, err := h.invService.CreateCategory(ctx, &models.CreateCategoryRequest{Name: name})
	if err != nil {
		return "", fmt.Errorf("БД відмовила у створенні категорії: %v", err)
	}
	if newCat != nil && newCat.ID != "" {
		newCat.Name = name
		*cats = append(*cats, *newCat) // Запам'ятовуємо, щоб наступні рядки її використовували
		return newCat.ID, nil
	}

	return "", fmt.Errorf("невідома помилка отримання ID нової категорії")
}

package handlers

import (
	"Omnilog_backend/internal/models"
	"Omnilog_backend/internal/services"
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type VehicleHandler struct {
	service      *services.VehicleService
	auditService *services.AuditService
}

func NewVehicleHandler(service *services.VehicleService, auditService *services.AuditService) *VehicleHandler {
	return &VehicleHandler{service: service, auditService: auditService}
}

func (h *VehicleHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id") // ДОДАНО: Дістаємо ID користувача для аудиту

	var req models.CreateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Перевірте правильність заповнення полів (марка, номер, тип, вантажопідйомність, бак та норма є обов'язковими)"})
		return
	}

	vehicle := &models.Vehicle{
		Brand:           req.Brand,
		Model:           req.Model,
		PlateNumber:     req.PlateNumber,
		Type:            req.Type,
		CapacityKg:      req.CapacityKg,
		TankCapacity:    req.TankCapacity,
		FuelNorm:        req.FuelNorm,
		DriverID:        req.DriverID,
		HomeWarehouseID: req.HomeWarehouseID,
		Status:          "ACTIVE",
	}

	err := h.service.CreateVehicle(c.Request.Context(), vehicle)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 🛡️ Аудит: Створення авто
	go func(uID string, entityID string) {
		_ = h.auditService.LogAction(
			context.Background(),
			uID,
			"CREATE",
			"VEHICLE",
			entityID,
			"Додано новий автомобіль: "+req.Brand+" ("+req.PlateNumber+")",
		)
	}(userID, fmt.Sprintf("%v", vehicle.ID)) // fmt.Sprintf захищає від будь-яких конфліктів типів

	c.JSON(http.StatusCreated, vehicle)
}

func (h *VehicleHandler) GetAll(c *gin.Context) {
	vehicles, err := h.service.GetAllVehicles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка завантаження списку автомобілів"})
		return
	}

	c.JSON(http.StatusOK, vehicles)
}

func (h *VehicleHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	vehicle, err := h.service.GetVehicleByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vehicle)
}

func (h *VehicleHandler) UpdateStatus(c *gin.Context) {
	userID := c.GetString("user_id") // ДОДАНО
	vehicleID := c.Param("id")

	var req models.VehicleStatusUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некоректні дані статусу. Перевірте правильність заповнення."})
		return
	}

	err := h.service.UpdateStatus(c.Request.Context(), vehicleID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 🛡️ Аудит: Зміна статусу
	go func(uID string, entityID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "VEHICLE", entityID, "Змінено статус автомобіля")
	}(userID, vehicleID)

	c.JSON(http.StatusOK, gin.H{"message": "Статус автомобіля успішно оновлено"})
}

type MaintenanceRequest struct {
	CurrentOdometer int     `json:"current_odometer" binding:"required,min=0"`
	Description     string  `json:"description" binding:"required"`
	PerformedBy     string  `json:"performed_by"`
	CostAmount      float64 `json:"cost_amount"`
}

func (h *VehicleHandler) PerformMaintenance(c *gin.Context) {
	vehicleID := c.Param("id")
	userID := c.GetString("user_id") // ДОДАНО: ID користувача для аудиту

	// 1. Отримуємо текстові поля з форми
	odometerStr := c.PostForm("current_odometer")
	description := c.PostForm("description")
	performedBy := c.PostForm("performed_by")
	costAmountStr := c.PostForm("cost_amount")

	// Валідація обов'язкових полів
	if odometerStr == "" || description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Одометр та опис є обов'язковими"})
		return
	}

	// Конвертуємо рядки в числа з перевіркою помилок
	odometer, err := strconv.Atoi(odometerStr)
	if err != nil || odometer < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Одометр повинен бути невід'ємним числом"})
		return
	}

	costAmount := 0.0
	if costAmountStr != "" {
		var parseErr error
		costAmount, parseErr = strconv.ParseFloat(costAmountStr, 64)
		if parseErr != nil || costAmount < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Вартість повинна бути невід'ємним числом"})
			return
		}
	}

	var documentURL string

	file, err := c.FormFile("document")
	if err == nil {
		// 🛡️ Ліміт розміру файлу: 10 МБ
		const maxFileSizeBytes = 10 << 20 // 10 MB
		if file.Size > maxFileSizeBytes {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Файл занадто великий. Максимум 10 МБ."})
			return
		}

		// 🛡️ Валідація розширення файлу
		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowedExts := map[string]bool{".pdf": true, ".jpg": true, ".jpeg": true, ".png": true}
		if !allowedExts[ext] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Допускаються тільки PDF, JPG та PNG"})
			return
		}

		// 🛡️ Генеруємо безпечне ім'я без user-input
		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		savePath := filepath.Join("uploads", "maintenance", newFileName)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не вдалося зберегти документ"})
			return
		}

		documentURL = "/api/" + filepath.ToSlash(savePath)
	}

	record := &models.MaintenanceRecord{
		VehicleID:   vehicleID,
		OdometerKm:  odometer,
		Description: description,
		PerformedBy: performedBy,
		CostAmount:  costAmount,
		DocumentURL: documentURL,
	}

	err = h.service.PerformMaintenance(c.Request.Context(), record)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 🛡️ Аудит: Проведення ТО
	go func(uID string, entityID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "VEHICLE", entityID, "Внесено запис про проходження ТО")
	}(userID, vehicleID)

	c.JSON(http.StatusOK, record)
}

func (h *VehicleHandler) ScheduleMaintenance(c *gin.Context) {
	vehicleID := c.Param("id")
	userID := c.GetString("user_id")

	var req struct {
		OdometerKm   int    `json:"odometer_km"`
		ServiceType  string `json:"service_type"`
		ScheduledFor string `json:"scheduled_for"`
		Description  string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некоректні дані планування ТО"})
		return
	}
	if req.OdometerKm < 0 || req.ServiceType == "" || req.ScheduledFor == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Одометр, тип ТО та дата планування є обов'язковими"})
		return
	}

	scheduledFor, err := time.Parse(time.RFC3339, req.ScheduledFor)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Дата планування має бути у форматі RFC3339"})
		return
	}

	description := strings.TrimSpace(req.Description)
	if description == "" {
		description = "Заплановане ТО"
	}

	record := &models.MaintenanceRecord{
		VehicleID:    vehicleID,
		OdometerKm:   req.OdometerKm,
		Description:  description,
		PerformedBy:  "Планування",
		ServiceType:  strings.ToUpper(strings.TrimSpace(req.ServiceType)),
		ScheduledFor: &scheduledFor,
	}

	if err := h.service.ScheduleMaintenance(c.Request.Context(), record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func(uID string, entityID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "CREATE", "VEHICLE", entityID, "Заплановано ТО")
	}(userID, vehicleID)

	c.JSON(http.StatusOK, record)
}

func (h *VehicleHandler) GetMaintenanceHistory(c *gin.Context) {
	vehicleID := c.Param("id")

	records, err := h.service.GetMaintenanceHistory(c.Request.Context(), vehicleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не вдалося завантажити історію ТО"})
		return
	}

	if records == nil {
		records = []*models.MaintenanceRecord{}
	}

	c.JSON(http.StatusOK, records)
}

func (h *VehicleHandler) AssignDriver(c *gin.Context) {
	userID := c.GetString("user_id") // ДОДАНО
	vehicleID := c.Param("id")

	var req struct {
		DriverID *string `json:"driver_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Невірний формат даних"})
		return
	}

	if err := h.service.AssignDriver(c.Request.Context(), vehicleID, req.DriverID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка оновлення водія: " + err.Error()})
		return
	}

	// 🛡️ Аудит: Призначення/Зняття водія
	go func(uID string, entityID string) {
		actionDesc := "Призначено нового водія на транспортний засіб"
		if req.DriverID == nil {
			actionDesc = "Знято водія з транспортного засобу"
		}
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "VEHICLE", entityID, actionDesc)
	}(userID, vehicleID)

	c.JSON(http.StatusOK, gin.H{"message": "Екіпаж успішно оновлено"})
}

func (h *VehicleHandler) GetDriverHistory(c *gin.Context) {
	vehicleID := c.Param("id")

	history, err := h.service.GetDriverHistory(c.Request.Context(), vehicleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка отримання історії екіпажів"})
		return
	}

	c.JSON(http.StatusOK, history)
}

func (h *VehicleHandler) Update(c *gin.Context) {
	vehicleID := c.Param("id")
	userID := c.GetString("user_id")

	var req struct {
		Brand       string  `json:"brand" binding:"required"`
		Model       string  `json:"model"`
		PlateNumber string  `json:"plate_number" binding:"required"`
		CapacityKg  float64 `json:"capacity_kg"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильні дані запиту"})
		return
	}

	if err := h.service.UpdateVehicle(c.Request.Context(), vehicleID, req.Brand, req.Model, req.PlateNumber, req.CapacityKg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка оновлення автомобіля"})
		return
	}

	// 🛡️ Аудит: Оновлення авто (привів до спільного стандарту)
	go func(uID string, entityID string, plate string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "VEHICLE", entityID, "Оновлено дані авто: "+plate)
	}(userID, vehicleID, req.PlateNumber)

	c.JSON(http.StatusOK, gin.H{"message": "Дані автомобіля оновлено"})
}

func (h *VehicleHandler) Delete(c *gin.Context) {
	vehicleID := c.Param("id")
	userID := c.GetString("user_id")

	if err := h.service.DeleteVehicle(c.Request.Context(), vehicleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 🛡️ Аудит: Списання авто
	go func(uID string, entityID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "DELETE", "VEHICLE", entityID, "Автомобіль списано/видалено")
	}(userID, vehicleID)

	c.JSON(http.StatusOK, gin.H{"message": "Автомобіль списано"})
}

func (h *VehicleHandler) GetAvailableForShipment(c *gin.Context) {
	fromWarehouseID := c.Query("from_warehouse_id")
	toWarehouseID := c.Query("to_warehouse_id")

	if fromWarehouseID == "" || toWarehouseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Вкажіть from_warehouse_id та to_warehouse_id"})
		return
	}

	vehicles, err := h.service.GetAvailableForRoute(c.Request.Context(), fromWarehouseID, toWarehouseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка завантаження автомобілів"})
		return
	}

	c.JSON(http.StatusOK, vehicles)
}

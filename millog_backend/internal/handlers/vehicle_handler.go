package handlers

import (
	"context"
	"fmt"
	"millog_backend/internal/models"
	"millog_backend/internal/services"
	"net/http"
	"path/filepath"
	"strconv"
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
		Brand:        req.Brand,
		Model:        req.Model,
		PlateNumber:  req.PlateNumber,
		Type:         req.Type,
		CapacityKg:   req.CapacityKg,
		TankCapacity: req.TankCapacity,
		FuelNorm:     req.FuelNorm,
		DriverID:     req.DriverID,
		Status:       "ACTIVE", // За замовчуванням вільна
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
	userID := c.GetString("user_id") // ДОДАНО
	vehicleID := c.Param("id")

	odometerStr := c.PostForm("current_odometer")
	description := c.PostForm("description")
	performedBy := c.PostForm("performed_by")
	costAmountStr := c.PostForm("cost_amount")

	if odometerStr == "" || description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Одометр та опис є обов'язковими"})
		return
	}

	odometer, _ := strconv.Atoi(odometerStr)
	costAmount, _ := strconv.ParseFloat(costAmountStr, 64)

	var documentURL string

	file, err := c.FormFile("document")
	if err == nil {
		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		savePath := filepath.Join("uploads", "maintenance", newFileName)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не вдалося зберегти документ"})
			return
		}

		documentURL = "/" + filepath.ToSlash(savePath)
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

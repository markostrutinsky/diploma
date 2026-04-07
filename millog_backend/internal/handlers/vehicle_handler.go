package handlers

import (
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
	service *services.VehicleService
}

func NewVehicleHandler(service *services.VehicleService) *VehicleHandler {
	return &VehicleHandler{service: service}
}

type CreateVehicleRequest struct {
	Brand        string  `json:"brand" binding:"required"`
	Model        string  `json:"model"`
	PlateNumber  string  `json:"plate_number" binding:"required"`
	TankCapacity float64 `json:"tank_capacity" binding:"required,gt=0"`
	FuelNorm     float64 `json:"fuel_norm" binding:"required,gt=0"`
}

func (h *VehicleHandler) Create(c *gin.Context) {
	var req CreateVehicleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Перевірте правильність заповнення полів (марка, номер, бак та норма є обов'язковими)"})
		return
	}

	vehicle := &models.Vehicle{
		Brand:        req.Brand,
		Model:        req.Model,
		PlateNumber:  req.PlateNumber,
		TankCapacity: req.TankCapacity,
		FuelNorm:     req.FuelNorm,
		Status:       models.VehicleActive, // Зверни увагу: константа має залишатись або бути "ACTIVE"
	}

	err := h.service.CreateVehicle(c.Request.Context(), vehicle)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

// ОНОВЛЕНО: Новий метод зміни статусу (з причиною)
func (h *VehicleHandler) UpdateStatus(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"message": "Статус автомобіля успішно оновлено"})
}

// Структура для запиту ТО (ОНОВЛЕНО: додали CostAmount)
type MaintenanceRequest struct {
	CurrentOdometer int     `json:"current_odometer" binding:"required,min=0"`
	Description     string  `json:"description" binding:"required"`
	PerformedBy     string  `json:"performed_by"`
	CostAmount      float64 `json:"cost_amount"` // Гроші
}

func (h *VehicleHandler) PerformMaintenance(c *gin.Context) {
	vehicleID := c.Param("id")

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

	// Конвертуємо рядки в числа
	odometer, _ := strconv.Atoi(odometerStr)
	costAmount, _ := strconv.ParseFloat(costAmountStr, 64)

	var documentURL string

	// 2. Намагаємось отримати файл (поле "document")
	file, err := c.FormFile("document")
	if err == nil {
		// Файл є! Створюємо унікальне ім'я, щоб не було конфліктів
		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		savePath := filepath.Join("uploads", "maintenance", newFileName)

		// Зберігаємо файл на диск
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не вдалося зберегти документ"})
			return
		}

		// Формуємо URL для бази даних (напр. /uploads/maintenance/1612345.pdf)
		documentURL = "/" + filepath.ToSlash(savePath)
	}

	// 3. Формуємо запис для бази
	record := &models.MaintenanceRecord{
		VehicleID:   vehicleID,
		OdometerKm:  odometer,
		Description: description,
		PerformedBy: performedBy,
		CostAmount:  costAmount,
		DocumentURL: documentURL, // Записуємо шлях до файлу (або порожньо, якщо файлу не було)
	}

	// 4. Зберігаємо все через сервіс
	err = h.service.PerformMaintenance(c.Request.Context(), record)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

// Отримання історії ТО
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
	vehicleID := c.Param("id")

	// Структура для прийому JSON. Вказівник (*string) потрібен, щоб прийняти null, якщо водія знімають.
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

package handlers

import (
	"millog_backend/internal/models"
	"millog_backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FuelHandler struct {
	FuelService *services.FuelService
}

type CreateFuelRequest struct {
	Liters     float64 `json:"liters" binding:"required,gt=0"`
	OdometerKm *int    `json:"odometer_km"`
	RecordType string  `json:"record_type" binding:"required,oneof=REFUEL EXPENSE"`
}

func NewFuelHandler(fuelService *services.FuelService) *FuelHandler {
	return &FuelHandler{FuelService: fuelService}
}

func (h *FuelHandler) CreateRecord(c *gin.Context) {
	vehicleID := c.Param("id")

	userID := c.GetString("userID")

	var req CreateFuelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некоректні дані: перевірте літри та тип запису"})
		return
	}

	record := &models.FuelRecord{
		VehicleID:  vehicleID,
		Liters:     req.Liters,
		OdometerKm: req.OdometerKm,
		RecordType: models.FuelRecordType(req.RecordType),
	}

	if userID != "" {
		record.CreatedBy = &userID
	}

	err := h.FuelService.AddFuelRecord(c.Request.Context(), record)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, record)
}

func (h *FuelHandler) GetHistory(c *gin.Context) {
	vehicleID := c.Param("id")

	records, err := h.FuelService.GetVehicleFuelHistory(c.Request.Context(), vehicleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не вдалося завантажити історію пального"})
		return
	}

	if records == nil {
		records = []*models.FuelRecord{}
	}

	c.JSON(http.StatusOK, records)
}

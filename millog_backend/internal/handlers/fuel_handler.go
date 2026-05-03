package handlers

import (
	"Omnilog_backend/internal/models"
	"Omnilog_backend/internal/services"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FuelHandler struct {
	FuelService  *services.FuelService
	auditService *services.AuditService
}

type CreateFuelRequest struct {
	Liters     float64 `json:"liters" binding:"required,gt=0"`
	OdometerKm *int    `json:"odometer_km"`
	RecordType string  `json:"record_type" binding:"required,oneof=REFUEL EXPENSE"`
}

func NewFuelHandler(fuelService *services.FuelService, auditService *services.AuditService) *FuelHandler {
	return &FuelHandler{FuelService: fuelService, auditService: auditService}
}

func (h *FuelHandler) CreateRecord(c *gin.Context) {
	vehicleID := c.Param("id")
	userID := c.GetString("user_id")

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

	go func(uID string, vID string, rType string, liters float64) {
		actionDesc := fmt.Sprintf("Додано запис про пальне (Тип: %s, Літри: %.2f)", rType, liters)
		_ = h.auditService.LogAction(context.Background(), uID, "CREATE", "FUEL_RECORD", vID, actionDesc)
	}(userID, vehicleID, req.RecordType, req.Liters)

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

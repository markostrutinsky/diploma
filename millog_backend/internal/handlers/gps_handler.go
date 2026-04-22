package handlers

import (
	"net/http"
	"strconv"
	"time"

	"millog_backend/internal/models"
	"millog_backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GPSTrackingHandler struct {
	gpsService *services.GPSTrackingService
	auditSvc   *services.AuditService
}

func NewGPSTrackingHandler(gpsService *services.GPSTrackingService, auditSvc *services.AuditService) *GPSTrackingHandler {
	return &GPSTrackingHandler{
		gpsService: gpsService,
		auditSvc:   auditSvc,
	}
}

// RecordVehicleLocation handles GPS updates from IoT devices or mobile apps
// POST /api/gps/locations
func (h *GPSTrackingHandler) RecordVehicleLocation(c *gin.Context) {
	unitID := c.GetInt64("unit_id")

	var req models.CreateGPSLocationRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	location := &models.GPSLocation{
		VehicleID: req.VehicleID,
		UnitID:    unitID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Altitude:  req.Altitude,
		Speed:     req.Speed,
		Heading:   req.Heading,
		Accuracy:  req.Accuracy,
		Timestamp: time.Now(),
	}

	if err := h.gpsService.RecordVehicleLocation(c.Request.Context(), location); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, location)
}

// GetFleetMap returns real-time positions of all vehicles (for live tracking dashboard)
// GET /api/gps/fleet-map
func (h *GPSTrackingHandler) GetFleetMap(c *gin.Context) {
	unitID := c.GetInt64("unit_id")

	locations, err := h.gpsService.GetFleetLocations(c.Request.Context(), unitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	plates, _ := h.gpsService.GetVehiclePlates(c.Request.Context())

	now := time.Now()
	vehicles := make([]gin.H, 0, len(locations))
	for _, loc := range locations {
		updatedSecondsAgo := int64(now.Sub(loc.Timestamp).Seconds())
		if updatedSecondsAgo < 0 {
			updatedSecondsAgo = 0
		}
		vehicles = append(vehicles, gin.H{
			"vehicle_id":          loc.VehicleID,
			"plate_number":        plates[loc.VehicleID],
			"latitude":            loc.Latitude,
			"longitude":           loc.Longitude,
			"speed":               loc.Speed,
			"heading":             loc.Heading,
			"timestamp":           loc.Timestamp,
			"updated_seconds_ago": updatedSecondsAgo,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"timestamp": now,
		"vehicles":  vehicles,
		"count":     len(vehicles),
	})
}

// GetVehicleTrajectory returns the path traveled by a vehicle in a time range
// GET /api/gps/trajectory?vehicle_id=<uuid>&start_time=2026-04-01T00:00:00Z&end_time=2026-04-22T23:59:59Z
func (h *GPSTrackingHandler) GetVehicleTrajectory(c *gin.Context) {
	vehicleID := c.Query("vehicle_id")
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	if vehicleID == "" || startTimeStr == "" || endTimeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query parameters: vehicle_id, start_time, end_time"})
		return
	}

	if _, err := uuid.Parse(vehicleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vehicle_id (expected UUID)"})
		return
	}

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_time format"})
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_time format"})
		return
	}

	locations, err := h.gpsService.GetVehicleTrajectory(c.Request.Context(), vehicleID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate distance traveled
	distance := h.gpsService.CalculateDistance(locations)

	c.JSON(http.StatusOK, gin.H{
		"vehicle_id":  vehicleID,
		"start_time":  startTime,
		"end_time":    endTime,
		"locations":   locations,
		"count":       len(locations),
		"distance_km": distance,
	})
}

// CreateGeofence creates a new alert boundary
// POST /api/gps/geofences
func (h *GPSTrackingHandler) CreateGeofence(c *gin.Context) {
	unitID := c.GetInt64("unit_id")

	var req models.Geofence
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UnitID = unitID
	if req.Active == false {
		req.Active = true // Default to active
	}

	if err := h.gpsService.CreateGeofence(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, req)
}

// GetGeofences returns all alert zones for the unit
// GET /api/gps/geofences
func (h *GPSTrackingHandler) GetGeofences(c *gin.Context) {
	unitID := c.GetInt64("unit_id")

	geofences, err := h.gpsService.GetGeofences(c.Request.Context(), unitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, geofences)
}

// GetGeofenceAlerts returns recent boundary breach events
// GET /api/gps/geofence-alerts?hours=24
func (h *GPSTrackingHandler) GetGeofenceAlerts(c *gin.Context) {
	unitID := c.GetInt64("unit_id")

	hours := 24 // default
	if hoursStr := c.Query("hours"); hoursStr != "" {
		if h, err := strconv.Atoi(hoursStr); err == nil {
			hours = h
		}
	}

	alerts, err := h.gpsService.GetGeofenceAlerts(c.Request.Context(), unitID, hours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, alerts)
}

// GetFleetStatus returns comprehensive fleet tracking data
// GET /api/gps/fleet-status
func (h *GPSTrackingHandler) GetFleetStatus(c *gin.Context) {
	unitID := c.GetInt64("unit_id")

	status, err := h.gpsService.GetDetailedFleetStatus(c.Request.Context(), unitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

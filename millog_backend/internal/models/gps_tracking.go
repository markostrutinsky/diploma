package models

import "time"

// GPSLocation represents a real-time GPS location of a vehicle
type GPSLocation struct {
	ID        int64     `json:"id"`
	VehicleID int64     `json:"vehicle_id"`
	UnitID    int64     `json:"unit_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Altitude  float64   `json:"altitude"`
	Speed     float64   `json:"speed"`           // km/h
	Heading   float64   `json:"heading"`         // degrees 0-360
	Accuracy  float64   `json:"accuracy"`        // meters
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

// VehicleTrack represents a historical track of a vehicle
type VehicleTrack struct {
	ID           int64          `json:"id"`
	VehicleID    int64          `json:"vehicle_id"`
	UnitID       int64          `json:"unit_id"`
	StartTime    time.Time      `json:"start_time"`
	EndTime      time.Time      `json:"end_time"`
	Distance     float64        `json:"distance"`         // km
	AverageSpeed float64        `json:"average_speed"`    // km/h
	MaxSpeed     float64        `json:"max_speed"`        // km/h
	LocationCount int           `json:"location_count"`
	CreatedAt    time.Time      `json:"created_at"`
}

// Geofence represents a geographic boundary for alerts
type Geofence struct {
	ID        int64     `json:"id"`
	UnitID    int64     `json:"unit_id"`
	Name      string    `json:"name"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Radius    float64   `json:"radius"` // meters
	Type      string    `json:"type"`   // WAREHOUSE, DANGER_ZONE, RESTRICTED_AREA
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GeofenceAlert represents an alert when vehicle enters/exits geofence
type GeofenceAlert struct {
	ID        int64     `json:"id"`
	VehicleID int64     `json:"vehicle_id"`
	GeofenceID int64    `json:"geofence_id"`
	EventType string    `json:"event_type"` // ENTER, EXIT
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateGPSLocationRequest is the request body for submitting GPS data
type CreateGPSLocationRequest struct {
	VehicleID int64   `json:"vehicle_id" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Altitude  float64 `json:"altitude"`
	Speed     float64 `json:"speed"`
	Heading   float64 `json:"heading"`
	Accuracy  float64 `json:"accuracy"`
}

// GetVehicleTrackRequest query parameters for vehicle history
type GetVehicleTrackRequest struct {
	VehicleID int64     `json:"vehicle_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}

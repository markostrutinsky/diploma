package models

import "time"

type VehicleStatus string

const (
	VehicleActive   VehicleStatus = "ACTIVE"
	VehicleInactive VehicleStatus = "INACTIVE"
)

type FuelRecordType string

const (
	FuelRefuel   FuelRecordType = "REFUEL"
	FuelExpense  FuelRecordType = "EXPENSE"
)

type Vehicle struct {
	ID         string        `json:"id"`
	Brand      string        `json:"brand"`
	Model      string        `json:"model"`
	PlateNumber string       `json:"plate_number"`
	Status     VehicleStatus `json:"status"`
	DriverID   *string       `json:"driver_id"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

type FuelRecord struct {
	ID         string        `json:"id"`
	VehicleID  string        `json:"vehicle_id"`
	Liters     float64       `json:"liters"`
	OdometerKm *int          `json:"odometer_km"`
	RecordType FuelRecordType `json:"record_type"`
	CreatedBy  *string       `json:"created_by"`
	CreatedAt  time.Time     `json:"created_at"`
}

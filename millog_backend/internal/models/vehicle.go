package models

import "time"

type VehicleStatus string

const (
	VehicleActive   VehicleStatus = "ACTIVE"
	VehicleInactive VehicleStatus = "INACTIVE"
	VehicleInRepair VehicleStatus = "IN_REPAIR" // ДОДАНО: Статус для ремонту/ТО
)

type FuelRecordType string

const (
	FuelRefuel  FuelRecordType = "REFUEL"
	FuelExpense FuelRecordType = "EXPENSE"
)

type Vehicle struct {
	ID           string        `json:"id"`
	Brand        string        `json:"brand"`
	Model        string        `json:"model"`
	PlateNumber  string        `json:"plate_number"`
	Status       VehicleStatus `json:"status"`
	DriverID     *string       `json:"driver_id"`
	TankCapacity float64       `json:"tank_capacity"`
	FuelNorm     float64       `json:"fuel_norm"`

	// ДОДАНО: Поля для бази даних
	MaintenanceIntervalKm   int `json:"maintenance_interval_km"`
	LastMaintenanceOdometer int `json:"last_maintenance_odometer"`

	// ДОДАНО: Обчислювані поля (заповнюються в Репозиторії на льоту)
	CurrentOdometer     int    `json:"current_odometer"`
	KmToNextMaintenance int    `json:"km_to_next_maintenance"`
	MaintenanceStatus   string `json:"maintenance_status"` // OK, WARNING, OVERDUE

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FuelRecord struct {
	ID            string         `json:"id"`
	VehicleID     string         `json:"vehicle_id"`
	Liters        float64        `json:"liters"`
	OdometerKm    *int           `json:"odometer_km"`
	RecordType    FuelRecordType `json:"record_type"`
	IsAnomaly     bool           `json:"is_anomaly"`
	AnomalyReason *string        `json:"anomaly_reason,omitempty"`
	CreatedBy     *string        `json:"created_by"`
	CreatedAt     time.Time      `json:"created_at"`
}

type MaintenanceRecord struct {
	ID          string    `json:"id"`
	VehicleID   string    `json:"vehicle_id"`
	OdometerKm  int       `json:"odometer_km"`
	Description string    `json:"description"`
	PerformedBy string    `json:"performed_by,omitempty"`
	CostAmount  float64   `json:"cost_amount"`
	DocumentURL string    `json:"document_url,omitempty"`
	DriverName  *string   `json:"driver_name"`
	CreatedAt   time.Time `json:"created_at"`
}

type VehicleStatusUpdate struct {
	Status string `json:"status" binding:"required"`
	Reason string `json:"reason"`
}

type VehicleDriverHistory struct {
	ID         string    `json:"id"`
	VehicleID  string    `json:"vehicle_id"`
	DriverID   *string   `json:"driver_id"`
	DriverName string    `json:"driver_name"`
	AssignedAt time.Time `json:"assigned_at"`
}

var FuelRecordCreatorRoles = []UserRole{
	RoleAdmin,
	RoleBrigadeCmdr,
	RoleBattalionCmdr,
	RoleCompanyCmdr,
	RolePlatoonCmdr,
	RoleBrigadeLogist,
	RoleBattalionLogist,
	RoleBrigadeStorekeeper,
	RoleBattalionStorekeeper,
	RoleCompanySergeant,
}

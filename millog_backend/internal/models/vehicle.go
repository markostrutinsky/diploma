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

type VehicleType string

const (
	VehicleTypePickup VehicleType = "PICKUP"
	VehicleTypeVan    VehicleType = "VAN"
	VehicleTypeTruck  VehicleType = "TRUCK"
)

type Vehicle struct {
	ID                      string        `json:"id"`
	Brand                   string        `json:"brand"`
	Model                   string        `json:"model"`
	PlateNumber             string        `json:"plate_number"`
	Type                    string        `json:"type"`        // ТИП
	CapacityKg              float64       `json:"capacity_kg"` // ВАГА
	Status                  VehicleStatus `json:"status"`
	DriverID                *string       `json:"driver_id"`
	DriverName              *string       `json:"driver_name"`            // ДОДАЙ ЦЕ ПОЛЕ!
	HomeWarehouseID         *string       `json:"home_warehouse_id"`      // ПОСТІЙНА БАЗА машини
	HomeWarehouseName       *string       `json:"home_warehouse_name"`    // назва для UI
	CurrentWarehouseID      *string       `json:"current_warehouse_id"`   // ПОТОЧНА ЛОКАЦІЯ (змінюється після рейсів)
	CurrentWarehouseName    *string       `json:"current_warehouse_name"` // назва для UI
	TankCapacity            float64       `json:"tank_capacity"`
	FuelNorm                float64       `json:"fuel_norm"`
	MaintenanceIntervalKm   int           `json:"maintenance_interval_km"`
	LastMaintenanceOdometer int           `json:"last_maintenance_odometer"`
	CurrentOdometer         int           `json:"current_odometer"`
	KmToNextMaintenance     int           `json:"km_to_next_maintenance"`
	MaintenanceStatus       string        `json:"maintenance_status"`
	AvgKmPerDay             float64       `json:"avg_km_per_day"`
	PredictedMaintDate      *time.Time    `json:"predicted_maint_date"`
	CurrentFuelLiters       float64       `json:"current_fuel_liters"` // Поточний залишок пального
	CreatedAt               time.Time     `json:"created_at"`
	UpdatedAt               time.Time     `json:"updated_at"`
}

type CreateVehicleRequest struct {
	Brand           string  `json:"brand" binding:"required"`
	Model           string  `json:"model"`
	PlateNumber     string  `json:"plate_number" binding:"required"`
	Type            string  `json:"type" binding:"required"`
	CapacityKg      float64 `json:"capacity_kg" binding:"required,gt=0"`
	TankCapacity    float64 `json:"tank_capacity" binding:"required,gt=0"`
	FuelNorm        float64 `json:"fuel_norm" binding:"required,gt=0"`
	DriverID        *string `json:"driver_id"`
	HomeWarehouseID *string `json:"home_warehouse_id"` // Обов'язковий базовий склад
}

type FuelRecord struct {
	ID            string         `json:"id"`
	VehicleID     string         `json:"vehicle_id"`
	Liters        float64        `json:"liters"`
	OdometerKm    *int           `json:"odometer_km"`
	RecordType    FuelRecordType `json:"record_type"`
	IsAnomaly     bool           `json:"is_anomaly"`
	AnomalyReason *string        `json:"anomaly_reason,omitempty"`
	// AnomalyExcessLiters — обсяг «зайвого» пального саме цього запису
	// (перевитрата понад норму або весь обсяг при витраті без руху).
	// Використовується антифрод-системою для розрахунку грошових втрат.
	AnomalyExcessLiters float64   `json:"anomaly_excess_liters"`
	CreatedBy           *string   `json:"created_by"`
	CreatedAt           time.Time `json:"created_at"`
}

type MaintenanceRecord struct {
	ID           string     `json:"id"`
	VehicleID    string     `json:"vehicle_id"`
	OdometerKm   int        `json:"odometer_km"`
	Description  string     `json:"description"`
	PerformedBy  string     `json:"performed_by,omitempty"`
	CostAmount   float64    `json:"cost_amount"`
	DocumentURL  string     `json:"document_url,omitempty"`
	DriverName   *string    `json:"driver_name"`
	Status       string     `json:"status"`
	ServiceType  string     `json:"service_type"`
	ScheduledFor *time.Time `json:"scheduled_for,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
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

type InventoryItem struct {
	ID          string    `json:"id"`
	WarehouseID string    `json:"warehouse_id"` // На якому складі це лежить
	Name        string    `json:"name"`         // "Квадрокоптер DJI Mavic 3T"
	Category    string    `json:"category"`     // "Техніка", "Медицина"
	Available   int       `json:"available"`    // Кількість (залишок)
	WeightKg    float64   `json:"weight_kg"`    // Вага 1 одиниці (для розрахунку перевантаження авто)
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Shipment (Накладна/Рейс) описує сам факт відправки машини з товаром
type Shipment struct {
	ID              string     `json:"id"`
	FromWarehouseID string     `json:"from_warehouse_id"`
	ToWarehouseID   string     `json:"to_warehouse_id"`
	VehicleID       string     `json:"vehicle_id"` // Яка машина поїхала
	Priority        string     `json:"priority"`   // "NORMAL" або "URGENT"
	Status          string     `json:"status"`     // "PENDING", "IN_TRANSIT", "DELIVERED", "CANCELLED"
	Direction       string     `json:"direction"`  // "DOWNSTREAM" (розподіл), "UPSTREAM" (консолідація/повернення), "LATERAL" (в межах підрозділу)
	CreatedAt       time.Time  `json:"created_at"`
	StartedAt       *time.Time `json:"started_at,omitempty"`   // Коли водій почав рейс
	DeliveredAt     *time.Time `json:"delivered_at,omitempty"` // Коли доставлено
}

// ShipmentItem (Вміст рейсу) описує, що саме поклали в машину
type ShipmentItem struct {
	ID          string `json:"id"`
	ShipmentID  string `json:"shipment_id"`
	InventoryID string `json:"inventory_id"` // Посилання на товар
	Quantity    int    `json:"quantity"`     // Скільки штук завантажили
}

type ShipmentRecord struct {
	ID            string    `json:"id"`
	FromWarehouse string    `json:"from_warehouse"`
	ToWarehouse   string    `json:"to_warehouse"`
	Vehicle       string    `json:"vehicle"`
	Priority      string    `json:"priority"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type IssueResourceRequest struct {
	ResourceID  string `json:"resource_id" binding:"required"`
	UserID      string `json:"user_id" binding:"required"` // Кому видаємо
	Quantity    int    `json:"quantity" binding:"required,gt=0"`
	Notes       string `json:"notes"`
	WarehouseID string `json:"warehouse_id"` // Явний склад (для адмін-ролей без unit_id)
}

type ShipmentItemRequest struct {
	ResourceID string  `json:"resource_id"`
	Quantity   int     `json:"quantity"`
	RequestID  *string `json:"request_id,omitempty"`
}

type CreateShipmentRequest struct {
	FromWarehouseID string                `json:"from_warehouse_id" binding:"required"`
	ToWarehouseID   string                `json:"to_warehouse_id" binding:"required"`
	VehicleID       string                `json:"vehicle_id" binding:"required"`
	Priority        string                `json:"priority" binding:"required"`
	Items           []ShipmentItemRequest `json:"items" binding:"required,min=1"`
	DistanceKm      float64               `json:"distance_km"` // Планова відстань маршруту (OSRM)
}

var FuelRecordCreatorRoles = []UserRole{
	RoleSystemAdmin,
	RoleTenantAdmin,
	RoleAdmin,
	RoleRegionDirector,
	RoleBranchManager,
	RoleDeptManager,
	RoleTeamLead,
	RoleRegionLogistician,
	RoleBranchLogistician,
	RoleRegionStorekeeper,
	RoleBranchStorekeeper,
	RoleDeptSupervisor,
}

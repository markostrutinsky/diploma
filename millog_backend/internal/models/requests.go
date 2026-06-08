package models

import "time"

type RequestStatus string

const (
	RequestPending    RequestStatus = "PENDING"
	RequestApproved   RequestStatus = "APPROVED"
	RequestLoading    RequestStatus = "LOADING" // Рейс сформовано, очікує завантаження в авто
	RequestDispatched RequestStatus = "DISPATCHED"
	RequestRejected   RequestStatus = "REJECTED"
	RequestCompleted  RequestStatus = "COMPLETED"
	RequestOpen       RequestStatus = "OPEN"
	RequestEscalated  RequestStatus = "ESCALATED"
)

type SupplyRequest struct {
	ID                 string        `json:"id"`
	CreatedBy          string        `json:"created_by"`
	ResourceID         *string       `json:"resource_id"` // Тепер nullable - може бути без прив'язки до конкретного складу
	ResourceName       string        `json:"resource_name"`
	ResourceCategoryID *string       `json:"resource_category_id"`
	Quantity           int           `json:"quantity"`
	Status             RequestStatus `json:"status"`
	TargetWarehouseID  string        `json:"target_warehouse_id"`
	ApprovedBy         *string       `json:"approved_by"`
	ApprovedAt         *time.Time    `json:"approved_at"`
	Comment            *string       `json:"comment"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`

	ManagerEmail string `json:"-"` // Додано для зручності при відправці SLA-сповіщень
}

type CreateSupplyRequest struct {
	ResourceID         *string `json:"resource_id"` // Тепер опціонально
	ResourceName       string  `json:"resource_name" binding:"required"`
	ResourceCategoryID *string `json:"resource_category_id"`
	Quantity           int     `json:"quantity" binding:"required,min=1"`
	TargetWarehouseID  string  `json:"target_warehouse_id" binding:"required,uuid"` // <-- ДОДАЙ ВАЛІДАЦІЮ UUID
}
type ApproveRequest struct {
	Approved bool   `json:"approved"`
	Comment  string `json:"comment"`
}

type SmartDispatchReq struct {
	RequestIDs      []string `json:"request_ids" binding:"required,min=1"`
	FromWarehouseID string   `json:"from_warehouse_id"`
}

type RequestItem struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	WeightKg          float64 `json:"weight_kg"`
	TargetWarehouseID string  `json:"target_warehouse_id,omitempty"`
}

type VehicleBin struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	MaxWeight    float64       `json:"max_weight"`
	UsedWeight   float64       `json:"used_weight"`
	Items        []RequestItem `json:"items"`
	FuelLiters   float64       `json:"fuel_liters"`
	FuelNorm     float64       `json:"fuel_norm"`
	TankCapacity float64       `json:"tank_capacity"`
}

type SmartDispatchResult struct {
	OptimizedRoutes []VehicleBin  `json:"routes"`
	Unassigned      []RequestItem `json:"unassigned"`
}

// SmartDispatchRoute — один рядок з "затвердження" результату Smart Розподілу:
// фронт передає, які саме заявки поїдуть якою машиною.
type SmartDispatchRoute struct {
	VehicleID  string   `json:"vehicle_id" binding:"required,uuid"`
	RequestIDs []string `json:"request_ids" binding:"required,min=1"`
}

// SmartDispatchConfirmReq — payload на /requests/smart-dispatch-confirm.
// Цільовий склад виводиться з target_warehouse_id заявок (усі мають збігатись),
// тому тут вказуємо лише звідки відправляти й розклад по машинам.
type SmartDispatchConfirmReq struct {
	FromWarehouseID string               `json:"from_warehouse_id" binding:"required,uuid"`
	Priority        string               `json:"priority"`
	DistanceKm      float64              `json:"distance_km"` // необов'язкове — якщо 0, розраховується на бекенді
	Routes          []SmartDispatchRoute `json:"routes" binding:"required,min=1"`
}

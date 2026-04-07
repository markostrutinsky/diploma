package models

import "time"

type ResourceCondition string

const (
	ConditionNew        ResourceCondition = "NEW"
	ConditionUsed       ResourceCondition = "USED"
	ConditionWrittenOff ResourceCondition = "WRITTEN_OFF"
)

type UnitMeasurement string

const (
	UnitPCS UnitMeasurement = "PCS"
	UnitKIT UnitMeasurement = "KIT"
	UnitKG  UnitMeasurement = "KG"
	UnitL   UnitMeasurement = "L"
)

type Warehouse struct {
	ID           string    `json:"id"`
	UnitID       int64     `json:"unit_id"`
	Name         string    `json:"name"`
	LocationType string    `json:"location_type"`
	Latitude     *float64  `json:"latitude"`
	Longitude    *float64  `json:"longitude"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateWarehouseRequest struct {
	UnitID       int64    `json:"unit_id" binding:"required"`
	Name         string   `json:"name" binding:"required"`
	LocationType string   `json:"location_type" binding:"required,oneof=STATIONARY MOBILE"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
}

type ResourceCategory struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Resource struct {
	ID                 string            `json:"id"`
	CategoryID         string            `json:"category_id"`
	UnitID             int64             `json:"unit_id"`
	WarehouseID        *string           `json:"warehouse_id"`
	Name               string            `json:"name"`
	Description        string            `json:"description"`
	Quantity           int               `json:"quantity"`
	UnitType           UnitMeasurement   `json:"unit_type"`
	SerialNumber       string            `json:"serial_number"`
	Condition          ResourceCondition `json:"condition"`
	MinQuantity        int               `json:"min_quantity"`
	AssignedToUserID   *string           `json:"assigned_to_user_id"`
	AssignedToUserName *string           `json:"assigned_to_user_name,omitempty" gorm:"-"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreateResourceRequest struct {
	CategoryID   string            `json:"category_id" binding:"required"`
	UnitID       int64             `json:"unit_id" binding:"required"`
	WarehouseID  *string           `json:"warehouse_id"`
	Name         string            `json:"name" binding:"required"`
	Description  string            `json:"description"`
	Quantity     int               `json:"quantity"`
	UnitType     UnitMeasurement   `json:"unit_type" binding:"required,oneof=PCS KIT KG L"`
	SerialNumber string            `json:"serial_number"`
	Condition    ResourceCondition `json:"condition" binding:"omitempty,oneof=NEW USED WRITTEN_OFF"`
	MinQuantity  int               `json:"min_quantity"`
}

type UpdateResourceRequest struct {
	Name         *string            `json:"name"`
	Description  *string            `json:"description"`
	Quantity     *int               `json:"quantity"`
	WarehouseID  *string            `json:"warehouse_id"`
	SerialNumber *string            `json:"serial_number"`
	Condition    *ResourceCondition `json:"condition"`
	MinQuantity  *int               `json:"min_quantity"`
}

type TransferResourceRequest struct {
	Quantity          int     `json:"quantity" binding:"required,gt=0"` // Скільки переміщуємо
	TargetWarehouseID *string `json:"target_warehouse_id"`              // Куди (може бути nil, якщо видаємо на руки)
	TargetUnitID      *int64  `json:"target_unit_id"`                   // Кому (якщо передаємо іншому підрозділу)
}

type WriteOffResourceRequest struct {
	Quantity int `json:"quantity" binding:"required,gt=0"`
}

type AssignResourceRequest struct {
	Quantity int    `json:"quantity" binding:"required,gt=0"`
	UserID   string `json:"user_id" binding:"required"`
}

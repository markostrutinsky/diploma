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
	ID             string            `json:"id"`
	CategoryID     string            `json:"category_id"`
	UnitID         int64             `json:"unit_id"`
	WarehouseID    *string           `json:"warehouse_id"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Quantity       int               `json:"quantity"`
	UnitType       UnitMeasurement   `json:"unit_type"`
	SerialNumber   string            `json:"serial_number"`
	Condition      ResourceCondition `json:"condition"`
	MinQuantity    int               `json:"min_quantity"`
	IssuedQuantity int               `json:"issued_quantity"`
	WeightKg       float64           `json:"weight_kg"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
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
	WeightKg     float64           `json:"weight_kg"`
}

type UpdateResourceRequest struct {
	Name         *string            `json:"name"`
	Description  *string            `json:"description"`
	Quantity     *int               `json:"quantity"`
	WarehouseID  *string            `json:"warehouse_id"`
	SerialNumber *string            `json:"serial_number"`
	Condition    *ResourceCondition `json:"condition"`
	MinQuantity  *int               `json:"min_quantity"`
	WeightKg     *float64           `json:"weight_kg"`
}

type WriteOffResourceRequest struct {
	Quantity int `json:"quantity" binding:"required,gt=0"`
}

type AssignResourceRequest struct {
	Quantity int    `json:"quantity" binding:"required,gt=0"`
	UserID   string `json:"user_id" binding:"required"`
}

type AssignmentStatus string

const (
	AssignmentActive     AssignmentStatus = "ACTIVE"
	AssignmentReturned   AssignmentStatus = "RETURNED"
	AssignmentLost       AssignmentStatus = "LOST"
	AssignmentWrittenOff AssignmentStatus = "WRITTEN_OFF"
)

type ResourceAssignment struct {
	ID         string           `json:"id"`
	ResourceID string           `json:"resource_id"`
	UserID     string           `json:"user_id"`
	UserName   string           `json:"user_name,omitempty"`
	Quantity   int              `json:"quantity"`
	Status     AssignmentStatus `json:"status"`
	IssuedAt   time.Time        `json:"issued_at"`
	ReturnedAt *time.Time       `json:"returned_at,omitempty"`
	Notes      *string          `json:"notes,omitempty"`
}

type ResourceListItem struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	WarehouseStock int    `json:"warehouse_stock"` // Скільки лежить на складі (quantity з resources)
	IssuedQuantity int    `json:"issued_quantity"` // Скільки на руках (сума з assignments)
	MinQuantity    int    `json:"min_quantity"`
	Condition      string `json:"condition"`
}

// MyEquipmentItem відображає одну позицію майна в особистому кабінеті користувача
type MyEquipmentItem struct {
	AssignmentID string    `json:"assignment_id"` // ID самої транзакції видачі (знадобиться для рапортів)
	ResourceID   string    `json:"resource_id"`
	ResourceName string    `json:"resource_name"`
	Quantity     int       `json:"quantity"`
	UnitType     string    `json:"unit_type"`
	IssuedAt     time.Time `json:"issued_at"`
	Status       string    `json:"status"` // 'ACTIVE' (на руках)
}

// Модель для конкретної розбіжності
type AuditDiscrepancy struct {
	ResourceID     string `json:"resource_id" binding:"required"`
	BookQuantity   int    `json:"book_quantity"`
	ActualQuantity int    `json:"actual_quantity"`
	Difference     int    `json:"difference"`
}

// Загальний запит на збереження результатів інвентаризації
type SubmitAuditRequest struct {
	WarehouseID   string             `json:"warehouse_id" binding:"required"`
	Discrepancies []AuditDiscrepancy `json:"discrepancies"`
}

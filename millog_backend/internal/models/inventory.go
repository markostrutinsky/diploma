package models

import "time"

type ResourceCondition string

const (
	ConditionNew       ResourceCondition = "NEW"
	ConditionUsed      ResourceCondition = "USED"
	ConditionWrittenOff ResourceCondition = "WRITTEN_OFF"
)

type ResourceCategory struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Resource struct {
	ID           string             `json:"id"`
	CategoryID   string             `json:"category_id"`
	UnitID       *int64             `json:"unit_id"` // склад (бригада/батальйон/рота)
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Quantity     int                `json:"quantity"`
	SerialNumber string             `json:"serial_number"`
	Location     string             `json:"location"`
	Condition    ResourceCondition  `json:"condition"`
	MinQuantity  int                `json:"min_quantity"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreateResourceRequest struct {
	CategoryID   string            `json:"category_id" binding:"required"`
	UnitID       *int64            `json:"unit_id"` // склад
	Name         string            `json:"name" binding:"required"`
	Description  string            `json:"description"`
	Quantity     int               `json:"quantity"`
	SerialNumber string            `json:"serial_number"`
	Location     string            `json:"location"`
	Condition    ResourceCondition `json:"condition" binding:"omitempty,oneof=NEW USED WRITTEN_OFF"`
	MinQuantity  int               `json:"min_quantity"`
}

type UpdateResourceRequest struct {
	Name         *string            `json:"name"`
	Description  *string            `json:"description"`
	Quantity     *int               `json:"quantity"`
	SerialNumber *string            `json:"serial_number"`
	Location     *string            `json:"location"`
	Condition    *ResourceCondition `json:"condition"`
	MinQuantity  *int               `json:"min_quantity"`
}

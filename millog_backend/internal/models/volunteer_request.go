package models

import "time"

type ContractorRequestStatus string

const (
	ContractorOpen      ContractorRequestStatus = "OPEN"
	ContractorTaken     ContractorRequestStatus = "TAKEN"
	ContractorDelivered ContractorRequestStatus = "DELIVERED"
	ContractorRejected  ContractorRequestStatus = "REJECTED"
	ContractorAccepted  ContractorRequestStatus = "ACCEPTED"
	ContractorCanceled  ContractorRequestStatus = "CANCELED"
)

type ContractorRequest struct {
	ID                string                  `json:"id"`
	CreatedBy         string                  `json:"created_by"`
	UnitID            *int64                  `json:"unit_id"`
	UnitName          *string                 `json:"unit_name"`
	TargetWarehouseID *string                 `json:"target_warehouse_id"`
	WarehouseName     *string                 `json:"warehouse_name"`
	Title             string                  `json:"title"`
	Description       string                  `json:"description"`
	Status            ContractorRequestStatus `json:"status"`
	TakenBy           *string                 `json:"taken_by"`
	TakenAt           *time.Time              `json:"taken_at"`
	CompletedAt       *time.Time              `json:"completed_at"`
	CreatedAt         time.Time               `json:"created_at"`

	// Організація-замовник. Для крос-tenant дошки підрядника важливо показувати,
	// від якої організації завдання (і чи має підрядник право його взяти).
	TenantID   *string `json:"tenant_id,omitempty"`
	TenantName *string `json:"tenant_name,omitempty"`
}

type CreateContractorRequest struct {
	Title             string  `json:"title" binding:"required"`
	Description       string  `json:"description"`
	UnitID            *int64  `json:"unit_id"`
	TargetWarehouseID *string `json:"target_warehouse_id"`
}

type AcceptContractorPayload struct {
	ResourceID   *string         `json:"resource_id"`
	CategoryID   string          `json:"category_id"`
	CategoryName string          `json:"category_name"`
	Name         string          `json:"name"`
	Quantity     int             `json:"quantity" binding:"required,min=1"`
	UnitType     UnitMeasurement `json:"unit_type" binding:"required,oneof=PCS KIT KG L"`
	UnitPrice    float64         `json:"unit_price" binding:"omitempty,gte=0"`
}

package models

import "time"

type CONTRACTORRequestStatus string

const (
	CONTRACTOROpen      CONTRACTORRequestStatus = "OPEN"
	CONTRACTORTaken     CONTRACTORRequestStatus = "TAKEN"
	CONTRACTORDelivered CONTRACTORRequestStatus = "DELIVERED"
	CONTRACTORRejected  CONTRACTORRequestStatus = "REJECTED"
	CONTRACTORAccepted  CONTRACTORRequestStatus = "ACCEPTED"
	CONTRACTORCanceled  CONTRACTORRequestStatus = "CANCELED"
)

type CONTRACTORRequest struct {
	ID          string                  `json:"id"`
	CreatedBy   string                  `json:"created_by"`
	UnitID      *int64                  `json:"unit_id"`
	UnitName    *string                 `json:"unit_name"`
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	Status      CONTRACTORRequestStatus `json:"status"`
	TakenBy     *string                 `json:"taken_by"`
	TakenAt     *time.Time              `json:"taken_at"`
	CompletedAt *time.Time              `json:"completed_at"`
	CreatedAt   time.Time               `json:"created_at"`
}

type CreateCONTRACTORRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	UnitID      *int64 `json:"unit_id"`
}

type AcceptCONTRACTORPayload struct {
	ResourceID *string         `json:"resource_id"`
	CategoryID string          `json:"category_id" binding:"required"`
	Name       string          `json:"name" binding:"required"`
	Quantity   int             `json:"quantity" binding:"required,min=1"`
	UnitType   UnitMeasurement `json:"unit_type" binding:"required,oneof=PCS KIT KG L"`
}

package models

import "time"

type VolunteerRequestStatus string

const (
	VolunteerOpen      VolunteerRequestStatus = "OPEN"
	VolunteerTaken     VolunteerRequestStatus = "TAKEN"
	VolunteerDelivered VolunteerRequestStatus = "DELIVERED"
	VolunteerRejected  VolunteerRequestStatus = "REJECTED"
	VolunteerAccepted  VolunteerRequestStatus = "ACCEPTED"
	VolunteerCanceled  VolunteerRequestStatus = "CANCELED"
)

type VolunteerRequest struct {
	ID          string                 `json:"id"`
	CreatedBy   string                 `json:"created_by"`
	UnitID      *int64                 `json:"unit_id"`
	UnitName    *string                `json:"unit_name"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Status      VolunteerRequestStatus `json:"status"`
	TakenBy     *string                `json:"taken_by"`
	TakenAt     *time.Time             `json:"taken_at"`
	CompletedAt *time.Time             `json:"completed_at"`
	CreatedAt   time.Time              `json:"created_at"`
}

type CreateVolunteerRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	UnitID      *int64 `json:"unit_id"`
}

type AcceptVolunteerPayload struct {
	ResourceID *string         `json:"resource_id"`
	CategoryID string          `json:"category_id" binding:"required"`
	Name       string          `json:"name" binding:"required"`
	Quantity   int             `json:"quantity" binding:"required,min=1"`
	UnitType   UnitMeasurement `json:"unit_type" binding:"required,oneof=PCS KIT KG L"`
}

package models

import "time"

type VolunteerRequestStatus string

const (
	VolunteerOpen      VolunteerRequestStatus = "OPEN"
	VolunteerTaken     VolunteerRequestStatus = "TAKEN"
	VolunteerCompleted VolunteerRequestStatus = "COMPLETED"
)

type VolunteerRequest struct {
	ID          string                 `json:"id"`
	CreatedBy   string                 `json:"created_by"`
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
}

package models

import "time"

type RequestStatus string

const (
	RequestPending  RequestStatus = "PENDING"
	RequestApproved RequestStatus = "APPROVED"
	RequestRejected RequestStatus = "REJECTED"
)

type SupplyRequest struct {
	ID         string         `json:"id"`
	CreatedBy  string         `json:"created_by"`
	ResourceID string         `json:"resource_id"`
	Quantity   int            `json:"quantity"`
	Status     RequestStatus  `json:"status"`
	ApprovedBy *string        `json:"approved_by"`
	ApprovedAt *time.Time     `json:"approved_at"`
	Comment    string         `json:"comment"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type CreateSupplyRequest struct {
	ResourceID string `json:"resource_id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
}

type ApproveRequest struct {
	Approved bool   `json:"approved"`
	Comment  string `json:"comment"`
}

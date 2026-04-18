package models

import "time"

type RequestStatus string

const (
	RequestPending    RequestStatus = "PENDING"
	RequestApproved   RequestStatus = "APPROVED"
	RequestDispatched RequestStatus = "DISPATCHED"
	RequestRejected   RequestStatus = "REJECTED"
	RequestCompleted  RequestStatus = "COMPLETED"
	RequestOpen       RequestStatus = "OPEN"
)

type SupplyRequest struct {
	ID                string        `json:"id"`
	CreatedBy         string        `json:"created_by"`
	ResourceID        string        `json:"resource_id"`
	Quantity          int           `json:"quantity"`
	Status            RequestStatus `json:"status"`
	TargetWarehouseID string        `json:"target_warehouse_id"`
	ApprovedBy        *string       `json:"approved_by"`
	ApprovedAt        *time.Time    `json:"approved_at"`
	Comment           *string       `json:"comment"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}

type CreateSupplyRequest struct {
	ResourceID        string `json:"resource_id" binding:"required"`
	Quantity          int    `json:"quantity" binding:"required,min=1"`
	TargetWarehouseID string `json:"target_warehouse_id" binding:"required,uuid"` // <-- ДОДАЙ ВАЛІДАЦІЮ UUID
}
type ApproveRequest struct {
	Approved bool   `json:"approved"`
	Comment  string `json:"comment"`
}

type SmartDispatchReq struct {
	RequestIDs []string `json:"request_ids" binding:"required,min=1"`
}

type RequestItem struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	WeightKg float64 `json:"weight_kg"`
}

type VehicleBin struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	MaxWeight  float64       `json:"max_weight"`
	UsedWeight float64       `json:"used_weight"`
	Items      []RequestItem `json:"items"`
}

type SmartDispatchResult struct {
	OptimizedRoutes []VehicleBin  `json:"routes"`
	Unassigned      []RequestItem `json:"unassigned"`
}

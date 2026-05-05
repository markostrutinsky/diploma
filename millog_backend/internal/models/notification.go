package models

import "time"

type NotificationType string

const (
	NotificationShipmentAssigned  NotificationType = "SHIPMENT_ASSIGNED"
	NotificationRequestApproved   NotificationType = "REQUEST_APPROVED"
	NotificationRequestRejected   NotificationType = "REQUEST_REJECTED"
	NotificationShipmentDelivered NotificationType = "SHIPMENT_DELIVERED"
	NotificationLowStock          NotificationType = "LOW_STOCK"
)

type Notification struct {
	ID        string           `json:"id"`
	UserID    string           `json:"user_id"`
	TenantID  string           `json:"tenant_id"`
	Type      NotificationType `json:"type"`
	Title     string           `json:"title"`
	Message   string           `json:"message"`
	RelatedID *string          `json:"related_id,omitempty"`
	IsRead    bool             `json:"is_read"`
	CreatedAt time.Time        `json:"created_at"`
	ReadAt    *time.Time       `json:"read_at,omitempty"`
}

type CreateNotificationRequest struct {
	UserID    string           `json:"user_id" binding:"required"`
	Type      NotificationType `json:"type" binding:"required"`
	Title     string           `json:"title" binding:"required"`
	Message   string           `json:"message" binding:"required"`
	RelatedID *string          `json:"related_id,omitempty"`
}

type NotificationListResponse struct {
	Notifications []Notification `json:"notifications"`
	UnreadCount   int            `json:"unread_count"`
}

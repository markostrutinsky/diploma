package models

import "time"

type AuditLog struct {
	ID         int       `json:"id"`
	UserEmail  string    `json:"user_email"`
	UserRole   string    `json:"user_role"`
	ActionType string    `json:"action_type"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	Details    string    `json:"details"`
	CreatedAt  time.Time `json:"created_at"`
}

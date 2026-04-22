package models

import "time"

// SubscriptionTier — тариф організації (tenant).
type SubscriptionTier string

const (
	TierFree       SubscriptionTier = "FREE"
	TierBasic      SubscriptionTier = "BASIC"
	TierPro        SubscriptionTier = "PRO"
	TierEnterprise SubscriptionTier = "ENTERPRISE"
)

// Tenant — окрема організація-клієнт платформи (Нова пошта, військова частина тощо).
type Tenant struct {
	ID                    string           `json:"id"`
	Name                  string           `json:"name"`
	Slug                  string           `json:"slug"`
	SubscriptionTier      SubscriptionTier `json:"subscription_tier"`
	SubscriptionExpiresAt *time.Time       `json:"subscription_expires_at,omitempty"`
	OwnerEmail            *string          `json:"owner_email,omitempty"`
	IsActive              bool             `json:"is_active"`
	CreatedAt             time.Time        `json:"created_at"`
	UpdatedAt             time.Time        `json:"updated_at"`
}

// CreateTenantRequest — запит на створення нової організації (self-service signup).
type CreateTenantRequest struct {
	OrganizationName string `json:"organization_name" binding:"required,min=2"`
	Slug             string `json:"slug" binding:"required,min=2,max=100"`
	OwnerEmail       string `json:"owner_email" binding:"required,email"`
	OwnerFullName    string `json:"owner_full_name" binding:"required"`
	OwnerPassword    string `json:"owner_password" binding:"required,min=8"`
}

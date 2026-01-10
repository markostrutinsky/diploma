package models

import "time"

type InviteToken struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string    `gorm:"not null;unique"`
	User      User      `gorm:"constraint:OnDelete:CASCADE"`
	TokenHash string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	UsedAt    *time.Time
	CreatedAt time.Time
}

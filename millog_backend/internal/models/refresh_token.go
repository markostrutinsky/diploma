package models

import "time"

type RefreshToken struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string    `gorm:"not null"`
	TokenHash string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	RevokedAt *time.Time
	CreatedAt time.Time
}

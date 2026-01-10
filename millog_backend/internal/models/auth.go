package models

import (
	"time"
)

type UserRole string

const (
	RoleAdmin     UserRole = "ADMIN"
	RoleWarehouse UserRole = "WAREHOUSE"
	RoleCommander UserRole = "COMMANDER"
	RoleVolunteer UserRole = "VOLUNTEER"
)

type UserStatus string

const (
	StatusPending UserStatus = "PENDING"
	StatusActive  UserStatus = "ACTIVE"
	StatusBlocked UserStatus = "BLOCKED"
)

type User struct {
	ID           string     `json:"id" gorm:"primaryKey"`
	Username     *string    `json:"username" gorm:"unique;not null;size:100"`
	Email        string     `json:"email" gorm:"unique;not null;size:255"`
	FullName     string     `json:"full_name" gorm:"size:255"`
	Phone        *string    `json:"phone" gorm:"size:50"`
	PasswordHash *string    `json:"-"`
	Role         UserRole   `json:"role" gorm:"type:varchar(20);default:'VOLUNTEER'"`
	Status       UserStatus `json:"status" gorm:"type:varchar(20);default:'PENDING'"`
	UnitID       *uint      `json:"unit_id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type CreateUserRequest struct {
	Username *string `json:"username"`
	Email    string  `json:"email" binding:"required,email"`
	FullName string  `json:"full_name" binding:"required"`
	Phone    *string `json:"phone"`
	Password *string `json:"password"`
	Role     string  `json:"role" binding:"required,oneof=ADMIN WAREHOUSE COMMANDER VOLUNTEER"`
	UnitID   *uint   `json:"unit_id"`
}

type CreateUserResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

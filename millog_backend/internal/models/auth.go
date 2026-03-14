package models

import (
	"time"
)

type UserRole string

const (
	RoleAdmin               UserRole = "ADMIN"
	RoleBrigadeCmdr         UserRole = "BRIGADE_CMDR"
	RoleBattalionCmdr       UserRole = "BATTALION_CMDR"
	RoleCompanyCmdr         UserRole = "COMPANY_CMDR"
	RolePlatoonCmdr         UserRole = "PLATOON_CMDR"
	RoleBrigadeLogist       UserRole = "BRIGADE_LOGIST"
	RoleBrigadeStorekeeper  UserRole = "BRIGADE_STOREKEEPER"
	RoleBattalionLogist     UserRole = "BATTALION_LOGIST"
	RoleBattalionStorekeeper UserRole = "BATTALION_STOREKEEPER"
	RoleCompanySergeant     UserRole = "COMPANY_SERGEANT"
	RoleVolunteer           UserRole = "VOLUNTEER"
)

// Roles that can create supply requests
var SupplyRequestCreatorRoles = []UserRole{
	RoleAdmin, RoleBrigadeCmdr, RoleBattalionCmdr, RoleCompanyCmdr, RolePlatoonCmdr,
	RoleBrigadeLogist, RoleBattalionLogist, RoleCompanySergeant,
}

// Roles that can approve supply requests
var SupplyRequestApproverRoles = []UserRole{
	RoleAdmin, RoleBrigadeCmdr, RoleBattalionCmdr, RoleCompanyCmdr,
	RoleBrigadeLogist, RoleBattalionLogist,
}

// Roles that can create volunteer requests
var VolunteerRequestCreatorRoles = []UserRole{
	RoleAdmin, RoleBrigadeCmdr, RoleBattalionCmdr, RoleCompanyCmdr,
	RoleBrigadeLogist, RoleBattalionLogist, RoleBrigadeStorekeeper, RoleBattalionStorekeeper, RoleCompanySergeant,
}

// Roles that can manage inventory (warehouses)
var InventoryManagerRoles = []UserRole{
	RoleAdmin, RoleBrigadeStorekeeper, RoleBattalionStorekeeper, RoleCompanySergeant,
}

// Roles that can manage units
var UnitManagerRoles = []UserRole{
	RoleAdmin, RoleBrigadeCmdr, RoleBattalionCmdr, RoleCompanyCmdr,
	RoleBrigadeLogist, RoleBattalionLogist, RoleBrigadeStorekeeper, RoleBattalionStorekeeper,
}

// Roles that can create users (except VOLUNTEER)
var UserCreatorRoles = []UserRole{
	RoleAdmin, RoleBrigadeCmdr, RoleBattalionCmdr, RoleCompanyCmdr, RolePlatoonCmdr,
	RoleBrigadeLogist, RoleBattalionLogist, RoleBrigadeStorekeeper, RoleBattalionStorekeeper, RoleCompanySergeant,
}

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
	Role         UserRole   `json:"role" gorm:"type:varchar(30);default:'VOLUNTEER'"`
	Status       UserStatus `json:"status" gorm:"type:varchar(20);default:'PENDING'"`
	UnitID       *int64     `json:"unit_id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type CreateUserRequest struct {
	Username *string `json:"username"`
	Email    string  `json:"email" binding:"required,email"`
	FullName string  `json:"full_name" binding:"required"`
	Phone    *string `json:"phone"`
	Password *string `json:"password"`
	Role     string  `json:"role" binding:"required"`
	UnitID   *int64  `json:"unit_id"`
}

type CreateUserResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	User         User   `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type SetupPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

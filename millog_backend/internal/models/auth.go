package models

import (
	"time"
)

type UserRole string

const (
	RoleAdmin                UserRole = "ADMIN"
	RoleBrigadeCmdr          UserRole = "BRIGADE_CMDR"
	RoleBattalionCmdr        UserRole = "BATTALION_CMDR"
	RoleCompanyCmdr          UserRole = "COMPANY_CMDR"
	RolePlatoonCmdr          UserRole = "PLATOON_CMDR"
	RoleBrigadeLogist        UserRole = "BRIGADE_LOGIST"
	RoleBrigadeStorekeeper   UserRole = "BRIGADE_STOREKEEPER"
	RoleBattalionLogist      UserRole = "BATTALION_LOGIST"
	RoleBattalionStorekeeper UserRole = "BATTALION_STOREKEEPER"
	RoleCompanySergeant      UserRole = "COMPANY_SERGEANT"
	RoleVolunteer            UserRole = "VOLUNTEER"
)

var MilitaryInventoryRoles = []UserRole{
	RoleAdmin,
	RoleBrigadeCmdr,
	RoleBattalionCmdr,
	RoleCompanyCmdr,
	RoleBrigadeLogist,
	RoleBattalionLogist,
	RoleBrigadeStorekeeper,
	RoleBattalionStorekeeper,
	RoleCompanySergeant,
}

// Roles that can create and manage warehouses (infrastructure)
var WarehouseManagerRoles = []UserRole{
	RoleAdmin,
	RoleBrigadeCmdr,
	RoleBrigadeLogist,
	RoleBattalionCmdr,
	RoleBattalionLogist,
	RoleCompanyCmdr,
	RolePlatoonCmdr,
}

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

var RoleCreationMap = map[UserRole][]UserRole{
	RoleAdmin: {
		RoleBrigadeCmdr, RoleBattalionCmdr, RoleCompanyCmdr, RolePlatoonCmdr,
		RoleBrigadeLogist, RoleBrigadeStorekeeper, RoleBattalionLogist,
		RoleBattalionStorekeeper, RoleCompanySergeant,
	},
	RoleBrigadeCmdr: {
		RoleBattalionCmdr, RoleCompanyCmdr, RolePlatoonCmdr,
		RoleBrigadeLogist, RoleBrigadeStorekeeper, RoleBattalionLogist,
		RoleBattalionStorekeeper, RoleCompanySergeant,
	},
	RoleBattalionCmdr: {
		RoleCompanyCmdr, RolePlatoonCmdr,
		RoleBattalionLogist, RoleBattalionStorekeeper, RoleCompanySergeant,
	},
	RoleCompanyCmdr: {
		RolePlatoonCmdr, RoleCompanySergeant,
	},
	RoleBrigadeLogist: {
		RoleBrigadeStorekeeper, RoleBattalionLogist, RoleBattalionStorekeeper,
	},
	RoleBattalionLogist: {
		RoleBattalionStorekeeper,
	},
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

func (r UserRole) CanCreate(targetRole UserRole) bool {
	if r == RoleAdmin {
		return true
	}

	allowedRoles, exists := RoleCreationMap[r]
	if !exists {
		return false
	}

	for _, allowed := range allowedRoles {
		if targetRole == allowed {
			return true
		}
	}
	return false
}

func (r UserRole) GetTargetUnitType() string {
	switch r {
	case RoleBrigadeCmdr, RoleBrigadeLogist, RoleBrigadeStorekeeper:
		return "BRIGADE"
	case RoleBattalionCmdr, RoleBattalionLogist, RoleBattalionStorekeeper:
		return "BATTALION"
	case RoleCompanyCmdr, RoleCompanySergeant:
		return "COMPANY"
	case RolePlatoonCmdr:
		return "PLATOON"
	default:
		return ""
	}
}

func (r UserRole) CanCreateUnitType(unitType string) bool {
	if r == RoleAdmin {
		return true
	}

	switch r {
	case RoleBrigadeCmdr, RoleBrigadeLogist, RoleBrigadeStorekeeper:
		return unitType == "BATTALION" || unitType == "COMPANY" || unitType == "PLATOON"
	case RoleBattalionCmdr, RoleBattalionLogist, RoleBattalionStorekeeper:
		return unitType == "COMPANY" || unitType == "PLATOON"
	case RoleCompanyCmdr, RoleCompanySergeant:
		return unitType == "PLATOON"
	default:
		return false
	}
}

type UpdateRoleRequest struct {
	Role   string `json:"role" binding:"required"`
	UnitID *int64 `json:"unit_id"`
}

var ApprovalMatrix = map[UserRole][]UserRole{
	RolePlatoonCmdr: {
		RoleCompanyCmdr, RoleCompanySergeant,
		RoleBattalionCmdr,
		RoleBrigadeCmdr,
		RoleAdmin,
	},

	RoleCompanyCmdr: {
		RoleBattalionCmdr,
		RoleBrigadeCmdr,
		RoleAdmin,
	},
	RoleCompanySergeant: {
		RoleCompanyCmdr,
		RoleBattalionCmdr,
		RoleBrigadeCmdr,
		RoleAdmin,
	},

	RoleBattalionCmdr: {
		RoleBrigadeCmdr,
		RoleAdmin,
	},

	RoleBattalionLogist: {
		RoleBattalionCmdr,
		RoleBrigadeCmdr, RoleBrigadeLogist,
		RoleAdmin,
	},

	RoleBattalionStorekeeper: {
		RoleBattalionCmdr, RoleBattalionLogist,
		RoleBrigadeCmdr, RoleBrigadeLogist,
		RoleAdmin,
	},

	RoleBrigadeStorekeeper: {
		RoleBrigadeCmdr, RoleBrigadeLogist,
		RoleAdmin,
	},

	RoleBrigadeCmdr: {
		RoleAdmin,
	},
	RoleBrigadeLogist: {
		RoleBrigadeCmdr, RoleAdmin,
	},

	"USER": {
		RolePlatoonCmdr, RoleCompanyCmdr, RoleCompanySergeant,
		RoleBattalionCmdr,
		RoleBrigadeCmdr,
		RoleAdmin,
	},
}

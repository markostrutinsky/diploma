package models

import (
	"time"
)

type UserRole string

const (
	RoleAdmin UserRole = "ADMIN"

	RoleRegionDirector    UserRole = "REGION_DIRECTOR"
	RoleRegionLogistician UserRole = "REGION_LOGISTICIAN"
	RoleRegionStorekeeper UserRole = "REGION_STOREKEEPER"

	RoleBranchManager     UserRole = "BRANCH_MANAGER"
	RoleBranchLogistician UserRole = "BRANCH_LOGISTICIAN"
	RoleBranchStorekeeper UserRole = "BRANCH_STOREKEEPER"

	RoleDeptManager    UserRole = "DEPT_MANAGER"
	RoleDeptSupervisor UserRole = "DEPT_SUPERVISOR"

	RoleTeamLead UserRole = "TEAM_LEAD"

	RoleContractor UserRole = "CONTRACTOR"
	RoleEmployee   UserRole = "EMPLOYEE"
)

var InternalInventoryRoles = []UserRole{
	RoleAdmin,
	RoleRegionDirector,
	RoleBranchManager,
	RoleDeptManager,
	RoleRegionLogistician,
	RoleBranchLogistician,
	RoleRegionStorekeeper,
	RoleBranchStorekeeper,
	RoleDeptSupervisor,
}

var WarehouseManagerRoles = []UserRole{
	RoleAdmin,
	RoleRegionDirector,
	RoleRegionLogistician,
	RoleBranchManager,
	RoleBranchLogistician,
	RoleDeptManager,
	RoleTeamLead,
}

var SupplyRequestCreatorRoles = []UserRole{
	RoleAdmin, RoleRegionDirector, RoleBranchManager, RoleDeptManager, RoleTeamLead,
	RoleRegionLogistician, RoleBranchLogistician, RoleDeptSupervisor,
}

var SupplyRequestApproverRoles = []UserRole{
	RoleAdmin, RoleRegionDirector, RoleBranchManager, RoleDeptManager,
	RoleRegionLogistician, RoleBranchLogistician,
}

var ContractorRequestCreatorRoles = []UserRole{
	RoleAdmin, RoleRegionDirector, RoleBranchManager, RoleDeptManager,
	RoleRegionLogistician, RoleBranchLogistician, RoleRegionStorekeeper, RoleBranchStorekeeper, RoleDeptSupervisor,
}

var InventoryManagerRoles = []UserRole{
	RoleAdmin, RoleRegionStorekeeper, RoleBranchStorekeeper, RoleDeptSupervisor,
}

var UnitManagerRoles = []UserRole{
	RoleAdmin, RoleRegionDirector, RoleBranchManager, RoleDeptManager,
	RoleRegionLogistician, RoleBranchLogistician, RoleRegionStorekeeper, RoleBranchStorekeeper,
}

var UserCreatorRoles = []UserRole{
	RoleAdmin, RoleRegionDirector, RoleBranchManager, RoleDeptManager, RoleTeamLead,
	RoleRegionLogistician, RoleBranchLogistician, RoleRegionStorekeeper, RoleBranchStorekeeper, RoleDeptSupervisor,
}

var RoleCreationMap = map[UserRole][]UserRole{
	RoleAdmin: {
		RoleRegionDirector, RoleBranchManager, RoleDeptManager, RoleTeamLead,
		RoleRegionLogistician, RoleRegionStorekeeper, RoleBranchLogistician,
		RoleBranchStorekeeper, RoleDeptSupervisor, RoleEmployee, RoleContractor,
	},
	RoleRegionDirector: {
		RoleBranchManager, RoleDeptManager, RoleTeamLead,
		RoleRegionLogistician, RoleRegionStorekeeper, RoleBranchLogistician,
		RoleBranchStorekeeper, RoleDeptSupervisor, RoleEmployee, RoleContractor,
	},
	RoleBranchManager: {
		RoleDeptManager, RoleTeamLead,
		RoleBranchLogistician, RoleBranchStorekeeper, RoleDeptSupervisor, RoleEmployee,
	},
	RoleDeptManager: {
		RoleTeamLead, RoleDeptSupervisor, RoleEmployee,
	},
	RoleRegionLogistician: {
		RoleRegionStorekeeper, RoleBranchLogistician, RoleBranchStorekeeper,
	},
	RoleBranchLogistician: {
		RoleBranchStorekeeper,
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
	Role         UserRole   `json:"role" gorm:"type:varchar(30);default:'CONTRACTOR'"`
	Status       UserStatus `json:"status" gorm:"type:varchar(20);default:'PENDING'"`
	UnitID       *int64     `json:"unit_id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// Обчислюване поле (не зберігається в users): найвищий тариф у ієрархії unit.
	// Заповнюється в UserRepository.GetByID, використовується фронтом для feature-gating.
	EffectiveSubscriptionTier string `json:"effective_subscription_tier,omitempty" gorm:"-"`
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

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type UpdateProfileRequest struct {
	FullName *string `json:"full_name"`
	Phone    *string `json:"phone"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
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
	case RoleRegionDirector, RoleRegionLogistician, RoleRegionStorekeeper:
		return "REGION"
	case RoleBranchManager, RoleBranchLogistician, RoleBranchStorekeeper:
		return "BRANCH"
	case RoleDeptManager, RoleDeptSupervisor:
		return "DEPARTMENT"
	case RoleTeamLead:
		return "TEAM"
	case RoleEmployee:
		return "ANY"
	default:
		return ""
	}
}

func (r UserRole) CanCreateUnitType(unitType string) bool {
	if r == RoleAdmin {
		return true
	}

	switch r {
	case RoleRegionDirector, RoleRegionLogistician, RoleRegionStorekeeper:
		return unitType == "BRANCH" || unitType == "DEPARTMENT" || unitType == "TEAM"
	case RoleBranchManager, RoleBranchLogistician, RoleBranchStorekeeper:
		return unitType == "DEPARTMENT" || unitType == "TEAM"
	case RoleDeptManager, RoleDeptSupervisor:
		return unitType == "TEAM"
	default:
		return false
	}
}

type UpdateRoleRequest struct {
	Role   string `json:"role" binding:"required"`
	UnitID *int64 `json:"unit_id"`
}

var ApprovalMatrix = map[UserRole][]UserRole{
	RoleTeamLead: {
		RoleDeptManager,
		RoleDeptSupervisor,
		RoleBranchManager,
		RoleBranchLogistician,
		RoleRegionLogistician,
		RoleRegionDirector,
		RoleAdmin,
	},

	RoleDeptManager: {
		RoleBranchManager,
		RoleBranchLogistician,
		RoleRegionLogistician,
		RoleRegionDirector,
		RoleAdmin,
	},
	RoleDeptSupervisor: {
		RoleDeptManager,
		RoleBranchManager,
		RoleBranchLogistician,
		RoleRegionLogistician,
		RoleRegionDirector,
		RoleAdmin,
	},

	RoleBranchManager: {
		RoleRegionDirector,
		RoleRegionLogistician,
		RoleAdmin,
	},

	RoleBranchLogistician: {
		RoleBranchManager,
		RoleRegionDirector, RoleRegionLogistician,
		RoleAdmin,
	},

	RoleBranchStorekeeper: {
		RoleBranchManager, RoleBranchLogistician,
		RoleRegionDirector, RoleRegionLogistician,
		RoleAdmin,
	},

	RoleRegionStorekeeper: {
		RoleRegionDirector, RoleRegionLogistician,
		RoleAdmin,
	},

	RoleRegionDirector: {
		RoleAdmin,
	},
	RoleRegionLogistician: {
		RoleRegionDirector, RoleAdmin,
	},

	// Звичайний співробітник подає заявку своєму ліду або менеджеру
	RoleEmployee: {
		RoleTeamLead, RoleDeptManager, RoleDeptSupervisor,
		RoleBranchManager,
		RoleBranchLogistician,
		RoleRegionLogistician,
		RoleRegionDirector,
		RoleAdmin,
	},
}

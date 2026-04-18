package models

type UnitType string

const (
	UnitRegion     UnitType = "REGION"     // Колишня BRIGADE (Регіон/Дирекція)
	UnitBranch     UnitType = "BRANCH"     // Колишній BATTALION (Філія)
	UnitDepartment UnitType = "DEPARTMENT" // Колишня COMPANY (Відділ)
	UnitTeam       UnitType = "TEAM"       // Колишній PLATOON (Команда)
)

// Unit залишаємо як назву структури (від "Business Unit" або "Організаційна одиниця"),
// щоб не переписувати половину бекенду, де використовується models.Unit
type Unit struct {
	ID               int64    `json:"id"`
	ParentID         *int64   `json:"parent_id"`
	Name             string   `json:"name"`
	UnitType         UnitType `json:"unit_type"`
	SubscriptionTier string   `json:"subscription_tier"`
}

type CreateUnitRequest struct {
	ParentID *int64 `json:"parent_id"`
	Name     string `json:"name" binding:"required"`
	// Оновлено валідацію: тепер дозволені лише нові корпоративні типи
	UnitType UnitType `json:"unit_type" binding:"required,oneof=REGION BRANCH DEPARTMENT TEAM"`
}

// Колишній ChangeCommanderRequest
type ChangeManagerRequest struct {
	NewManagerID string `json:"new_manager_id" binding:"required"` // Колишній new_commander_id
}

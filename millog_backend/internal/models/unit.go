package models

type UnitType string

const (
	UnitBrigade   UnitType = "BRIGADE"
	UnitBattalion UnitType = "BATTALION"
	UnitCompany   UnitType = "COMPANY"
	UnitPlatoon   UnitType = "PLATOON"
)

type Unit struct {
	ID       int64    `json:"id"`
	ParentID *int64   `json:"parent_id"`
	Name     string   `json:"name"`
	UnitType UnitType `json:"unit_type"`
}

type CreateUnitRequest struct {
	ParentID *int64   `json:"parent_id"`
	Name     string   `json:"name" binding:"required"`
	UnitType UnitType `json:"unit_type" binding:"required,oneof=BRIGADE BATTALION COMPANY PLATOON"`
}

type ChangeCommanderRequest struct {
	NewCommanderID string `json:"new_commander_id" binding:"required"`
}

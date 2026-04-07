package models

type PredictStat struct {
	ResourceName  string  `json:"resource_name"`
	CurrentStock  int     `json:"current_stock"`
	DailyBurnRate float64 `json:"daily_burn_rate"`
	DaysLeft      int     `json:"days_left"`
}

type FleetRiskStat struct {
	VehicleName  string `json:"vehicle_name"`
	TotalRefuels int    `json:"total_refuels"`
	Anomalies    int    `json:"anomalies"`
	RiskScore    int    `json:"risk_score"`
}

type UnitReadinessStat struct {
	UnitName       string `json:"unit_name"`
	TotalResources int    `json:"total_resources"`
	ReadyResources int    `json:"ready_resources"`
	ReadinessScore int    `json:"readiness_score"`
}

type FuelMonthlyStat struct {
	Month       string  `json:"month"`
	TotalLiters float64 `json:"total_liters"`
	Anomalies   int     `json:"anomalies"`
}

type MaintenancePredictStat struct {
	VehicleName string `json:"vehicle_name"`
	CurrentOdo  int    `json:"current_odo"`
	NextMaint   int    `json:"next_maint"`
	KmLeft      int    `json:"km_left"`
}

type VolunteerSLAStat struct {
	AverageDays    float64 `json:"average_days"`
	CompletedCount int     `json:"completed_count"`
}

type FleetTCOStat struct {
	VehicleBrand string  `json:"vehicle_brand"`
	TotalCost    float64 `json:"total_cost"`
}

// НОВЕ: Статуси заявок
type VolunteerRequestStat struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

// НОВЕ: Таймлайн заявок
type VolunteerTimelineStat struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// Головна структура
type DashboardAnalytics struct {
	ActiveVehicles    int `json:"active_vehicles"`
	CriticalResources int `json:"critical_resources"`
	FuelAnomalies     int `json:"fuel_anomalies"`

	PredictiveBurnRate []PredictStat            `json:"predictive_burn_rate"`
	FleetRisk          []FleetRiskStat          `json:"fleet_risk"`
	UnitReadiness      []UnitReadinessStat      `json:"unit_readiness"`
	FuelHistory        []FuelMonthlyStat        `json:"fuel_history"`
	MaintenancePredict []MaintenancePredictStat `json:"maintenance_predict"`
	VolunteerSLA       VolunteerSLAStat         `json:"volunteer_sla"`
	FleetTCO           []FleetTCOStat           `json:"fleet_tco"`

	// Поля для вкладки "Волонтери"
	VolunteerFunnel   []VolunteerRequestStat  `json:"volunteer_funnel"`
	VolunteerTimeline []VolunteerTimelineStat `json:"volunteer_timeline"`
}

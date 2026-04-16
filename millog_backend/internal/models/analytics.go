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

type VolunteerRequestStat struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type VolunteerTimelineStat struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// --- НОВЕ ДЛЯ SMART ПОПОВНЕННЯ (МОДАЛЬНЕ ВІКНО) ---

// Товар, якого не вистачає (для відображення у таблиці)
type DeficitResource struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Current int    `json:"current"`
	Min     int    `json:"min"`
	Needed  int    `json:"needed"` // Скільки треба дозамовити
}

// Структура одного елемента, який юзер вибрав для замовлення
type SmartReplenishItem struct {
	ResourceID string `json:"resource_id"`
	Name       string `json:"name"`
	Quantity   int    `json:"quantity"`
	Target     string `json:"target"` // "WAREHOUSE" (Зі складу) або "VOLUNTEER" (Волонтери)
}

// Загальне тіло запиту (Payload), яке прилетить з фронтенду
type SmartReplenishRequest struct {
	Items []SmartReplenishItem `json:"items"`
}

type WarehouseLoadStat struct {
	WarehouseName string `json:"warehouse_name"`
	TotalItems    int    `json:"total_items"`
}

type TopResourceStat struct {
	ResourceName string `json:"resource_name"`
	TotalOrdered int    `json:"total_ordered"`
}

// --- ГОЛОВНА СТРУКТУРА ДАШБОРДУ ---
type DashboardAnalytics struct {
	ActiveVehicles    int `json:"active_vehicles"`
	CriticalResources int `json:"critical_resources"`
	FuelAnomalies     int `json:"fuel_anomalies"`

	WrittenOffResources int `json:"written_off_resources"` // Списане майно
	CompletedRequests   int `json:"completed_requests"`    // Виконані заявки на склади (отримане майно)
	InRepairVehicles    int `json:"in_repair_vehicles"`    // Машини в ремонті
	InactiveVehicles    int `json:"inactive_vehicles"`

	// НОВЕ: Графіки для вкладки Логістика
	WarehouseLoad []WarehouseLoadStat `json:"warehouse_load"`
	TopResources  []TopResourceStat   `json:"top_resources"`

	PredictiveBurnRate []PredictStat            `json:"predictive_burn_rate"`
	FleetRisk          []FleetRiskStat          `json:"fleet_risk"`
	UnitReadiness      []UnitReadinessStat      `json:"unit_readiness"`
	FuelHistory        []FuelMonthlyStat        `json:"fuel_history"`
	MaintenancePredict []MaintenancePredictStat `json:"maintenance_predict"`
	VolunteerSLA       VolunteerSLAStat         `json:"volunteer_sla"`
	FleetTCO           []FleetTCOStat           `json:"fleet_tco"`

	VolunteerFunnel   []VolunteerRequestStat  `json:"volunteer_funnel"`
	VolunteerTimeline []VolunteerTimelineStat `json:"volunteer_timeline"`

	// НОВЕ: Поле для списку дефіциту (для модалки та віджета)
	DeficitResources []DeficitResource `json:"deficit_resources"`
}

package models

type PredictStat struct {
	ResourceName  string  `json:"resource_name"`
	CurrentStock  int     `json:"current_stock"`
	MinQuantity   int     `json:"min_quantity"`
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

type ContractorSLA struct {
	AverageDays    float64 `json:"average_days"`
	CompletedCount int     `json:"completed_count"`
	// --- ДОДАЄМО НОВІ ПОЛЯ ---
	OTDPercentage int     `json:"otd_percentage"` // On-Time Delivery (%)
	OverdueCount  int     `json:"overdue_count"`  // Прострочені завдання (шт)
	FastestDays   float64 `json:"fastest_days"`   // Найшвидше виконання (днів)
}

type FleetTCOStat struct {
	VehicleBrand string  `json:"vehicle_brand"`
	TotalCost    float64 `json:"total_cost"`
}

type CONTRACTORRequestStat struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type CONTRACTORTimelineStat struct {
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
	Target     string `json:"target"` // "WAREHOUSE" (Зі складу) або "CONTRACTOR" (Волонтери)
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

// Фінансова вартість залишків по складах
type WarehouseValueStat struct {
	WarehouseName string  `json:"warehouse_name"`
	TotalValue    float64 `json:"total_value"`
}

// Найдорожчі позиції майна
type TopCostlyResourceStat struct {
	ResourceName string  `json:"resource_name"`
	TotalValue   float64 `json:"total_value"`
	Quantity     int     `json:"quantity"`
	UnitPrice    float64 `json:"unit_price"`
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

	// Фінансові метрики
	InventoryTotalValue float64 `json:"inventory_total_value"` // Загальна вартість залишків (грн)
	WriteOffTotalValue  float64 `json:"write_off_total_value"` // Вартість списаного за період (грн)

	// НОВЕ: Графіки для вкладки Логістика
	WarehouseLoad []WarehouseLoadStat `json:"warehouse_load"`
	TopResources  []TopResourceStat   `json:"top_resources"`

	// Фінансові графіки
	WarehouseValueStats []WarehouseValueStat    `json:"warehouse_value_stats"` // Вартість по складах
	TopCostlyResources  []TopCostlyResourceStat `json:"top_costly_resources"`  // Найдорожчі позиції

	PredictiveBurnRate []PredictStat            `json:"predictive_burn_rate"`
	FleetRisk          []FleetRiskStat          `json:"fleet_risk"`
	UnitReadiness      []UnitReadinessStat      `json:"unit_readiness"`
	FuelHistory        []FuelMonthlyStat        `json:"fuel_history"`
	MaintenancePredict []MaintenancePredictStat `json:"maintenance_predict"`
	CONTRACTORSLA      ContractorSLA            `json:"CONTRACTOR_sla"`
	FleetTCO           []FleetTCOStat           `json:"fleet_tco"`

	CONTRACTORFunnel   []CONTRACTORRequestStat  `json:"CONTRACTOR_funnel"`
	CONTRACTORTimeline []CONTRACTORTimelineStat `json:"CONTRACTOR_timeline"`

	// НОВЕ: Поле для списку дефіциту (для модалки та віджета)
	DeficitResources []DeficitResource `json:"deficit_resources"`
}

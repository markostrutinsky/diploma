package models

// --- Базові структури для метрик ---
type ConditionStat struct {
	Condition string `json:"condition"`
	Count     int    `json:"count"`
}

type LocationStat struct {
	LocationType string `json:"location_type"`
	Count        int    `json:"count"`
}

type VolunteerRequestStat struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

// --- Нові структури для розширеної аналітики ---
type FuelMonthlyStat struct {
	Month       string  `json:"month"`
	TotalLiters float64 `json:"total_liters"`
	Anomalies   int     `json:"anomalies"`
}

type TacticalStat struct {
	LocationType string `json:"location_type"`
	NewItems     int    `json:"new_items"`
	UsedItems    int    `json:"used_items"`
}

// --- Головна структура, яка збирає все разом ---
type DashboardAnalytics struct {
	// ТОП Картки
	ActiveVehicles    int `json:"active_vehicles"`
	CriticalResources int `json:"critical_resources"`
	FuelAnomalies     int `json:"fuel_anomalies"`

	// Дані для графіків
	TacticalStats   []TacticalStat         `json:"tactical_stats"`
	BurnRate        []ConditionStat        `json:"burn_rate"`
	CriticalItems   []Resource             `json:"critical_items"`
	FuelHistory     []FuelMonthlyStat      `json:"fuel_history"`
	FleetHealth     []ConditionStat        `json:"fleet_health"`
	VolunteerFunnel []VolunteerRequestStat `json:"volunteer_funnel"`
}

package repositories

import (
	"context"

	"millog_backend/internal/models"
)

type AnalyticsRepository struct{}

func NewAnalyticsRepository() *AnalyticsRepository {
	return &AnalyticsRepository{}
}

func (r *AnalyticsRepository) GetDashboardStats(ctx context.Context, db DBExecutor) (*models.DashboardAnalytics, error) {
	var stats models.DashboardAnalytics

	// 1. ТОП Метрики
	metricsQuery := `
		SELECT 
			(SELECT COUNT(*) FROM vehicles WHERE status = 'ACTIVE') as active_vehicles,
			(SELECT COUNT(*) FROM resources WHERE quantity <= min_quantity AND quantity > 0 AND condition != 'WRITTEN_OFF') as critical_resources,
			(SELECT COUNT(*) FROM fuel_records WHERE is_anomaly = true) as fuel_anomalies
	`
	db.QueryRow(ctx, metricsQuery).Scan(&stats.ActiveVehicles, &stats.CriticalResources, &stats.FuelAnomalies)

	// 2. Ешелонування (TacticalStats)
	queryTactical := `
		SELECT 
			COALESCE(w.location_type, 'UNASSIGNED'),
			COUNT(r.id) FILTER (WHERE r.condition = 'NEW') as new_items,
			COUNT(r.id) FILTER (WHERE r.condition = 'USED') as used_items
		FROM resources r
		LEFT JOIN warehouses w ON r.warehouse_id = w.id
		WHERE r.condition != 'WRITTEN_OFF'
		GROUP BY w.location_type
	`
	tRows, err := db.Query(ctx, queryTactical)
	if err == nil {
		for tRows.Next() {
			var s models.TacticalStat
			tRows.Scan(&s.LocationType, &s.NewItems, &s.UsedItems)
			stats.TacticalStats = append(stats.TacticalStats, s)
		}
		tRows.Close()
	}

	// 3. Зношеність (BurnRate)
	brRows, err := db.Query(ctx, `SELECT condition, COUNT(id) FROM resources GROUP BY condition`)
	if err == nil {
		for brRows.Next() {
			var s models.ConditionStat
			brRows.Scan(&s.Condition, &s.Count)
			stats.BurnRate = append(stats.BurnRate, s)
		}
		brRows.Close()
	}

	// 4. Критичні залишки (CriticalItems)
	queryCrit := `SELECT id, name, quantity, min_quantity FROM resources WHERE quantity <= min_quantity AND condition != 'WRITTEN_OFF' LIMIT 5`
	cRows, err := db.Query(ctx, queryCrit)
	if err == nil {
		for cRows.Next() {
			var res models.Resource
			cRows.Scan(&res.ID, &res.Name, &res.Quantity, &res.MinQuantity)
			stats.CriticalItems = append(stats.CriticalItems, res)
		}
		cRows.Close()
	}

	// 5. Історія пального (FuelHistory)
	queryFuel := `
		SELECT 
			TO_CHAR(created_at, 'Mon') as month,
			SUM(liters) as total,
			COUNT(id) FILTER (WHERE is_anomaly = true) as anomalies
		FROM fuel_records
		WHERE created_at > NOW() - INTERVAL '6 months'
		GROUP BY month, DATE_TRUNC('month', created_at)
		ORDER BY DATE_TRUNC('month', created_at)
	`
	fRows, err := db.Query(ctx, queryFuel)
	if err == nil {
		for fRows.Next() {
			var f models.FuelMonthlyStat
			fRows.Scan(&f.Month, &f.TotalLiters, &f.Anomalies)
			stats.FuelHistory = append(stats.FuelHistory, f)
		}
		fRows.Close()
	}

	// 6. Стан автопарку (FleetHealth)
	// Використовуємо існуючу структуру ConditionStat, де Condition - це статус авто
	fhRows, err := db.Query(ctx, `SELECT status, COUNT(id) FROM vehicles GROUP BY status`)
	if err == nil {
		for fhRows.Next() {
			var s models.ConditionStat
			fhRows.Scan(&s.Condition, &s.Count)
			stats.FleetHealth = append(stats.FleetHealth, s)
		}
		fhRows.Close()
	}

	// 7. Волонтерські заявки (VolunteerFunnel)
	vRows, err := db.Query(ctx, `SELECT status, COUNT(id) FROM volunteer_requests GROUP BY status`)
	if err == nil {
		for vRows.Next() {
			var v models.VolunteerRequestStat
			vRows.Scan(&v.Status, &v.Count)
			stats.VolunteerFunnel = append(stats.VolunteerFunnel, v)
		}
		vRows.Close()
	}

	return &stats, nil
}

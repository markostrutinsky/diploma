package repositories

import (
	"context"
	"time"

	"millog_backend/internal/models"
)

type AnalyticsRepository struct{}

func NewAnalyticsRepository() *AnalyticsRepository {
	return &AnalyticsRepository{}
}

func (r *AnalyticsRepository) GetDashboardStats(ctx context.Context, db DBExecutor, startDateStr, endDateStr string) (*models.DashboardAnalytics, error) {
	var stats models.DashboardAnalytics

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		startDate = time.Now().AddDate(0, 0, -30)
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		endDate = time.Now()
	}
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	// 1. ТОП Метрики
	db.QueryRow(ctx, `
		SELECT 
			(SELECT COUNT(*) FROM vehicles WHERE status = 'ACTIVE'),
			(SELECT COUNT(*) FROM resources WHERE quantity <= min_quantity AND quantity > 0 AND condition != 'WRITTEN_OFF'),
			(SELECT COUNT(*) FROM fuel_records WHERE is_anomaly = true AND created_at BETWEEN $1 AND $2)
	`, startDate, endDate).Scan(&stats.ActiveVehicles, &stats.CriticalResources, &stats.FuelAnomalies)

	// 2. ПРОГНОЗ ВИЧЕРПАННЯ
	queryPredict := `
		WITH consumption AS (
			SELECT resource_id, SUM(quantity) as consumed FROM supply_requests
			WHERE created_at BETWEEN $1 AND $2 AND status IN ('APPROVED', 'COMPLETED') GROUP BY resource_id
		)
		SELECT r.name, r.quantity, c.consumed / NULLIF(EXTRACT(EPOCH FROM ($2 - $1))/86400, 0) as daily_burn,
			(r.quantity / NULLIF(c.consumed / NULLIF(EXTRACT(EPOCH FROM ($2 - $1))/86400, 0), 0))::int as days_left
		FROM resources r JOIN consumption c ON r.id = c.resource_id
		WHERE r.condition != 'WRITTEN_OFF' AND c.consumed > 0 ORDER BY days_left ASC LIMIT 6
	`
	pRows, _ := db.Query(ctx, queryPredict, startDate, endDate)
	defer pRows.Close()
	for pRows.Next() {
		var p models.PredictStat
		pRows.Scan(&p.ResourceName, &p.CurrentStock, &p.DailyBurnRate, &p.DaysLeft)
		stats.PredictiveBurnRate = append(stats.PredictiveBurnRate, p)
	}

	// 3. АНТИКОРУПЦІЙНИЙ ІНДЕКС (ГСМ)
	queryRisk := `
		SELECT v.brand || ' ' || v.model || ' (' || v.plate_number || ')', COUNT(f.id),
			COUNT(f.id) FILTER (WHERE f.is_anomaly = true),
			CASE WHEN COUNT(f.id) > 0 THEN (COUNT(f.id) FILTER (WHERE f.is_anomaly = true) * 100 / COUNT(f.id)) ELSE 0 END as score
		FROM vehicles v JOIN fuel_records f ON v.id = f.vehicle_id WHERE f.created_at BETWEEN $1 AND $2
		GROUP BY v.id, v.brand, v.model, v.plate_number HAVING COUNT(f.id) FILTER (WHERE f.is_anomaly = true) > 0 ORDER BY score DESC LIMIT 5
	`
	rRows, _ := db.Query(ctx, queryRisk, startDate, endDate)
	defer rRows.Close()
	for rRows.Next() {
		var r models.FleetRiskStat
		rRows.Scan(&r.VehicleName, &r.TotalRefuels, &r.Anomalies, &r.RiskScore)
		stats.FleetRisk = append(stats.FleetRisk, r)
	}

	// 4. БОЄГОТОВНІСТЬ ПІДРОЗДІЛІВ
	queryReadiness := `
		SELECT u.name, COUNT(r.id), COUNT(r.id) FILTER (WHERE r.quantity >= r.min_quantity AND r.quantity > 0),
			CASE WHEN COUNT(r.id) > 0 THEN (COUNT(r.id) FILTER (WHERE r.quantity >= r.min_quantity AND r.quantity > 0) * 100 / COUNT(r.id)) ELSE 0 END as score
		FROM units u LEFT JOIN resources r ON u.id = r.unit_id AND r.condition != 'WRITTEN_OFF'
		GROUP BY u.id, u.name HAVING COUNT(r.id) > 0 ORDER BY score ASC
	`
	uRows, _ := db.Query(ctx, queryReadiness)
	defer uRows.Close()
	for uRows.Next() {
		var u models.UnitReadinessStat
		uRows.Scan(&u.UnitName, &u.TotalResources, &u.ReadyResources, &u.ReadinessScore)
		stats.UnitReadiness = append(stats.UnitReadiness, u)
	}

	// 5. КАРДІОЛІНІЯ ГСМ
	queryFuelCardio := `
		SELECT TO_CHAR(DATE(created_at), 'DD.MM'), SUM(liters), COUNT(id) FILTER (WHERE is_anomaly = true)
		FROM fuel_records WHERE created_at BETWEEN $1 AND $2 GROUP BY DATE(created_at) ORDER BY DATE(created_at)
	`
	fRows, _ := db.Query(ctx, queryFuelCardio, startDate, endDate)
	defer fRows.Close()
	for fRows.Next() {
		var f models.FuelMonthlyStat
		fRows.Scan(&f.Month, &f.TotalLiters, &f.Anomalies)
		stats.FuelHistory = append(stats.FuelHistory, f)
	}

	// 6. ПРОГНОЗ ТО (Скільки км залишилось)
	queryMaint := `
		SELECT 
			v.brand || ' ' || v.plate_number,
			COALESCE((SELECT MAX(odometer_km) FROM fuel_records WHERE vehicle_id = v.id), v.last_maintenance_odometer) as current_odo,
			(v.last_maintenance_odometer + v.maintenance_interval_km) as next_maint,
			(v.last_maintenance_odometer + v.maintenance_interval_km) - COALESCE((SELECT MAX(odometer_km) FROM fuel_records WHERE vehicle_id = v.id), v.last_maintenance_odometer) as km_left
		FROM vehicles v WHERE v.status = 'ACTIVE' AND v.maintenance_interval_km > 0
		ORDER BY km_left ASC LIMIT 6
	`
	mRows, _ := db.Query(ctx, queryMaint)
	defer mRows.Close()
	for mRows.Next() {
		var m models.MaintenancePredictStat
		mRows.Scan(&m.VehicleName, &m.CurrentOdo, &m.NextMaint, &m.KmLeft)
		stats.MaintenancePredict = append(stats.MaintenancePredict, m)
	}

	// 7. ЕФЕКТИВНІСТЬ ВОЛОНТЕРІВ (SLA)
	querySLA := `
		SELECT 
			COALESCE(AVG(EXTRACT(EPOCH FROM (completed_at - created_at))/86400), 0) as avg_days,
			COUNT(id)
		FROM volunteer_requests WHERE status = 'COMPLETED' AND completed_at IS NOT NULL AND created_at BETWEEN $1 AND $2
	`
	db.QueryRow(ctx, querySLA, startDate, endDate).Scan(&stats.VolunteerSLA.AverageDays, &stats.VolunteerSLA.CompletedCount)

	// 8. ВИТРАТИ НА РЕМОНТИ ПО МАРКАХ (TCO)
	queryTCO := `
		SELECT COALESCE(v.brand, 'Інше'), SUM(m.cost_amount) as total_cost
		FROM maintenance_records m JOIN vehicles v ON m.vehicle_id = v.id
		WHERE m.created_at BETWEEN $1 AND $2 GROUP BY v.brand ORDER BY total_cost DESC LIMIT 5
	`
	tcoRows, _ := db.Query(ctx, queryTCO, startDate, endDate)
	defer tcoRows.Close()
	for tcoRows.Next() {
		var t models.FleetTCOStat
		tcoRows.Scan(&t.VehicleBrand, &t.TotalCost)
		stats.FleetTCO = append(stats.FleetTCO, t)
	}

	// 9. ВОЛОНТЕРСЬКА ВОРОНКА (Статуси)
	queryFunnel := `
		SELECT status, COUNT(id) FROM volunteer_requests 
		WHERE created_at BETWEEN $1 AND $2 
		GROUP BY status
	`
	vRows, _ := db.Query(ctx, queryFunnel, startDate, endDate)
	defer vRows.Close()
	for vRows.Next() {
		var v models.VolunteerRequestStat
		vRows.Scan(&v.Status, &v.Count)
		stats.VolunteerFunnel = append(stats.VolunteerFunnel, v)
	}

	// 10. ДИНАМІКА ВОЛОНТЕРСЬКИХ ЗАЯВОК (Графік)
	queryTimeline := `
		SELECT TO_CHAR(DATE(created_at), 'DD.MM'), COUNT(id) 
		FROM volunteer_requests 
		WHERE created_at BETWEEN $1 AND $2 
		GROUP BY DATE(created_at) 
		ORDER BY DATE(created_at)
	`
	tRows, _ := db.Query(ctx, queryTimeline, startDate, endDate)
	defer tRows.Close()
	for tRows.Next() {
		var vt models.VolunteerTimelineStat
		tRows.Scan(&vt.Date, &vt.Count)
		stats.VolunteerTimeline = append(stats.VolunteerTimeline, vt)
	}

	return &stats, nil
}

func (r *AnalyticsRepository) CreateAutoReplenishRequests(ctx context.Context, db DBExecutor) (int, error) {
	// SQL-запит: Знайти всі дефіцитні ресурси, для яких немає активних (OPEN/IN_PROGRESS) заявок
	// І автоматично створити нові заявки типу 'VOLUNTEER'
	query := `
		INSERT INTO supply_requests (resource_id, quantity, status, type, description, created_at, updated_at)
		SELECT 
			r.id, 
			(r.min_quantity * 2 - r.quantity) as needed_qty, -- Поповнюємо з запасом
			'OPEN', 
			'VOLUNTEER', 
			'Автоматично згенерована потреба на основі дефіциту залишків',
			NOW(), 
			NOW()
		FROM resources r
		LEFT JOIN supply_requests sr ON r.id = sr.resource_id AND sr.status IN ('OPEN', 'IN_PROGRESS')
		WHERE r.quantity <= r.min_quantity 
		  AND r.condition != 'WRITTEN_OFF'
		  AND sr.id IS NULL -- Тільки якщо по цьому ресурсу ще немає відкритої заявки
		RETURNING id;
	`
	rows, err := db.Query(ctx, query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
	}
	return count, nil
}

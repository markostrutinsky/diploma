package repositories

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"millog_backend/internal/models"
)

type AnalyticsRepository struct{}

func NewAnalyticsRepository() *AnalyticsRepository {
	return &AnalyticsRepository{}
}

func (r *AnalyticsRepository) GetDashboardStats(ctx context.Context, db DBExecutor, startDateStr, endDateStr, unitID string) (*models.DashboardAnalytics, error) {
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

	resFilter := ""
	resFilterPrefix := ""
	volFilter := ""
	unitFilter := ""

	if unitID != "" {
		if _, err := strconv.Atoi(unitID); err == nil {
			resFilter = fmt.Sprintf(" AND unit_id = %s", unitID)
			resFilterPrefix = fmt.Sprintf(" AND r.unit_id = %s", unitID)
			volFilter = fmt.Sprintf(" AND unit_id = %s", unitID)
			unitFilter = fmt.Sprintf(" WHERE u.id = %s", unitID)
		}
	}

	// 1. ТОП Метрики
	queryMetrics := fmt.Sprintf(`
		SELECT 
			(SELECT COUNT(*) FROM vehicles WHERE status = 'ACTIVE'),
			(SELECT COUNT(*) FROM resources WHERE quantity <= min_quantity AND quantity > 0 AND condition != 'WRITTEN_OFF' %s),
			(SELECT COUNT(*) FROM fuel_records WHERE is_anomaly = true AND created_at BETWEEN $1 AND $2)
	`, resFilter)
	db.QueryRow(ctx, queryMetrics, startDate, endDate).Scan(&stats.ActiveVehicles, &stats.CriticalResources, &stats.FuelAnomalies)

	// ====================================================================
	// ЖИТТЄВИЙ ЦИКЛ (Виправлено під статус APPROVED)
	// ====================================================================
	queryWrittenOff := fmt.Sprintf(`SELECT COUNT(*) FROM resources WHERE condition = 'WRITTEN_OFF' %s`, resFilter)
	db.QueryRow(ctx, queryWrittenOff).Scan(&stats.WrittenOffResources)

	queryCompletedReqs := `
		SELECT COUNT(*) FROM supply_requests 
		WHERE status = 'APPROVED' AND updated_at BETWEEN $1 AND $2
	`
	db.QueryRow(ctx, queryCompletedReqs, startDate, endDate).Scan(&stats.CompletedRequests)

	db.QueryRow(ctx, `SELECT COUNT(*) FROM vehicles WHERE status = 'IN_REPAIR'`).Scan(&stats.InRepairVehicles)
	db.QueryRow(ctx, `SELECT COUNT(*) FROM vehicles WHERE status = 'INACTIVE'`).Scan(&stats.InactiveVehicles)

	// ====================================================================
	// 2. ПРОГНОЗ ВИЧЕРПАННЯ (Виправлено під статус APPROVED)
	// ====================================================================
	queryPredict := fmt.Sprintf(`
		WITH consumption AS (
			SELECT resource_id, SUM(quantity) as consumed FROM supply_requests
			WHERE created_at BETWEEN $1 AND $2 AND status = 'APPROVED' GROUP BY resource_id
		)
		SELECT r.name, r.quantity, c.consumed / NULLIF(EXTRACT(EPOCH FROM ($2 - $1))/86400, 0) as daily_burn,
			(r.quantity / NULLIF(c.consumed / NULLIF(EXTRACT(EPOCH FROM ($2 - $1))/86400, 0), 0))::int as days_left
		FROM resources r JOIN consumption c ON r.id = c.resource_id
		WHERE r.condition != 'WRITTEN_OFF' AND c.consumed > 0 %s 
		ORDER BY days_left ASC
	`, resFilterPrefix)
	pRows, _ := db.Query(ctx, queryPredict, startDate, endDate)
	defer pRows.Close()
	for pRows.Next() {
		var p models.PredictStat
		pRows.Scan(&p.ResourceName, &p.CurrentStock, &p.DailyBurnRate, &p.DaysLeft)
		stats.PredictiveBurnRate = append(stats.PredictiveBurnRate, p)
	}

	// 3. АНТИКОРУПЦІЙНИЙ ІНДЕКС
	queryRisk := `
		SELECT v.brand || ' ' || v.model || ' (' || v.plate_number || ')', COUNT(f.id),
			COUNT(f.id) FILTER (WHERE f.is_anomaly = true),
			CASE WHEN COUNT(f.id) > 0 THEN (COUNT(f.id) FILTER (WHERE f.is_anomaly = true) * 100 / COUNT(f.id)) ELSE 0 END as score
		FROM vehicles v JOIN fuel_records f ON v.id = f.vehicle_id WHERE f.created_at BETWEEN $1 AND $2
		GROUP BY v.id, v.brand, v.model, v.plate_number HAVING COUNT(f.id) FILTER (WHERE f.is_anomaly = true) > 0 ORDER BY score DESC
	`
	rRows, _ := db.Query(ctx, queryRisk, startDate, endDate)
	defer rRows.Close()
	for rRows.Next() {
		var r models.FleetRiskStat
		rRows.Scan(&r.VehicleName, &r.TotalRefuels, &r.Anomalies, &r.RiskScore)
		stats.FleetRisk = append(stats.FleetRisk, r)
	}

	// 4. ЗАБЕЗПЕЧЕНІСТЬ ПІДРОЗДІЛІВ
	queryReadiness := fmt.Sprintf(`
		SELECT u.name, COUNT(r.id), COUNT(r.id) FILTER (WHERE r.quantity >= r.min_quantity AND r.quantity > 0),
			CASE WHEN COUNT(r.id) > 0 THEN (COUNT(r.id) FILTER (WHERE r.quantity >= r.min_quantity AND r.quantity > 0) * 100 / COUNT(r.id)) ELSE 0 END as score
		FROM units u LEFT JOIN resources r ON u.id = r.unit_id AND r.condition != 'WRITTEN_OFF'
		%s
		GROUP BY u.id, u.name HAVING COUNT(r.id) > 0 ORDER BY score ASC
	`, unitFilter)
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

	// 6. ПРОГНОЗ ТО
	queryMaint := `
		SELECT 
			v.brand || ' ' || v.plate_number,
			COALESCE((SELECT MAX(odometer_km) FROM fuel_records WHERE vehicle_id = v.id), v.last_maintenance_odometer) as current_odo,
			(v.last_maintenance_odometer + v.maintenance_interval_km) as next_maint,
			(v.last_maintenance_odometer + v.maintenance_interval_km) - COALESCE((SELECT MAX(odometer_km) FROM fuel_records WHERE vehicle_id = v.id), v.last_maintenance_odometer) as km_left
		FROM vehicles v WHERE v.status = 'ACTIVE' AND v.maintenance_interval_km > 0
		ORDER BY km_left ASC
	`
	mRows, _ := db.Query(ctx, queryMaint)
	defer mRows.Close()
	for mRows.Next() {
		var m models.MaintenancePredictStat
		mRows.Scan(&m.VehicleName, &m.CurrentOdo, &m.NextMaint, &m.KmLeft)
		stats.MaintenancePredict = append(stats.MaintenancePredict, m)
	}

	// 7. ЕФЕКТИВНІСТЬ ВОЛОНТЕРІВ (Якщо тут теж треба APPROVED, зміни COMPLETED на APPROVED)
	querySLA := fmt.Sprintf(`
		SELECT 
			COALESCE(AVG(EXTRACT(EPOCH FROM (completed_at - created_at))/86400), 0) as avg_days,
			COUNT(id)
		FROM volunteer_requests WHERE status = 'COMPLETED' AND completed_at IS NOT NULL AND created_at BETWEEN $1 AND $2 %s
	`, volFilter)
	db.QueryRow(ctx, querySLA, startDate, endDate).Scan(&stats.VolunteerSLA.AverageDays, &stats.VolunteerSLA.CompletedCount)

	// 8. ВИТРАТИ НА РЕМОНТИ
	queryTCO := `
		SELECT COALESCE(v.brand, 'Інше'), SUM(m.cost_amount) as total_cost
		FROM maintenance_records m JOIN vehicles v ON m.vehicle_id = v.id
		WHERE m.created_at BETWEEN $1 AND $2 GROUP BY v.brand ORDER BY total_cost DESC
	`
	tcoRows, _ := db.Query(ctx, queryTCO, startDate, endDate)
	defer tcoRows.Close()
	for tcoRows.Next() {
		var t models.FleetTCOStat
		tcoRows.Scan(&t.VehicleBrand, &t.TotalCost)
		stats.FleetTCO = append(stats.FleetTCO, t)
	}

	// 9. ВОЛОНТЕРСЬКА ВОРОНКА
	queryFunnel := fmt.Sprintf(`
		SELECT status, COUNT(id) FROM volunteer_requests 
		WHERE created_at BETWEEN $1 AND $2 %s
		GROUP BY status
	`, volFilter)
	vRows, _ := db.Query(ctx, queryFunnel, startDate, endDate)
	defer vRows.Close()
	for vRows.Next() {
		var v models.VolunteerRequestStat
		vRows.Scan(&v.Status, &v.Count)
		stats.VolunteerFunnel = append(stats.VolunteerFunnel, v)
	}

	// 10. ДИНАМІКА ВОЛОНТЕРСЬКИХ ЗАЯВОК
	queryTimeline := fmt.Sprintf(`
		SELECT TO_CHAR(DATE(created_at), 'DD.MM'), COUNT(id) 
		FROM volunteer_requests 
		WHERE created_at BETWEEN $1 AND $2 %s
		GROUP BY DATE(created_at) 
		ORDER BY DATE(created_at)
	`, volFilter)
	tRows, _ := db.Query(ctx, queryTimeline, startDate, endDate)
	defer tRows.Close()
	for tRows.Next() {
		var vt models.VolunteerTimelineStat
		tRows.Scan(&vt.Date, &vt.Count)
		stats.VolunteerTimeline = append(stats.VolunteerTimeline, vt)
	}

	// 11. СПИСОК ДЕФІЦИТУ ДЛЯ SMART-ПОПОВНЕННЯ (З урахуванням майна в дорозі)
	queryDeficit := fmt.Sprintf(`
		WITH PendingOrders AS (
			-- Рахуємо скільки майна ВЖЕ замовлено (висить у відкритих заявках)
			SELECT resource_id, SUM(quantity) as pending_qty
			FROM supply_requests
			WHERE status IN ('OPEN', 'IN_PROGRESS', 'APPROVED') 
			GROUP BY resource_id
		)
		SELECT 
			r.id, 
			r.name, 
			r.quantity, 
			r.min_quantity, 
			-- Формула: (Мінімум * 2) - (Фактичний залишок + Вже замовлено)
			(r.min_quantity * 2 - (r.quantity + COALESCE(p.pending_qty, 0))) as needed
		FROM resources r
		LEFT JOIN PendingOrders p ON r.id = p.resource_id
		-- Показуємо тільки те, де ФАКТ + В ДОРОЗІ все ще менше або дорівнює мінімуму
		WHERE (r.quantity + COALESCE(p.pending_qty, 0)) <= r.min_quantity 
		  AND r.condition != 'WRITTEN_OFF' %s
	`, resFilterPrefix) // УВАГА: тут ми змінили resFilter на resFilterPrefix, бо додали аліас 'r.'

	dRows, _ := db.Query(ctx, queryDeficit)
	defer dRows.Close()
	for dRows.Next() {
		var d models.DeficitResource
		dRows.Scan(&d.ID, &d.Name, &d.Current, &d.Min, &d.Needed)
		stats.DeficitResources = append(stats.DeficitResources, d)
	}

	return &stats, nil
}

// НОВА ФУНКЦІЯ: Обробка вибраних логістом позицій
func (r *AnalyticsRepository) ProcessSmartReplenish(ctx context.Context, db DBExecutor, req models.SmartReplenishRequest, userID string) (int, error) {
	count := 0

	for _, item := range req.Items {
		if item.Target == "WAREHOUSE" {
			// Створюємо офіційну заявку на забезпечення (на склад)
			_, err := db.Exec(ctx, `
				INSERT INTO supply_requests (resource_id, quantity, status, created_by, comment, created_at)
				VALUES ($1, $2, 'OPEN', $3, 'Автоматичне замовлення через Smart-модуль', NOW())
			`, item.ResourceID, item.Quantity, userID)
			if err == nil {
				count++
			}
		} else if item.Target == "VOLUNTEER" {
			// Створюємо запит для волонтерів
			title := "Потреба: " + item.Name
			desc := fmt.Sprintf("Автоматично сформована потреба для підрозділу на %s (Кількість: %d)", item.Name, item.Quantity)

			_, err := db.Exec(ctx, `
				INSERT INTO volunteer_requests (created_by, title, description, status, created_at)
				VALUES ($1, $2, $3, 'OPEN', NOW())
			`, userID, title, desc)
			if err == nil {
				count++
			}
		}
	}

	return count, nil
}

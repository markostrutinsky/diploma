package repositories

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"Omnilog_backend/internal/models"

	"github.com/google/uuid"
)

type AnalyticsRepository struct{}

func NewAnalyticsRepository() *AnalyticsRepository {
	return &AnalyticsRepository{}
}

// tcondBuilder повертає функцію-хелпер, яка формує " AND {alias}.tenant_id = '{uuid}'"
// з ВАЛІДОВАНОГО UUID (tid походить з JWT, але ми додатково парсимо).
// Для SYSTEM_ADMIN (tid == "") повертає "".
func tcondBuilder(ctx context.Context) func(alias string) string {
	tid := TenantFromCtx(ctx)
	if tid == "" {
		return func(string) string { return "" }
	}
	if _, err := uuid.Parse(tid); err != nil {
		return func(string) string { return "" }
	}
	return func(alias string) string {
		if alias == "" {
			return fmt.Sprintf(" AND tenant_id = '%s'", tid)
		}
		return fmt.Sprintf(" AND %s.tenant_id = '%s'", alias, tid)
	}
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

	tcond := tcondBuilder(ctx)

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
			(SELECT COUNT(*) FROM vehicles WHERE status = 'ACTIVE'%s),
			(SELECT COUNT(*) FROM resources WHERE quantity <= min_quantity AND quantity > 0 AND condition != 'WRITTEN_OFF' %s%s),
			(SELECT COUNT(*) FROM fuel_records WHERE is_anomaly = true AND created_at BETWEEN $1 AND $2%s)
	`, tcond(""), resFilter, tcond(""), tcond(""))
	db.QueryRow(ctx, queryMetrics, startDate, endDate).Scan(&stats.ActiveVehicles, &stats.CriticalResources, &stats.FuelAnomalies)

	// ====================================================================
	// ЖИТТЄВИЙ ЦИКЛ (Виправлено під статус APPROVED)
	// ====================================================================
	queryWrittenOff := fmt.Sprintf(`SELECT COUNT(*) FROM resources WHERE condition = 'WRITTEN_OFF' %s%s`, resFilter, tcond(""))
	db.QueryRow(ctx, queryWrittenOff).Scan(&stats.WrittenOffResources)

	queryCompletedReqs := fmt.Sprintf(`
		SELECT COUNT(*) FROM supply_requests 
		WHERE status = 'APPROVED' AND updated_at BETWEEN $1 AND $2%s
	`, tcond(""))
	db.QueryRow(ctx, queryCompletedReqs, startDate, endDate).Scan(&stats.CompletedRequests)

	db.QueryRow(ctx, fmt.Sprintf(`SELECT COUNT(*) FROM vehicles WHERE status = 'IN_REPAIR'%s`, tcond(""))).Scan(&stats.InRepairVehicles)
	db.QueryRow(ctx, fmt.Sprintf(`SELECT COUNT(*) FROM vehicles WHERE status = 'INACTIVE'%s`, tcond(""))).Scan(&stats.InactiveVehicles)

	// ====================================================================
	// 2. ПРОГНОЗ ВИЧЕРПАННЯ (Виправлено під статус APPROVED)
	// ====================================================================
	queryPredict := fmt.Sprintf(`
		WITH consumption AS (
			SELECT resource_id, SUM(quantity) as consumed FROM supply_requests
			WHERE created_at BETWEEN $1 AND $2 AND status = 'APPROVED'%s GROUP BY resource_id
		)
		SELECT r.name, r.quantity, c.consumed / NULLIF(EXTRACT(EPOCH FROM ($2 - $1))/86400, 0) as daily_burn,
			(r.quantity / NULLIF(c.consumed / NULLIF(EXTRACT(EPOCH FROM ($2 - $1))/86400, 0), 0))::int as days_left
		FROM resources r JOIN consumption c ON r.id = c.resource_id
		WHERE r.condition != 'WRITTEN_OFF' AND c.consumed > 0 %s%s 
		ORDER BY days_left ASC
	`, tcond(""), resFilterPrefix, tcond("r"))
	pRows, _ := db.Query(ctx, queryPredict, startDate, endDate)
	defer pRows.Close()
	for pRows.Next() {
		var p models.PredictStat
		pRows.Scan(&p.ResourceName, &p.CurrentStock, &p.DailyBurnRate, &p.DaysLeft)
		stats.PredictiveBurnRate = append(stats.PredictiveBurnRate, p)
	}

	// 3. АНТИКОРУПЦІЙНИЙ ІНДЕКС
	queryRisk := fmt.Sprintf(`
		SELECT v.brand || ' ' || v.model || ' (' || v.plate_number || ')', COUNT(f.id),
			COUNT(f.id) FILTER (WHERE f.is_anomaly = true),
			CASE WHEN COUNT(f.id) > 0 THEN (COUNT(f.id) FILTER (WHERE f.is_anomaly = true) * 100 / COUNT(f.id)) ELSE 0 END as score
		FROM vehicles v JOIN fuel_records f ON v.id = f.vehicle_id WHERE f.created_at BETWEEN $1 AND $2%s
		GROUP BY v.id, v.brand, v.model, v.plate_number HAVING COUNT(f.id) FILTER (WHERE f.is_anomaly = true) > 0 ORDER BY score DESC
	`, tcond("v"))
	rRows, _ := db.Query(ctx, queryRisk, startDate, endDate)
	defer rRows.Close()
	for rRows.Next() {
		var r models.FleetRiskStat
		rRows.Scan(&r.VehicleName, &r.TotalRefuels, &r.Anomalies, &r.RiskScore)
		stats.FleetRisk = append(stats.FleetRisk, r)
	}

	// 4. ЗАБЕЗПЕЧЕНІСТЬ ПІДРОЗДІЛІВ
	// unitFilter вже починається з " WHERE" (або "") — треба акуратно долучити tenant
	uWhere := unitFilter
	if uWhere == "" && tcond("u") != "" {
		uWhere = " WHERE 1=1" + tcond("u")
	} else if uWhere != "" {
		uWhere += tcond("u")
	}
	queryReadiness := fmt.Sprintf(`
		SELECT u.name, COUNT(r.id), COUNT(r.id) FILTER (WHERE r.quantity >= r.min_quantity AND r.quantity > 0),
			CASE WHEN COUNT(r.id) > 0 THEN (COUNT(r.id) FILTER (WHERE r.quantity >= r.min_quantity AND r.quantity > 0) * 100 / COUNT(r.id)) ELSE 0 END as score
		FROM units u LEFT JOIN resources r ON u.id = r.unit_id AND r.condition != 'WRITTEN_OFF'
		%s
		GROUP BY u.id, u.name HAVING COUNT(r.id) > 0 ORDER BY score ASC
	`, uWhere)
	uRows, _ := db.Query(ctx, queryReadiness)
	defer uRows.Close()
	for uRows.Next() {
		var u models.UnitReadinessStat
		uRows.Scan(&u.UnitName, &u.TotalResources, &u.ReadyResources, &u.ReadinessScore)
		stats.UnitReadiness = append(stats.UnitReadiness, u)
	}

	// 5. КАРДІОЛІНІЯ ГСМ
	queryFuelCardio := fmt.Sprintf(`
		SELECT TO_CHAR(DATE(created_at), 'DD.MM'), SUM(liters), COUNT(id) FILTER (WHERE is_anomaly = true)
		FROM fuel_records WHERE created_at BETWEEN $1 AND $2%s GROUP BY DATE(created_at) ORDER BY DATE(created_at)
	`, tcond(""))
	fRows, _ := db.Query(ctx, queryFuelCardio, startDate, endDate)
	defer fRows.Close()
	for fRows.Next() {
		var f models.FuelMonthlyStat
		fRows.Scan(&f.Month, &f.TotalLiters, &f.Anomalies)
		stats.FuelHistory = append(stats.FuelHistory, f)
	}

	// 6. ПРОГНОЗ ТО
	queryMaint := fmt.Sprintf(`
		SELECT 
			v.brand || ' ' || v.plate_number,
			COALESCE((SELECT MAX(odometer_km) FROM fuel_records WHERE vehicle_id = v.id), v.last_maintenance_odometer) as current_odo,
			(v.last_maintenance_odometer + v.maintenance_interval_km) as next_maint,
			(v.last_maintenance_odometer + v.maintenance_interval_km) - COALESCE((SELECT MAX(odometer_km) FROM fuel_records WHERE vehicle_id = v.id), v.last_maintenance_odometer) as km_left
		FROM vehicles v WHERE v.status = 'ACTIVE' AND v.maintenance_interval_km > 0%s
		ORDER BY km_left ASC
	`, tcond("v"))
	mRows, _ := db.Query(ctx, queryMaint)
	defer mRows.Close()
	for mRows.Next() {
		var m models.MaintenancePredictStat
		mRows.Scan(&m.VehicleName, &m.CurrentOdo, &m.NextMaint, &m.KmLeft)
		stats.MaintenancePredict = append(stats.MaintenancePredict, m)
	}

	// ====================================================================
	// 7. ЕФЕКТИВНІСТЬ ПІДРЯДНИКІВ (SLA) - Enterprise SQL Logic
	// ====================================================================

	// 7.1. Рахуємо статистику по виконаних завданнях (Середній час, Найшвидший час, OTD)
	// Якщо в базі немає колонки deadline, ми розумно припускаємо базовий SLA у 7 днів (created_at + INTERVAL '7 days')
	querySLA := fmt.Sprintf(`
		SELECT 
			COALESCE(AVG(EXTRACT(EPOCH FROM (completed_at - created_at))/86400), 0) as avg_days,
			COUNT(id) as completed_count,
			COALESCE(MIN(EXTRACT(EPOCH FROM (completed_at - created_at))/86400), 0) as fastest_days,
			COALESCE(
				ROUND(
					(COUNT(id) FILTER (WHERE completed_at <= COALESCE(deadline, created_at + INTERVAL '7 days'))::numeric 
					/ NULLIF(COUNT(id), 0)) * 100
				), 0
			) as otd_percentage
		FROM CONTRACTOR_requests 
		WHERE status = 'COMPLETED' AND completed_at IS NOT NULL AND created_at BETWEEN $1 AND $2 %s%s
	`, volFilter, tcond(""))

	var otd float64
	errSLA := db.QueryRow(ctx, querySLA, startDate, endDate).Scan(
		&stats.CONTRACTORSLA.AverageDays,
		&stats.CONTRACTORSLA.CompletedCount,
		&stats.CONTRACTORSLA.FastestDays,
		&otd,
	)
	if errSLA == nil {
		stats.CONTRACTORSLA.OTDPercentage = int(otd)
	}

	// 7.2. Рахуємо прострочені завдання (ті, що досі в роботі, але дедлайн вже минув)
	queryOverdue := fmt.Sprintf(`
		SELECT COUNT(id) 
		FROM CONTRACTOR_requests 
		WHERE status IN ('OPEN', 'IN_PROGRESS', 'TAKEN') 
		  AND COALESCE(deadline, created_at + INTERVAL '7 days') < NOW() %s%s
	`, volFilter, tcond(""))

	_ = db.QueryRow(ctx, queryOverdue).Scan(&stats.CONTRACTORSLA.OverdueCount)
	// 8. ВИТРАТИ НА РЕМОНТИ
	queryTCO := fmt.Sprintf(`
		SELECT COALESCE(v.brand, 'Інше'), SUM(m.cost_amount) as total_cost
		FROM maintenance_records m JOIN vehicles v ON m.vehicle_id = v.id
		WHERE m.created_at BETWEEN $1 AND $2%s GROUP BY v.brand ORDER BY total_cost DESC
	`, tcond("v"))
	tcoRows, _ := db.Query(ctx, queryTCO, startDate, endDate)
	defer tcoRows.Close()
	for tcoRows.Next() {
		var t models.FleetTCOStat
		tcoRows.Scan(&t.VehicleBrand, &t.TotalCost)
		stats.FleetTCO = append(stats.FleetTCO, t)
	}

	// 9. ВОЛОНТЕРСЬКА ВОРОНКА
	queryFunnel := fmt.Sprintf(`
		SELECT status, COUNT(id) FROM CONTRACTOR_requests 
		WHERE created_at BETWEEN $1 AND $2 %s%s
		GROUP BY status
	`, volFilter, tcond(""))
	vRows, _ := db.Query(ctx, queryFunnel, startDate, endDate)
	defer vRows.Close()
	for vRows.Next() {
		var v models.CONTRACTORRequestStat
		vRows.Scan(&v.Status, &v.Count)
		stats.CONTRACTORFunnel = append(stats.CONTRACTORFunnel, v)
	}

	// 10. ДИНАМІКА ВОЛОНТЕРСЬКИХ ЗАЯВОК
	queryTimeline := fmt.Sprintf(`
		SELECT TO_CHAR(DATE(created_at), 'DD.MM'), COUNT(id) 
		FROM CONTRACTOR_requests 
		WHERE created_at BETWEEN $1 AND $2 %s%s
		GROUP BY DATE(created_at) 
		ORDER BY DATE(created_at)
	`, volFilter, tcond(""))
	tRows, _ := db.Query(ctx, queryTimeline, startDate, endDate)
	defer tRows.Close()
	for tRows.Next() {
		var vt models.CONTRACTORTimelineStat
		tRows.Scan(&vt.Date, &vt.Count)
		stats.CONTRACTORTimeline = append(stats.CONTRACTORTimeline, vt)
	}

	// 11. СПИСОК ДЕФІЦИТУ ДЛЯ SMART-ПОПОВНЕННЯ (З урахуванням майна в дорозі)
	queryDeficit := fmt.Sprintf(`
		WITH PendingOrders AS (
			-- Рахуємо скільки майна ВЖЕ замовлено (висить у відкритих заявках)
			SELECT resource_id, SUM(quantity) as pending_qty
			FROM supply_requests
			WHERE status IN ('OPEN', 'IN_PROGRESS', 'APPROVED')%s 
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
		  AND r.condition != 'WRITTEN_OFF' %s%s
	`, tcond(""), resFilterPrefix, tcond("r"))

	dRows, _ := db.Query(ctx, queryDeficit)
	defer dRows.Close()
	for dRows.Next() {
		var d models.DeficitResource
		dRows.Scan(&d.ID, &d.Name, &d.Current, &d.Min, &d.Needed)
		stats.DeficitResources = append(stats.DeficitResources, d)
	}

	// 12. НАЙЗАВАНТАЖЕНІШІ СКЛАДИ (ТОП-5)
	queryWarehouseLoad := fmt.Sprintf(`
		SELECT COALESCE(w.name, 'Без складу / В дорозі'), SUM(r.quantity) as total_items
		FROM resources r
		LEFT JOIN warehouses w ON r.warehouse_id = w.id
		WHERE r.condition != 'WRITTEN_OFF' AND r.quantity > 0 %s%s
		GROUP BY w.id, w.name
		ORDER BY total_items DESC
		LIMIT 5
	`, resFilterPrefix, tcond("r"))

	wlRows, _ := db.Query(ctx, queryWarehouseLoad)
	defer wlRows.Close()
	for wlRows.Next() {
		var wl models.WarehouseLoadStat
		wlRows.Scan(&wl.WarehouseName, &wl.TotalItems)
		stats.WarehouseLoad = append(stats.WarehouseLoad, wl)
	}

	// 13. ТОП-5 ЗАТРЕБУВАНИХ РЕСУРСІВ
	queryTopResources := fmt.Sprintf(`
		SELECT r.name, SUM(sr.quantity) as total_ordered
		FROM supply_requests sr
		JOIN resources r ON sr.resource_id = r.id
		WHERE sr.status = 'APPROVED' AND sr.created_at BETWEEN $1 AND $2 %s%s
		GROUP BY r.id, r.name
		ORDER BY total_ordered DESC
		LIMIT 5
	`, resFilterPrefix, tcond("r"))

	trRows, _ := db.Query(ctx, queryTopResources, startDate, endDate)
	defer trRows.Close()
	for trRows.Next() {
		var tr models.TopResourceStat
		trRows.Scan(&tr.ResourceName, &tr.TotalOrdered)
		stats.TopResources = append(stats.TopResources, tr)
	}

	return &stats, nil
}

// НОВА ФУНКЦІЯ: Обробка вибраних логістом позицій
func (r *AnalyticsRepository) ProcessSmartReplenish(ctx context.Context, db DBExecutor, req models.SmartReplenishRequest, userID string) (int, error) {
	count := 0
	tid := TenantFromCtx(ctx)

	for _, item := range req.Items {
		if item.Target == "WAREHOUSE" {
			// Створюємо офіційну заявку на забезпечення (на склад)
			_, err := db.Exec(ctx, `
				INSERT INTO supply_requests (resource_id, quantity, status, created_by, comment, created_at, tenant_id)
				VALUES ($1, $2, 'OPEN', $3, 'Автоматичне замовлення через Smart-модуль', NOW(),
					COALESCE(NULLIF($4, '')::uuid, (SELECT tenant_id FROM users WHERE id = $3)))
			`, item.ResourceID, item.Quantity, userID, tid)
			if err == nil {
				count++
			}
		} else if item.Target == "CONTRACTOR" {
			// Створюємо запит для волонтерів
			title := "Потреба: " + item.Name
			desc := fmt.Sprintf("Автоматично сформована потреба для підрозділу на %s (Кількість: %d)", item.Name, item.Quantity)

			_, err := db.Exec(ctx, `
				INSERT INTO CONTRACTOR_requests (created_by, title, description, status, created_at, tenant_id)
				VALUES ($1, $2, $3, 'OPEN', NOW(),
					COALESCE(NULLIF($4, '')::uuid, (SELECT tenant_id FROM users WHERE id = $1)))
			`, userID, title, desc, tid)
			if err == nil {
				count++
			}
		}
	}

	return count, nil
}

// --- СТРУКТУРИ ДЛЯ ЕКСПОРТУ ---

type ExportInventoryRow struct {
	Category  string
	ItemName  string
	UnitName  string // Філія / Підрозділ
	Warehouse string
	Quantity  int
	UnitType  string
	Condition string
}

type ExportFuelRow struct {
	Date       time.Time
	Vehicle    string
	Plate      string
	RecordType string
	Liters     float64
	Driver     string
}

func (r *AnalyticsRepository) GetInventoryForExport(ctx context.Context, db DBExecutor, unitID *int) ([]ExportInventoryRow, error) {
	tcond := tcondBuilder(ctx)
	query := fmt.Sprintf(`
		SELECT 
			c.name as category, 
			r.name as item_name, 
			COALESCE(u.name, 'Без підрозділу') as unit_name, 
			COALESCE(w.name, 'Без складу') as warehouse, 
			r.quantity, 
			r.unit_type, 
			r.condition
		FROM resources r
		JOIN resource_categories c ON r.category_id = c.id
		LEFT JOIN units u ON r.unit_id = u.id
		LEFT JOIN warehouses w ON r.warehouse_id = w.id
		WHERE ($1::int IS NULL OR r.unit_id = $1)%s
		ORDER BY u.name, c.name, r.name;
	`, tcond("r"))

	rows, err := db.Query(ctx, query, unitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []ExportInventoryRow
	for rows.Next() {
		var row ExportInventoryRow
		// Тепер Scan ніколи не впаде, бо ми гарантовано отримаємо текст
		if err := rows.Scan(&row.Category, &row.ItemName, &row.UnitName, &row.Warehouse, &row.Quantity, &row.UnitType, &row.Condition); err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	return result, nil
}

// GetFuelForExport витягує історію пального за період
func (r *AnalyticsRepository) GetFuelForExport(ctx context.Context, db DBExecutor, startDate, endDate time.Time) ([]ExportFuelRow, error) {
	tcond := tcondBuilder(ctx)
	query := fmt.Sprintf(`
		SELECT 
			f.created_at, v.brand || ' ' || COALESCE(v.model, ''), v.plate_number, 
			f.record_type, f.liters, COALESCE(u.full_name, 'Невідомо')
		FROM fuel_records f
		JOIN vehicles v ON f.vehicle_id = v.id
		LEFT JOIN users u ON f.created_by = u.id
		WHERE f.created_at >= $1 AND f.created_at <= $2%s
		ORDER BY f.created_at DESC;
	`, tcond("f"))

	rows, err := db.Query(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []ExportFuelRow
	for rows.Next() {
		var row ExportFuelRow
		if err := rows.Scan(&row.Date, &row.Vehicle, &row.Plate, &row.RecordType, &row.Liters, &row.Driver); err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	return result, nil
}

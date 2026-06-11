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

	unitFilter := ""
	validUnitID := ""

	if unitID != "" {
		if parsedUnitID, err := strconv.Atoi(unitID); err == nil {
			validUnitID = strconv.Itoa(parsedUnitID)
			unitFilter = fmt.Sprintf(" WHERE u.id = %s", validUnitID)
		}
	}

	resourceUnitFilter := func(alias string) string {
		if validUnitID == "" {
			return ""
		}
		prefix := ""
		if alias != "" {
			prefix = alias + "."
		}
		tenantExpr := "tenant_id"
		if alias != "" {
			tenantExpr = alias + ".tenant_id"
		}
		return fmt.Sprintf(`
		 AND (
			%sunit_id = %s
			OR %swarehouse_id IN (
				SELECT w.id
				FROM warehouses w
				WHERE w.unit_id = %s
				  AND w.tenant_id = %s
			)
		)`, prefix, validUnitID, prefix, validUnitID, tenantExpr)
	}

	vehicleUnitFilter := func(alias string) string {
		if validUnitID == "" {
			return ""
		}
		return fmt.Sprintf(`
		 AND EXISTS (
			SELECT 1
			FROM warehouses w
			WHERE (w.id = %s.current_warehouse_id OR w.id = %s.home_warehouse_id)
			  AND w.unit_id = %s
			  AND w.tenant_id = %s.tenant_id
		)`, alias, alias, validUnitID, alias)
	}

	fuelUnitFilter := func(alias string) string {
		if validUnitID == "" {
			return ""
		}
		return fmt.Sprintf(`
		 AND EXISTS (
			SELECT 1
			FROM vehicles v
			JOIN warehouses w
			  ON (w.id = v.current_warehouse_id OR w.id = v.home_warehouse_id)
			 AND w.tenant_id = v.tenant_id
			WHERE v.id = %s.vehicle_id
			  AND v.tenant_id = %s.tenant_id
			  AND w.unit_id = %s
		)`, alias, alias, validUnitID)
	}

	supplyUnitFilter := func(alias string) string {
		if validUnitID == "" {
			return ""
		}
		return fmt.Sprintf(`
		 AND (
			EXISTS (
				SELECT 1
				FROM warehouses w
				WHERE w.id = %s.target_warehouse_id
				  AND w.unit_id = %s
				  AND w.tenant_id = %s.tenant_id
			)
			OR EXISTS (
				SELECT 1
				FROM resources res
				LEFT JOIN warehouses w
				  ON w.id = res.warehouse_id
				 AND w.tenant_id = res.tenant_id
				WHERE res.id = %s.resource_id
				  AND (res.unit_id = %s OR w.unit_id = %s)
				  AND res.tenant_id = %s.tenant_id
			)
		)`, alias, validUnitID, alias, alias, validUnitID, validUnitID, alias)
	}

	contractorUnitFilter := func(alias string) string {
		if validUnitID == "" {
			return ""
		}
		return fmt.Sprintf(`
		 AND (
			%s.unit_id = %s
			OR EXISTS (
				SELECT 1
				FROM warehouses w
				WHERE w.id = %s.target_warehouse_id
				  AND w.unit_id = %s
				  AND w.tenant_id = %s.tenant_id
			)
		)`, alias, validUnitID, alias, validUnitID, alias)
	}

	// 1. ТОП Метрики
	queryMetrics := fmt.Sprintf(`
		SELECT 
			(SELECT COUNT(*) FROM vehicles v WHERE v.status IN ('ACTIVE', 'ON_MISSION')%s%s),
			(SELECT COUNT(*) FROM resources r WHERE r.quantity < r.min_quantity AND r.min_quantity > 0 AND r.condition != 'WRITTEN_OFF'%s%s),
			(SELECT COUNT(*) FROM fuel_records fr WHERE fr.is_anomaly = true AND fr.created_at BETWEEN $1 AND $2%s%s)
	`, vehicleUnitFilter("v"), tcond("v"), resourceUnitFilter("r"), tcond("r"), fuelUnitFilter("fr"), tcond("fr"))
	db.QueryRow(ctx, queryMetrics, startDate, endDate).Scan(&stats.ActiveVehicles, &stats.CriticalResources, &stats.FuelAnomalies)

	// ====================================================================
	// ЖИТТЄВИЙ ЦИКЛ (Виправлено під статус APPROVED)
	// ====================================================================
	queryWrittenOff := fmt.Sprintf(`SELECT COUNT(*) FROM resources r WHERE r.condition = 'WRITTEN_OFF' %s%s`, resourceUnitFilter("r"), tcond("r"))
	db.QueryRow(ctx, queryWrittenOff).Scan(&stats.WrittenOffResources)

	queryCompletedReqs := fmt.Sprintf(`
		SELECT COUNT(*) FROM supply_requests sr
		WHERE sr.status = 'COMPLETED' AND sr.updated_at BETWEEN $1 AND $2%s%s
	`, supplyUnitFilter("sr"), tcond("sr"))
	db.QueryRow(ctx, queryCompletedReqs, startDate, endDate).Scan(&stats.CompletedRequests)

	db.QueryRow(ctx, fmt.Sprintf(`SELECT COUNT(*) FROM vehicles v WHERE v.status = 'IN_REPAIR'%s%s`, vehicleUnitFilter("v"), tcond("v"))).Scan(&stats.InRepairVehicles)
	db.QueryRow(ctx, fmt.Sprintf(`SELECT COUNT(*) FROM vehicles v WHERE v.status = 'INACTIVE'%s%s`, vehicleUnitFilter("v"), tcond("v"))).Scan(&stats.InactiveVehicles)

	// ====================================================================
	// 2. ПРОГНОЗ ВИЧЕРПАННЯ (Виправлено під статус APPROVED)
	// ====================================================================
	queryPredict := fmt.Sprintf(`
		WITH consumption AS (
			SELECT
				sr.resource_id,
				sr.resource_name,
				sr.target_warehouse_id,
				SUM(sr.quantity) as consumed
			FROM supply_requests sr
			WHERE sr.updated_at BETWEEN $1 AND $2
			  AND sr.status = 'COMPLETED'%s%s
			GROUP BY sr.resource_id, sr.resource_name, sr.target_warehouse_id
		)
		SELECT r.name, r.quantity, r.min_quantity, c.consumed / NULLIF(EXTRACT(EPOCH FROM ($2 - $1))/86400, 0) as daily_burn,
			(r.quantity / NULLIF(c.consumed / NULLIF(EXTRACT(EPOCH FROM ($2 - $1))/86400, 0), 0))::int as days_left
		FROM resources r
		JOIN consumption c
		  ON c.resource_id = r.id
		  OR (
			c.resource_id IS NULL
			AND c.target_warehouse_id IS NOT NULL
			AND c.resource_name = r.name
			AND c.target_warehouse_id = r.warehouse_id
		  )
		WHERE r.condition != 'WRITTEN_OFF' AND c.consumed > 0 %s%s 
		GROUP BY r.id, r.name, r.quantity, r.min_quantity, c.consumed
		ORDER BY days_left ASC
	`, supplyUnitFilter("sr"), tcond("sr"), resourceUnitFilter("r"), tcond("r"))
	pRows, _ := db.Query(ctx, queryPredict, startDate, endDate)
	defer pRows.Close()
	for pRows.Next() {
		var p models.PredictStat
		pRows.Scan(&p.ResourceName, &p.CurrentStock, &p.MinQuantity, &p.DailyBurnRate, &p.DaysLeft)
		stats.PredictiveBurnRate = append(stats.PredictiveBurnRate, p)
	}

	// 3. АНТИКОРУПЦІЙНИЙ ІНДЕКС
	queryRisk := fmt.Sprintf(`
		SELECT TRIM(CONCAT(COALESCE(v.brand, 'Інше'), ' ', COALESCE(v.model, ''), ' (', v.plate_number, ')')), COUNT(f.id),
			COUNT(f.id) FILTER (WHERE f.is_anomaly = true),
			CASE WHEN COUNT(f.id) > 0 THEN ROUND((COUNT(f.id) FILTER (WHERE f.is_anomaly = true)::numeric * 100) / COUNT(f.id))::int ELSE 0 END as score
		FROM vehicles v JOIN fuel_records f ON v.id = f.vehicle_id AND f.tenant_id = v.tenant_id WHERE f.created_at BETWEEN $1 AND $2%s%s
		GROUP BY v.id, v.brand, v.model, v.plate_number
		HAVING COUNT(f.id) FILTER (WHERE f.is_anomaly = true) > 0
		ORDER BY score DESC, COUNT(f.id) FILTER (WHERE f.is_anomaly = true) DESC, COUNT(f.id) DESC
	`, vehicleUnitFilter("v"), tcond("v"))
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
		SELECT u.name, COUNT(r.id), COUNT(r.id) FILTER (WHERE r.quantity >= r.min_quantity),
			CASE WHEN COUNT(r.id) > 0 THEN (COUNT(r.id) FILTER (WHERE r.quantity >= r.min_quantity) * 100 / COUNT(r.id)) ELSE 0 END as score
		FROM units u LEFT JOIN resources r
		  ON (
			u.id = r.unit_id
			OR EXISTS (
				SELECT 1
				FROM warehouses w
				WHERE w.id = r.warehouse_id
				  AND w.unit_id = u.id
				  AND w.tenant_id = r.tenant_id
			)
		  )
		  AND r.condition != 'WRITTEN_OFF'%s
		%s
		GROUP BY u.id, u.name HAVING COUNT(r.id) > 0 ORDER BY score ASC
	`, tcond("r"), uWhere)
	uRows, _ := db.Query(ctx, queryReadiness)
	defer uRows.Close()
	for uRows.Next() {
		var u models.UnitReadinessStat
		uRows.Scan(&u.UnitName, &u.TotalResources, &u.ReadyResources, &u.ReadinessScore)
		stats.UnitReadiness = append(stats.UnitReadiness, u)
	}

	// 5. КАРДІОЛІНІЯ ГСМ
	queryFuelCardio := fmt.Sprintf(`
		SELECT TO_CHAR(DATE(fr.created_at), 'DD.MM'), SUM(fr.liters), COUNT(fr.id) FILTER (WHERE fr.is_anomaly = true)
		FROM fuel_records fr
		WHERE fr.record_type IN ('EXPENSE', 'CONSUMPTION') AND fr.created_at BETWEEN $1 AND $2%s%s
		GROUP BY DATE(fr.created_at) ORDER BY DATE(fr.created_at)
	`, fuelUnitFilter("fr"), tcond("fr"))
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
			TRIM(CONCAT(COALESCE(v.brand, 'Інше'), ' ', v.plate_number)),
			COALESCE((SELECT MAX(odometer_km) FROM fuel_records WHERE vehicle_id = v.id AND tenant_id = v.tenant_id), COALESCE(v.last_maintenance_odometer, 0)) as current_odo,
			(COALESCE(v.last_maintenance_odometer, 0) + COALESCE(NULLIF(v.maintenance_interval_km, 0), 10000)) as next_maint,
			(COALESCE(v.last_maintenance_odometer, 0) + COALESCE(NULLIF(v.maintenance_interval_km, 0), 10000)) - COALESCE((SELECT MAX(odometer_km) FROM fuel_records WHERE vehicle_id = v.id AND tenant_id = v.tenant_id), COALESCE(v.last_maintenance_odometer, 0)) as km_left
		FROM vehicles v WHERE v.status != 'INACTIVE'%s%s
		ORDER BY km_left ASC
	`, vehicleUnitFilter("v"), tcond("v"))
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
			COALESCE(AVG(EXTRACT(EPOCH FROM (cr.completed_at - cr.created_at))/86400), 0) as avg_days,
			COUNT(cr.id) as completed_count,
			COALESCE(MIN(EXTRACT(EPOCH FROM (cr.completed_at - cr.created_at))/86400), 0) as fastest_days,
			COALESCE(
				ROUND(
					(COUNT(cr.id) FILTER (WHERE cr.completed_at <= COALESCE(cr.deadline, cr.created_at + INTERVAL '7 days'))::numeric 
					/ NULLIF(COUNT(cr.id), 0)) * 100
				), 0
			) as otd_percentage
		FROM contractor_requests cr
		WHERE cr.status IN ('COMPLETED', 'ACCEPTED')
		  AND cr.completed_at IS NOT NULL
		  AND cr.completed_at BETWEEN $1 AND $2%s%s
	`, contractorUnitFilter("cr"), tcond("cr"))

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
		SELECT COUNT(cr.id) 
		FROM contractor_requests cr
		WHERE cr.status IN ('OPEN', 'IN_PROGRESS', 'TAKEN') 
		  AND cr.created_at <= $2
		  AND COALESCE(cr.deadline, cr.created_at + INTERVAL '7 days') < $2%s%s
	`, contractorUnitFilter("cr"), tcond("cr"))

	_ = db.QueryRow(ctx, queryOverdue, startDate, endDate).Scan(&stats.CONTRACTORSLA.OverdueCount)
	// 8. TCO автопарку: ремонт + фактична вартість дозаправок у рейсах
	queryTCO := fmt.Sprintf(`
		WITH vehicle_scope AS (
			SELECT
				v.id,
				v.tenant_id,
				TRIM(CONCAT(COALESCE(v.brand, 'Інше'), ' ', COALESCE(v.model, ''), ' (', v.plate_number, ')')) AS vehicle_label
			FROM vehicles v
			WHERE 1=1%s%s
		),
		maintenance_costs AS (
			SELECT m.vehicle_id, COALESCE(SUM(m.cost_amount), 0) AS total_cost
			FROM maintenance_records m
			JOIN vehicle_scope vs ON vs.id = m.vehicle_id
			WHERE m.created_at BETWEEN $1 AND $2
			GROUP BY m.vehicle_id
		),
		fuel_costs AS (
			SELECT sr.vehicle_id, COALESCE(SUM(sr.cost_uah), 0) AS total_cost
			FROM shipment_refuels sr
			JOIN vehicle_scope vs ON vs.id = sr.vehicle_id AND vs.tenant_id = sr.tenant_id
			WHERE sr.created_at BETWEEN $1 AND $2
			  AND sr.cost_uah IS NOT NULL
			GROUP BY sr.vehicle_id
		)
		SELECT
			vs.vehicle_label,
			COALESCE(mc.total_cost, 0) + COALESCE(fc.total_cost, 0) AS total_cost
		FROM vehicle_scope vs
		LEFT JOIN maintenance_costs mc ON mc.vehicle_id = vs.id
		LEFT JOIN fuel_costs fc ON fc.vehicle_id = vs.id
		WHERE COALESCE(mc.total_cost, 0) + COALESCE(fc.total_cost, 0) > 0
		ORDER BY total_cost DESC, vs.vehicle_label
		LIMIT 5
	`, vehicleUnitFilter("v"), tcond("v"))
	tcoRows, _ := db.Query(ctx, queryTCO, startDate, endDate)
	defer tcoRows.Close()
	for tcoRows.Next() {
		var t models.FleetTCOStat
		tcoRows.Scan(&t.VehicleBrand, &t.TotalCost)
		stats.FleetTCO = append(stats.FleetTCO, t)
	}

	// 9. ВОЛОНТЕРСЬКА ВОРОНКА
	queryFunnel := fmt.Sprintf(`
		SELECT cr.status, COUNT(cr.id) FROM contractor_requests cr
		WHERE cr.created_at BETWEEN $1 AND $2%s%s
		GROUP BY cr.status
	`, contractorUnitFilter("cr"), tcond("cr"))
	vRows, _ := db.Query(ctx, queryFunnel, startDate, endDate)
	defer vRows.Close()
	for vRows.Next() {
		var v models.CONTRACTORRequestStat
		vRows.Scan(&v.Status, &v.Count)
		stats.CONTRACTORFunnel = append(stats.CONTRACTORFunnel, v)
	}

	// 10. ДИНАМІКА ВОЛОНТЕРСЬКИХ ЗАЯВОК
	queryTimeline := fmt.Sprintf(`
		SELECT TO_CHAR(DATE(cr.created_at), 'DD.MM'), COUNT(cr.id) 
		FROM contractor_requests cr
		WHERE cr.created_at BETWEEN $1 AND $2%s%s
		GROUP BY DATE(cr.created_at) 
		ORDER BY DATE(cr.created_at)
	`, contractorUnitFilter("cr"), tcond("cr"))
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
			-- Рахуємо, скільки майна вже замовлено або знаходиться в дорозі
			SELECT sr.resource_id, sr.resource_name, sr.target_warehouse_id, SUM(sr.quantity) as pending_qty
			FROM supply_requests sr
			WHERE sr.status IN ('PENDING', 'ESCALATED', 'APPROVED', 'LOADING', 'DISPATCHED', 'OPEN', 'IN_PROGRESS')%s%s
			GROUP BY sr.resource_id, sr.resource_name, sr.target_warehouse_id
		)
		SELECT 
			r.id, 
			r.name, 
			r.quantity, 
			r.min_quantity, 
			-- Формула: (Мінімум * 2) - (Фактичний залишок + Вже замовлено)
			(r.min_quantity * 2 - (r.quantity + COALESCE(SUM(p.pending_qty), 0))) as needed
		FROM resources r
		LEFT JOIN PendingOrders p
		  ON p.resource_id = r.id
		  OR (
			p.resource_id IS NULL
			AND p.target_warehouse_id IS NOT NULL
			AND p.resource_name = r.name
			AND p.target_warehouse_id = r.warehouse_id
		  )
		-- Показуємо тільки те, де ФАКТ + В ДОРОЗІ все ще менше або дорівнює мінімуму
		WHERE r.condition != 'WRITTEN_OFF' AND r.min_quantity > 0 %s%s
		GROUP BY r.id, r.name, r.quantity, r.min_quantity
		HAVING (r.quantity + COALESCE(SUM(p.pending_qty), 0)) <= r.min_quantity
	`, supplyUnitFilter("sr"), tcond("sr"), resourceUnitFilter("r"), tcond("r"))

	dRows, _ := db.Query(ctx, queryDeficit)
	defer dRows.Close()
	for dRows.Next() {
		var d models.DeficitResource
		dRows.Scan(&d.ID, &d.Name, &d.Current, &d.Min, &d.Needed)
		stats.DeficitResources = append(stats.DeficitResources, d)
	}
	stats.CriticalResources = len(stats.DeficitResources)

	// 12. НАЙЗАВАНТАЖЕНІШІ СКЛАДИ (ТОП-5)
	queryWarehouseLoad := fmt.Sprintf(`
		SELECT COALESCE(w.name, 'Без складу / В дорозі'), SUM(r.quantity) as total_items
		FROM resources r
		LEFT JOIN warehouses w ON r.warehouse_id = w.id AND w.tenant_id = r.tenant_id
		WHERE r.condition != 'WRITTEN_OFF' AND r.quantity > 0 %s%s
		GROUP BY w.id, w.name
		ORDER BY total_items DESC
		LIMIT 5
	`, resourceUnitFilter("r"), tcond("r"))

	wlRows, _ := db.Query(ctx, queryWarehouseLoad)
	defer wlRows.Close()
	for wlRows.Next() {
		var wl models.WarehouseLoadStat
		wlRows.Scan(&wl.WarehouseName, &wl.TotalItems)
		stats.WarehouseLoad = append(stats.WarehouseLoad, wl)
	}

	// 13. ТОП-5 ЗАТРЕБУВАНИХ РЕСУРСІВ
	queryTopResources := fmt.Sprintf(`
		SELECT COALESCE(NULLIF(sr.resource_name, ''), r.name, 'Невідомий ресурс'), SUM(sr.quantity) as total_ordered
		FROM supply_requests sr
		LEFT JOIN resources r ON sr.resource_id = r.id AND r.tenant_id = sr.tenant_id
		WHERE sr.status NOT IN ('REJECTED', 'CANCELED', 'CANCELLED') AND sr.created_at BETWEEN $1 AND $2%s%s
		GROUP BY COALESCE(NULLIF(sr.resource_name, ''), r.name, 'Невідомий ресурс')
		ORDER BY total_ordered DESC
		LIMIT 5
	`, supplyUnitFilter("sr"), tcond("sr"))

	trRows, _ := db.Query(ctx, queryTopResources, startDate, endDate)
	defer trRows.Close()
	for trRows.Next() {
		var tr models.TopResourceStat
		trRows.Scan(&tr.ResourceName, &tr.TotalOrdered)
		stats.TopResources = append(stats.TopResources, tr)
	}

	// 14. ЗАГАЛЬНА ВАРТІСТЬ ЗАЛИШКІВ
	queryInventoryValue := fmt.Sprintf(`
		SELECT COALESCE(SUM(r.quantity * r.unit_price), 0)
		FROM resources r
		WHERE r.condition != 'WRITTEN_OFF' AND r.unit_price > 0 %s%s
	`, resourceUnitFilter("r"), tcond("r"))
	db.QueryRow(ctx, queryInventoryValue).Scan(&stats.InventoryTotalValue)

	// 15. ВАРТІСТЬ СПИСАНОГО ЗА ПЕРІОД
	queryWriteOffValue := fmt.Sprintf(`
		SELECT COALESCE(SUM(r.quantity * r.unit_price), 0)
		FROM resources r
		WHERE r.condition = 'WRITTEN_OFF' AND r.unit_price > 0
		  AND r.updated_at BETWEEN $1 AND $2 %s%s
	`, resourceUnitFilter("r"), tcond("r"))
	db.QueryRow(ctx, queryWriteOffValue, startDate, endDate).Scan(&stats.WriteOffTotalValue)

	// 16. ВАРТІСТЬ ПО СКЛАДАХ (ТОП-5)
	queryWarehouseValue := fmt.Sprintf(`
		SELECT COALESCE(w.name, 'Без складу / В дорозі'), COALESCE(SUM(r.quantity * r.unit_price), 0) as total_value
		FROM resources r
		LEFT JOIN warehouses w ON r.warehouse_id = w.id AND w.tenant_id = r.tenant_id
		WHERE r.condition != 'WRITTEN_OFF' AND r.unit_price > 0 AND r.quantity > 0 %s%s
		GROUP BY w.id, w.name
		ORDER BY total_value DESC
		LIMIT 5
	`, resourceUnitFilter("r"), tcond("r"))

	wvRows, _ := db.Query(ctx, queryWarehouseValue)
	defer wvRows.Close()
	for wvRows.Next() {
		var wv models.WarehouseValueStat
		wvRows.Scan(&wv.WarehouseName, &wv.TotalValue)
		stats.WarehouseValueStats = append(stats.WarehouseValueStats, wv)
	}

	// 17. ТОП-5 НАЙДОРОЖЧИХ ПОЗИЦІЙ
	queryTopCostly := fmt.Sprintf(`
		SELECT r.name, COALESCE(SUM(r.quantity * r.unit_price), 0) as total_value, r.quantity, r.unit_price
		FROM resources r
		WHERE r.condition != 'WRITTEN_OFF' AND r.unit_price > 0 AND r.quantity > 0 %s%s
		GROUP BY r.id, r.name, r.quantity, r.unit_price
		ORDER BY total_value DESC
		LIMIT 5
	`, resourceUnitFilter("r"), tcond("r"))

	tcRows, _ := db.Query(ctx, queryTopCostly)
	defer tcRows.Close()
	for tcRows.Next() {
		var tc models.TopCostlyResourceStat
		tcRows.Scan(&tc.ResourceName, &tc.TotalValue, &tc.Quantity, &tc.UnitPrice)
		stats.TopCostlyResources = append(stats.TopCostlyResources, tc)
	}

	return &stats, nil
}

// НОВА ФУНКЦІЯ: Обробка вибраних логістом позицій
func (r *AnalyticsRepository) ProcessSmartReplenish(ctx context.Context, db DBExecutor, req models.SmartReplenishRequest, userID string) (int, error) {
	count := 0
	tid := TenantFromCtx(ctx)

	for _, item := range req.Items {
		if item.Target == "WAREHOUSE" {
			// Створюємо автоматичну заявку одразу зі статусом APPROVED —
			// це системна дія адміна, не потребує ручного погодження.
			_, err := db.Exec(ctx, `
				INSERT INTO supply_requests (resource_id, resource_name, quantity, status, created_by, comment, created_at, tenant_id)
				VALUES ($1, $2, $3, 'APPROVED', $4, 'Автоматичне замовлення через Smart-модуль', NOW(),
					COALESCE(NULLIF($5, '')::uuid, (SELECT tenant_id FROM users WHERE id = $4)))
			`, item.ResourceID, item.Name, item.Quantity, userID, tid)
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
	Category   string
	ItemName   string
	UnitName   string // Філія / Підрозділ
	Warehouse  string
	Quantity   int
	UnitType   string
	Condition  string
	UnitPrice  float64
	TotalValue float64
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
			r.condition,
			r.unit_price,
			r.quantity * r.unit_price as total_value
		FROM resources r
		JOIN resource_categories c ON r.category_id = c.id AND c.tenant_id = r.tenant_id
		LEFT JOIN units u ON r.unit_id = u.id AND u.tenant_id = r.tenant_id
		LEFT JOIN warehouses w ON r.warehouse_id = w.id AND w.tenant_id = r.tenant_id
		WHERE ($1::int IS NULL OR r.unit_id = $1 OR w.unit_id = $1)%s
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
		if err := rows.Scan(&row.Category, &row.ItemName, &row.UnitName, &row.Warehouse, &row.Quantity, &row.UnitType, &row.Condition, &row.UnitPrice, &row.TotalValue); err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	return result, nil
}

// GetFuelForExport витягує історію пального за період
func (r *AnalyticsRepository) GetFuelForExport(ctx context.Context, db DBExecutor, startDate, endDate time.Time, unitID *int) ([]ExportFuelRow, error) {
	tcond := tcondBuilder(ctx)
	unitFilter := ""
	args := []interface{}{startDate, endDate}
	if unitID != nil {
		args = append(args, *unitID)
		unitFilter = fmt.Sprintf(`
		AND EXISTS (
			SELECT 1
			FROM warehouses w
			WHERE (w.id = v.current_warehouse_id OR w.id = v.home_warehouse_id)
			  AND w.unit_id = $%d
			  AND w.tenant_id = v.tenant_id
		)`, len(args))
	}
	query := fmt.Sprintf(`
		SELECT 
			f.created_at, v.brand || ' ' || COALESCE(v.model, ''), v.plate_number, 
			f.record_type, f.liters, COALESCE(u.full_name, 'Невідомо')
		FROM fuel_records f
		JOIN vehicles v ON f.vehicle_id = v.id AND v.tenant_id = f.tenant_id
		LEFT JOIN users u ON f.created_by = u.id
		WHERE f.created_at >= $1 AND f.created_at <= $2%s%s
		ORDER BY f.created_at DESC;
	`, unitFilter, tcond("f"))

	rows, err := db.Query(ctx, query, args...)
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

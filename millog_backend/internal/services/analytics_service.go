package services

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"Omnilog_backend/internal/models"
	"Omnilog_backend/internal/repositories"
)

type AnalyticsService struct {
	repo *repositories.AnalyticsRepository
	db   repositories.DBExecutor
}

func NewAnalyticsService(repo *repositories.AnalyticsRepository, db repositories.DBExecutor) *AnalyticsService {
	return &AnalyticsService{repo: repo, db: db}
}

func (s *AnalyticsService) GetDashboardAnalytics(ctx context.Context, start, end, unitID string) (*models.DashboardAnalytics, error) {
	return s.repo.GetDashboardStats(ctx, s.db, start, end, unitID)
}

// НОВА ФУНКЦІЯ: Приймає налаштування замовлення з фронтенду та ID користувача
func (s *AnalyticsService) RunSmartReplenish(ctx context.Context, req models.SmartReplenishRequest, userID string) (int, error) {
	return s.repo.ProcessSmartReplenish(ctx, s.db, req, userID)
}

// GenerateInventoryExcel формує XLSX з поточними залишками
func (s *AnalyticsService) GenerateInventoryExcel(ctx context.Context, unitID *int) ([]byte, error) {
	data, err := s.repo.GetInventoryForExport(ctx, s.db, unitID)
	if err != nil {
		return nil, fmt.Errorf("помилка отримання даних: %w", err)
	}

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Залишки на складах"
	f.SetSheetName("Sheet1", sheetName)

	// Стиль для шапки (жирний шрифт, сірий фон)
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#4F81BD"}, Pattern: 1},
	})

	headers := []string{"Філія / Підрозділ", "Склад", "Категорія", "Найменування майна", "Кількість", "Од. вим.", "Стан", "Ціна, грн", "Вартість, грн"}
	for col, header := range headers {
		colName, _ := excelize.ColumnNumberToName(col + 1)
		cell := fmt.Sprintf("%s1", colName)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
		f.SetColWidth(sheetName, colName, colName, 20) // Задаємо ширину колонок
	}
	f.SetColWidth(sheetName, "D", "D", 40) // Назва майна ширша

	// Заповнюємо дані
	for i, row := range data {
		rowIndex := i + 2 // Рядок 1 - це шапка

		// Переклад значень для бухгалтерії
		unitType := row.UnitType
		switch unitType {
		case "PCS":
			unitType = "шт"
		case "KIT":
			unitType = "компл"
		case "KG":
			unitType = "кг"
		case "L":
			unitType = "л"
		}

		condition := "Нове"
		if row.Condition == "USED" {
			condition = "Вживане"
		} else if row.Condition == "WRITTEN_OFF" {
			condition = "Списано"
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIndex), row.UnitName)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIndex), row.Warehouse)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIndex), row.Category)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIndex), row.ItemName)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowIndex), row.Quantity)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowIndex), unitType)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowIndex), condition)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowIndex), row.UnitPrice)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", rowIndex), row.TotalValue)
	}

	// Додаємо автофільтр на всі дані
	f.AutoFilter(sheetName, fmt.Sprintf("A1:I%d", len(data)+1), nil)

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GenerateFuelExcel формує XLSX зі звітом по пальному
func (s *AnalyticsService) GenerateFuelExcel(ctx context.Context, startDate, endDate time.Time, unitID *int) ([]byte, error) {
	data, err := s.repo.GetFuelForExport(ctx, s.db, startDate, endDate, unitID)
	if err != nil {
		return nil, fmt.Errorf("помилка отримання даних: %w", err)
	}

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Звіт по пальному"
	f.SetSheetName("Sheet1", sheetName)

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E26B0A"}, Pattern: 1},
	})

	headers := []string{"Дата та Час", "Транспортний засіб", "Держ. номер", "Тип операції", "Літри", "Відповідальний"}
	for col, header := range headers {
		colName, _ := excelize.ColumnNumberToName(col + 1)
		cell := fmt.Sprintf("%s1", colName)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
		f.SetColWidth(sheetName, colName, colName, 22)
	}

	loc, _ := time.LoadLocation("Europe/Kyiv")

	for i, row := range data {
		rowIndex := i + 2

		opType := "Витрата (Списання)"
		if row.RecordType == "REFUEL" {
			opType = "Заправка (Надходження)"
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIndex), row.Date.In(loc).Format("02.01.2006 15:04"))
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIndex), row.Vehicle)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIndex), row.Plate)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIndex), opType)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowIndex), row.Liters)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowIndex), row.Driver)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 🚀 PRO FEATURE #1: GetAdvancedKPIs повертає розширену аналітику з KPI
// Формат відповіді синхронізовано з компонентом KPIDashboard на фронтенді.
func (s *AnalyticsService) GetAdvancedKPIs(ctx context.Context, startDate, endDate string, unitID int64) (map[string]interface{}, error) {
	tenantID := repositories.TenantFromCtx(ctx)
	addTenantFilter := func(args *[]interface{}, alias string) string {
		if tenantID == "" {
			return ""
		}
		*args = append(*args, tenantID)
		column := "tenant_id"
		if alias != "" {
			column = alias + ".tenant_id"
		}
		return fmt.Sprintf(" AND %s = $%d::uuid", column, len(*args))
	}

	// ---------- SLA ----------
	// Заявка "вчасна", якщо її апрувнули за ≤ 24 години.
	slaArgs := []interface{}{startDate, endDate}
	slaFilter := ""
	if unitID > 0 {
		slaArgs = append(slaArgs, unitID)
		slaFilter += fmt.Sprintf(" AND w.unit_id = $%d", len(slaArgs))
	}
	slaFilter += addTenantFilter(&slaArgs, "sr")

	slaQuery := fmt.Sprintf(`
		WITH req AS (
			SELECT sr.id, sr.status, sr.created_at, sr.approved_at
			FROM supply_requests sr
			LEFT JOIN warehouses w ON w.id = sr.target_warehouse_id AND w.tenant_id = sr.tenant_id
			WHERE sr.created_at::date >= $1::date
			  AND sr.created_at::date <= $2::date%s
		)
		SELECT
			COUNT(*) AS total,
			COUNT(*) FILTER (
				WHERE approved_at IS NOT NULL
				  AND EXTRACT(EPOCH FROM (approved_at - created_at))/3600 <= 24
			) AS on_time,
			COALESCE(AVG(
				CASE WHEN approved_at IS NOT NULL
					THEN EXTRACT(EPOCH FROM (approved_at - created_at))/3600
				END
			), 0) AS avg_hours
		FROM req
	`, slaFilter)

	var totalReq, onTime int64
	var avgHours float64
	if err := s.db.QueryRow(ctx, slaQuery, slaArgs...).Scan(&totalReq, &onTime, &avgHours); err != nil {
		totalReq, onTime, avgHours = 0, 0, 0
	}
	onTimePercent := 100.0
	if totalReq > 0 {
		onTimePercent = float64(onTime) * 100.0 / float64(totalReq)
	}

	// ---------- TCO ----------
	// Парсимо дати для підрахунку попереднього періоду (для trend)
	parsedStart, errPS := time.Parse("2006-01-02", startDate)
	parsedEnd, errPE := time.Parse("2006-01-02", endDate)
	periodDays := 30
	if errPS == nil && errPE == nil {
		d := int(parsedEnd.Sub(parsedStart).Hours()/24) + 1
		if d > 0 {
			periodDays = d
		}
	}
	prevEnd := parsedStart.AddDate(0, 0, -1).Format("2006-01-02")
	prevStart := parsedStart.AddDate(0, 0, -periodDays).Format("2006-01-02")
	_ = errPS
	_ = errPE

	buildFuelCostArgs := func(sd, ed string) ([]interface{}, string) {
		a := []interface{}{sd, ed}
		f := addTenantFilter(&a, "sr")
		if unitID > 0 {
			a = append(a, unitID)
			f += `
			AND EXISTS (
				SELECT 1
				FROM vehicles v
				JOIN warehouses w
				  ON (w.id = v.current_warehouse_id OR w.id = v.home_warehouse_id)
				 AND w.tenant_id = v.tenant_id
				WHERE v.id = sr.vehicle_id
				  AND v.tenant_id = sr.tenant_id
				  AND w.unit_id = $` + fmt.Sprintf("%d", len(a)) + `
			)`
		}
		return a, f
	}

	queryFuelCost := func(sd, ed string) float64 {
		a, vuf := buildFuelCostArgs(sd, ed)
		q := fmt.Sprintf(`
			SELECT COALESCE(SUM(sr.cost_uah), 0)
			FROM shipment_refuels sr
			WHERE sr.created_at::date >= $1::date
			  AND sr.created_at::date <= $2::date
			  AND sr.cost_uah IS NOT NULL%s
		`, vuf)
		var cost float64
		_ = s.db.QueryRow(ctx, q, a...).Scan(&cost)
		return cost
	}

	totalFuelCost := queryFuelCost(startDate, endDate)
	prevFuelCost := queryFuelCost(prevStart, prevEnd)

	// Кількість відвантажених одиниць
	shipArgs := []interface{}{startDate, endDate}
	shipUnitFilter := ""
	if unitID > 0 {
		shipArgs = append(shipArgs, unitID)
		shipUnitFilter = fmt.Sprintf(" AND (wf.unit_id = $%d OR wt.unit_id = $%d)", len(shipArgs), len(shipArgs))
	}
	shipTenantFilter := addTenantFilter(&shipArgs, "s")
	shipQuery := fmt.Sprintf(`
		SELECT COALESCE(SUM(si.quantity), 0)
		FROM shipments s
		JOIN shipment_items si ON si.shipment_id = s.id
		LEFT JOIN warehouses wf ON wf.id = s.from_warehouse_id AND wf.tenant_id = s.tenant_id
		LEFT JOIN warehouses wt ON wt.id = s.to_warehouse_id AND wt.tenant_id = s.tenant_id
		WHERE s.created_at::date >= $1::date
		  AND s.created_at::date <= $2::date%s%s
	`, shipUnitFilter, shipTenantFilter)
	var unitsShipped int64
	_ = s.db.QueryRow(ctx, shipQuery, shipArgs...).Scan(&unitsShipped)

	costPerUnit := 0.0
	if unitsShipped > 0 {
		costPerUnit = totalFuelCost / float64(unitsShipped)
	}

	// Trend: порівнюємо поточний і попередній період (>5% різниця = up/down)
	tcoTrend := "stable"
	if prevFuelCost > 0 {
		changePct := (totalFuelCost - prevFuelCost) / prevFuelCost * 100
		if changePct > 5 {
			tcoTrend = "up"
		} else if changePct < -5 {
			tcoTrend = "down"
		}
	} else if totalFuelCost > 0 {
		tcoTrend = "up" // попередніх даних не було — зростання від нуля
	}

	// ---------- RISK ----------
	riskArgs := []interface{}{}
	riskFilter := ""
	if unitID > 0 {
		riskArgs = append(riskArgs, unitID)
		riskFilter += fmt.Sprintf(
			` AND (r.unit_id = $%d OR r.warehouse_id IN (SELECT id FROM warehouses WHERE unit_id = $%d AND tenant_id = r.tenant_id))`,
			len(riskArgs),
			len(riskArgs),
		)
	}
	riskFilter += addTenantFilter(&riskArgs, "r")
	riskQuery := fmt.Sprintf(`
		SELECT r.id, r.name, r.quantity, r.min_quantity
		FROM resources r
		WHERE r.min_quantity > 0%s
	`, riskFilter)

	rows, err := s.db.Query(ctx, riskQuery, riskArgs...)
	if err != nil {
		return nil, fmt.Errorf("kpi risk query: %w", err)
	}
	defer rows.Close()

	totalResources := 0
	atRisk := 0
	criticalNames := make([]string, 0)
	within7 := make([]string, 0)
	within14 := make([]string, 0)
	within30 := make([]string, 0)

	for rows.Next() {
		var id, name string
		var qty, minQty int
		if err := rows.Scan(&id, &name, &qty, &minQty); err != nil {
			continue
		}
		totalResources++
		if minQty <= 0 {
			continue
		}
		ratio := float64(qty) / float64(minQty)
		switch {
		case ratio < 0.3:
			within7 = append(within7, name)
			atRisk++
			if len(criticalNames) < 10 {
				criticalNames = append(criticalNames, name)
			}
		case ratio < 0.6:
			within14 = append(within14, name)
			atRisk++
		case ratio < 1.0:
			within30 = append(within30, name)
		}
	}

	criticalPercent := 0.0
	if totalResources > 0 {
		criticalPercent = float64(atRisk) * 100.0 / float64(totalResources)
	}

	// ---------- INVENTORY VALUE ----------
	invArgs := []interface{}{}
	invFilter := ""
	if unitID > 0 {
		invArgs = append(invArgs, unitID)
		invFilter += fmt.Sprintf(
			` AND (r.unit_id = $%d OR r.warehouse_id IN (SELECT id FROM warehouses WHERE unit_id = $%d AND tenant_id = r.tenant_id))`,
			len(invArgs),
			len(invArgs),
		)
	}
	invFilter += addTenantFilter(&invArgs, "r")
	invQuery := fmt.Sprintf(`
		SELECT COALESCE(SUM(r.quantity * COALESCE(r.unit_price, 0)), 0)
		FROM resources r
		WHERE r.condition != 'WRITTEN_OFF'%s
	`, invFilter)
	var inventoryTotalValue float64
	_ = s.db.QueryRow(ctx, invQuery, invArgs...).Scan(&inventoryTotalValue)

	return map[string]interface{}{
		"reporting_period": fmt.Sprintf("%s — %s", startDate, endDate),
		"sla": map[string]interface{}{
			"on_time_percent": onTimePercent,
			"total_requests":  totalReq,
			"on_time_count":   onTime,
			"late_count":      totalReq - onTime,
			"avg_delay_hours": avgHours,
		},
		"tco": map[string]interface{}{
			"total_fuel_cost":     totalFuelCost,
			"total_units_shipped": unitsShipped,
			"cost_per_unit":       costPerUnit,
			"trend":               tcoTrend,
		},
		"risk": map[string]interface{}{
			"critical_resources_percent": criticalPercent,
			"critical_resources":         criticalNames,
			"at_risk_count":              atRisk,
			"total_resources":            totalResources,
		},
		"depletion_forecast": map[string]interface{}{
			"within_7_days":   within7,
			"within_14_days":  within14,
			"within_30_days":  within30,
			"action_required": len(within7) > 0,
		},
		"inventory_total_value": inventoryTotalValue,
	}, nil
}

// 🚀 PRO FEATURE #2: GetDemandForecast прогнозує попит на наступні 3 місяці
func (s *AnalyticsService) GetDemandForecast(ctx context.Context, unitID int64) (map[string]interface{}, error) {
	args := []interface{}{}
	filter := "WHERE sr.created_at >= NOW() - INTERVAL '90 days' AND sr.status NOT IN ('REJECTED', 'CANCELED', 'CANCELLED')"
	if tenantID := repositories.TenantFromCtx(ctx); tenantID != "" {
		args = append(args, tenantID)
		filter += fmt.Sprintf(" AND sr.tenant_id = $%d::uuid", len(args))
	}
	if unitID > 0 {
		args = append(args, unitID)
		filter += fmt.Sprintf(`
			AND (
				EXISTS (
					SELECT 1
					FROM warehouses w
					WHERE w.id = sr.target_warehouse_id
					  AND w.unit_id = $%d
					  AND w.tenant_id = sr.tenant_id
				)
				OR EXISTS (
					SELECT 1
					FROM resources res
					LEFT JOIN warehouses rw
					  ON rw.id = res.warehouse_id
					 AND rw.tenant_id = res.tenant_id
					WHERE res.id = sr.resource_id
					  AND res.tenant_id = sr.tenant_id
					  AND (res.unit_id = $%d OR rw.unit_id = $%d)
				)
			)`, len(args), len(args), len(args))
	}

	query := `
		WITH monthly_usage AS (
			SELECT 
				DATE_TRUNC('month', sr.created_at)::date as month,
				COALESCE(sr.resource_id::text, NULLIF(sr.resource_name, ''), sr.id::text) as resource_key,
				SUM(sr.quantity) as requested_quantity
			FROM supply_requests sr
			` + filter + `
			GROUP BY DATE_TRUNC('month', sr.created_at), COALESCE(sr.resource_id::text, NULLIF(sr.resource_name, ''), sr.id::text)
		),
		resource_trends AS (
			SELECT 
				resource_key,
				AVG(requested_quantity) as avg_monthly_demand,
				MAX(requested_quantity) as peak_demand
			FROM monthly_usage
			GROUP BY resource_key
		)
		SELECT 
			COUNT(*) as resources_analyzed,
			ROUND(COALESCE(AVG(avg_monthly_demand), 0)::numeric, 2) as avg_demand,
			COALESCE(MAX(peak_demand), 0) as peak_demand_observed
		FROM resource_trends
	`

	var resourcesAnalyzed int64
	var avgDemand, peakDemand sql.NullFloat64
	_ = s.db.QueryRow(ctx, query, args...).Scan(&resourcesAnalyzed, &avgDemand, &peakDemand)

	if resourcesAnalyzed == 0 {
		return map[string]interface{}{
			"resources_analyzed": 0,
			"avg_monthly_demand": 0,
			"peak_demand":        0,
			"forecast_3_months": map[string]interface{}{
				"month_1": 0,
				"month_2": 0,
				"month_3": 0,
			},
		}, nil
	}

	avgValue := 0.0
	if avgDemand.Valid {
		avgValue = avgDemand.Float64
	}
	peakValue := 0.0
	if peakDemand.Valid {
		peakValue = peakDemand.Float64
	}

	return map[string]interface{}{
		"resources_analyzed": resourcesAnalyzed,
		"avg_monthly_demand": avgValue,
		"peak_demand":        peakValue,
		"forecast_3_months": map[string]interface{}{
			"month_1": avgValue,
			"month_2": avgValue,
			"month_3": avgValue,
		},
	}, nil
}

// GetPredictiveMaintenanceSchedule returns predicted maintenance schedule for vehicles
// Реалізація використовує реальну схему: vehicles + maintenance_records + fuel_records (для пробігу).
func (s *AnalyticsService) GetPredictiveMaintenanceSchedule(ctx context.Context, unitID int64) (map[string]interface{}, error) {
	tenantID := repositories.TenantFromCtx(ctx)
	if tenantID == "" {
		return nil, fmt.Errorf("tenant_id is required for predictive maintenance")
	}

	args := []interface{}{tenantID}
	where := "WHERE v.status != 'INACTIVE' AND v.tenant_id = $1::uuid"

	if unitID > 0 {
		args = append(args, unitID)
		where += fmt.Sprintf(`
			AND EXISTS (
				SELECT 1
				FROM warehouses w
				WHERE (w.id = v.current_warehouse_id OR w.id = v.home_warehouse_id)
				  AND w.unit_id = $%d
				  AND w.tenant_id = v.tenant_id
			)`, len(args))
	}

	query := `
		SELECT
			v.id,
			v.plate_number,
			COALESCE(NULLIF(v.maintenance_interval_km, 0), 10000) AS interval_km,
			COALESCE(v.last_maintenance_odometer, 0) AS last_odo,
			GREATEST(COALESCE(v.last_maintenance_odometer, 0), COALESCE((
				SELECT MAX(fr.odometer_km)
				FROM fuel_records fr
				WHERE fr.vehicle_id = v.id
				  AND fr.tenant_id = v.tenant_id
			), 0)) AS current_odo,
			COALESCE((
				SELECT MAX(mr.created_at)
				FROM maintenance_records mr
				WHERE mr.vehicle_id = v.id
				  AND COALESCE(mr.status, 'COMPLETED') = 'COMPLETED'
			), v.created_at) AS last_service_at,
			sr.id::text AS scheduled_record_id,
			sr.service_type AS scheduled_service_type,
			sr.scheduled_for
		FROM vehicles v
		LEFT JOIN LATERAL (
			SELECT mr.id, mr.service_type, mr.scheduled_for
			FROM maintenance_records mr
			WHERE mr.vehicle_id = v.id
			  AND COALESCE(mr.status, 'COMPLETED') = 'SCHEDULED'
			ORDER BY COALESCE(mr.scheduled_for, mr.created_at) DESC
			LIMIT 1
		) sr ON TRUE
		` + where + `
		ORDER BY v.plate_number
	`

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("maintenance query: %w", err)
	}
	defer rows.Close()

	schedules := make([]map[string]interface{}, 0)
	overdue := 0
	dueSoon := 0
	compliantVehicles := 0
	totalVehicles := 0
	idCounter := 0

	for rows.Next() {
		var (
			vehicleID, plate     string
			intervalKM           int
			lastOdo              int
			currentOdo           int
			lastServiceAt        time.Time
			scheduledRecordID    sql.NullString
			scheduledServiceType sql.NullString
			scheduledFor         sql.NullTime
		)
		if err := rows.Scan(&vehicleID, &plate, &intervalKM, &lastOdo, &currentOdo, &lastServiceAt, &scheduledRecordID, &scheduledServiceType, &scheduledFor); err != nil {
			continue
		}
		totalVehicles++

		mileageSince := currentOdo - lastOdo
		if mileageSince < 0 {
			mileageSince = 0
		}

		nextDate := lastServiceAt.AddDate(0, 12, 0)
		if scheduledFor.Valid {
			nextDate = scheduledFor.Time
		}
		daysRemaining := int(time.Until(nextDate).Hours() / 24)
		recommendedKM := intervalKM
		if recommendedKM < 1000 {
			recommendedKM = 10000
		}
		mileageRatio := float64(mileageSince) / float64(recommendedKM)
		status := "DUE"
		if scheduledRecordID.Valid {
			status = "SCHEDULED"
		} else if daysRemaining < 0 || mileageSince >= recommendedKM {
			status = "OVERDUE"
		}
		priority := "LOW"
		if status == "OVERDUE" || mileageRatio >= 1 || daysRemaining <= 7 {
			priority = "HIGH"
		} else if mileageRatio >= 0.8 || daysRemaining <= 30 {
			priority = "MEDIUM"
		}
		serviceType := "INSPECTION"
		if scheduledServiceType.Valid && scheduledServiceType.String != "" {
			serviceType = scheduledServiceType.String
		}

		idCounter++
		if status == "OVERDUE" {
			overdue++
		} else if daysRemaining <= 30 {
			dueSoon++
			compliantVehicles++
		} else {
			compliantVehicles++
		}

		schedules = append(schedules, map[string]interface{}{
			"id":                    idCounter,
			"vehicle_id":            vehicleID,
			"vehicle_plate":         plate,
			"scheduled_record_id":   scheduledRecordID.String,
			"service_type":          serviceType,
			"last_service_date":     lastServiceAt.Format(time.RFC3339),
			"next_service_date":     nextDate.Format(time.RFC3339),
			"days_remaining":        daysRemaining,
			"mileage_since_service": mileageSince,
			"recommended_mileage":   recommendedKM,
			"priority":              priority,
			"status":                status,
			"current_odometer":      currentOdo,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("maintenance rows: %w", err)
	}

	avgCompliance := 100.0
	if totalVehicles > 0 {
		avgCompliance = float64(compliantVehicles) * 100.0 / float64(totalVehicles)
	}

	return map[string]interface{}{
		"schedules":          schedules,
		"total_overdue":      overdue,
		"total_due_soon":     dueSoon,
		"average_compliance": avgCompliance,
	}, nil
}

// GetFuelAnomalyDetection detects fuel consumption anomalies and potential fraud
// Використовує fuel_records.is_anomaly + статистику по записах.
func (s *AnalyticsService) GetFuelAnomalyDetection(ctx context.Context, unitID int64) (map[string]interface{}, error) {
	args := []interface{}{}
	where := "WHERE 1=1"

	if tenantID := repositories.TenantFromCtx(ctx); tenantID != "" {
		args = append(args, tenantID)
		where += fmt.Sprintf(" AND v.tenant_id = $%d::uuid", len(args))
	}

	if unitID > 0 {
		args = append(args, unitID)
		where += fmt.Sprintf(`
			AND EXISTS (
				SELECT 1
				FROM warehouses w
				WHERE (w.id = v.current_warehouse_id OR w.id = v.home_warehouse_id)
				  AND w.unit_id = $%d
				  AND w.tenant_id = v.tenant_id
			)`, len(args))
	}

	// Базова статистика по машинах за 90 днів
	statsQuery := `
		SELECT
			v.id,
			v.plate_number,
			COUNT(fr.id) FILTER (WHERE fr.created_at >= NOW() - INTERVAL '90 days') AS records_90d,
			COUNT(fr.id) FILTER (WHERE fr.is_anomaly AND fr.created_at >= NOW() - INTERVAL '90 days') AS anomaly_count,
			COALESCE(AVG(fr.liters) FILTER (WHERE fr.is_anomaly AND fr.created_at >= NOW() - INTERVAL '90 days'), 0) AS avg_anomaly_liters,
			COALESCE(SUM(fr.anomaly_excess_liters) FILTER (WHERE fr.is_anomaly AND fr.created_at >= NOW() - INTERVAL '90 days'), 0) AS total_excess_liters,
			MAX(fr.created_at) FILTER (WHERE fr.is_anomaly) AS last_anomaly_at,
			STRING_AGG(DISTINCT fr.anomaly_reason, '; ') FILTER (WHERE fr.is_anomaly AND fr.anomaly_reason IS NOT NULL AND fr.created_at >= NOW() - INTERVAL '90 days') AS reasons
		FROM vehicles v
		LEFT JOIN fuel_records fr ON fr.vehicle_id = v.id AND fr.tenant_id = v.tenant_id
		` + where + `
		GROUP BY v.id, v.plate_number
		HAVING COUNT(fr.id) > 0
		ORDER BY anomaly_count DESC, v.plate_number
	`

	rows, err := s.db.Query(ctx, statsQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("fuel anomalies query: %w", err)
	}
	defer rows.Close()

	priceArgs := []interface{}{}
	priceFilter := ""
	if tenantID := repositories.TenantFromCtx(ctx); tenantID != "" {
		priceArgs = append(priceArgs, tenantID)
		priceFilter += fmt.Sprintf(" AND sr.tenant_id = $%d::uuid", len(priceArgs))
	}
	if unitID > 0 {
		priceArgs = append(priceArgs, unitID)
		priceFilter += fmt.Sprintf(`
			AND EXISTS (
				SELECT 1
				FROM vehicles v
				JOIN warehouses w
				  ON (w.id = v.current_warehouse_id OR w.id = v.home_warehouse_id)
				 AND w.tenant_id = v.tenant_id
				WHERE v.id = sr.vehicle_id
				  AND v.tenant_id = sr.tenant_id
				  AND w.unit_id = $%d
			)`, len(priceArgs))
	}
	priceQuery := fmt.Sprintf(`
		SELECT COALESCE(SUM(sr.cost_uah) / NULLIF(SUM(sr.liters), 0), 0)
		FROM shipment_refuels sr
		WHERE sr.cost_uah IS NOT NULL%s
	`, priceFilter)
	avgFuelUAHPerLiter := 0.0
	_ = s.db.QueryRow(ctx, priceQuery, priceArgs...).Scan(&avgFuelUAHPerLiter)

	anomalies := make([]map[string]interface{}, 0)
	totalMonitored := 0
	withAnomalies := 0
	totalLoss := 0.0
	highRisk := 0
	id := 0

	for rows.Next() {
		var (
			vehicleID, plate  string
			records90d        int64
			anomalyCount      int64
			avgAnomalyLiters  float64
			totalExcessLiters float64
			lastAt            *time.Time
			reasons           *string
		)
		if err := rows.Scan(&vehicleID, &plate, &records90d, &anomalyCount, &avgAnomalyLiters, &totalExcessLiters, &lastAt, &reasons); err != nil {
			continue
		}
		totalMonitored++
		if anomalyCount == 0 {
			continue
		}
		withAnomalies++

		// Ризик = доля аномалій * 100, обмежена зверху
		ratio := 0.0
		if records90d > 0 {
			ratio = float64(anomalyCount) / float64(records90d)
		}
		riskScore := minFloat(ratio*100.0, 100.0)

		level := "LOW"
		switch {
		case riskScore >= 50:
			level = "HIGH"
			highRisk++
		case riskScore >= 20:
			level = "MEDIUM"
		}

		// Тип аномалії визначаємо з причини
		anomalyType := "ABNORMAL_CONSUMPTION"
		details := "Виявлено аномалії у заправках та списаннях"
		if reasons != nil && *reasons != "" {
			details = *reasons
			lower := strings.ToLower(*reasons)
			switch {
			case strings.Contains(lower, "крадіж"), strings.Contains(lower, "злив"):
				anomalyType = "FREQUENT_SMALL_REFILLS"
			case strings.Contains(lower, "ціна"), strings.Contains(lower, "вартість"):
				anomalyType = "PRICE_ANOMALY"
			case strings.Contains(lower, "екстрем"), strings.Contains(lower, "переповн"):
				anomalyType = "EXTREME_REFILL"
			case strings.Contains(lower, "перевитрата"), strings.Contains(lower, "витрата"):
				anomalyType = "ABNORMAL_CONSUMPTION"
			}
		}

		// Потенційні втрати на місяць (екстраполюємо з 90 днів).
		// Рахуємо лише ЗАЙВЕ пальне (перевитрата понад норму + витрата без руху),
		// а не весь обсяг аномальних записів. total_excess_liters — сума за 90 днів,
		// тож ділимо на 3, щоб отримати оцінку на місяць.
		potentialLossPerMonth := (totalExcessLiters * avgFuelUAHPerLiter) / 3.0
		// Підстраховка для історичних записів без excess (стара схема): якщо аномалії є,
		// а зайвих літрів не зафіксовано — беремо консервативну оцінку від середнього обсягу.
		if potentialLossPerMonth == 0 && anomalyCount > 0 {
			potentialLossPerMonth = (avgAnomalyLiters * float64(anomalyCount) * avgFuelUAHPerLiter) / 3.0
		}
		totalLoss += potentialLossPerMonth

		lastDetected := time.Now().Format(time.RFC3339)
		if lastAt != nil {
			lastDetected = lastAt.Format(time.RFC3339)
		}

		// Впевненість: базова 50% + додатково залежно від кількості даних
		// Чим більше записів, тим вища впевненість (макс 95%)
		confidence := minFloat(50.0+float64(records90d)/2.0, 95.0)

		id++
		anomalies = append(anomalies, map[string]interface{}{
			"id":                  id,
			"vehicle_id":          vehicleID,
			"vehicle_plate":       plate,
			"anomaly_type":        anomalyType,
			"risk_score":          riskScore,
			"investigation_level": level,
			"last_detected":       lastDetected,
			"details":             details,
			"confidence":          confidence,
			"potential_loss":      potentialLossPerMonth,
		})
	}

	return map[string]interface{}{
		"anomalies":                anomalies,
		"total_vehicles_monitored": totalMonitored,
		"vehicles_with_anomalies":  withAnomalies,
		"total_potential_loss":     totalLoss,
		"high_risk_count":          highRisk,
	}, nil
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

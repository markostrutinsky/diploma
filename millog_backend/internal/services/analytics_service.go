package services

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"

	"millog_backend/internal/models"
	"millog_backend/internal/repositories"
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

	headers := []string{"Філія / Підрозділ", "Склад", "Категорія", "Найменування майна", "Кількість", "Од. вим.", "Стан"}
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
	}

	// Додаємо автофільтр на всі дані
	f.AutoFilter(sheetName, fmt.Sprintf("A1:G%d", len(data)+1), nil)

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GenerateFuelExcel формує XLSX зі звітом по пальному
func (s *AnalyticsService) GenerateFuelExcel(ctx context.Context, startDate, endDate time.Time) ([]byte, error) {
	data, err := s.repo.GetFuelForExport(ctx, s.db, startDate, endDate)
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
	// unit-фільтр опціональний: 0 означає "всі підрозділи"
	unitFilter := ""
	args := []interface{}{startDate, endDate}
	if unitID > 0 {
		unitFilter = " AND w.unit_id = $3"
		args = append(args, unitID)
	}

	// ---------- SLA ----------
	// Заявка "вчасна", якщо її апрувнули за ≤ 24 години.
	slaQuery := fmt.Sprintf(`
		WITH req AS (
			SELECT sr.id, sr.status, sr.created_at, sr.approved_at
			FROM supply_requests sr
			LEFT JOIN warehouses w ON w.id = sr.target_warehouse_id
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
	`, unitFilter)

	var totalReq, onTime int64
	var avgHours float64
	if err := s.db.QueryRow(ctx, slaQuery, args...).Scan(&totalReq, &onTime, &avgHours); err != nil {
		totalReq, onTime, avgHours = 0, 0, 0
	}
	onTimePercent := 100.0
	if totalReq > 0 {
		onTimePercent = float64(onTime) * 100.0 / float64(totalReq)
	}

	// ---------- TCO ----------
	// cost per liter (UAH) — константа, бо в схемі немає колонки з ціною пального.
	const fuelUAHPerLiter = 55.0
	fuelArgs := []interface{}{startDate, endDate}
	vehicleUnitFilter := ""
	if unitID > 0 {
		vehicleUnitFilter = `
			AND fr.vehicle_id IN (
				SELECT DISTINCT s.vehicle_id FROM shipments s
				JOIN warehouses w ON w.id = s.from_warehouse_id OR w.id = s.to_warehouse_id
				WHERE w.unit_id = $3
			)`
		fuelArgs = append(fuelArgs, unitID)
	}
	fuelQuery := fmt.Sprintf(`
		SELECT COALESCE(SUM(liters), 0)
		FROM fuel_records fr
		WHERE fr.record_type = 'CONSUMPTION'
		  AND fr.created_at::date >= $1::date
		  AND fr.created_at::date <= $2::date%s
	`, vehicleUnitFilter)
	var totalLiters float64
	_ = s.db.QueryRow(ctx, fuelQuery, fuelArgs...).Scan(&totalLiters)

	// Кількість відвантажених одиниць
	shipArgs := []interface{}{startDate, endDate}
	shipUnitFilter := ""
	if unitID > 0 {
		shipUnitFilter = " AND (wf.unit_id = $3 OR wt.unit_id = $3)"
		shipArgs = append(shipArgs, unitID)
	}
	shipQuery := fmt.Sprintf(`
		SELECT COALESCE(SUM(si.quantity), 0)
		FROM shipments s
		JOIN shipment_items si ON si.shipment_id = s.id
		LEFT JOIN warehouses wf ON wf.id = s.from_warehouse_id
		LEFT JOIN warehouses wt ON wt.id = s.to_warehouse_id
		WHERE s.created_at::date >= $1::date
		  AND s.created_at::date <= $2::date%s
	`, shipUnitFilter)
	var unitsShipped int64
	_ = s.db.QueryRow(ctx, shipQuery, shipArgs...).Scan(&unitsShipped)

	totalFuelCost := totalLiters * fuelUAHPerLiter
	costPerUnit := 0.0
	if unitsShipped > 0 {
		costPerUnit = totalFuelCost / float64(unitsShipped)
	}

	// ---------- RISK ----------
	riskArgs := []interface{}{}
	resourceUnitFilter := ""
	if unitID > 0 {
		resourceUnitFilter = ` AND (r.unit_id = $1 OR r.warehouse_id IN (SELECT id FROM warehouses WHERE unit_id = $1))`
		riskArgs = append(riskArgs, unitID)
	}
	riskQuery := fmt.Sprintf(`
		SELECT r.id, r.name, r.quantity, r.min_quantity
		FROM resources r
		WHERE r.min_quantity > 0%s
	`, resourceUnitFilter)

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
			"trend":               "stable",
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
	}, nil
}

// 🚀 PRO FEATURE #2: GetDemandForecast прогнозує попит на наступні 3 місяці
func (s *AnalyticsService) GetDemandForecast(ctx context.Context, unitID int64) (map[string]interface{}, error) {
	// Простий алгоритм: беремо середню 90 днів попиту та проектуємо вперед
	query := `
		WITH monthly_usage AS (
			SELECT 
				DATE_TRUNC('month', sr.created_at)::date as month,
				sr.resource_id,
				COUNT(*) as request_count,
				AVG(r.quantity) as avg_quantity
			FROM supply_requests sr
			JOIN resources r ON sr.resource_id = r.id
			WHERE sr.unit_id = $1
			AND sr.created_at >= NOW() - INTERVAL '90 days'
			GROUP BY DATE_TRUNC('month', sr.created_at), sr.resource_id
		),
		resource_trends AS (
			SELECT 
				resource_id,
				ROUND(AVG(request_count), 2) as avg_monthly_demand,
				MAX(request_count) as peak_demand
			FROM monthly_usage
			GROUP BY resource_id
		)
		SELECT 
			COUNT(*) as resources_analyzed,
			ROUND(AVG(avg_monthly_demand), 2) as avg_demand,
			MAX(peak_demand) as peak_demand_observed
		FROM resource_trends
	`

	var resourcesAnalyzed, avgDemand, peakDemand interface{}
	_ = s.db.QueryRow(ctx, query, unitID).Scan(&resourcesAnalyzed, &avgDemand, &peakDemand)

	if resourcesAnalyzed == nil {
		return map[string]interface{}{
			"resources_analyzed": 0,
			"forecast_3_months": map[string]interface{}{
				"month_1": 0,
				"month_2": 0,
				"month_3": 0,
			},
		}, nil
	}

	return map[string]interface{}{
		"resources_analyzed": resourcesAnalyzed,
		"avg_monthly_demand": avgDemand,
		"peak_demand":        peakDemand,
		"forecast_3_months": map[string]interface{}{
			"month_1": avgDemand,
			"month_2": avgDemand,
			"month_3": avgDemand,
		},
	}, nil
}

// GetPredictiveMaintenanceSchedule returns predicted maintenance schedule for vehicles
// Реалізація використовує реальну схему: vehicles + maintenance_records + fuel_records (для пробігу).
func (s *AnalyticsService) GetPredictiveMaintenanceSchedule(ctx context.Context, unitID int64) (map[string]interface{}, error) {
	// У схемі немає vehicles.unit_id, тому фільтр по підрозділу наразі ігнорується.
	_ = unitID

	query := `
		SELECT
			v.id,
			v.plate_number,
			COALESCE(NULLIF(v.maintenance_interval_km, 0), 10000) AS interval_km,
			COALESCE(v.last_maintenance_odometer, 0) AS last_odo,
			COALESCE((SELECT MAX(fr.odometer_km) FROM fuel_records fr WHERE fr.vehicle_id = v.id), 0) AS current_odo,
			(SELECT MAX(mr.created_at) FROM maintenance_records mr WHERE mr.vehicle_id = v.id) AS last_service_at
		FROM vehicles v
		WHERE v.status <> 'WRITTEN_OFF'
		ORDER BY v.plate_number
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("maintenance query: %w", err)
	}
	defer rows.Close()

	serviceTypes := []struct {
		code     string
		kmFactor float64 // частка інтервалу ТО
		monthly  int     // періодичність у місяцях
		priority string
	}{
		{"OIL_CHANGE", 0.5, 6, "HIGH"},
		{"TIRE_ROTATION", 1.0, 9, "MEDIUM"},
		{"FILTER_REPLACEMENT", 1.0, 12, "MEDIUM"},
		{"INSPECTION", 2.0, 12, "LOW"},
	}

	schedules := make([]map[string]interface{}, 0)
	overdue := 0
	dueSoon := 0
	compliantVehicles := 0
	totalVehicles := 0
	now := time.Now()
	idCounter := 0

	for rows.Next() {
		var (
			vehicleID, plate string
			intervalKM       int
			lastOdo          int
			currentOdo       int
			lastServiceAt    *time.Time
		)
		if err := rows.Scan(&vehicleID, &plate, &intervalKM, &lastOdo, &currentOdo, &lastServiceAt); err != nil {
			continue
		}
		totalVehicles++

		mileageSince := currentOdo - lastOdo
		if mileageSince < 0 {
			mileageSince = 0
		}

		lastServiceDate := now.AddDate(0, -6, 0)
		if lastServiceAt != nil {
			lastServiceDate = *lastServiceAt
		}

		vehicleCompliant := true

		for _, st := range serviceTypes {
			idCounter++
			recommendedKM := int(float64(intervalKM) * st.kmFactor)
			if recommendedKM < 1000 {
				recommendedKM = 1000
			}
			nextDate := lastServiceDate.AddDate(0, st.monthly, 0)
			daysRemaining := int(time.Until(nextDate).Hours() / 24)

			status := "SCHEDULED"
			if daysRemaining < 0 || mileageSince > recommendedKM {
				status = "OVERDUE"
				overdue++
				vehicleCompliant = false
			} else if daysRemaining <= 30 {
				dueSoon++
			}

			schedules = append(schedules, map[string]interface{}{
				"id":                    idCounter,
				"vehicle_id":            vehicleID,
				"vehicle_plate":         plate,
				"service_type":          st.code,
				"last_service_date":     lastServiceDate.Format(time.RFC3339),
				"next_service_date":     nextDate.Format(time.RFC3339),
				"days_remaining":        daysRemaining,
				"mileage_since_service": mileageSince,
				"recommended_mileage":   recommendedKM,
				"priority":              st.priority,
				"status":                status,
			})
		}

		if vehicleCompliant {
			compliantVehicles++
		}
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
	_ = unitID // немає vehicles.unit_id у схемі

	// Базова статистика по машинах за 90 днів
	statsQuery := `
		SELECT
			v.id,
			v.plate_number,
			COUNT(fr.id) FILTER (WHERE fr.created_at >= NOW() - INTERVAL '90 days') AS records_90d,
			COUNT(fr.id) FILTER (WHERE fr.is_anomaly AND fr.created_at >= NOW() - INTERVAL '90 days') AS anomaly_count,
			COALESCE(AVG(fr.liters) FILTER (WHERE fr.is_anomaly AND fr.created_at >= NOW() - INTERVAL '90 days'), 0) AS avg_anomaly_liters,
			MAX(fr.created_at) FILTER (WHERE fr.is_anomaly) AS last_anomaly_at,
			STRING_AGG(DISTINCT fr.anomaly_reason, '; ') FILTER (WHERE fr.is_anomaly AND fr.anomaly_reason IS NOT NULL AND fr.created_at >= NOW() - INTERVAL '90 days') AS reasons
		FROM vehicles v
		LEFT JOIN fuel_records fr ON fr.vehicle_id = v.id
		GROUP BY v.id, v.plate_number
		ORDER BY anomaly_count DESC, v.plate_number
	`

	rows, err := s.db.Query(ctx, statsQuery)
	if err != nil {
		return nil, fmt.Errorf("fuel anomalies query: %w", err)
	}
	defer rows.Close()

	const fuelUAHPerLiter = 55.0
	anomalies := make([]map[string]interface{}, 0)
	totalMonitored := 0
	withAnomalies := 0
	totalLoss := 0.0
	highRisk := 0
	id := 0

	for rows.Next() {
		var (
			vehicleID, plate string
			records90d       int64
			anomalyCount     int64
			avgAnomalyLiters float64
			lastAt           *time.Time
			reasons          *string
		)
		if err := rows.Scan(&vehicleID, &plate, &records90d, &anomalyCount, &avgAnomalyLiters, &lastAt, &reasons); err != nil {
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
		riskScore := ratio * 100.0
		if riskScore > 100 {
			riskScore = 100
		}

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
			lower := *reasons
			switch {
			case contains(lower, "крадіж"), contains(lower, "злив"):
				anomalyType = "FREQUENT_SMALL_REFILLS"
			case contains(lower, "ціна"), contains(lower, "вартість"):
				anomalyType = "PRICE_ANOMALY"
			case contains(lower, "екстрем"), contains(lower, "переповн"):
				anomalyType = "EXTREME_REFILL"
			}
		}

		potentialLoss := avgAnomalyLiters * float64(anomalyCount) * fuelUAHPerLiter
		totalLoss += potentialLoss

		lastDetected := time.Now().Format(time.RFC3339)
		if lastAt != nil {
			lastDetected = lastAt.Format(time.RFC3339)
		}

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
			"confidence":          minFloat(50+ratio*100, 99),
			"potential_loss":      potentialLoss,
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

func contains(s, sub string) bool {
	return len(sub) > 0 && len(s) >= len(sub) && (indexFold(s, sub) >= 0)
}

// indexFold — просте регістронезалежне підрядок-пошук (без зовнішніх залежностей).
func indexFold(s, sub string) int {
	ls := []rune(s)
	lsub := []rune(sub)
	for i := 0; i+len(lsub) <= len(ls); i++ {
		match := true
		for j := 0; j < len(lsub); j++ {
			a := ls[i+j]
			b := lsub[j]
			if a >= 'A' && a <= 'Z' {
				a += 'a' - 'A'
			}
			if b >= 'A' && b <= 'Z' {
				b += 'a' - 'A'
			}
			if a != b {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

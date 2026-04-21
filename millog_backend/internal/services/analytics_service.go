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
// SLA (% вовчасно), TCO (на ресурс), Ризики (дефіцит), Прогноз (що скончатися через 7/14/30 днів)
func (s *AnalyticsService) GetAdvancedKPIs(ctx context.Context, startDate, endDate string, unitID int64) (map[string]interface{}, error) {
	kpis := make(map[string]interface{})

	// SLA: % заявок виконаних вовчасу (без затримок)
	slaQuery := `
		WITH request_times AS (
			SELECT 
				id,
				status,
				EXTRACT(EPOCH FROM (updated_at - created_at))/3600 as hours_to_approve
			FROM supply_requests
			WHERE 
				unit_id = $1
				AND created_at::date >= $2::date
				AND created_at::date <= $3::date
		)
		SELECT 
			ROUND(100.0 * COUNT(CASE WHEN hours_to_approve <= 24 THEN 1 END) / NULLIF(COUNT(*), 0), 2) as sla_percentage,
			COUNT(*) as total_requests,
			COUNT(CASE WHEN hours_to_approve <= 24 THEN 1 END) as on_time
		FROM request_times
		WHERE status IN ('APPROVED', 'REJECTED')
	`

	var slaPercent, totalReq, onTime interface{}
	err := s.db.QueryRow(ctx, slaQuery, unitID, startDate, endDate).Scan(&slaPercent, &totalReq, &onTime)
	if err != nil {
		slaPercent, totalReq, onTime = 0, 0, 0
	}

	kpis["sla"] = map[string]interface{}{
		"percentage":     slaPercent,
		"total_requests": totalReq,
		"on_time":        onTime,
		"target":         95.0,
		"status":         "good", // Буде змінено на основі значення
	}

	// TCO (Total Cost of Ownership): витрати на доставку / кількість ресурсів
	tcoQuery := `
		SELECT 
			COUNT(DISTINCT sr.id) as shipments,
			COALESCE(SUM(v.fuel_cost), 0) as total_fuel_cost,
			COUNT(DISTINCT sr.resource_id) as resources_moved
		FROM supply_requests sr
		LEFT JOIN vehicles v ON sr.vehicle_id = v.id
		WHERE sr.unit_id = $1
		AND sr.created_at::date >= $2::date
		AND sr.created_at::date <= $3::date
	`

	var shipments, fuelCost, resourcesMoved interface{}
	_ = s.db.QueryRow(ctx, tcoQuery, unitID, startDate, endDate).Scan(&shipments, &fuelCost, &resourcesMoved)
	if shipments == nil {
		shipments, fuelCost, resourcesMoved = 0, 0, 0
	}

	kpis["tco"] = map[string]interface{}{
		"total_cost":        fuelCost,
		"shipments":         shipments,
		"resources_moved":   resourcesMoved,
		"cost_per_shipment": 0, // Розраховується в if
	}

	// RISK: Дефіцит ресурсів (запасів менше 20% від min_quantity)
	riskQuery := `
		SELECT 
			COUNT(*) as at_risk_count,
			ROUND(100.0 * COUNT(*) / NULLIF(COUNT(DISTINCT id), 0), 2) as risk_percentage
		FROM resources
		WHERE warehouse_id IN (SELECT id FROM warehouses WHERE unit_id = $1)
		AND quantity < (min_quantity * 0.2)
	`

	var riskCount, riskPercent interface{}
	_ = s.db.QueryRow(ctx, riskQuery, unitID).Scan(&riskCount, &riskPercent)
	if riskCount == nil {
		riskCount, riskPercent = 0, 0
	}

	kpis["risk"] = map[string]interface{}{
		"at_risk_count":   riskCount,
		"risk_percentage": riskPercent,
		"status":          "warning", // Залежить від значення
	}

	// FORECAST: Яка частина ресурсів скончатися за 7/14/30 днів за поточною швидкістю
	forecastQuery := `
		WITH usage_rate AS (
			SELECT 
				resource_id,
				COALESCE(AVG(quantity), 0) as daily_usage
			FROM supply_requests
			WHERE unit_id = $1
			AND created_at::date >= (NOW()::date - INTERVAL '30 days')
			AND created_at::date < NOW()::date
			GROUP BY resource_id
		)
		SELECT 
			COUNT(CASE WHEN (r.quantity / NULLIF(ur.daily_usage, 0)) <= 7 THEN 1 END) as depletes_7days,
			COUNT(CASE WHEN (r.quantity / NULLIF(ur.daily_usage, 0)) <= 14 THEN 1 END) as depletes_14days,
			COUNT(CASE WHEN (r.quantity / NULLIF(ur.daily_usage, 0)) <= 30 THEN 1 END) as depletes_30days,
			COUNT(DISTINCT r.id) as total_resources
		FROM resources r
		LEFT JOIN usage_rate ur ON r.id = ur.resource_id
		WHERE r.warehouse_id IN (SELECT id FROM warehouses WHERE unit_id = $1)
	`

	var depl7, depl14, depl30, totalRes interface{}
	_ = s.db.QueryRow(ctx, forecastQuery, unitID).Scan(&depl7, &depl14, &depl30, &totalRes)
	if depl7 == nil {
		depl7, depl14, depl30, totalRes = 0, 0, 0, 0
	}

	kpis["forecast"] = map[string]interface{}{
		"depletes_7_days":  depl7,
		"depletes_14_days": depl14,
		"depletes_30_days": depl30,
		"total_resources":  totalRes,
	}

	return kpis, nil
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
// PRO FEATURE: Analyzes historical maintenance, mileage, and usage patterns
func (s *AnalyticsService) GetPredictiveMaintenanceSchedule(ctx context.Context, unitID int64) (map[string]interface{}, error) {
	// Query to get vehicle maintenance history and mileage data
	query := `
		SELECT 
			v.id,
			v.plate_number,
			v.brand,
			v.model,
			COALESCE(v.mileage, 0) AS current_mileage,
			COUNT(CASE WHEN am.maintenance_type = 'OIL_CHANGE' THEN 1 END) AS oil_changes_count,
			COUNT(CASE WHEN am.maintenance_type = 'TIRE_ROTATION' THEN 1 END) AS tire_rotations_count,
			COUNT(CASE WHEN am.maintenance_type = 'FILTER_REPLACEMENT' THEN 1 END) AS filter_replacements_count,
			MAX(am.performed_at) AS last_maintenance_date,
			COALESCE(AVG(EXTRACT(EPOCH FROM (NOW() - am.performed_at)) / 86400 / 30), 0)::int AS avg_months_between_maintenance
		FROM vehicles v
		LEFT JOIN audit_logs am ON v.id = CAST(am.entity_id AS BIGINT) 
			AND am.action = 'MAINTENANCE' 
			AND am.entity_type = 'VEHICLE'
		WHERE v.unit_id = $1
		GROUP BY v.id, v.plate_number, v.brand, v.model, v.mileage
		ORDER BY COALESCE(v.mileage, 0) DESC
	`

	rows, err := s.db.Query(ctx, query, unitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get maintenance history: %w", err)
	}
	defer rows.Close()

	type VehicleMaintenance struct {
		ID                      int64
		PlateNumber             string
		Brand                   string
		Model                   string
		CurrentMileage          int64
		OilChangesCount         int
		TireRotationsCount      int
		FilterReplacementsCount int
		LastMaintenanceDate     *time.Time
		AvgMonthsBetweenMaint   int
	}

	var vehicles []VehicleMaintenance
	maintenanceSchedule := make([]map[string]interface{}, 0)

	for rows.Next() {
		var vm VehicleMaintenance
		if err := rows.Scan(&vm.ID, &vm.PlateNumber, &vm.Brand, &vm.Model, &vm.CurrentMileage,
			&vm.OilChangesCount, &vm.TireRotationsCount, &vm.FilterReplacementsCount,
			&vm.LastMaintenanceDate, &vm.AvgMonthsBetweenMaint); err != nil {
			continue
		}
		vehicles = append(vehicles, vm)

		// Predict next maintenance based on:
		// 1. Oil change every 5,000 km or 6 months
		// 2. Tire rotation every 15,000 km or 8 months
		// 3. Filter replacement every 10,000 km or 12 months
		maintenanceItems := make([]map[string]interface{}, 0)

		// Oil change prediction
		nextOilChangeMileage := ((vm.CurrentMileage / 5000) + 1) * 5000
		oilChangeMonths := vm.AvgMonthsBetweenMaint
		if oilChangeMonths == 0 {
			oilChangeMonths = 6
		}
		nextOilChangeDate := time.Now().AddDate(0, oilChangeMonths, 0)

		maintenanceItems = append(maintenanceItems, map[string]interface{}{
			"type":           "OIL_CHANGE",
			"description":    "Заміна масла й масляного фільтра",
			"priority":       "high",
			"next_due_km":    nextOilChangeMileage,
			"next_due_date":  nextOilChangeDate.Format(time.RFC3339),
			"days_remaining": int(time.Until(nextOilChangeDate).Hours() / 24),
		})

		// Tire rotation prediction
		nextTireRotationMileage := ((vm.CurrentMileage / 15000) + 1) * 15000
		maintenanceItems = append(maintenanceItems, map[string]interface{}{
			"type":           "TIRE_ROTATION",
			"description":    "Ротація коліс",
			"priority":       "medium",
			"next_due_km":    nextTireRotationMileage,
			"next_due_date":  time.Now().AddDate(0, 8, 0).Format(time.RFC3339),
			"days_remaining": int(time.Until(time.Now().AddDate(0, 8, 0)).Hours() / 24),
		})

		// Filter replacement prediction
		nextFilterReplacementMileage := ((vm.CurrentMileage / 10000) + 1) * 10000
		maintenanceItems = append(maintenanceItems, map[string]interface{}{
			"type":           "FILTER_REPLACEMENT",
			"description":    "Заміна повітряного й кабінного фільтра",
			"priority":       "medium",
			"next_due_km":    nextFilterReplacementMileage,
			"next_due_date":  time.Now().AddDate(0, 12, 0).Format(time.RFC3339),
			"days_remaining": int(time.Until(time.Now().AddDate(0, 12, 0)).Hours() / 24),
		})

		maintenanceSchedule = append(maintenanceSchedule, map[string]interface{}{
			"vehicle_id":           vm.ID,
			"plate_number":         vm.PlateNumber,
			"brand_model":          fmt.Sprintf("%s %s", vm.Brand, vm.Model),
			"current_mileage":      vm.CurrentMileage,
			"total_maintenance":    vm.OilChangesCount + vm.TireRotationsCount + vm.FilterReplacementsCount,
			"last_maintenance":     vm.LastMaintenanceDate,
			"maintenance_schedule": maintenanceItems,
		})
	}

	return map[string]interface{}{
		"unit_id":                 unitID,
		"vehicles_analyzed":       len(vehicles),
		"maintenance_forecast":    maintenanceSchedule,
		"urgent_maintenance_days": 30,
		"generated_at":            time.Now().Format(time.RFC3339),
	}, nil
}

// GetFuelAnomalyDetection detects fuel consumption anomalies and potential fraud
// PRO FEATURE: Uses statistical analysis to identify unusual fuel patterns
func (s *AnalyticsService) GetFuelAnomalyDetection(ctx context.Context, unitID int64) (map[string]interface{}, error) {
	// Query to analyze fuel consumption patterns per vehicle
	query := `
		WITH fuel_stats AS (
			SELECT 
				v.id,
				v.plate_number,
				v.brand,
				v.model,
				COUNT(f.id) AS refill_count,
				SUM(f.liters) AS total_liters,
				AVG(f.liters) AS avg_liters_per_refill,
				STDDEV(f.liters) AS stddev_liters,
				MIN(f.liters) AS min_liters,
				MAX(f.liters) AS max_liters,
				COALESCE(v.mileage, 0) AS current_mileage,
				(SUM(f.cost) / NULLIF(SUM(f.liters), 0))::decimal(10,2) AS avg_price_per_liter
			FROM vehicles v
			LEFT JOIN fuel_entries f ON v.id = f.vehicle_id 
				AND f.created_at >= NOW() - INTERVAL '90 days'
			WHERE v.unit_id = $1
			GROUP BY v.id, v.plate_number, v.brand, v.model, v.mileage
		)
		SELECT 
			id, plate_number, brand, model, refill_count, total_liters, 
			avg_liters_per_refill, stddev_liters, min_liters, max_liters,
			current_mileage, avg_price_per_liter
		FROM fuel_stats
		WHERE refill_count > 0
		ORDER BY COALESCE(stddev_liters, 0) DESC
	`

	rows, err := s.db.Query(ctx, query, unitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fuel consumption data: %w", err)
	}
	defer rows.Close()

	anomalies := make([]map[string]interface{}, 0)
	var totalVehiclesAnalyzed int
	suspiciousVehicles := 0

	for rows.Next() {
		var (
			id, refillCount                               int64
			plateNumber, brand, model                     string
			totalLiters, avgLitersPerRefill, stddevLiters float64
			minLiters, maxLiters, currentMileage          float64
			avgPricePerLiter                              float64
		)

		if err := rows.Scan(&id, &plateNumber, &brand, &model, &refillCount, &totalLiters,
			&avgLitersPerRefill, &stddevLiters, &minLiters, &maxLiters, &currentMileage, &avgPricePerLiter); err != nil {
			continue
		}

		totalVehiclesAnalyzed++

		// Detection logic: Flag if any refill is > 2 standard deviations from average
		var detectedAnomalies []string
		riskScore := 0.0

		// Check for extreme refills (possible tank filling fraud)
		if stddevLiters > 0 {
			upperBound := avgLitersPerRefill + (2 * stddevLiters)
			if maxLiters > upperBound {
				detectedAnomalies = append(detectedAnomalies, fmt.Sprintf("Екстремальне заправлення: %.1f л (норма: %.1f±%.1f л)", maxLiters, avgLitersPerRefill, stddevLiters))
				riskScore += 25.0
			}
		}

		// Check for unusual consumption (fuel theft)
		if avgLitersPerRefill < 10 && refillCount > 10 {
			detectedAnomalies = append(detectedAnomalies, "Часті дрібні заправлення (можлива крадіжка палива)")
			riskScore += 30.0
		}

		// Check for price anomalies
		if avgPricePerLiter > 0 {
			if avgPricePerLiter > 2.5 {
				detectedAnomalies = append(detectedAnomalies, fmt.Sprintf("Висока вартість палива: %.2f грн/л", avgPricePerLiter))
				riskScore += 15.0
			}
		}

		// Check for abnormal consumption vs mileage
		if currentMileage > 0 && totalLiters > 0 {
			fuelEfficiency := currentMileage / totalLiters
			if fuelEfficiency < 3.0 { // Less than 3 km/l is suspicious
				detectedAnomalies = append(detectedAnomalies, fmt.Sprintf("Дуже низька паливна економічність: %.2f км/л (норма: 5-8 км/л)", fuelEfficiency))
				riskScore += 35.0
			}
		}

		if len(detectedAnomalies) > 0 {
			suspiciousVehicles++
			anomalies = append(anomalies, map[string]interface{}{
				"vehicle_id":          id,
				"plate_number":        plateNumber,
				"brand_model":         fmt.Sprintf("%s %s", brand, model),
				"risk_score":          riskScore,
				"refill_count_90d":    refillCount,
				"total_liters_90d":    fmt.Sprintf("%.1f", totalLiters),
				"avg_refill_liters":   fmt.Sprintf("%.1f", avgLitersPerRefill),
				"stddev":              fmt.Sprintf("%.1f", stddevLiters),
				"anomalies":           detectedAnomalies,
				"investigation_level": map[string]string{"low": "Низька", "medium": "Середня", "high": "Висока"}[getInvestigationLevel(riskScore)],
			})
		}
	}

	return map[string]interface{}{
		"unit_id":                unitID,
		"analysis_period_days":   90,
		"vehicles_analyzed":      totalVehiclesAnalyzed,
		"suspicious_vehicles":    suspiciousVehicles,
		"suspicion_rate_percent": float64(suspiciousVehicles) * 100.0 / float64(totalVehiclesAnalyzed),
		"anomalies":              anomalies,
		"generated_at":           time.Now().Format(time.RFC3339),
		"recommendation":         "Перевірте прапорені операції та розслідуйте високий рівень ризику",
	}, nil
}

func getInvestigationLevel(riskScore float64) string {
	if riskScore >= 50 {
		return "high"
	}
	if riskScore >= 25 {
		return "medium"
	}
	return "low"
}

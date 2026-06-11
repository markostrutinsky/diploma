package repositories

import (
	"Omnilog_backend/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type VehicleRepository struct{}

func NewVehicleRepository() *VehicleRepository {
	return &VehicleRepository{}
}

func (r *VehicleRepository) Create(ctx context.Context, db DBExecutor, v *models.Vehicle) error {
	tid := TenantFromCtx(ctx)
	if tid == "" {
		return fmt.Errorf("tenant_id is required for vehicles")
	}

	if v.HomeWarehouseID != nil && *v.HomeWarehouseID != "" {
		var ok bool
		if err := db.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM warehouses WHERE id = $1 AND tenant_id = $2
			)
		`, *v.HomeWarehouseID, tid).Scan(&ok); err != nil {
			return fmt.Errorf("помилка перевірки складу: %w", err)
		}
		if !ok {
			return fmt.Errorf("склад не знайдено у вашій організації")
		}
	}

	if v.DriverID != nil && *v.DriverID != "" {
		var ok bool
		if err := db.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM users
				WHERE id = $1 AND tenant_id = $2 AND status = $3
			)
		`, *v.DriverID, tid, models.StatusActive).Scan(&ok); err != nil {
			return fmt.Errorf("помилка перевірки водія: %w", err)
		}
		if !ok {
			return fmt.Errorf("водія не знайдено у вашій організації")
		}
	}

	query := `
        INSERT INTO vehicles (
            brand, model, plate_number, type, capacity_kg, status, driver_id,
            tank_capacity, fuel_norm, home_warehouse_id, current_warehouse_id, tenant_id
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $10, $11)
        RETURNING id, maintenance_interval_km, last_maintenance_odometer, created_at, updated_at`

	err := db.QueryRow(ctx, query,
		v.Brand, v.Model, v.PlateNumber, v.Type, v.CapacityKg, v.Status, v.DriverID,
		v.TankCapacity, v.FuelNorm, v.HomeWarehouseID, tid,
	).Scan(&v.ID, &v.MaintenanceIntervalKm, &v.LastMaintenanceOdometer, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return fmt.Errorf("помилка створення авто: %w", err)
	}
	return nil
}

func (r *VehicleRepository) GetAll(ctx context.Context, db DBExecutor) ([]models.Vehicle, error) {
	args := []any{}
	tFilter := tenantFilter(ctx, "v", "WHERE", &args)
	query := `
        WITH VehicleStats AS (
            SELECT 
                vehicle_id,
                MIN(created_at) as first_date,
                MAX(created_at) as last_date,
                MIN(odometer_km) as first_odo,
                MAX(odometer_km) as last_odo
            FROM fuel_records
            WHERE created_at >= CURRENT_DATE - INTERVAL '30 days' AND odometer_km IS NOT NULL
            GROUP BY vehicle_id
        )
        SELECT 
            v.id, v.brand, v.model, v.plate_number, v.type, v.capacity_kg, v.status, v.driver_id, u.full_name as driver_name, 
            v.home_warehouse_id, wh.name as home_warehouse_name,
            v.current_warehouse_id, wc.name as current_warehouse_name,
            v.tank_capacity, v.fuel_norm, 
            v.maintenance_interval_km, v.last_maintenance_odometer, v.created_at, v.updated_at,
            COALESCE((
                SELECT odometer_km FROM fuel_records 
                WHERE vehicle_id = v.id AND odometer_km IS NOT NULL 
                ORDER BY created_at DESC LIMIT 1
            ), 0) AS current_odometer,
            COALESCE(
                (vs.last_odo - vs.first_odo)::float / 
                NULLIF(EXTRACT(EPOCH FROM (vs.last_date - vs.first_date))/86400, 0), 
            0) as avg_km_per_day,
            GREATEST(0, COALESCE((
                SELECT SUM(CASE WHEN record_type = 'REFUEL' THEN liters ELSE -liters END)
                FROM fuel_records WHERE vehicle_id = v.id
            ), 0)) AS current_fuel_liters

        FROM vehicles v
        LEFT JOIN users u ON v.driver_id = u.id
        LEFT JOIN warehouses wh ON wh.id = v.home_warehouse_id
        LEFT JOIN warehouses wc ON wc.id = v.current_warehouse_id
        LEFT JOIN VehicleStats vs ON v.id = vs.vehicle_id` + tFilter + `
        ORDER BY v.created_at DESC
    `

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("помилка отримання списку авто: %w", err)
	}
	defer rows.Close()

	var vehicles []models.Vehicle
	for rows.Next() {
		var v models.Vehicle
		err := rows.Scan(
			&v.ID, &v.Brand, &v.Model, &v.PlateNumber, &v.Type, &v.CapacityKg, &v.Status, &v.DriverID, &v.DriverName,
			&v.HomeWarehouseID, &v.HomeWarehouseName,
			&v.CurrentWarehouseID, &v.CurrentWarehouseName,
			&v.TankCapacity, &v.FuelNorm, &v.MaintenanceIntervalKm, &v.LastMaintenanceOdometer,
			&v.CreatedAt, &v.UpdatedAt, &v.CurrentOdometer, &v.AvgKmPerDay, &v.CurrentFuelLiters,
		)
		if err != nil {
			return nil, fmt.Errorf("помилка сканування рядка авто: %w", err)
		}

		if v.CurrentOdometer < v.LastMaintenanceOdometer {
			v.CurrentOdometer = v.LastMaintenanceOdometer
		}

		distanceSinceMaintenance := v.CurrentOdometer - v.LastMaintenanceOdometer
		v.KmToNextMaintenance = v.MaintenanceIntervalKm - distanceSinceMaintenance

		if v.KmToNextMaintenance < 0 {
			v.MaintenanceStatus = "OVERDUE"
		} else if v.KmToNextMaintenance <= 1000 {
			v.MaintenanceStatus = "WARNING"
		} else {
			v.MaintenanceStatus = "OK"
		}

		// 🔥 РОЗРАХУНОК ДАТИ В GO:
		if v.AvgKmPerDay > 0 && v.KmToNextMaintenance > 0 {
			daysLeft := float64(v.KmToNextMaintenance) / v.AvgKmPerDay
			predictedDate := time.Now().Add(time.Duration(daysLeft*24) * time.Hour)
			v.PredictedMaintDate = &predictedDate
		}

		vehicles = append(vehicles, v)
	}

	return vehicles, rows.Err()
}

func (r *VehicleRepository) GetByID(ctx context.Context, id string, db DBExecutor) (*models.Vehicle, error) {
	args := []any{id}
	tFilter := tenantFilter(ctx, "v", "AND", &args)
	query := `
        WITH VehicleStats AS (
            SELECT 
                vehicle_id,
                MIN(created_at) as first_date,
                MAX(created_at) as last_date,
                MIN(odometer_km) as first_odo,
                MAX(odometer_km) as last_odo
            FROM fuel_records
            WHERE vehicle_id = $1 AND created_at >= CURRENT_DATE - INTERVAL '30 days' AND odometer_km IS NOT NULL
            GROUP BY vehicle_id
        )
        SELECT 
            v.id, v.brand, v.model, v.plate_number, v.type, v.capacity_kg, v.status, v.driver_id, u.full_name as driver_name, 
            v.home_warehouse_id, wh.name as home_warehouse_name,
            v.current_warehouse_id, wc.name as current_warehouse_name,
            v.tank_capacity, v.fuel_norm, 
            v.maintenance_interval_km, v.last_maintenance_odometer, v.created_at, v.updated_at,
            COALESCE((
                SELECT odometer_km FROM fuel_records 
                WHERE vehicle_id = v.id AND odometer_km IS NOT NULL 
                ORDER BY created_at DESC LIMIT 1
            ), 0) AS current_odometer,
            COALESCE(
                (vs.last_odo - vs.first_odo)::float / 
                NULLIF(EXTRACT(EPOCH FROM (vs.last_date - vs.first_date))/86400, 0), 
            0) as avg_km_per_day
        FROM vehicles v
        LEFT JOIN users u ON v.driver_id = u.id
        LEFT JOIN warehouses wh ON wh.id = v.home_warehouse_id
        LEFT JOIN warehouses wc ON wc.id = v.current_warehouse_id
        LEFT JOIN VehicleStats vs ON v.id = vs.vehicle_id
        WHERE v.id = $1` + tFilter + `
    `

	var v models.Vehicle
	err := db.QueryRow(ctx, query, args...).Scan(
		&v.ID, &v.Brand, &v.Model, &v.PlateNumber, &v.Type, &v.CapacityKg, &v.Status, &v.DriverID, &v.DriverName,
		&v.HomeWarehouseID, &v.HomeWarehouseName,
		&v.CurrentWarehouseID, &v.CurrentWarehouseName,
		&v.TankCapacity, &v.FuelNorm, &v.MaintenanceIntervalKm, &v.LastMaintenanceOdometer,
		&v.CreatedAt, &v.UpdatedAt, &v.CurrentOdometer, &v.AvgKmPerDay,
	)

	if err != nil {
		return nil, err
	}

	if v.CurrentOdometer < v.LastMaintenanceOdometer {
		v.CurrentOdometer = v.LastMaintenanceOdometer
	}
	distanceSinceMaintenance := v.CurrentOdometer - v.LastMaintenanceOdometer
	v.KmToNextMaintenance = v.MaintenanceIntervalKm - distanceSinceMaintenance

	if v.KmToNextMaintenance < 0 {
		v.MaintenanceStatus = "OVERDUE"
	} else if v.KmToNextMaintenance <= 1000 {
		v.MaintenanceStatus = "WARNING"
	} else {
		v.MaintenanceStatus = "OK"
	}

	// 🔥 РОЗРАХУНОК ДАТИ
	if v.AvgKmPerDay > 0 && v.KmToNextMaintenance > 0 {
		daysLeft := float64(v.KmToNextMaintenance) / v.AvgKmPerDay
		predictedDate := time.Now().Add(time.Duration(daysLeft*24) * time.Hour)
		v.PredictedMaintDate = &predictedDate
	}

	return &v, nil
}

func (r *VehicleRepository) UpdateStatus(ctx context.Context, id string, status string, reason string, db DBExecutor) error {
	args := []any{status, reason, id}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `
        UPDATE vehicles 
        SET status = $1, status_reason = $2, updated_at = CURRENT_TIMESTAMP
        WHERE id = $3` + tFilter

	cmdTag, err := db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("помилка оновлення статусу: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("авто з ID %s не знайдено для оновлення", id)
	}

	return nil
}

func (r *VehicleRepository) PerformMaintenance(ctx context.Context, record *models.MaintenanceRecord, db DBExecutor) error {
	b, ok := db.(interface {
		Begin(context.Context) (pgx.Tx, error)
	})
	if !ok {
		return fmt.Errorf("db не підтримує транзакції")
	}

	tx, err := b.Begin(ctx)
	if err != nil {
		return fmt.Errorf("помилка старту транзакції: %w", err)
	}
	defer tx.Rollback(ctx)

	tid := TenantFromCtx(ctx)
	if tid == "" {
		return fmt.Errorf("tenant_id is required for vehicle maintenance")
	}

	args := []any{record.OdometerKm, record.VehicleID}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	updateQuery := `
        UPDATE vehicles 
        SET last_maintenance_odometer = $1, status = 'ACTIVE', status_reason = NULL, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2` + tFilter
	cmdTag, err := tx.Exec(ctx, updateQuery, args...)
	if err != nil {
		return fmt.Errorf("помилка оновлення статусу авто: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("авто з ID %s не знайдено", record.VehicleID)
	}

	insertQuery := `
        INSERT INTO maintenance_records (vehicle_id, odometer_km, description, performed_by, cost_amount, document_url, driver_id, status, service_type)
        VALUES ($1, $2, $3, $4, $5, $6, (SELECT driver_id FROM vehicles WHERE id = $1 AND tenant_id = $7::uuid), 'COMPLETED', 'INSPECTION')
        RETURNING id, created_at
    `
	err = tx.QueryRow(ctx, insertQuery,
		record.VehicleID,
		record.OdometerKm,
		record.Description,
		record.PerformedBy,
		record.CostAmount,
		record.DocumentURL,
		tid,
	).Scan(&record.ID, &record.CreatedAt)

	if err != nil {
		return fmt.Errorf("помилка збереження історії робіт: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("помилка коміту транзакції ТО: %w", err)
	}

	return nil
}

func (r *VehicleRepository) ScheduleMaintenance(ctx context.Context, record *models.MaintenanceRecord, db DBExecutor) error {
	tid := TenantFromCtx(ctx)
	if tid == "" {
		return fmt.Errorf("tenant_id is required for vehicle maintenance")
	}
	if record.ScheduledFor == nil {
		return fmt.Errorf("scheduled_for is required")
	}

	query := `
        INSERT INTO maintenance_records (vehicle_id, odometer_km, description, performed_by, cost_amount, document_url, driver_id, status, service_type, scheduled_for)
        SELECT v.id, $2, $3, $4, 0, '', v.driver_id, 'SCHEDULED', $5, $6
        FROM vehicles v
        WHERE v.id = $1 AND v.tenant_id = $7::uuid
        RETURNING id, created_at
    `
	err := db.QueryRow(ctx, query,
		record.VehicleID,
		record.OdometerKm,
		record.Description,
		record.PerformedBy,
		record.ServiceType,
		record.ScheduledFor,
		tid,
	).Scan(&record.ID, &record.CreatedAt)
	if err != nil {
		return fmt.Errorf("помилка планування ТО: %w", err)
	}

	record.Status = "SCHEDULED"
	record.CostAmount = 0
	record.DocumentURL = ""
	return nil
}

func (r *VehicleRepository) GetMaintenanceHistory(ctx context.Context, vehicleID string, db DBExecutor) ([]*models.MaintenanceRecord, error) {
	args := []any{vehicleID}
	tFilter := tenantFilter(ctx, "v", "AND", &args)
	query := `
        SELECT m.id, m.vehicle_id, m.odometer_km, m.description, m.performed_by, m.cost_amount, COALESCE(m.document_url, ''), m.created_at, u.full_name,
               COALESCE(m.status, 'COMPLETED'), COALESCE(m.service_type, 'INSPECTION'), m.scheduled_for
        FROM maintenance_records m
        JOIN vehicles v ON v.id = m.vehicle_id
        LEFT JOIN users u ON m.driver_id = u.id AND u.tenant_id = v.tenant_id
        WHERE m.vehicle_id = $1` + tFilter + `
        ORDER BY m.created_at DESC
    `
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*models.MaintenanceRecord
	for rows.Next() {
		var rec models.MaintenanceRecord
		err := rows.Scan(
			&rec.ID,
			&rec.VehicleID,
			&rec.OdometerKm,
			&rec.Description,
			&rec.PerformedBy,
			&rec.CostAmount,
			&rec.DocumentURL,
			&rec.CreatedAt,
			&rec.DriverName,
			&rec.Status,
			&rec.ServiceType,
			&rec.ScheduledFor,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, &rec)
	}

	return records, nil
}

func (r *VehicleRepository) AssignDriver(ctx context.Context, vehicleID string, driverID *string, db DBExecutor) error {
	tid := TenantFromCtx(ctx)
	if tid == "" {
		return fmt.Errorf("tenant_id is required for vehicles")
	}
	if driverID != nil && *driverID != "" {
		var ok bool
		if err := db.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM users
				WHERE id = $1 AND tenant_id = $2 AND status = $3
			)
		`, *driverID, tid, models.StatusActive).Scan(&ok); err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("водія не знайдено у вашій організації")
		}
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	args := []any{driverID, vehicleID}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	tag, err := tx.Exec(ctx, `UPDATE vehicles SET driver_id = $1 WHERE id = $2`+tFilter, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("авто не знайдено у вашій організації")
	}

	_, err = tx.Exec(ctx, `INSERT INTO vehicle_driver_history (vehicle_id, driver_id) VALUES ($1, $2)`, vehicleID, driverID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *VehicleRepository) GetDriverHistory(ctx context.Context, vehicleID string, db DBExecutor) ([]models.VehicleDriverHistory, error) {
	args := []any{vehicleID}
	tFilter := tenantFilter(ctx, "v", "AND", &args)
	query := `
		SELECT h.id, h.vehicle_id, h.driver_id, COALESCE(u.full_name, 'Без закріплення (Резерв)'), h.assigned_at
		FROM vehicle_driver_history h
		JOIN vehicles v ON v.id = h.vehicle_id
		LEFT JOIN users u ON h.driver_id = u.id
		WHERE h.vehicle_id = $1` + tFilter + `
		ORDER BY h.assigned_at DESC
	`
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.VehicleDriverHistory
	for rows.Next() {
		var rec models.VehicleDriverHistory
		if err := rows.Scan(&rec.ID, &rec.VehicleID, &rec.DriverID, &rec.DriverName, &rec.AssignedAt); err != nil {
			return nil, err
		}
		history = append(history, rec)
	}
	return history, nil
}

// Update оновлює базові параметри автомобіля
func (r *VehicleRepository) Update(ctx context.Context, db DBExecutor, id string, brand string, model string, plateNumber string, capacityKg float64) error {
	args := []any{brand, model, plateNumber, capacityKg, id}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `UPDATE vehicles SET brand = $1, model = $2, plate_number = $3, capacity_kg = $4 WHERE id = $5` + tFilter
	_, err := db.Exec(ctx, query, args...)
	return err
}

// Delete безповоротно видаляє транспортний засіб
func (r *VehicleRepository) Delete(ctx context.Context, db DBExecutor, id string) error {
	args := []any{id}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `DELETE FROM vehicles WHERE id = $1` + tFilter
	_, err := db.Exec(ctx, query, args...)
	return err
}

// GetAvailableForRoute повертає вільні авто, що фізично знаходяться на складі відправника АБО отримувача.
// Також включає машини, чий home_warehouse — один з цих складів (навіть якщо current = NULL).
func (r *VehicleRepository) GetAvailableForRoute(ctx context.Context, db DBExecutor, fromWarehouseID string, toWarehouseID string) ([]models.Vehicle, error) {
	args := []any{fromWarehouseID, toWarehouseID}
	tFilter := tenantFilter(ctx, "v", "AND", &args)
	query := `
		SELECT v.id, v.brand, v.model, v.plate_number, v.type, v.capacity_kg, 
		       v.status, v.driver_id, v.tank_capacity, v.fuel_norm, 
		       v.maintenance_interval_km, v.last_maintenance_odometer, 
		       v.created_at, v.updated_at,
		       GREATEST(0, COALESCE((
		         SELECT SUM(CASE WHEN record_type = 'REFUEL' THEN liters ELSE -liters END)
		         FROM fuel_records WHERE vehicle_id = v.id
		       ), 0)) AS current_fuel_liters
		FROM vehicles v
		WHERE v.status = 'ACTIVE'
		  AND v.driver_id IS NOT NULL
		  AND v.type IN ('VAN', 'TRUCK', 'PICKUP')
		  AND (
		    v.current_warehouse_id IN ($1, $2)
		    OR (v.current_warehouse_id IS NULL AND v.home_warehouse_id IN ($1, $2))
		  )` + tFilter + `
		ORDER BY v.capacity_kg ASC
	`

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Vehicle
	for rows.Next() {
		var v models.Vehicle
		err := rows.Scan(
			&v.ID, &v.Brand, &v.Model, &v.PlateNumber, &v.Type, &v.CapacityKg,
			&v.Status, &v.DriverID, &v.TankCapacity, &v.FuelNorm,
			&v.MaintenanceIntervalKm, &v.LastMaintenanceOdometer,
			&v.CreatedAt, &v.UpdatedAt, &v.CurrentFuelLiters,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, v)
	}

	return list, rows.Err()
}

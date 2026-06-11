package repositories

import (
	"Omnilog_backend/internal/models"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type FuelRepository struct{}

func NewFuelRepository() *FuelRepository {
	return &FuelRepository{}
}

// Допоміжний інтерфейс для старту транзакцій з DBExecutor
type beginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

func (r *FuelRepository) CreateFuelRecord(ctx context.Context, record *models.FuelRecord, db DBExecutor) error {
	// 🛡️ ДОДАНО: Context timeout (10 сек) щоб запобігти infinite waits
	var cancel context.CancelFunc
	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	b, ok := db.(beginner)
	if !ok {
		return fmt.Errorf("переданий db не підтримує старт транзакцій")
	}

	tx, err := b.Begin(ctx)
	if err != nil {
		return fmt.Errorf("помилка старту транзакції: %w", err)
	}
	defer tx.Rollback(ctx)

	tid := TenantFromCtx(ctx)
	if tid == "" {
		return fmt.Errorf("tenant_id is required for fuel records")
	}

	// 1. Отримуємо дані авто та БЛОКУЄМО рядок (FOR UPDATE), щоб ніхто не змінив дані паралельно
	var fuelNorm float64
	var tankCapacity float64
	err = tx.QueryRow(ctx, "SELECT fuel_norm, tank_capacity FROM vehicles WHERE id = $1 AND tenant_id = $2::uuid FOR UPDATE", record.VehicleID, tid).Scan(&fuelNorm, &tankCapacity)
	if err != nil {
		return fmt.Errorf("помилка отримання даних авто: %w", err)
	}

	// 2. Вираховуємо поточний залишок пального "Віртуальний бак"
	var currentBalance float64
	balanceQuery := `
		SELECT COALESCE(SUM(CASE WHEN record_type = 'REFUEL' THEN liters ELSE -liters END), 0)
		FROM fuel_records
		WHERE vehicle_id = $1 AND tenant_id = $2::uuid
	`
	err = tx.QueryRow(ctx, balanceQuery, record.VehicleID, tid).Scan(&currentBalance)
	if err != nil {
		return fmt.Errorf("помилка розрахунку залишку пального: %w", err)
	}

	// 3. БІЗНЕС-ЛОГІКА: Захист від "мінуса" та "переливу" (Hard Stop)
	if record.RecordType == models.FuelExpense {
		if record.Liters > currentBalance {
			return fmt.Errorf("недостатньо пального. У баку: %.1f л, спроба списати: %.1f л", currentBalance, record.Liters)
		}
	} else if record.RecordType == models.FuelRefuel {
		// Якщо заправка переповнює бак — обрізаємо до місткості (не повертаємо помилку)
		maxRefuel := tankCapacity - currentBalance
		if maxRefuel <= 0 {
			return fmt.Errorf("бак повний! Поточний залишок: %.1f л, максимум: %.1f л", currentBalance, tankCapacity)
		}
		if record.Liters > maxRefuel {
			record.Liters = maxRefuel
		}
	}

	// 3.5 БІЗНЕС-ЛОГІКА: Захист від скручування одометра (Hard Stop)
	var lastOdometer int
	var hasLastOdometer bool

	if record.OdometerKm != nil {
		err = tx.QueryRow(ctx,
			`SELECT odometer_km FROM fuel_records 
             WHERE vehicle_id = $1 AND tenant_id = $2::uuid AND odometer_km IS NOT NULL 
             ORDER BY created_at DESC LIMIT 1`,
			record.VehicleID, tid,
		).Scan(&lastOdometer)

		if err == nil {
			hasLastOdometer = true
			if *record.OdometerKm < lastOdometer {
				// Жорстко блокуємо запис, якщо одометр пішов у мінус
				return fmt.Errorf("помилка одометра: поточний (%d км) не може бути меншим за попередній (%d км)", *record.OdometerKm, lastOdometer)
			}
		} else if err != pgx.ErrNoRows {
			log.Printf("Помилка отримання попереднього одометра: %v", err)
		}
	}

	// 4. ДЕТЕКТОР АНОМАЛІЙ (Soft Stop - підозріла поведінка)
	if record.RecordType == models.FuelExpense && record.OdometerKm != nil && hasLastOdometer {
		distance := *record.OdometerKm - lastOdometer

		if distance == 0 {
			// Пальне спалено, але машина не рухалась → весь обсяг вважаємо втратою
			record.IsAnomaly = true
			record.AnomalyExcessLiters = record.Liters
			reason := fmt.Sprintf("Витрата пального без руху (холостий хід / обігрів). Одометр не змінився: %d км", lastOdometer)
			record.AnomalyReason = &reason
		} else if fuelNorm > 0 {
			actualConsumption := (record.Liters / float64(distance)) * 100

			if actualConsumption > (fuelNorm * 1.3) {
				// Перевитрата більше ніж на 30% від норми.
				// Втратою вважаємо лише обсяг ПОНАД норму, а не всі літри.
				record.IsAnomaly = true
				expectedLiters := fuelNorm * float64(distance) / 100.0
				excess := record.Liters - expectedLiters
				if excess < 0 {
					excess = 0
				}
				record.AnomalyExcessLiters = excess
				reason := fmt.Sprintf("Перевитрата пального! Факт: %.1f л/100км (Норма: %.1f). Дистанція: %d км. Зайве: %.1f л", actualConsumption, fuelNorm, distance, excess)
				record.AnomalyReason = &reason
			}
		}
	}
	// 5. Запис у базу
	query := `
        INSERT INTO fuel_records (
            vehicle_id, liters, odometer_km, record_type, created_by, is_anomaly, anomaly_reason, anomaly_excess_liters, tenant_id
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, created_at`
	err = tx.QueryRow(ctx, query,
		record.VehicleID, record.Liters, record.OdometerKm, record.RecordType,
		record.CreatedBy, record.IsAnomaly, record.AnomalyReason, record.AnomalyExcessLiters, tid,
	).Scan(&record.ID, &record.CreatedAt)

	if err != nil {
		return fmt.Errorf("помилка створення запису пального: %w", err)
	}

	// 6. Фіксуємо транзакцію (знімаємо блокування FOR UPDATE)
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("помилка коміту транзакції: %w", err)
	}

	return nil
}

func (r *FuelRepository) GetRecordsByVehicleID(ctx context.Context, vehicleID string, db DBExecutor) ([]*models.FuelRecord, error) {
	args := []any{vehicleID}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `
        SELECT id, vehicle_id, liters, odometer_km, record_type, is_anomaly, anomaly_reason, created_by, created_at 
        FROM fuel_records 
        WHERE vehicle_id = $1` + tFilter + `
        ORDER BY created_at DESC
    `

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*models.FuelRecord
	for rows.Next() {
		var record models.FuelRecord
		err := rows.Scan(
			&record.ID,
			&record.VehicleID,
			&record.Liters,
			&record.OdometerKm,
			&record.RecordType,
			&record.IsAnomaly,
			&record.AnomalyReason,
			&record.CreatedBy,
			&record.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	return records, nil
}

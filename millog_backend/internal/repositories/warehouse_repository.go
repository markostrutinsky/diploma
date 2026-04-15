package repositories

import (
	"context"
	"errors"
	"millog_backend/internal/models"
)

type WarehouseRepository struct{}

func NewWarehouseRepository() *WarehouseRepository {
	return &WarehouseRepository{}
}

func (r *WarehouseRepository) Create(ctx context.Context, db DBExecutor, w *models.Warehouse) error {
	query := `INSERT INTO warehouses (unit_id, name, location_type, latitude, longitude)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`

	return db.QueryRow(ctx, query, w.UnitID, w.Name, w.LocationType, w.Latitude, w.Longitude).
		Scan(&w.ID, &w.CreatedAt, &w.UpdatedAt)
}
func (r *WarehouseRepository) ListByUnit(ctx context.Context, db DBExecutor, unitID int64) ([]models.Warehouse, error) {
	var list []models.Warehouse

	if unitID == 0 {
		// Для Адміна - показуємо взагалі всі склади без ієрархії
		query := `
            SELECT id, unit_id, name, location_type, latitude, longitude, created_at, updated_at
            FROM warehouses
            ORDER BY name`

		rows, err := db.Query(ctx, query)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var w models.Warehouse
			if err := rows.Scan(&w.ID, &w.UnitID, &w.Name, &w.LocationType, &w.Latitude, &w.Longitude, &w.CreatedAt, &w.UpdatedAt); err != nil {
				return nil, err
			}
			list = append(list, w)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}

	} else {
		// Для командирів - склади їхнього підрозділу ТА підлеглих (рекурсія)
		query := `
            WITH RECURSIVE unit_hierarchy AS (
                SELECT id FROM units WHERE id = $1
                UNION ALL
                SELECT u.id FROM units u
                JOIN unit_hierarchy uh ON u.parent_id = uh.id
            )
            SELECT id, unit_id, name, location_type, latitude, longitude, created_at, updated_at
            FROM warehouses
            WHERE unit_id IN (SELECT id FROM unit_hierarchy)
            ORDER BY name`

		rows, err := db.Query(ctx, query, unitID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var w models.Warehouse
			if err := rows.Scan(&w.ID, &w.UnitID, &w.Name, &w.LocationType, &w.Latitude, &w.Longitude, &w.CreatedAt, &w.UpdatedAt); err != nil {
				return nil, err
			}
			list = append(list, w)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	if list == nil {
		list = []models.Warehouse{}
	}

	return list, nil
}

// Додай цей метод до WarehouseRepository
func (r *WarehouseRepository) UpdateLocation(ctx context.Context, db DBExecutor, warehouseID string, lat, lng float64) error {
	query := `
		UPDATE warehouses 
		SET latitude = $1, longitude = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $3
	`

	result, err := db.Exec(ctx, query, lat, lng, warehouseID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("склад не знайдено")
	}

	return nil
}

// Update оновлює базові параметри складу
func (r *WarehouseRepository) Update(ctx context.Context, db DBExecutor, id string, name string, capacityLevel string, zoneType string) error {
	query := `UPDATE warehouses SET name = $1, capacity_level = $2, zone_type = $3 WHERE id = $4`
	_, err := db.Exec(ctx, query, name, capacityLevel, zoneType, id)
	return err
}

// Delete безповоротно видаляє склад
func (r *WarehouseRepository) Delete(ctx context.Context, db DBExecutor, id string) error {
	query := `DELETE FROM warehouses WHERE id = $1`
	_, err := db.Exec(ctx, query, id)
	return err
}

package repositories

import (
	"context"
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

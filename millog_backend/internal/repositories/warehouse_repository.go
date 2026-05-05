package repositories

import (
	"Omnilog_backend/internal/models"
	"context"
	"errors"
)

type WarehouseRepository struct{}

func NewWarehouseRepository() *WarehouseRepository {
	return &WarehouseRepository{}
}

func (r *WarehouseRepository) Create(ctx context.Context, db DBExecutor, w *models.Warehouse) error {
	tid := TenantFromCtx(ctx)
	if tid == "" {
		query := `INSERT INTO warehouses (unit_id, name, location_type, latitude, longitude)
		          VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
		return db.QueryRow(ctx, query, w.UnitID, w.Name, w.LocationType, w.Latitude, w.Longitude).
			Scan(&w.ID, &w.CreatedAt, &w.UpdatedAt)
	}
	query := `INSERT INTO warehouses (unit_id, name, location_type, latitude, longitude, tenant_id)
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`
	return db.QueryRow(ctx, query, w.UnitID, w.Name, w.LocationType, w.Latitude, w.Longitude, tid).
		Scan(&w.ID, &w.CreatedAt, &w.UpdatedAt)
}
func (r *WarehouseRepository) ListByUnit(ctx context.Context, db DBExecutor, unitID int64) ([]models.Warehouse, error) {
	var list []models.Warehouse

	if unitID == 0 {
		// Для Адміна - всі склади поточного tenant
		args := []any{}
		tFilter := tenantFilter(ctx, "", "WHERE", &args)
		query := `
            SELECT id, unit_id, name, location_type, latitude, longitude, created_at, updated_at
            FROM warehouses` + tFilter + `
            ORDER BY name`

		rows, err := db.Query(ctx, query, args...)
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
		args := []any{unitID}
		tFilter := tenantFilter(ctx, "w", "AND", &args)
		query := `
            WITH RECURSIVE 
            -- Нащадки (вниз по ієрархії)
            descendants AS (
                SELECT id FROM units WHERE id = $1
                UNION ALL
                SELECT u.id FROM units u
                JOIN descendants d ON u.parent_id = d.id
            ),
            -- Предки (вгору по ієрархії) 
            ancestors AS (
                SELECT id, parent_id FROM units WHERE id = $1
                UNION ALL
                SELECT u.id, u.parent_id FROM units u
                JOIN ancestors a ON u.id = a.parent_id
            )
            SELECT w.id, w.unit_id, w.name, w.location_type, w.latitude, w.longitude, w.created_at, w.updated_at
            FROM warehouses w
            WHERE w.unit_id IN (
                SELECT id FROM descendants
                UNION
                SELECT id FROM ancestors
            )` + tFilter + `
            ORDER BY w.name`

		rows, err := db.Query(ctx, query, args...)
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
	args := []any{lat, lng, warehouseID}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `
		UPDATE warehouses 
		SET latitude = $1, longitude = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $3` + tFilter

	result, err := db.Exec(ctx, query, args...)
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
	args := []any{name, capacityLevel, zoneType, id}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `UPDATE warehouses SET name = $1, capacity_level = $2, zone_type = $3 WHERE id = $4` + tFilter
	_, err := db.Exec(ctx, query, args...)
	return err
}

// Delete безповоротно видаляє склад
func (r *WarehouseRepository) Delete(ctx context.Context, db DBExecutor, id string) error {
	args := []any{id}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `DELETE FROM warehouses WHERE id = $1` + tFilter
	_, err := db.Exec(ctx, query, args...)
	return err
}

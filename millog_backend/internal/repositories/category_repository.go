package repositories

import (
	"context"

	"millog_backend/internal/models"
)

type CategoryRepository struct{}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

func (r *CategoryRepository) Create(ctx context.Context, db DBExecutor, c *models.ResourceCategory) error {
	tid := TenantFromCtx(ctx)
	if tid == "" {
		query := `INSERT INTO resource_categories (name, description) VALUES ($1, $2) RETURNING id, created_at`
		return db.QueryRow(ctx, query, c.Name, c.Description).Scan(&c.ID, &c.CreatedAt)
	}
	query := `INSERT INTO resource_categories (name, description, tenant_id) VALUES ($1, $2, $3) RETURNING id, created_at`
	return db.QueryRow(ctx, query, c.Name, c.Description, tid).Scan(&c.ID, &c.CreatedAt)
}

func (r *CategoryRepository) List(ctx context.Context, db DBExecutor) ([]models.ResourceCategory, error) {
	args := []any{}
	tFilter := tenantFilter(ctx, "", "WHERE", &args)
	q := `SELECT id, name, description, created_at FROM resource_categories` + tFilter + ` ORDER BY name`
	rows, err := db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.ResourceCategory
	for rows.Next() {
		var c models.ResourceCategory
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (r *CategoryRepository) Update(ctx context.Context, db DBExecutor, id string, name string, description string) error {
	args := []any{name, description, id}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `UPDATE resource_categories SET name = $1, description = $2 WHERE id = $3` + tFilter
	_, err := db.Exec(ctx, query, args...)
	return err
}

func (r *CategoryRepository) Delete(ctx context.Context, db DBExecutor, id string) error {
	args := []any{id}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `DELETE FROM resource_categories WHERE id = $1` + tFilter
	_, err := db.Exec(ctx, query, args...)
	return err
}

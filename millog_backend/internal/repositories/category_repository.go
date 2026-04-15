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
	query := `INSERT INTO resource_categories (name, description) VALUES ($1, $2) RETURNING id, created_at`
	return db.QueryRow(ctx, query, c.Name, c.Description).Scan(&c.ID, &c.CreatedAt)
}

func (r *CategoryRepository) List(ctx context.Context, db DBExecutor) ([]models.ResourceCategory, error) {
	rows, err := db.Query(ctx, `SELECT id, name, description, created_at FROM resource_categories ORDER BY name`)
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

// Update оновлює дані категорії
func (r *CategoryRepository) Update(ctx context.Context, db DBExecutor, id string, name string, description string) error {
	query := `UPDATE resource_categories SET name = $1, description = $2 WHERE id = $3`
	_, err := db.Exec(ctx, query, name, description, id)
	return err
}

// Delete безповоротно видаляє категорію
func (r *CategoryRepository) Delete(ctx context.Context, db DBExecutor, id string) error {
	query := `DELETE FROM resource_categories WHERE id = $1`
	_, err := db.Exec(ctx, query, id)
	return err
}

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

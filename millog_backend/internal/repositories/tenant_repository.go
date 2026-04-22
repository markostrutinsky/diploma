package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"millog_backend/internal/models"
)

type TenantRepository struct{}

func NewTenantRepository() *TenantRepository {
	return &TenantRepository{}
}

// List повертає всі tenant-и платформи (для SYSTEM_ADMIN).
// Підтримує фільтр за пошуком (name, slug).
func (r *TenantRepository) List(ctx context.Context, db DBExecutor, search string) ([]models.Tenant, error) {
	args := []any{}
	where := ""
	if s := strings.TrimSpace(search); s != "" {
		args = append(args, "%"+strings.ToLower(s)+"%")
		where = fmt.Sprintf(" WHERE LOWER(name) LIKE $%d OR LOWER(slug) LIKE $%d", len(args), len(args))
	}
	q := `SELECT id, name, slug, subscription_tier, subscription_expires_at,
			owner_email, is_active, created_at, updated_at
		FROM tenants` + where + ` ORDER BY created_at DESC`
	rows, err := db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Tenant
	for rows.Next() {
		var t models.Tenant
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.SubscriptionTier,
			&t.SubscriptionExpiresAt, &t.OwnerEmail, &t.IsActive,
			&t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *TenantRepository) GetByID(ctx context.Context, db DBExecutor, id string) (*models.Tenant, error) {
	q := `SELECT id, name, slug, subscription_tier, subscription_expires_at,
			owner_email, is_active, created_at, updated_at
		FROM tenants WHERE id = $1`
	var t models.Tenant
	err := db.QueryRow(ctx, q, id).Scan(&t.ID, &t.Name, &t.Slug, &t.SubscriptionTier,
		&t.SubscriptionExpiresAt, &t.OwnerEmail, &t.IsActive,
		&t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TenantRepository) GetBySlug(ctx context.Context, db DBExecutor, slug string) (*models.Tenant, error) {
	q := `SELECT id, name, slug, subscription_tier, subscription_expires_at,
			owner_email, is_active, created_at, updated_at
		FROM tenants WHERE slug = $1`
	var t models.Tenant
	err := db.QueryRow(ctx, q, slug).Scan(&t.ID, &t.Name, &t.Slug, &t.SubscriptionTier,
		&t.SubscriptionExpiresAt, &t.OwnerEmail, &t.IsActive,
		&t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// UpdateTier встановлює тариф тенанта і опційно дату завершення підписки.
func (r *TenantRepository) UpdateTier(ctx context.Context, db DBExecutor, id string,
	tier models.SubscriptionTier, expiresAt *time.Time) error {
	q := `UPDATE tenants SET subscription_tier = $1, subscription_expires_at = $2,
			updated_at = CURRENT_TIMESTAMP WHERE id = $3`
	_, err := db.Exec(ctx, q, tier, expiresAt, id)
	return err
}

// SetActive призупиняє (false) або відновлює (true) організацію.
func (r *TenantRepository) SetActive(ctx context.Context, db DBExecutor, id string, active bool) error {
	_, err := db.Exec(ctx, `UPDATE tenants SET is_active = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, active, id)
	return err
}

// CountUsers — кількість активних користувачів у tenant-і.
func (r *TenantRepository) CountUsers(ctx context.Context, db DBExecutor, tenantID string) (int, error) {
	var n int
	err := db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE tenant_id = $1`, tenantID).Scan(&n)
	return n, err
}

// Delete видаляє tenant разом з усіма пов'язаними даними (через ON DELETE CASCADE).
// НЕБЕЗПЕЧНА операція — тільки SYSTEM_ADMIN.
func (r *TenantRepository) Delete(ctx context.Context, db DBExecutor, id string) error {
	_, err := db.Exec(ctx, `DELETE FROM tenants WHERE id = $1`, id)
	return err
}

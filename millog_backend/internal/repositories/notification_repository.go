package repositories

import (
	"context"
	"fmt"
	"log"

	"Omnilog_backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepository struct{}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{}
}

// Create створює нове сповіщення
func (r *NotificationRepository) Create(ctx context.Context, db *pgxpool.Pool, notif *models.Notification) error {
	tid := TenantFromCtx(ctx)
	if tid == "" {
		return fmt.Errorf("tenant_id is required for creating notifications")
	}

	log.Printf("DEBUG: Creating notification - userID=%s, tenantID=%s, type=%s, title=%s", notif.UserID, tid, notif.Type, notif.Title)

	query := `
		INSERT INTO notifications (user_id, tenant_id, type, title, message, related_id, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)
		RETURNING id, created_at
	`

	err := db.QueryRow(ctx, query,
		notif.UserID, tid, notif.Type, notif.Title, notif.Message, notif.RelatedID, notif.IsRead,
	).Scan(&notif.ID, &notif.CreatedAt)

	log.Printf("DEBUG: Notification created - id=%s, err=%v", notif.ID, err)
	return err
}

// ListByUser отримує всі сповіщення користувача
func (r *NotificationRepository) ListByUser(ctx context.Context, db *pgxpool.Pool, userID string, limit int) ([]models.Notification, error) {
	args := []any{userID}
	tFilter := tenantFilter(ctx, "n", "AND", &args)

	if limit <= 0 {
		limit = 50
	}
	args = append(args, limit)

	query := fmt.Sprintf(`
		SELECT id, user_id, tenant_id, type, title, message, related_id, is_read, created_at, read_at
		FROM notifications n
		WHERE user_id = $1 %s
		ORDER BY created_at DESC
		LIMIT $%d
	`, tFilter, len(args))

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		err := rows.Scan(
			&n.ID, &n.UserID, &n.TenantID, &n.Type, &n.Title, &n.Message,
			&n.RelatedID, &n.IsRead, &n.CreatedAt, &n.ReadAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

// GetUnreadCount повертає кількість непрочитаних сповіщень
func (r *NotificationRepository) GetUnreadCount(ctx context.Context, db *pgxpool.Pool, userID string) (int, error) {
	args := []any{userID}
	tFilter := tenantFilter(ctx, "n", "AND", &args)

	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM notifications n
		WHERE user_id = $1 AND is_read = FALSE %s
	`, tFilter)

	var count int
	err := db.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

// MarkAsRead позначає сповіщення як прочитане
func (r *NotificationRepository) MarkAsRead(ctx context.Context, db *pgxpool.Pool, notificationID string, userID string) error {
	args := []any{notificationID, userID}
	tFilter := tenantFilter(ctx, "", "AND", &args)

	query := fmt.Sprintf(`
		UPDATE notifications
		SET is_read = TRUE, read_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND user_id = $2 %s
	`, tFilter)

	result, err := db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

// MarkAllAsRead позначає всі сповіщення користувача як прочитані
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, db *pgxpool.Pool, userID string) error {
	args := []any{userID}
	tFilter := tenantFilter(ctx, "", "AND", &args)

	query := fmt.Sprintf(`
		UPDATE notifications
		SET is_read = TRUE, read_at = CURRENT_TIMESTAMP
		WHERE user_id = $1 AND is_read = FALSE %s
	`, tFilter)

	_, err := db.Exec(ctx, query, args...)
	return err
}

// Delete видаляє сповіщення
func (r *NotificationRepository) Delete(ctx context.Context, db *pgxpool.Pool, notificationID string, userID string) error {
	args := []any{notificationID, userID}
	tFilter := tenantFilter(ctx, "", "AND", &args)

	query := fmt.Sprintf(`
		DELETE FROM notifications
		WHERE id = $1 AND user_id = $2 %s
	`, tFilter)

	result, err := db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

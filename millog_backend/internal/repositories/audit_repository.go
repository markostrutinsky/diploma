package repositories

import (
	"context"
	"fmt"
	"millog_backend/internal/models"
)

type AuditLogRepository struct{}

func NewAuditLogRepository() *AuditLogRepository {
	return &AuditLogRepository{}
}

// LogAction - універсальна функція для запису будь-якої дії
func (r *AuditLogRepository) LogAction(ctx context.Context, db DBExecutor, userID, actionType, entityType, entityID, details string) error {
	tid := TenantFromCtx(ctx)
	if tid == "" {
		query := `
			INSERT INTO audit_logs (user_id, action_type, entity_type, entity_id, details, tenant_id)
			SELECT $1, $2, $3, $4, $5, u.tenant_id FROM users u WHERE u.id = $1
		`
		_, err := db.Exec(ctx, query, userID, actionType, entityType, entityID, details)
		return err
	}
	query := `
		INSERT INTO audit_logs (user_id, action_type, entity_type, entity_id, details, tenant_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := db.Exec(ctx, query, userID, actionType, entityType, entityID, details, tid)
	return err
}

// GetLogs - отримуємо історію для адмінів (з прив'язкою до юзерів)
func (r *AuditLogRepository) GetLogs(ctx context.Context, db DBExecutor, limit int) ([]models.AuditLog, error) {
	args := []any{}
	tFilter := tenantFilter(ctx, "a", "WHERE", &args)
	args = append(args, limit)
	query := `
		SELECT 
			a.id, COALESCE(u.email, 'Видалений користувач'), COALESCE(u.role, 'UNKNOWN'), 
			a.action_type, a.entity_type, a.entity_id, a.details, a.created_at
		FROM audit_logs a
		LEFT JOIN users u ON a.user_id = u.id` + tFilter + `
		ORDER BY a.created_at DESC
		LIMIT $` + fmt.Sprintf("%d", len(args)) + `
	`
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var l models.AuditLog
		if err := rows.Scan(&l.ID, &l.UserEmail, &l.UserRole, &l.ActionType, &l.EntityType, &l.EntityID, &l.Details, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

package repositories

import (
	"context"
	"errors"
	"fmt"

	"Omnilog_backend/internal/models"

	"github.com/jackc/pgx/v5"
)

type ContractorMembershipRepository struct{}

func NewContractorMembershipRepository() *ContractorMembershipRepository {
	return &ContractorMembershipRepository{}
}

// Apply створює (за потреби) заявку підрядника на співпрацю з організацією.
// Ідемпотентна: якщо запис уже існує — не змінює його статус (повторна спроба «взяти»
// завдання не скидає вже наявне APPROVED/REJECTED у PENDING).
// Повертає поточний статус членства після виклику.
func (r *ContractorMembershipRepository) Apply(ctx context.Context, db DBExecutor, contractorID, tenantID string) (models.ContractorMembershipStatus, error) {
	var ok bool
	if err := db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM users u
			JOIN tenants t ON t.id = $2
			WHERE u.id = $1
			  AND u.role = $3
			  AND u.status = $4
			  AND t.is_active = TRUE
		)
	`, contractorID, tenantID, models.RoleContractor, models.StatusActive).Scan(&ok); err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("підрядника або активну організацію не знайдено")
	}

	_, err := db.Exec(ctx, `
		INSERT INTO contractor_memberships (contractor_id, tenant_id, status)
		VALUES ($1, $2, 'PENDING')
		ON CONFLICT (contractor_id, tenant_id) DO NOTHING
	`, contractorID, tenantID)
	if err != nil {
		return "", err
	}
	return r.GetStatus(ctx, db, contractorID, tenantID)
}

// GetStatus повертає статус членства (або "" якщо запис відсутній).
func (r *ContractorMembershipRepository) GetStatus(ctx context.Context, db DBExecutor, contractorID, tenantID string) (models.ContractorMembershipStatus, error) {
	var status models.ContractorMembershipStatus
	err := db.QueryRow(ctx, `
		SELECT status FROM contractor_memberships
		WHERE contractor_id = $1 AND tenant_id = $2
	`, contractorID, tenantID).Scan(&status)
	if err != nil {
		// Немає запису членства — для нас це валідний стан «ще не подавався».
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return status, nil
}

// IsApproved — швидка перевірка, чи може підрядник брати завдання організації.
func (r *ContractorMembershipRepository) IsApproved(ctx context.Context, db DBExecutor, contractorID, tenantID string) (bool, error) {
	status, err := r.GetStatus(ctx, db, contractorID, tenantID)
	if err != nil {
		return false, err
	}
	return status == models.MembershipApproved, nil
}

// ListByTenant повертає всі членства організації (для адмін-панелі), з даними підрядника.
// Якщо status != "" — фільтрує за статусом.
func (r *ContractorMembershipRepository) ListByTenant(ctx context.Context, db DBExecutor, tenantID string, status models.ContractorMembershipStatus) ([]models.ContractorMembership, error) {
	query := `
		SELECT m.id, m.contractor_id, m.tenant_id, m.status, m.note,
		       m.requested_at, m.decided_at, m.decided_by,
		       u.full_name, u.email, u.phone
		FROM contractor_memberships m
		JOIN users u ON u.id = m.contractor_id
		WHERE m.tenant_id = $1
	`
	args := []any{tenantID}
	if status != "" {
		query += " AND m.status = $2"
		args = append(args, status)
	}
	query += " ORDER BY m.requested_at DESC"

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.ContractorMembership
	for rows.Next() {
		var m models.ContractorMembership
		if err := rows.Scan(
			&m.ID, &m.ContractorID, &m.TenantID, &m.Status, &m.Note,
			&m.RequestedAt, &m.DecidedAt, &m.DecidedBy,
			&m.ContractorName, &m.ContractorEmail, &m.ContractorPhone,
		); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

// ListByContractor повертає членства конкретного підрядника (self-view), з назвами організацій.
func (r *ContractorMembershipRepository) ListByContractor(ctx context.Context, db DBExecutor, contractorID string) ([]models.ContractorMembership, error) {
	rows, err := db.Query(ctx, `
		SELECT m.id, m.contractor_id, m.tenant_id, m.status, m.note,
		       m.requested_at, m.decided_at, m.decided_by, t.name
		FROM contractor_memberships m
		JOIN tenants t ON t.id = m.tenant_id
		WHERE m.contractor_id = $1 AND t.is_active = TRUE
		ORDER BY m.requested_at DESC
	`, contractorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.ContractorMembership
	for rows.Next() {
		var m models.ContractorMembership
		if err := rows.Scan(
			&m.ID, &m.ContractorID, &m.TenantID, &m.Status, &m.Note,
			&m.RequestedAt, &m.DecidedAt, &m.DecidedBy, &m.TenantName,
		); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

// Decide змінює статус членства (APPROVED/REJECTED) у межах організації.
// tenantID гарантує, що адмін однієї організації не керує членствами іншої.
func (r *ContractorMembershipRepository) Decide(ctx context.Context, db DBExecutor, membershipID, tenantID, deciderID string, status models.ContractorMembershipStatus) error {
	tag, err := db.Exec(ctx, `
		UPDATE contractor_memberships
		SET status = $1, decided_at = CURRENT_TIMESTAMP, decided_by = $2
		WHERE id = $3 AND tenant_id = $4
	`, status, deciderID, membershipID, tenantID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("членство не знайдено у вашій організації")
	}
	return nil
}

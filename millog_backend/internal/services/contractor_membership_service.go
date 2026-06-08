package services

import (
	"context"

	"Omnilog_backend/internal/models"
	"Omnilog_backend/internal/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ContractorMembershipService інкапсулює бізнес-логіку схвалення підрядників організаціями.
type ContractorMembershipService struct {
	repo   *repositories.ContractorMembershipRepository
	dbPool *pgxpool.Pool
}

func NewContractorMembershipService(repo *repositories.ContractorMembershipRepository, db *pgxpool.Pool) *ContractorMembershipService {
	return &ContractorMembershipService{repo: repo, dbPool: db}
}

// withRLSBypass виконує fn у транзакції, де app.tenant_id порожній, тимчасово знімаючи
// tenant-isolation RLS на час читання. Це потрібно, бо членства зв'язують ГЛОБАЛЬНОГО
// підрядника (tenant_id IS NULL) з конкретною організацією, і JOIN до users/tenants інакше
// відфільтровується політикою. Безпека зберігається: самі запити явно скоупляться по
// contractor_memberships (m.tenant_id / m.contractor_id), тож крос-tenant витоку немає.
func (s *ContractorMembershipService) withRLSBypass(ctx context.Context, fn func(tx repositories.DBExecutor) error) error {
	tx, err := s.dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // read-only: завжди відкочуємо
	// set_config(..., is_local=true) діє лише в межах цієї транзакції.
	if _, err := tx.Exec(ctx, "SELECT set_config('app.tenant_id', '', true)"); err != nil {
		return err
	}
	return fn(tx)
}

// ListByTenant — для адмін-панелі організації: хто подався/схвалений/відхилений.
func (s *ContractorMembershipService) ListByTenant(ctx context.Context, tenantID string, status models.ContractorMembershipStatus) ([]models.ContractorMembership, error) {
	var out []models.ContractorMembership
	err := s.withRLSBypass(ctx, func(tx repositories.DBExecutor) error {
		list, err := s.repo.ListByTenant(ctx, tx, tenantID, status)
		if err != nil {
			return err
		}
		out = list
		return nil
	})
	return out, err
}

// ListByContractor — self-view підрядника: з якими організаціями він співпрацює.
func (s *ContractorMembershipService) ListByContractor(ctx context.Context, contractorID string) ([]models.ContractorMembership, error) {
	var out []models.ContractorMembership
	err := s.withRLSBypass(ctx, func(tx repositories.DBExecutor) error {
		list, err := s.repo.ListByContractor(ctx, tx, contractorID)
		if err != nil {
			return err
		}
		out = list
		return nil
	})
	return out, err
}

// Approve — організація підтверджує підрядника.
func (s *ContractorMembershipService) Approve(ctx context.Context, membershipID, tenantID, deciderID string) error {
	return s.repo.Decide(ctx, s.dbPool, membershipID, tenantID, deciderID, models.MembershipApproved)
}

// Reject — організація відхиляє підрядника.
func (s *ContractorMembershipService) Reject(ctx context.Context, membershipID, tenantID, deciderID string) error {
	return s.repo.Decide(ctx, s.dbPool, membershipID, tenantID, deciderID, models.MembershipRejected)
}

// Apply — підрядник самостійно надсилає заявку на співпрацю з організацією (без взяття завдання).
// Ідемпотентна: якщо запис уже існує (PENDING/APPROVED/REJECTED) — статус не скидається.
// Повертає поточний статус членства після виклику.
func (s *ContractorMembershipService) Apply(ctx context.Context, contractorID, tenantID string) (models.ContractorMembershipStatus, error) {
	return s.repo.Apply(ctx, s.dbPool, contractorID, tenantID)
}

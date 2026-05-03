package services

import (
	"Omnilog_backend/internal/models"
	"Omnilog_backend/internal/repositories"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditService struct {
	repo *repositories.AuditLogRepository
	Pool *pgxpool.Pool
}

func NewAuditService(repo *repositories.AuditLogRepository, pool *pgxpool.Pool) *AuditService {
	return &AuditService{repo: repo, Pool: pool}
}

func (s *AuditService) LogAction(ctx context.Context, userID, actionType, entityType, entityID, details string) error {
	return s.repo.LogAction(ctx, s.Pool, userID, actionType, entityType, entityID, details)
}

func (s *AuditService) GetLogs(ctx context.Context, limit int) ([]models.AuditLog, error) {
	return s.repo.GetLogs(ctx, s.Pool, limit)
}

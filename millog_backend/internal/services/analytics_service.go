package services

import (
	"context"

	"millog_backend/internal/models"
	"millog_backend/internal/repositories"
)

type AnalyticsService struct {
	repo *repositories.AnalyticsRepository
	db   repositories.DBExecutor
}

func NewAnalyticsService(repo *repositories.AnalyticsRepository, db repositories.DBExecutor) *AnalyticsService {
	return &AnalyticsService{repo: repo, db: db}
}

func (s *AnalyticsService) GetDashboardAnalytics(ctx context.Context, start, end, unitID string) (*models.DashboardAnalytics, error) {
	return s.repo.GetDashboardStats(ctx, s.db, start, end, unitID)
}

// НОВА ФУНКЦІЯ: Приймає налаштування замовлення з фронтенду та ID користувача
func (s *AnalyticsService) RunSmartReplenish(ctx context.Context, req models.SmartReplenishRequest, userID string) (int, error) {
	return s.repo.ProcessSmartReplenish(ctx, s.db, req, userID)
}

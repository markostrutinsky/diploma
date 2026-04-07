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

func (s *AnalyticsService) GetDashboardAnalytics(ctx context.Context, start, end string) (*models.DashboardAnalytics, error) {
	return s.repo.GetDashboardStats(ctx, s.db, start, end)
}

func (s *AnalyticsService) RunAutoReplenish(ctx context.Context) (int, error) {
	return s.repo.CreateAutoReplenishRequests(ctx, s.db)
}

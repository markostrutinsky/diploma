package services

import (
	"context"
	"fmt"
	"log"

	"Omnilog_backend/internal/models"
	"Omnilog_backend/internal/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationService struct {
	repo   *repositories.NotificationRepository
	dbPool *pgxpool.Pool
}

func NewNotificationService(repo *repositories.NotificationRepository, db *pgxpool.Pool) *NotificationService {
	return &NotificationService{
		repo:   repo,
		dbPool: db,
	}
}

// CreateNotification створює нове сповіщення для користувача
func (s *NotificationService) CreateNotification(ctx context.Context, req *models.CreateNotificationRequest) error {
	notif := &models.Notification{
		UserID:    req.UserID,
		Type:      req.Type,
		Title:     req.Title,
		Message:   req.Message,
		RelatedID: req.RelatedID,
		IsRead:    false,
	}

	return s.repo.Create(ctx, s.dbPool, notif)
}

// NotifyDriverAboutShipment створює сповіщення для водія про призначення рейсу
func (s *NotificationService) NotifyDriverAboutShipment(ctx context.Context, driverID string, shipmentID string, fromWarehouse string, toWarehouse string) error {
	title := "🚚 Новий рейс призначено"
	message := fmt.Sprintf("Вам призначено рейс: %s → %s. Перевірте деталі в розділі «Транспорт».", fromWarehouse, toWarehouse)

	log.Printf("DEBUG: NotifyDriverAboutShipment - driverID=%s, shipmentID=%s, title=%s", driverID, shipmentID, title)

	req := &models.CreateNotificationRequest{
		UserID:    driverID,
		Type:      models.NotificationShipmentAssigned,
		Title:     title,
		Message:   message,
		RelatedID: &shipmentID,
	}

	err := s.CreateNotification(ctx, req)
	log.Printf("DEBUG: CreateNotification result - err=%v", err)
	return err
}

// ListNotifications отримує список сповіщень користувача
func (s *NotificationService) ListNotifications(ctx context.Context, userID string, limit int) (*models.NotificationListResponse, error) {
	notifications, err := s.repo.ListByUser(ctx, s.dbPool, userID, limit)
	if err != nil {
		return nil, err
	}

	if notifications == nil {
		notifications = []models.Notification{}
	}

	unreadCount, err := s.repo.GetUnreadCount(ctx, s.dbPool, userID)
	if err != nil {
		return nil, err
	}

	return &models.NotificationListResponse{
		Notifications: notifications,
		UnreadCount:   unreadCount,
	}, nil
}

// MarkAsRead позначає сповіщення як прочитане
func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID string, userID string) error {
	return s.repo.MarkAsRead(ctx, s.dbPool, notificationID, userID)
}

// MarkAllAsRead позначає всі сповіщення користувача як прочитані
func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	return s.repo.MarkAllAsRead(ctx, s.dbPool, userID)
}

// DeleteNotification видаляє сповіщення
func (s *NotificationService) DeleteNotification(ctx context.Context, notificationID string, userID string) error {
	return s.repo.Delete(ctx, s.dbPool, notificationID, userID)
}

// GetUnreadCount повертає кількість непрочитаних сповіщень
func (s *NotificationService) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	return s.repo.GetUnreadCount(ctx, s.dbPool, userID)
}

package services

import (
	"context"
	"fmt"

	"millog_backend/internal/models"
	"millog_backend/internal/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CONTRACTORRequestService struct {
	repo   *repositories.CONTRACTORRequestRepository
	dbPool *pgxpool.Pool
}

func NewCONTRACTORRequestService(repo *repositories.CONTRACTORRequestRepository, db *pgxpool.Pool) *CONTRACTORRequestService {
	return &CONTRACTORRequestService{repo: repo, dbPool: db}
}

func (s *CONTRACTORRequestService) Create(ctx context.Context, userID string, unitID int64, req *models.CreateCONTRACTORRequest) (*models.CONTRACTORRequest, error) {
	var finalUnitID *int64
	if unitID != 0 {
		id := unitID
		finalUnitID = &id
	}

	vr := &models.CONTRACTORRequest{
		CreatedBy:   userID,
		UnitID:      finalUnitID,
		Title:       req.Title,
		Description: req.Description,
		Status:      models.CONTRACTOROpen,
	}

	if err := s.repo.Create(ctx, s.dbPool, vr); err != nil {
		return nil, fmt.Errorf("failed to create CONTRACTOR request: %w", err)
	}
	return vr, nil
}

// Отримання списку заявок (з фільтром по статусу)
func (s *CONTRACTORRequestService) List(ctx context.Context, status models.CONTRACTORRequestStatus) ([]models.CONTRACTORRequest, error) {
	return s.repo.List(ctx, s.dbPool, status)
}

// Універсальний метод для оновлення статусів (Взяти в роботу, Доставити, Скасувати, Відхилити)
func (s *CONTRACTORRequestService) UpdateStatus(ctx context.Context, requestID string, userID string, newStatus models.CONTRACTORRequestStatus) error {
	return s.repo.UpdateStatus(ctx, s.dbPool, requestID, userID, newStatus)
}

func (s *CONTRACTORRequestService) AcceptAndStore(ctx context.Context, requestID string, commanderID string, unitID int64, payload models.AcceptCONTRACTORPayload) error {

	// 1. Починаємо транзакцію
	tx, err := s.dbPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("не вдалося почати транзакцію: %w", err)
	}
	defer tx.Rollback(ctx)

	err = s.repo.UpdateStatus(ctx, tx, requestID, commanderID, models.CONTRACTORAccepted)
	if err != nil {
		return fmt.Errorf("не вдалося оновити статус заявки: %w", err)
	}

	if payload.ResourceID != nil && *payload.ResourceID != "" {
		updateQuery := `
            UPDATE resources 
            SET quantity = quantity + $1, updated_at = CURRENT_TIMESTAMP 
            WHERE id = $2 AND unit_id = $3
        `
		res, err := tx.Exec(ctx, updateQuery, payload.Quantity, *payload.ResourceID, unitID)
		if err != nil {
			return fmt.Errorf("помилка оновлення кількості майна: %w", err)
		}

		if res.RowsAffected() == 0 {
			return fmt.Errorf("ресурс не знайдено на балансі вашого підрозділу")
		}

	} else {
		insertQuery := `
            INSERT INTO resources (unit_id, category_id, name, quantity, unit_type, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        `
		_, err = tx.Exec(ctx, insertQuery, unitID, payload.CategoryID, payload.Name, payload.Quantity, payload.UnitType)
		if err != nil {
			return fmt.Errorf("помилка додавання нового майна на склад: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("помилка збереження транзакції: %w", err)
	}

	return nil
}

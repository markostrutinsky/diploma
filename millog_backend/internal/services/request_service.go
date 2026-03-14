package services

import (
	"context"
	"fmt"

	"millog_backend/internal/models"
	"millog_backend/internal/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RequestService struct {
	requestRepo  *repositories.SupplyRequestRepository
	resourceRepo *repositories.ResourceRepository
	dbPool       *pgxpool.Pool
}

func NewRequestService(reqRepo *repositories.SupplyRequestRepository, resRepo *repositories.ResourceRepository, db *pgxpool.Pool) *RequestService {
	return &RequestService{requestRepo: reqRepo, resourceRepo: resRepo, dbPool: db}
}

func (s *RequestService) Create(ctx context.Context, userID string, req *models.CreateSupplyRequest) (*models.SupplyRequest, error) {
	sr := &models.SupplyRequest{
		CreatedBy:  userID,
		ResourceID: req.ResourceID,
		Quantity:   req.Quantity,
		Status:     models.RequestPending,
	}
	if err := s.requestRepo.Create(ctx, s.dbPool, sr); err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	return sr, nil
}

func (s *RequestService) List(ctx context.Context) ([]models.SupplyRequest, error) {
	return s.requestRepo.List(ctx, s.dbPool)
}

func (s *RequestService) Approve(ctx context.Context, requestID, approverID string, approved bool, comment string) error {
	tx, err := s.dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	req, err := s.requestRepo.GetByID(ctx, tx, requestID)
	if err != nil {
		return fmt.Errorf("request not found")
	}
	if req.Status != models.RequestPending {
		return fmt.Errorf("request already processed")
	}

	if err := s.requestRepo.Approve(ctx, tx, requestID, approverID, approved, comment); err != nil {
		return err
	}

	if approved {
		resource, err := s.resourceRepo.GetByID(ctx, tx, req.ResourceID)
		if err != nil {
			return err
		}
		newQty := resource.Quantity + req.Quantity
		if err := s.resourceRepo.UpdateQuantity(ctx, tx, req.ResourceID, newQty); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

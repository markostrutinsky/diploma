package services

import (
	"context"
	"errors"
	"fmt"

	"millog_backend/internal/models"
	"millog_backend/internal/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CanApproveRequest(creatorRole models.UserRole, approverRole models.UserRole) bool {
	if approverRole == models.RoleAdmin {
		return true
	}
	allowedApprovers, exists := models.ApprovalMatrix[creatorRole]
	if !exists {
		return false
	}
	for _, role := range allowedApprovers {
		if role == approverRole {
			return true
		}
	}
	return false
}

type RequestService struct {
	requestRepo  *repositories.SupplyRequestRepository
	resourceRepo *repositories.ResourceRepository
	userRepo     *repositories.UserRepository
	dbPool       *pgxpool.Pool
}

func NewRequestService(reqRepo *repositories.SupplyRequestRepository, resRepo *repositories.ResourceRepository, userRepo *repositories.UserRepository, db *pgxpool.Pool) *RequestService {
	return &RequestService{requestRepo: reqRepo, resourceRepo: resRepo, userRepo: userRepo, dbPool: db}
}

func (s *RequestService) Create(ctx context.Context, userID string, req *models.CreateSupplyRequest) (*models.SupplyRequest, error) {
	sr := &models.SupplyRequest{
		CreatedBy:         userID,
		ResourceID:        req.ResourceID,
		Quantity:          req.Quantity,
		Status:            models.RequestPending,
		TargetWarehouseID: req.TargetWarehouseID,
	}

	if err := s.requestRepo.Create(ctx, s.dbPool, sr); err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return sr, nil
}

// Додаємо аргументи userRole та userUnitID
func (s *RequestService) List(ctx context.Context, userRole string, userUnitID *int64) ([]models.SupplyRequest, error) {
	return s.requestRepo.List(ctx, s.dbPool, userRole, userUnitID)
}

func (s *RequestService) Approve(ctx context.Context, requestID, approverID string, approverRole models.UserRole, approved bool, comment string) error {
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

	if req.CreatedBy == approverID {
		return errors.New("неможливо погодити власну заявку (конфлікт інтересів)")
	}

	creator, err := s.userRepo.GetByID(ctx, tx, req.CreatedBy)
	if err != nil {
		return fmt.Errorf("failed to get creator details: %w", err)
	}

	if !CanApproveRequest(models.UserRole(creator.Role), approverRole) {
		return errors.New("недостатньо прав для погодження заявки цього рівня (порушення субординації)")
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

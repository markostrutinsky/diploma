package services

import (
	"context"
	"fmt"

	"millog_backend/internal/models"
	"millog_backend/internal/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UnitService struct {
	repo     *repositories.UnitRepository
	userRepo *repositories.UserRepository
	dbPool   *pgxpool.Pool
}

func NewUnitService(repo *repositories.UnitRepository, userRepo *repositories.UserRepository, db *pgxpool.Pool) *UnitService {
	return &UnitService{repo: repo, userRepo: userRepo, dbPool: db}
}

func (s *UnitService) Create(ctx context.Context, req *models.CreateUnitRequest, creatorRole models.UserRole) (*models.Unit, error) {

	if !creatorRole.CanCreateUnitType(string(req.UnitType)) {
		return nil, fmt.Errorf("недостатньо прав для створення підрозділу типу %s", req.UnitType)
	}

	expectedParentType := ""
	switch req.UnitType {
	case "BATTALION":
		expectedParentType = "BRIGADE"
	case "COMPANY":
		expectedParentType = "BATTALION"
	case "PLATOON":
		expectedParentType = "COMPANY"
	case "BRIGADE":
		if req.ParentID != nil {
			return nil, fmt.Errorf("бригада не може бути підпорядкована іншому підрозділу")
		}
	default:
		return nil, fmt.Errorf("невідомий тип підрозділу: %s", req.UnitType)
	}

	if expectedParentType != "" {
		if req.ParentID == nil {
			return nil, fmt.Errorf("для підрозділу типу %s обов'язково потрібно вказати батьківський підрозділ", req.UnitType)
		}

		parentUnit, err := s.repo.GetByID(ctx, s.dbPool, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("батьківський підрозділ не знайдено: %w", err)
		}

		if parentUnit.UnitType != models.UnitType(expectedParentType) {
			return nil, fmt.Errorf("порушення ієрархії: підрозділ типу %s може підпорядковуватися тільки типу %s (обрано %s)",
				req.UnitType, expectedParentType, parentUnit.UnitType)
		}
	}

	u := &models.Unit{
		ParentID: req.ParentID,
		Name:     req.Name,
		UnitType: req.UnitType,
	}

	if err := s.repo.Create(ctx, s.dbPool, u); err != nil {
		return nil, fmt.Errorf("failed to create unit: %w", err)
	}

	return u, nil
}

func (s *UnitService) List(ctx context.Context) ([]models.Unit, error) {
	return s.repo.List(ctx, s.dbPool)
}

func (s *UnitService) GetAvailableForRole(ctx context.Context, roleStr string) ([]models.Unit, error) {
	role := models.UserRole(roleStr)
	unitType := role.GetTargetUnitType()

	if unitType == "" {
		// Якщо для цієї ролі підрозділ не потрібен (наприклад, Адмін), повертаємо пустий список
		return []models.Unit{}, nil
	}

	return s.repo.GetAvailableUnitsForRole(ctx, s.dbPool, unitType, role)
}

func (s *UnitService) GetVisibleUnits(ctx context.Context, userID string, role models.UserRole) ([]models.Unit, error) {
	if role == models.RoleAdmin {
		return s.repo.List(ctx, s.dbPool)
	}

	user, err := s.userRepo.GetByID(ctx, s.dbPool, userID)
	if err != nil {
		return nil, fmt.Errorf("не вдалося ідентифікувати користувача: %w", err)
	}

	if user.UnitID == nil {
		return []models.Unit{}, nil
	}

	return s.repo.GetUnitsHierarchy(ctx, s.dbPool, *user.UnitID)
}

func (s *UnitService) ChangeCommander(ctx context.Context, targetUnitID int64, newCommanderID string, requesterRole string, requesterUnitID int64) error {
	if requesterRole != "ADMIN" {
		hasAccess, err := s.repo.CheckHierarchy(ctx, s.dbPool, requesterUnitID, targetUnitID)
		if err != nil {
			return fmt.Errorf("помилка перевірки прав доступу: %v", err)
		}
		if !hasAccess {
			return fmt.Errorf("відмовлено в доступі: ви не можете змінити командира цього підрозділу")
		}
	}

	return s.repo.ChangeCommanderTx(ctx, s.dbPool, targetUnitID, newCommanderID)
}

func (s *UnitService) GetMyHierarchyForRole(ctx context.Context, roleStr string, commanderUnitID int64) ([]models.Unit, error) {
	role := models.UserRole(roleStr)
	unitType := role.GetTargetUnitType()

	if unitType == "" {
		// Якщо посаді не потрібен підрозділ (хоча таких тут не має бути)
		return []models.Unit{}, nil
	}

	return s.repo.GetAvailableUnitsInHierarchy(ctx, s.dbPool, unitType, role, commanderUnitID)
}

package services

import (
	"context"
	"fmt"

	"Omnilog_backend/internal/models"
	"Omnilog_backend/internal/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ContractorRequestService struct {
	repo       *repositories.ContractorRequestRepository
	memberRepo *repositories.ContractorMembershipRepository
	dbPool     *pgxpool.Pool
}

func NewContractorRequestService(repo *repositories.ContractorRequestRepository, memberRepo *repositories.ContractorMembershipRepository, db *pgxpool.Pool) *ContractorRequestService {
	return &ContractorRequestService{repo: repo, memberRepo: memberRepo, dbPool: db}
}

// MembershipRequiredError повідомляє, що підрядник не має схвалення від організації,
// якій належить завдання. Handler перетворює це на 403 із машинно-зчитуваним кодом,
// щоб фронтенд показав коректний стан («заявку надіслано» / «очікує» / «відхилено»).
type MembershipRequiredError struct {
	TenantID string
	Status   models.ContractorMembershipStatus
	Message  string
}

func (e *MembershipRequiredError) Error() string { return e.Message }

func (s *ContractorRequestService) Create(ctx context.Context, userID string, unitID int64, req *models.CreateContractorRequest) (*models.ContractorRequest, error) {
	var finalUnitID *int64
	if unitID != 0 {
		id := unitID
		finalUnitID = &id
	}

	vr := &models.ContractorRequest{
		CreatedBy:         userID,
		UnitID:            finalUnitID,
		TargetWarehouseID: req.TargetWarehouseID,
		Title:             req.Title,
		Description:       req.Description,
		Status:            models.ContractorOpen,
	}

	if err := s.repo.Create(ctx, s.dbPool, vr); err != nil {
		return nil, fmt.Errorf("failed to create Contractor request: %w", err)
	}
	return vr, nil
}

// Отримання списку заявок (з фільтром по статусу).
// isContractor=true → крос-tenant marketplace для підрядника; інакше — tenant-scoped.
func (s *ContractorRequestService) List(ctx context.Context, status models.ContractorRequestStatus, isContractor bool, contractorID string) ([]models.ContractorRequest, error) {
	return s.repo.List(ctx, s.dbPool, status, isContractor, contractorID)
}

// Універсальний метод для оновлення статусів (Взяти в роботу, Доставити, Скасувати, Відхилити)
func (s *ContractorRequestService) UpdateStatus(ctx context.Context, requestID string, userID string, newStatus models.ContractorRequestStatus) error {
	return s.repo.UpdateStatus(ctx, s.dbPool, requestID, userID, newStatus)
}

// Take — підрядник бере завдання в роботу, але лише якщо організація-замовник його схвалила.
// Якщо схвалення немає — автоматично подаємо заявку на співпрацю (PENDING) і повертаємо
// MembershipRequiredError, щоб користувач побачив зрозуміле пояснення.
func (s *ContractorRequestService) Take(ctx context.Context, requestID string, contractorID string) error {
	tenantID, err := s.repo.GetTenantID(ctx, s.dbPool, requestID)
	if err != nil {
		return fmt.Errorf("не вдалося знайти завдання: %w", err)
	}
	if tenantID == "" {
		return fmt.Errorf("не вдалося визначити організацію цього завдання")
	}

	approved, err := s.memberRepo.IsApproved(ctx, s.dbPool, contractorID, tenantID)
	if err != nil {
		return fmt.Errorf("не вдалося перевірити доступ до організації: %w", err)
	}

	if !approved {
		// Перша спроба = автоматична заявка на співпрацю.
		status, applyErr := s.memberRepo.Apply(ctx, s.dbPool, contractorID, tenantID)
		if applyErr != nil {
			return fmt.Errorf("не вдалося надіслати заявку на співпрацю: %w", applyErr)
		}

		msg := "Щоб брати завдання цієї організації, потрібне підтвердження. Заявку на співпрацю надіслано — очікуйте рішення адміністратора організації."
		if status == models.MembershipPending {
			msg = "Вашу заявку на співпрацю з організацією ще розглядають. Очікуйте підтвердження адміністратора."
		} else if status == models.MembershipRejected {
			msg = "Організація відхилила вашу заявку на співпрацю, тож ви не можете брати її завдання."
		}

		return &MembershipRequiredError{TenantID: tenantID, Status: status, Message: msg}
	}

	return s.repo.UpdateStatus(ctx, s.dbPool, requestID, contractorID, models.ContractorTaken)
}

func (s *ContractorRequestService) AcceptAndStore(ctx context.Context, requestID string, commanderID string, unitID int64, payload models.AcceptContractorPayload) error {

	// 1. Починаємо транзакцію
	tx, err := s.dbPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("не вдалося почати транзакцію: %w", err)
	}
	defer tx.Rollback(ctx)

	err = s.repo.UpdateStatus(ctx, tx, requestID, commanderID, models.ContractorAccepted)
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
		tid := repositories.TenantFromCtx(ctx)
		if tid == "" {
			return fmt.Errorf("tenant_id is required for creating resources")
		}

		insertQuery := `
            INSERT INTO resources (unit_id, category_id, name, quantity, unit_type, tenant_id, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        `
		_, err = tx.Exec(ctx, insertQuery, unitID, payload.CategoryID, payload.Name, payload.Quantity, payload.UnitType, tid)
		if err != nil {
			return fmt.Errorf("помилка додавання нового майна на склад: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("помилка збереження транзакції: %w", err)
	}

	return nil
}

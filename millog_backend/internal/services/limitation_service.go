package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SubscriptionLimits визначає обмеження для кожного тарифу
type SubscriptionLimits struct {
	MaxWarehouses int
	MaxResources  int
	MaxUsers      int
	MaxVehicles   int
	Unlimited     bool
}

var LimitsByTier = map[string]SubscriptionLimits{
	"FREE": {
		MaxWarehouses: 1,
		MaxResources:  20,
		MaxUsers:      5,
		MaxVehicles:   1,
		Unlimited:     false,
	},
	"BASIC": {
		MaxWarehouses: 10,
		MaxResources:  100,
		MaxUsers:      50,
		MaxVehicles:   5,
		Unlimited:     false,
	},
	"PRO": {
		MaxWarehouses: 100,
		MaxResources:  1000,
		MaxUsers:      500,
		MaxVehicles:   50,
		Unlimited:     false,
	},
	"ENTERPRISE": {
		MaxWarehouses: 999999,
		MaxResources:  999999,
		MaxUsers:      999999,
		MaxVehicles:   999999,
		Unlimited:     true,
	},
}

// LimitationService перевіряє обмеження підписки при створенні об'єктів
type LimitationService struct {
	dbPool *pgxpool.Pool
}

func NewLimitationService(db *pgxpool.Pool) *LimitationService {
	return &LimitationService{dbPool: db}
}

// CheckWarehouseLimit перевіряє, чи може користувач створити ще один склад
func (s *LimitationService) CheckWarehouseLimit(ctx context.Context, unitID int64) error {
	tier, err := s.getUserTier(ctx, unitID)
	if err != nil {
		return err
	}

	if LimitsByTier[tier].Unlimited {
		return nil
	}

	limit := LimitsByTier[tier].MaxWarehouses

	// Лічимо кількість складів для цієї одиниці
	query := `SELECT COUNT(*) FROM warehouses WHERE unit_id = $1`
	var count int
	err = s.dbPool.QueryRow(ctx, query, unitID).Scan(&count)
	if err != nil {
		return err
	}

	if count >= limit {
		return fmt.Errorf(
			"Ліміт складів досягнут: %d/%d. Обновіте підписку на PRO для більшого ліміту.",
			count, limit,
		)
	}

	return nil
}

// CheckResourceLimit перевіряє, чи може користувач створити ще один ресурс
func (s *LimitationService) CheckResourceLimit(ctx context.Context, unitID int64) error {
	tier, err := s.getUserTier(ctx, unitID)
	if err != nil {
		return err
	}

	if LimitsByTier[tier].Unlimited {
		return nil
	}

	limit := LimitsByTier[tier].MaxResources

	// Лічимо ресурси всіх складів цієї одиниці
	query := `
		SELECT COUNT(DISTINCT r.id) 
		FROM resources r
		JOIN warehouses w ON r.warehouse_id = w.id
		WHERE w.unit_id = $1
	`
	var count int
	err = s.dbPool.QueryRow(ctx, query, unitID).Scan(&count)
	if err != nil {
		return err
	}

	if count >= limit {
		return fmt.Errorf(
			"Ліміт ресурсів досягнут: %d/%d. Обновіте підписку на PRO для більшого ліміту.",
			count, limit,
		)
	}

	return nil
}

// CheckUserLimit перевіряє, чи може користувач додати ще одного користувача
func (s *LimitationService) CheckUserLimit(ctx context.Context, creatorRole string, creatorUnitID int64) error {
	// Отримуємо тариф за ролеюю користувача який створює
	tier, err := s.getUserTier(ctx, creatorUnitID)
	if err != nil {
		return err
	}

	if LimitsByTier[tier].Unlimited {
		return nil
	}

	limit := LimitsByTier[tier].MaxUsers

	// Лічимо активних користувачів для цієї одиниці та її нижчих рівнів
	query := `
		WITH RECURSIVE unit_tree AS (
			SELECT id FROM units WHERE id = $1
			UNION ALL
			SELECT u.id FROM units u
			JOIN unit_tree ut ON u.parent_id = ut.id
		)
		SELECT COUNT(*) FROM users u
		WHERE u.unit_id IN (SELECT id FROM unit_tree)
		AND u.status != 'BLOCKED'
	`
	var count int
	err = s.dbPool.QueryRow(ctx, query, creatorUnitID).Scan(&count)
	if err != nil {
		return err
	}

	if count >= limit {
		return fmt.Errorf(
			"Ліміт користувачів досягнут: %d/%d. Обновіте підписку на PRO для більшого ліміту.",
			count, limit,
		)
	}

	return nil
}

// CheckVehicleLimit перевіряє, чи може користувач додати ще один транспортний засіб
func (s *LimitationService) CheckVehicleLimit(ctx context.Context, unitID int64) error {
	tier, err := s.getUserTier(ctx, unitID)
	if err != nil {
		return err
	}

	if LimitsByTier[tier].Unlimited {
		return nil
	}

	limit := LimitsByTier[tier].MaxVehicles

	query := `SELECT COUNT(*) FROM vehicles WHERE unit_id = $1`
	var count int
	err = s.dbPool.QueryRow(ctx, query, unitID).Scan(&count)
	if err != nil {
		return err
	}

	if count >= limit {
		return fmt.Errorf(
			"Ліміт автомобілів досягнут: %d/%d. Обновіте підписку на PRO для більшого ліміту.",
			count, limit,
		)
	}

	return nil
}

// getUserTier отримує тариф користувача за unitID
func (s *LimitationService) getUserTier(ctx context.Context, unitID int64) (string, error) {
	if unitID == 0 {
		return "BASIC", nil
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		WITH RECURSIVE unit_hierarchy AS (
			SELECT id, parent_id, subscription_tier, 1 as depth
			FROM units
			WHERE id = $1
			UNION ALL
			SELECT u.id, u.parent_id, u.subscription_tier, uh.depth + 1
			FROM units u
			JOIN unit_hierarchy uh ON u.parent_id = uh.id
		)
		SELECT subscription_tier 
		FROM unit_hierarchy
		ORDER BY (CASE 
			WHEN subscription_tier = 'ENTERPRISE' THEN 2
			WHEN subscription_tier = 'PRO' THEN 1
			ELSE 0 
		END) DESC
		LIMIT 1
	`

	var tier string
	err := s.dbPool.QueryRow(ctx, query, unitID).Scan(&tier)
	if err != nil {
		return "BASIC", nil
	}

	if tier == "" {
		return "BASIC", nil
	}

	return tier, nil
}

// GetLimits повертає поточні ліміти для користувача
func (s *LimitationService) GetLimits(ctx context.Context, unitID int64) (map[string]interface{}, error) {
	tier, err := s.getUserTier(ctx, unitID)
	if err != nil {
		tier = "BASIC"
	}

	limits := LimitsByTier[tier]

	// Лічимо поточні використання
	warehouseCount := 0
	resourceCount := 0
	userCount := 0
	vehicleCount := 0

	if unitID != 0 {
		// Warehouses
		_ = s.dbPool.QueryRow(ctx, "SELECT COUNT(*) FROM warehouses WHERE unit_id = $1", unitID).Scan(&warehouseCount)

		// Resources
		_ = s.dbPool.QueryRow(ctx, `
			SELECT COUNT(DISTINCT r.id) 
			FROM resources r
			JOIN warehouses w ON r.warehouse_id = w.id
			WHERE w.unit_id = $1
		`, unitID).Scan(&resourceCount)

		// Users (this unit + children)
		_ = s.dbPool.QueryRow(ctx, `
			WITH RECURSIVE unit_tree AS (
				SELECT id FROM units WHERE id = $1
				UNION ALL
				SELECT u.id FROM units u
				JOIN unit_tree ut ON u.parent_id = ut.id
			)
			SELECT COUNT(*) FROM users u
			WHERE u.unit_id IN (SELECT id FROM unit_tree)
			AND u.status != 'BLOCKED'
		`, unitID).Scan(&userCount)

		// Vehicles
		_ = s.dbPool.QueryRow(ctx, "SELECT COUNT(*) FROM vehicles WHERE unit_id = $1", unitID).Scan(&vehicleCount)
	}

	return map[string]interface{}{
		"subscription_tier": tier,
		"warehouses": map[string]interface{}{
			"current": warehouseCount,
			"limit":   limits.MaxWarehouses,
			"percent": (warehouseCount * 100) / limits.MaxWarehouses,
		},
		"resources": map[string]interface{}{
			"current": resourceCount,
			"limit":   limits.MaxResources,
			"percent": (resourceCount * 100) / limits.MaxResources,
		},
		"users": map[string]interface{}{
			"current": userCount,
			"limit":   limits.MaxUsers,
			"percent": (userCount * 100) / limits.MaxUsers,
		},
		"vehicles": map[string]interface{}{
			"current": vehicleCount,
			"limit":   limits.MaxVehicles,
			"percent": (vehicleCount * 100) / limits.MaxVehicles,
		},
	}, nil
}

package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"millog_backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SubscriptionTierWeight визначає іерархію тарифів
var SubscriptionTierWeight = map[string]int{
	"FREE":       0,
	"BASIC":      1,
	"PRO":        2,
	"ENTERPRISE": 3,
}

// RequireSubscriptionTier перевіряє, чи користувач має достатній тариф
// Приклад використання: middleware.RequireSubscriptionTier("PRO", dbPool)
func RequireSubscriptionTier(minTier string, dbPool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsVal, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Не знайдено токен авторизації"})
			c.Abort()
			return
		}

		claims := claimsVal.(*Claims)

		// SYSTEM_ADMIN завжди має доступ до всіх фіч (підтримка платформи).
		// Застаріле ADMIN також трактуємо як tenant-адміна — для нього тариф перевіряється нижче.
		if claims.Role == models.RoleSystemAdmin {
			c.Next()
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		userTier, err := getTenantSubscriptionTier(ctx, dbPool, claims.TenantID, claims.UnitID)
		if err != nil {
			slog.Error("Failed to fetch subscription tier", "error", err, "tenantID", claims.TenantID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка отримання тарифу"})
			c.Abort()
			return
		}

		// Перевіряємо чи достатньо тарифу
		requiredWeight := SubscriptionTierWeight[minTier]
		userWeight := SubscriptionTierWeight[userTier]

		if userWeight < requiredWeight {
			// Логуємо спробу несанкціонованого доступу
			logUnauthorizedPremiumAccess(c, dbPool, claims.UserID, minTier, userTier)

			c.JSON(http.StatusPaymentRequired, gin.H{
				"error":         "Функція доступна тільки на тарифі: " + minTier,
				"current_tier":  userTier,
				"required_tier": minTier,
				"upgrade_url":   "/billing",
				"message":       "Оновіть підписку для доступу до цієї функції",
			})
			c.Abort()
			return
		}

		// Передаємо тариф далі в handler для логування
		c.Set("user_tier", userTier)
		c.Next()
	}
}

// getTenantSubscriptionTier бере тариф з tenants.subscription_tier.
// Fallback на units.subscription_tier — сумісність зі старими даними, де tenant ще не виставлений.
func getTenantSubscriptionTier(ctx context.Context, dbPool *pgxpool.Pool, tenantID string, unitID int64) (string, error) {
	// 1. Основне джерело істини — tenants.
	if tenantID != "" {
		var tier string
		err := dbPool.QueryRow(ctx, `SELECT subscription_tier FROM tenants WHERE id = $1`, tenantID).Scan(&tier)
		if err == nil && tier != "" {
			return tier, nil
		}
	}

	// 2. Fallback для legacy даних — рекурсивно по units
	if unitID == 0 {
		return "FREE", nil
	}
	query := `
		WITH RECURSIVE unit_hierarchy AS (
			SELECT id, parent_id, subscription_tier, 1 as depth
			FROM units WHERE id = $1
			UNION ALL
			SELECT u.id, u.parent_id, u.subscription_tier, uh.depth + 1
			FROM units u JOIN unit_hierarchy uh ON u.parent_id = uh.id
		)
		SELECT subscription_tier FROM unit_hierarchy
		ORDER BY (CASE
			WHEN subscription_tier = 'ENTERPRISE' THEN 3
			WHEN subscription_tier = 'PRO' THEN 2
			WHEN subscription_tier = 'BASIC' THEN 1
			ELSE 0 END) DESC
		LIMIT 1`

	var tier string
	err := dbPool.QueryRow(ctx, query, unitID).Scan(&tier)
	if err != nil || tier == "" {
		return "FREE", nil
	}
	return tier, nil
}

// logUnauthorizedPremiumAccess логує спробу доступу без дозволу до платної фічи
func logUnauthorizedPremiumAccess(c *gin.Context, dbPool *pgxpool.Pool, userID string, requiredTier, currentTier string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		query := `
			INSERT INTO audit_logs (user_id, action, entity_type, entity_id, details, created_at)
			VALUES ($1, $2, $3, $4, $5, NOW())
		`

		details := fmt.Sprintf(
			"Спроба доступу до функції %s. Потрібен %s, має %s. Метод: %s, Путь: %s",
			requiredTier,
			requiredTier,
			currentTier,
			c.Request.Method,
			c.Request.URL.Path,
		)

		_, _ = dbPool.Exec(ctx, query,
			userID,
			"UNAUTHORIZED_PREMIUM_ACCESS",
			"SECURITY",
			c.Request.URL.Path,
			details,
		)
	}()
}

package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"Omnilog_backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Claims struct {
	UserID   string          `json:"user_id"`
	TenantID string          `json:"tenant_id"`
	Email    string          `json:"email"`
	Role     models.UserRole `json:"role"`
	UnitID   int64           `json:"unit_id"`
	jwt.RegisteredClaims
}

func AuthMiddleware(jwtSecret string, db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string

		// Пріоритет 1: Authorization header (Bearer token)
		auth := c.GetHeader("Authorization")
		if auth != "" {
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenStr = parts[1]
			}
		}

		// Пріоритет 2: access_token cookie (httpOnly)
		if tokenStr == "" {
			if cookieVal, err := c.Cookie("access_token"); err == nil && cookieVal != "" {
				tokenStr = cookieVal
			}
		}

		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}
		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		// === НОВИЙ БЛОК: ПЕРЕВІРКА В БАЗІ ДАНИХ ===
		var status string
		var tenantID *string
		// 🛡️ Додаємо таймаут до контексту (2 сек) щоб не висіти на БД
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		err = db.QueryRow(ctx, "SELECT status, tenant_id FROM users WHERE id = $1", claims.UserID).Scan(&status, &tenantID)

		// Якщо користувача видалили з БД або його статус BLOCKED - відхиляємо запит
		if err != nil || status == "BLOCKED" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Ваш профіль деактивовано. Сесію завершено."})
			c.Abort()
			return
		}
		// ==========================================

		// Якщо tenant_id у БД є, а у claims пустий/інший — беремо з БД (свіжіший source of truth)
		if tenantID != nil && *tenantID != "" {
			claims.TenantID = *tenantID
		}

		// Якщо все добре і статус ACTIVE — пускаємо далі
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("unit_id", claims.UnitID)
		c.Set("tenant_id", claims.TenantID)
		c.Set("claims", claims)

		// Пробрасуємо tenant_id у context.Context, щоб repositories могли
		// автоматично застосовувати tenant scoping без зміни сигнатур.
		if claims.TenantID != "" {
			c.Request = c.Request.WithContext(WithTenant(c.Request.Context(), claims.TenantID))
		}

		c.Next()
	}
}

func RequireRoles(roles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		role := roleVal.(models.UserRole)
		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		c.Abort()
	}
}

// RequireAnyRole checks if user has any of the given roles
func RequireAnyRole(roles []models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		role := roleVal.(models.UserRole)
		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		c.Abort()
	}
}

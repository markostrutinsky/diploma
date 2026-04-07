package middleware

import (
	"net/http"
	"strings"

	"millog_backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Claims struct {
	UserID string          `json:"user_id"`
	Email  string          `json:"email"`
	Role   models.UserRole `json:"role"`
	UnitID int64           `json:"unit_id"`
	jwt.RegisteredClaims
}

func AuthMiddleware(jwtSecret string, db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		tokenStr := parts[1]
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
		// Робимо швидкий запит до БД, щоб перевірити актуальний статус користувача
		err = db.QueryRow(c.Request.Context(), "SELECT status FROM users WHERE id = $1", claims.UserID).Scan(&status)

		// Якщо користувача видалили з БД або його статус BLOCKED - відхиляємо запит
		if err != nil || status == "BLOCKED" {
			// Можеш написати текст помилки українською, щоб фронтенд його красиво показав
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Ваш профіль деактивовано. Сесію завершено."})
			c.Abort()
			return
		}
		// ==========================================

		// Якщо все добре і статус ACTIVE — пускаємо далі
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("unit_id", claims.UnitID)
		c.Set("claims", claims)

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

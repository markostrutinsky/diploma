package middleware

import (
	"net/http"

	"millog_backend/internal/models"

	"github.com/gin-gonic/gin"
)

// TenantIDFromContext повертає tenant_id поточного користувача, або "" якщо відсутній.
func TenantIDFromContext(c *gin.Context) string {
	v, ok := c.Get("tenant_id")
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}

// UserIDFromContext повертає user_id поточного користувача.
func UserIDFromContext(c *gin.Context) string {
	v, ok := c.Get("user_id")
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}

// UserRoleFromContext повертає роль користувача.
func UserRoleFromContext(c *gin.Context) models.UserRole {
	v, ok := c.Get("user_role")
	if !ok {
		return ""
	}
	r, _ := v.(models.UserRole)
	return r
}

// RequireSystemAdmin — endpoint доступний лише для платформного адміна.
// Використовується на /admin/platform/* (tenants management, billing, cross-tenant stats).
func RequireSystemAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if UserRoleFromContext(c) != models.RoleSystemAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Потрібні права SYSTEM_ADMIN"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireTenant — endpoint доступний лише користувачам із прив'язкою до tenant.
// Блокує SYSTEM_ADMIN без active tenant_id (бо у них немає даних для ізоляції).
func RequireTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		if TenantIDFromContext(c) == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Ця операція доступна лише в контексті організації"})
			c.Abort()
			return
		}
		c.Next()
	}
}

package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"millog_backend/internal/models"
	"millog_backend/internal/repositories"
)

// PlatformHandler обробляє ендпойнти платформного адміна (SYSTEM_ADMIN).
// Тут немає tenant-фільтрації — ми свідомо дивимось крос-тенантно.
type PlatformHandler struct {
	tenantRepo *repositories.TenantRepository
	userRepo   *repositories.UserRepository
	db         *pgxpool.Pool
}

func NewPlatformHandler(tr *repositories.TenantRepository, ur *repositories.UserRepository, db *pgxpool.Pool) *PlatformHandler {
	return &PlatformHandler{tenantRepo: tr, userRepo: ur, db: db}
}

// GET /api/platform/tenants?search=...
func (h *PlatformHandler) ListTenants(c *gin.Context) {
	search := c.Query("search")
	// SYSTEM_ADMIN — ctx без tenant-а → репо поверне всі
	tenants, err := h.tenantRepo.List(c.Request.Context(), h.db, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Додаємо кількість юзерів окремо (аби не ускладнювати SQL join)
	type tenantWithStats struct {
		models.Tenant
		UserCount int `json:"user_count"`
	}
	out := make([]tenantWithStats, 0, len(tenants))
	for _, t := range tenants {
		n, _ := h.tenantRepo.CountUsers(c.Request.Context(), h.db, t.ID)
		out = append(out, tenantWithStats{Tenant: t, UserCount: n})
	}
	c.JSON(http.StatusOK, out)
}

// GET /api/platform/tenants/:id
func (h *PlatformHandler) GetTenant(c *gin.Context) {
	id := c.Param("id")
	t, err := h.tenantRepo.GetByID(c.Request.Context(), h.db, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}
	n, _ := h.tenantRepo.CountUsers(c.Request.Context(), h.db, id)
	c.JSON(http.StatusOK, gin.H{"tenant": t, "user_count": n})
}

type updateTierRequest struct {
	Tier      string  `json:"tier" binding:"required,oneof=FREE BASIC PRO ENTERPRISE"`
	ExpiresAt *string `json:"expires_at,omitempty"` // ISO-8601 або null
}

// PATCH /api/platform/tenants/:id/tier
func (h *PlatformHandler) UpdateTier(c *gin.Context) {
	id := c.Param("id")
	var req updateTierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var exp *time.Time
	if req.ExpiresAt != nil && *req.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "expires_at must be RFC3339"})
			return
		}
		exp = &t
	}
	if err := h.tenantRepo.UpdateTier(c.Request.Context(), h.db, id,
		models.SubscriptionTier(req.Tier), exp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

type setActiveRequest struct {
	Active bool `json:"active"`
}

// PATCH /api/platform/tenants/:id/active
func (h *PlatformHandler) SetActive(c *gin.Context) {
	id := c.Param("id")
	var req setActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.tenantRepo.SetActive(c.Request.Context(), h.db, id, req.Active); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// DELETE /api/platform/tenants/:id — НЕЗВОРОТНО видаляє tenant з усіма даними.
func (h *PlatformHandler) DeleteTenant(c *gin.Context) {
	id := c.Param("id")
	if err := h.tenantRepo.Delete(c.Request.Context(), h.db, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// GET /api/platform/stats — overview для платформного дашборду.
func (h *PlatformHandler) Stats(c *gin.Context) {
	ctx := c.Request.Context()
	var stats struct {
		TotalTenants     int            `json:"total_tenants"`
		ActiveTenants    int            `json:"active_tenants"`
		TotalUsers       int            `json:"total_users"`
		TenantsByTier    map[string]int `json:"tenants_by_tier"`
		NewTenants30Days int            `json:"new_tenants_30_days"`
	}
	stats.TenantsByTier = map[string]int{}
	_ = h.db.QueryRow(ctx, `SELECT COUNT(*) FROM tenants`).Scan(&stats.TotalTenants)
	_ = h.db.QueryRow(ctx, `SELECT COUNT(*) FROM tenants WHERE is_active = TRUE`).Scan(&stats.ActiveTenants)
	_ = h.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&stats.TotalUsers)
	_ = h.db.QueryRow(ctx, `SELECT COUNT(*) FROM tenants WHERE created_at > NOW() - INTERVAL '30 days'`).Scan(&stats.NewTenants30Days)
	rows, err := h.db.Query(ctx, `SELECT subscription_tier, COUNT(*) FROM tenants GROUP BY subscription_tier`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tier string
			var cnt int
			if err := rows.Scan(&tier, &cnt); err == nil {
				stats.TenantsByTier[tier] = cnt
			}
		}
	}
	c.JSON(http.StatusOK, stats)
}

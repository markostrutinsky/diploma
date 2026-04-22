package handlers

import (
	"context"
	"fmt"
	"millog_backend/internal/middleware"
	"millog_backend/internal/models"
	"millog_backend/internal/services"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService  *services.AuthService
	auditService *services.AuditService
}

func NewAuthHandler(authService *services.AuthService, auditService *services.AuditService) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		auditService: auditService,
	}
}

// SignupTenant — публічний endpoint для self-service реєстрації організації.
// Створює tenant + першого TENANT_ADMIN.
func (h *AuthHandler) SignupTenant(c *gin.Context) {
	var req models.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, ownerID, err := h.authService.CreateTenant(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func(tID, oID, name string) {
		_ = h.auditService.LogAction(context.Background(), oID, "CREATE", "TENANT", tID, "Створено нову організацію: "+name)
	}(tenantID, ownerID, req.OrganizationName)

	c.JSON(http.StatusCreated, gin.H{
		"tenant_id": tenantID,
		"owner_id":  ownerID,
		"message":   "Організацію створено. Тепер ви можете увійти як адміністратор.",
	})
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	userID := c.GetString("user_id")

	var request models.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	creatorRoleVal, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "не вдалося визначити права користувача"})
		return
	}
	creatorRole := creatorRoleVal.(models.UserRole)
	creatorTenantID, _ := c.Get("tenant_id")
	creatorTenantIDStr, _ := creatorTenantID.(string)
	response, err := h.authService.RegisterUser(c.Request.Context(), &request, creatorRole, creatorTenantIDStr)
	if err != nil {
		if strings.Contains(err.Error(), "не має прав для створення") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func(uID string, entityID string, newUserName string) {
		_ = h.auditService.LogAction(context.Background(), uID, "CREATE", "USER", entityID, "Зареєстровано нового користувача: "+newUserName)
	}(userID, response.ID, request.FullName)

	c.JSON(http.StatusCreated, gin.H{
		"id":      response.ID,
		"message": "User registered successfully. Please wait for admin approval.",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request models.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.authService.Login(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	go func(email string) {
		_ = h.auditService.LogAction(context.Background(), "SYSTEM", "LOGIN", "USER", email, "Успішна авторизація користувача")
	}(request.Email)

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) ListCommanders(c *gin.Context) {
	commanders, err := h.authService.GetCommanders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не вдалося завантажити список керівників"})
		return
	}

	if commanders == nil {
		commanders = []*models.User{}
	}

	c.JSON(http.StatusOK, commanders)
}

func (h *AuthHandler) GetVisibleUsers(c *gin.Context) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не знайдено токен авторизації"})
		return
	}
	claims := claimsVal.(*middleware.Claims)

	var reqUnitID *int64
	if claims.UnitID != 0 {
		val := claims.UnitID
		reqUnitID = &val
	}

	users, err := h.authService.GetVisibleUsers(c.Request.Context(), string(claims.Role), reqUnitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка завантаження списку користувачів"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *AuthHandler) BootstrapAdmin(c *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		FullName string `json:"full_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.authService.BootstrapAdmin(c.Request.Context(), request.Email, request.Password, request.FullName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func(email string) {
		_ = h.auditService.LogAction(context.Background(), "SYSTEM", "CREATE", "USER", email, "Ініціалізація головного адміністратора")
	}(request.Email)

	c.JSON(http.StatusCreated, gin.H{"message": "Admin created. You can now log in."})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := h.authService.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) SetupPassword(c *gin.Context) {
	var request models.SetupPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.authService.SetupPassword(c.Request.Context(), request.Token, request.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func(token string) {
		_ = h.auditService.LogAction(context.Background(), "SYSTEM", "UPDATE", "USER", "", "Встановлення пароля за токеном: "+token)
	}(request.Token)

	c.JSON(http.StatusOK, gin.H{"message": "Password set successfully. You can now log in."})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var request models.RefreshRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.authService.Refresh(c.Request.Context(), request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	go func() {
		_ = h.auditService.LogAction(context.Background(), "SYSTEM", "REFRESH", "TOKEN", "", "Оновлення токена доступу")
	}()

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) RegisterContractor(c *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		FullName string `json:"full_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.authService.RegisterCONTRACTOR(c.Request.Context(), request.Email, request.Password, request.FullName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func(email string) {
		_ = h.auditService.LogAction(context.Background(), "SYSTEM", "CREATE", "USER", email, "Реєстрація зовнішнього підрядника")
	}(request.Email)

	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) UpdateRoleAndUnit(c *gin.Context) {
	targetUserID := c.Param("id")
	commanderID := c.GetString("user_id")

	var req models.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильний формат даних"})
		return
	}

	err := h.authService.UpdateRoleAndUnit(c.Request.Context(), commanderID, targetUserID, string(req.Role), req.UnitID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	go func(uID string, entityID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "USER", entityID, "Оновлено посаду/відділ користувача")
	}(commanderID, targetUserID)

	c.JSON(http.StatusOK, gin.H{"message": "Кадрове переміщення виконано успішно"})
}

func (h *AuthHandler) BlockUser(c *gin.Context) {
	userID := c.GetString("user_id")
	targetUserID := c.Param("id")

	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не знайдено токен"})
		return
	}
	claims := claimsVal.(*middleware.Claims)

	if fmt.Sprint(claims.UserID) == targetUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неможливо заблокувати власний профіль"})
		return
	}

	isAuthorized := false

	// 1. Визначаємо, чи є поточна роль керівною
	isManager := false
	switch models.UserRole(claims.Role) {
	case models.RoleRegionDirector,
		models.RoleBranchManager,
		models.RoleDeptManager,
		models.RoleDeptSupervisor,
		models.RoleTeamLead:
		isManager = true
	}

	// 2. Перевіряємо права на блокування
	if claims.Role == models.RoleAdmin {
		isAuthorized = true
	} else if isManager {
		if claims.UnitID != 0 {
			isSubordinate, err := h.authService.CheckSubordination(c.Request.Context(), claims.UnitID, targetUserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка перевірки підпорядкування"})
				return
			}
			if isSubordinate {
				isAuthorized = true
			}
		}
	}

	if !isAuthorized {
		c.JSON(http.StatusForbidden, gin.H{"error": "У вас немає прав для блокування цього користувача (він не є вашим підлеглим)"})
		return
	}

	err := h.authService.BlockUser(c.Request.Context(), targetUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не вдалося заблокувати користувача"})
		return
	}

	go func(uID string, entityID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "USER", entityID, "Заблоковано обліковий запис користувача")
	}(userID, targetUserID)

	c.JSON(http.StatusOK, gin.H{"message": "Користувача переведено в резерв (заблоковано)"})
}

func (h *AuthHandler) UnblockUser(c *gin.Context) {
	userID := c.GetString("user_id")
	targetUserID := c.Param("id")

	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не знайдено токен"})
		return
	}
	claims := claimsVal.(*middleware.Claims)

	if claims.Role != "ADMIN" && claims.Role != "REGION_DIRECTOR" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Тільки Адміністратор або Керівник регіону може розблоковувати людей"})
		return
	}

	err := h.authService.UnblockUser(c.Request.Context(), targetUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не вдалося розблокувати користувача"})
		return
	}

	go func(uID string, entityID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "USER", entityID, "Розблоковано обліковий запис користувача")
	}(userID, targetUserID)

	c.JSON(http.StatusOK, gin.H{"message": "Користувача розблоковано та переведено в кадровий резерв"})
}

func (h *AuthHandler) UpdateMyProfile(c *gin.Context) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не знайдено токен авторизації"})
		return
	}
	claims := claimsVal.(*middleware.Claims)
	userID := claims.UserID

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Невірний формат даних: " + err.Error()})
		return
	}

	err := h.authService.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		if strings.Contains(err.Error(), "зайнятий") || strings.Contains(err.Error(), "використовується") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func(uID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "USER", uID, "Оновлено власний профіль")
	}(userID)

	c.JSON(http.StatusOK, gin.H{"message": "Профіль успішно оновлено"})
}

func (h *AuthHandler) UpdateMyPassword(c *gin.Context) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не знайдено токен авторизації"})
		return
	}
	claims := claimsVal.(*middleware.Claims)

	var req models.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пароль має містити щонайменше 8 символів"})
		return
	}

	err := h.authService.UpdateMyPassword(c.Request.Context(), claims.UserID, req.OldPassword, req.NewPassword)
	if err != nil {
		if err.Error() == "невірний поточний пароль" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func(uID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "USER", uID, "Змінено власний пароль")
	}(claims.UserID)

	c.JSON(http.StatusOK, gin.H{"message": "Пароль успішно змінено"})
}

func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Введіть коректний email"})
		return
	}

	_ = h.authService.RequestPasswordReset(c.Request.Context(), req.Email)

	go func(email string) {
		_ = h.auditService.LogAction(context.Background(), "SYSTEM", "UPDATE", "USER", email, "Запит на скидання пароля")
	}(req.Email)

	c.JSON(http.StatusOK, gin.H{
		"message": "Якщо цей email зареєстрований у системі, ми надіслали на нього інструкції з відновлення пароля.",
	})
}

// GetUserLimits returns subscription tier information and resource usage limits for the current user
func (h *AuthHandler) GetUserLimits(c *gin.Context) {
	userID := c.GetString("user_id")
	unitIDVal, exists := c.Get("unit_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не вдалося визначити підрозділ користувача"})
		return
	}
	unitID := unitIDVal.(string)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Get subscription tier
	tier, err := h.authService.GetUserSubscriptionTier(ctx, unitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не вдалося отримати інформацію про підписку"})
		return
	}

	// Define limits by tier
	limitsByTier := map[string]map[string]interface{}{
		"BASIC": {
			"maxWarehouses": 10,
			"maxResources":  100,
			"maxUsers":      50,
			"maxVehicles":   5,
			"unlimited":     false,
		},
		"PRO": {
			"maxWarehouses": 100,
			"maxResources":  1000,
			"maxUsers":      500,
			"maxVehicles":   50,
			"unlimited":     false,
		},
		"ENTERPRISE": {
			"maxWarehouses": -1,
			"maxResources":  -1,
			"maxUsers":      -1,
			"maxVehicles":   -1,
			"unlimited":     true,
		},
	}

	limits := limitsByTier[tier]
	if limits == nil {
		limits = limitsByTier["BASIC"]
	}

	go func() {
		_ = h.auditService.LogAction(context.Background(), userID, "READ", "USER_LIMITS", userID, fmt.Sprintf("Перегляд інформації про ліміти (tier: %s)", tier))
	}()

	c.JSON(http.StatusOK, gin.H{
		"subscriptionTier": tier,
		"limits":           limits,
	})
}

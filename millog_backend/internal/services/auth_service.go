package services

import (
	"context"
	"errors"
	"fmt"
	"millog_backend/internal/middleware"
	"millog_backend/internal/models"
	"millog_backend/internal/repositories"
	"millog_backend/internal/tokens"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepository         *repositories.UserRepository
	unitRepository         *repositories.UnitRepository
	inviteTokenRepository  *repositories.InviteTokenRepository
	refreshTokenRepository *repositories.RefreshTokenRepository
	dbPool                 *pgxpool.Pool
	emailService           EmailService
	jwtSecret              string
}

func NewAuthService(
	userRepository *repositories.UserRepository,
	unitRepository *repositories.UnitRepository,
	inviteTokenRepository *repositories.InviteTokenRepository,
	refreshTokenRepository *repositories.RefreshTokenRepository,
	dbPool *pgxpool.Pool,
	emailService EmailService,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		userRepository:         userRepository,
		unitRepository:         unitRepository,
		inviteTokenRepository:  inviteTokenRepository,
		refreshTokenRepository: refreshTokenRepository,
		dbPool:                 dbPool,
		emailService:           emailService,
		jwtSecret:              jwtSecret,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, request *models.CreateUserRequest, creatorRole models.UserRole) (*models.CreateUserResponse, error) {
	role := s.parseRole(request.Role)
	if role == models.RoleContractor {
		return nil, fmt.Errorf("волонтери реєструються самостійно через сторінку реєстрації")
	}

	if !creatorRole.CanCreate(role) {
		return nil, fmt.Errorf("ваша посада не має прав для створення користувача з роллю %s", role)
	}

	tx, err := s.dbPool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	username := s.deriveUsername(request)

	var phone *string
	if request.Phone != nil {
		phone = request.Phone
	}

	now := time.Now()
	user := &models.User{
		Username:  &username,
		Email:     request.Email,
		FullName:  request.FullName,
		Phone:     phone,
		Status:    models.StatusPending,
		Role:      role,
		UnitID:    request.UnitID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.userRepository.CreateUser(ctx, tx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	rawToken, err := tokens.GenerateInviteToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate invite token: %w", err)
	}
	tokenHash := tokens.HashToken(rawToken)
	expiresAt := time.Now().Add(24 * time.Hour)

	inviteToken := &models.InviteToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	if err := s.inviteTokenRepository.CreateInviteToken(ctx, tx, inviteToken); err != nil {
		return nil, fmt.Errorf("failed to create invite token: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	baseURL := os.Getenv("FRONTEND_URL")
	if baseURL == "" {
		baseURL = "http://localhost"
	}
	inviteLink := fmt.Sprintf("%s/setup-password?token=%s", strings.TrimSuffix(baseURL, "/"), rawToken)

	go func(email, link string) {
		sendErr := s.emailService.SendInviteEmail(email, link)
		if sendErr != nil {
			fmt.Printf("ERROR: Failed to send email to %s: %v\n", email, sendErr)
		}
	}(user.Email, inviteLink)

	return &models.CreateUserResponse{
		ID:      user.ID,
		Message: "Користувача створено. На пошту надіслано лист для встановлення паролю.",
	}, nil
}

func (s *AuthService) RegisterCONTRACTOR(ctx context.Context, email, password, fullName string) (*models.CreateUserResponse, error) {
	existing, _ := s.userRepository.GetByEmail(ctx, s.dbPool, email)
	if existing != nil {
		return nil, fmt.Errorf("користувач з таким email вже існує")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	hashStr := string(hash)
	username := s.deriveUsername(&models.CreateUserRequest{Email: email})

	now := time.Now()
	user := &models.User{
		Username:     &username,
		Email:        email,
		FullName:     fullName,
		PasswordHash: &hashStr,
		Role:         models.RoleContractor,
		Status:       models.StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepository.CreateUser(ctx, s.dbPool, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &models.CreateUserResponse{
		ID:      user.ID,
		Message: "Реєстрація успішна. Можете увійти.",
	}, nil
}

func (s *AuthService) deriveUsername(request *models.CreateUserRequest) string {
	if request.Username != nil && strings.TrimSpace(*request.Username) != "" {
		return strings.TrimSpace(*request.Username)
	}
	at := strings.Index(request.Email, "@")
	if at > 0 {
		return request.Email[:at]
	}
	return request.Email
}

func (s *AuthService) GetCommanders(ctx context.Context) ([]*models.User, error) {
	commanders, err := s.userRepository.GetCommanders(ctx, s.dbPool)
	if err != nil {
		return nil, fmt.Errorf("помилка отримання командирів: %w", err)
	}

	for _, cmdr := range commanders {
		cmdr.PasswordHash = nil
	}

	return commanders, nil
}

func (s *AuthService) BootstrapAdmin(ctx context.Context, email, password, fullName string) error {
	if os.Getenv("ALLOW_BOOTSTRAP_OVERRIDE") != "true" {
		var count int
		err := s.dbPool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check users: %w", err)
		}
		if count > 0 {
			return fmt.Errorf("bootstrap only allowed when no users exist")
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	hashStr := string(hash)
	username := s.deriveUsername(&models.CreateUserRequest{Email: email})

	now := time.Now()
	user := &models.User{
		Username:     &username,
		Email:        email,
		FullName:     fullName,
		PasswordHash: &hashStr,
		Role:         models.RoleAdmin,
		Status:       models.StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return s.userRepository.CreateUser(ctx, s.dbPool, user)
}

func (s *AuthService) parseRole(r string) models.UserRole {
	switch r {
	case "ADMIN":
		return models.RoleAdmin
	case "REGION_DIRECTOR":
		return models.RoleRegionDirector
	case "BRANCH_MANAGER":
		return models.RoleBranchManager
	case "DEPT_MANAGER":
		return models.RoleDeptManager
	case "TEAM_LEAD":
		return models.RoleTeamLead
	case "REGION_LOGISTICIAN":
		return models.RoleRegionLogistician
	case "REGION_STOREKEEPER":
		return models.RoleRegionStorekeeper
	case "BRANCH_LOGISTICIAN":
		return models.RoleBranchLogistician
	case "BRANCH_STOREKEEPER":
		return models.RoleBranchStorekeeper
	case "DEPT_SUPERVISOR":
		return models.RoleDeptSupervisor
	case "CONTRACTOR":
		return models.RoleContractor
	default:
		return models.RoleContractor
	}
}

func (s *AuthService) SetupPassword(ctx context.Context, token, password string) error {
	tokenHash := tokens.HashToken(token)
	tx, err := s.dbPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	invToken, err := s.inviteTokenRepository.FindByTokenHash(ctx, tx, tokenHash)
	if err != nil {
		return fmt.Errorf("invalid or expired token")
	}
	if invToken.UsedAt != nil {
		return fmt.Errorf("token already used")
	}
	if time.Now().After(invToken.ExpiresAt) {
		return fmt.Errorf("token expired")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	hashStr := string(hash)

	if err := s.userRepository.UpdatePassword(ctx, tx, invToken.UserID, hashStr); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	if err := s.userRepository.UpdateStatus(ctx, tx, invToken.UserID, models.StatusActive); err != nil {
		return fmt.Errorf("failed to activate user: %w", err)
	}
	if err := s.inviteTokenRepository.MarkAsUsed(ctx, tx, invToken.ID); err != nil {
		return fmt.Errorf("failed to mark token used: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.LoginResponse, error) {
	user, err := s.userRepository.GetByEmail(ctx, s.dbPool, email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}
	if user.PasswordHash == nil || *user.PasswordHash == "" {
		return nil, fmt.Errorf("please set up your password first (check your email)")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}
	if user.Status != models.StatusActive {
		return nil, fmt.Errorf("account is not active")
	}

	var unitID int64

	if user.UnitID != nil {
		unitID = *user.UnitID
	}
	accessExpiresAt := time.Now().Add(24 * time.Hour)
	claims := &middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		UnitID: unitID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	refreshRaw, err := tokens.GenerateInviteToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	refreshHash := tokens.HashToken(refreshRaw)
	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	rt := &models.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		TokenHash: refreshHash,
		ExpiresAt: refreshExpiresAt,
		CreatedAt: time.Now(),
	}
	if err := s.refreshTokenRepository.Create(ctx, s.dbPool, rt); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &models.LoginResponse{
		Token:        tokenStr,
		RefreshToken: refreshRaw,
		ExpiresAt:    accessExpiresAt.Unix(),
		User:         *user,
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*models.LoginResponse, error) {
	tokenHash := tokens.HashToken(refreshToken)
	rt, err := s.refreshTokenRepository.FindByTokenHash(ctx, s.dbPool, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired refresh token")
	}
	if time.Now().After(rt.ExpiresAt) {
		return nil, fmt.Errorf("refresh token expired")
	}

	user, err := s.userRepository.GetByID(ctx, s.dbPool, rt.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	if user.Status != models.StatusActive {
		return nil, fmt.Errorf("account is not active")
	}

	accessExpiresAt := time.Now().Add(24 * time.Hour)
	claims := &middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	return &models.LoginResponse{
		Token:        tokenStr,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiresAt.Unix(),
		User:         *user,
	}, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	return s.userRepository.GetByID(ctx, s.dbPool, userID)
}

func (s *AuthService) GetVisibleUsers(ctx context.Context, requesterRole string, requesterUnitID *int64) ([]*models.User, error) {
	return s.userRepository.GetVisibleUsers(ctx, s.dbPool, requesterRole, requesterUnitID)
}

func (s *AuthService) UpdateRoleAndUnit(ctx context.Context, commanderID string, targetUserID string, newRole string, newUnitID *int64) error {

	commander, err := s.userRepository.GetByID(ctx, s.dbPool, commanderID)
	if err != nil {
		return fmt.Errorf("помилка авторизації командира")
	}

	if !s.isRoleChangePermitted(string(commander.Role), newRole) {
		return fmt.Errorf("субординація: ви не маєте прав призначати на посаду %s", newRole)
	}

	if newUnitID != nil && commander.Role != models.RoleAdmin {
		if commander.UnitID == nil {
			return fmt.Errorf("ви не прив'язані до підрозділу")
		}

		allowedUnits, err := s.unitRepository.GetSubordinateUnitIDs(ctx, s.dbPool, *commander.UnitID)
		if err != nil {
			return fmt.Errorf("помилка бази даних при перевірці ієрархії")
		}

		hasAccess := false
		for _, id := range allowedUnits {
			if id == *newUnitID {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			return fmt.Errorf("помилка безпеки: підрозділ не входить до вашого формування")
		}
	}

	return s.userRepository.UpdateRoleAndUnit(ctx, s.dbPool, targetUserID, newRole, newUnitID)
}

func (s *AuthService) isRoleChangePermitted(commanderRole string, targetRole string) bool {
	if commanderRole == string(models.RoleAdmin) {
		return true
	}

	allowedCommanders, exists := models.ApprovalMatrix[models.UserRole(targetRole)]
	if !exists {
		return false
	}

	for _, allowedRole := range allowedCommanders {
		if commanderRole == string(allowedRole) {
			return true
		}
	}

	return false
}

func (s *AuthService) BlockUser(ctx context.Context, userID string) error {
	return s.userRepository.BlockUser(ctx, s.dbPool, userID)
}

func (s *AuthService) CheckSubordination(ctx context.Context, commanderUnitID int64, targetUserID string) (bool, error) {
	return s.userRepository.CheckSubordination(ctx, s.dbPool, commanderUnitID, targetUserID)
}

func (s *AuthService) UnblockUser(ctx context.Context, userID string) error {
	return s.userRepository.UnblockUser(ctx, s.dbPool, userID)
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID string, req models.UpdateProfileRequest) error {
	if userID == "" {
		return errors.New("не вказано ID користувача")
	}

	if req.FullName == nil && req.Phone == nil && req.Username == nil && req.Email == nil {
		return errors.New("немає даних для оновлення")
	}

	err := s.userRepository.UpdateProfile(ctx, s.dbPool, userID, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) UpdateMyPassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	// 1. Отримуємо користувача
	user, err := s.userRepository.GetByID(ctx, s.dbPool, userID)
	if err != nil {
		return errors.New("користувача не знайдено")
	}

	if user.PasswordHash == nil || *user.PasswordHash == "" {
		return errors.New("у вас ще не встановлено пароль")
	}

	// 2. Перевіряємо старий пароль
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("невірний поточний пароль")
	}

	// 3. Хешуємо новий пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("помилка хешування: %w", err)
	}

	// 4. Зберігаємо (використовуємо твій існуючий метод UpdatePassword)
	if err := s.userRepository.UpdatePassword(ctx, s.dbPool, userID, string(hash)); err != nil {
		return fmt.Errorf("помилка збереження пароля: %w", err)
	}

	go func(email string) {
		_ = s.emailService.SendPasswordChangedAlert(email)
	}(user.Email)

	return nil
}

func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) error {
	// 1. Шукаємо користувача
	user, err := s.userRepository.GetByEmail(ctx, s.dbPool, email)
	if err != nil {
		// БЕЗПЕКА: Якщо email не знайдено, ми не повертаємо помилку!
		// Інакше хакери зможуть перебирати пошти і дізнаватися, хто зареєстрований.
		return nil
	}

	// Якщо юзер заблокований - нічого не робимо
	if user.Status == models.StatusBlocked {
		return nil
	}

	// 2. Генеруємо токен (перевикористовуємо твою інфраструктуру інвайтів)
	rawToken, err := tokens.GenerateInviteToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}
	tokenHash := tokens.HashToken(rawToken)

	resetToken := &models.InviteToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(1 * time.Hour), // Лінк на скидання діє лише 1 годину!
		CreatedAt: time.Now(),
	}

	if err := s.inviteTokenRepository.CreateInviteToken(ctx, s.dbPool, resetToken); err != nil {
		return fmt.Errorf("failed to save reset token: %w", err)
	}

	// 3. Формуємо лінк. Ми геніально перевикористовуємо фронтенд-сторінку setup-password!
	baseURL := os.Getenv("FRONTEND_URL")
	if baseURL == "" {
		baseURL = "http://localhost"
	}
	resetLink := fmt.Sprintf("%s/setup-password?token=%s", strings.TrimSuffix(baseURL, "/"), rawToken)

	// 4. Відправляємо лист у фоні
	go func(userEmail, link string) {
		_ = s.emailService.SendPasswordResetEmail(userEmail, link)
	}(user.Email, resetLink)

	return nil
}

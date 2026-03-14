package services

import (
	"context"
	"fmt"
	"millog_backend/internal/middleware"
	"os"
	"millog_backend/internal/models"
	"millog_backend/internal/repositories"
	"millog_backend/internal/tokens"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepository         *repositories.UserRepository
	inviteTokenRepository  *repositories.InviteTokenRepository
	refreshTokenRepository *repositories.RefreshTokenRepository
	dbPool                 *pgxpool.Pool
	emailService           EmailService
	jwtSecret              string
}

func NewAuthService(
	userRepository *repositories.UserRepository,
	inviteTokenRepository *repositories.InviteTokenRepository,
	refreshTokenRepository *repositories.RefreshTokenRepository,
	dbPool *pgxpool.Pool,
	emailService EmailService,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		userRepository:         userRepository,
		inviteTokenRepository:  inviteTokenRepository,
		refreshTokenRepository: refreshTokenRepository,
		dbPool:                 dbPool,
		emailService:          emailService,
		jwtSecret:              jwtSecret,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, request *models.CreateUserRequest) (*models.CreateUserResponse, error) {
	role := s.parseRole(request.Role)
	if role == models.RoleVolunteer {
		return nil, fmt.Errorf("волонтери реєструються самостійно через сторінку реєстрації")
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

func (s *AuthService) RegisterVolunteer(ctx context.Context, email, password, fullName string) (*models.CreateUserResponse, error) {
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
		Role:         models.RoleVolunteer,
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
	case "BRIGADE_CMDR":
		return models.RoleBrigadeCmdr
	case "BATTALION_CMDR":
		return models.RoleBattalionCmdr
	case "COMPANY_CMDR":
		return models.RoleCompanyCmdr
	case "PLATOON_CMDR":
		return models.RolePlatoonCmdr
	case "BRIGADE_LOGIST":
		return models.RoleBrigadeLogist
	case "BRIGADE_STOREKEEPER":
		return models.RoleBrigadeStorekeeper
	case "BATTALION_LOGIST":
		return models.RoleBattalionLogist
	case "BATTALION_STOREKEEPER":
		return models.RoleBattalionStorekeeper
	case "COMPANY_SERGEANT":
		return models.RoleCompanySergeant
	case "VOLUNTEER":
		return models.RoleVolunteer
	default:
		return models.RoleVolunteer
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
		User:        *user,
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
		User:        *user,
	}, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	return s.userRepository.GetByID(ctx, s.dbPool, userID)
}

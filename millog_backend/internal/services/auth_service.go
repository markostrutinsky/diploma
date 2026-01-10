package services

import (
	"context"
	"fmt"
	"millog_backend/internal/models"
	"millog_backend/internal/repositories"
	"millog_backend/internal/tokens"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthService struct {
	userRepository        *repositories.UserRepository
	inviteTokenRepository *repositories.InviteTokenRepository
	dbPool                *pgxpool.Pool
	emailService          EmailService
}

func NewAuthService(userRepository *repositories.UserRepository, inviteTokenRepository *repositories.InviteTokenRepository, dbPool *pgxpool.Pool, emailService EmailService) *AuthService {
	return &AuthService{
		userRepository:        userRepository,
		inviteTokenRepository: inviteTokenRepository,
		dbPool:                dbPool,
		emailService:          emailService,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, request *models.CreateUserRequest) (*models.CreateUserResponse, error) {
	tx, err := s.dbPool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var username *string
	if request.Username != nil {
		username = request.Username
	}

	var phone *string
	if request.Phone != nil {
		phone = request.Phone
	}
	user := &models.User{
		Username: username,
		Email:    request.Email,
		FullName: request.FullName,
		Phone:    phone,
		Status:   models.StatusPending,
		Role:     models.RoleVolunteer,
		UnitID:   request.UnitID,
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

	// Create the invite token with database-generated ID
	inviteToken := &models.InviteToken{
		UserID:    user.ID,
		User:      *user,
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

	inviteLink := fmt.Sprintf("https://millog.system/setup-password?token=%s", rawToken)

	go func(email, link string) {
		// Тут бажано створити новий контекст, наприклад context.Background()
		// бо ctx запиту може закритися раніше, ніж пошта відправиться
		sendErr := s.emailService.SendInviteEmail(email, link)
		if sendErr != nil {
			fmt.Printf("ERROR: Failed to send email to %s: %v\n", email, sendErr)
		}
	}(user.Email, inviteLink)

	return &models.CreateUserResponse{
		ID:      user.ID,
		Message: "User registered successfully. Please wait for admin approval.",
	}, nil
}

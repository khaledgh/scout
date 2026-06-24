package services

import (
	"context"
	"kashfi/internal/config"
	"kashfi/internal/models"
	"kashfi/internal/utils"
	"time"

	"gorm.io/gorm"
)

type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	User         *models.User
}

func (s *AuthService) Login(ctx context.Context, phone, password string) (*LoginResult, error) {
	phone = utils.NormalizePhone(phone)
	var user models.User
	if err := s.db.WithContext(ctx).Where("phone = ? AND is_active = true", phone).First(&user).Error; err != nil {
		return nil, ErrUnauthorized
	}
	if !utils.CheckPassword(password, user.PasswordHash) {
		return nil, ErrUnauthorized
	}

	now := time.Now()
	s.db.WithContext(ctx).Model(&user).Update("last_login_at", &now)

	return s.buildTokens(&user)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*LoginResult, error) {
	claims, err := utils.ParseToken(refreshToken, s.cfg.JWT.Secret)
	if err != nil || claims.Type != "refresh" {
		return nil, ErrUnauthorized
	}

	var user models.User
	if err := s.db.WithContext(ctx).Where("id = ? AND is_active = true", claims.UserID).First(&user).Error; err != nil {
		return nil, ErrUnauthorized
	}

	return s.buildTokens(&user)
}

func (s *AuthService) Me(ctx context.Context, userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return nil, ErrNotFound
	}
	return &user, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uint, current, newPass string) error {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return ErrNotFound
	}
	if !utils.CheckPassword(current, user.PasswordHash) {
		return ErrUnauthorized
	}
	hash, err := utils.HashPassword(newPass, s.cfg.JWT.BcryptCost)
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Model(&user).Update("password_hash", hash).Error
}

func (s *AuthService) buildTokens(user *models.User) (*LoginResult, error) {
	access, err := utils.GenerateAccessToken(user.ID, string(user.Role), s.cfg.JWT.Secret, s.cfg.JWT.AccessTTL)
	if err != nil {
		return nil, err
	}
	refresh, err := utils.GenerateRefreshToken(user.ID, string(user.Role), s.cfg.JWT.Secret, s.cfg.JWT.RefreshTTL)
	if err != nil {
		return nil, err
	}
	return &LoginResult{AccessToken: access, RefreshToken: refresh, User: user}, nil
}

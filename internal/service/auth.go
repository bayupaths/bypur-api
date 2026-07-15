package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/bayupaths/bypur-api/internal/config"
	"github.com/bayupaths/bypur-api/internal/model"
	"github.com/bayupaths/bypur-api/internal/repository"
	"github.com/bayupaths/bypur-api/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type UserResponse struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	FullName *string `json:"fullName"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type LoginResponse struct {
	User   UserResponse  `json:"user"`
	Tokens TokenResponse `json:"tokens"`
}

type AuthProfileUpdateData struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	FullName *string `json:"fullName"`
}

type AuthService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.RefreshTokenRepository
	cfg       *config.Config
}

const (
	MaxFailedAttempts = 5
	LockoutDuration   = 30 * time.Minute
)

func NewAuthService(userRepo repository.UserRepository, tokenRepo repository.RefreshTokenRepository, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, tokenRepo: tokenRepo, cfg: cfg}
}

// HashPassword melakukan hashing password menggunakan bcrypt
func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// VerifyPassword mencocokkan password mentah dengan hash
func (s *AuthService) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateAccessToken membuat JWT access token menggunakan pkg/jwt
func (s *AuthService) GenerateAccessToken(user *model.User) (string, error) {
	duration, err := time.ParseDuration(s.cfg.JWT.AccessExpire)
	if err != nil {
		duration = 15 * time.Minute
	}

	return jwt.GenerateToken(user.ID, user.Email, s.cfg.JWT.Secret, duration, "bypur-api")
}

// GenerateRefreshToken membuat token refresh, menyimpannya di DB dan mengembalikan token string
func (s *AuthService) GenerateRefreshToken(ctx context.Context, userID string) (string, error) {
	duration, err := time.ParseDuration(s.cfg.JWT.RefreshExpire)
	if err != nil {
		duration = 7 * 24 * time.Hour
	}

	tokenStr, err := jwt.GenerateToken(userID, "", s.cfg.JWT.Secret+"refresh", duration, "bypur-api")
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(duration)
	rt := model.RefreshToken{
		Token:     tokenStr,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	if err := s.tokenRepo.Create(ctx, &rt); err != nil {
		return "", err
	}

	return tokenStr, nil
}

// Login melakukan proses autentikasi admin
func (s *AuthService) Login(ctx context.Context, identifier, password, ip string) (*LoginResponse, error) {
	slog.Info("Login attempt", "identifier", identifier, "ip", ip)

	user, err := s.userRepo.GetByUsernameOrEmail(ctx, identifier)
	if err != nil {
		slog.Warn("Login failed - user not found", "identifier", identifier)
		// Timing attack mitigation
		_ = bcrypt.CompareHashAndPassword([]byte("$2b$12$FakeHashForTimingAttacksDontLeakExistence"), []byte(password))
		return nil, errors.New("user not found")
	}

	// Periksa apakah akun dikunci
	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		slog.Warn("Login blocked - account locked", "userId", user.ID, "lockedUntil", user.LockedUntil)
		return nil, errors.New("account is temporarily locked due to too many login failures")
	}

	// Verifikasi sandi
	if !s.VerifyPassword(password, user.PasswordHash) {
		user.FailedLoginAttempts++
		if user.FailedLoginAttempts >= MaxFailedAttempts {
			lockTime := time.Now().Add(LockoutDuration)
			user.LockedUntil = &lockTime
		}
		_ = s.userRepo.Update(ctx, user)

		slog.Warn("Login failed - incorrect password", "identifier", identifier)
		return nil, errors.New("incorrect password")
	}

	// Reset kegagalan login & catat sukses login
	now := time.Now()
	user.FailedLoginAttempts = 0
	user.LockedUntil = nil
	user.LastLoginAt = &now
	user.LastLoginIp = &ip
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Generate token
	accessToken, err := s.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	slog.Info("User logged in successfully", "userId", user.ID)

	return &LoginResponse{
		User: UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
		},
		Tokens: TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

// RefreshAccessToken memperbarui access token menggunakan refresh token yang masih valid
func (s *AuthService) RefreshAccessToken(ctx context.Context, tokenStr string) (string, error) {
	// Verifikasi refresh token menggunakan pkg/jwt
	claims, err := jwt.VerifyToken(tokenStr, s.cfg.JWT.Secret+"refresh")
	if err != nil {
		return "", errors.New("invalid or expired refresh token")
	}

	// Cari token di DB untuk memastikan belum di-revoke
	_, err = s.tokenRepo.GetByTokenAndUser(ctx, tokenStr, claims.UserID)
	if err != nil {
		return "", errors.New("refresh token is not registered or has been revoked")
	}

	// Ambil data user
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return "", errors.New("user not found")
	}

	if user.Status != "active" {
		return "", errors.New("account status is inactive")
	}

	// Buat access token baru
	accessToken, err := s.GenerateAccessToken(user)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// Logout mencabut validitas refresh token
func (s *AuthService) Logout(ctx context.Context, tokenStr string) error {
	now := time.Now()
	err := s.tokenRepo.Revoke(ctx, tokenStr, now)
	if err != nil {
		slog.Error("Logout failed to revoke token", "error", err)
		return err
	}

	slog.Info("User logged out successfully")
	return nil
}

// GetUserByID mengambil data user yang disanitasi
func (s *AuthService) GetUserByID(ctx context.Context, id string) (*UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
	}, nil
}

// UpdateProfile memperbarui profil admin
func (s *AuthService) UpdateProfile(ctx context.Context, id string, data *AuthProfileUpdateData) (*UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if data.Username != nil {
		user.Username = *data.Username
	}
	if data.Email != nil {
		user.Email = *data.Email
	}
	if data.FullName != nil {
		user.FullName = data.FullName
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
	}, nil
}

// ChangePassword mengganti password admin
func (s *AuthService) ChangePassword(ctx context.Context, id string, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}

	if !s.VerifyPassword(currentPassword, user.PasswordHash) {
		return errors.New("current password is incorrect")
	}

	if currentPassword == newPassword {
		return errors.New("new password cannot be the same as the old password")
	}

	newHash, err := s.HashPassword(newPassword)
	if err != nil {
		return err
	}

	now := time.Now()
	user.PasswordHash = newHash
	user.PasswordChangedAt = &now
	user.MustChangePassword = false

	return s.userRepo.Update(ctx, user)
}

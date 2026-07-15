package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bayupaths/bypur-api/internal/config"
	"github.com/bayupaths/bypur-api/internal/model"
	"github.com/bayupaths/bypur-api/internal/service"
)

type fakeAuthUserRepo struct {
	users        map[string]*model.User
	byIdentifier map[string]string
	updateErr    error
	lastUpdated  *model.User
}

func (f *fakeAuthUserRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	user := f.users[id]
	if user == nil {
		return nil, errors.New("not found")
	}
	return user, nil
}

func (f *fakeAuthUserRepo) GetByUsernameOrEmail(ctx context.Context, identifier string) (*model.User, error) {
	id := f.byIdentifier[identifier]
	if id == "" {
		return nil, errors.New("not found")
	}
	return f.GetByID(ctx, id)
}

func (f *fakeAuthUserRepo) Create(ctx context.Context, user *model.User) error {
	f.users[user.ID] = user
	f.byIdentifier[user.Username] = user.ID
	f.byIdentifier[user.Email] = user.ID
	return nil
}

func (f *fakeAuthUserRepo) Update(ctx context.Context, user *model.User) error {
	f.lastUpdated = user
	if f.updateErr != nil {
		return f.updateErr
	}
	f.users[user.ID] = user
	return nil
}

type fakeAuthRefreshTokenRepo struct {
	tokens       map[string]*model.RefreshToken
	createErr    error
	getErr       error
	revokeErr    error
	lastCreated  *model.RefreshToken
	revokedToken string
}

func (f *fakeAuthRefreshTokenRepo) Create(ctx context.Context, rt *model.RefreshToken) error {
	f.lastCreated = rt
	if f.createErr != nil {
		return f.createErr
	}
	f.tokens[rt.Token] = rt
	return nil
}

func (f *fakeAuthRefreshTokenRepo) GetByTokenAndUser(ctx context.Context, tokenStr string, userID string) (*model.RefreshToken, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	rt := f.tokens[tokenStr]
	if rt == nil || rt.UserID != userID || rt.RevokedAt != nil || time.Now().After(rt.ExpiresAt) {
		return nil, errors.New("not found")
	}
	return rt, nil
}

func (f *fakeAuthRefreshTokenRepo) Revoke(ctx context.Context, tokenStr string, revokedAt time.Time) error {
	f.revokedToken = tokenStr
	if f.revokeErr != nil {
		return f.revokeErr
	}
	if rt := f.tokens[tokenStr]; rt != nil {
		rt.RevokedAt = &revokedAt
	}
	return nil
}

func authTestConfig() *config.Config {
	return &config.Config{
		JWT: config.JWTConfig{
			Secret:        "01234567890123456789012345678901",
			AccessExpire:  "15m",
			RefreshExpire: "24h",
		},
	}
}

func newAuthTestService(user *model.User) (*service.AuthService, *fakeAuthUserRepo, *fakeAuthRefreshTokenRepo) {
	userRepo := &fakeAuthUserRepo{
		users:        map[string]*model.User{user.ID: user},
		byIdentifier: map[string]string{user.Username: user.ID, user.Email: user.ID},
	}
	tokenRepo := &fakeAuthRefreshTokenRepo{tokens: map[string]*model.RefreshToken{}}
	return service.NewAuthService(userRepo, tokenRepo, authTestConfig()), userRepo, tokenRepo
}

func TestAuthServiceLoginSuccessResetsFailuresAndStoresRefreshToken(t *testing.T) {
	svc := service.NewAuthService(nil, nil, authTestConfig())
	hash, err := svc.HashPassword("Correct1")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	user := &model.User{
		ID:                  "user-1",
		Username:            "bayu",
		Email:               "bayu@example.com",
		PasswordHash:        hash,
		Status:              "active",
		FailedLoginAttempts: 3,
	}
	svc, userRepo, tokenRepo := newAuthTestService(user)

	result, err := svc.Login(context.Background(), "bayu", "Correct1", "127.0.0.1")
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}
	if result.User.ID != "user-1" || result.Tokens.AccessToken == "" || result.Tokens.RefreshToken == "" {
		t.Fatalf("unexpected login result: %+v", result)
	}
	if userRepo.lastUpdated == nil || userRepo.lastUpdated.FailedLoginAttempts != 0 || userRepo.lastUpdated.LastLoginIp == nil {
		t.Fatalf("successful login did not reset audit fields: %+v", userRepo.lastUpdated)
	}
	if tokenRepo.lastCreated == nil || tokenRepo.lastCreated.UserID != "user-1" {
		t.Fatalf("refresh token was not persisted correctly: %+v", tokenRepo.lastCreated)
	}
}

func TestAuthServiceLoginLocksAccountAfterMaxFailures(t *testing.T) {
	svc := service.NewAuthService(nil, nil, authTestConfig())
	hash, err := svc.HashPassword("Correct1")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	user := &model.User{
		ID:                  "user-1",
		Username:            "bayu",
		Email:               "bayu@example.com",
		PasswordHash:        hash,
		Status:              "active",
		FailedLoginAttempts: service.MaxFailedAttempts - 1,
	}
	svc, userRepo, _ := newAuthTestService(user)

	if _, err := svc.Login(context.Background(), "bayu", "Wrong1", "127.0.0.1"); err == nil {
		t.Fatal("expected incorrect password error")
	}
	if userRepo.lastUpdated == nil || userRepo.lastUpdated.FailedLoginAttempts != service.MaxFailedAttempts || userRepo.lastUpdated.LockedUntil == nil {
		t.Fatalf("failed login did not lock account: %+v", userRepo.lastUpdated)
	}
}

func TestAuthServiceLoginRejectsCurrentlyLockedAccount(t *testing.T) {
	svc := service.NewAuthService(nil, nil, authTestConfig())
	hash, err := svc.HashPassword("Correct1")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	lockedUntil := time.Now().Add(time.Hour)
	user := &model.User{
		ID:           "user-1",
		Username:     "bayu",
		Email:        "bayu@example.com",
		PasswordHash: hash,
		Status:       "active",
		LockedUntil:  &lockedUntil,
	}
	svc, _, tokenRepo := newAuthTestService(user)

	if _, err := svc.Login(context.Background(), "bayu", "Correct1", "127.0.0.1"); err == nil {
		t.Fatal("expected locked account error")
	}
	if tokenRepo.lastCreated != nil {
		t.Fatal("locked account should not receive refresh token")
	}
}

func TestAuthServiceRefreshAccessTokenRequiresRegisteredActiveUser(t *testing.T) {
	user := &model.User{ID: "user-1", Username: "bayu", Email: "bayu@example.com", Status: "active"}
	svc, _, tokenRepo := newAuthTestService(user)

	refreshToken, err := svc.GenerateRefreshToken(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("GenerateRefreshToken returned error: %v", err)
	}

	accessToken, err := svc.RefreshAccessToken(context.Background(), refreshToken)
	if err != nil {
		t.Fatalf("RefreshAccessToken returned error: %v", err)
	}
	if accessToken == "" {
		t.Fatal("expected access token")
	}

	if err := tokenRepo.Revoke(context.Background(), refreshToken, time.Now()); err != nil {
		t.Fatalf("Revoke returned error: %v", err)
	}
	if _, err := svc.RefreshAccessToken(context.Background(), refreshToken); err == nil {
		t.Fatal("expected revoked refresh token to be rejected")
	}
}

func TestAuthServiceRefreshAccessTokenRejectsInactiveUser(t *testing.T) {
	user := &model.User{ID: "user-1", Username: "bayu", Email: "bayu@example.com", Status: "suspended"}
	svc, _, _ := newAuthTestService(user)

	refreshToken, err := svc.GenerateRefreshToken(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("GenerateRefreshToken returned error: %v", err)
	}

	if _, err := svc.RefreshAccessToken(context.Background(), refreshToken); err == nil {
		t.Fatal("expected inactive account to be rejected")
	}
}

func TestAuthServiceChangePasswordValidatesCurrentAndNewPassword(t *testing.T) {
	svc := service.NewAuthService(nil, nil, authTestConfig())
	hash, err := svc.HashPassword("Current1")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	user := &model.User{ID: "user-1", Username: "bayu", Email: "bayu@example.com", PasswordHash: hash, MustChangePassword: true}
	svc, userRepo, _ := newAuthTestService(user)

	if err := svc.ChangePassword(context.Background(), "user-1", "Wrong1", "NextPass1"); err == nil {
		t.Fatal("expected current password validation error")
	}
	if err := svc.ChangePassword(context.Background(), "user-1", "Current1", "Current1"); err == nil {
		t.Fatal("expected same password validation error")
	}
	if err := svc.ChangePassword(context.Background(), "user-1", "Current1", "NextPass1"); err != nil {
		t.Fatalf("ChangePassword returned error: %v", err)
	}

	if userRepo.lastUpdated == nil || userRepo.lastUpdated.PasswordChangedAt == nil || userRepo.lastUpdated.MustChangePassword {
		t.Fatalf("password metadata was not updated: %+v", userRepo.lastUpdated)
	}
	if !svc.VerifyPassword("NextPass1", userRepo.lastUpdated.PasswordHash) {
		t.Fatal("new password hash does not verify")
	}
}

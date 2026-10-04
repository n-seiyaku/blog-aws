package auth

import (
	"blog-aws-backend/internal/domain/session"
	"blog-aws-backend/internal/domain/user"
	"blog-aws-backend/internal/infrastructure/token"
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepository struct {
	usersByEmail map[string]user.User
	usersByID    map[string]user.User

	createErr error
	findErr   error
	getErr    error

	created *user.User
}

func (r *fakeUserRepository) Create(ctx context.Context, newUser user.User) error {
	if r.createErr != nil {
		return r.createErr
	}
	r.created = &newUser
	return nil
}

func (r *fakeUserRepository) Update(ctx context.Context, updatedUser user.User) error {
	return nil
}

func (r *fakeUserRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *fakeUserRepository) GetUserByID(ctx context.Context, id string) (user.User, error) {
	if r.getErr != nil {
		return user.User{}, r.getErr
	}
	foundUser, ok := r.usersByID[id]
	if !ok {
		return user.User{}, user.ErrorNotFound
	}
	return foundUser, nil
}

func (r *fakeUserRepository) FindByEmail(ctx context.Context, email string) (user.User, error) {
	if r.findErr != nil {
		return user.User{}, r.findErr
	}
	foundUser, ok := r.usersByEmail[email]
	if !ok {
		return user.User{}, user.ErrorNotFound
	}
	return foundUser, nil
}

type fakeSessionRepository struct {
	sessionsByHash map[string]session.Session

	createErr error
	findErr   error
	rotateErr error
	revokeErr error

	created []session.Session
	rotated *rotateCall
	revoked []string
}

type rotateCall struct {
	sessionID        string
	refreshTokenHash string
	expiresAt        time.Time
}

func (r *fakeSessionRepository) Create(ctx context.Context, newSession session.Session) error {
	if r.createErr != nil {
		return r.createErr
	}
	r.created = append(r.created, newSession)
	return nil
}

func (r *fakeSessionRepository) FindByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (session.Session, error) {
	if r.findErr != nil {
		return session.Session{}, r.findErr
	}
	foundSession, ok := r.sessionsByHash[refreshTokenHash]
	if !ok {
		return session.Session{}, session.ErrorNotFound
	}
	return foundSession, nil
}

func (r *fakeSessionRepository) Rotate(ctx context.Context, sessionID string, refreshTokenHash string, expiresAt time.Time) error {
	if r.rotateErr != nil {
		return r.rotateErr
	}
	r.rotated = &rotateCall{
		sessionID:        sessionID,
		refreshTokenHash: refreshTokenHash,
		expiresAt:        expiresAt,
	}
	return nil
}

func (r *fakeSessionRepository) Revoke(ctx context.Context, sessionID string) error {
	if r.revokeErr != nil {
		return r.revokeErr
	}
	r.revoked = append(r.revoked, sessionID)
	return nil
}

type fakeTokenGenerator struct {
	token string
	err   error

	userID string
	email  string
	role   string
}

func (g *fakeTokenGenerator) GenerateToken(userID string, email string, role string) (string, error) {
	g.userID = userID
	g.email = email
	g.role = role

	if g.err != nil {
		return "", g.err
	}
	return g.token, nil
}

func TestRegisterUsecaseCreatesNormalizedUserWithHashedPassword(t *testing.T) {
	repository := &fakeUserRepository{
		usersByEmail: map[string]user.User{},
	}
	usecase := NewRegisterUsecase(repository)

	newUser, err := usecase.Execute(context.Background(), RegisterRequest{
		Email:    " User@Example.COM ",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if newUser.Email != "user@example.com" {
		t.Fatalf("expected normalized email, got %q", newUser.Email)
	}
	if newUser.ID == "" {
		t.Fatal("expected generated user id")
	}
	if newUser.Role != user.RoleUser {
		t.Fatalf("expected role %q, got %q", user.RoleUser, newUser.Role)
	}
	if newUser.PasswordHash == "password123" {
		t.Fatal("expected password to be hashed")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(newUser.PasswordHash), []byte("password123")); err != nil {
		t.Fatalf("password hash does not match password: %v", err)
	}
	if repository.created == nil {
		t.Fatal("expected repository Create to be called")
	}
}

func TestRegisterUsecaseRejectsDuplicateEmail(t *testing.T) {
	repository := &fakeUserRepository{
		usersByEmail: map[string]user.User{
			"taken@example.com": {ID: "user-1", Email: "taken@example.com"},
		},
	}
	usecase := NewRegisterUsecase(repository)

	_, err := usecase.Execute(context.Background(), RegisterRequest{
		Email:    " taken@example.com ",
		Password: "password123",
	})

	if !errors.Is(err, user.ErrorEmailAlreadyExists) {
		t.Fatalf("expected ErrorEmailAlreadyExists, got %v", err)
	}
}

func TestLoginUsecaseReturnsTokensAndCreatesRefreshSession(t *testing.T) {
	passwordHash := hashPassword(t, "password123")
	repository := &fakeUserRepository{
		usersByEmail: map[string]user.User{
			"user@example.com": {
				ID:           "user-1",
				Email:        "user@example.com",
				PasswordHash: passwordHash,
				Role:         user.RoleAdmin,
			},
		},
	}
	sessionRepository := &fakeSessionRepository{}
	tokenGenerator := &fakeTokenGenerator{token: "access-token"}
	usecase := NewLoginUsecase(repository, sessionRepository, tokenGenerator)

	response, err := usecase.Execute(context.Background(), LoginRequest{
		Email:    " User@Example.COM ",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if response.AccessToken != "access-token" {
		t.Fatalf("expected access token, got %q", response.AccessToken)
	}
	if response.RefreshToken == "" {
		t.Fatal("expected refresh token")
	}
	if tokenGenerator.userID != "user-1" || tokenGenerator.email != "user@example.com" || tokenGenerator.role != "ADMIN" {
		t.Fatalf("unexpected token claims: userID=%q email=%q role=%q", tokenGenerator.userID, tokenGenerator.email, tokenGenerator.role)
	}
	if len(sessionRepository.created) != 1 {
		t.Fatalf("expected 1 session to be created, got %d", len(sessionRepository.created))
	}

	createdSession := sessionRepository.created[0]
	if createdSession.UserID != "user-1" {
		t.Fatalf("expected session user id user-1, got %q", createdSession.UserID)
	}
	if createdSession.RefreshTokenHash != token.HashRefreshToken(response.RefreshToken) {
		t.Fatal("expected session refresh token hash to match returned refresh token")
	}
	assertExpiresInAboutThirtyDays(t, createdSession.ExpiresAt)
}

func TestLoginUsecaseRejectsInvalidCredentials(t *testing.T) {
	passwordHash := hashPassword(t, "password123")
	repository := &fakeUserRepository{
		usersByEmail: map[string]user.User{
			"user@example.com": {
				ID:           "user-1",
				Email:        "user@example.com",
				PasswordHash: passwordHash,
				Role:         user.RoleUser,
			},
		},
	}
	usecase := NewLoginUsecase(repository, &fakeSessionRepository{}, &fakeTokenGenerator{token: "access-token"})

	_, err := usecase.Execute(context.Background(), LoginRequest{
		Email:    "user@example.com",
		Password: "wrong-password",
	})

	if !errors.Is(err, ErrorInvalidCredentials) {
		t.Fatalf("expected ErrorInvalidCredentials, got %v", err)
	}
}

func TestRefreshUsecaseRotatesRefreshToken(t *testing.T) {
	oldRefreshToken := "old-refresh-token"
	oldRefreshTokenHash := token.HashRefreshToken(oldRefreshToken)
	sessionRepository := &fakeSessionRepository{
		sessionsByHash: map[string]session.Session{
			oldRefreshTokenHash: {
				ID:               "session-1",
				UserID:           "user-1",
				RefreshTokenHash: oldRefreshTokenHash,
				ExpiresAt:        time.Now().Add(time.Hour),
			},
		},
	}
	userRepository := &fakeUserRepository{
		usersByID: map[string]user.User{
			"user-1": {
				ID:    "user-1",
				Email: "user@example.com",
				Role:  user.RoleUser,
			},
		},
	}
	usecase := NewRefreshUsecase(userRepository, sessionRepository, &fakeTokenGenerator{token: "new-access-token"})

	response, err := usecase.Execute(context.Background(), RefreshRequest{
		RefreshToken: oldRefreshToken,
	})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if response.AccessToken != "new-access-token" {
		t.Fatalf("expected new access token, got %q", response.AccessToken)
	}
	if response.RefreshToken == "" || response.RefreshToken == oldRefreshToken {
		t.Fatalf("expected new refresh token, got %q", response.RefreshToken)
	}
	if sessionRepository.rotated == nil {
		t.Fatal("expected session to be rotated")
	}
	if sessionRepository.rotated.sessionID != "session-1" {
		t.Fatalf("expected session-1 to be rotated, got %q", sessionRepository.rotated.sessionID)
	}
	if sessionRepository.rotated.refreshTokenHash != token.HashRefreshToken(response.RefreshToken) {
		t.Fatal("expected rotated hash to match new refresh token")
	}
	assertExpiresInAboutThirtyDays(t, sessionRepository.rotated.expiresAt)
}

func TestRefreshUsecaseRevokesExpiredSession(t *testing.T) {
	refreshToken := "expired-refresh-token"
	refreshTokenHash := token.HashRefreshToken(refreshToken)
	sessionRepository := &fakeSessionRepository{
		sessionsByHash: map[string]session.Session{
			refreshTokenHash: {
				ID:               "session-1",
				UserID:           "user-1",
				RefreshTokenHash: refreshTokenHash,
				ExpiresAt:        time.Now().Add(-time.Minute),
			},
		},
	}
	usecase := NewRefreshUsecase(&fakeUserRepository{}, sessionRepository, &fakeTokenGenerator{token: "access-token"})

	_, err := usecase.Execute(context.Background(), RefreshRequest{
		RefreshToken: refreshToken,
	})

	if !errors.Is(err, ErrorInvalidCredentials) {
		t.Fatalf("expected ErrorInvalidCredentials, got %v", err)
	}
	if len(sessionRepository.revoked) != 1 || sessionRepository.revoked[0] != "session-1" {
		t.Fatalf("expected expired session to be revoked, got %#v", sessionRepository.revoked)
	}
}

func TestLogoutUsecaseRevokesExistingSession(t *testing.T) {
	refreshToken := "refresh-token"
	refreshTokenHash := token.HashRefreshToken(refreshToken)
	sessionRepository := &fakeSessionRepository{
		sessionsByHash: map[string]session.Session{
			refreshTokenHash: {
				ID:               "session-1",
				UserID:           "user-1",
				RefreshTokenHash: refreshTokenHash,
			},
		},
	}
	usecase := NewLogoutUsecase(sessionRepository)

	if err := usecase.Execute(context.Background(), LogoutRequest{RefreshToken: refreshToken}); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if len(sessionRepository.revoked) != 1 || sessionRepository.revoked[0] != "session-1" {
		t.Fatalf("expected session to be revoked, got %#v", sessionRepository.revoked)
	}
}

func TestLogoutUsecaseIgnoresMissingSession(t *testing.T) {
	usecase := NewLogoutUsecase(&fakeSessionRepository{
		sessionsByHash: map[string]session.Session{},
	})

	if err := usecase.Execute(context.Background(), LogoutRequest{RefreshToken: "missing-refresh-token"}); err != nil {
		t.Fatalf("expected missing session to be ignored, got %v", err)
	}
}

func hashPassword(t *testing.T, password string) string {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	return string(hash)
}

func assertExpiresInAboutThirtyDays(t *testing.T, expiresAt time.Time) {
	t.Helper()

	duration := time.Until(expiresAt)
	if duration < 29*24*time.Hour || duration > 31*24*time.Hour {
		t.Fatalf("expected expiry about 30 days from now, got %s", duration)
	}
}

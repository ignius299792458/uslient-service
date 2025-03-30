package user

import (
	"context"
	"errors"
	"had-service/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService interface {

	// Auth
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	Logout(ctx context.Context, sessionID string) error

	// Auth helpers
	ValidateToken(tokenString string) (string, error)
	RefreshSession(ctx context.Context, sessionID string) (*LoginResponse, error)
	RequestPasswordReset(ctx context.Context, req PasswordResetRequest) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}

type AuthServiceConfig struct {
	JWTSecret        string
	TokenExpiration  time.Duration
	UploadDir        string
	DefaultAvatarURL string
}

// ServiceImpl implements the Service interface
type AuthServiceImpl struct {
	userRepo   UserRepository
	authConfig AuthServiceConfig
}

func NewAuthService(repo UserRepository, config AuthServiceConfig) AuthService {
	return &AuthServiceImpl{
		userRepo:   repo,
		authConfig: config,
	}
}

func (s *AuthServiceImpl) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email")
	}

	// Check password
	if !CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid password")
	}

	// Generate token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Create session
	session := &UserSession{
		UserID:    user.ID,
		TokenHash: token, // In production, you'd hash this
		ExpiresAt: time.Now().Add(s.config.TokenExpiration),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, errors.New("failed to create session")
	}

	// Return response
	return &LoginResponse{
		Token:     token,
		ExpiresAt: session.ExpiresAt,
		User:      user.ToResponse(),
	}, nil
}

func (s *AuthServiceImpl) Logout(ctx context.Context, sessionID string) error {
	//TODO implement me
	panic("implement me")
}

// ValidateToken validates a JWT token and returns the user ID
func (s *AuthServiceImpl) ValidateToken(tokenString string) (string, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(s.authConfig.JWTSecret), nil
	})

	if err != nil {
		return "", err
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check token expiration
		exp, ok := claims["exp"].(float64)
		if !ok {
			return "", errors.New("invalid token expiration")
		}

		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return "", errors.New("token expired")
		}

		// Get user ID
		userID, ok := claims["sub"].(string)
		if !ok {
			return "", errors.New("invalid user ID in token")
		}

		return userID, nil
	}

	return "", errors.New("invalid token")
}

func (s *AuthServiceImpl) RefreshSession(ctx context.Context, sessionID string) (*LoginResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *AuthServiceImpl) RequestPasswordReset(ctx context.Context, req PasswordResetRequest) error {
	//TODO implement me
	panic("implement me")
}

func (s *AuthServiceImpl) ResetPassword(ctx context.Context, token, newPassword string) error {
	//TODO implement me
	panic("implement me")
}

package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"had-service/constants"
	"log"
	"time"
)

type UserRepository interface {
	// User Management
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, username string) error

	//

	// Session Management
	CreateSession(ctx context.Context, session *UserSession) (*UserSession, error)
	GetSession(ctx context.Context, userId string) (*User, error)
	CleanSession(ctx context.Context, userId string) error
}

type UserPersistRepository struct {
	db    *gorm.DB
	redis redis.Client
}

func NewUserPersistRepository(db *gorm.DB, redis redis.Client) *UserPersistRepository {
	return &UserPersistRepository{
		db:    db,
		redis: redis,
	}
}

// Implementations

// Create User
func (r *UserPersistRepository) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Get User
func (r *UserPersistRepository) GetByID(ctx context.Context, id string) (*User, error) {
	user := &User{}
	if err := r.db.WithContext(ctx).First(user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *UserPersistRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	user := &User{}
	if err := r.db.WithContext(ctx).First(user, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *UserPersistRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	if err := r.db.WithContext(ctx).First(user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

// Update User
func (r *UserPersistRepository) Update(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete User
func (r *UserPersistRepository) Delete(ctx context.Context, username string) error {
	return r.db.WithContext(ctx).Delete(&User{}, "username = ?", username).Error
}

// session management
func (r *UserPersistRepository) CreateSession(ctx context.Context, session *UserSession) (*UserSession, error) {
	// Set creation time if not provided
	if session.CreatedAt.IsZero() {
		session.CreatedAt = time.Now()
	}
	session.UpdatedAt = time.Now()

	// Calculate TTL from expiration time
	ttl := session.ExpiresAt.Sub(time.Now()) * 2
	if ttl <= 0 {
		return nil, errors.New("session already expired")
	}

	// Store each field separately in the hash
	sessionKey := fmt.Sprintf(constants.UserSessionRK, session.UserID)
	fields := map[string]interface{}{
		"user_id":    session.UserID,
		"user_agent": session.UserAgent,
		"ip":         session.IP,
		"expires_at": session.ExpiresAt.Format(time.RFC3339),
		"created_at": session.CreatedAt.Format(time.RFC3339),
		"updated_at": session.UpdatedAt.Format(time.RFC3339),
	}

	if err := r.redis.HSet(ctx, sessionKey, fields).Err(); err != nil {
		return nil, fmt.Errorf("failed to store session in Redis: %w", err)
	}

	// Set expiration on the hash
	if err := r.redis.Expire(ctx, sessionKey, ttl).Err(); err != nil {
		return nil, fmt.Errorf("failed to set session expiration: %w", err)
	}

	// Use a single set for all active users
	const activeUsersKey = constants.ActiveUsersRK

	// Add this user to the active users set
	if err := r.redis.SAdd(ctx, activeUsersKey, session.UserID).Err(); err != nil {
		return nil, fmt.Errorf("failed to add user to active users set: %w", err)
	}

	return session, nil
}

func (r *UserPersistRepository) GetSession(ctx context.Context, userId string) (*UserSession, error) {
	sessionKey := fmt.Sprintf("session:%s", userId)

	// Get all fields from the hash
	sessionData, err := r.redis.HGetAll(ctx, sessionKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get session data: %w", err)
	}

	// Check if session exists
	if len(sessionData) == 0 {
		return nil, fmt.Errorf("session not found for user %s", userId)
	}

	// Reconstruct the session object from fields
	session := &UserSession{
		UserID:    userId,
		UserAgent: sessionData["user_agent"],
		IP:        sessionData["ip"],
	}

	// Parse time fields
	if expiresAt, err := time.Parse(time.RFC3339, sessionData["expires_at"]); err == nil {
		session.ExpiresAt = expiresAt
	}
	if createdAt, err := time.Parse(time.RFC3339, sessionData["created_at"]); err == nil {
		session.CreatedAt = createdAt
	}
	if updatedAt, err := time.Parse(time.RFC3339, sessionData["updated_at"]); err == nil {
		session.UpdatedAt = updatedAt
	}

	// Check if session has expired
	if time.Now().After(session.ExpiresAt) {
		// Delete the expired session
		if err := r.redis.Del(ctx, sessionKey).Err(); err != nil {
			log.Printf("Failed to delete expired session: %v", err)
		}
		return nil, fmt.Errorf("session has expired")
	}

	return session, nil
}

package user

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents the user model
type User struct {
	ID           string `json:"id" gorm:"primaryKey;type:uuid"`
	Username     string `json:"username" gorm:"uniqueIndex;not null"`
	Email        string `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string `json:"-" gorm:"not null"`
	LuckyNumber  string `json:"lucky_number" gorm:"not null"` // for password reset safety

	FirstName string `json:"first_name" gorm:"not null"`
	LastName  string `json:"last_name" gorm:"not null"`
	Location  string `json:"location" gorm:"not null"`

	ProfilePicture string   `json:"profile_picture" gorm:"not null"`
	CoverPicture   string   `json:"cover_picture" gorm:"not null"`
	Ft3Pictures    []string `json:"featured_3_pictures" gorm:"not null"`

	IsPremium  bool       `json:"is_premium" gorm:"default:false"`
	PremiumExp *time.Time `json:"premium_expiration"`

	IsPrivate  bool `json:"is_private" gorm:"default:false"`
	IsVerified bool `json:"is_verified" gorm:"default:false"`

	ProfileLikes int `json:"profile_likes" gorm:"default:0"`

	LastLoginAt *time.Time     `json:"last_login_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// UserSession represents a login session (stored in Redis)
type UserSession struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TokenHash string    `json:"-"`
	UserAgent string    `json:"user_agent"`
	IP        string    `json:"ip"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=30"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6,max=20"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`

	ProfilePicture string   `json:"profile_picture"`
	CoverPicture   string   `json:"cover_picture"`
	Ft3Pictures    []string `json:"featured_3_pictures" binding:"required"`

	IsPremium  bool `json:"is_premium" gorm:"default:false"`
	IsPrivate  bool `json:"is_private" gorm:"default:false"`
	IsVerified bool `json:"is_verified" gorm:"default:false"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	FirstName string `json:"first_name" binding:"min=3,max=30"`
	LastName  string `json:"last_name" binding:"min=3,max=20"`
	Email     string `json:"email" binding:"email"`

	ProfilePicture string   `json:"profile_picture"`
	CoverPicture   string   `json:"cover_picture"`
	Ft3Pictures    []string `json:"featured_3_pictures"`

	IsPremium  *bool      `json:"is_premium"`
	PremiumExp *time.Time `json:"premium_expiration"`
	IsPrivate  *bool      `json:"is_private"`
	IsVerified *bool      `json:"is_verified"`
}

// LoginRequest represents the login credentials
type LoginRequest struct {
	Email    string `json:"Email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// PasswordResetRequest represents a password reset request
type PasswordResetRequest struct {
	Email       string `json:"email" binding:"required,email"`
	LuckyNumber int    `json:"lucky_number" binding:"required"`
}

// PasswordResetConfirmRequest represents a password reset confirmation
type PasswordResetConfirmRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// UserProfile represents the user profile data
type UserProfile struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Location  string `json:"location"`

	ProfilePicture string   `json:"profile_picture"`
	CoverPicture   string   `json:"cover_picture"`
	Ft3Pictures    []string `json:"featured_3_pictures"`

	ProfileLikes int `json:"profile_likes"`

	IsPremium  bool       `json:"is_premium"`
	PremiumExp *time.Time `json:"premium_expiration,omitempty"`
	IsPrivate  bool       `json:"is_private"`
	IsVerified bool       `json:"is_verified"`

	PostsCount int `json:"posts_count" gorm:"-"`
}

// UserResponse is the standard user response
type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Location  string `json:"location"`

	ProfilePicture string   `json:"profile_picture"`
	CoverPicture   string   `json:"cover_picture"`
	Ft3Pictures    []string `json:"featured_3_pictures"`

	ProfileLikes int `json:"profile_likes"`

	IsPremium  bool       `json:"is_premium"`
	PremiumExp *time.Time `json:"premium_expiration,omitempty"`
	IsPrivate  bool       `json:"is_private"`
	IsVerified bool       `json:"is_verified"`
	CreatedAt  time.Time  `json:"created_at"`
}

// LoginResponse contains the token and user info
type LoginResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      UserResponse `json:"user"`
}

// HashPassword creates a hashed password from a plaintext password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a plaintext password against a hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ToResponse converts a User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:             u.ID,
		Username:       u.Username,
		Email:          u.Email,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Location:       u.Location,
		ProfilePicture: u.ProfilePicture,
		CoverPicture:   u.CoverPicture,
		Ft3Pictures:    u.Ft3Pictures,
		ProfileLikes:   u.ProfileLikes,
		IsPremium:      u.IsPremium,
		PremiumExp:     u.PremiumExp,
		IsPrivate:      u.IsPrivate,
		IsVerified:     u.IsVerified,
		CreatedAt:      u.CreatedAt,
	}
}

// ToProfile converts a User to UserProfile
func (u *User) ToProfile() UserProfile {
	return UserProfile{
		ID:             u.ID,
		Username:       u.Username,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Location:       u.Location,
		ProfilePicture: u.ProfilePicture,
		CoverPicture:   u.CoverPicture,
		Ft3Pictures:    u.Ft3Pictures,
		ProfileLikes:   u.ProfileLikes,
		IsPremium:      u.IsPremium,
		PremiumExp:     u.PremiumExp,
		IsPrivate:      u.IsPrivate,
		IsVerified:     u.IsVerified,
	}
}

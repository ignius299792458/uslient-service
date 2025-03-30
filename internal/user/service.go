package user

import (
	"context"
	"errors"
	"time"
)

type UserService interface {
	// User management
	Register(ctx context.Context, req CreateUserRequest) (*User, string, error)
	GetProfile(ctx context.Context, userID string) (*UserProfile, error)
	CheckEmailAndUsername(ctx context.Context, req CheckEmailAndUsernameRequest) (bool, error)
	UpdateProfile(ctx context.Context, userID string, req UpdateUserRequest) (*User, error)
	DeleteProfile(ctx context.Context, userID string) error
}

// ServiceImpl implements the Service interface
type UserServiceImpl struct {
	userRepo UserRepository
}

// NewService creates a new user service
func NewService(repo UserRepository) UserService {
	return &UserServiceImpl{
		userRepo: repo,
	}
}

// Register creates a new user account
func (s *UserServiceImpl) Register(ctx context.Context, req CreateUserRequest) (*User, string, error) {
	// check username
	existingUsername, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err == nil && existingUsername != nil {
		return nil, "", errors.New("username already exists")
	}

	// check email
	existingEmail, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingEmail != nil {
		return nil, "", errors.New("email already exists")
	}

	// Hash the password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, "", errors.New("failed to hash password")
	}

	// Create user
	user := &User{
		Username:       req.Username,
		Email:          req.Email,
		PasswordHash:   hashedPassword,
		LuckyNumber:    req.LuckyNumber,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		ProfilePicture: req.ProfilePicture,
		CoverPicture:   req.CoverPicture,
		Ft3Pictures:    req.Ft3Pictures,
		IsPrivate:      req.IsPrivate,
		IsVerified:     req.IsVerified,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Save to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, "", err
	}

	// Generate token
	//token, err := s.generateToken(user.ID)
	//if err != nil {
	//	return nil, "", errors.New("failed to generate token")
	//}
	return user, "", nil
}

func (s *UserServiceImpl) CheckEmailAndUsername(ctx context.Context, req CheckEmailAndUsernameRequest) (bool, error) {
	// check username
	existingUsername, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err == nil && existingUsername != nil {
		return true, errors.New("username already exists")
	}

	// check email
	existingEmail, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingEmail != nil {
		return true, errors.New("email already exists")
	}

	return false, nil
}

func (s *UserServiceImpl) GetProfile(ctx context.Context, userID string) (*UserProfile, error) {
	// Get the user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convert to profile
	profile := user.ToProfile()

	return &profile, nil
}

// UpdateProfile updates a user's profile information
func (s *UserServiceImpl) UpdateProfile(ctx context.Context, userID string, req UpdateUserRequest) (*User, error) {
	// Get the user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}

	if req.LastName != "" {
		user.LastName = req.LastName
	}

	if req.ProfilePicture != "" {
		user.ProfilePicture = req.ProfilePicture
	}

	if req.CoverPicture != "" {
		user.CoverPicture = req.CoverPicture
	}

	if req.Ft3Pictures != nil {
		user.Ft3Pictures = req.Ft3Pictures
	}

	user.UpdatedAt = time.Now()

	// Save to database
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Delete Profile
func (s *UserServiceImpl) DeleteProfile(ctx context.Context, userID string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user doesn't exists")
	}

	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return err
	}
	return nil
}

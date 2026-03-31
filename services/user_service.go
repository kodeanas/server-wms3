package services

import (
	"context"
	"wms/models"
	"wms/repositories"
)

// UserService defines user business logic
type UserService interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*models.User, error)
	ListUsers(ctx context.Context, limit, offset int) ([]models.User, int64, error)
	UpdateUser(ctx context.Context, id string, user *models.User) error
	DeleteUser(ctx context.Context, id string) error
}

type userService struct {
	repo repositories.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

// CreateUser creates a new user
func (s *userService) CreateUser(ctx context.Context, user *models.User) error {
	return s.repo.Create(ctx, user)
}

// GetUser gets user by ID
func (s *userService) GetUser(ctx context.Context, id string) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

// GetUserByEmail gets user by email
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

// GetUserByPhone gets user by phone
func (s *userService) GetUserByPhone(ctx context.Context, phone string) (*models.User, error) {
	return s.repo.GetByPhone(ctx, phone)
}

// ListUsers lists all users
func (s *userService) ListUsers(ctx context.Context, limit, offset int) ([]models.User, int64, error) {
	return s.repo.GetAll(ctx, limit, offset)
}

// UpdateUser updates user
func (s *userService) UpdateUser(ctx context.Context, id string, user *models.User) error {
	return s.repo.Update(ctx, id, user)
}

// DeleteUser deletes user
func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

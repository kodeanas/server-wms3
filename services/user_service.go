package services

import (
	"time"
	"wms/models"
	"wms/repositories"
)

type UserService interface {
	CreateUser(input CreateUserPayload) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	ListUsers() ([]models.User, error)
	UpdateUser(id string, input UpdateUserPayload) (*models.User, error)
	DeleteUser(id string) error
	UpdatePassword(id string, password string) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

type CreateUserPayload struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Status   bool   `json:"status"`
	Role     string `json:"role"`
}

type UpdateUserPayload struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Status bool   `json:"status"`
	Role   string `json:"role"`
}

func (s *userService) CreateUser(input CreateUserPayload) (*models.User, error) {
	user := &models.User{
		Name:      input.Name,
		Email:     input.Email,
		Phone:     input.Phone,
		Password:  input.Password,
		Status:    input.Status,
		Role:      input.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserByID(id string) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) ListUsers() ([]models.User, error) {
	return s.repo.List()
}

func (s *userService) UpdateUser(id string, input UpdateUserPayload) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	user.Name = input.Name
	user.Email = input.Email
	user.Phone = input.Phone
	user.Status = input.Status
	user.Role = input.Role
	user.UpdatedAt = time.Now()
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) DeleteUser(id string) error {
	return s.repo.Delete(id)
}

func (s *userService) UpdatePassword(id string, password string) error {
	return s.repo.UpdatePassword(id, password)
}

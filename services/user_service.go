package services

import (
	"time"
	"wms/models"
	"wms/repositories"
)

type UserService interface {
	CreateUser(input models.CreateUserPayload) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	ListUsersPaginated(page, limit int, search string) ([]models.User, int64, error)
	UpdateUser(id string, input models.UpdateUserPayload) (*models.User, error)
	DeleteUser(id string) error
	UpdatePassword(id string, password string) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(input models.CreateUserPayload) (*models.User, error) {
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

func (s *userService) ListUsersPaginated(page, limit int, search string) ([]models.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.ListPaginated(limit, offset, search)
}

func (s *userService) UpdateUser(id string, input models.UpdateUserPayload) (*models.User, error) {
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

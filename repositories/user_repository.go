package repositories

import (
	"context"
	"wms/config"
	"wms/models"
)

type userRepository struct{}

// NewUserRepository creates a new user repository
func NewUserRepository() UserRepository {
	return &userRepository{}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return config.DB.WithContext(ctx).Create(user).Error
}

// GetByID gets user by ID
func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := config.DB.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail gets user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := config.DB.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByPhone gets user by phone
func (r *userRepository) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User
	err := config.DB.WithContext(ctx).First(&user, "phone = ?", phone).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll gets all users with pagination
func (r *userRepository) GetAll(ctx context.Context, limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	err := config.DB.WithContext(ctx).Model(&models.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = config.DB.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Update updates user
func (r *userRepository) Update(ctx context.Context, id string, user *models.User) error {
	return config.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(user).Error
}

// Delete deletes user
func (r *userRepository) Delete(ctx context.Context, id string) error {
	return config.DB.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
}

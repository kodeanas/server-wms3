package repositories

import (
	"wms/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
	List() ([]models.User, error)
	Update(user *models.User) error
	Delete(id string) error
	UpdatePassword(id string, password string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) List() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Model(&models.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"name":       user.Name,
		"email":      user.Email,
		"phone":      user.Phone,
		"status":     user.Status,
		"role":       user.Role,
		"updated_at": user.UpdatedAt,
	}).Error
}

func (r *userRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.User{}).Error
}

func (r *userRepository) UpdatePassword(id string, password string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("password", password).Error
}

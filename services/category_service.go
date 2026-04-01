package services

import (
	"errors"
	"fmt"
	"strings"
	"wms/models"
	"wms/repositories"

	"github.com/google/uuid"
)

// CategoryService defines business logic for categories.
type CategoryService interface {
	CreateCategory(input CreateCategoryPayload) (*models.Category, error)
	ListCategories() ([]models.Category, error)
}

// CreateCategoryPayload request payload.
type CreateCategoryPayload struct {
	Name     string  `json:"name" binding:"required"`
	Slug     string  `json:"slug" binding:"required"`
	Discount int     `json:"discount"`
	MinPrice float64 `json:"min_price"`
	MaxPrice float64 `json:"max_price"`
}

type categoryService struct {
	repo repositories.CategoryRepository
}

// NewCategoryService constructor.
func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) CreateCategory(input CreateCategoryPayload) (*models.Category, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, errors.New("name is required")
	}
	if strings.TrimSpace(input.Slug) == "" {
		return nil, errors.New("slug is required")
	}
	if input.Discount < 0 || input.Discount > 100 {
		return nil, errors.New("discount must be between 0 and 100")
	}
	if input.MinPrice < 0 || input.MaxPrice < 0 {
		return nil, errors.New("price values must be non-negative")
	}
	if input.MaxPrice > 0 && input.MinPrice > input.MaxPrice {
		return nil, errors.New("min_price cannot be greater than max_price")
	}

	slug := strings.ToLower(strings.TrimSpace(input.Slug))
	slug = strings.ReplaceAll(slug, " ", "-")

	if _, err := s.repo.GetBySlug(slug); err == nil {
		return nil, fmt.Errorf("category slug '%s' already exists", slug)
	}

	category := &models.Category{
		ID:       uuid.New(),
		Name:     strings.TrimSpace(input.Name),
		Slug:     slug,
		Discount: input.Discount,
		MinPrice: input.MinPrice,
		MaxPrice: input.MaxPrice,
		Status:   "active",
	}

	if err := s.repo.Create(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *categoryService) ListCategories() ([]models.Category, error) {
	return s.repo.List()
}

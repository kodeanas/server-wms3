package services

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"wms/models"
	"wms/repositories"

	"github.com/google/uuid"
)

// CategoryService defines business logic for categories.
type CategoryService interface {
	CreateCategory(input CreateCategoryPayload) (*models.Category, error)
	ListCategories() ([]models.Category, error)
	GetCategoryByID(id string) (*models.Category, error)
	UpdateCategory(id string, input UpdateCategoryPayload) (*models.Category, error)
	DeleteCategory(id string) error
}

// UpdateCategoryPayload request payload for update.
type UpdateCategoryPayload struct {
	Name     string   `json:"name"`
	Slug     string   `json:"slug"`
	Discount *int     `json:"discount"`
	MinPrice *float64 `json:"min_price"`
	MaxPrice *float64 `json:"max_price"`
	Status   string   `json:"status"`
}

// DeleteCategory deletes a category by its ID.
func (s *categoryService) DeleteCategory(id string) error {
	return s.repo.Delete(id)
}

// UpdateCategory updates a category by its ID.
func (s *categoryService) UpdateCategory(id string, input UpdateCategoryPayload) (*models.Category, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		category.Name = input.Name
	}
	if input.Slug != "" {
		slug := s.generateSlugFromName(input.Slug)
		uniqueSlug, err := s.generateUniqueSlug(slug)
		if err != nil {
			return nil, err
		}
		category.Slug = uniqueSlug
	}
	if input.Discount != nil {
		category.Discount = input.Discount
	}
	if input.MinPrice != nil {
		p := models.Price(*input.MinPrice)
		category.MinPrice = &p
	}
	if input.MaxPrice != nil {
		p := models.Price(*input.MaxPrice)
		category.MaxPrice = &p
	}
	if input.Status != "" {
		category.Status = input.Status
	}
	if err := s.repo.Update(category); err != nil {
		return nil, err
	}
	return category, nil
}

// GetCategoryByID returns a category by its ID.
func (s *categoryService) GetCategoryByID(id string) (*models.Category, error) {
	return s.repo.GetByID(id)
}

// CreateCategoryPayload request payload.
type CreateCategoryPayload struct {
	Name     string   `json:"name" binding:"required"`
	Slug     string   `json:"slug"`
	Discount *int     `json:"discount"`
	MinPrice *float64 `json:"min_price"`
	MaxPrice *float64 `json:"max_price"`
}

type categoryService struct {
	repo repositories.CategoryRepository
}

// NewCategoryService constructor.
func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

// generateSlugFromName converts name to slug format (lowercase with dashes)
func (s *categoryService) generateSlugFromName(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

// generateUniqueSlug generates a unique slug by checking existing ones
func (s *categoryService) generateUniqueSlug(baseSlug string) (string, error) {
	_, err := s.repo.GetBySlug(baseSlug)
	if err != nil {
		// Slug doesn't exist, so we can use it.
		return baseSlug, nil
	}

	similarSlugs, err := s.repo.GetSlugLike(baseSlug)
	if err != nil {
		return "", err
	}

	// Find the highest suffix number among the similar slugs.
	maxNum := 0
	for _, cat := range similarSlugs {
		slug := cat.Slug

		// We only care about slugs that start with our base slug.
		if !strings.HasPrefix(slug, baseSlug) {
			continue
		}

		// Get the suffix part (e.g., "1", "2" from "otomotif1", "otomotif2").
		suffixStr := strings.TrimPrefix(slug, baseSlug)

		// The existence of the base slug itself implies that the next number will be at least 1.
		if suffixStr == "" {
			continue
		}

		// Try to convert the suffix to an integer.
		num, err := strconv.Atoi(suffixStr)
		if err != nil {
			// Suffix is not a number (e.g., "otomotif-baru"), so we ignore it.
			continue
		}

		// If this suffix number is the highest we've seen, remember it.
		if num > maxNum {
			maxNum = num
		}
	}

	newSlug := fmt.Sprintf("%s%d", baseSlug, maxNum+1)
	return newSlug, nil
}

func (s *categoryService) CreateCategory(input CreateCategoryPayload) (*models.Category, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, errors.New("name is required")
	}
	if input.Discount != nil && (*input.Discount < 0 || *input.Discount > 100) {
		return nil, errors.New("discount must be between 0 and 100")
	}
	if input.MinPrice != nil && *input.MinPrice < 0 {
		return nil, errors.New("price values must be non-negative")
	}
	if input.MaxPrice != nil && *input.MaxPrice < 0 {
		return nil, errors.New("price values must be non-negative")
	}
	if input.MaxPrice != nil && input.MinPrice != nil && *input.MaxPrice > 0 && *input.MinPrice > *input.MaxPrice {
		return nil, errors.New("min_price cannot be greater than max_price")
	}

	// Generate slug from name
	var slug string
	if strings.TrimSpace(input.Slug) != "" {
		// If slug provided, use it and make it unique if needed
		slug = s.generateSlugFromName(input.Slug)
	} else {
		// Generate slug from name
		slug = s.generateSlugFromName(input.Name)
	}

	// Make slug unique if duplicate
	uniqueSlug, err := s.generateUniqueSlug(slug)
	if err != nil {
		return nil, err
	}

	// Set default values if not provided
	discount := input.Discount
	if discount == nil {
		defaultDiscount := 0
		discount = &defaultDiscount
	}

	// Convert float64 prices to Price type
	var minPrice *models.Price
	if input.MinPrice != nil {
		p := models.Price(*input.MinPrice)
		minPrice = &p
	}

	var maxPrice *models.Price
	if input.MaxPrice != nil {
		p := models.Price(*input.MaxPrice)
		maxPrice = &p
	}

	category := &models.Category{
		ID:       uuid.New(),
		Name:     strings.TrimSpace(input.Name),
		Slug:     uniqueSlug,
		Discount: discount,
		MinPrice: minPrice,
		MaxPrice: maxPrice,
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

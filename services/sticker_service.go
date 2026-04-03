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

// StickerService defines business logic for stickers.
type StickerService interface {
	CreateSticker(input CreateStickerPayload) (*models.Sticker, error)
	GetStickerBySlug(slug string) (*models.Sticker, error)
	GetStickerByID(id string) (*models.Sticker, error)
	ListStickers() ([]models.Sticker, error)
	UpdateSticker(id string, input UpdateStickerPayload) (*models.Sticker, error)
	DeleteSticker(id string) error
}

// CreateStickerPayload request payload.
type CreateStickerPayload struct {
	CodeHex    string   `json:"code_hex" binding:"required"`
	Name       string   `json:"name" binding:"required"`
	Slug       string   `json:"slug"`
	Type       string   `json:"type"`
	FixedPrice *int     `json:"fixed_price"`
	MinPrice   *float64 `json:"min_price"`
	MaxPrice   *float64 `json:"max_price"`
}

// UpdateStickerPayload request payload for update.
type UpdateStickerPayload struct {
	CodeHex    string   `json:"code_hex"`
	Name       string   `json:"name"`
	Slug       string   `json:"slug"`
	Type       string   `json:"type"`
	FixedPrice *int     `json:"fixed_price"`
	MinPrice   *float64 `json:"min_price"`
	MaxPrice   *float64 `json:"max_price"`
}

type stickerService struct {
	repo repositories.StickerRepository
}

// NewStickerService constructor.
func NewStickerService(repo repositories.StickerRepository) StickerService {
	return &stickerService{repo: repo}
}

// generateSlugFromName converts name to slug format (lowercase with dashes)
func (s *stickerService) generateSlugFromName(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

// generateUniqueSlug generates a unique slug by checking existing ones
func (s *stickerService) generateUniqueSlug(baseSlug string) (string, error) {
	// Check if base slug exists
	_, err := s.repo.GetBySlug(baseSlug)
	if err != nil {
		// Slug doesn't exist, return it
		return baseSlug, nil
	}

	// Get all similar slugs
	similarSlugs, err := s.repo.GetSlugLike(baseSlug)
	if err != nil {
		return "", err
	}

	// Extract numbers from similar slugs and find the highest number
	maxNum := 0
	for _, sticker := range similarSlugs {
		if sticker.Slug == baseSlug {
			maxNum = 0
			continue
		}
		// Try to extract number from slug like "sticker1", "sticker2"
		parts := strings.Split(sticker.Slug, baseSlug)
		if len(parts) > 1 && parts[1] != "" {
			if num, err := strconv.Atoi(parts[1]); err == nil && num > maxNum {
				maxNum = num
			}
		}
	}

	// Generate new slug with incremented number
	newSlug := fmt.Sprintf("%s%d", baseSlug, maxNum+1)
	return newSlug, nil
}

func (s *stickerService) CreateSticker(input CreateStickerPayload) (*models.Sticker, error) {
	// Validation
	if strings.TrimSpace(input.CodeHex) == "" {
		return nil, errors.New("code_hex is required")
	}
	if strings.TrimSpace(input.Name) == "" {
		return nil, errors.New("name is required")
	}
	if input.MinPrice != nil && input.MaxPrice != nil && *input.MinPrice < 0 || (*input.MaxPrice < 0) {
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

	sticker := &models.Sticker{
		ID:         uuid.New(),
		CodeHex:    input.CodeHex,
		Name:       input.Name,
		Slug:       uniqueSlug,
		Type:       input.Type,
		FixedPrice: input.FixedPrice,
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
	}

	if err := s.repo.Create(sticker); err != nil {
		return nil, err
	}

	return sticker, nil
}

func (s *stickerService) GetStickerBySlug(slug string) (*models.Sticker, error) {
	return s.repo.GetBySlug(slug)
}

func (s *stickerService) GetStickerByID(id string) (*models.Sticker, error) {
	return s.repo.GetByID(id)
}

func (s *stickerService) ListStickers() ([]models.Sticker, error) {
	return s.repo.List()
}

func (s *stickerService) UpdateSticker(id string, input UpdateStickerPayload) (*models.Sticker, error) {
	sticker, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("sticker not found")
	}

	// Validate input if provided
	if input.MinPrice != nil && *input.MinPrice < 0 {
		return nil, errors.New("price values must be non-negative")
	}
	if input.MaxPrice != nil && *input.MaxPrice < 0 {
		return nil, errors.New("price values must be non-negative")
	}
	if input.MaxPrice != nil && input.MinPrice != nil && *input.MaxPrice > 0 && *input.MinPrice > *input.MaxPrice {
		return nil, errors.New("min_price cannot be greater than max_price")
	}

	// Update fields if provided
	if input.CodeHex != "" {
		sticker.CodeHex = input.CodeHex
	}
	if input.Name != "" {
		sticker.Name = input.Name
		// Jika nama berubah, update slug
		newSlug := s.generateSlugFromName(input.Name)
		if newSlug != sticker.Slug {
			existing, _ := s.repo.GetBySlug(newSlug)
			if existing != nil && existing.ID != sticker.ID {
				uniqueSlug, err := s.generateUniqueSlug(newSlug)
				if err != nil {
					return nil, err
				}
				sticker.Slug = uniqueSlug
			} else {
				sticker.Slug = newSlug
			}
		}
	}
	if input.Slug != "" {
		// Generate slug from input and check uniqueness
		newSlug := s.generateSlugFromName(input.Slug)
		if newSlug != sticker.Slug {
			// Slug is different, check if it's unique
			existing, _ := s.repo.GetBySlug(newSlug)
			if existing != nil && existing.ID != sticker.ID {
				// Slug exists for another record, generate unique slug
				uniqueSlug, err := s.generateUniqueSlug(newSlug)
				if err != nil {
					return nil, err
				}
				sticker.Slug = uniqueSlug
			} else {
				sticker.Slug = newSlug
			}
		}
	}
	if input.Type != "" {
		sticker.Type = input.Type
	}
	if input.FixedPrice != nil {
		sticker.FixedPrice = input.FixedPrice
	}
	if input.MinPrice != nil {
		p := models.Price(*input.MinPrice)
		sticker.MinPrice = &p
	}
	if input.MaxPrice != nil {
		p := models.Price(*input.MaxPrice)
		sticker.MaxPrice = &p
	}

	if err := s.repo.Update(sticker); err != nil {
		return nil, err
	}

	return sticker, nil
}

func (s *stickerService) DeleteSticker(id string) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("sticker not found")
	}
	return s.repo.Delete(id)
}

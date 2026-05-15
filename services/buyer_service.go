package services

import (
	"wms/models"
	"wms/repositories"

	"fmt"
	"time"

	"github.com/google/uuid"
)

type BuyerWithClass struct {
	models.Buyer
	Class *ClassWithDecimal `json:"class,omitempty"`
}

type ClassWithDecimal struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	MinOrder            int    `json:"min_order"`
	Disc                int    `json:"disc"`
	MinTransactionValue string `json:"min_transaction_value"`
	Week                int    `json:"week"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}

// BuyerService defines business logic for buyers.
type BuyerService interface {
	CreateBuyer(input models.CreateBuyerPayload) (*models.Buyer, error)
	GetBuyerByID(id string) (*models.Buyer, error)
	GetBuyerDetail(id string) (*BuyerWithClass, error)
	ListBuyersPaginated(page, limit int, search string) ([]BuyerWithClass, int64, error)
	ListBuyersByClass(classID string, page, limit int, search string) ([]BuyerWithClass, *ClassWithDecimal, int64, error)
	UpdateBuyer(id string, input models.UpdateBuyerPayload) (*models.Buyer, error)
	DeleteBuyer(id string) error
}

type buyerService struct {
	repo      repositories.BuyerRepository
	classRepo repositories.ClassRepository
}

// NewBuyerService constructor.
func NewBuyerService(repo repositories.BuyerRepository, classRepo repositories.ClassRepository) BuyerService {
	return &buyerService{repo: repo, classRepo: classRepo}
}

func (s *buyerService) CreateBuyer(input models.CreateBuyerPayload) (*models.Buyer, error) {
	buyer := &models.Buyer{
		ID:      uuid.New(),
		Name:    input.Name,
		Email:   input.Email,
		Phone:   input.Phone,
		ClassID: input.ClassID,
		Address: input.Address,
	}
	if err := s.repo.Create(buyer); err != nil {
		return nil, err
	}
	return buyer, nil
}

func (s *buyerService) GetBuyerByID(id string) (*models.Buyer, error) {
	return s.repo.GetByID(id)
}

func (s *buyerService) GetBuyerDetail(id string) (*BuyerWithClass, error) {
	buyer, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	var class *ClassWithDecimal
	if buyer.ClassID != "" {
		if c, err := s.classRepo.GetByID(buyer.ClassID); err == nil {
			class = classToWithDecimal(c)
		}
	}
	return &BuyerWithClass{Buyer: *buyer, Class: class}, nil
}

func classToWithDecimal(c *models.Class) *ClassWithDecimal {
	if c == nil {
		return nil
	}
	return &ClassWithDecimal{
		ID:                  c.ID.String(),
		Name:                c.Name,
		MinOrder:            c.MinOrder,
		Disc:                c.Disc,
		MinTransactionValue: fmt.Sprintf("%.2f", c.MinTransactionValue),
		Week:                c.Week,
		CreatedAt:           c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           c.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *buyerService) ListBuyersPaginated(page, limit int, search string) ([]BuyerWithClass, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	buyers, total, err := s.repo.ListPaginated(limit, offset, search)
	if err != nil {
		return nil, 0, err
	}
	result := make([]BuyerWithClass, 0, len(buyers))
	for _, b := range buyers {
		var class *ClassWithDecimal
		if b.ClassID != "" {
			c, err := s.classRepo.GetByID(b.ClassID)
			if err == nil {
				class = classToWithDecimal(c)
			}
		}
		result = append(result, BuyerWithClass{Buyer: b, Class: class})
	}
	return result, total, nil
}

func (s *buyerService) ListBuyersByClass(classID string, page, limit int, search string) ([]BuyerWithClass, *ClassWithDecimal, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Validate class exists
	classModel, err := s.classRepo.GetByID(classID)
	if err != nil {
		return nil, nil, 0, err
	}
	classInfo := classToWithDecimal(classModel)

	buyers, total, err := s.repo.ListByClassID(classID, limit, offset, search)
	if err != nil {
		return nil, nil, 0, err
	}
	result := make([]BuyerWithClass, 0, len(buyers))
	for _, b := range buyers {
		result = append(result, BuyerWithClass{Buyer: b, Class: classInfo})
	}
	return result, classInfo, total, nil
}

func (s *buyerService) UpdateBuyer(id string, input models.UpdateBuyerPayload) (*models.Buyer, error) {
	buyer, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		buyer.Name = input.Name
	}
	if input.Email != "" {
		buyer.Email = input.Email
	}
	if input.Phone != "" {
		buyer.Phone = input.Phone
	}
	if input.ClassID != "" {
		buyer.ClassID = input.ClassID
	}
	if input.Address != "" {
		buyer.Address = input.Address
	}
	if err := s.repo.Update(buyer); err != nil {
		return nil, err
	}
	return buyer, nil
}

func (s *buyerService) DeleteBuyer(id string) error {
	return s.repo.Delete(id)
}

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
	ListBuyers() ([]BuyerWithClass, error)
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

func (s *buyerService) ListBuyers() ([]BuyerWithClass, error) {
	buyers, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	var result []BuyerWithClass
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
	return result, nil
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

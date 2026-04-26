package services

import (
	"wms/models"
	"wms/repositories"

	"github.com/google/uuid"
)

type ClassService interface {
	CreateClass(input CreateClassPayload) (*models.Class, error)
	GetClassByID(id string) (*models.Class, error)
	ListClasses() ([]models.Class, error)
	UpdateClass(id string, input UpdateClassPayload) (*models.Class, error)
	DeleteClass(id string) error
	MoveUp(id string) error
	MoveDown(id string) error
}

type CreateClassPayload struct {
	Name                string  `json:"name" binding:"required"`
	MinOrder            int     `json:"min_order" binding:"required"`
	Disc                int     `json:"disc"`
	MinTransactionValue float64 `json:"min_transaction_value" binding:"required"`
	Week                int     `json:"week"`
	Iteration           int     `json:"iteration"`
	Status              string  `json:"status"`
}

type UpdateClassPayload struct {
	Name                string  `json:"name"`
	MinOrder            int     `json:"min_order"`
	Disc                int     `json:"disc"`
	MinTransactionValue float64 `json:"min_transaction_value"`
	Week                int     `json:"week"`
	Iteration           int     `json:"iteration"`
	Status              string  `json:"status"`
}

type classService struct {
	repo repositories.ClassRepository
}

func NewClassService(repo repositories.ClassRepository) ClassService {
	return &classService{repo: repo}
}

func (s *classService) CreateClass(input CreateClassPayload) (*models.Class, error) {
       maxIter, err := s.repo.GetMaxIteration()
       if err != nil {
	       return nil, err
       }
       class := &models.Class{
	       ID:                  uuid.New(),
	       Name:                input.Name,
	       MinOrder:            input.MinOrder,
	       Disc:                input.Disc,
	       MinTransactionValue: input.MinTransactionValue,
	       Week:                input.Week,
	       Iteration:           maxIter + 1,
	       Status:              input.Status,
       }
       if err := s.repo.Create(class); err != nil {
	       return nil, err
       }
       return class, nil
}

func (s *classService) GetClassByID(id string) (*models.Class, error) {
	return s.repo.GetByID(id)
}

func (s *classService) ListClasses() ([]models.Class, error) {
	return s.repo.List()
}

func (s *classService) UpdateClass(id string, input UpdateClassPayload) (*models.Class, error) {
	class, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		class.Name = input.Name
	}
	if input.MinOrder != 0 {
		class.MinOrder = input.MinOrder
	}
	if input.Disc != 0 {
		class.Disc = input.Disc
	}
	if input.MinTransactionValue != 0 {
		class.MinTransactionValue = input.MinTransactionValue
	}
	if input.Week != 0 {
		class.Week = input.Week
	}
	if err := s.repo.Update(class); err != nil {
		return nil, err
	}
	return class, nil
}

func (s *classService) DeleteClass(id string) error {
       class, err := s.repo.GetByID(id)
       if err != nil {
	       return err
       }
       deletedIter := class.Iteration
       if err := s.repo.Delete(id); err != nil {
	       return err
       }
       return s.repo.DecrementIterationAbove(deletedIter)
}

func (s *classService) MoveUp(id string) error {
       class, err := s.repo.GetByID(id)
       if err != nil {
	       return err
       }
       prev, err := s.repo.GetPrevByIteration(class.Iteration)
       if err != nil {
	       return err // Sudah paling atas atau error lain
       }
       return s.repo.SwapIteration(class, prev)
}

func (s *classService) MoveDown(id string) error {
       class, err := s.repo.GetByID(id)
       if err != nil {
	       return err
       }
       next, err := s.repo.GetNextByIteration(class.Iteration)
       if err != nil {
	       return err // Sudah paling bawah atau error lain
       }
       return s.repo.SwapIteration(class, next)
}

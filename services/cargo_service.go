package services

import (
	"fmt"
	"time"
	dto "wms/dto/response"
	"wms/models"
	"wms/repositories"

	"github.com/google/uuid"
)

type CargoService interface {
	CreateCargo(userID string) (*models.Cargo, error)
	ListCargosPaginated(page, limit int, search string) ([]dto.CargoListResponse, int64, error)
	GetCargoDetail(cargoID string) (*dto.CargoDetailResponse, error)
	ScanBag(cargoID, bagCode string) error
	FinishCargo(cargoID string) error
}

type cargoService struct {
	cargoRepo repositories.CargoRepository
	bagRepo   repositories.BagRepository
}

func NewCargoService(cargoRepo repositories.CargoRepository, bagRepo repositories.BagRepository) CargoService {
	return &cargoService{
		cargoRepo: cargoRepo,
		bagRepo:   bagRepo,
	}
}

func (s *cargoService) CreateCargo(userID string) (*models.Cargo, error) {
	cargo := &models.Cargo{
		Code:      fmt.Sprintf("CARGO-%d", time.Now().UnixNano()),
		Status:    "open",
		IsSale:    false,
		IsOnline:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if userID != "" {
		if uid, err := uuid.Parse(userID); err == nil && uid != uuid.Nil {
			cargo.UserID = &uid
		}
	}
	if err := s.cargoRepo.Create(cargo); err != nil {
		return nil, err
	}
	return cargo, nil
}

func (s *cargoService) ListCargosPaginated(page, limit int, search string) ([]dto.CargoListResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	cargos, total, err := s.cargoRepo.FindAllPaginated(limit, offset, search)
	if err != nil {
		return nil, 0, err
	}

	resp := make([]dto.CargoListResponse, 0, len(cargos))
	for _, c := range cargos {
		resp = append(resp, dto.CargoListResponse{
			ID:        c.ID.String(),
			Code:      c.Code,
			Status:    c.Status,
			IsSale:    c.IsSale,
			IsOnline:  c.IsOnline,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		})
	}
	return resp, total, nil
}

func (s *cargoService) GetCargoDetail(cargoID string) (*dto.CargoDetailResponse, error) {
	cargo, err := s.cargoRepo.FindByID(cargoID)
	if err != nil {
		return nil, err
	}

	// Ambil semua bag yang punya cargo_id ini
	bags, err := s.bagRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var cargoDetail dto.CargoDetailResponse
	cargoDetail.ID = cargo.ID.String()
	cargoDetail.Code = cargo.Code
	cargoDetail.Status = cargo.Status
	cargoDetail.IsSale = cargo.IsSale
	cargoDetail.IsOnline = cargo.IsOnline
	cargoDetail.CreatedAt = cargo.CreatedAt.Format(time.RFC3339)
	cargoDetail.TotalBag = 0
	cargoDetail.TotalItem = 0
	cargoDetail.TotalPrice = 0

	// Filter bag yang sesuai dengan cargo_id
	for _, b := range bags {
		if b.CargoID != nil && b.CargoID.String() == cargoID {
			cargoDetail.TotalBag++
			cargoDetail.Bags = append(cargoDetail.Bags, dto.BagItemResponse{
				ID:      b.ID.String(),
				Code:    b.Code,
				Type:    b.Type,
				IsMoved: b.IsMoved,
			})
		}
	}

	return &cargoDetail, nil
}

func (s *cargoService) ScanBag(cargoID, bagCode string) error {
	// Validasi cargo exists
	cargo, err := s.cargoRepo.FindByID(cargoID)
	if err != nil {
		return fmt.Errorf("cargo tidak ditemukan")
	}

	// Cari bag berdasarkan code
	bags, err := s.bagRepo.FindAll()
	if err != nil {
		return err
	}

	var targetBag *models.Bag
	for _, b := range bags {
		if b.Code == bagCode {
			targetBag = &b
			break
		}
	}

	if targetBag == nil {
		return fmt.Errorf("bag dengan code %s tidak ditemukan", bagCode)
	}

	// Validasi bag: harus is_moved=false dan belum punya cargo
	if targetBag.IsMoved {
		return fmt.Errorf("bag sudah moved, tidak dapat di-scan")
	}
	if targetBag.CargoID != nil {
		return fmt.Errorf("bag sudah diassign ke cargo lain")
	}

	// Validasi status cargo harus "open"
	if cargo.Status != "open" {
		return fmt.Errorf("cargo status bukan open")
	}

	// Assign bag ke cargo
	cargoUUID := cargo.ID
	targetBag.CargoID = &cargoUUID
	if err := s.bagRepo.UpdateCargoID(targetBag.ID.String(), cargoID); err != nil {
		return err
	}

	return nil
}

func (s *cargoService) FinishCargo(cargoID string) error {
	// Validasi cargo exists dan status open
	cargo, err := s.cargoRepo.FindByID(cargoID)
	if err != nil {
		return fmt.Errorf("cargo tidak ditemukan")
	}

	if cargo.Status != "open" {
		return fmt.Errorf("cargo status bukan open")
	}

	// Set status menjadi lock
	if err := s.cargoRepo.SetStatus(cargoID, "lock"); err != nil {
		return err
	}

	if err := s.cargoRepo.SetIsSale(cargoID, true); err != nil {
		return err
	}

	// Set is_moved = true untuk semua bag dalam cargo
	if err := s.cargoRepo.SetIsMoved(cargoID); err != nil {
		return err
	}

	return nil
}

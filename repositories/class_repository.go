package repositories

import (
	"database/sql"
	"wms/models"

	"gorm.io/gorm"
)

// ClassRepository defines interface for class CRUD.
type ClassRepository interface {
	Create(class *models.Class) error
	GetByID(id string) (*models.Class, error)
	List() ([]models.Class, error)
	Update(class *models.Class) error
	Delete(id string) error
	GetMaxIteration() (int, error)
	GetByIteration(iteration int) (*models.Class, error)
	GetPrevByIteration(iteration int) (*models.Class, error)
	GetNextByIteration(iteration int) (*models.Class, error)
	SwapIteration(class1, class2 *models.Class) error
	DecrementIterationAbove(deletedIteration int) error
}

func (r *classRepository) GetMaxIteration() (int, error) {
	var max sql.NullInt64
	err := r.db.Model(&models.Class{}).Select("MAX(iteration)").Scan(&max).Error
	if err != nil {
		return 0, err
	}
	if !max.Valid {
		return 0, nil
	}
	return int(max.Int64), nil
}

func (r *classRepository) GetByIteration(iteration int) (*models.Class, error) {
	var class models.Class
	if err := r.db.Where("iteration = ?", iteration).First(&class).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

func (r *classRepository) GetPrevByIteration(iteration int) (*models.Class, error) {
	var class models.Class
	err := r.db.Where("iteration < ?", iteration).Order("iteration DESC").First(&class).Error
	if err != nil {
		return nil, err
	}
	return &class, nil
}

func (r *classRepository) GetNextByIteration(iteration int) (*models.Class, error) {
	var class models.Class
	err := r.db.Where("iteration > ?", iteration).Order("iteration ASC").First(&class).Error
	if err != nil {
		return nil, err
	}
	return &class, nil
}

func (r *classRepository) SwapIteration(class1, class2 *models.Class) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		tmp := class1.Iteration
		if err := tx.Model(class1).Update("iteration", class2.Iteration).Error; err != nil {
			return err
		}
		if err := tx.Model(class2).Update("iteration", tmp).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *classRepository) DecrementIterationAbove(deletedIteration int) error {
	return r.db.Model(&models.Class{}).
		Where("iteration > ?", deletedIteration).
		Update("iteration", gorm.Expr("iteration - 1")).Error
}

type classRepository struct {
	db *gorm.DB
}

func NewClassRepository(db *gorm.DB) ClassRepository {
	return &classRepository{db: db}
}

func (r *classRepository) Create(class *models.Class) error {
	return r.db.Create(class).Error
}

func (r *classRepository) GetByID(id string) (*models.Class, error) {
	var class models.Class
	if err := r.db.Where("id = ?", id).First(&class).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

func (r *classRepository) List() ([]models.Class, error) {
       var classes []models.Class
       if err := r.db.Order("iteration ASC").Find(&classes).Error; err != nil {
	       return nil, err
       }
       return classes, nil
}

func (r *classRepository) Update(class *models.Class) error {
	return r.db.Save(class).Error
}

func (r *classRepository) Delete(id string) error {
	return r.db.Delete(&models.Class{}, "id = ?", id).Error
}

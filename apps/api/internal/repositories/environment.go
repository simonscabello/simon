package repositories

import (
	"errors"

	"gorm.io/gorm"

	"simon/apps/api/internal/models"
)

type EnvironmentRepository struct {
	db *gorm.DB
}

func NewEnvironmentRepository(db *gorm.DB) *EnvironmentRepository {
	return &EnvironmentRepository{db: db}
}

func (r *EnvironmentRepository) Create(e *models.Environment) error {
	return r.db.Create(e).Error
}

func (r *EnvironmentRepository) ListByUserID(userID uint) ([]models.Environment, error) {
	var list []models.Environment
	err := r.db.Where("user_id = ?", userID).Order("created_at asc").Find(&list).Error
	return list, err
}

func (r *EnvironmentRepository) FindByIDForUser(id, userID uint) (*models.Environment, error) {
	var e models.Environment
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (r *EnvironmentRepository) Update(e *models.Environment) error {
	return r.db.Save(e).Error
}

func (r *EnvironmentRepository) Delete(id, userID uint) error {
	res := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Environment{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

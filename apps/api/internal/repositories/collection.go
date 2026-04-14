package repositories

import (
	"errors"

	"gorm.io/gorm"

	"simon/apps/api/internal/models"
)

type CollectionRepository struct {
	db *gorm.DB
}

func NewCollectionRepository(db *gorm.DB) *CollectionRepository {
	return &CollectionRepository{db: db}
}

func (r *CollectionRepository) Create(c *models.Collection) error {
	return r.db.Create(c).Error
}

func (r *CollectionRepository) ListByUserID(userID uint) ([]models.Collection, error) {
	var list []models.Collection
	err := r.db.Where("user_id = ?", userID).Order("created_at asc").Find(&list).Error
	return list, err
}

func (r *CollectionRepository) FindByIDForUser(id, userID uint) (*models.Collection, error) {
	var col models.Collection
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&col).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &col, nil
}

func (r *CollectionRepository) Update(c *models.Collection) error {
	return r.db.Save(c).Error
}

func (r *CollectionRepository) Delete(id, userID uint) error {
	res := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Collection{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

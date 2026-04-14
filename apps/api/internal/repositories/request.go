package repositories

import (
	"errors"

	"gorm.io/gorm"

	"simon/apps/api/internal/models"
)

type RequestRepository struct {
	db *gorm.DB
}

func NewRequestRepository(db *gorm.DB) *RequestRepository {
	return &RequestRepository{db: db}
}

func (r *RequestRepository) Create(req *models.Request) error {
	return r.db.Create(req).Error
}

func (r *RequestRepository) ListByCollectionID(collectionID uint) ([]models.Request, error) {
	var list []models.Request
	err := r.db.Where("collection_id = ?", collectionID).Order("created_at asc").Find(&list).Error
	return list, err
}

func (r *RequestRepository) FindByIDForUser(id, userID uint) (*models.Request, error) {
	var req models.Request
	err := r.db.
		Joins("JOIN collections ON collections.id = requests.collection_id AND collections.user_id = ?", userID).
		Where("requests.id = ?", id).
		First(&req).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &req, nil
}

func (r *RequestRepository) Update(req *models.Request) error {
	return r.db.Save(req).Error
}

func (r *RequestRepository) Delete(id, userID uint) error {
	req, err := r.FindByIDForUser(id, userID)
	if err != nil {
		return err
	}
	if req == nil {
		return gorm.ErrRecordNotFound
	}
	return r.db.Delete(&models.Request{}, req.ID).Error
}

func (r *RequestRepository) CollectionOwnedByUser(collectionID, userID uint) (bool, error) {
	var n int64
	err := r.db.Model(&models.Collection{}).Where("id = ? AND user_id = ?", collectionID, userID).Count(&n).Error
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

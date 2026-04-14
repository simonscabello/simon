package repositories

import (
	"errors"

	"gorm.io/gorm"

	"simon/apps/api/internal/models"
)

type RequestHeaderRepository struct {
	db *gorm.DB
}

func NewRequestHeaderRepository(db *gorm.DB) *RequestHeaderRepository {
	return &RequestHeaderRepository{db: db}
}

func (r *RequestHeaderRepository) Create(h *models.RequestHeader) error {
	return r.db.Create(h).Error
}

func (r *RequestHeaderRepository) ListByRequestID(requestID uint) ([]models.RequestHeader, error) {
	var list []models.RequestHeader
	err := r.db.Where("request_id = ?", requestID).Order("created_at asc").Find(&list).Error
	return list, err
}

func (r *RequestHeaderRepository) FindByIDForUser(id, userID uint) (*models.RequestHeader, error) {
	var h models.RequestHeader
	err := r.db.
		Joins("JOIN requests ON requests.id = request_headers.request_id").
		Joins("JOIN collections ON collections.id = requests.collection_id AND collections.user_id = ?", userID).
		Where("request_headers.id = ?", id).
		First(&h).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &h, nil
}

func (r *RequestHeaderRepository) Update(h *models.RequestHeader) error {
	return r.db.Save(h).Error
}

func (r *RequestHeaderRepository) Delete(id, userID uint) error {
	h, err := r.FindByIDForUser(id, userID)
	if err != nil {
		return err
	}
	if h == nil {
		return gorm.ErrRecordNotFound
	}
	return r.db.Delete(&models.RequestHeader{}, h.ID).Error
}

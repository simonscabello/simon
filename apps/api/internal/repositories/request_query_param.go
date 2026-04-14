package repositories

import (
	"errors"

	"gorm.io/gorm"

	"simon/apps/api/internal/models"
)

type RequestQueryParamRepository struct {
	db *gorm.DB
}

func NewRequestQueryParamRepository(db *gorm.DB) *RequestQueryParamRepository {
	return &RequestQueryParamRepository{db: db}
}

func (r *RequestQueryParamRepository) Create(q *models.RequestQueryParam) error {
	return r.db.Create(q).Error
}

func (r *RequestQueryParamRepository) ListByRequestID(requestID uint) ([]models.RequestQueryParam, error) {
	var list []models.RequestQueryParam
	err := r.db.Where("request_id = ?", requestID).Order("created_at asc").Find(&list).Error
	return list, err
}

func (r *RequestQueryParamRepository) FindByIDForUser(id, userID uint) (*models.RequestQueryParam, error) {
	var q models.RequestQueryParam
	err := r.db.
		Joins("JOIN requests ON requests.id = request_query_params.request_id").
		Joins("JOIN collections ON collections.id = requests.collection_id AND collections.user_id = ?", userID).
		Where("request_query_params.id = ?", id).
		First(&q).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &q, nil
}

func (r *RequestQueryParamRepository) Update(q *models.RequestQueryParam) error {
	return r.db.Save(q).Error
}

func (r *RequestQueryParamRepository) Delete(id, userID uint) error {
	row, err := r.FindByIDForUser(id, userID)
	if err != nil {
		return err
	}
	if row == nil {
		return gorm.ErrRecordNotFound
	}
	return r.db.Delete(&models.RequestQueryParam{}, row.ID).Error
}

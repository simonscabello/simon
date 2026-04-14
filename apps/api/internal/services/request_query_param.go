package services

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"simon/apps/api/internal/models"
	"simon/apps/api/internal/repositories"
)

type RequestQueryParamService struct {
	requests *repositories.RequestRepository
	params   *repositories.RequestQueryParamRepository
}

func NewRequestQueryParamService(db *gorm.DB) *RequestQueryParamService {
	return &RequestQueryParamService{
		requests: repositories.NewRequestRepository(db),
		params:   repositories.NewRequestQueryParamRepository(db),
	}
}

func (s *RequestQueryParamService) Create(userID, requestID uint, key, value string, enabled *bool) (*models.RequestQueryParam, error) {
	req, err := s.requests.FindByIDForUser(requestID, userID)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, ErrNotFound
	}
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, ErrInvalidKey
	}
	q := &models.RequestQueryParam{
		RequestID: requestID,
		Key:       key,
		Value:     value,
		Enabled:   enabledOrDefault(enabled),
	}
	if err := s.params.Create(q); err != nil {
		return nil, err
	}
	return q, nil
}

func (s *RequestQueryParamService) List(userID, requestID uint) ([]models.RequestQueryParam, error) {
	req, err := s.requests.FindByIDForUser(requestID, userID)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, ErrNotFound
	}
	return s.params.ListByRequestID(requestID)
}

func (s *RequestQueryParamService) Update(userID, id uint, key, value string, enabled *bool) (*models.RequestQueryParam, error) {
	q, err := s.params.FindByIDForUser(id, userID)
	if err != nil {
		return nil, err
	}
	if q == nil {
		return nil, ErrNotFound
	}
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, ErrInvalidKey
	}
	q.Key = key
	q.Value = value
	if enabled != nil {
		q.Enabled = *enabled
	}
	if err := s.params.Update(q); err != nil {
		return nil, err
	}
	return q, nil
}

func (s *RequestQueryParamService) Delete(userID, id uint) error {
	err := s.params.Delete(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

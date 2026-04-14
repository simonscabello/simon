package services

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"simon/apps/api/internal/models"
	"simon/apps/api/internal/repositories"
)

type RequestHeaderService struct {
	requests *repositories.RequestRepository
	headers  *repositories.RequestHeaderRepository
}

func NewRequestHeaderService(db *gorm.DB) *RequestHeaderService {
	return &RequestHeaderService{
		requests: repositories.NewRequestRepository(db),
		headers:  repositories.NewRequestHeaderRepository(db),
	}
}

func enabledOrDefault(p *bool) bool {
	if p == nil {
		return true
	}
	return *p
}

func (s *RequestHeaderService) Create(userID, requestID uint, key, value string, enabled *bool) (*models.RequestHeader, error) {
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
	h := &models.RequestHeader{
		RequestID: requestID,
		Key:       key,
		Value:     value,
		Enabled:   enabledOrDefault(enabled),
	}
	if err := s.headers.Create(h); err != nil {
		return nil, err
	}
	return h, nil
}

func (s *RequestHeaderService) List(userID, requestID uint) ([]models.RequestHeader, error) {
	req, err := s.requests.FindByIDForUser(requestID, userID)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, ErrNotFound
	}
	return s.headers.ListByRequestID(requestID)
}

func (s *RequestHeaderService) Update(userID, id uint, key, value string, enabled *bool) (*models.RequestHeader, error) {
	h, err := s.headers.FindByIDForUser(id, userID)
	if err != nil {
		return nil, err
	}
	if h == nil {
		return nil, ErrNotFound
	}
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, ErrInvalidKey
	}
	h.Key = key
	h.Value = value
	if enabled != nil {
		h.Enabled = *enabled
	}
	if err := s.headers.Update(h); err != nil {
		return nil, err
	}
	return h, nil
}

func (s *RequestHeaderService) Delete(userID, id uint) error {
	err := s.headers.Delete(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

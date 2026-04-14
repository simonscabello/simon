package services

import (
	"errors"
	"net/url"
	"strings"

	"gorm.io/gorm"

	"simon/apps/api/internal/models"
	"simon/apps/api/internal/repositories"
)

var allowedMethods = map[string]struct{}{
	"GET": {}, "POST": {}, "PUT": {}, "DELETE": {},
}

type RequestService struct {
	reqs *repositories.RequestRepository
}

func NewRequestService(db *gorm.DB) *RequestService {
	return &RequestService{reqs: repositories.NewRequestRepository(db)}
}

func ValidateRequestURL(raw string) bool {
	s := strings.TrimSpace(raw)
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if u.Host == "" {
		return false
	}
	return true
}

func NormalizeMethod(m string) string {
	return strings.ToUpper(strings.TrimSpace(m))
}

func IsAllowedHTTPMethod(m string) bool {
	_, ok := allowedMethods[NormalizeMethod(m)]
	return ok
}

func (s *RequestService) Create(userID, collectionID uint, name, method, rawURL, body string) (*models.Request, error) {
	ok, err := s.reqs.CollectionOwnedByUser(collectionID, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNotFound
	}
	method = NormalizeMethod(method)
	if !IsAllowedHTTPMethod(method) {
		return nil, ErrInvalidMethod
	}
	if !ValidateRequestURL(rawURL) {
		return nil, ErrInvalidURL
	}
	req := &models.Request{
		Name:         strings.TrimSpace(name),
		Method:       method,
		URL:          strings.TrimSpace(rawURL),
		Body:         body,
		CollectionID: collectionID,
	}
	if req.Name == "" {
		return nil, ErrInvalidName
	}
	if err := s.reqs.Create(req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *RequestService) List(userID, collectionID uint) ([]models.Request, error) {
	ok, err := s.reqs.CollectionOwnedByUser(collectionID, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNotFound
	}
	return s.reqs.ListByCollectionID(collectionID)
}

func (s *RequestService) Update(userID, id uint, name, method, rawURL, body string) (*models.Request, error) {
	req, err := s.reqs.FindByIDForUser(id, userID)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, ErrNotFound
	}
	method = NormalizeMethod(method)
	if !IsAllowedHTTPMethod(method) {
		return nil, ErrInvalidMethod
	}
	if !ValidateRequestURL(rawURL) {
		return nil, ErrInvalidURL
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidName
	}
	req.Name = name
	req.Method = method
	req.URL = strings.TrimSpace(rawURL)
	req.Body = body
	if err := s.reqs.Update(req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *RequestService) Delete(userID, id uint) error {
	err := s.reqs.Delete(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

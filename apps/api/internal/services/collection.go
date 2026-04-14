package services

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"simon/apps/api/internal/models"
	"simon/apps/api/internal/repositories"
)

type CollectionService struct {
	cols *repositories.CollectionRepository
}

func NewCollectionService(db *gorm.DB) *CollectionService {
	return &CollectionService{cols: repositories.NewCollectionRepository(db)}
}

func (s *CollectionService) Create(userID uint, name string) (*models.Collection, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidName
	}
	col := &models.Collection{
		Name:   name,
		UserID: userID,
	}
	if err := s.cols.Create(col); err != nil {
		return nil, err
	}
	return col, nil
}

func (s *CollectionService) List(userID uint) ([]models.Collection, error) {
	return s.cols.ListByUserID(userID)
}

func (s *CollectionService) Update(userID, id uint, name string) (*models.Collection, error) {
	col, err := s.cols.FindByIDForUser(id, userID)
	if err != nil {
		return nil, err
	}
	if col == nil {
		return nil, ErrNotFound
	}
	col.Name = strings.TrimSpace(name)
	if err := s.cols.Update(col); err != nil {
		return nil, err
	}
	return col, nil
}

func (s *CollectionService) Delete(userID, id uint) error {
	err := s.cols.Delete(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *CollectionService) EnsureOwned(userID, id uint) (*models.Collection, error) {
	return s.cols.FindByIDForUser(id, userID)
}

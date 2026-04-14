package services

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"simon/apps/api/internal/models"
	"simon/apps/api/internal/repositories"
)

type EnvironmentService struct {
	envs *repositories.EnvironmentRepository
}

func NewEnvironmentService(db *gorm.DB) *EnvironmentService {
	return &EnvironmentService{envs: repositories.NewEnvironmentRepository(db)}
}

func (s *EnvironmentService) Create(userID uint, name string) (*models.Environment, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidName
	}
	e := &models.Environment{Name: name, UserID: userID}
	if err := s.envs.Create(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *EnvironmentService) List(userID uint) ([]models.Environment, error) {
	return s.envs.ListByUserID(userID)
}

func (s *EnvironmentService) Update(userID, id uint, name string) (*models.Environment, error) {
	e, err := s.envs.FindByIDForUser(id, userID)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, ErrNotFound
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidName
	}
	e.Name = name
	if err := s.envs.Update(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *EnvironmentService) Delete(userID, id uint) error {
	err := s.envs.Delete(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

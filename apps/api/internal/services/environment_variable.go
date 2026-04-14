package services

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"simon/apps/api/internal/models"
	"simon/apps/api/internal/repositories"
)

type EnvironmentVariableService struct {
	envs  *repositories.EnvironmentRepository
	vars  *repositories.EnvironmentVariableRepository
}

func NewEnvironmentVariableService(db *gorm.DB) *EnvironmentVariableService {
	return &EnvironmentVariableService{
		envs: repositories.NewEnvironmentRepository(db),
		vars: repositories.NewEnvironmentVariableRepository(db),
	}
}

func (s *EnvironmentVariableService) Create(userID, environmentID uint, key, value string, enabled *bool) (*models.EnvironmentVariable, error) {
	env, err := s.envs.FindByIDForUser(environmentID, userID)
	if err != nil {
		return nil, err
	}
	if env == nil {
		return nil, ErrNotFound
	}
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, ErrInvalidKey
	}
	dup, err := s.vars.FindByEnvironmentAndKey(environmentID, key)
	if err != nil {
		return nil, err
	}
	if dup != nil {
		return nil, ErrDuplicateKey
	}
	v := &models.EnvironmentVariable{
		EnvironmentID: environmentID,
		Key:           key,
		Value:         value,
		Enabled:       enabledOrDefault(enabled),
	}
	if err := s.vars.Create(v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *EnvironmentVariableService) List(userID, environmentID uint) ([]models.EnvironmentVariable, error) {
	env, err := s.envs.FindByIDForUser(environmentID, userID)
	if err != nil {
		return nil, err
	}
	if env == nil {
		return nil, ErrNotFound
	}
	return s.vars.ListByEnvironmentID(environmentID)
}

func (s *EnvironmentVariableService) Update(userID, id uint, key, value string, enabled *bool) (*models.EnvironmentVariable, error) {
	v, err := s.vars.FindByIDForUser(id, userID)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, ErrNotFound
	}
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, ErrInvalidKey
	}
	if key != v.Key {
		dup, err := s.vars.FindByEnvironmentAndKey(v.EnvironmentID, key)
		if err != nil {
			return nil, err
		}
		if dup != nil {
			return nil, ErrDuplicateKey
		}
	}
	v.Key = key
	v.Value = value
	if enabled != nil {
		v.Enabled = *enabled
	}
	if err := s.vars.Update(v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *EnvironmentVariableService) Delete(userID, id uint) error {
	err := s.vars.Delete(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

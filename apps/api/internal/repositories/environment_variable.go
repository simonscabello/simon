package repositories

import (
	"errors"

	"gorm.io/gorm"

	"simon/apps/api/internal/models"
)

type EnvironmentVariableRepository struct {
	db *gorm.DB
}

func NewEnvironmentVariableRepository(db *gorm.DB) *EnvironmentVariableRepository {
	return &EnvironmentVariableRepository{db: db}
}

func (r *EnvironmentVariableRepository) Create(v *models.EnvironmentVariable) error {
	return r.db.Create(v).Error
}

func (r *EnvironmentVariableRepository) ListByEnvironmentID(environmentID uint) ([]models.EnvironmentVariable, error) {
	var list []models.EnvironmentVariable
	err := r.db.Where("environment_id = ?", environmentID).Order("created_at asc").Find(&list).Error
	return list, err
}

func (r *EnvironmentVariableRepository) FindByIDForUser(id, userID uint) (*models.EnvironmentVariable, error) {
	var v models.EnvironmentVariable
	err := r.db.
		Joins("JOIN environments ON environments.id = environment_variables.environment_id AND environments.user_id = ?", userID).
		Where("environment_variables.id = ?", id).
		First(&v).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &v, nil
}

func (r *EnvironmentVariableRepository) FindByEnvironmentAndKey(environmentID uint, key string) (*models.EnvironmentVariable, error) {
	var v models.EnvironmentVariable
	err := r.db.Where("environment_id = ? AND key = ?", environmentID, key).First(&v).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &v, nil
}

func (r *EnvironmentVariableRepository) Update(v *models.EnvironmentVariable) error {
	return r.db.Save(v).Error
}

func (r *EnvironmentVariableRepository) Delete(id, userID uint) error {
	row, err := r.FindByIDForUser(id, userID)
	if err != nil {
		return err
	}
	if row == nil {
		return gorm.ErrRecordNotFound
	}
	return r.db.Delete(&models.EnvironmentVariable{}, row.ID).Error
}

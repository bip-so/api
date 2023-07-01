package shortner

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
)

func (r shortRepo) Create(instance *models.Short) (*models.Short, error) {
	fmt.Println(instance)
	results := r.db.Create(&instance)
	return instance, results.Error
}

func (r shortRepo) Get(query map[string]interface{}) (*models.Short, error) {
	var instance models.Short
	err := r.db.Where(query).First(&instance).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &instance, nil
}

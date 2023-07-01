package core

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

type StudioStatsRepository struct {
}

//err := postgres.GetDB().Table(tableName).Where("id", entityID).UpdateColumn("reel_count", gorm.Expr("reel_count  - ?", 1)).Error
// return err

func (s StudioStatsRepository) CounterPlusPlus(id uint64, key string, updateMap map[string]interface{}) error {

	err := postgres.GetDB().Table(models.STUDIO).Where("id", id).Updates(updateMap).Error
	return err
}

func (s StudioStatsRepository) CounterMinusMinus(id uint64, key string, updateMap map[string]interface{}) error {
	err := postgres.GetDB().Table(models.STUDIO).Where("id", id).Updates(updateMap).Error
	return err
}

func (s StudioStatsRepository) UpdateStudioCounter() {
	//var err error

}

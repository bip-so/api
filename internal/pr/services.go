package pr

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

func (s prService) GetPublishRequestsByStudio(studioID uint64, userId uint64) (*[]models.PublishRequest, error) {
	instances, err := App.Repo.GetPublishRequestsByStudio(map[string]interface{}{"studio_id": studioID, "status": models.PUBLISH_REQUEST_PENDING})
	if err != nil {
		return nil, err
	}
	return instances, nil
}

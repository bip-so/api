package shortner

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

func (s shortService) Get(code string) (*models.Short, error) {
	instance, err := App.Repo.Get(map[string]interface{}{"short_code": code})
	//instance2 := instance
	//instance2.ID = 0
	//instance2.UUID = uuid.New()
	//App.Repo.Create(instance2)

	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (s shortService) Create(url string) (*models.Short, error) {
	var shorty models.Short
	shorty.OriginalURL = url
	shorty.ShortCode = utils.NewShortNanoid()
	shorty.Count = 1
	instance, err := App.Repo.Create(&shorty)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

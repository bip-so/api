package mailers

import "gitlab.com/phonepost/bip-be-platform/internal/models"

func (r mailersRepo) GetUsersByIDs(userIDs []uint64) (users []models.User, err error) {
	err = r.db.Model(models.User{}).Select("email").Where("id in ?", userIDs).Find(&users).Error
	return nil, err
}

package queries

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

func (q studioInviteEmailQuery) GetListOfInvitedUserViaEmail(email string) []models.StudioInviteViaEmail {
	var studioInvites []models.StudioInviteViaEmail
	postgres.GetDB().Model(models.StudioInviteViaEmail{}).Where("email = ?", email).Find(&studioInvites)
	return studioInvites
}
func (q studioInviteEmailQuery) UpdateStudioViaEmailAccepted(id uint64) error {
	return postgres.GetDB().Model(models.StudioInviteViaEmail{}).Where("id = ?", id).Updates("has_accepted = true").Error
}

package bat

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
)

func (c batController) InviteViaEmailController(data CreateEmailInvite, user *models.User, canvasBranchID uint64, studioID uint64) error {
	for _, value := range data.Invites {
		hasBipAccount := queries.App.UserQueries.HasBIPAccount(value.Email)
		if hasBipAccount {
			invitedUserInstance := queries.App.UserQueries.GetUserInstanceByEmail(value.Email)
			_ = App.Service.InviteViaEmailExistingAccount(invitedUserInstance, user, value.CanvasPermissionsGroup, canvasBranchID, studioID)
		} else {
			_ = App.Service.InviteViaEmailNewAccount(value.Email, value.CanvasPermissionsGroup, user, canvasBranchID, studioID)
		}
	}
	return nil
}

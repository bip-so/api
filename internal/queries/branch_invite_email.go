package queries

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

func (q branchInviteViaEmailQuery) GetListOfInvitedUserViaEmail(email string) []models.BranchInviteViaEmail {
	var branchInvites []models.BranchInviteViaEmail
	postgres.GetDB().Model(models.BranchInviteViaEmail{}).Where("email = ?", email).Find(&branchInvites)
	return branchInvites
}

func (q branchInviteViaEmailQuery) BranchInviteViaEmailBranchInviteViaEmailExists(email string) bool {
	var exists bool
	_ = postgres.GetDB().Model(models.BranchInviteViaEmail{}).
		Select("count(*) > 0").
		Where("email = ?", email).
		Find(&exists).
		Error
	if exists {
		return true
	}

	return false
}

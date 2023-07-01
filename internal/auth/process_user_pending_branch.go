package auth

import (
	"gitlab.com/phonepost/bip-be-platform/internal/feed"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/internal/workflows"
)

func (s *authService) AddUserToStudioBranch(studioID uint64, branchID uint64, user *models.User, pg string) {
	// Add user to Studio as a member
	errJoinPersonToStudio := workflows.WorkflowJoinUserToStudio(user, studioID)
	if errJoinPersonToStudio != nil {
		return
	}

	//branch
	branch, errGetBranch := queries.App.BranchQuery.GetBranch(map[string]interface{}{"id": branchID}, true)

	if errGetBranch != nil {
		return
	}
	// Add permissions for this user
	memberInstance, err := queries.App.MemberQuery.GetMember(map[string]interface{}{"user_id": user.ID, "studio_id": studioID})

	if err != nil {
		return
	}
	errAddingPerms := permissions.App.Service.CreateCustomCanvasBranchPermission(branch.CanvasRepository.CollectionID, user.ID, studioID, branch.CanvasRepository.ID, branch.ID, branch.CanvasRepository.ParentCanvasRepositoryID, pg, memberInstance.ID)
	if errAddingPerms != nil {
		return
	}

	go func() {
		feed.App.Service.LeaveStudio(studioID, user.ID)
		// Adding view metadata to parent if not present
		permissions.App.Service.AddMemberToCollectionIfNotPresent(user.ID, memberInstance.ID, branch.CanvasRepository.CollectionID, studioID)
		if branch.CanvasRepository.ParentCanvasRepositoryID != nil {
			permissions.App.Service.AddMemberToCanvasIfNotPresent(user.ID, memberInstance.ID, *branch.CanvasRepository.ParentCanvasRepositoryID, studioID)
		}
	}()

	return
}

func (s *authService) ProcessUserPendingInvitesToBranch(user *models.User) {
	// Query branch_invite_via_emails to check if any pending invites are there
	if !queries.App.BranchInviteViaEmailQuery.BranchInviteViaEmailBranchInviteViaEmailExists(user.Email.String) {
		return
	}
	// BranchInviteViaEmail
	//var invited *[]models.BranchInviteViaEmail
	//_ = s.db.Model(&models.BranchInviteViaEmail{}).Where("email = ?", user.Email.String).Find(&invited).Error
	invited := queries.App.BranchInviteViaEmailQuery.GetListOfInvitedUserViaEmail(user.Email.String)
	// if found we need to do the operations basically
	for _, key := range invited {
		App.Service.AddUserToStudioBranch(key.StudioID, key.BranchID, user, key.PermissionGroup)
	}
}

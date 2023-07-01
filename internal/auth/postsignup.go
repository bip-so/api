package auth

import "gitlab.com/phonepost/bip-be-platform/internal/models"

func (s *authService) PostSignupWorkflow(user *models.User, isLegacy bool) {

	// When a new user signs up we are checking if There are any invites pending
	App.Service.ProcessInviteViaEmail(user)
	App.PostSignup.CreateDefaultStudio(user)
	if isLegacy {
		go App.Service.ProcessUserPendingInvitesToBranch(user)
	}

}

package auth

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/internal/workflows"

	"gitlab.com/phonepost/bip-be-platform/internal/feed"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

func (s *authService) ProcessInviteViaEmail(user *models.User) {
	// Query branch_invite_via_emails to check if any pending invites are there
	//var studioInvites []models.StudioInviteViaEmail
	//s.db.Model(models.StudioInviteViaEmail{}).
	//	Where("email = ?", user.Email.String).
	//	Find(&studioInvites)
	studioInvites := queries.App.StudioInviteEmailQuery.GetListOfInvitedUserViaEmail(user.Email.String)
	for k, inviteInstance := range studioInvites {
		fmt.Println("processing Invite #", k)
		_ = s.UserInviteWhenUserExists(user, inviteInstance)
	}

}

func (s *authService) UserInviteWhenUserExists(NewlyCreatedUser *models.User, inviteInstance models.StudioInviteViaEmail) error {
	// Get the user instance/ Can't have error here
	// Create a Member for the studio if not exists
	errJoinStudioSafe := workflows.WorkflowJoinUserToStudio(NewlyCreatedUser, inviteInstance.StudioID)
	if errJoinStudioSafe != nil {
		println(errJoinStudioSafe.Error(), NewlyCreatedUser.Email.String)
	}
	// This needs to be called after above is called as that is creating the member and stuff
	member1, _ := queries.App.MemberQuery.GetMember(map[string]interface{}{"user_id": NewlyCreatedUser.ID, "studio_id": inviteInstance.StudioID})
	fakeMemberArray := []models.Member{*member1}
	// Finally, add the required role ids to This User(Member)
	// We can have many roles
	for _, roleID := range inviteInstance.Roles2 {
		// We need to add invitedUser to th roleID
		// Add UserID (Member to Role)
		//postgres.GetDB().Model(&models.Role{}).Where("id = ?", utils.Uint64(roleID)).First(&ThisRole)
		role1, _ := queries.App.RoleQuery.GetRole(utils.Uint64(roleID))
		// Incase the role was deleted
		if role1.ID != 0 {
			queries.App.MemberQuery.BulkAddMembersToRole(&fakeMemberArray, role1)
		}
	}
	// Invited user is not member and also added to the Roles.
	// We'll also need to Create and Instance for the Studio Invited instance and Set as accepted for the record.
	// Delete
	//s.db.Model(models.StudioInviteViaEmail{}).Where("id = ?", inviteInstance.ID).Updates("has_accepted = true")
	_ = queries.App.StudioInviteEmailQuery.UpdateStudioViaEmailAccepted(inviteInstance.ID)
	// Updating User Associated Studios
	//queries.App.StudioQueries.DeleteUserAssociatedStudioDataByUserID(NewlyCreatedUser.ID)
	//supabase.UpdateUserSupabase(NewlyCreatedUser.ID, true)
	feed.App.Service.JoinStudio(inviteInstance.StudioID, NewlyCreatedUser.ID)

	return nil
}

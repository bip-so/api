package studio

import (
	"errors"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"

	"github.com/lib/pq"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	user2 "gitlab.com/phonepost/bip-be-platform/internal/user"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/supabase"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

/*
	We need to allow people to invite people to invite to studio
	User sends an email to Many Emails.
	We'll Check if the email already exists on Platform,
	We'll add them
	Add to studio as a Member
	Specific Role
	(Roles[] and Emails[])
*/

func (sr studioRepo) IsUserAlreadyOnBIP(email string) bool {
	var user models.User
	var result int64
	err := sr.db.Model(&user).Where("email = ?", email).Count(&result).Error
	if err != nil {
		return false
	}
	if result == 1 {
		return true
	} else {
		return false
	}
}

//

type NewInvitePostOne struct {
	Email string   `json:"email"`
	Roles []uint64 `json:"roles"`
}

// Entry Point
func (sc *studioController) CreateStudioInvites(data []NewInvitePostOne, studioID uint64, currentUser *models.User) []string {
	// Loop through data
	var errorResponse []string
	for _, entry := range data {
		// Check if the user is already part of studio
		// This function is a stub need to write this right now always false
		// Anyways contunie
		if App.StudioService.DoesUserAlreadyBelongsToStudio(entry.Email, studioID) {
			errorResponse = append(errorResponse, entry.Email+" is already in the Studio")
			continue
		}
		// Check if User with Email Exists or not
		if App.StudioRepo.IsUserAlreadyOnBIP(entry.Email) {
			// User with this email found
			errorAddingEmail := App.StudioService.UserInviteWhenUserExists(entry.Email, entry, studioID, currentUser)
			if errorAddingEmail != nil {
				errorResponse = append(errorResponse, errorAddingEmail.Error())
			} else {
				// Informing user he has been added to Studio
				App.StudioService.InformUserAddedToStudioSendMailer(entry.Email, currentUser, studioID)
			}

		} else {
			// User with this email not found
			App.StudioService.UserInviteWhenUserDoesNotExists(entry, studioID, currentUser)
			// Todo : Email Needs to be Formatted
			App.StudioService.InviteNewUserSendMailer(entry.Email, currentUser, studioID)
		}

	}

	return errorResponse
}

func (ss studioService) UserInviteWhenUserExists(email string, data NewInvitePostOne, studioID uint64, currentUser *models.User) error {
	// Get the user instance/ Can't have error here
	invitedUser, _ := user2.App.Repo.GetUser(map[string]interface{}{"email": email})
	// Create a Member for the studio if not exists
	errJoinStudioSafe := ss.JoinStudioSafe(invitedUser, studioID)
	if errJoinStudioSafe != nil {
		return errors.New(email + " " + errJoinStudioSafe.Error())
	}
	// This needs to be called after above is called as that is creating the member and stuff
	member1, _ := queries.App.MemberQuery.GetMember(map[string]interface{}{"user_id": invitedUser.ID, "studio_id": studioID})
	fakeMemberArray := []models.Member{*member1}
	// Finally, add the required role ids to This User(Member)
	// We can have many roles
	for _, roleID := range data.Roles {
		// We need to add invitedUser to th roleID
		// Add UserID (Member to Role)
		role1, _ := queries.App.RoleQuery.GetRole(roleID)
		//postgres.GetDB().Model(&models.Role{}).Where("id = ?", roleID).First(&ThisRole)
		//member.App.Repo.AddMembersInRole(fakeMemberArray, ThisRole)
		queries.App.MemberQuery.BulkAddMembersToRole(&fakeMemberArray, role1)
	}
	// Invited user is not member and also added to the Roles.
	// We'll also need to Create and Instance for the Studio Invited instance and Set as accepted for the record.
	App.StudioRepo.CreateNewStudioInviteInstance(studioID, data, currentUser.ID, true)
	// Todo: We need to send notifications to the UserInviteWhenUserExists users

	// Updating User Associated Studios
	queries.App.StudioQueries.DeleteUserAssociatedStudioDataByUserID(invitedUser.ID)
	supabase.UpdateUserSupabase(invitedUser.ID, true)

	return nil
}

// CreateNewStudioInviteInstance: Only creates Instance
func (ss studioRepo) CreateNewStudioInviteInstance(studioID uint64, data NewInvitePostOne, userID uint64, hasAccepted bool) {
	strRoleIDs := pq.StringArray{}
	for _, roleID := range data.Roles {
		strRoleIDs = append(strRoleIDs, utils.String(roleID))
	}
	NewStudioInviteViaEmailObject := models.StudioInviteViaEmail{}
	NewStudioInviteViaEmailObject.Email = data.Email
	NewStudioInviteViaEmailObject.StudioID = studioID
	// NewStudioInviteViaEmailObject.Roles = data.Roles
	NewStudioInviteViaEmailObject.Roles2 = strRoleIDs
	NewStudioInviteViaEmailObject.PermissionGroup = "NA"
	NewStudioInviteViaEmailObject.CreatedByID = userID
	NewStudioInviteViaEmailObject.HasAccepted = hasAccepted
	result := ss.db.Create(&NewStudioInviteViaEmailObject)
	fmt.Println(result.Error)
}

func (ss studioService) UserInviteWhenUserDoesNotExists(data NewInvitePostOne, studioID uint64, currentUser *models.User) {
	App.StudioRepo.CreateNewStudioInviteInstance(studioID, data, currentUser.ID, false)
}

//
func (s *studioService) JoinStudioSafe(user *models.User, studioId uint64) error {
	var err error
	mem, err := queries.App.MemberQuery.GetMember(map[string]interface{}{"user_id": user.ID, "studio_id": studioId})
	if err == nil && mem != nil {
		if mem.IsRemoved {
			return errors.New(" User was banned")
		} else if mem.HasLeft {
			return queries.App.MemberQuery.JoinStudio([]uint64{user.ID}, studioId)
		}
	} else {
		//memberId := member.App.Service.AddMemberToStudio(user.ID, studioId)
		member2 := queries.App.MemberQuery.AddUserIDToStudio(user.ID, studioId)
		if member2 == nil {
			return errors.New(" Could not add the user to Studio")
		}
		memberObj, err := queries.App.MemberQuery.GetMember(map[string]interface{}{"id": member2.ID})
		if err != nil {
			return err
		}
		err = queries.App.MemberQuery.AddMembersToStudioInMemberRole(studioId, []models.Member{*memberObj})
		if err != nil {
			return err
		}
	}
	return nil
}

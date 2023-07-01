package bat

import (
	"errors"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"

	"gitlab.com/phonepost/bip-be-platform/internal/feed"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

func (s batService) CreateBranchAccessToken(user *models.User, branchID uint64, pg string) (*models.BranchAccessToken, error) {
	var instance models.BranchAccessToken
	instance.InviteCode = utils.NewShortNanoid()
	instance.BranchID = branchID
	instance.CreatedByID = user.ID
	instance.IsActive = true
	instance.PermissionGroup = pg
	mi, err := queries.App.BranchAccessTokenQuery.CreateBranchTokenInstance(&instance)
	if err != nil {
		return nil, err
	}
	return mi, nil
}

func (s batService) GetBranchAccessToken(code string) (*models.BranchAccessToken, error) {
	instance, err := queries.App.BranchAccessTokenQuery.GetBranchTokenInstance(map[string]interface{}{"invite_code": code})
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (s batService) JoinUserStudioBranch(studioID uint64, branchID uint64, user *models.User, bat *models.BranchAccessToken) error {
	// Add user to Studio as a member
	errJoinPersonToStudio := s.JoinStudioSafe(user, studioID)
	if errJoinPersonToStudio != nil {
		return errJoinPersonToStudio
	}

	//branch
	branch, errGetBranch := queries.App.BranchQuery.GetBranchWithRepoAndStudio(branchID)
	if errGetBranch != nil {
		return errGetBranch
	}
	// Add permissions for this user
	memberInstance, err := queries.App.MemberQuery.GetMember(map[string]interface{}{"user_id": user.ID, "studio_id": studioID})

	if err != nil {
		return err
	}
	errAddingPerms := permissions.App.Service.CreateCustomCanvasBranchPermission(branch.CanvasRepository.CollectionID, user.ID, studioID, branch.CanvasRepository.ID, branch.ID, branch.CanvasRepository.ParentCanvasRepositoryID, bat.PermissionGroup, memberInstance.ID)
	if errAddingPerms != nil {
		return errAddingPerms
	}
	// Adding permissions to parent if not present.
	permissions.App.Service.AddMemberToCollectionIfNotPresent(user.ID, memberInstance.ID, branch.CanvasRepository.CollectionID, studioID)
	if branch.CanvasRepository.ParentCanvasRepositoryID != nil {
		permissions.App.Service.AddMemberToCanvasIfNotPresent(user.ID, memberInstance.ID, *branch.CanvasRepository.ParentCanvasRepositoryID, studioID)
	}

	go func() {
		//queries.App.StudioQueries.DeleteUserAssociatedStudioDataByUserID(user.ID)
		//supabase.UpdateUserSupabase(user.ID, true)
		feed.App.Service.LeaveStudio(studioID, user.ID)
	}()

	return nil
}

// Todo: Refactor
func (s *batService) JoinStudioSafe(user *models.User, studioId uint64) error {
	var err error
	mem, err := queries.App.MemberQuery.GetMember(map[string]interface{}{"user_id": user.ID, "studio_id": studioId})
	if err == nil && mem != nil {
		if mem.IsRemoved {
			return errors.New("Can't join as you were banned")
		} else if mem.HasLeft {
			return queries.App.MemberQuery.JoinStudio([]uint64{user.ID}, studioId)
		}
	} else {
		//memberId := member.App.Service.AddMemberToStudio(user.ID, studioId)
		member2 := queries.App.MemberQuery.AddUserIDToStudio(user.ID, studioId)
		if member2 == nil {
			return errors.New("error in adding user to studio")
		}
		memberObj, err := queries.App.MemberQuery.GetMember(map[string]interface{}{"id": member2.ID})
		if err != nil {
			return err
		}
		//err = member.App.Service.AddMembersToStudioMemberRole(studioId, []models.Member{*memberObj})
		err = queries.App.MemberQuery.AddMembersToStudioInMemberRole(studioId, []models.Member{*memberObj})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *batService) DeleteBAT(id uint64) error {
	return queries.App.BranchAccessTokenQuery.DeleteBranchToken(id)
}

//(invitedUserInstance, value.CanvasPermissionsGroup, user, canvasBranchID, studioID, invitedUserInstance)
func (s *batService) InviteViaEmailExistingAccount(inviteduser *models.User, invitedByuser *models.User, pg string, canvasBranchID uint64, studioID uint64) error {
	/*
		Steps
			- Add Person to Studio if not already
			- Create PG for Person
	*/
	App.Service.AddUserToStudioBranch(studioID, canvasBranchID, inviteduser, pg)
	return nil
}

func (s *batService) AddUserToStudioBranch(studioID uint64, branchID uint64, user *models.User, pg string) error {
	// Add user to Studio as a member
	_, errJoinPersonToStudio := queries.App.StudioQueries.SafeAddUserToStudio(user.ID, studioID)
	if errJoinPersonToStudio != nil {
		return errJoinPersonToStudio
	}

	//branch
	//branch, errGetBranch := App.Repo.GetBranchWithRepo(map[string]interface{}{"id": branchID})
	branch, errGetBranch := queries.App.BranchQuery.GetBranchWithRepoAndStudio(branchID)

	if errGetBranch != nil {
		return errGetBranch
	}
	// Add permissions for this user
	memberInstance, err := queries.App.MemberQuery.GetMember(map[string]interface{}{"user_id": user.ID, "studio_id": studioID})
	if err != nil {
		return err
	}
	errAddingPerms := permissions.App.Service.CreateCustomCanvasBranchPermission(branch.CanvasRepository.CollectionID, user.ID, studioID, branch.CanvasRepository.ID, branch.ID, branch.CanvasRepository.ParentCanvasRepositoryID, pg, memberInstance.ID)
	if errAddingPerms != nil {
		return errAddingPerms
	}

	go func() {
		//queries.App.StudioQueries.DeleteUserAssociatedStudioDataByUserID(user.ID)
		//supabase.UpdateUserSupabase(user.ID, true)
		// @todo: SR Ask
		feed.App.Service.LeaveStudio(studioID, user.ID)
	}()

	return nil
}

func (s *batService) InviteViaEmailNewAccount(email string, pg string, user *models.User, canvasBranchID uint64, studioID uint64) error {
	// We don't have this email in our system
	// We have to create instance of BranchInviteViaEmail
	newInvite := models.BranchInviteViaEmail{}
	//newInvite.Email = data.Invites
	newInvite.Email = email
	newInvite.PermissionGroup = pg
	newInvite.CreatedByID = user.ID
	newInvite.BranchID = canvasBranchID
	newInvite.StudioID = studioID
	_ = postgres.GetDB().Create(&newInvite).Error
	// Resend Email may be needed
	go App.Service.InviteNewUserSendMailer(email, user, studioID, canvasBranchID)
	return nil
}

func (s *batService) InviteNewUserSendMailer(email string, InvitedByUser *models.User, studioID, canvasBranchID uint64) {
	// get studio instance
	var studio models.Studio
	//ENVMAILER := configs.GetConfigString("ENV")
	err := postgres.GetDB().Model(models.Studio{}).Where("id = ?", studioID).First(&studio).Error
	//studioName := studio.Handle
	//studioURL := models.MailerRouterPaths[ENVMAILER]["BASE_URL"] + "@" + studioName
	canvasUrl := notifications.App.Service.GenerateCanvasBranchUrlByID(canvasBranchID)
	body := "Hello, <br> You have been invite to canvas in " + canvasUrl + " by " + InvitedByUser.Username + "</strong></div>"
	bodyPlainText := "Hello, <br> You have been invite to canvas in " + canvasUrl + " by " + InvitedByUser.Username + "</strong></div>"
	subject := "Invited to BIP Studio Canvas"
	var mailer pkg.BipMailer
	toList := []string{email}
	emptyList := []string{}
	err = pkg.BipMailer.SendEmail(mailer, toList, emptyList, emptyList, subject, body, bodyPlainText)
	if err != nil {
		fmt.Println(err)
	}
}

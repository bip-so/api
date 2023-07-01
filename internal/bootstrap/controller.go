package bootstrap

import (
	"encoding/json"

	"gitlab.com/phonepost/bip-be-platform/internal/queries"

	"gitlab.com/phonepost/bip-be-platform/internal/follow"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/studio"
	"gitlab.com/phonepost/bip-be-platform/internal/user"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gorm.io/gorm"
)

var (
	bootstrapController = BootstrapController{}
)

type BootstrapController struct{}

func (c BootstrapController) GetUserAssociatedStudios(userID uint64) (*models.UserAssociatedStudio, error) {
	userStudios, err := studio.App.UserAssociatedStudioRepo.GetUserAssociatedStudioDataByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			members, err := queries.App.MemberQuery.GetAllStudiosUserMemberOf(userID)
			if err != nil {
				return nil, err
			}
			studioIDsMap := map[uint64]bool{}
			studioIDs := []uint64{}
			for _, member := range members {
				if !studioIDsMap[member.StudioID] {
					studioIDsMap[member.StudioID] = true
					studioIDs = append(studioIDs, member.StudioID)
				}
			}

			allStudios, err := studio.App.StudioRepo.GetStudiosByIDs(studioIDs)
			if err != nil {
				return nil, err
			}

			studioData := []AssociatedStudio{}
			for _, stdio := range *allStudios {
				associatedStdio := AssociatedStudio{
					ID:          stdio.ID,
					UUID:        string(stdio.UUID.String()),
					DisplayName: stdio.DisplayName,
					Handle:      stdio.Handle,
					ImageURL:    stdio.ImageURL,
					CreatedByID: stdio.CreatedByID,
				}
				studioData = append(studioData, associatedStdio)
			}
			userCreatedStudios := bootstrapController.GetSortedUserAssociatedStudios(studioData, userID)
			data, err := json.Marshal(userCreatedStudios)
			if err != nil {
				return nil, err
			}
			userStudios := models.NewUserAssociatedStudio(userID, string(data))
			err = studio.App.UserAssociatedStudioRepo.CreateUserAssociatedStudio(userStudios)
			if err != nil {
				logger.Error("Failed to create UserAssociatedStudio.")
			}
			return userStudios, nil
		} else {
			return nil, err
		}
	}
	return userStudios, err
}

func (c BootstrapController) HandleController(handle string, authUser *models.User) (*models.Studio, *[]models.Member, *models.User, *[]models.FollowUser, error) {
	studioInstance, members, err := studio.App.Controller.GetStudioByHandleController(handle, authUser)
	if err != nil && err == gorm.ErrRecordNotFound {
		userInstance, err := user.App.Controller.GetUserByHandleController(handle)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		if authUser == nil {
			return nil, nil, userInstance, nil, nil
		}
		fUser, err := follow.App.Repo.GetIsUserFollowingUser(authUser.ID, []uint64{userInstance.ID})
		if err != nil {
			return nil, nil, nil, nil, err
		}
		return nil, nil, userInstance, fUser, nil
	}
	return studioInstance, members, nil, nil, err
}
func (c BootstrapController) GetSortedUserAssociatedStudios(studioData []AssociatedStudio, userID uint64) []AssociatedStudio {
	var userAdminStudios []AssociatedStudio
	var userNonAdminStudios []AssociatedStudio
	var userCreatedStudios []AssociatedStudio
	for _, studioInstance := range studioData {
		if studioInstance.CreatedByID == userID {
			userCreatedStudios = append(userCreatedStudios, studioInstance)
			continue
		}
		flag, _ := queries.App.UserQueries.IsUserAdminInStudio(userID, studioInstance.ID)
		if flag {
			for _, t := range studioData {
				if t.ID == studioInstance.ID {
					userAdminStudios = append(userAdminStudios, t)
				}
			}
		} else {
			for _, t := range studioData {
				if t.ID == studioInstance.ID {
					userNonAdminStudios = append(userNonAdminStudios, t)
				}
			}
		}
	}
	userCreatedStudios = append(userCreatedStudios, userAdminStudios...)
	userCreatedStudios = append(userCreatedStudios, userNonAdminStudios...)
	return userCreatedStudios
}

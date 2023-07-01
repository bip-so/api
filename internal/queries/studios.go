package queries

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

// Preloads Topics
func (q *studioQuery) GetStudioQuery(query map[string]interface{}) (*models.Studio, error) {
	var object *models.Studio
	result := postgres.GetDB().Model(models.Studio{}).Where(query).Preload("Topics").First(&object)
	return object, result.Error
}

func (q *studioQuery) GetStudioInstanceByID(studioID uint64) (*models.Studio, error) {
	var studio models.Studio
	err := postgres.GetDB().Model(models.Studio{}).Where("id = ? and is_archived = ?", studioID, false).Preload("CreatedByUser").First(&studio).Error
	return &studio, err
}

func (q *studioQuery) UpdateStudioByID(studioID uint64, updates map[string]interface{}) error {
	var studioInstance *models.Studio
	err := postgres.GetDB().Model(models.Studio{}).Where(models.Studio{BaseModel: models.BaseModel{ID: studioID}}).Updates(updates).Find(&studioInstance).Error
	return err
}

func (q *studioQuery) StudioMemberCount(id uint64) int64 {
	var result int64
	postgres.GetDB().Model(&models.Member{}).Where("studio_id = ? and has_left = ? and is_removed = ?", id, false, false).Count(&result)
	return result
}

func (q *studioQuery) NewStudioInstance(name, handle, description, website, imageURL string, userID uint64) *models.Studio {
	return &models.Studio{
		DisplayName: name,
		Handle:      handle,
		Description: description,
		Website:     website,
		ImageURL:    imageURL,
		CreatedByID: userID,
		UpdatedByID: userID,
	}
}

//DeleteUserAssociatedStudioDataByUserID : Updated the Associated User Studio Table on Member Removal
func (q *studioQuery) DeleteUserAssociatedStudioDataByUserID(userID uint64) (*models.UserAssociatedStudio, error) {
	var userStudios models.UserAssociatedStudio
	err := postgres.GetDB().Where(models.UserAssociatedStudio{UserID: userID}).Delete(&userStudios).Error
	return &userStudios, err
}

// Adds a User to Studio As a Member and is idempotent
func (r studioQuery) SafeAddUserToStudio(userID uint64, studioID uint64) (*models.Member, error) {

	// Ok this is an edge case we are implementing here.
	// We are checking if this user has a request pending.
	var quickcheck models.StudioMembersRequest
	_ = postgres.GetDB().Model(&models.StudioMembersRequest{}).Where(map[string]interface{}{
		"studio_id": studioID,
		"user_id":   userID,
	}).First(&quickcheck).Error
	if quickcheck.ID != 0 {
		// Delete this Instance
		_ = postgres.GetDB().Delete(quickcheck)

	}
	// End of Edge Case

	// Check if user is already a member
	var count int64
	_ = postgres.GetDB().Model(&models.Member{}).Where("studio_id = ? AND user_id = ?", studioID, userID).Count(&count).Error
	if count == 0 {
		// Not a member
		memberObject := App.MemberQuery.NewMemberObject(userID, studioID)
		errCreatingMemeber := postgres.GetDB().Create(memberObject).Error
		if errCreatingMemeber != nil {
			return nil, errCreatingMemeber
		}
		return memberObject, nil
	}
	var member models.Member
	errGettingMember := postgres.GetDB().Model(&models.Member{}).Where("studio_id = ? AND user_id = ?", studioID, userID).First(&member).Error
	if errGettingMember != nil {
		return nil, errGettingMember
	}
	return &member, nil
}

func (q *studioQuery) IsPersonalStudio(studioID uint64) bool {
	isPersonalStudio := false
	studio, err := q.GetStudioQuery(map[string]interface{}{"id": studioID})
	if err != nil {
		return false
	}
	user, err := App.UserQueries.GetUserByID(studio.CreatedByID)
	if err != nil {
		return false
	}
	if user.DefaultStudioID == studio.ID {
		isPersonalStudio = true
	}
	return isPersonalStudio
}

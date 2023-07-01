package queries

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"log"
)

func (q memberQuery) NewMemberObject(userId uint64, studioId uint64) *models.Member {
	return &models.Member{
		UserID:      userId,
		StudioID:    studioId,
		CreatedByID: userId,
		UpdatedByID: userId,
	}
}
func (q *memberQuery) GetMember(query map[string]interface{}) (*models.Member, error) {
	var member models.Member
	err := postgres.GetDB().Model(&models.Member{}).Where(query).First(&member).Error
	if err != nil {
		//log.Fatalln(err)
		return nil, err
	}
	return &member, nil
}

func (m *memberQuery) GetMembers(query map[string]interface{}) ([]models.Member, error) {
	var members []models.Member
	result := postgres.GetDB().Model(&models.Member{}).
		Where(query).
		Find(&members)

	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}
func (m memberQuery) GetMemberInstanceWithdUserInstance(query map[string]interface{}) (*models.Member, error) {
	var member models.Member
	err := postgres.GetDB().Model(&models.Member{}).Where(query).Preload("User").First(&member).Error
	if err != nil {
		//log.Fatalln(err)
		return nil, err
	}
	return &member, nil
}

func (q memberQuery) CreateMember(member *models.Member) *models.Member {
	if err := postgres.GetDB().Create(&member).Error; err != nil {
		log.Fatalln(err)
		return nil
	}
	return member
}

// AddUserIDToStudio Returns : MemberID
func (q memberQuery) AddUserIDToStudio(userdId uint64, studioId uint64) *models.Member {
	// Workflow Removed any pending requests for this user to join the studio
	App.StudioMemberRequestQuery.CheckAndDeleteStudioMembersRequestIfExists(studioId, userdId)
	memberObject := q.NewMemberObject(userdId, studioId)
	createdMember := q.CreateMember(memberObject)
	if createdMember == nil {
		return nil
	}
	return createdMember
}

func (q memberQuery) AddMembersToStudioInMemberRole(studioID uint64, members []models.Member) error {
	memberRole := App.RoleQuery.GetStudioMemberRole(studioID)
	err := q.BulkAddMembersToRole(&members, memberRole)
	return err
}

func (q memberQuery) BulkAddMembersToRole(addMembers *[]models.Member, role *models.Role) error {
	err := postgres.GetDB().Model(&role).Association("Members").Append(addMembers)
	return err
}

func (q memberQuery) GetMemberCountForStudio(studioId uint64) (int64, error) {
	var result int64
	err := postgres.GetDB().Model(&models.Member{}).Where("studio_id = ? and has_left = ? and is_removed = ?", studioId, false, false).Count(&result).Error
	return result, err
}

func (q memberQuery) UpdateHasLeft(userIDs []uint64, studioID uint64, hasLeft bool) error {
	updates := map[string]interface{}{
		"has_left": hasLeft,
	}
	err := postgres.GetDB().Model(&models.Member{}).Where("user_id IN ? and studio_id = ?", userIDs, studioID).Updates(updates).Error
	return err
}

// GetMembersByUserIDs : Returns list of Members given UserID slice and studioID
func (q memberQuery) GetMembersByUserIDs(userIds []uint64, studioID uint64) ([]models.Member, error) {
	var members []models.Member
	result := postgres.GetDB().Model(&models.Member{}).Where("user_id IN ? and studio_id = ?", userIds, studioID).Find(&members)

	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

// GetMembersByStudioIDPaginated : Paginated query to get members of a studio
func (q memberQuery) GetMembersByStudioIDPaginated(studioID uint64, skip int) ([]models.Member, error) {
	var members []models.Member
	result := postgres.GetDB().Model(&models.Member{}).
		Where("studio_id = ? AND has_left = false AND is_removed = false", studioID).
		Preload("User").
		Preload("Roles").
		Find(&members)

	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

func (m memberQuery) BanUser(userID uint64, studioID uint64, banReason string, removedByID uint64) error {
	updates := map[string]interface{}{
		"removed_reason": banReason,
		"is_removed":     true,
		"removed_by_id":  removedByID,
	}
	err := postgres.GetDB().Model(&models.Member{}).Where("user_id = ? and studio_id = ?", userID, studioID).Updates(updates).Error
	return err
}

func (q *memberQuery) RemoveMembersInRole(removeMembers []models.Member, role *models.Role) error {
	err := postgres.GetDB().Model(&role).Association("Members").Delete(removeMembers)
	return err
}

func (m *memberQuery) MembersSearch(username string, studioID uint64) ([]models.Member, error) {
	username = username + "%"
	var members []models.Member
	err := postgres.GetDB().Table("members").
		Joins("left join users on users.id = members.user_id").
		Where(" members.studio_id = ? and (users.username ILIKE ? or users.full_name ILIKE ? or users.email ILIKE ?)", studioID, username, username, username).
		Preload("User").
		Find(&members).Error
	return members, err
}

func (q *memberQuery) GetMultipleMembersByID(memberIDs []uint64) ([]models.Member, error) {
	var members []models.Member
	result := postgres.GetDB().Model(&models.Member{}).Where("id IN ?", memberIDs).Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

func (q *memberQuery) AddMembersInBulkWithUserID(userIds []uint64, studioId uint64) ([]models.Member, error) {
	var memberObjects []models.Member
	for _, userId := range userIds {
		membObject := q.NewMemberObject(userId, studioId)
		memberObjects = append(memberObjects, *membObject)
	}
	err := q.CreateBulkMembers(memberObjects)
	if err != nil {
		return nil, err
	}
	err = q.AddMembersToStudioInMemberRole(studioId, memberObjects)
	return memberObjects, err
}

func (m *memberQuery) CreateBulkMembers(members []models.Member) error {
	err := postgres.GetDB().Create(&members).Error
	return err
}

func (m *memberQuery) GetMembersByStudioIDandUserIDs(studioID uint64, userIDs []uint64, skip int) ([]models.Member, error) {
	var members []models.Member
	result := postgres.GetDB().Model(&models.Member{}).
		Where("studio_id = ? AND has_left = false AND is_removed = false and user_id in ?", studioID, userIDs).
		Order("id asc").Offset(skip).
		Preload("User").
		Preload("Roles").
		Find(&members)

	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

func (m *memberQuery) GetMembersOfUserInMultipleStudios(studioIDs []uint64, userID uint64) (*[]models.Member, error) {
	var members []models.Member
	err := postgres.GetDB().Model(&models.Member{}).Where("studio_id IN ? and user_id = ? and is_removed = ? and has_left = ?", studioIDs, userID, false, false).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return &members, nil
}
func (m *memberQuery) GetAllStudiosUserMemberOf(userID uint64) ([]models.Member, error) {
	var members []models.Member
	result := postgres.GetDB().Model(&models.Member{}).Where("user_id = ? and is_removed = false and has_left = false", userID).Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

func (q *memberQuery) JoinStudio(userIds []uint64, studioId uint64) error {
	// This is in case user had left the studio
	err := q.UpdateHasLeft(userIds, studioId, false)
	members, err := q.GetMembersByUserIDs(userIds, studioId)
	if err != nil {
		return err
	}
	err = q.AddMembersToStudioInMemberRole(studioId, members)
	return err
}

func (m *memberQuery) GetMembersWithPreload(query map[string]interface{}) ([]models.Member, error) {
	var members []models.Member
	result := postgres.GetDB().Model(&models.Member{}).
		Where(query).
		Preload("User").
		Find(&members)

	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

func (m *memberQuery) RoleMembersSearch(username string, roleID uint64, skip, limit int) ([]models.RoleMember, error) {
	username = "%" + username + "%"
	var roleMembers []models.RoleMember
	err := postgres.GetDB().Table("role_members").
		Joins("left join members on members.id = role_members.member_id").
		Joins("left join users on users.id = members.user_id").
		Where(" role_members.role_id = ? and (users.username ILIKE ? or users.full_name ILIKE ? or users.email ILIKE ?)", roleID, username, username, username).
		Offset(skip).
		Limit(limit).
		Find(&roleMembers).Error
	return roleMembers, err
}

package permissions

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

func (r permissionRepo) getMembersByUserIDs(userIds []uint64, studioID uint64) ([]models.Member, error) {
	var members []models.Member
	result := postgres.GetDB().Model(&models.Member{}).
		Where("user_id IN ? and studio_id = ? AND has_left = false AND is_removed = false", userIds, studioID).
		Preload("Roles").
		Find(&members)

	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

func (r permissionRepo) getCollectionPermissionsForUser(memberIds, roleIds []uint64) ([]models.CollectionPermission, error) {
	var collectionPerms []models.CollectionPermission
	err := postgres.GetDB().Model(&collectionPerms).
		Where("member_id IN ? OR role_id IN ?", memberIds, roleIds).
		Preload("Studio").
		Preload("Role").
		Preload("Collection").
		Preload("Member").
		Preload("Member.User").
		Preload("Role.Members").
		Preload("Collection").Find(&collectionPerms).Error
	return collectionPerms, err
}

func getAllBranchesUserMemberOf(userID uint64) ([]models.Member, error) {
	var members []models.Member
	result := postgres.GetDB().Model(&models.Member{}).Where("user_id = ?", userID).Preload("Roles").Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

func (r permissionRepo) getAllStudiosUserMemberOf(userID uint64) ([]models.Member, error) {
	var members []models.Member
	result := postgres.GetDB().Model(&models.Member{}).Where("user_id = ? AND has_left = false AND is_removed = false", userID).Preload("Roles").Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

func (r permissionRepo) getStudioPermissionsByMemberIds(memberIds []uint64, roleIds []uint64) ([]models.StudioPermission, error) {
	var studioPerms []models.StudioPermission
	err := postgres.GetDB().Model(&studioPerms).Where("(member_id IN ? OR role_id IN ?)", memberIds, roleIds).Order("studio_id ASC").Find(&studioPerms).Error
	return studioPerms, err
}

func (r permissionRepo) GetMember(query map[string]interface{}) (*models.Member, error) {
	var member models.Member
	err := postgres.GetDB().Model(&models.Member{}).Where(query).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r permissionRepo) createCollectionPermission(collectionId uint64, studioId uint64, roleId uint64, isOverridden bool, permissionGroup string) (*models.CollectionPermission, error) {
	collectionPerms := &models.CollectionPermission{}
	collectionPerms.CollectionId = collectionId
	collectionPerms.RoleId = &roleId
	collectionPerms.StudioID = studioId
	collectionPerms.IsOverridden = isOverridden
	collectionPerms.PermissionGroup = permissionGroup
	err := r.db.Create(collectionPerms).Error
	if err != nil {
		return nil, err
	}
	return collectionPerms, nil
}

func (r permissionRepo) getStudioPermission(query map[string]interface{}) (models.StudioPermission, error) {
	var studioPerms models.StudioPermission
	err := postgres.GetDB().Model(&studioPerms).Where(query).First(&studioPerms).Error
	return studioPerms, err
}

func (r permissionRepo) GetRoleByID(roleID uint64) (*models.Role, error) {
	var role models.Role
	err := r.db.Model(&models.Role{}).Where("id = ?", roleID).Preload("Members").Find(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r permissionRepo) getCollectionPermissionsByID(query map[string]interface{}) ([]models.CollectionPermission, error) {
	var collectionPerms []models.CollectionPermission
	err := postgres.GetDB().Model(&collectionPerms).
		Where(query).
		Preload("Member").
		Find(&collectionPerms).Error
	return collectionPerms, err
}

func (r permissionRepo) getCollectionPermissions(query map[string]interface{}) ([]models.CollectionPermission, error) {
	var collectionPerms []models.CollectionPermission
	err := postgres.GetDB().Model(&collectionPerms).Where(query).Find(&collectionPerms).Error
	return collectionPerms, err
}

func (r permissionRepo) GetCollection(query map[string]interface{}) (*models.Collection, error) {
	var collection models.Collection
	err := r.db.Model(&models.Collection{}).Where(query).First(&collection).Error
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

type RoleMembersSerializer struct {
	MemberID   uint64 `json:"memberId"`
	Id         uint64 `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	IsSystem   bool   `json:"isSystem"`
	IsNonPerms bool   `json:"isNonPerms"`
}

func (r permissionRepo) GetMemberRolesByID(studioID uint64, memberID uint64) ([]RoleMembersSerializer, error) {
	var roles []RoleMembersSerializer
	result := postgres.GetDB().Raw(`
		select * from roles 
		left join role_members on role_members.role_id = roles.id
		where studio_id = ? and role_members.member_id = ?
	`, studioID, memberID).Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

func (r permissionRepo) CreateCollectionPermissionByMemberID(collectionId uint64, studioId uint64, memberId uint64, isOverridden bool, permissionGroup string, authUserId uint64) (*models.CollectionPermission, error) {
	collectionPerms := &models.CollectionPermission{}
	collectionPerms.CollectionId = collectionId
	if memberId != 0 {
		collectionPerms.MemberId = &memberId
	}
	collectionPerms.StudioID = studioId
	collectionPerms.IsOverridden = isOverridden
	collectionPerms.PermissionGroup = permissionGroup
	err := r.db.Create(&collectionPerms).Error
	if err != nil {
		return nil, err
	}

	var collectionPerm models.CollectionPermission
	err = postgres.GetDB().Model(&collectionPerm).Where("collection_id = ? and member_id = ?", collectionPerms.CollectionId, collectionPerms.MemberId).Preload("Member").Find(&collectionPerm).Error
	if err != nil {
		fmt.Println("Error on fetching collectionpermission", err.Error())
	}
	extraData := notifications.NotificationExtraData{
		CanvasRepoID: collectionPerms.CollectionId,
	}
	contentObject := models.COLLECTION
	member, _ := App.Repo.GetMember(map[string]interface{}{"id": memberId})
	notifications.App.Service.PublishNewNotification(notifications.CollectionInviteByName, authUserId, []uint64{member.UserID}, &collectionPerms.StudioID,
		nil, extraData, &collectionPerms.ID, &contentObject)

	return collectionPerms, nil
}

func (r permissionRepo) UpdateCollectionPermissions(query map[string]interface{}, collectionId uint64, studioId uint64, memberId uint64, roleId uint64, isOverridden bool, permissionGroup string, authUserId uint64) (*models.CollectionPermission, error) {
	var collectionPerms *models.CollectionPermission

	err := postgres.GetDB().Model(&collectionPerms).Where(query).First(&collectionPerms).Error

	// create flow
	if err != nil {
		if roleId == 0 {
			collectionPerms.RoleId = nil
		} else {
			collectionPerms.RoleId = &roleId
		}

		if memberId == 0 {
			collectionPerms.MemberId = nil
		} else {
			collectionPerms.MemberId = &memberId
		}
		collectionPerms.StudioID = studioId
		collectionPerms.CollectionId = collectionId
		collectionPerms.PermissionGroup = permissionGroup
	}

	if err == nil {
		if roleId != 0 {
			var role *models.Role
			r.db.Model(&role).Where("id = ?", roleId).First(&role)
			if !(role.Name == models.SYSTEM_ADMIN_ROLE) {
				collectionPerms.PermissionGroup = permissionGroup
			}
		} else {
			collectionPerms.PermissionGroup = permissionGroup
		}
	}

	// we are always updating these field if record if found or not found
	collectionPerms.IsOverridden = isOverridden

	// save will create a new record if it doesn't finds or update the existing record
	err2 := postgres.GetDB().Save(&collectionPerms).Error

	if err2 != nil {
		return nil, err2
	}
	if err != nil {
		// new record created we need to send preloaded data
		//err = postgres.GetDB().Model(&collectionPerms).Where(query).Preload("Member.User").First(&collectionPerms).Error
		// Add a notification of collection invite
		// For auto adding permission notification is not needed to send so commenting here.
		//extraData := notifications.NotificationExtraData{
		//	CollectionID: collectionPerms.CollectionId,
		//}
		//contentObject := models.COLLECTION
		//if collectionPerms.RoleId == nil {
		//	notifications.App.Service.PublishNewNotification(notifications.CollectionInviteByName, authUserId, []uint64{collectionPerms.Member.UserID}, &collectionPerms.StudioID,
		//		nil, extraData, &collectionPerms.CollectionId, &contentObject)
		//} else {
		//	notifications.App.Service.PublishNewNotification(notifications.CollectionInviteByGroup, authUserId, nil, &collectionPerms.StudioID,
		//		&[]uint64{*collectionPerms.RoleId}, extraData, &collectionPerms.CollectionId, &contentObject)
		//}
		//return collectionPerms, err
	}

	return collectionPerms, nil
}

func (r permissionRepo) GetBranchAccessToken(query map[string]interface{}) (*models.BranchAccessToken, error) {
	var instance *models.BranchAccessToken
	err := r.db.Model(&models.BranchAccessToken{}).Where(query).First(&instance).Error
	return instance, err
}

func (r permissionRepo) GetBranchWithRepo(query map[string]interface{}) (*models.CanvasBranch, error) {
	var branch models.CanvasBranch
	err := r.db.Model(&models.CanvasBranch{}).Where(query).Preload("CanvasRepository").First(&branch).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &branch, nil
}

func (r permissionRepo) CreateCollectionPermissionByMemberIDWithoutNotification(collectionId uint64, studioId uint64, memberId uint64, isOverridden bool, permissionGroup string, authUserId uint64) (*models.CollectionPermission, error) {
	collectionPerms := &models.CollectionPermission{}
	collectionPerms.CollectionId = collectionId
	if memberId != 0 {
		collectionPerms.MemberId = &memberId
	}
	collectionPerms.StudioID = studioId
	collectionPerms.IsOverridden = isOverridden
	collectionPerms.PermissionGroup = permissionGroup
	err := r.db.Create(&collectionPerms).Error
	if err != nil {
		return nil, err
	}
	return collectionPerms, nil
}

func (r permissionRepo) CreateCollectionPermissionByRoleIDWithoutNotification(collectionId uint64, studioId uint64, roleID uint64, isOverridden bool, permissionGroup string, authUserId uint64) (*models.CollectionPermission, error) {
	collectionPerms := &models.CollectionPermission{}
	collectionPerms.CollectionId = collectionId
	collectionPerms.RoleId = &roleID
	collectionPerms.StudioID = studioId
	collectionPerms.IsOverridden = isOverridden
	collectionPerms.PermissionGroup = permissionGroup
	err := r.db.Create(&collectionPerms).Error
	if err != nil {
		return nil, err
	}
	return collectionPerms, nil
}

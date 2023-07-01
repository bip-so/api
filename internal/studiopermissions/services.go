package studiopermissions

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

var (
	StudioPermissionService studioPermissionService
)

// Create a new studoio permissions object
//func (m studiopermissionService) MakeNewStudioPermission(studioId uint64, pg string, roleId uint64, memberId uint64, isOverride bool) *models.StudioPermission {
//	//emptyStudioPermsObject := StudioPermissionService.NewStudioPermission(studioId, pg, roleId, memberId, isOverride)
//	//return emptyStudioPermsObject
//}

func (m studioPermissionService) NewStudioPermission(studioId uint64, permsgroup string, roleId *uint64, memberId *uint64, isOverriddenFlag bool) (*models.StudioPermission, error) {
	//fmt.Println(emptyStudioPermsObject)
	sp := &models.StudioPermission{
		StudioID:        studioId,
		PermissionGroup: permsgroup,
		RoleId:          roleId,
		MemberId:        memberId,
		IsOverridden:    isOverriddenFlag,
	}
	repo := NewStudioPermissionsRepo()
	_, err := StudioPermissionRepo.CreateStudioPermission(repo, sp)
	if err != nil {
		return nil, err
	}

	return sp, nil
}

func (m studioPermissionService) UpdateStudioPermissions(query map[string]interface{}, body CreateStudioPermissionsPost, studioId uint64) (*models.StudioPermission, error) {

	repo := NewStudioPermissionsRepo()
	studioperms, err := repo.UpdateStudioPermissions(query, body, studioId)
	if err != nil {
		return nil, err
	}

	return studioperms, nil
}

func (m studioPermissionService) GetStudioPermissions(query map[string]interface{}) ([]models.StudioPermission, error) {

	repo := NewStudioPermissionsRepo()
	studioperms, err := repo.GetStudioPermissions(query)
	if err != nil {
		return nil, err
	}

	return studioperms, nil
}

// Function will return loggedIn users permissions_group string for this session on studios

/* Read Me

//var perms MAP
// Check is this user is Member of this studio.
// if not Member found
var perms = "pg_studio_none"
return perms

*/
func (m studioPermissionService) GetPemrs(userID uint64, studioID uint64) string {
	permissionMap := make(map[string]uint)
	fmt.Println(permissionMap)

	return ""
}

package role

import (
	"errors"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"

	"gitlab.com/phonepost/bip-be-platform/internal/member"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/supabase"
)

func (r roleService) NewRole(studioid uint64, userId uint64, roleName string, roleColor string, isSystemRole bool) *models.Role {

	return &models.Role{
		StudioID: studioid,
		Name:     roleName,
		Color:    roleColor,
		IsSystem: isSystemRole,
	}
}

func (r roleService) DefaultStudioAdminRole() *models.Role {

	return &models.Role{
		StudioID: 0,
		Name:     models.SYSTEM_ADMIN_ROLE,
		Color:    "#ffffff",
		IsSystem: true,
	}
}

func (r roleService) GetRole(Id uint64) (*models.Role, error) {

	repo := NewRoleRepo()
	role, err := repo.GetRole(Id)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r roleService) GetRolesByStudio(studioId uint64) ([]models.Role, error) {

	repo := NewRoleRepo()
	studioperms, err := repo.GetRolesByStudioID(studioId)
	if err != nil {
		return nil, err
	}

	return studioperms, nil
}

func (r roleService) CreateNewRole(studioId uint64, rp CreateRolePost) (*models.Role, error) {
	repo := NewRoleRepo()
	roleObject := &models.Role{
		StudioID: studioId,
		Name:     rp.Name,
		Color:    "#ffffff",
		IsSystem: false,
		Icon:     "",
	}

	createdRole, err := RoleRepo.CreateRole(repo, roleObject)
	if err != nil {
		return nil, err
	}

	return createdRole, nil
}

func (r roleService) DeleteRole(roleId uint64) error {
	repo := NewRoleRepo()
	err := RoleRepo.DeleteRole(repo, roleId)
	if err != nil {
		return err
	}
	return nil
}

func (r roleService) UpdateMembershipRole(ump UpdateManagementPost) ([]uint64, error) {

	// FE can send both memberID or UserId
	repo := NewRoleRepo()
	role, err := repo.GetRole(ump.RoleId)

	if err != nil {
		return nil, err
	}

	studioMembers, err := member.App.Service.GetMembersByStudio(role.StudioID, 0)

	if err != nil {
		return nil, err
	}

	members, err := repo.FindMembersInRole(role)
	if err != nil {
		return nil, err
	}

	addMemberIDs := []uint64{}
	addMembers := []models.Member{}
	updateMemberIDs := []uint64{}
	for _, mID := range ump.MembersAdded {
		found := false
		for _, eMember := range members {
			if mID == eMember.ID || mID == eMember.UserID {
				found = true
				break
			}
		}
		if !found {
			// given ID not found in role
			// check if the member exists, if not create a new member
			// given ID can be userId, or memberId
			exists := false
			var memberID uint64
			for _, memb := range studioMembers {
				if mID == memb.ID || mID == memb.UserID {
					exists = true
					memberID = memb.ID
					break
				}
			}
			if !exists {
				//create new member
				NewMember := queries.App.MemberQuery.AddUserIDToStudio(mID, role.StudioID)
				memberID = NewMember.ID
				if memberID == 0 {
					return nil, errors.New("couldn't add user to studio")
				}
			}
			addMembers = append(addMembers, models.Member{BaseModel: models.BaseModel{ID: memberID}})
			addMemberIDs = append(addMemberIDs, memberID)
			updateMemberIDs = append(updateMemberIDs, memberID)
		}
	}

	removeMembers := []models.Member{}
	if role.Name != models.SYSTEM_ADMIN_ROLE || (role.Name == models.SYSTEM_ADMIN_ROLE && len(members) > 1) {
		for _, mID := range ump.MembersRemoved {
			found := false
			var memberID uint64
			for _, eMember := range members {
				if mID == eMember.ID || mID == eMember.UserID {
					found = true
					memberID = eMember.ID
					break
				}
			}
			if found {
				removeMembers = append(removeMembers, models.Member{BaseModel: models.BaseModel{ID: memberID}})
				updateMemberIDs = append(updateMemberIDs, memberID)
			}
		}
	}

	if len(addMembers) > 0 {
		err = repo.AddMembersInRole(addMembers, role)

		if err != nil {
			return nil, err
		}
	}

	if len(removeMembers) > 0 {
		err = repo.RemoveMembersInRole(removeMembers, role)
	}

	// Invalidating the user cache here.
	go func() {
		memberInstances, err := queries.App.MemberQuery.GetMultipleMembersByID(updateMemberIDs)
		if err != nil {
			logger.Error(err.Error())
		}
		for _, memberInstance := range memberInstances {
			// @todo later we can update supabase and delete user associated studios only on some specific roles like Administrator, Member
			supabase.UpdateUserSupabase(memberInstance.UserID, true)
			queries.App.StudioQueries.DeleteUserAssociatedStudioDataByUserID(memberInstance.UserID)
			permissions.App.Service.InvalidateUserPermissionCache(memberInstance.UserID, role.StudioID)
		}
	}()
	return addMemberIDs, err
}

func (r roleService) UpdateRole(rp UpdateRolePost) error {
	repo := NewRoleRepo()
	updates := map[string]interface{}{
		"name": rp.Name,
	}

	err := repo.UpdateRole(rp.RoleId, updates)
	return err
}

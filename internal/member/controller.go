package member

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
)

func (c memberController) CreateStudioMemberController(userID uint64, stdID uint64) uint64 {
	//memberId := App.Service.AddMemberToStudio(userID, stdID)
	member := queries.App.MemberQuery.AddUserIDToStudio(userID, stdID)
	if member == nil {
		logger.Info("Failed to create a member while creating studio. ")
		return 0
	}
	return member.ID
}

func (c memberController) GetMemberInstance(id uint64) *models.Member {
	member, err := queries.App.MemberQuery.GetMember(map[string]interface{}{"id": id})
	if err != nil {
		return nil
	}
	return member
}

func (c memberController) CanvasBranchMembersAndRoles(canvasBranchID uint64) (*[]MembersCanvasBranch, error) {
	perms, err := queries.App.CanvasBranchPermissionQuery.GetCanvasBranchPermissions(map[string]interface{}{"canvas_branch_id": canvasBranchID})
	if err != nil {
		return nil, err
	}
	canvasBranchMembers := &[]MembersCanvasBranch{}
	for _, perm := range perms {
		if perm.Member != nil {
			branchMember := &MembersCanvasBranch{
				ID:              perm.ID,
				Type:            "member",
				PermissionGroup: perm.PermissionGroup,
				MemberID:        *perm.MemberId,
				User:            SerializeBranchUser(*perm.Member.User),
			}
			*canvasBranchMembers = append(*canvasBranchMembers, *branchMember)
		} else if perm.Role != nil {
			branchMember := &MembersCanvasBranch{
				ID:              perm.ID,
				Type:            "role",
				PermissionGroup: perm.PermissionGroup,
				RoleID:          *perm.RoleId,
				Role:            SerializeBranchRole(*perm.Role),
			}
			*canvasBranchMembers = append(*canvasBranchMembers, *branchMember)
		}
	}
	return canvasBranchMembers, nil
}

func (c memberController) RoleMembers(roleID uint64, skip, limit int) ([]BranchUserSerializer, error) {
	roleMembers, err := queries.App.RoleQuery.GetRoleMembers(roleID, skip, limit)
	if err != nil {
		return nil, err
	}
	memberIDs := []uint64{}
	for _, roleMember := range roleMembers {
		memberIDs = append(memberIDs, roleMember.MemberID)
	}
	members, _ := queries.App.MemberQuery.GetMembersWithPreload(map[string]interface{}{"id": memberIDs})
	users := []BranchUserSerializer{}
	for _, member := range members {
		users = append(users, SerializeBranchUser(*member.User))
	}
	return users, nil
}

func (c memberController) RoleMembersSearch(search string, roleID uint64, skip, limit int) ([]BranchUserSerializer, error) {
	roleMembers, err := queries.App.MemberQuery.RoleMembersSearch(search, roleID, skip, limit)
	if err != nil {
		return nil, err
	}
	memberIDs := []uint64{}
	for _, roleMember := range roleMembers {
		memberIDs = append(memberIDs, roleMember.MemberID)
	}
	members, _ := queries.App.MemberQuery.GetMembersWithPreload(map[string]interface{}{"id": memberIDs})
	users := []BranchUserSerializer{}
	for _, member := range members {
		users = append(users, SerializeBranchUser(*member.User))
	}
	return users, nil
}

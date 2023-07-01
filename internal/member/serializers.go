package member

import (
	"gitlab.com/phonepost/bip-be-platform/internal/follow"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/user"
)

type MemberSerializer struct {
	ID        uint64                  `json:"id"`
	UUID      string                  `json:"uuid"`
	UserID    uint64                  `json:"userID"`
	StudioID  uint64                  `json:"studioID"`
	IsRemoved bool                    `json:"isRemoved"`
	User      user.UserMiniSerializer `json:"user"`
}

func SerializeMember(member *models.Member) *MemberSerializer {
	view := &MemberSerializer{
		ID:        member.ID,
		UUID:      member.UUID.String(),
		UserID:    member.UserID,
		StudioID:  member.StudioID,
		IsRemoved: member.IsRemoved,
	}
	if member.User != nil {
		view.User = user.UserMiniSerializerData(member.User)
		resp, _ := follow.App.Controller.GetUserFollowFollowCountHandler(view.UserID)
		view.User.Followers = resp.Followers
		view.User.Following = resp.Following
	}
	return view
}

type MembersCanvasBranch struct {
	ID              uint64               `json:"id"`
	Type            string               `json:"type"`
	PermissionGroup string               `json:"permissionGroup"`
	MemberID        uint64               `json:"memberId"`
	RoleID          uint64               `json:"roleId"`
	User            BranchUserSerializer `json:"user"`
	Role            BranchRoleSerializer `json:"role"`
}

type BranchUserSerializer struct {
	Id        uint64 `json:"id"`
	UUID      string `json:"uuid"`
	FullName  string `json:"fullName"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatarUrl"`
}

func SerializeBranchUser(user models.User) BranchUserSerializer {
	return BranchUserSerializer{
		Id:        user.ID,
		UUID:      user.UUID.String(),
		FullName:  user.FullName,
		Username:  user.Username,
		AvatarUrl: user.AvatarUrl,
	}
}

type BranchRoleSerializer struct {
	Id           uint64 `json:"id"`
	UUID         string `json:"uuid"`
	Name         string `json:"name"`
	MembersCount int    `json:"membersCount"`
}

func SerializeBranchRole(role models.Role) BranchRoleSerializer {
	branchRole := BranchRoleSerializer{
		Id:           role.ID,
		UUID:         role.UUID.String(),
		Name:         role.Name,
		MembersCount: len(role.Members),
	}
	return branchRole
}

func BulkSerializeMembers(members []models.Member) []MemberSerializer {
	views := []MemberSerializer{}
	for _, member := range members {
		views = append(views, *SerializeMember(&member))
	}
	return views
}

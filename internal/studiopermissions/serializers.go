package studiopermissions

import (
	"gitlab.com/phonepost/bip-be-platform/internal/member"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/role"
	"gitlab.com/phonepost/bip-be-platform/internal/studio"
)

type StudioPermissionsSerializer struct {
	ID              uint64                   `json:"id"`
	UUID            string                   `json:"uuid"`
	StudioID        uint64                   `json:"studioID"`
	PermissionGroup string                   `json:"permissionGroup"`
	RoleID          *uint64                  `json:"roleID"`
	MemberID        *uint64                  `json:"memberID"`
	IsOverridden    bool                     `json:"isOverridden"`
	Role            *role.RoleSerializer     `json:"role"`
	Studio          *studio.StudioSerializer `json:"studio"`
	Member          *member.MemberSerializer `json:"member"`
}

func SerializeStudioPermission(studioperm *models.StudioPermission) *StudioPermissionsSerializer {
	view := StudioPermissionsSerializer{
		ID:              studioperm.ID,
		UUID:            studioperm.UUID.String(),
		StudioID:        studioperm.StudioID,
		PermissionGroup: studioperm.PermissionGroup,
		RoleID:          studioperm.RoleId,
		MemberID:        studioperm.MemberId,
		IsOverridden:    studioperm.IsOverridden,
	}
	if studioperm.Member != nil {
		view.Member = member.SerializeMember(studioperm.Member)
	}
	if studioperm.Role != nil {
		view.Role = role.SerializeRole(studioperm.Role)
	}
	if studioperm.Studio != nil {
		view.Studio = studio.SerializeStudio(studioperm.Studio)
	}
	return &view
}

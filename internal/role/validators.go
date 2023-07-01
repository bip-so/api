package role

type CreateRolePost struct {
	Name  string `json:"name" binding:"required"` // Role Name
	Color string `json:"color"`                   // Color, we'll use #ffffff by default
	Icon  string `json:"icon"`                    // icons for future
}

type DeleteRolePost struct {
	RoleID uint64 `json:"role_id" binding:"required"`
}

type UpdateManagementPost struct {
	RoleId         uint64   `json:"roleId"`
	MembersAdded   []uint64 `json:"membersAdded"`   // array id of MemberIDs/UserIDs to be added
	MembersRemoved []uint64 `json:"membersRemoved"` // array id of MemberIDs/userIDs to be removed
}

type UpdateRolePost struct {
	RoleId uint64 `json:"roleId" binding:"required"`
	Name   string `json:"name" binding:"required"` // Role Name
}

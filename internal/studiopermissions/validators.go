package studiopermissions

type CreateStudioPermissionsPost struct {
	PermsGroup       string `json:"permsGroup"`
	RoleId           uint64 `json:"roleId"`
	MemberId         uint64 `json:"memberId"`
	IsOverriddenFlag bool   `json:"isOverriddenFlag" default:"false"`
}

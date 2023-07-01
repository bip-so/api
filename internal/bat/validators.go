package bat

type CreateAccessTokenPost struct {
	PermissionGroup string `json:"permissionGroup"`
}
type CreateEmailInvite struct {
	Invites []EmailInvitePerEmail `json:"invites"`
}
type EmailInvitePerEmail struct {
	Email                  string `json:"email"`
	CanvasPermissionsGroup string `json:"canvasPermissionsGroup"`
}
type PlaceHolder struct {
}

package member

type CreateStudioMemberValidator struct {
	UserID   uint64
	StudioID uint64
}

type NewDiscordMember struct {
	StudioId uint64
	UserId   uint64
}

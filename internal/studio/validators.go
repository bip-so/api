package studio

type CreateStudioValidator struct {
	Name        string   `json:"name"`
	Handle      string   `json:"handle"`
	Description string   `json:"description"`
	Website     string   `json:"website"`
	Topics      []string `json:"topics"`
}

type UpdateStudioValidator struct {
	Name        string   `json:"name"`
	Handle      string   `json:"handle"`
	Description string   `json:"description"`
	ImageURL    string   `json:"-"`
	Website     string   `json:"website"`
	Topics      []string `json:"topics"`
}

type BanUserValidator struct {
	UserID    uint64 `json:"userId"`
	BanReason string `json:"banReason"`
}

type JoinStudioBulkPost struct {
	UsersAdded []uint64 `json:"usersAdded"` // array id of UserIDs to be added
}

// Empty Post
type StudioMembershipRequestNew struct {
	//UserID uint64 `json:"userId"`
}

// Empty Post
type StudioMembershipRequestReject struct {
	//UserID uint64 `json:"userId"`
}

type RegisterNewStudioIntegrationValidator struct {
	GuildId   string `json:"guildId" binding:"required"`
	GuildName string `json:"guildName" binding:"required"`
}

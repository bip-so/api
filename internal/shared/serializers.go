package shared

type CommonUserMiniSerializer struct {
	Id        uint64 `json:"id"`
	UUID      string `json:"uuid"`
	FullName  string `json:"fullName"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatarUrl"`
}

type ReactedCount struct {
	Emoji   string `json:"emoji"`
	Count   int    `json:"count"`
	Reacted bool   `json:"reacted"`
}

type CommonStudioMiniSerializer struct {
	ID                    uint64 `json:"id"`
	UUID                  string `json:"uuid"`
	DisplayName           string `json:"displayName"`
	Handle                string `json:"handle"`
	ImageURL              string `json:"imageUrl"`
	CreatedByID           uint64 `json:"createdById"`
	AllowPublicMembership bool   `json:"allowPublicMembership"`
	IsRequested           bool   `json:"isRequested"`
}

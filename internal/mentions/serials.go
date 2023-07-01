package mentions

type MentionedUserSerializer struct {
	Type      string `json:"type"`
	ID        uint64 `json:"id"`
	UUID      string `json:"uuid"`
	FullName  string `json:"fullName"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatarUrl"`
}

type MentionedCanvasSerializer struct {
	Type     string `json:"type"`
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	UUID     string `json:"uuid"`
	Key      string `json:"key"`
	RepoID   uint64 `json:"repoID"`
	RepoKey  string `json:"repoKey"`
	RepoName string `json:"repoName"`
	RepoUUID string `json:"repoUUID"`
}
type MentionedRolesSerializer struct {
	Type string `json:"type"`
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

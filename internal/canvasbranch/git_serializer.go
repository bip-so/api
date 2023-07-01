package canvasbranch

import (
	"gitlab.com/phonepost/bip-be-platform/internal/shared"
	"time"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	stores "gitlab.com/phonepost/bip-be-platform/pkg/stores/git"
)

type GitLogView struct {
	CommitID  string                          `json:"commitId"`
	Message   string                          `json:"message"`
	User      shared.CommonUserMiniSerializer `json:"user"`
	CreatedAt time.Time                       `json:"createdAt"`
}

func GitLogsSerializerData(logs []*stores.GitLog, users *[]models.User) *[]GitLogView {

	gitLogViews := []GitLogView{}

	userMap := map[string]*models.User{}

	for i := range *users {
		user := (*users)[i]
		userMap[user.UUID.String()] = &user
	}

	for _, log := range logs {
		gitLogViews = append(gitLogViews, GitLogView{
			CommitID: log.ID,
			Message:  log.Message,
			User: shared.CommonUserMiniSerializer{
				Id:        userMap[log.UserID].ID,
				UUID:      log.UserID,
				Username:  userMap[log.UserID].Username,
				FullName:  userMap[log.UserID].FullName,
				AvatarUrl: userMap[log.UserID].AvatarUrl,
			},
			CreatedAt: log.CreatedAt,
		})
	}

	return &gitLogViews
}

type GitAttributionView struct {
	Edits int                             `json:"edits"`
	User  shared.CommonUserMiniSerializer `json:"user"`
}

func GitAttributionsSerializerData(attrs *[]models.Attribution) *[]GitAttributionView {

	gitAttrViews := []GitAttributionView{}

	for _, attr := range *attrs {
		gitAttrViews = append(gitAttrViews, GitAttributionView{
			Edits: attr.Edits,
			User: shared.CommonUserMiniSerializer{
				Id:        attr.User.ID,
				UUID:      attr.User.UUID.String(),
				Username:  attr.User.Username,
				FullName:  attr.User.FullName,
				AvatarUrl: attr.User.AvatarUrl,
			},
		})
	}

	return &gitAttrViews
}

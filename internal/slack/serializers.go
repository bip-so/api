package slack2

import (
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"gorm.io/datatypes"
)

type UserGetSerializer struct {
	Id              uint64              `json:"id"`
	UUID            string              `json:"uuid"`
	FirstName       string              `json:"firstName"`
	LastName        string              `json:"lastName"`
	FullName        string              `json:"fullName"`
	Username        string              `json:"username"`
	HasEmail        bool                `json:"hasEmail"`
	IsSuperuser     bool                `json:"isSuperuser"`
	IsSetupDone     bool                `json:"isSetupDone"`
	IsEmailVerified bool                `json:"isEmailVerified"`
	AvatarUrl       string              `json:"avatarUrl"`
	Followers       uint64              `json:"followers"`
	Following       uint64              `json:"following"`
	IsFollowing     *bool               `json:"isFollowing"`
	UserProfile     *models.UserProfile `json:"userProfile"`
	CreatedAt       time.Time           `json:"createdAt"`
	UpdatedAt       time.Time           `json:"updatedAt"`
}

func UserGetSerializerData(user *models.User) UserGetSerializer {

	return UserGetSerializer{
		Id:              user.ID,
		UUID:            user.UUID.String(),
		Username:        user.Username,
		FullName:        user.FullName,
		HasEmail:        user.Email.Valid,
		AvatarUrl:       user.AvatarUrl,
		IsSuperuser:     user.IsSuperuser,
		IsSetupDone:     user.IsSetupDone,
		IsEmailVerified: user.IsEmailVerified,
		CreatedAt:       user.CreatedAt,
		UserProfile:     user.UserProfile,
		UpdatedAt:       user.UpdatedAt,
	}
}

type AccountView struct {
	UserGetSerializer
	ProviderID   string                 `json:"providerID"`
	ProviderName string                 `json:"providerName"`
	Metadata     datatypes.JSON         `json:"metadata"`
	Extras       map[string]interface{} `json:"extras"`
}

func RenderAccountWithExtras(c *gin.Context, user *models.UserSocialAuth, extras map[string]interface{}) {
	accountView := AccountView{
		UserGetSerializer: UserGetSerializerData(user.User),
		Extras:            extras,
		Metadata:          user.Metadata,
		ProviderID:        user.ProviderID,
		ProviderName:      user.ProviderName,
	}
	response.RenderResponse(c, accountView)
}

type SlackMessagePayload struct {
	Text      string `json:"text"`
	ChannelID string `json:"channel"`
	ThreadTs  string `json:"thread_ts"`
}

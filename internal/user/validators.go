package user

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type UpdateUserValidator struct {
	FullName    string               `json:"fullName" form:"fullName"`
	File        multipart.FileHeader `json:"file" form:"file"`
	Bio         string               `json:"bio" form:"bio"`
	Username    string               `json:"username" form:"username"`
	AvatarUrl   string               `json:"avatarUrl" form:"avatarUrl"`
	TwitterUrl  string               `json:"twitterUrl" form:"twitterUrl"`
	IsSetupDone bool                 `json:"isSetupDone"`
	Website     string               `json:"website" form:"website"`
	Location    string               `json:"location" form:"location"`
}

type UpdateUserSettingsValidator struct {
	Type                    string `json:"type" binding:"required"`
	AllComments             bool   `json:"allComments" binding:"required"`
	RepliesToMe             bool   `json:"repliesToMe" binding:"required"`
	Mentions                bool   `json:"mentions" binding:"required"`
	Reactions               bool   `json:"reactions" binding:"required"`
	Invite                  bool   `json:"invite" binding:"required"`
	FollowedMe              bool   `json:"followedMe" binding:"required"`
	FollowedMyStudio        bool   `json:"followedMyStudio" binding:"required"`
	PublishAndMergeRequests bool   `json:"publishAndMergeRequests" binding:"required"`
	ResponseToMyRequests    bool   `json:"responseToMyRequests" binding:"required"`
	SystemNotifications     bool   `json:"systemNotifications" binding:"required"`
	DarkMode                bool   `json:"darkMode" binding:"required"`
}

type PatchUserSettingsValidator struct {
	Data []UpdateUserSettingsValidator `json:"data"`
}

type UserSearchParams struct {
	Text      string
	ProductId string
}

func UserSearchValidator(c *gin.Context) UserSearchParams {
	query := c.Request.URL.Query()
	return UserSearchParams{
		Text:      query.Get("text"),
		ProductId: query.Get("productId"),
	}
}

type UserUpdateFollowers struct {
	FollowerID uint64
}

type FollowUserFollowCountResponse struct {
	Followers uint64 `json:"followers"`
	Following uint64 `json:"following"`
}

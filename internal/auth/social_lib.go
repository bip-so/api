package auth

import (
	"errors"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

// Login with Discord, Twitter, Slack
type SocialAuthPost struct {
	Provider                   string `json:"provider" binding:"required"`
	ProviderID                 string `json:"providerID" binding:"required"`
	Image                      string `json:"image" binding:"required"`
	FullName                   string `json:"fullName" binding:"required"`
	Email                      string `json:"email"`
	UserName                   string `json:"userName"`
	AccessToken                string `json:"access_token"`
	UserViaDiscordAppDirectory bool   `json:"userViaDiscordAppDirectory"`
	ClientReferenceId          string `json:"clientReferenceId"`
}

func (obj SocialAuthPost) Validate() error {
	allowedScope := []string{"discord", "twitter", "slack"}
	if !utils.SliceContainsItem(allowedScope, obj.Provider) {
		return errors.New("Please only send  \"discord\",\"twitter\", \"slack\"")
	}
	return nil
}

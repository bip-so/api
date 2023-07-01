package auth

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"time"
)

type AuthLegacySerializer struct {
	c *gin.Context
}

type AuthSerializerWithTokens struct {
	Id                uint64    `json:"id"`
	UUID              string    `json:"uuid"`
	FullName          string    `json:"fullName"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	RefreshToken      string    `json:"refreshToken"`
	AccessToken       string    `json:"accessToken"`
	AccessTokenId     string    `json:"accessTokenId"`
	RefreshTokenId    string    `json:"refreshTokenId"`
	IsSocialLogin     bool      `json:"isSocialLogin"`
	LastLogin         time.Time `json:"lastLogin"`
	UserProfile       string    `json:"userProfile"`
	IsSetupDone       bool      `json:"isSetupDone"`
	IsEmailVerified   bool      `json:"isEmailVerified"`
	AvatarUrl         string    `json:"avatarUrl"`
	DefaultStudioID   uint64    `json:"defaultStudioID"`
	ViaDiscordStore   bool      `json:"viaDiscordStore"`
	IsNewUser         bool      `json:"isNewUser"`
	IsSignUp          bool      `json:"isSignUp"`
	ClientReferenceId string    `json:"clientReferenceId"`
}

// Get the tokens
func (self *AuthLegacySerializer) GetOtpDetails(otpObject ResponseNewOtp) ResponseNewOtp {
	return ResponseNewOtp{
		Email: otpObject.Email,
		//Otp:   otpObject.Otp,
	}
}

// TokenDetails is Token.go
func (self *AuthLegacySerializer) GetRefreshToken(tokens TokenDetails) TokenDetails {

	return TokenDetails{
		AccessToken:    tokens.AccessToken,
		RefreshToken:   tokens.RefreshToken,
		AccessTokenID:  tokens.AccessTokenID,
		RefreshTokenID: tokens.AccessTokenID,
	}
}

func (self *AuthLegacySerializer) GetAuthSerializerWithTokens(tokens TokenDetails) AuthSerializerWithTokens {
	user := self.c.MustGet("currentUser").(*models.User)
	return AuthSerializerWithTokens{
		Id:                user.ID,
		UUID:              user.UUID.String(),
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email.String,
		AccessToken:       tokens.AccessToken,
		RefreshToken:      tokens.RefreshToken,
		AccessTokenId:     tokens.AccessTokenID,
		RefreshTokenId:    tokens.RefreshTokenID,
		IsSocialLogin:     false,
		IsSetupDone:       user.IsSetupDone,
		IsEmailVerified:   user.IsEmailVerified,
		AvatarUrl:         user.AvatarUrl,
		LastLogin:         *user.LastLogin,
		DefaultStudioID:   user.DefaultStudioID,
		ViaDiscordStore:   user.ViaDiscordStore,
		ClientReferenceId: user.ClientReferenceId,
	}

}

func (self *AuthLegacySerializer) GetAuthSerializerWithTokensAndSocials(tokens TokenDetails) AuthSerializerWithTokens {
	user := self.c.MustGet("currentUser").(*models.User)
	return AuthSerializerWithTokens{
		Id:                user.ID,
		UUID:              user.UUID.String(),
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email.String,
		AccessToken:       tokens.AccessToken,
		RefreshToken:      tokens.RefreshToken,
		AccessTokenId:     tokens.AccessTokenID,
		RefreshTokenId:    tokens.RefreshTokenID,
		IsSocialLogin:     true,
		IsSetupDone:       user.IsSetupDone,
		IsEmailVerified:   user.IsEmailVerified,
		AvatarUrl:         user.AvatarUrl,
		DefaultStudioID:   user.DefaultStudioID,
		ViaDiscordStore:   user.ViaDiscordStore,
		ClientReferenceId: user.ClientReferenceId,
	}
}

type AuthSocialSerializer struct {
	c *gin.Context
}
type AuthSocialSerializerToken struct {
	Id   uint64 `json:"id"`
	UUID string `json:"uuid"`

	FullName          string `json:"fullName"`
	Username          string `json:"username"`
	Email             string `json:"email"`
	RefreshToken      string `json:"refreshToken"`
	AccessToken       string `json:"accessToken"`
	IsSocialLogin     bool   `json:"isSocialLogin"`
	LastLogin         string `json:"lastLogin"`
	UserProfile       string `json:"userProfile"`
	AvatarUrl         string `json:"avatarUrl"`
	DefaultStudioID   uint64 `json:"defaultStudioID"`
	ViaDiscordStore   bool   `json:"viaDiscordStore"`
	ClientReferenceId string `json:"clientReferenceId"`
}

func (self *AuthSocialSerializer) GetAuthSocialSerializer(accessToken string, refreshToken string) AuthSocialSerializerToken {
	user := self.c.MustGet("currentUser").(*models.User)
	return AuthSocialSerializerToken{
		Id:                user.ID,
		UUID:              user.UUID.String(),
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email.String,
		AccessToken:       accessToken,
		RefreshToken:      refreshToken,
		IsSocialLogin:     true,
		AvatarUrl:         user.AvatarUrl,
		DefaultStudioID:   user.DefaultStudioID,
		ViaDiscordStore:   user.ViaDiscordStore,
		ClientReferenceId: user.ClientReferenceId,
	}
}

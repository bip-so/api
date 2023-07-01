package auth

import (
	"fmt"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type expectedLoginPostData struct {
	Account      string `json:"username" binding:"required"` // Can be Email or Username
	Password     string `json:"password"`
	Otp          string `json:"otp"`
	IsGhostLogin bool   `json:"isGhostLogin"`
	GhostToken   string `json:"ghostToken"`
	GhostSecret  string `json:"ghostSecret"`
}

type expectedRefreshTokenPostData struct {
	//AccessToken   string `json:"accessToken"  binding:"required"` // Can be Email or Username
	RefreshToken  string `json:"refreshToken"  binding:"required" `
	AccessTokenID string `json:"AccessTokenID" binding:"required"`
}

type RefreshTokenResponseObject struct {
}

type expectedLegacySignupPostData struct {
	//Username              string `json:"username"` // Can be Email or Username
	Email             string `json:"email"`
	Password          string `json:"password"`
	ClientReferenceId string `json:"clientReferenceId"`
	//IsSocial              bool   `json:"is_social"`
	//SocialProvider        string `json:"social_provider"` // providers = ("discord", "slack", "twitter")
	//SocialProviderID      string `json:"social_provider_id"`
	//SocialProviderMetadat string `json:"social_provider_metadata"`
}

type expectedExistingEmailBody struct {
	Email string `json:"email" binding:"required"`
}

type ResetPasswordData struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
	//NewPassword string `json:"newPassword"`
	//OldPassword string `json:"oldPassword"`
}

type forgotPasswordData struct {
	Email string `json:"email" binding:"required"`
	//NewPassword string `json:"newPassword"`
	//OldPassword string `json:"oldPassword"`
}

type SpecialLogin struct {
	Email string `json:"email" binding:"required"`
}

type changePasswordData struct {
	NewPassword string `json:"newPassword" binding:"required"`
	OldPassword string `json:"oldPassword" binding:"required"`
}

func signupLegacyPostValidation(email string, password string) error {
	var user *models.User
	if password == "" || email == "" {
		return AuthErrorsEmailOrPasswordAreEmpty
	}

	if len(password) < 6 {
		return AuthErrorsPasswordCharLimit
	}

	isEmailValidFlag := utils.IsEmailValid(email)
	if !isEmailValidFlag {
		return AuthErrorsEmailNotValid
	}
	// check email uniqueness

	user, _ = App.Service.getUserWithEmail(email)
	if user != nil {
		return AuthErrorsUserExisitsWithThisEmail
	}
	return nil
}

func loginValidations(username string, password string) error {
	var err error
	var user *models.User

	if password == "" || username == "" {
		return AuthErrorsEmailOrPasswordAreEmpty
	}

	if len(password) < 6 {
		return AuthErrorsPasswordCharLimit
	}

	isEmailValidFlag := utils.IsEmailValid(username)
	if !isEmailValidFlag {
		return AuthErrorsEmailNotValid
	}
	user, err = App.Service.getUserWithEmail(username)
	fmt.Println(user.Password)
	if user == nil {
		return AuthErrorsUsernotFound
	} else {
		// if user was found checking the passwords.
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			fmt.Println(err)
			return AuthErrorsUserPasswordNotMatch
		}
	}
	return nil
}

type expectedExistingUsernameBody struct {
	Username string `json:"username" binding:"required"`
}

// generateUserOtp - Post
type PostNewOtpRequest struct {
	Email string `json:"email" binding:"required,email" `
}

// generateUserOtp - Response
type ResponseNewOtp struct {
	Otp   string
	Email string
}

type NewGhostToken struct {
	Token string
	Email string
}
type ResponseForgotPass struct {
	Otp   string
	Email string
}

func loginValidationsOtp(username string, otp string) error {
	var user *models.User

	if otp == "" || username == "" {
		return AuthErrorsEmailOrOtpAreEmpty
	}

	isEmailValidFlag := utils.IsEmailValid(username)
	if !isEmailValidFlag {
		return AuthErrorsEmailNotValid
	}
	user, _ = App.Service.getUserWithEmail(username)
	if user == nil {
		return AuthErrorsUsernotFound
	}
	return nil
}

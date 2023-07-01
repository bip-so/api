package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	user2 "gitlab.com/phonepost/bip-be-platform/internal/user"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/supabase"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

// EmailVerificationFlow: We are sharing this as an email link.
// We check is the token is correct or not.
func (h *authController) EmailVerificationFlow(token string) bool {
	// updateUserEmailConfirmationFlag
	user, _ := queries.App.UserQueries.GetUser(map[string]interface{}{"uuid": token})
	fmt.Println(user)
	if user == nil {
		return false
	} else {
		// only update when we actually need to so on False only.
		if !user.IsEmailVerified {
			queries.App.UserQueries.UpdateUser(user.ID, map[string]interface{}{"is_email_verified": true}, false)
			go func() {
				supabase.UpdateUserEmailIsVerifiedSupabase(user.ID)
			}()
		}
		// Verifying the Emails
		return true
	}
}

func (h *authController) signupLegacy(postData expectedLegacySignupPostData) (*models.User, error) {
	// unpacking
	email := strings.ToLower(postData.Email)
	password := postData.Password

	var err error
	var user *models.User

	// basic validation
	err = signupLegacyPostValidation(email, password)
	if err != nil {
		return nil, err
	}

	// validation for username/stusio name check
	// discuss with Paras.
	// global.GlobalController.CheckHandleAvailable()

	user, err = TryLegacySignup(email, password, postData.ClientReferenceId)

	if err != nil {
		return nil, err
	}
	return user, nil

}

func (h *authController) loginLegacy(username string, password string) (*models.User, error) {

	var err error
	var user *models.User

	err = loginValidations(username, password)
	if err != nil {
		return nil, err
	}

	user, _ = App.Service.getUserWithEmail(username)
	return user, err

}

func (h *authController) resetPasswordWithToken(email string, newPassword string) error {
	user, _ := queries.App.UserQueries.GetUser(map[string]interface{}{"email": email})
	var err error
	// Update or create Password
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	err = user2.App.Repo.UpdatePassword(*user, string(newHashedPassword))
	App.Service.PasswordChangedMailer(user)
	return err
}

func (h *authController) resetPassword(user models.User, oldPassword string, newPassword string) error {
	var err error
	if user.HasPassword {
		if oldPassword == "" || newPassword == "" {
			return errors.New("passwords is empty")
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
		if err != nil {
			return AuthErrorsUserPasswordNotMatch
		}
	}
	// Update or create Password
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	err = user2.App.Repo.UpdatePassword(user, string(newHashedPassword))
	App.Service.PasswordChangedMailer(&user)
	return err

}

func (h *authController) checkExistingEmail(email string) bool {
	var err error

	_, err = App.Service.getUserWithEmail(email)
	if err != nil {
		return false
	} else {
		return true
	}

}

func (h *authController) checkExistingUsername(email string) bool {
	var err error

	_, err = App.Service.getUserWithUserName(email)
	if err != nil {
		return false
	} else {
		return true
	}

}

// Legacy login with OTP
func (h *authController) loginLegacyWithOtp(username string, otp string) (*models.User, error) {

	var err error
	var user *models.User

	err = loginValidationsOtp(username, otp)
	if err != nil {
		return nil, err
	}

	// Need to Get and Check the OTP

	flg := App.Service.IsValidOtp(username, otp)
	fmt.Println("AuthService.IsValidOtp")
	fmt.Println(flg)

	if flg == false {
		return nil, errors.New("Token did not match")
	}

	user, _ = App.Service.getUserWithEmail(username)
	return user, err

}

// Login with Token Only
func (h *authController) loginLegacyWithToken(data expectedLoginPostData) (*models.User, error) {
	var err error
	var user *models.User
	secret := "XRRZGWQFCWKHFHHDBYWYT"
	rc := redis.RedisClient()
	ctx := redis.GetBgContext()
	val, err := rc.Get(ctx, models.GhostLogin+data.Account).Result()
	if err == nil {
		return nil, err
	}
	data2 := NewGhostToken{}
	_ = json.Unmarshal([]byte(val), &data2)

	if data.GhostToken == data2.Token && data.GhostSecret == secret {
		user, _ = App.Service.getUserWithEmail(data.Account)
		return user, err
	}
	return nil, errors.New("Lol, not OK ")
}

func (h *authController) LogoutProper(token string) {
	rc := redis.RedisClient()
	ctx := redis.GetBgContext()
	_ = rc.Set(ctx, "token-sessions:"+token, "false", time.Hour*24*7*30).Err()
}

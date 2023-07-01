package auth

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	user2 "gitlab.com/phonepost/bip-be-platform/internal/user"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/kafka"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// We need to ve
func (r *authRoutes) VerifyUserEmail(c *gin.Context) {
	//basePath, _ := c.Get("appBaseUrl")
	emailVerificationToken := c.Param("verificationKey")
	// Check if the token is valid
	tokenStatus := App.Controller.EmailVerificationFlow(emailVerificationToken)
	// if the token is found we need to update the user instance
	if tokenStatus {
		c.JSON(http.StatusOK, gin.H{"message": "Your email has been verified."})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Your Token is used or expired. Please request for new verification emails."})
		return
	}
}

// Get User details
// @Summary 	Registers a new user with Email and Password
// @Description
// @Tags		Auth
// @Accept 		json
// @Produce 	json
// @Param 		body  	body 		expectedLegacySignupPostData true "Legacy user Signup"
// @Router 		/v1/auth/signup [POST]
// @Success 	200 		{object} 	response.ApiResponse
func (r *authRoutes) signupLegacy(c *gin.Context) {
	basePath, errBasePath := c.Get("appBaseUrl")
	if !errBasePath {
		fmt.Println(errBasePath)
	}

	var body expectedLegacySignupPostData
	// Check is POST data is correct
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		response.RenderEntityNotUnprocessable(c, response.BadRequestError(c.Request.Context()))
		return
	}
	user, err := App.Controller.signupLegacy(body)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// We have a user.
	// make tokens
	tokenDetails, err := CreateToken(user.ID)
	if err != nil {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "err in creating token",
		})
		return
	}
	basePathStr := basePath.(string)
	serializer := AuthLegacySerializer{c}
	App.Service.EmailVerificationMailer(user.Email.String, user.UUID.String(), basePathStr)
	App.Service.PostSignupWorkflow(user, true)
	user, _ = user2.App.Repo.GetUser(map[string]interface{}{"id": user.ID})
	nowTime := time.Now()
	user.LastLogin = &nowTime
	c.Set("currentUser", user)
	// Sending email for Account verification
	// Send otp to user via email
	// CC: Commented 18 Aug 2022
	go func() {
		// Experimental:
		// @todo commented this because this giving some error
		//	App.Service.PostSignupWorkflow(user)

		// Sending the emails.
		//mailData := mailers.SendEmail{
		//	ToEmails: []string{user.Email.String},
		//}
		//mailers.App.Service.SendEmailEvent(mailers.VerifyEmail, mailData)
	}()

	response.RenderResponse(c, serializer.GetAuthSerializerWithTokens(*tokenDetails))
	return
}

// Get User details
// @Summary 	Legacy Login User
// @Description
// @Tags		Auth
// @Accept 		json
// @Produce 	json
// @Param 		body  	body 		expectedLoginPostData true "Login User Data"
// @Router 		/v1/auth/login [POST]
// @Success 	200 		{object} 	response.ApiResponse
func (r *authRoutes) loginLegacy(c *gin.Context) {
	var body expectedLoginPostData
	var userInstance *models.User
	var usernameAsEmail string
	var gerr error
	// Check is POST data is correct
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		response.RenderEntityNotUnprocessable(c, response.BadRequestError(c.Request.Context()))
		return
	}
	if body.IsGhostLogin {
		// This Ghost Login
		user, err := App.Controller.loginLegacyWithToken(body)
		userInstance = user
		gerr = err
	} else {
		// Fork in the flow now. If OTP has value then we call Different controller. As we have only Email and OTP.
		if body.Otp == "" {
			logger.Info("App.Controller.loginLegacy Path ")
			logger.Debug(body.Password)
			usernameAsEmail = strings.ToLower(body.Account)
			user, err := App.Controller.loginLegacy(usernameAsEmail, body.Password)
			userInstance = user
			gerr = err
		} else {
			logger.Info(" App.Controller.loginLegacyWithOtp Found. ")
			usernameAsEmail = strings.ToLower(body.Account)
			user, err := App.Controller.loginLegacyWithOtp(usernameAsEmail, body.Otp)
			userInstance = user
			gerr = err
		}
	}

	if gerr != nil {
		c.JSON(400, gin.H{"error": gerr.Error()})
		return
	}
	tokenDetails, errToken := CreateToken(userInstance.ID)
	if errToken != nil {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "Error in creating token",
		})
		return
	}
	c.Set("currentUser", userInstance)
	serializer := AuthLegacySerializer{c}
	//App.Repo.updateLastLogin(userInstance.ID)
	queries.App.UserQueries.UpdateUser(userInstance.ID, map[string]interface{}{"last_login": time.Now()}, false)

	nowTime := time.Now()
	userInstance.LastLogin = &nowTime
	response.RenderResponse(c, serializer.GetAuthSerializerWithTokens(*tokenDetails))
	// Delete the Token after This
	deletedOtpDataFromCache(strconv.FormatUint(userInstance.ID, 10))
	go func() {
		var bipkafaka kafka.BipKafka
		bipkafaka.Publish(
			configs.KAKFA_TOPICS_CALCULATE_PERMISSIONS,
			utils.String(userInstance.ID),
			[]byte(utils.String(userInstance.ID)))

		//payload, _ := json.Marshal(userInstance)
		//apiClient.AddToQueue(tasks.LoginEmail, payload, apiClient.DEFAULT, apiClient.CommonRetry)
	}()
	return
}

// Logout User
// @Summary 	This will logout a given user token.
// @Description Future we'll integrate with Redis
// @Tags		Auth
// @Accept 		json
// @Produce 	json
// @Router 		/v1/auth/logout [POST]
// @Success 	200 		{object} 	response.ApiResponse
func (r *authRoutes) logout(c *gin.Context) {
	// Frontend Removes the token
	//user, _ := c.Get("currentUser")
	accessTokenRaw := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(accessTokenRaw, "Bearer ")
	App.Controller.LogoutProper(token)
	response.RenderSuccessResponse(c, "User logged out successfully.")
	return
}

// Get User details
// @Summary 	Check if user exists with this email or not
// @Description
// @Tags		Auth
// @Accept 		json
// @Produce 	json
// @Param 		body  	body 		expectedExistingEmailBody true "Exisitng Emails Checker Data"
// @Router 		/v1/auth/existing-email [POST]
// @Success 	200 		{object} 	response.ApiResponse
func (r *authRoutes) existingEmail(c *gin.Context) {
	var body expectedExistingEmailBody
	// Check is POST data is correct
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		response.RenderEntityNotUnprocessable(c, response.BadRequestError(c.Request.Context()))
		return
	}
	flag := App.Controller.checkExistingEmail(body.Email)
	if flag {
		response.RenderSuccessResponse(c, "existing")
		return
	}

	response.RenderNotFoundResponse(c, "new")
	return
}

// @Summary 	Refresh Tokens
// @Description
// @Tags		Auth
// @Param 		body  	body 		expectedRefreshTokenPostData true "Refresh Token Post"
// @Router 		/v1/auth/refresh-token [POST]
// @Success 	200 		{object} 	TokenDetails
func (r *authRoutes) refreshToken(c *gin.Context) {
	// @todo: Refactor
	var body expectedRefreshTokenPostData
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// New Code
	refreshToken := body.RefreshToken
	//verify the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(configs.GetConfigString("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Refresh token expired")
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	fmt.Println(claims)
	if ok && token.Valid {
		//refreshUuid, ok := claims["refreshTokenId"].(string) //convert the interface to string
		//fmt.Println("refreshTokenIdrefreshTokenId", refreshUuid)
		//if !ok {
		//	c.JSON(http.StatusUnprocessableEntity, err)
		//	return
		//}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["userId"]), 10, 64)
		fmt.Println(userId)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "Error occurred")
			return
		}
		//Delete the previous Refresh Token
		rc := redis.RedisClient()
		rctx := redis.GetBgContext()
		rc.Del(rctx, "authtokens:"+body.AccessTokenID)
		//Create new pairs of refresh and access tokens
		newTokens, _ := CreateToken(userId)
		serializer := AuthLegacySerializer{c}
		response.RenderOkWithData(c, serializer.GetRefreshToken(*newTokens))
		return
	} else {
		c.JSON(http.StatusUnauthorized, "refresh expired")
	}

}

// Get User details
// @Summary 	Check if user exists with this username or not
// @Description
// @Tags		Auth
// @Accept 		json
// @Produce 	json
// @Param 		body  	body 		expectedExistingUsernameBody true "Exisitng Username Checker Data"
// @Router 		/v1/auth/existing-username [POST]
// @Success 	200 		{object} 	response.ApiResponse
func (r *authRoutes) existingUsername(c *gin.Context) {
	var body expectedExistingUsernameBody
	// Check is POST data is correct
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		response.RenderEntityNotUnprocessable(c, response.BadRequestError(c.Request.Context()))
		return
	}
	flag := App.Controller.checkExistingUsername(body.Username)
	if flag {
		response.RenderSuccessResponse(c, "existing")
		return
	}

	response.RenderNotFoundResponse(c, "new")
	return
}

// Get User details
// @Summary 	Generate OTP for a USER
// @Tags		Auth
// @Accept 		json
// @Produce 	json
// @Param 		body  	body 		PostNewOtpRequest true "New OTP Request"
// @Router 		/v1/auth/otp [POST]
// @Success 	200 		{object} 	ResponseNewOtp
func (r *authRoutes) generateUserOtp(c *gin.Context) {
	var body PostNewOtpRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if email exists : Send Error
	status := user2.App.Service.CheckUserExistsWithEmail(body.Email)
	if status == false {
		response.RenderNotFoundResponse(c, "No user found  with this email")
		return
	}

	// Check is Token Exists in Redis  : Send Error
	// RedisUserOtpNS
	var userInstance *models.User

	userInstance, _ = App.Service.getUserWithEmail(body.Email)
	if userInstance == nil {
		response.RenderNotFoundResponse(c, "No user found  with this email")
		return
	}

	userIdStr := strconv.FormatUint(userInstance.ID, 10)
	rc := redis.RedisClient()
	rctx := redis.GetBgContext()
	val, err := rc.Get(rctx, models.RedisUserOtpNS+userIdStr).Result()
	fmt.Println(val)

	// if err is nil meaning Key was found
	if err == nil {
		rc.Del(rctx, models.RedisUserOtpNS+userIdStr)
		//response.RenderNotFoundResponse(c, "OTP Token already sent to user.")
		//return
	}

	// Generate OTP and Save in Redis : Send Email
	newOtp, _ := utils.GenerateOTP(6)
	newOptObject := ResponseNewOtp{
		Email: body.Email,
		Otp:   newOtp,
	}
	OtpObjectjson, _ := json.Marshal(newOptObject)
	// set OTP for 24 Hours
	rc.Set(rctx, models.RedisUserOtpNS+userIdStr, OtpObjectjson, 24*time.Hour)

	// Respond with Email and OTP for now.
	App.Service.OTPMailer(newOptObject.Email, newOptObject.Otp)
	serializer := AuthLegacySerializer{c}
	response.RenderOkWithData(c, serializer.GetOtpDetails(newOptObject))
	return
}

// Social Login
// We will do login based on social
// Here the workflow is w do things,
// When a person try to login we will check if we are getting "email"
// Email will be primary KEY
// @Summary 	Social Login V1
// @Tags		Auth
// @Accept 		json
// @Produce 	json
// @Param 		body  	body 		SocialAuthPost true "Auth Login"
// @Router 		/v1/auth/social-login [POST]
func (r *authRoutes) SocialLogin(c *gin.Context) {
	// Please read this carefully as this function does a lots of things.
	// Get the POST from the FE
	var body SocialAuthPost
	var createdUser *models.User
	var errAppUserNotFoundWithEmailFlow error
	var errUserFoundFlow error
	byPassCheckUserWithEmailFlow := false
	isSignUp := false
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	validate := body.Validate()
	if validate != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), validate))
		return
	}
	// Note: We are changing the order of social login.
	// We are first checking if the user has a Social AuthInstance
	// If that is found we will not do any check with the email
	createdUser = App.Service.UserNotFoundWithEmailMayHaveSocialAuthInstanceFlow(body)
	if createdUser != nil {
		byPassCheckUserWithEmailFlow = true
	}
	// if email is empty in the POST
	if body.Email == "" && createdUser == nil {
		// We did not find any Email from this and also UserNotFoundWithEmailMayHaveSocialAuthInstanceFlow
		// Check is the user exists
		// Means we need to create a USER and User Auth Object
		createdUser, errAppUserNotFoundWithEmailFlow = App.Service.UserNotFoundWithEmailFlow(body)
		if errAppUserNotFoundWithEmailFlow != nil {
			response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), errAppUserNotFoundWithEmailFlow))
			return
		}
		//App.Repo.updateUserEmailConfirmationFlag(createdUser.ID)
		queries.App.UserQueries.UpdateUser(createdUser.ID, map[string]interface{}{"is_email_verified": true}, false)
	}

	// We are skipping the whole email workflow if the user was found in the earlier social flow.
	if !byPassCheckUserWithEmailFlow {
		isSignUp = true
		// Do we have a user with this email?
		user, _ := App.Service.getUserWithEmail(body.Email)
		// User found with given Email
		if user != nil {
			// User found with this email flow
			createdUser, errUserFoundFlow = App.Service.UserFoundWithEmailFlow(user, body)
			// Setting Email as Verified
			//App.Repo.updateUserEmailConfirmationFlag(createdUser.ID)
			queries.App.UserQueries.UpdateUser(user.ID, map[string]interface{}{"is_email_verified": true}, false)
			if errUserFoundFlow != nil {
				response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), errUserFoundFlow))
				return
			}
		} else {
			// We are making Emails optional we need to check is a user Exists for
			// ProviderID and Provider if Instance is found with this meaning we have a user
			createdUser = App.Service.UserNotFoundWithEmailMayHaveSocialAuthInstanceFlow(body)
			if createdUser == nil {
				// User Not found with this email flow
				createdUser, errAppUserNotFoundWithEmailFlow = App.Service.UserNotFoundWithEmailFlow(body)
				if errAppUserNotFoundWithEmailFlow != nil {
					response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), errAppUserNotFoundWithEmailFlow))
					return
				}
			}
			userProfile := user2.UpdateUserValidator{
				FullName:    createdUser.FullName,
				Username:    createdUser.Username,
				IsSetupDone: true,
				AvatarUrl:   body.Image,
			}
			_, err := user2.App.Controller.SetupUserController(&userProfile, createdUser)
			// setting email as verified as social login
			//App.Repo.updateUserEmailConfirmationFlag(createdUser.ID)
			queries.App.UserQueries.UpdateUser(createdUser.ID, map[string]interface{}{"is_email_verified": true}, false)
			// start the post signup flow as this is a new user.
			App.Service.PostSignupWorkflow(createdUser, false)
			createdUser, _ = queries.App.UserQueries.GetUser(map[string]interface{}{"id": createdUser.ID})
			if err != nil {
				response.RenderCustomErrorResponse(c, err)
				return
			}
		}
	}
	isNewUser := App.Service.CheckIsNewUser(createdUser)
	tokenDetails, err := CreateToken(createdUser.ID)
	if err != nil {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "err in creating token",
		})
		return
	}
	c.Set("currentUser", createdUser)
	nowTime := time.Now()
	createdUser.LastLogin = &nowTime
	createdUser.IsSetupDone = true
	//App.Repo.updateLastLogin(createdUser.ID)
	queries.App.UserQueries.UpdateUser(createdUser.ID, map[string]interface{}{
		"last_login":        time.Now(),
		"via_discord_store": body.UserViaDiscordAppDirectory,
	}, false)
	//App.Repo.UpdateUserViaDiscordAppDirectory(createdUser.ID, body.UserViaDiscordAppDirectory)
	serializer := AuthLegacySerializer{c}
	// Experimental
	resp := serializer.GetAuthSerializerWithTokensAndSocials(*tokenDetails)
	resp.IsNewUser = isNewUser
	resp.IsSignUp = isSignUp
	response.RenderResponse(c, resp)
	go func() {
		var bipkafaka kafka.BipKafka
		bipkafaka.Publish(
			configs.KAKFA_TOPICS_CALCULATE_PERMISSIONS,
			utils.String(createdUser.ID),
			[]byte(utils.String(createdUser.ID)))
	}()
	return

}

// @todo: this is BROOKEN

// @Summary 	user has forgotton the password, we need an email
// @Tags		Auth
// @Accept 		json
// @Produce 	json
// @Param 		body  	body 		forgotPasswordData true "Reset Password"
// @Router 		/v1/auth/forgot-password [POST]
func (r *authRoutes) forgotPasswordInit(c *gin.Context) {
	var body forgotPasswordData
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		response.RenderEntityNotUnprocessable(c, response.BadRequestError(c.Request.Context()))
		return
	}
	// Lowercase email
	email := strings.ToLower(body.Email)
	user, _ := queries.App.UserQueries.GetUser(map[string]interface{}{"email": email})
	if user == nil {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "No user was found with this email",
		})
		return
	}

	newOtp, _ := utils.GenerateOTP(6)
	userIdStr := strconv.FormatUint(user.ID, 10)
	rc := redis.RedisClient()
	rctx := redis.GetBgContext()
	val, err := rc.Get(rctx, models.RedisUserForgotPassword+userIdStr).Result()
	fmt.Println(val)
	// if err is nil meaning Key was found
	if err == nil {
		// Todo: We send the email again
		c.JSON(http.StatusOK, gin.H{
			"message": "We have re-sent the email.",
		})
		return
	}

	newOptObject := ResponseNewOtp{
		Email: body.Email,
		Otp:   newOtp,
	}
	OtpJson, _ := json.Marshal(newOptObject)
	// set OTP for 24 Hours
	rc.Set(rctx, models.RedisUserForgotPassword+userIdStr, OtpJson, 24*time.Hour)

	App.Service.ForgotPasswordMailer(newOptObject.Email, newOptObject.Otp)

	c.JSON(http.StatusOK, gin.H{
		"message": "We have reset password email.",
	})
	return
}

// @Summary 	Reset Password
// @Tags		Auth
// @Accept 		json
// @Produce 	json
// @Param 		body  	body 		ResetPasswordData true "Reset Password"
// @Router 		/v1/auth/reset-password [POST]
func (r *authRoutes) resetPassword(c *gin.Context) {
	var body ResetPasswordData
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		response.RenderEntityNotUnprocessable(c, response.BadRequestError(c.Request.Context()))
		return
	}

	rc := redis.RedisClient()
	ctx := redis.GetBgContext()
	var FinalResponse ResponseNewOtp
	var FoundKey string
	iter := rc.Scan(ctx, 0, "user-forgot-pass:*", 0).Iterator()
	for iter.Next(ctx) {
		var thisResponseObject ResponseNewOtp
		value := iter.Val()
		val, _ := rc.Get(ctx, value).Result()
		_ = json.Unmarshal([]byte(val), &thisResponseObject)
		if thisResponseObject.Otp == body.Token {
			FinalResponse.Otp = thisResponseObject.Otp
			FinalResponse.Email = thisResponseObject.Email
			FoundKey = value
			break
		}
		FoundKey = "noKeyFound"
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}

	if FinalResponse.Otp == "" {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "We could not find the token details",
		})
		return
	}

	if FoundKey == "noKeyFound" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Incorrect Token",
		})
	}

	err4 := App.Controller.resetPasswordWithToken(FinalResponse.Email, body.Password)
	if err4 != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"errpr": err4.Error(),
		})
	} else {
		rc.Del(ctx, FoundKey)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "New Password is set, re-login.",
	})
	return
}

// @Summary 	Change password for Loggedin User
// Description Send Old and New password, Check user object, some users will have hasPassword false, send empty Old password
// @Tags		Auth
// @Accept 		json
// @Produce 	json
// @Param 		body  	body 		changePasswordData true "Change Password"
// @Router 		/v1/auth/change-password [POST]
func (r *authRoutes) changePassword(c *gin.Context) {
	var body changePasswordData
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		response.RenderEntityNotUnprocessable(c, response.BadRequestError(c.Request.Context()))
		return
	}
	user, _ := r.GetLoggedInUser(c)
	errReset := App.Controller.resetPassword(*user, body.OldPassword, body.NewPassword)
	if errReset != nil {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "error in updating password",
		})
		return
	}

	response.RenderSuccessResponse(c, "Password is updated.")
	return
}

func (r *authRoutes) SpecialGhostLogin(c *gin.Context) {
	email := c.Query("email")
	email = strings.ToLower(email)
	if email == "" {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "No Email was found ",
		})
		return
	}
	// Lowercase email
	user, _ := queries.App.UserQueries.GetUser(map[string]interface{}{"email": email})
	if user == nil {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "No user was found with this email",
		})
		return
	}

	id2 := uuid.New()
	rc := redis.RedisClient()
	ctx := redis.GetBgContext()
	_, err := rc.Get(ctx, models.GhostLogin+email).Result()
	if err == nil {
		// Todo: We send the email again
		rc.Del(ctx, models.GhostLogin+email)
	}

	newOptObject := NewGhostToken{
		Email: email,
		Token: id2.String(),
	}
	TokenJson, _ := json.Marshal(newOptObject)
	// set OTP for 24 Hours
	rc.Set(ctx, models.GhostLogin+email, TokenJson, 24*time.Hour)
	c.JSON(http.StatusOK, gin.H{
		"token": id2.String(),
	})
	return
}

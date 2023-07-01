package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"strconv"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/search"

	"gorm.io/datatypes"
)

type UserService interface {
	// Add your methods here
	EmptyUser(email, password, username, avatarURL string) *models.User
	NewUserSocialAuth(userId uint64, providerName, providerId string, metadata datatypes.JSON) *models.UserSocialAuth
	CreateNewUser(email, password, username, avatarURL string) *models.User
	CheckUserExistsWithEmail(email string) bool
	CreateNewUserProfile(userID uint64, bio string, twitter_url *string, website *string, location *string) (*models.UserProfile, error)
	AddAllUsersToAlgolia() error
	AddUserToAlgolia(userID uint64) error
}

func (us userService) EmptyUser(email, password, username, avatarURL string) *models.User {
	return &models.User{
		Email: sql.NullString{
			String: email,
			Valid:  true,
		},
		Password: password,
		// last login
		Username: username,
		//AvatarUrl: avatarURL,
		// is super user
		// is empty
		//timezone
		// datejoined

	}
}

func (us userService) NewUserSocialAuth(userId uint64, providerName, providerId string, metadata datatypes.JSON) *models.UserSocialAuth {
	return &models.UserSocialAuth{
		UserID:       userId,
		ProviderName: providerName,
		ProviderID:   providerId,
		Metadata:     metadata,
	}
}

// Add error handling
func (us userService) CreateNewUser(email, password, username, avatarURL string, clientReferenceId string) *models.User {
	user := &models.User{
		Email: sql.NullString{
			String: email,
			Valid:  true,
		},
		Password:          password,
		Username:          username,
		ClientReferenceId: clientReferenceId,
	}
	createdUser := queries.App.UserQueries.CreateUser(user)
	return createdUser
}

// Create user profile
func (us userService) CreateNewUserProfile(userID uint64, bio string, twitter_url *string, website *string, location *string) (*models.UserProfile, error) {

	var userProfile *models.UserProfile
	userProfile, err := App.Repo.GetUserProfile(userID)
	if err == nil {
		return userProfile, nil
	}
	userProfile = &models.UserProfile{
		Bio:        bio,
		TwitterUrl: twitter_url,
		Website:    website,
		Location:   location,
		UserID:     userID,
	}

	err = App.Repo.CreateUserProfile(userProfile)
	if err != nil {
		return nil, err
	}
	return userProfile, nil
}

func (us userService) CheckUserExistsWithEmail(email string) bool {
	return App.Repo.CheckIfUserPresentWithEmail(email)
}

func (us userService) AddAllUsersToAlgolia() error {
	users, err := App.Repo.GetAllUsers()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	userDocs := []interface{}{}
	for _, usr := range *users {
		userDocs = append(userDocs, *UserModelToUserDocument(&usr))
	}
	err = search.GetIndex(search.UserDocumentIndexName).SaveRecords(userDocs)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func (ss userService) AddUserToAlgolia(userID uint64) error {
	//usr, err := App.Repo.GetUser(map[string]interface{}{"id": userID})
	usr, err := queries.App.UserQueries.GetUserByIDFromDB(userID)

	if err != nil {
		logger.Error(err.Error())
		return err
	}
	usrDoc := UserModelToUserDocument(usr)
	err = search.GetIndex(search.UserDocumentIndexName).SaveRecord(usrDoc)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

type FollowUserService interface {
	FollowerCountUse(user models.User) FollowUserFollowCountResponse
}

type followService struct{}

// NewBasicGoodsService returns a naive, stateless implementation of GoodsService.
func NewFollowService() FollowUserService {
	return &followService{}
}

func (fs followService) FollowerCountUse(user models.User) FollowUserFollowCountResponse {
	loggedInUserIdString := strconv.FormatUint(user.ID, 10)
	rc := redis.RedisClient()
	rctx := redis.GetBgContext()
	fmt.Println(models.RedisFollowUserNS + loggedInUserIdString)
	val, err := rc.Get(rctx, models.RedisFollowUserNS+loggedInUserIdString).Result()
	if err == nil {
		data := FollowUserFollowCountResponse{}
		_ = json.Unmarshal([]byte(val), &data)
		return data
	}
	followersCount, _ := App.FollowRepo.UserCountFollowing(user.ID)
	followingCount, _ := App.FollowRepo.UserCountFollower(user.ID)

	fc := FollowUserFollowCountResponse{
		Followers: followersCount,
		Following: followingCount,
	}
	// update redis
	updateCahceFollowerCount(fc, user.ID)
	return fc
}

func (us userService) CreateNewSocialUser(email, username, fullName, avatarURL string, clientReferenceId string) (*models.User, error) {
	newSocialUser := models.User{}
	newSocialUser.Username = username
	if email != "" {
		newSocialUser.Email = sql.NullString{
			String: email,
			Valid:  true,
		}
	} else {
		newSocialUser.Email = sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	newSocialUser.FullName = fullName
	newSocialUser.AvatarUrl = avatarURL
	newSocialUser.HasPassword = false
	newSocialUser.ClientReferenceId = clientReferenceId
	// Though we are setting a RANDOM Password, Primary Login will only Allow Social. They can set a password later.
	unHashedPassword := utils.NewNanoid()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(unHashedPassword), bcrypt.DefaultCost)
	newSocialUser.Password = string(hashedPassword)
	newSocialUserCreated := queries.App.UserQueries.CreateUser(&newSocialUser)
	return newSocialUserCreated, nil
}

func (us userService) CreateUserSocialAuth(userID uint64, providerName string, providerID string, metadata []byte) error {
	usa := models.UserSocialAuth{}
	usa.UserID = userID
	usa.ProviderName = providerName
	usa.ProviderID = providerID
	usa.Metadata = metadata
	err := App.Repo.CreateUserSocialAuth(&usa)
	if err != nil {
		return err
	}
	return nil
}

// DoesUSAExitsWithOutUserID : Return UserID or 0 from UserSocialAuth
func (us userService) DoesUSAExitsWithOutUserID(providerName string, providerID string) uint64 {
	instance, _ := App.Repo.GetUserSocialAuth(map[string]interface{}{"provider_name": providerName, "provider_id": providerID})
	if instance == nil {
		return 0
	}
	return instance.UserID
}

func (us userService) DoesUSAExits(userID uint64, providerName string, providerID string) bool {
	instance, _ := App.Repo.GetUserSocialAuth(map[string]interface{}{"user_id": userID, "provider_name": providerName, "provider_id": providerID})
	if instance == nil {
		return false
	}

	return true
}

func (us userService) UpdateUSA(query map[string]interface{}, updates map[string]interface{}) error {
	err := App.Repo.db.Model(models.UserSocialAuth{}).Where(query).Updates(updates).Error
	return err
}

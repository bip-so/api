package queries

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/supabase"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"strconv"
	"strings"
	"time"
)

// GetUser: Preloads with User Profile
func (q *userQuery) GetUser(query map[string]interface{}) (*models.User, error) {
	var user models.User
	err := postgres.GetDB().Model(&models.User{}).Where(query).Preload("UserProfile").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//CreateUser
func (q *userQuery) CreateUser(user *models.User) *models.User {
	_ = postgres.GetDB().Create(user).Error
	return user
}

func (q *userQuery) UpdateUser(userID uint64, updates map[string]interface{}, supabaseUpdate bool) (*models.User, error) {
	var user *models.User
	err := postgres.GetDB().Model(models.User{}).Where(
		models.User{BaseModel: models.BaseModel{ID: userID}}).Updates(updates).Find(&user).Error

	if supabaseUpdate {
		go func() {
			supabase.UpdateUserDefaultStudioSupabase(userID, updates["default_studio_id"].(uint64))
		}()
	}
	return user, err
}

func (q *userQuery) GetUsersByIDs(userIDs []uint64) ([]models.User, error) {
	var users []models.User
	err := postgres.GetDB().Model(&models.User{}).Where("id in ?", userIDs).Preload("UserProfile").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (q *userQuery) GetUserSocialAuthByIDs(userIDs []uint64) ([]models.UserSocialAuth, error) {
	var usersSocialAuths []models.UserSocialAuth
	err := postgres.GetDB().Model(&models.UserSocialAuth{}).Where("user_id IN ?", userIDs).Find(&usersSocialAuths).Error
	if err != nil {
		return nil, err
	}
	return usersSocialAuths, nil
}

const UserPreLoader = "users-preload:"

// GetUserByID: This one request the USER Objects with Profile included
// The function here will set the cache for some time we can change this
// The idea is simple we save things in cache

func (q *userQuery) GetUserByID(id uint64) (*models.User, error) {
	var user models.User
	userIDStr := strconv.FormatUint(id, 10)
	key := UserPreLoader + userIDStr
	val, _ := redis.RedisClient().Get(context.Background(), key).Result()
	if val == "" {
		fmt.Println("Sending DB Value: " + key)
		err := postgres.GetDB().Model(&models.User{}).Where("id = ?", id).Preload("UserProfile").First(&user).Error
		if err != nil {
			return nil, err
		}
		userObjectJson, _ := json.Marshal(user)
		_ = redis.RedisClient().Set(context.Background(), key, userObjectJson, 120*time.Second).Err()
	} else {
		fmt.Println("Sending Cached Value: " + key)
		val2, _ := redis.RedisClient().Get(context.Background(), key).Result()
		_ = json.Unmarshal([]byte(val2), &user)
	}
	fmt.Println(&user)
	defer utils.TimeTrack(time.Now())
	return &user, nil
}

func (q *userQuery) GetUserByIDFromDB(id uint64) (*models.User, error) {
	var user models.User
	err := postgres.GetDB().Model(&models.User{}).Where("id = ?", id).Preload("UserProfile").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (q *userQuery) IsUserAdminInStudio(userId uint64, studioId uint64) (bool, error) {
	studioAdminRole, _ := App.RoleQuery.GetStudioAdminRole(studioId)
	// For this user and this studio get member id
	memberInstance, _ := App.MemberQuery.GetMember(map[string]interface{}{"user_id": userId, "studio_id": studioId})
	var isAdmin bool
	for _, mem := range studioAdminRole.Members {
		if memberInstance.ID == mem.ID {
			isAdmin = true
			break
		}
	}
	return isAdmin, nil
}

func (q *userQuery) GetUserSocialAuth(query map[string]interface{}) (*models.UserSocialAuth, error) {
	var user models.UserSocialAuth
	err := postgres.GetDB().Model(&models.UserSocialAuth{}).Where(query).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (q *userQuery) HasBIPAccount(email string) bool {
	var user models.User
	result := postgres.GetDB().Model(&user).Where("email = ?", email)
	if result.RowsAffected == 0 {
		return false
	}
	return true
}

func (q userQuery) GetUserInstanceByEmail(email string) *models.User {
	var user *models.User
	_ = postgres.GetDB().Model(&models.User{}).Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error
	return user
}

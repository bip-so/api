package follow

import (
	"encoding/json"
	"fmt"
	"strconv"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/user"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
)

func (c followController) GetUserFollowFollowCountHandler(userID uint64) (FollowUserFollowCountResponse, error) {
	// get a json of followers and following
	// get user instance
	// redis empty call api
	// return
	// check redis
	loggedInUserIdString := strconv.FormatUint(userID, 10)
	rc := redis.RedisClient()
	rctx := redis.GetBgContext()
	val, err := rc.Get(rctx, models.RedisFollowUserNS+loggedInUserIdString).Result()
	if err == nil {
		data := FollowUserFollowCountResponse{}
		_ = json.Unmarshal([]byte(val), &data)
		return data, nil
	}

	ffcrObject := App.Service.FollowerCountUse(userID)
	return ffcrObject, nil
}

func (c followController) FollowUserHandle(user *models.User, toFollowUser uint64) {

}

func (c followController) GetStudioFollowersCountHandler(studioID uint64) (FollowUserStudioCountResponse, error) {
	studioIDStr := strconv.FormatUint(studioID, 10)
	rc := redis.RedisClient()
	rctx := redis.GetBgContext()
	val, err := rc.Get(rctx, models.RedisFollowUserStudioNS+studioIDStr).Result()
	if err == nil {
		data := FollowUserStudioCountResponse{}
		_ = json.Unmarshal([]byte(val), &data)
		return data, nil
	}
	ffcrObject := App.Service.StudioFollowCount(studioID)
	return ffcrObject, nil
}

func (c followController) GetUserFollowers(userID uint64) ([]user.UserMiniSerializer, error) {
	followers, err := App.Repo.GetFollowUsers(map[string]interface{}{"user_id": userID})
	if err != nil {
		return nil, err
	}
	fmt.Println(len(followers))
	users := []models.User{}
	for _, follower := range followers {
		users = append(users, *follower.FollowerUser)
	}
	serializerUser := []user.UserMiniSerializer{}
	for _, usr := range users {
		serialized := user.UserMiniSerializerData(&usr)
		resp, _ := c.GetUserFollowFollowCountHandler(usr.ID)
		serialized.Followers = resp.Followers
		serialized.Following = resp.Following
		serializerUser = append(serializerUser, serialized)
	}
	return serializerUser, nil
}

func (c followController) GetUserFollowing(userID uint64) ([]user.UserMiniSerializer, error) {
	followers, err := App.Repo.GetFollowUsers(map[string]interface{}{"follower_id": userID})
	if err != nil {
		return nil, err
	}
	users := []models.User{}
	for _, follower := range followers {
		users = append(users, *follower.User)
	}
	serializerUser := []user.UserMiniSerializer{}
	for _, usr := range users {
		serialized := user.UserMiniSerializerData(&usr)
		resp, _ := c.GetUserFollowFollowCountHandler(usr.ID)
		serialized.Followers = resp.Followers
		serialized.Following = resp.Following
		serializerUser = append(serializerUser, serialized)
	}
	return serializerUser, nil
}

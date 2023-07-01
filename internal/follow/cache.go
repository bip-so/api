package follow

import (
	"encoding/json"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"strconv"
)

// Redis Update -? User Followers
func updateCahceFollowerCount(fc FollowUserFollowCountResponse, userid uint64) {
	loggedInUserIdString := strconv.FormatUint(userid, 10)
	rc := redis.RedisClient()
	rctx := redis.GetBgContext()
	fcjson, _ := json.Marshal(fc)
	rc.Set(rctx, models.RedisFollowUserNS+loggedInUserIdString, fcjson, 0)
}

// Redis Delete Cache -? User Followers
func deleteCahceFollowerCount(userid uint64) {
	loggedInUserIdString := strconv.FormatUint(userid, 10)
	rc := redis.RedisClient()
	rc.Del(redis.GetBgContext(), models.RedisFollowUserNS+loggedInUserIdString)
}

// Redis update for Studio -. Follow
func updateCahceStudioFollowerCount(fc FollowUserStudioCountResponse, studioid uint64) {
	studioIDStr := strconv.FormatUint(studioid, 10)
	rc := redis.RedisClient()
	rctx := redis.GetBgContext()
	fcjson, _ := json.Marshal(fc)
	rc.Set(rctx, models.RedisFollowUserStudioNS+studioIDStr, fcjson, 0)
}

func deleteCahceStudioFollowerCount(studioid uint64) {
	studioIDStr := strconv.FormatUint(studioid, 10)
	rc := redis.RedisClient()
	rc.Del(redis.GetBgContext(), models.RedisFollowUserStudioNS+studioIDStr)
}

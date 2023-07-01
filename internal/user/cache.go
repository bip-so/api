package user

import (
	"encoding/json"
	"strconv"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
)

// Redis Update -? User Followers
func updateCahceFollowerCount(fc FollowUserFollowCountResponse, userid uint64) {
	loggedInUserIdString := strconv.FormatUint(userid, 10)
	rc := redis.RedisClient()
	rctx := redis.GetBgContext()
	fcjson, _ := json.Marshal(fc)
	rc.Set(rctx, models.RedisFollowUserNS+loggedInUserIdString, fcjson, 0)
}

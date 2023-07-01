package follow

import (
	"github.com/gin-gonic/gin"
)

type FollowUserSerializer struct {
	c *gin.Context
}

type FollowStudioSerializer struct {
	c *gin.Context
}

func (self *FollowUserSerializer) GetUserFollowCounts(countRes FollowUserFollowCountResponse) FollowUserFollowCountResponse {
	return FollowUserFollowCountResponse{
		Followers: countRes.Followers,
		Following: countRes.Following,
	}
}

func (self *FollowStudioSerializer) GetStudioFollowCounts(countRes FollowUserStudioCountResponse) FollowUserStudioCountResponse {
	return FollowUserStudioCountResponse{
		Followers: countRes.Followers,
	}
}

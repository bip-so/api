package canvasrepo

import (
	"context"
	"fmt"
	"strconv"
)

const UserRepoSetAll = "repo-access-history:"

func (crc *repoUserCachingService) AddToUserRepoToSet(userID, repoID uint64) {
	fmt.Println("REPO SET")
	userIDStr := strconv.Itoa(int(userID))
	repoIDStr := strconv.Itoa(int(repoID))
	fmt.Println(userIDStr)
	fmt.Println(repoIDStr)
	result := crc.redisClient.SAdd(context.Background(), UserRepoSetAll+userIDStr, repoIDStr)
	fmt.Println(result)
}

package xpcontribs

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"strconv"
	"time"
)

func (s *xpcontribService) XPSummarizer(studioID uint64) {
	studio, _ := queries.App.StudioQueries.GetStudioQuery(map[string]interface{}{"id": studioID})
	if !studio.FeatureFlagHasXP {
		return
	}
	//
	studioIntegration, err := queries.App.StudioIntegrationQuery.GetDiscordStudioIntegration(studioID)
	if err != nil {
		fmt.Println("Error in getting studio Integration", err)
		return
	}
	guildID := studioIntegration.TeamID

	// We need to do following now.
	// With this studioID
	studioIDStr := strconv.Itoa(int(studioID))
	studioKey := MainStudioLogNameSpace + studioIDStr
	// RecentUserList is string slice
	RecentUserList := redis.RedisClient().SMembers(context.Background(), MainStudioAuditLogNameSpace+studioIDStr).Val()
	fmt.Println(RecentUserList)

	UserAddedKeyMap := redis.RedisClient().HGetAll(context.Background(), studioKey).Val()
	fmt.Println("UserAddedKeyExists", UserAddedKeyMap)

	userPoints := map[string]int{}
	userIDs := []uint64{}
	// Loop through the RecentUserList (user_id) as string
	for _, userIDstr := range RecentUserList {
		fmt.Println(userIDstr)
		userID, _ := strconv.ParseUint(userIDstr, 10, 64)
		fmt.Println(userID)
		userIDs = append(userIDs, userID)
		// userID need to get discord id
	}
	userSocialAuths, _ := queries.App.UserQueries.GetUserSocialAuthByIDs(userIDs)
	for _, userSocialAuth := range userSocialAuths {
		userIDStr := utils.String(userSocialAuth.UserID)
		addedPoints, err := strconv.Atoi(UserAddedKeyMap[userIDStr+"-added"])
		if err != nil {
			addedPoints = 0
		}
		updatedPoints, _ := strconv.Atoi(UserAddedKeyMap[userIDStr+"-updated"])
		if err != nil {
			updatedPoints = 0
		}
		userPoints[userSocialAuth.ProviderID] = addedPoints + updatedPoints
	}
	studioUserPoints := map[string]map[string]int{
		guildID: userPoints,
	}
	studioUserPointsStr, _ := json.Marshal(studioUserPoints)
	s.cache.HSet(context.Background(), MainStudioPointsNameSpace+"start_it", guildID, studioUserPointsStr)
	for _, userID := range userIDs {
		userIDStr := utils.String(userID)
		s.cache.HDelete(context.Background(), studioKey, userIDStr+"-added")
		s.cache.HDelete(context.Background(), studioKey, userIDStr+"-updated")
		s.cache.HDelete(context.Background(), studioKey, userIDStr+"-deleted")
	}
	// Loop through the UserAddedKeyExists Get / Summarizer
	// Get Discord ID for This Studio
	// Get DiscordID for a User
	defer utils.TimeTrack(time.Now())
}

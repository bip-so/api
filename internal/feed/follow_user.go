package feed

import (
	"context"
	"fmt"

	"github.com/GetStream/stream-go2/v7"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	bipStream "gitlab.com/phonepost/bip-be-platform/pkg/stores/stream"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

func (s *feedService) SelfFollowUser(userID string) {
	followee, err := s.Stream.FlatFeed(bipStream.FlatFeedName, userID)
	if err != nil {
		logger.Error(err.Error())
	}
	followerTimeline, err := s.Stream.FlatFeed(bipStream.FlatTimelineName, userID)
	if err != nil {
		logger.Error(err.Error())
	}
	_, err = followerTimeline.Follow(context.Background(), followee, stream.WithFollowFeedActivityCopyLimit(100))
	if err != nil {
		logger.Error(err.Error())
	}
}

func (s *feedService) BulkSelfFollowUser(users []uint64) {
	calls := [][]stream.FollowRelationship{}
	follows := []stream.FollowRelationship{}
	for i, userID := range users {
		followee, err := s.Stream.FlatFeed(bipStream.FlatFeedName, utils.String(userID))
		if err != nil {
			logger.Error(err.Error())
		}
		followerTimeline, err := s.Stream.FlatFeed(bipStream.FlatTimelineName, utils.String(userID))
		follows = append(follows, stream.NewFollowRelationship(followee, followerTimeline))
		if i == len(users)-1 || (len(follows)%1000 == 0) { //2500 is limit in per call of followMany
			calls = append(calls, follows)
			follows = []stream.FollowRelationship{}
		}
	}
	s.Stream.FollowMany(context.Background(), follows)
	for _, follow := range calls {
		err := s.Stream.FollowMany(context.Background(), follow, stream.WithFollowManyActivityCopyLimit(100))
		if err != nil {
			fmt.Println("error on bulk follow many", err)
		}
	}
}

func (s *feedService) FollowUser(followUser *models.FollowUser) {
	followee, err := s.Stream.FlatFeed(bipStream.FlatFeedName, utils.String(followUser.UserId))
	if err != nil {
		logger.Error(err.Error())
	}
	followerTimeline, err := s.Stream.FlatFeed(bipStream.FlatTimelineName, utils.String(followUser.FollowerId))
	if err != nil {
		logger.Error(err.Error())
	}
	_, err = followerTimeline.Follow(context.Background(), followee, stream.WithFollowFeedActivityCopyLimit(100))
	if err != nil {
		logger.Error(err.Error())
	}
}

func (s *feedService) UnfollowUser(userId uint64, followeeUserId uint64) {
	followee, err := s.Stream.FlatFeed(bipStream.FlatFeedName, utils.String(followeeUserId))
	if err != nil {
		logger.Error(err.Error())
	}
	followerTimeline, err := s.Stream.FlatFeed(bipStream.FlatTimelineName, utils.String(userId))
	if err != nil {
		logger.Error(err.Error())
	}
	response, err := followerTimeline.Unfollow(context.Background(), followee)
	if err != nil {
		logger.Error(err.Error())
	}
	fmt.Println("Unfollow user response:", response)
}

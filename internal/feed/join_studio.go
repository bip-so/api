package feed

import (
	"context"
	"fmt"

	"github.com/GetStream/stream-go2/v7"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	bipStream "gitlab.com/phonepost/bip-be-platform/pkg/stores/stream"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

func (s *feedService) JoinStudio(studioID uint64, userID uint64) {
	followee, err := s.Stream.FlatFeed(bipStream.FlatFeedName, utils.String(studioID))
	if err != nil {
		logger.Error(err.Error())
	}
	followerTimeline, err := s.Stream.FlatFeed(bipStream.FlatTimelineName, utils.String(userID))
	if err != nil {
		logger.Error(err.Error())
	}
	_, err = followerTimeline.Follow(context.Background(), followee, stream.WithFollowFeedActivityCopyLimit(100))
	if err != nil {
		fmt.Println(err)
		logger.Error(err.Error())
	}
}

func (s *feedService) LeaveStudio(studioID uint64, userID uint64) {
	followee, err := s.Stream.FlatFeed(bipStream.FlatFeedName, utils.String(studioID))
	if err != nil {
		logger.Error(err.Error())
	}
	followerTimeline, err := s.Stream.FlatFeed(bipStream.FlatTimelineName, utils.String(userID))
	if err != nil {
		logger.Error(err.Error())
	}
	response, err := followerTimeline.Unfollow(context.Background(), followee)
	if err != nil {
		logger.Error(err.Error())
	}
	fmt.Println("Unfollow studio response:", response)
}

func (s *feedService) BulkJoinStudio(studioID uint64, userIDs []uint64) {
	followee, err := s.Stream.FlatFeed(bipStream.FlatFeedName, utils.String(studioID))
	if err != nil {
		logger.Error(err.Error())
	}
	calls := [][]stream.FollowRelationship{}
	follows := []stream.FollowRelationship{}
	for i, userID := range userIDs {
		followerTimeline, err := s.Stream.FlatFeed(bipStream.FlatTimelineName, utils.String(userID))
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		follows = append(follows, stream.NewFollowRelationship(followerTimeline, followee))
		if i == len(userIDs)-1 || (len(follows)%1000 == 0) { //2500 is limit in per call of followMany
			calls = append(calls, follows)
			follows = []stream.FollowRelationship{}
		}
	}
	for _, follow := range calls {
		err = s.Stream.FollowMany(context.Background(), follow, stream.WithFollowManyActivityCopyLimit(100))
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

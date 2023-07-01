package workflows

import (
	"context"
	"errors"
	"fmt"
	"github.com/GetStream/stream-go2/v7"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	bipStream "gitlab.com/phonepost/bip-be-platform/pkg/stores/stream"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/supabase"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

func WorkflowJoinUserToStudio(user *models.User, studioId uint64) error {
	var err error
	mem, err := queries.App.MemberQuery.GetMember(map[string]interface{}{"user_id": user.ID, "studio_id": studioId})
	if err == nil && mem != nil {
		if mem.IsRemoved {
			return errors.New("can't join as you were banned")
		} else if mem.HasLeft {
			return queries.App.MemberQuery.JoinStudio([]uint64{user.ID}, studioId)
		}
	} else {
		newlyCreatedMember := queries.App.MemberQuery.AddUserIDToStudio(user.ID, studioId)
		if newlyCreatedMember == nil {
			return errors.New("error in adding user to studio")
		}
		err = queries.App.MemberQuery.AddMembersToStudioInMemberRole(studioId, []models.Member{*newlyCreatedMember})
		if err != nil {
			return err
		}
	}

	// post user joins a studio
	go func() {
		queries.App.StudioQueries.DeleteUserAssociatedStudioDataByUserID(user.ID)
		supabase.UpdateUserSupabase(user.ID, true)
	}()
	return nil
}

const (
	NotificationFeed = "notification"
	FlatFeedName     = "bip_feed"
	FlatTimelineName = "timeline"
)

// WorkflowHelper
func WorkflowHelperFeedUpdateOnJoinStudio(studioID uint64, userID uint64) {
	var streamClient *stream.Client
	followee, err := streamClient.FlatFeed(FlatFeedName, utils.String(studioID))
	if err != nil {
		logger.Error(err.Error())
	}
	followerTimeline, err := streamClient.FlatFeed(FlatTimelineName, utils.String(userID))
	if err != nil {
		logger.Error(err.Error())
	}
	_, err = followerTimeline.Follow(context.Background(), followee, stream.WithFollowFeedActivityCopyLimit(100))
	if err != nil {
		fmt.Println(err)
		logger.Error(err.Error())
	}
}

func WorkflowHelperFeedUpdateOnLeaveStudio(studioID uint64, userID uint64) {
	var streamClient *stream.Client
	followee, err := streamClient.FlatFeed(bipStream.FlatFeedName, utils.String(studioID))
	if err != nil {
		logger.Error(err.Error())
	}
	followerTimeline, err := streamClient.FlatFeed(FlatTimelineName, utils.String(userID))
	if err != nil {
		logger.Error(err.Error())
	}
	response, err := followerTimeline.Unfollow(context.Background(), followee)
	if err != nil {
		logger.Error(err.Error())
	}
	fmt.Println("Unfollow studio response:", response)
}

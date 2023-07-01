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

func (s feedService) AddReelActivity(reel *models.Reel) {
	App.Repo.db.Where("id = ?", reel.ID).Preload("CanvasBranch").First(&reel)
	fmt.Println("reel id", reel.ID)
	streamActivity := stream.Activity{
		Actor:     utils.String(reel.CreatedByID),
		Verb:      "Reel",
		Object:    utils.String(reel.ID),
		ForeignID: utils.String(reel.ID),
		Time: stream.Time{
			Time: reel.CreatedAt,
		},
	}

	if reel.CanvasBranch.PublicAccess == "private" {
		feeds := []stream.Feed{}
		canvasBranchPerms, err := App.Repo.GetCanvasBranchPermission(map[string]interface{}{"canvas_branch_id": reel.CanvasBranchID})
		if err != nil {
			logger.Error(err.Error())
		}
		for _, branchPerm := range canvasBranchPerms {
			if branchPerm.Role != nil {
				for _, member := range branchPerm.Role.Members {
					feed, _ := s.Stream.FlatFeed(bipStream.FlatTimelineName, utils.String(member.UserID))
					feeds = append(feeds, feed)
				}
			} else {
				feed, _ := s.Stream.FlatFeed(bipStream.FlatTimelineName, utils.String(branchPerm.Member.UserID))
				feeds = append(feeds, feed)
			}
		}
		err = bipStream.AddFeedToMany(streamActivity, feeds)
		if err != nil {
			fmt.Println("Error on adding private reelActivity To feed", err)
		}
		fmt.Println("private reels feed added ")
	} else {
		userFeed, err := bipStream.GetFlatFeed(bipStream.FlatFeedName, utils.String(reel.CreatedByID))
		if err != nil {
			logger.Error(err.Error())
			return
		}
		streamActivity.To = []string{bipStream.FlatFeedName + ":" + utils.String(reel.StudioID)}
		response, err := userFeed.AddActivity(context.Background(), streamActivity)
		if err != nil {
			fmt.Println("Error on adding private reelActivity To feed", err)
			return
		}
		fmt.Println("user add feed response ", response)
	}
}

func (s feedService) GetReelActivities(userID uint64, offset int, limit int) (*stream.FlatFeedResponse, error) {
	feeds, err := bipStream.GetFlatFeed(bipStream.FlatTimelineName, utils.String(userID))
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	var response *stream.FlatFeedResponse
	if offset < 1 {
		response, err = feeds.GetActivities(context.Background(),
			stream.WithActivitiesLimit(limit),
			stream.WithActivitiesOffset(offset))
	} else {
		response, err = feeds.GetActivities(
			context.Background(),
			stream.WithActivitiesLimit(limit),
			stream.WithActivitiesOffset(offset),
		)
	}
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return response, nil
}

func (s feedService) RemoveReelActivity(reel *models.Reel, userID uint64) {
	if reel.CanvasBranch.PublicAccess == "private" {
		canvasBranchPerms, err := App.Repo.GetCanvasBranchPermission(map[string]interface{}{"canvas_branch_id": reel.CanvasBranchID})
		if err != nil {
			logger.Error(err.Error())
		}
		for _, branchPerm := range canvasBranchPerms {
			if branchPerm.Role != nil {
				for _, member := range branchPerm.Role.Members {
					userFeed, _ := s.Stream.FlatFeed(bipStream.FlatTimelineName, utils.String(member.UserID))
					_, err = userFeed.RemoveActivityByForeignID(context.Background(), utils.String(reel.ID))
					if err != nil {
						logger.Error(err.Error())
					}
				}
			} else {
				userFeed, _ := s.Stream.FlatFeed(bipStream.FlatTimelineName, utils.String(branchPerm.Member.UserID))
				_, err = userFeed.RemoveActivityByForeignID(context.Background(), utils.String(reel.ID))
				if err != nil {
					logger.Error(err.Error())
				}
			}
		}
	} else {
		userFeed, err := bipStream.GetFlatFeed(bipStream.FlatFeedName, utils.String(userID))
		if err != nil {
			logger.Error(err.Error())
			return
		}
		_, err = userFeed.RemoveActivityByForeignID(context.Background(), utils.String(reel.ID))
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}
}

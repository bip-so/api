package notifications

import (
	"context"
	"fmt"
	"time"

	"github.com/GetStream/stream-go2/v7"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	stream2 "gitlab.com/phonepost/bip-be-platform/pkg/stores/stream"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

func (s notificationService) AddNotificationToStream(notification *models.Notification) error {
	userID := utils.String(notification.NotifierID)
	notificationFeed, err := stream2.StreamClient().NotificationFeed(stream2.NotificationFeed, userID)
	if err != nil {
		fmt.Println("stream client error", err)
		return err
	}
	notificationData := SerializeNotification(*notification)
	activity := stream.Activity{
		Actor:     "User:" + userID,
		Object:    notification.Entity + "-" + notification.Event, // We need to update this
		Verb:      "User:" + userID,
		Time:      stream.Time{Time: notification.CreatedAt.Round(60 * time.Second)},
		ForeignID: utils.String(notification.ID),
		Extra: map[string]interface{}{
			"data": notificationData,
		},
	}
	var activitiesTo []string
	if notification.IsPersonal {
		activitiesTo = append(activitiesTo, stream2.NotificationFeed+":"+userID+"-personal")
	}
	if notification.StudioID != nil {
		activitiesTo = append(activitiesTo, stream2.NotificationFeed+":"+userID+"-"+utils.String(*notification.StudioID))
	}
	activity.To = activitiesTo
	_, err = notificationFeed.AddActivity(context.TODO(), activity)
	if err != nil {
		fmt.Println("Error while adding notification ", err)
		return err
	}
	fmt.Println("Successfully added notification to the stream", notification.ID)
	return nil
}

func (s notificationService) GetStreamNotificationForUser(query string, skip int, limit int) ([]interface{}, error) {
	notificationsFeed, err := stream2.GetNotificationFeed(stream2.NotificationFeed, query)
	if err != nil {
		fmt.Println("stream client error: ", err)
		return nil, err
	}
	var resp *stream.NotificationFeedResponse
	if skip < 1 {
		resp, err = notificationsFeed.GetActivities(context.Background())
	} else {
		resp, err = notificationsFeed.GetActivities(context.Background(),
			stream.WithNotificationsMarkRead(true),
			stream.WithNotificationsMarkSeen(true),
			stream.WithActivitiesLimit(limit),
			stream.WithActivitiesOffset(skip),
		)
	}
	if err != nil {
		fmt.Println("stream client error: ", err)
		return nil, err
	}
	data := []interface{}{}
	for _, result := range resp.Results {
		for _, activity := range result.Activities {
			activityData, isok := activity.Extra["data"].(interface{})
			if isok {
				data = append(data, activityData)
			}
		}
	}
	return data, err
}

func (s notificationService) UpdateActivity(notification *models.Notification) error {
	activities, err := stream2.GetActivitiesByForeignID(stream.NewForeignIDTimePair(utils.String(notification.ID), stream.Time{Time: notification.CreatedAt.Round(60 * time.Second)}))
	if err != nil {
		fmt.Println("stream client error: ", err)
		return err
	}
	for _, activity := range activities.Results {
		activity.Extra["data"] = SerializeNotification(*notification)
		_, err = stream2.UpdateActivities(activity)
		if err != nil {
			fmt.Println("stream client error: ", err)
		}
	}
	return nil
}

func (s notificationService) RemoveActivity(notification *models.Notification) {
	fmt.Println("Removing feed from activity", notification.NotifierID, notification.ID)
	notificationsFeed, err := stream2.GetNotificationFeed(stream2.NotificationFeed, utils.String(notification.NotifierID))
	if err != nil {
		fmt.Println("stream client error: ", err)
		return
	}
	_, err = notificationsFeed.RemoveActivityByForeignID(context.Background(), utils.String(notification.ID))
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func (s notificationService) BulkAddNotificationToStream(notifications []models.Notification) error {
	calls := [][]stream.Activity{}
	activities := []stream.Activity{}
	var notificationFeed *stream.NotificationFeed
	var err error
	for i, notification := range notifications {
		// s.AddNotificationToStream(&notification)
		userID := utils.String(notification.NotifierID)
		notificationFeed, err = stream2.StreamClient().NotificationFeed(stream2.NotificationFeed, userID)
		if err != nil {
			fmt.Println("stream client error", err)
			return err
		}
		notificationData := SerializeNotification(notification)
		activity := stream.Activity{
			Actor:     "User:" + userID,
			Object:    notification.Entity + "-" + notification.Event, // We need to update this
			Verb:      "User:" + userID,
			Time:      stream.Time{Time: notification.CreatedAt.Round(60 * time.Second)},
			ForeignID: utils.String(notification.ID),
			Extra: map[string]interface{}{
				"data": notificationData,
			},
		}
		var activitiesTo []string
		if notification.IsPersonal {
			activitiesTo = append(activitiesTo, stream2.NotificationFeed+":"+userID+"-personal")
		}
		if notification.StudioID != nil {
			activitiesTo = append(activitiesTo, stream2.NotificationFeed+":"+userID+"-"+utils.String(*notification.StudioID))
		}
		activity.To = activitiesTo
		activities = append(activities, activity)
		if i == len(notifications)-1 || (len(activities)%500 == 0) {
			calls = append(calls, activities)
			activities = []stream.Activity{}
		}
	}
	fmt.Println("Came here after for loop")
	for _, activity := range activities {
		fmt.Println(activity.Actor, activity.Verb, activity.Object)
	}
	for _, activities := range calls {
		_, err = notificationFeed.AddActivities(context.TODO(), activities...)
		if err != nil {
			fmt.Println("Error while adding notification Trying again after sleep", err)
			time.Sleep(60)
			_, err = notificationFeed.AddActivities(context.TODO(), activities...)
			if err != nil {
				fmt.Println("Error while adding notification ", err)
				return err
			}
		}
		fmt.Println("Added notifications to stream")
	}
	return nil
}

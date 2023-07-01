package stream

import (
	"context"
	"fmt"
	"github.com/GetStream/stream-go2/v7"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"time"
)

type Stream struct {
	client *stream.Client
}

func NewStream() *Stream {
	return &Stream{
		client: StreamClient(),
	}
}

var streamClient *stream.Client
var ctx = context.Background()

func InitStreamClient() {
	var err error
	streamClient, err = stream.New(
		configs.GetStreamConfig().ApiKey,
		configs.GetStreamConfig().ApiSecret,
		stream.WithAPIRegion("us-east"),
		stream.WithAPIVersion("1.0"),
		stream.WithTimeout(5*time.Second),
	)
	if err != nil {
		logger.Error(err.Error())
	}
	fmt.Println("GetStream Initiated successfully")
}

func GetFlatFeed(slug string, id string) (*stream.FlatFeed, error) {
	return streamClient.FlatFeed(slug, id)
}

func GetNotificationFeed(slug string, id string) (*stream.NotificationFeed, error) {
	return streamClient.NotificationFeed(slug, id)
}

func GetActivitiesByForeignID(slug stream.ForeignIDTimePair) (*stream.GetActivitiesResponse, error) {
	response, err := streamClient.GetActivitiesByForeignID(ctx, slug)
	return response, err
}

func UpdateActivities(slug stream.Activity) (*stream.BaseResponse, error) {
	response, err := streamClient.UpdateActivities(ctx, slug)
	return response, err
}

func GetAggregatedFeed(slug string, id string) (*stream.AggregatedFeed, error) {
	return streamClient.AggregatedFeed(slug, id)
}

func AddFeedToMany(activity stream.Activity, feeds []stream.Feed) error {
	return streamClient.AddToMany(ctx, activity, feeds...)
}

func FollowMany(relationships []stream.FollowRelationship, opts []stream.FollowManyOption) error {
	return streamClient.FollowMany(ctx, relationships, opts...)
}

func StreamClient() *stream.Client {
	return streamClient
}

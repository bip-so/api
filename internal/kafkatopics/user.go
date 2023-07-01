package kafkatopics

import (
	"github.com/segmentio/kafka-go"
	"gitlab.com/phonepost/bip-be-platform/internal/feed"
	"gitlab.com/phonepost/bip-be-platform/internal/user"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

func NewUserConsumer(msg *kafka.Message) {
	user.App.Service.AddUserToAlgolia(utils.Uint64(string(msg.Key)))
}

func UpdateUserConsumer(msg *kafka.Message) {
	user.App.Service.AddUserToAlgolia(utils.Uint64(string(msg.Key)))
}

func AddUserToFeedStream(msg *kafka.Message) {
	feed.App.Service.SelfFollowUser(string(msg.Key))
}

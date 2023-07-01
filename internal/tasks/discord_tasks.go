package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"gitlab.com/phonepost/bip-be-platform/cmd/discord_worker/eventHandler"
	"gitlab.com/phonepost/bip-be-platform/internal/canvasrepo"
	cache "gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"strings"
	"time"
)

func (s taskService) UpdateDiscordTreeMessage(ctx context.Context, task *asynq.Task) {
	var payload map[string]uint64
	json.Unmarshal(task.Payload(), &payload)
	collectionID := payload["collectionId"]
	canvasrepo.App.Service.SendCollectionTreeToDiscord(collectionID)
}

func (s taskService) RunFailedDiscordEventsCron(ctx context.Context, task *asynq.Task) {
	s.cache = cache.NewCache()
	iter := s.cache.GetAllMatchingKeys(ctx, eventHandler.RedisDiscordNamespaceFailed+"*")
	for iter.Next(ctx) {
		value := s.cache.Get(ctx, iter.Val())
		keyStr := iter.Val()
		valStr := value.(string)
		fullKeySplit := strings.Split(keyStr, ":")
		keyName := fullKeySplit[1] // getting 2nd part
		fmt.Println("Processing: ", keyName)

		data := eventHandler.EventSerializer{}
		err := json.Unmarshal([]byte(valStr), &data)
		if err != nil {
			s.cache.Set(ctx, eventHandler.RedisDiscordNamespaceFailed+keyName, value, &cache.Options{Expiration: time.Hour * 168})
		}

		// process this value
		err = eventHandler.Process(data, valStr)
		if err != nil {
			s.cache.Set(ctx, eventHandler.RedisDiscordNamespaceFailed+keyName, value, &cache.Options{Expiration: time.Hour * 24})
		}

		// success
		//s.cache.Set(ctx, eventHandler.RedisDiscordNamespaceProcessed+keyName, value, nil)
		s.cache.Delete(ctx, iter.Val())
	}
}

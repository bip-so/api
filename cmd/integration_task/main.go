package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/cmd/api"
	"gitlab.com/phonepost/bip-be-platform/cmd/integration_task/discord_integration"
	"gitlab.com/phonepost/bip-be-platform/cmd/integration_task/slack_integration"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/lambda/connect_discord/connect"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/core"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	cache "gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"strings"
	"time"
)

func appSetup() {
	//fileName := ".env"
	fileName := ".env"
	// Init config from the env file.
	configs.InitConfig(fileName, ".")
	// Start everything : Logger / DB / Redis : Need Err.
	core.InitCore(fileName, ".")
	api.InitAllApps()
}

const (
	RedisIntegrationTaskNamespaceAll       = "integrationtask:*"
	RedisIntegrationTaskNamespaceProcessed = "integrationtask-success:"
	RedisIntegrationTaskNamespaceFailed    = "integrationtask-failed:"
)

func main() {
	appSetup()
	cacheInstance := cache.NewCache()
	ctx := context.Background()

	for {
		iter := cacheInstance.GetAllMatchingKeys(ctx, RedisIntegrationTaskNamespaceAll)
		for iter.Next(ctx) {
			value := cacheInstance.Get(ctx, iter.Val())
			keyStr := iter.Val()
			valStr := value.(string)
			fullKeySplit := strings.Split(keyStr, ":")
			keyName := fullKeySplit[1] // getting 2nd part
			fmt.Println(valStr)
			fmt.Println("Processing: ", keyName)

			data := map[string]uint64{}
			err := json.Unmarshal([]byte(valStr), &data)
			if err != nil {
				errValue := map[string]interface{}{
					"key":   keyStr,
					"value": value,
					"error": err,
				}
				errValueBytes, _ := json.Marshal(errValue)
				cacheInstance.Set(ctx, RedisIntegrationTaskNamespaceFailed+keyName, errValueBytes, &cache.Options{Expiration: time.Hour * 168})
				cacheInstance.Delete(ctx, iter.Val())
				continue
			}
			integration, err := connect.GetIntegrationByID(data["id"])
			if integration.TeamID == "" {
				fmt.Println("Error: invalid discord integration")
				cacheInstance.Set(ctx, RedisIntegrationTaskNamespaceFailed+keyName, err, &cache.Options{Expiration: time.Hour * 168})
				cacheInstance.Delete(ctx, iter.Val())
				continue
			}
			if integration.Type == models.DISCORD_INTEGRATION_TYPE {
				// Start discord integration task
				// turning on the dm notification on
				postgres.GetDB().Model(models.Studio{}).Where("id = ?", integration.StudioID).Update("discord_notifications_enabled", true)
				err = discord_integration.DiscordIntegrationTask(integration)
				if err != nil {
					errValue := map[string]interface{}{
						"key":   keyStr,
						"value": value,
						"error": err,
					}
					errValueBytes, _ := json.Marshal(errValue)
					cacheInstance.Set(ctx, RedisIntegrationTaskNamespaceFailed+keyName, errValueBytes, &cache.Options{Expiration: time.Hour * 24})
					cacheInstance.Delete(ctx, iter.Val())
					UpdateStudioIntegration(integration.ID, models.StudioIntegrationFailed)
					notifications.App.Service.PublishNewNotification(notifications.DiscordIntegrationTaskFailed, 0, []uint64{}, &integration.StudioID,
						nil, notifications.NotificationExtraData{}, nil, nil)
					continue
				}
			} else if integration.Type == models.SLACK_INTEGRATION_TYPE {
				postgres.GetDB().Model(models.Studio{}).Where("id = ?", integration.StudioID).Update("slack_notifications_enabled", true)
				err = slack_integration.SlackIntegrationTask(integration)
				if err != nil {
					errValue := map[string]interface{}{
						"key":   keyStr,
						"value": value,
						"error": err,
					}
					errValueBytes, _ := json.Marshal(errValue)
					cacheInstance.Set(ctx, RedisIntegrationTaskNamespaceFailed+keyName, errValueBytes, &cache.Options{Expiration: time.Hour * 24})
					cacheInstance.Delete(ctx, iter.Val())
					UpdateStudioIntegration(integration.ID, models.StudioIntegrationFailed)
					// @todo change the notification to slack integration task failed
					notifications.App.Service.PublishNewNotification(notifications.DiscordIntegrationTaskFailed, 0, []uint64{}, &integration.StudioID,
						nil, notifications.NotificationExtraData{}, nil, nil)
					continue
				}
			}
			// success
			cacheInstance.Set(ctx, RedisIntegrationTaskNamespaceProcessed+keyName, value, nil)
			cacheInstance.Delete(ctx, iter.Val())
			UpdateStudioIntegration(integration.ID, models.StudioIntegrationSuccess)
		}
	}
}

func UpdateStudioIntegration(integrationID uint64, status string) {
	err := postgres.GetDB().Model(models.StudioIntegration{}).Where("id = ?", integrationID).Updates(map[string]interface{}{
		"integration_status": status,
	}).Error
	if err != nil {
		fmt.Println("Error in updating studio integration status: ", err)
	}
}

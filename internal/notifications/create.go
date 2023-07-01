package notifications

import (
	"context"
	"encoding/json"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"reflect"
)

func (s notificationService) PublishNewNotification(
	event string, createdByID uint64, notifierIDs []uint64, studioID *uint64, roleIDs *[]uint64, extraData NotificationExtraData, objectID *uint64, contentObject *string) {
	notification := PostNotificationInstance(event, createdByID, notifierIDs, studioID, roleIDs, extraData, objectID, contentObject)
	notificationString, _ := json.Marshal(notification)
	s.kafka.Publish(configs.KAFKA_TOPICS_NOTIFICATIONS, event, notificationString)
}

// CreateNotification entity param is string and, it should be matched with the function names
// It calls the specific method based on category to create the different types of notification.
// Probably we only create the stuff from here.
// Question: If notification settings are off for all the types for user. Should we have to create notification record in our db?
func (s notificationService) CreateNotification(event string, notification *PostNotification) {
	eventHandler := event + "Handler"
	t := notificationService{}
	method := reflect.ValueOf(t).MethodByName(eventHandler)
	params := []reflect.Value{reflect.ValueOf(notification)}
	method.Call(params)
}

func (s notificationService) ExecuteAllRoughBranchNotifications(branchID uint64) {
	notifications := s.cache.HGetAll(context.Background(), s.GetRoughBranchNotificationsRedisKey(branchID))
	for _, value := range notifications {
		var notification *PostNotification
		json.Unmarshal([]byte(value), &notification)
		eventHandler := notification.Event + "Handler"
		t := notificationService{}
		method := reflect.ValueOf(t).MethodByName(eventHandler)
		params := []reflect.Value{reflect.ValueOf(notification)}
		method.Call(params)
	}
	s.cache.Delete(context.Background(), s.GetRoughBranchNotificationsRedisKey(branchID))
}

package notifications

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gorm.io/gorm"
)

func (s notificationService) updateNotificationCountTable(notification *models.Notification) {
	notificationCount, err := App.Repo.GetNotificationCount(map[string]interface{}{"user_id": notification.NotifierID})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.CreateNewNotificationCount(notification)
		return
	}
	notificationCount.All += 1
	if notification.IsPersonal {
		notificationCount.Personal += 1
	}
	if notification.StudioID != nil && notification.IsPersonal == false {
		studioCountList := []StudioNotificationCountView{}
		json.Unmarshal([]byte(notificationCount.Studio), &studioCountList)
		studioPresent := false
		for index, studioCount := range studioCountList {
			if studioCount.Studio.ID == *notification.StudioID {
				studioPresent = true
				studioCountList[index].Count += 1
				break
			}
		}
		if studioPresent == false {
			studio, _ := App.Repo.GetStudioByID(*notification.StudioID)
			fmt.Println("Building the studio and studio count", studio)
			studioCountData := StudioNotificationCountView{
				Count:  1,
				Studio: BuildStudio(studio),
			}
			studioCountList = append(studioCountList, studioCountData)
		}
		parsedStudioCountList, err := json.Marshal(studioCountList)
		if err != nil {
			fmt.Println("[updateNotificationCountTable] Error on parsedStudioCountList parsing", err)
		}
		notificationCount.Studio = string(parsedStudioCountList)
	}
	err = App.Repo.SaveNotificationCount(notificationCount)
	if err != nil {
		fmt.Println("[updateNotificationCountTable] Error on updating the notification count", err)
	}
}

// CreateNewNotificationCount if user is new user or old user without notification count record
func (s notificationService) CreateNewNotificationCount(notification *models.Notification) {
	countList := App.Repo.GetNotificationCountForUser(notification.NotifierID)
	StudioCountData := s.BuildStudioNotificationCount(countList["studio"].([]map[string]interface{}))
	parsedStudioCountData, err := json.Marshal(StudioCountData)
	if err != nil {
		fmt.Println("Error on create new notification Count", parsedStudioCountData)
		return
	}
	notificationCount := models.NewNotificationCount(
		notification.NotifierID, countList["all"].(int64), countList["personal"].(int64), string(parsedStudioCountData))
	App.Repo.CreateNotificationCount(notificationCount)
}

// BuildStudioNotificationCount for the NotificationCount model
func (s notificationService) BuildStudioNotificationCount(countList []map[string]interface{}) (StudioCountData []StudioNotificationCountView) {
	for _, stdCount := range countList {
		studio := stdCount["studio"].(models.Studio)
		StudioCountData = append(StudioCountData, StudioNotificationCountView{
			Count:  stdCount["count"].(int),
			Studio: BuildStudio(&studio),
		})
	}
	return StudioCountData
}

func (s notificationService) mapStudioWithCount(studioList *[]models.Studio, studioCountList []studioCount) (res []map[string]interface{}) {
	result := []map[string]interface{}{}
	for _, studio := range *studioList {
		count := 0
		for _, stdCount := range studioCountList {
			if stdCount.StudioID == studio.ID {
				break
			}
		}
		result = append(result, map[string]interface{}{
			"count":  count,
			"studio": studio,
		})
	}
	return result
}

// MakeNotificationCountAllRead This will reset all the notificationCount to zero
func (s notificationService) MakeNotificationCountAllRead(userID uint64) error {
	notificationCount, err := App.Repo.GetNotificationCount(map[string]interface{}{"user_id": userID})
	if err != nil {
		fmt.Println("Error on MakeNotificationCountAllRead", err)
		return err
	}
	notificationCount.All = 0
	notificationCount.Personal = 0
	studioCountList := []StudioNotificationCountView{}
	json.Unmarshal([]byte(notificationCount.Studio), &studioCountList)
	for index := range studioCountList {
		studioCountList[index].Count = 0
	}
	parsedStudioCountList, err := json.Marshal(studioCountList)
	if err != nil {
		fmt.Println("[MakeNotificationCountAllRead] Error on parsing studioCountList", err)
		return err
	}
	notificationCount.Studio = string(parsedStudioCountList)
	err = App.Repo.SaveNotificationCount(notificationCount)
	if err != nil {
		fmt.Println("[MakeNotificationCountAllRead] Error on updating the notification count", err)
		return err
	}
	return nil
}

/*
	Triggers when studio is name, image is updated.
	Update Notifications of all the users present in that studio.

	To get the list of user Ids we are Grouping the studio_id and user_id and fetching the notifications.
	From list of notifications we loop to get the UserIDs.
*/
func (s notificationService) UpdateNotificationCountAfterStudioSave(studio *models.Studio) {
	members, err := App.Repo.GetMembersByStudioID(studio.ID)
	if err != nil {
		fmt.Println(err)
	}
	UserIDs := []uint64{}
	for _, member := range members {
		UserIDs = append(UserIDs, member.UserID)
	}
	notificationCountInstances, err := App.Repo.GetNotificationCountByUserIDs(UserIDs)
	if err != nil {
		fmt.Println(err)
	}
	for _, notificationCount := range *notificationCountInstances {
		studioCountList := []StudioNotificationCountView{}
		json.Unmarshal([]byte(notificationCount.Studio), &studioCountList)
		studioPresent := false
		for index, studioCount := range studioCountList {
			if studioCount.Studio.ID == studio.ID {
				studioPresent = true
				studioCountList[index].Studio.DisplayName = studio.DisplayName
				studioCountList[index].Studio.Handle = studio.Handle
				studioCountList[index].Studio.ImageURL = studio.ImageURL
				break
			}
		}
		if studioPresent == false {
			continue
		}
		parsedStudioCountList, err := json.Marshal(studioCountList)
		if err != nil {
			fmt.Println("[MakeNotificationCountAllSeen] Error on parsing studioCountList", err)
		}
		notificationCount.Studio = string(parsedStudioCountList)
		err = App.Repo.SaveNotificationCount(&notificationCount)
		if err != nil {
			fmt.Println("[UpdateNotificationCountAfterStudioSave] Error on updating the notification count", err)
		}
	}
}

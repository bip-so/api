package slack2

import (
	"fmt"
	ar "gitlab.com/phonepost/bip-be-platform/internal/accessrequest"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/internal/studio_integration"
)

func (s slackService) BlockActionsHandler(body *SlackAppMentionPayload) {
	studioIntegration, err := App.Repo.GetStudioIntegration(body.Team.Id)
	if err != nil {
		fmt.Println("Error in getting studio integration", err)
		return
	}
	for _, action := range body.Actions {
		if action.Type == "static_select" && action.ActionId == "grantPermission" {
			permission := action.SelectedOption.Value
			notification, err := notifications.App.Repo.GetNotification(map[string]interface{}{"slack_dm_id": body.Message.Ts})
			if err != nil {
				fmt.Println("Error in getting notification by slack messages ts", err)
				s.SendSlackCommonErrorResponse(body.Channel.Id, body.Message.Ts, studioIntegration.AccessKey)
				return
			}
			mergeRequestBody := ar.ManageAccessRequestPost{
				Status:                      models.ACCESS_REQUEST_ACCEPTED,
				CanvasBranchPermissionGroup: permission,
			}
			err = ar.App.Service.ManageAccessRequest(*notification.ObjectId, mergeRequestBody, notification.NotifierID)
			if err != nil {
				fmt.Println("Error in processing access request", err)
				s.SendSlackCommonErrorResponse(body.Channel.Id, body.Message.Ts, studioIntegration.AccessKey)
				return
			}
		} else if action.Type == "static_select" && action.ActionId == "dmNotifications" {
			status := action.SelectedOption.Value
			dmStatus := true
			if status == "disable" {
				dmStatus = false
			}
			studio_integration.App.Repo.UpdateSlackDmNotification(studioIntegration.StudioID, dmStatus)
		}
	}
}

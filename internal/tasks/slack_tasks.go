package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/slack-go/slack"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	bipSlack "gitlab.com/phonepost/bip-be-platform/internal/slack"
	"gitlab.com/phonepost/bip-be-platform/lambda/slack/connect"
	"gitlab.com/phonepost/bip-be-platform/pkg/integrations"
)

func (s taskService) SlackIntegrationTask(ctx context.Context, task *asynq.Task) {
	var body struct {
		Response       *slack.OAuthResponse `json:"response"`
		StudioInstance *models.Studio       `json:"studioInstance"`
	}
	err := json.Unmarshal(task.Payload(), &body)
	if err != nil {
		fmt.Println("Error in parsing body", err)
	}
	response := body.Response
	studioInstance := body.StudioInstance

	slackMemberList, err := integrations.GetSlackTeamMembers(response.AccessToken, response.TeamID)
	if err != nil {
		fmt.Println("Error in getting slack members", slackMemberList)
	}
	fmt.Println("slackMemberList", slackMemberList)

	slackIDs := []string{}
	for _, slackMember := range slackMemberList.Members {
		slackIDs = append(slackIDs, slackMember.Id)
	}
	fmt.Println("slackIDs", slackIDs)

	existingUsers, err := connect.FindUsersBySlackIDs(slackIDs)
	if err != nil {
		fmt.Println("error in getting user by slackIDs", err)
	}
	fmt.Println("Existing user", existingUsers)

	members, err := connect.GetMembers(studioInstance.ID)
	if err != nil {
		fmt.Println("Error in getting studio members", slackMemberList)
	}
	fmt.Println("Studio Members", members)

	usersSocialAuth := []models.UserSocialAuth{}
	users := []models.User{}
	memberStudioUserIds := []uint64{}
	for _, member := range slackMemberList.Members {
		if member.IsBot || member.Deleted || member.Id == "USLACKBOT" {
			continue
		}
		var foundUser *models.UserSocialAuth
		var foundMember *models.Member
		for _, eUser := range existingUsers {
			if eUser.ProviderID == member.Id {
				foundUser = &eUser
				break
			}
		}

		for _, studioMember := range members {
			if foundUser != nil && foundUser.UserID == studioMember.UserID {
				foundMember = &studioMember
				break
			}
		}

		if foundUser == nil {
			slackUserName := member.Profile.DisplayName + connect.GenerateRandomKey(4)
			nickName := member.Profile.DisplayName
			user := connect.CreateNewUser("", "", slackUserName, nickName, member.Profile.Avatar)
			slackUser := connect.NewSlackUser(user.ID, member.Id, nil)
			users = append(users, *user)
			usersSocialAuth = append(usersSocialAuth, *slackUser)
			foundUser = slackUser
		} else {
			queries.App.StudioQueries.DeleteUserAssociatedStudioDataByUserID(foundUser.UserID)
		}
		if foundMember == nil {
			memberStudioUserIds = append(memberStudioUserIds, foundUser.UserID)
		}
	}
	if len(users) > 0 {
		err = connect.CreateUsers(&users, &usersSocialAuth)
		if err != nil {
			fmt.Println("Error in creating social auths")
			return
		}
	}
	connect.BulkJoinStudio(studioInstance.ID, memberStudioUserIds)

	slackUserGroupResponse, err := integrations.GetSlackUserGroups(response.AccessToken, response.TeamID)
	if err != nil {
		fmt.Println("Error in getting slack user groups")
		return
	}
	allRoleNames := []string{}
	for _, userGroup := range slackUserGroupResponse.UserGroups {
		allRoleNames = append(allRoleNames, userGroup.Name)
	}
	existingMemberGroups, _ := connect.GetRolesByNames(allRoleNames, studioInstance.ID)
	for _, userGroup := range slackUserGroupResponse.UserGroups {
		var foundMemberGroup *models.Role
		for _, mGroup := range existingMemberGroups {
			if mGroup.SlackRoleID.String == userGroup.Id {
				foundMemberGroup = &mGroup
				break
			}
		}
		if foundMemberGroup == nil {
			group, err := connect.CreateNewRole(studioInstance.ID, userGroup.Name, userGroup.Id)
			if err != nil {
				fmt.Println(err, group.Name, group.ID)
				continue
			}
			fmt.Println("new role created======>", group)
			foundMemberGroup = group
		}
		userIDs := []uint64{}
		userSocialAuths, _ := connect.FindUsersBySlackIDs(userGroup.Users)
		for _, usa := range userSocialAuths {
			userIDs = append(userIDs, usa.UserID)
		}
		connect.UpdateMembershipRole(foundMemberGroup.ID, userIDs, []uint64{})
	}
}

func (s taskService) SlackBipMarkAction(ctx context.Context, task *asynq.Task) {
	var body *bipSlack.SlackAppMentionPayload
	err := json.Unmarshal(task.Payload(), &body)
	if err != nil {
		fmt.Println("Error in parsing body", err)
	}
	fmt.Println("Slack Event payload", body)
	bipSlack.App.Service.HandleSlackAppMentions(body)
}

func (s taskService) SlackEventSubscriptions(ctx context.Context, task *asynq.Task) {
	var body *bipSlack.SlackEventTypePayload
	err := json.Unmarshal(task.Payload(), &body)
	if err != nil {
		fmt.Println("Error in parsing body", err)
	}
	fmt.Println("Slack Event payload", body)
	if body.Event.Type == "team_join" {
		var teamJoinBody *bipSlack.SlackUserTeamJoinEventPayload
		err := json.Unmarshal(task.Payload(), &teamJoinBody)
		if err != nil {
			fmt.Println("Error in parsing body", err)
		}
		bipSlack.App.Service.HandleTeamJoinEvent(teamJoinBody)
		return
	} else {
		var eventBody *bipSlack.SlackEventPayload
		err := json.Unmarshal(task.Payload(), &eventBody)
		if err != nil {
			fmt.Println("Error in parsing body", err)
		}
		bipSlack.App.Service.HandleSlackEvents(eventBody)
	}
}

func (s taskService) SlackSlashCommands(ctx context.Context, task *asynq.Task) {
	var slashCommand slack.SlashCommand
	json.Unmarshal(task.Payload(), &slashCommand)
	userSocialAuthInstance, err := queries.App.UserQueries.GetUserSocialAuth(map[string]interface{}{"provider_id": slashCommand.UserID})
	if err != nil {
		fmt.Println("UserSocialAuthInstance Not found", err)
		return
	}
	userInstance, err := queries.App.UserQueries.GetUserByID(userSocialAuthInstance.UserID)
	if err != nil {
		fmt.Println("userInstance Not found", err)
		return
	}
	studioIntegration, err := queries.App.StudioIntegrationQuery.GetStudioIntegration(map[string]interface{}{"team_id": slashCommand.TeamID, "type": models.SLACK_INTEGRATION_TYPE})
	if err != nil {
		fmt.Println("studioIntegration Not found", err)
		return
	}
	studioInstance, err := queries.App.StudioQueries.GetStudioQuery(map[string]interface{}{"id": studioIntegration.StudioID})
	if err != nil {
		fmt.Println("studioInstance Not found", err)
		return
	}
	api := slack.New(studioIntegration.AccessKey)
	var text string
	switch slashCommand.Command {
	case "/bip-search":
		text, err = bipSlack.SlackBipSearchTreeBuilderHandler(slashCommand.Text, studioInstance, userInstance)
		if err != nil {
			fmt.Println("Error in generating text for bip-search  ==> ", err)
			api.PostEphemeral(slashCommand.ChannelID, slashCommand.UserID, slack.MsgOptionText(err.Error(), false))
			return
		}
	case "/bip-new":
		text, err = bipSlack.SlackBipNewSlashCommandHandler(slashCommand, studioInstance, userInstance)
		if err != nil {
			fmt.Println("Error in generating text for bip-new  ==> ", err)
			api.PostEphemeral(slashCommand.ChannelID, slashCommand.UserID, slack.MsgOptionText("User doesn't have permission to create canvases", false))
			return
		}
	default:
		fmt.Println("Slash command received ===>", slashCommand.Command)
		api.PostEphemeral(slashCommand.ChannelID, slashCommand.UserID, slack.MsgOptionText(err.Error(), false))
		return
	}

	msgOptions := slack.MsgOptionText(text, false)
	resp, err := api.PostEphemeral(slashCommand.ChannelID, slashCommand.UserID, msgOptions)
	if err != nil {
		fmt.Println("Error in sending post ephemera message", err)
		return
	}
	fmt.Println("Success response", resp)
	return
}

package slack2

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/role"
	"gitlab.com/phonepost/bip-be-platform/lambda/slack/connect"
	"gitlab.com/phonepost/bip-be-platform/pkg/integrations"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/s3"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"time"
)

func (s slackService) HandleSlackAppMentions(body *SlackAppMentionPayload) {
	if body.Type == "message_action" && body.CallbackId == "bipThis" {
		studioIntegration, err := App.Repo.GetStudioIntegration(body.Team.Id)
		if err != nil {
			fmt.Println("Error in getting studio integration", err)
			return
		}
		fmt.Println("Found studio integration", studioIntegration)
		if len(body.Message.Blocks) > 0 || len(body.Message.Files) > 0 {
			user, err := App.Repo.GetUserSocialAuth(body.User.Id)
			if err == gorm.ErrRecordNotFound {
				user, err = App.Service.CreateSlackUser(studioIntegration, body.User.Id)
				if err != nil {
					s.SendSlackCommonErrorResponse(body.Channel.Id, body.MessageTs, studioIntegration.AccessKey)
					return
				}
			}
			author, err := App.Repo.GetUserSocialAuth(body.Message.User)
			if err == gorm.ErrRecordNotFound {
				author, err = App.Service.CreateSlackUser(studioIntegration, body.Message.User)
				if err != nil {
					s.SendSlackCommonErrorResponse(body.CallbackId, body.MessageTs, studioIntegration.AccessKey)
					return
				}
			}
			fmt.Println("Bip User got from slack userID", user, author)
			content := ""
			for _, block := range body.Message.Blocks {
				for _, element := range block.Elements {
					for _, textElement := range element.Elements {
						if textElement.Type == "text" {
							content += textElement.Text
						} else if textElement.Type == "link" {
							content += textElement.Url
						}
					}
				}
			}
			attachments := []string{}
			for _, file := range body.Message.Files {
				if file.UrlPrivate == "" {
					continue
				}
				fileUrl := s.UploadSlackAttachmentTos3(file.UrlPrivate, studioIntegration.AccessKey, body.Team.Id, file.Id, file.Name)
				if fileUrl != "" {
					attachments = append(attachments, fileUrl)
				}
			}
			slackMessage := models.NewSlackMessage(body.Message.Ts, content, author.UserID, user.UserID, time.Now(), attachments, body.Channel.Id, body.Channel.Name, "", "")
			err = App.Repo.CreateMessage(slackMessage)
			if err != nil {
				fmt.Println("Error in creating slack message")
				s.SendSlackCommonErrorResponse(body.Channel.Id, body.MessageTs, studioIntegration.AccessKey)
				return
			}
			api := slack.New(studioIntegration.AccessKey)
			err = api.AddReaction("white_check_mark", slack.ItemRef{Channel: body.Channel.Id, Timestamp: body.Message.Ts})
			fmt.Println(err)
			messageReply := SlackMessagePayload{
				Text:      "Successfully captured the message. You can place it in an appropriate canvas by typing '//' on the canvas",
				ChannelID: body.Channel.Id,
				ThreadTs:  body.Message.Ts,
			}
			App.Service.SendSlackMessageReply(messageReply, studioIntegration.AccessKey)
		}
	} else if body.Type == "block_actions" && len(body.Actions) > 0 {
		s.BlockActionsHandler(body)
		return
	}
}

func (s slackService) UploadSlackAttachmentTos3(url, accessToken, teamID, fileID, fileName string) string {
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println("[UploadSlackAttachmentTos3] Error on new request method", err)
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("[UploadSlackAttachmentTos3] Error on new request method", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("[UploadSlackAttachmentTos3] Error on parsing the response", err)
	}
	reader := bytes.NewReader(body)
	s3Path := fmt.Sprintf("%s-%s-%s", teamID, fileID, fileName)
	response, err := s3.UploadObjectToBucket(fmt.Sprintf("%s/%s", "slack", s3Path), reader, true)
	if err != nil {
		fmt.Println("error in uploading object", err)
	}
	return response
}

func (s slackService) CreateSlackUser(studioIntegration *models.StudioIntegration, slackUserID string) (*models.UserSocialAuth, error) {
	slackProfile, err := integrations.GetSlackProfile(studioIntegration.AccessKey, slackUserID)
	if err != nil {
		fmt.Println("Error in getting slack profile", err)
		return nil, err
	}
	slackUserName := slackProfile.UserProfile.DisplayName + connect.GenerateRandomKey(4)
	nickName := slackProfile.UserProfile.DisplayName
	userInstance := connect.CreateNewUser("", "", slackUserName, nickName, slackProfile.UserProfile.Avatar)
	slackUser := connect.NewSlackUser(userInstance.ID, slackUserID, nil)
	users := []models.User{*userInstance}
	usersSocialAuth := []models.UserSocialAuth{*slackUser}
	err = connect.CreateUsers(&users, &usersSocialAuth)
	if err != nil {
		fmt.Println("Error in creating social auths")
		return nil, err
	}
	connect.BulkJoinStudio(studioIntegration.StudioID, []uint64{userInstance.ID})
	return slackUser, nil
}

func (s slackService) HandleSlackEvents(body *SlackEventPayload) {
	teamID := body.TeamId
	if teamID == "" && body.Event.Team != "" {
		teamID = body.Event.Team
	} else if teamID == "" {
		teamID = body.Event.Subteam.TeamId
	}
	studioIntegration, err := App.Repo.GetStudioIntegration(teamID)
	if err != nil {
		fmt.Println("Error in getting studio integration", err)
		return
	}
	if body.Event.Type == "app_mention" {
		response := s.GetSlackThreadMessages(body.Event.Channel, body.Event.ThreadTs, studioIntegration.AccessKey)
		var messagesData *SlackMessagesPayload
		json.Unmarshal(response, &messagesData)
		messages := messagesData.Messages
		messages = messages[:len(messages)-1]
		for _, message := range messages {
			user, err := App.Repo.GetUserSocialAuth(body.Event.User)
			if err == gorm.ErrRecordNotFound {
				user, err = App.Service.CreateSlackUser(studioIntegration, body.Event.User)
				if err != nil {
					s.SendSlackCommonErrorResponse(body.Event.Channel, message.Ts, studioIntegration.AccessKey)
					return
				}
			}
			author, err := App.Repo.GetUserSocialAuth(message.User)
			if err == gorm.ErrRecordNotFound {
				author, err = App.Service.CreateSlackUser(studioIntegration, message.User)
				if err != nil {
					s.SendSlackCommonErrorResponse(body.Event.Channel, message.Ts, studioIntegration.AccessKey)
					return
				}
			}
			fmt.Println("Bip User got from slack userID", user, author)
			content := ""
			for _, block := range message.Blocks {
				for _, element := range block.Elements {
					for _, textElement := range element.Elements {
						content += textElement.Text
					}
				}
			}
			attachments := []string{}
			for _, file := range message.Files {
				fileUrl := s.UploadSlackAttachmentTos3(file.UrlPrivate, studioIntegration.AccessKey, body.Event.Team, file.Id, file.Name)
				if fileUrl != "" {
					attachments = append(attachments, fileUrl)
				}
			}
			slackMessage := models.NewSlackMessage(message.Ts, content, author.UserID, user.UserID, time.Now(), attachments, body.Event.Channel, "", "", "")
			err = App.Repo.CreateMessage(slackMessage)
			if err != nil {
				fmt.Println("Error in creating slack message")
				s.SendSlackCommonErrorResponse(body.Event.Channel, message.Ts, studioIntegration.AccessKey)
				return
			}
			api := slack.New(studioIntegration.AccessKey)
			err = api.AddReaction("white_check_mark", slack.ItemRef{Channel: body.Event.Channel, Timestamp: message.Ts})
			fmt.Println(err)
			messageReply := SlackMessagePayload{
				Text:      "Successfully captured the message. You can place it in an appropriate canvas by typing '//' on the canvas",
				ChannelID: body.Event.Channel,
				ThreadTs:  message.Ts,
			}
			App.Service.SendSlackMessageReply(messageReply, studioIntegration.AccessKey)
		}
	} else if body.Event.Type == "subteam_created" {
		role, err := s.CreateSlackRole(studioIntegration, body.Event.Subteam.Name, body.Event.Subteam.Id, body.Event.Subteam.Users)
		if err != nil {
			fmt.Println("Error in creating new role", err)
		}
		fmt.Println("Created new role successfully", role)
	} else if body.Event.Type == "subteam_updated" {
		s.HandleSubTeamUpdated(studioIntegration, body)
	} else if body.Event.Type == "subteam_members_changed" {
		s.SlackRoleMembersChanged(studioIntegration, body)
	}
}

func (s slackService) CreateSlackRole(studioIntegration *models.StudioIntegration, roleName string, slackRoleID string, slackUserIDs []string) (*models.Role, error) {
	repo := role.NewRoleRepo()
	roleObject := &models.Role{
		StudioID:    studioIntegration.StudioID,
		Name:        roleName,
		Color:       "#ffffff",
		IsSystem:    false,
		Icon:        "",
		SlackRoleID: sql.NullString{Valid: true, String: slackRoleID},
	}

	createdRole, err := role.RoleRepo.CreateRole(repo, roleObject)
	if err != nil {
		return nil, err
	}
	fmt.Println("[CreateSlackRole] Role Created", createdRole.SlackRoleID.String, createdRole.ID)

	userIDs := []uint64{}
	userSocialAuths, _ := connect.FindUsersBySlackIDs(slackUserIDs)
	for _, usa := range userSocialAuths {
		userIDs = append(userIDs, usa.UserID)
	}
	connect.UpdateMembershipRole(roleObject.ID, userIDs, []uint64{})
	return createdRole, nil
}

func (s slackService) UpdateSlackRole(roleName string, slackRoleID string, studioID uint64) {
	App.Repo.db.Model(models.Role{}).Where("slack_role_id = ? and studio_id  = ?", slackRoleID, studioID).Updates(map[string]interface{}{
		"name": roleName,
	})
}

func (s slackService) SlackRoleMembersChanged(studioIntegration *models.StudioIntegration, body *SlackEventPayload) {
	roleUpdates := role.UpdateManagementPost{}
	var roleInstance models.Role
	err := App.Repo.db.Model(models.Role{}).Where("slack_role_id = ? and studio_id = ?", body.Event.SubteamId, studioIntegration.StudioID).First(&roleInstance).Error
	if err != nil {
		fmt.Println("error in getting role by slack Role ID", roleInstance, body.Event.SubteamId, err)
		return
	}
	roleUpdates.RoleId = roleInstance.ID

	var addUsers []models.UserSocialAuth
	var addUserIDs []uint64
	err = App.Repo.db.Model(models.UserSocialAuth{}).Where("provider_id in ?", body.Event.AddedUsers).Find(&addUsers).Error
	if err != nil {
		fmt.Println("Error in getting add users", err)
		return
	}
	for _, addUser := range addUsers {
		addUserIDs = append(addUserIDs, addUser.UserID)
	}
	roleUpdates.MembersAdded = addUserIDs

	var removeUsers []models.UserSocialAuth
	var RemoveUserIDs []uint64
	err = App.Repo.db.Model(models.UserSocialAuth{}).Where("provider_id in ?", body.Event.RemovedUsers).Find(&removeUsers).Error
	if err != nil {
		fmt.Println("Error in getting add users", err)
		return
	}
	for _, removeUser := range removeUsers {
		RemoveUserIDs = append(RemoveUserIDs, removeUser.UserID)
	}
	roleUpdates.MembersRemoved = RemoveUserIDs

	_, err = role.App.Service.UpdateMembershipRole(roleUpdates)
	if err != nil {
		fmt.Println("error in updating membership", err)
	}
}

func (s slackService) HandleSubTeamUpdated(studioIntegration *models.StudioIntegration, body *SlackEventPayload) {
	var roleInstance *models.Role
	err := App.Repo.db.Model(models.Role{}).Where("slack_role_id = ? and studio_id = ?", body.Event.Subteam.Id, studioIntegration.StudioID).First(&roleInstance).Error
	if err == gorm.ErrRecordNotFound && body.Event.Subteam.DeletedBy == "" {
		s.CreateSlackRole(studioIntegration, body.Event.Subteam.Name, body.Event.Subteam.Id, body.Event.Subteam.Users)
		return
	}
	if body.Event.Subteam.DeletedBy != "" {
		role.App.Service.DeleteRole(roleInstance.ID)
	} else {
		s.UpdateSlackRole(body.Event.Subteam.Name, body.Event.Subteam.Id, studioIntegration.StudioID)
	}
}

func (s slackService) HandleTeamJoinEvent(body *SlackUserTeamJoinEventPayload) {
	fmt.Println("temid", body.TeamId)
	studioIntegration, err := App.Repo.GetStudioIntegration(body.TeamId)
	if err != nil {
		fmt.Println("Error in getting studio integration", err)
		return
	}
	socialAuth, err := s.CreateSlackUser(studioIntegration, body.Event.User.Id)
	if err != nil {
		fmt.Println("Error in creating new user", err)
	}
	fmt.Println("Created new user socialauth done", socialAuth)
}

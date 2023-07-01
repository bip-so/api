package slack2

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"io/ioutil"
	"net/http"
	"strings"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	user2 "gitlab.com/phonepost/bip-be-platform/internal/user"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

const (
	SLACK_PROVIDER         = "slack"
	SLACK_INTEGRATION_TYPE = "slack"
)

// UpsertGoogleUser upserts slack user
func UpsertSlackUser(slackID, slackName, slackHandle, email, profileImage, slackToken, slackTeamId, sessionSecret string) (*models.UserSocialAuth, error) {
	var user models.UserSocialAuth
	var err error
	err = postgres.GetDB().Table("users").Where("UPPER(email) = ? ", strings.ToUpper(email)).First(&user).Error

	if err != nil {
		tempUser, err := NewSlackUserfromAuth(slackHandle, slackName, email, profileImage, slackID, slackTeamId, slackToken)

		if err != nil {
			return nil, err
		}
		user = *tempUser
	} else {

		if len(user.User.AvatarUrl) == 0 {
			user.User.AvatarUrl = profileImage
		}
		if len(user.User.FullName) == 0 {
			user.User.FullName = slackName
		}
		err = postgres.GetDB().Save(&user.User).Error
		if err != nil {
			return nil, err
		}
	}
	if len(profileImage) != 0 {
		//go ResizeAvtarUrl(user.ID, profileImage)
	}

	return &user, nil
}

// NewSlackUserfromAuth returns new slack user object from oauth
func NewSlackUserfromAuth(userName, name, email, profileImage, slackUserId, slackTeamId, slackToken string) (*models.UserSocialAuth, error) {
	var err error
	user := &models.User{
		Username: userName,
		FullName: name,
		Email: sql.NullString{
			String: email,
			Valid:  email != "",
		},
		AvatarUrl: profileImage,
	}
	newUserInstance := queries.App.UserQueries.CreateUser(user)
	if newUserInstance.ID != 0 {
		return nil, errors.New("Failed: queries.App.UserQueries.CreateUser /. NewSlackUserfromAuth")
	}

	metadata, _ := json.Marshal(slackToken)
	slackUser := &models.UserSocialAuth{
		UserID:       user.ID,
		ProviderName: SLACK_PROVIDER,
		ProviderID:   slackUserId + ":" + slackTeamId,
		Metadata:     metadata,
	}
	err = user2.App.Repo.CreateUserSocialAuth(slackUser)
	if err != nil {
		return nil, err
	}
	return slackUser, err
}

func (s *slackService) SendSlackMessageReply(message SlackMessagePayload, accessToken string) {
	url := "https://slack.com/api/chat.postMessage"
	method := "POST"

	payload, _ := json.Marshal(message)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(string(payload)))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func (s slackService) SendSlackCommonErrorResponse(channelID, threadTs, accessToken string) {
	messagePayload := SlackMessagePayload{
		Text:      "Error capturing message",
		ChannelID: channelID,
		ThreadTs:  threadTs,
	}
	App.Service.SendSlackMessageReply(messagePayload, accessToken)
}

func (s slackService) GetSlackThreadMessages(channelID, threadTs, accessToken string) []byte {
	url := fmt.Sprintf("https://slack.com/api/conversations.replies?channel=%s&ts=%s", channelID, threadTs)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(string(body))
	return body
}

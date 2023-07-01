package integrations

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type AuthUser struct {
	ID          string `json:"id"`
	AccessToken string `json:"access_token"`
}

type TeamData struct {
	ID string `json:"id"`
}

type SlackAuthData struct {
	Sub       string `json:"sub"`
	Nonce     string `json:"nonce"`
	Email     string `json:"email"`
	Team      string `json:"https://slack.com/team_id"`
	AvatarURL string `json:"picture"`
	Name      string `json:"name"`
}

type SlackAuthResponse struct {
	Data SlackAuthData `json:"data"`
}

type UserType struct {
	Email       string `json:"email"`
	DisplayName string `json:"real_name"`
	Avatar      string `json:"image_192"`
}
type SlackUserType struct {
	OK          bool     `json:"ok"`
	UserProfile UserType `json:"profile"`
}

type SlackTeamResponse struct {
	OK   bool      `json:"ok"`
	Team SlackTeam `json:"team"`
}

type SlackTeam struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon struct {
		Image132 string `json:"image_132"`
	} `json:"icon"`
}

type SlackResponse struct {
	Ok          bool   `json:"ok"`
	Error       string `json:"error"`
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
}
type SlackTeamMember struct {
	Id      string   `json:"id"`
	IsAdmin string   `json:"is_admin"`
	IsBot   bool     `json:"is_bot"`
	Deleted bool     `json:"deleted"`
	Profile UserType `json:"profile"`
}
type cursor struct {
	NextCursor string `json:"next_cursor"`
}
type TeamUsersList struct {
	Ok               bool              `json:"ok"`
	Members          []SlackTeamMember `json:"members"`
	ResponseMetadata cursor            `json:"response_metadata"`
	ErrorMessage     string            `json:"error"`
}

type SlackPostMessageResponse struct {
	Ok      bool                   `json:"ok"`
	Message map[string]interface{} `json:"message"`
	Channel string                 `json:"channel"`
	Ts      string                 `json:"ts"`
}

func SendMessageToSlackChannel(ctx context.Context, message, accessToken, threadTs, channelID string, blocks []map[string]interface{}) (*SlackPostMessageResponse, error) {

	// var response *slack.OAuthResponse
	//sendMessageToSlackChannel()

	postMessageUrl := "https://slack.com/api/chat.postMessage"
	contentType := "application/json" //"application/json"

	values := map[string]interface{}{
		"text":      message,
		"blocks":    blocks,
		"channel":   channelID,
		"thread_ts": threadTs,
	}
	fmt.Println("values", values)
	valuesStr, _ := json.Marshal(values)
	payload := strings.NewReader(string(valuesStr))
	client := &http.Client{}
	req, err := http.NewRequest("POST", postMessageUrl, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	result, _ := ioutil.ReadAll(res.Body)
	var slackResponse *SlackPostMessageResponse
	err = json.Unmarshal(result, &slackResponse)
	if err != nil {
		println("err", err.Error())
		return nil, err
	}
	slackMessage := slackResponse.Message
	println("slackMessage:", slackMessage)
	println("slack...", string(result))
	println("===>>>>>>>", slackResponse.Ok, string(result))
	if slackResponse.Ok {
		return slackResponse, nil
	}
	return nil, errors.New("sending slack message failed")
}

func GetSlackProfile(accessToken string, slackUserId string) (*SlackUserType, error) {
	req1, _ := http.NewRequest("GET", "https://slack.com/api/users.profile.get?user="+slackUserId+"&pretty=1", nil)
	req1.Header.Add("Authorization", "Bearer "+accessToken)

	client1 := &http.Client{}
	resp1, err := client1.Do(req1)
	if err != nil {
		return nil, errors.New("error while requesting for User Profile")
	}
	var thisUser SlackUserType

	if err := json.NewDecoder(resp1.Body).Decode(&thisUser); err != nil {
		return nil, errors.New("error while reading response")
	}
	return &thisUser, nil
}

func GetSlackTeam(accessToken string, slackTeamID string) (*SlackTeamResponse, error) {
	req1, _ := http.NewRequest("GET", "https://slack.com/api/team.info?team="+slackTeamID, nil)
	req1.Header.Add("Authorization", "Bearer "+accessToken)

	client1 := &http.Client{}
	resp1, err := client1.Do(req1)
	if err != nil {
		return nil, errors.New("error while requesting for User Profile")
	}
	var response SlackTeamResponse

	if err := json.NewDecoder(resp1.Body).Decode(&response); err != nil {
		return nil, errors.New("error while reading response")
	}
	return &response, nil
}

func GetSlackTeamMembers(accessToken string, slackTeamID string) (*TeamUsersList, error) {
	var teamMembers TeamUsersList
	for {
		var req1 *http.Request
		if teamMembers.ResponseMetadata.NextCursor == "" {
			req1, _ = http.NewRequest("GET", "https://slack.com/api/users.list?limit=200&team_id="+slackTeamID+"&pretty=1", nil)
		} else {
			req1, _ = http.NewRequest("GET", "https://slack.com/api/users.list?cursor="+teamMembers.ResponseMetadata.NextCursor+"&limit=200&team_id="+slackTeamID+"&pretty=1", nil)
		}

		req1.Header.Add("Authorization", "Bearer "+accessToken)

		client := &http.Client{}
		resp1, err := client.Do(req1)
		if err != nil {
			fmt.Println("error while team members")
			return nil, err
		}
		var responseTeamMembers TeamUsersList

		if err := json.NewDecoder(resp1.Body).Decode(&responseTeamMembers); err != nil {
			fmt.Println("error while reading response")
		}
		fmt.Println("Members list from team ====>", responseTeamMembers)
		teamMembers.Members = append(teamMembers.Members, responseTeamMembers.Members...)
		teamMembers.ResponseMetadata.NextCursor = responseTeamMembers.ResponseMetadata.NextCursor
		if responseTeamMembers.ResponseMetadata.NextCursor == "" {
			break
		}
	}
	return &teamMembers, nil
}

type SlackUserGroups struct {
	Id                  string      `json:"id"`
	TeamId              string      `json:"team_id"`
	IsUsergroup         bool        `json:"is_usergroup"`
	IsSubteam           bool        `json:"is_subteam"`
	Name                string      `json:"name"`
	Description         string      `json:"description"`
	Handle              string      `json:"handle"`
	IsExternal          bool        `json:"is_external"`
	DateCreate          int         `json:"date_create"`
	DateUpdate          int         `json:"date_update"`
	DateDelete          int         `json:"date_delete"`
	AutoType            interface{} `json:"auto_type"`
	AutoProvision       bool        `json:"auto_provision"`
	EnterpriseSubteamId string      `json:"enterprise_subteam_id"`
	CreatedBy           string      `json:"created_by"`
	UpdatedBy           string      `json:"updated_by"`
	DeletedBy           interface{} `json:"deleted_by"`
	Users               []string    `json:"users"`
	UserCount           int         `json:"user_count"`
	ChannelCount        int         `json:"channel_count"`
}

type SlackUserGroupResponse struct {
	Ok         bool              `json:"ok"`
	UserGroups []SlackUserGroups `json:"usergroups"`
}

func GetSlackUserGroups(accessToken, slackTeamID string) (*SlackUserGroupResponse, error) {
	req1, _ := http.NewRequest("GET", "https://slack.com/api/usergroups.list?include_count=true&include_disabled=false&include_users=true&team_id="+slackTeamID, nil)
	req1.Header.Add("Authorization", "Bearer "+accessToken)

	client1 := &http.Client{}
	resp1, err := client1.Do(req1)
	if err != nil {
		return nil, errors.New("error while requesting for User Groups")
	}
	var response SlackUserGroupResponse
	if err := json.NewDecoder(resp1.Body).Decode(&response); err != nil {
		return nil, errors.New("error while reading response")
	}
	return &response, nil
}

func DeleteSlackMessage(accessToken, channelID, messageTs string) error {
	url := fmt.Sprintf("https://slack.com/api/chat.delete?channel=%s&ts=%s", channelID, messageTs)
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))
	return nil
}

func SlackUpdateMessage(channelID string, ts string, accessToken string, blocks interface{}) {
	postMessageUrl := "https://slack.com/api/chat.update"
	contentType := "application/json" //"application/json"
	values := map[string]interface{}{
		"text":    "",
		"blocks":  blocks,
		"channel": channelID,
		"ts":      ts,
	}
	fmt.Println("values", values)
	valuesStr, _ := json.Marshal(values)
	payload := strings.NewReader(string(valuesStr))
	client := &http.Client{}
	req, err := http.NewRequest("POST", postMessageUrl, payload)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", "Bearer "+accessToken)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error on slack update message", err)
	}
	defer res.Body.Close()
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error in reading slack update message response", err)
	}
	fmt.Println("slack message update...", string(result))
	return
}

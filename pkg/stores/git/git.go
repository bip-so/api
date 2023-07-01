package stores

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
)

type GitUserRepo struct {
	UserID    string
	UserName  string
	ProductID string
	PageID    string
}

type GitBlockV1 struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Text       string `json:"text"`
	URL        string `json:"url"`
	UserID     string `json:"userId"`
	MessageID  string `json:"messageId"`
	TweetID    string `json:"tweetId"`
	PageID     string `json:"pageId"`
	Position   int    `json:"position"`
	Properties string `json:"properties"`
}

type GitBlockV2 struct {
	ID         string                   `json:"-"` //has to be uuid
	UUID       string                   `json:"uuid"`
	Version    uint                     `json:"version"`
	Type       string                   `json:"type"`
	Rank       int32                    `json:"rank"`
	Children   []map[string]interface{} `json:"children"`
	Attributes map[string]interface{}   `json:"attributes"`
	CreatedAt  time.Time                `json:"createdAt"`
	UpdatedAt  time.Time                `json:"updatedAt"`
}

type GitLog struct {
	ID          string    `json:"id"`
	Message     string    `json:"message"`
	UserID      string    `json:"-"`
	AuthorEmail string    `json:"authorEmail"`
	CreatedAt   time.Time `json:"createdAt"`
}

type GitAttribution struct {
	UserID      string `json:"-"`
	AuthorEmail string `json:"authorEmail"`
	Edits       int    `json:"edits"`
}

type GitBlockAttribution struct {
	UserID      string    `json:"-"`
	AuthorEmail string    `json:"authorEmail"`
	BlockID     string    `json:"blockId"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func CreateGitUserRepo(userID string, userName string, productID string, pageID string) (*GitUserRepo, error) {
	if userID == "" || userName == "" || productID == "" || pageID == "" {
		return nil, errors.New("CreateGitUserRepo Failed: Fields should not be empty")
	}
	return &GitUserRepo{
		UserID:    userID,
		UserName:  userName,
		ProductID: productID,
		PageID:    pageID,
	}, nil
}

func (gitUserRepo *GitUserRepo) getUserEmail() string {
	return gitUserRepo.UserID + "@bip.so"
}

var freeGitAPIHostIdx *int

func getFreeGitAPIHost() string {
	totalAPIHosts := len(configs.GetGitConfig().Hosts)
	if freeGitAPIHostIdx == nil {
		zero := 0
		freeGitAPIHostIdx = &zero
	} else {
		inc := *freeGitAPIHostIdx + 1
		if inc >= totalAPIHosts {
			inc = 0
		}
		freeGitAPIHostIdx = &inc
	}
	gitAPIHost := configs.GetGitConfig().Hosts[*freeGitAPIHostIdx]
	fmt.Println("gitAPIHost", gitAPIHost, *freeGitAPIHostIdx)
	return gitAPIHost
}

func (gitUserRepo *GitUserRepo) getGitRequest(reqType string, path string, body *bytes.Buffer) (*http.Request, error) {
	gitAPIHost := getFreeGitAPIHost()
	req, err := http.NewRequest(reqType,
		gitAPIHost+path,
		body)
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+configs.GetGitConfig().Secret)
	req.Header.Set("x-studio-id", gitUserRepo.ProductID)
	req.Header.Set("x-page-id", gitUserRepo.PageID)
	req.Header.Set("x-user-email", gitUserRepo.getUserEmail())
	req.Header.Set("x-user-name", gitUserRepo.UserName)
	return req, nil
}

func (gitUserRepo *GitUserRepo) CreateSnapshot(blocks *[]interface{}, fromBranchName, branchName string, message string) error {
	fmt.Println("Inside Create Snapshot")
	//fmt.Println(blocks)
	sendBody := map[string]interface{}{
		"blocks":         blocks,
		"branchName":     branchName,
		"fromBranchName": fromBranchName,
		"message":        message,
	}
	jsonData, _ := json.Marshal(sendBody)
	//fmt.Println("jsonData as string : ???")
	//fmt.Println(string(jsonData))
	client := http.Client{}
	req, err := gitUserRepo.getGitRequest("POST", "/api/snapshot/create", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		resp, err = client.Do(req)
	}
	fmt.Println(resp.StatusCode)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var respBody struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return err
	}
	if respBody.Success {
		return nil
	}
	return errors.New("there was some problem in commiting")
}

func (gitUserRepo *GitUserRepo) MergeMergeRequest(toBranchName string, fromBranchName string, mergeStatus string, fromBranchCreatedCommitID string, changesAccepted *map[string]interface{}) (string, string, string, error) {

	sendBody := map[string]interface{}{
		"toBranchName":   toBranchName,
		"fromBranchName": fromBranchName,
		"mergeType":      mergeStatus,
		// "fromBranchCreatedCommitId": fromBranchCreatedCommitID,
		"changesAccepted": changesAccepted,
	}

	jsonData, _ := json.Marshal(sendBody)

	client := http.Client{}
	req, err := gitUserRepo.getGitRequest("POST", "/api/mergereq/merge", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", "", err
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		resp, err = client.Do(req)
	}
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()
	var respBody struct {
		Success bool `json:"success"`
		Data    struct {
			CommitID     string `json:"commitId"`
			SrcCommitID  string `json:"srcCommitId"`
			DestCommitID string `json:"destCommitId"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return "", "", "", err
	}
	if respBody.Success {
		return respBody.Data.CommitID, respBody.Data.SrcCommitID, fromBranchCreatedCommitID, nil
	}
	return "", "", "", errors.New("there was some problem in merging merge request")
}

func (gitUserRepo *GitUserRepo) FetchBranchFromGit(branchName string) ([]*GitBlockV2, error) {

	client := http.Client{}
	req, err := gitUserRepo.getGitRequest("GET", "/api/snapshot/branch/"+branchName, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		resp, err = client.Do(req)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var respBody struct {
		Success bool          `json:"success"`
		Data    []*GitBlockV2 `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}
	if respBody.Success {
		return respBody.Data, nil
	}
	return nil, errors.New("there was some problem in fetching branch")
}

func (gitUserRepo *GitUserRepo) FetchCommitFromGit(commitID string) ([]*GitBlockV2, error) {

	client := http.Client{}
	req, err := gitUserRepo.getGitRequest("GET", "/api/snapshot/get/"+commitID, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		resp, err = client.Do(req)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var respBody struct {
		Success bool          `json:"success"`
		Data    []*GitBlockV2 `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}
	if respBody.Success {
		return respBody.Data, nil
	}
	return nil, errors.New("Version issue on the blocks data")
}

func (gitUserRepo *GitUserRepo) CreateBranch(fromBranchName, branchName string) (string, error) {

	sendBody := map[string]interface{}{
		"branchName": branchName,
	}

	jsonData, _ := json.Marshal(sendBody)

	client := http.Client{}
	req, err := gitUserRepo.getGitRequest("POST", "/api/branch/create/"+fromBranchName, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		resp, err = client.Do(req)
	}
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var respBody struct {
		Success bool   `json:"success"`
		Data    string `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return "", err
	}
	if respBody.Success {
		return respBody.Data, nil
	}
	return "", errors.New("there was some problem in deleting branch")
}

func (gitUserRepo *GitUserRepo) DeleteBranch(branchName string) error {

	client := http.Client{}
	req, err := gitUserRepo.getGitRequest("DELETE", "/api/branch/"+branchName, bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		resp, err = client.Do(req)
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var respBody struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return err
	}
	if respBody.Success {
		return nil
	}
	return errors.New("there was some problem in deleting branch")
}

func (gitUserRepo *GitUserRepo) FetchBranchHistoryFromGit(branchName string, startCommitID string) ([]*GitLog, string, error) {

	query := ""
	if startCommitID != "" {
		query = "?start=" + startCommitID
	}

	client := http.Client{}
	//fmt.Println("branchName+query", branchName+query)
	req, err := gitUserRepo.getGitRequest("GET", "/api/history/branch/"+branchName+query, bytes.NewBuffer([]byte{}))
	//fmt.Println("req")
	//fmt.Println(req)
	if err != nil {
		return nil, "", err
	}
	resp, err := client.Do(req)
	//fmt.Println("resp")
	//fmt.Println(resp)
	//fmt.Println(resp.Status)
	if err != nil || resp.StatusCode != http.StatusOK {
		resp, err = client.Do(req)
	}
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	var respBody struct {
		Success bool      `json:"success"`
		Data    []*GitLog `json:"data"`
		Next    string    `json:"next"`
	}
	fmt.Println(respBody)
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, "", err
	}
	if respBody.Success {
		return respBody.Data, respBody.Next, nil
	}
	fmt.Println(respBody.Success)
	return nil, "", errors.New("there was some problem in fetching history")
}

func (gitUserRepo *GitUserRepo) FetchAllAttributionsForBranch(branchName string, lastSyncedCommitID string) ([]*GitAttribution, string, error) {

	client := http.Client{}
	req, err := gitUserRepo.getGitRequest("GET", "/api/attribution/all/branch/"+branchName+"?lastSyncedCommitId="+lastSyncedCommitID, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, "", err
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		resp, err = client.Do(req)
	}
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	var respBody struct {
		Success       bool              `json:"success"`
		Data          []*GitAttribution `json:"data"`
		StartCommitID string            `json:"startCommitId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, "", err
	}
	if respBody.Success {
		return respBody.Data, respBody.StartCommitID, nil
	}
	return nil, "", errors.New("there was some problem in fetching all attributes")
}

func (gitUserRepo *GitUserRepo) FetchLatestBlockAttributionsForBranch(branchName string) ([]*GitBlockAttribution, error) {

	client := http.Client{}
	req, err := gitUserRepo.getGitRequest("GET", "/api/attribution/branch/"+branchName, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		resp, err = client.Do(req)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var respBody struct {
		Success bool                   `json:"success"`
		Data    []*GitBlockAttribution `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}
	if respBody.Success {
		return respBody.Data, nil
	}
	return nil, errors.New("there was some problem in fetching all attributes")
}

func (gitUserRepo *GitUserRepo) CreateSnapshotForMessageBlock(blocks *[]interface{}, branchName, messageBlockId string) error {

	sendBody := map[string]interface{}{
		"blocks":         blocks,
		"branchName":     branchName,
		"messageBlockId": messageBlockId,
	}

	jsonData, _ := json.Marshal(sendBody)

	client := http.Client{}
	req, err := gitUserRepo.getGitRequest("POST", "/api/snapshot/block/create", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		resp, err = client.Do(req)
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var respBody struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return err
	}
	if respBody.Success {
		return nil
	}
	return errors.New("there was some problem in commiting")
}

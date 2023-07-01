package slack2

import (
	"errors"
	"fmt"
	"github.com/slack-go/slack"
	"gitlab.com/phonepost/bip-be-platform/internal/canvasbranch"
	"gitlab.com/phonepost/bip-be-platform/internal/canvasrepo"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"net/url"
	"strings"
)

func SlackBipSearchTreeBuilderHandler(searchString string, studio *models.Studio, user *models.User) (string, error) {
	repoAndCollectionRows := canvasbranch.App.Repo.QueryDB(searchString, studio.ID)
	resp := canvasbranch.App.Service.ProcessSearchDump(repoAndCollectionRows, user.ID, studio.ID, "false")
	treeMessage := ""
	repos := resp["repos"]
	for i, repo := range repos {
		if searchString == "" || strings.Contains(strings.ToLower(repo["name"].(string)), searchString) {
			canvasUrlTitle := notifications.App.Service.GenerateCanvasUrlTitle(repo["name"].(string), repo["defaultBranchID"].(uint64))
			canvasRepoName := repo["name"].(string)
			canvasRepoName = strings.ReplaceAll(canvasRepoName, ">", "")
			canvasRepoName = strings.ReplaceAll(canvasRepoName, "<", "")
			canvasRepoName = strings.ReplaceAll(canvasRepoName, "\n", "")
			treeMessage += fmt.Sprintf("*<%s|%s>*", fmt.Sprintf("%s/%s/%s", configs.GetAppInfoConfig().FrontendHost, studio.Handle, url.QueryEscape(canvasUrlTitle)), canvasRepoName)
			branch, err := queries.App.BranchQuery.GetBranchByID(repo["defaultBranchID"].(uint64))
			if err != nil {
				continue
			}
			if branch.PublicAccess == "private" {
				treeMessage += ":lock:"
			}
			treeMessage += "\n"
		}
		if i >= 39 {
			if (len(repos) - i) > 0 {
				treeMessage += "\n" + fmt.Sprintf("<%s|+%d more>", fmt.Sprintf("%s/%s", configs.GetAppInfoConfig().FrontendHost, studio.Handle), len(repos)-i)
			}
			break
		}
	}
	var slackText string
	if len(repos) > 0 {
		slackText = fmt.Sprintf("*Canvases List*\n\n%s", treeMessage)
	} else {
		slackText = "No Canvases found"
	}
	return slackText, nil
}

func SlackBipNewSlashCommandHandler(slashCommand slack.SlashCommand, studioInstance *models.Studio, userInstance *models.User) (string, error) {
	collectionInstance, err := LastUserAccessCollectionInStudio(studioInstance.ID, userInstance.ID)
	if err != nil || collectionInstance == nil {
		fmt.Println("Error in creating collection instance", err)
		return "", errors.New("User doesn't have permission to create canvases")
	}
	newCanvasPost := canvasrepo.NewCanvasRepoPost{
		Name:         slashCommand.Text,
		CollectionID: collectionInstance.ID,
		Position:     1,
	}
	repo, err := canvasrepo.App.Controller.CreateCanvasRepo(newCanvasPost, userInstance.ID, studioInstance.ID, *userInstance, collectionInstance.PublicAccess)
	if err != nil {
		fmt.Println("Error in creating canvas Repo", err)
		return "", err
	}
	var text string
	canvasUrl := notifications.App.Service.GenerateCanvasBranchUrl(repo.Key, repo.Name, studioInstance.ID, *repo.DefaultBranchID)
	text = "New Canvas Created\n" + canvasUrl
	return text, nil
}

func LastUserAccessCollectionInStudio(studioID, userID uint64) (*models.Collection, error) {
	var collection *models.Collection
	collections, err := queries.App.CollectionQuery.GetCollections(map[string]interface{}{"studio_id": studioID, "is_archived": false})
	if err != nil {
		return nil, err
	}
	for i, col := range collections {
		col = collections[len(collections)-i-1]
		hasPermission, err := permissions.App.Service.CanUserDoThisOnCollection(userID, studioID, col.ID, permissiongroup.COLLECTION_VIEW_METADATA)
		if err != nil {
			continue
		}
		if hasPermission {
			collection = &col
			break
		}
	}
	return collection, nil
}

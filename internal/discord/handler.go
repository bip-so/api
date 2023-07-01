package discord

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/phonepost/bip-be-platform/internal/canvasbranch"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

func AutoCompleteBipNewHandler(body Interaction, activeDiscord []models.StudioIntegration, user *models.User) ([]map[string]interface{}, error) {
	studio := activeDiscord[0].Studio
	focusedOption := Option{}
	collectionOption := Option{}
	for _, option := range body.Data.Options {
		if option.Focused {
			focusedOption = option
		}
		if option.Name == "collection" {
			collectionOption = option
		}
	}
	searchString := focusedOption.Value
	suggestions := []map[string]interface{}{}
	rootCollections := canvasbranch.RootSerializedByStudioID(studio.ID, user.ID)
	if focusedOption.Name == "collection" {
		for _, col := range *rootCollections {
			if searchString == "" || strings.Contains(strings.ToLower(col.Name), searchString) {
				colName := col.Name
				if len(colName) > 40 {
					colName = colName[:40] + "..."
				}
				suggestion := map[string]interface{}{
					"name":  colName,
					"value": utils.String(col.Id),
				}
				suggestions = append(suggestions, suggestion)
			}
			if len(suggestions) >= 25 {
				break
			}
		}
	} else if focusedOption.Name == "parent-canvas" {
		collectionID := utils.Uint64(collectionOption.Value)
		for _, col := range *rootCollections {
			if col.Id == collectionID {
				for _, repo := range col.Repos {
					if repo.IsLanguageCanvas {
						continue
					}
					fmt.Println(repo.Name, searchString)
					if searchString == "" || strings.Contains(strings.ToLower(repo.Name), searchString) {
						name := repo.Name
						if len(name) > 40 {
							name = name[:40] + "..."
						}
						suggestion := map[string]interface{}{
							"name":  name,
							"value": utils.String(repo.ID),
						}
						suggestions = append(suggestions, suggestion)
					}
				}
			}
			if len(suggestions) >= 25 {
				break
			}
		}
	}
	fmt.Println("suggestions here", suggestions)
	return suggestions, nil
}

func AutoCompleteBipSearchHandler(body Interaction, activeDiscord []models.StudioIntegration, user *models.User) ([]map[string]interface{}, error) {
	option := body.Data.Options[0]
	studio := activeDiscord[0].Studio
	suggestions := []map[string]interface{}{}
	searchString := option.Value
	repoAndCollectionRows := canvasbranch.App.Repo.QueryDB(searchString, studio.ID)
	resp := canvasbranch.App.Service.ProcessSearchDump(repoAndCollectionRows, user.ID, studio.ID, "false")
	repos := resp["repos"]
	for _, repo := range repos {
		if searchString == "" || strings.Contains(strings.ToLower(repo["name"].(string)), searchString) {
			name := repo["name"].(string)
			if len(name) > 40 {
				name = name[:40] + "..."
			}
			suggestion := map[string]interface{}{
				"name":  name,
				"value": utils.String(repo["id"].(uint64)),
			}
			suggestions = append(suggestions, suggestion)
		}
		if len(suggestions) >= 25 {
			break
		}
	}
	return suggestions, nil
}

func BipSearchTreeBuilderHandler(body Interaction, activeDiscord []models.StudioIntegration, user *models.User) (*discordgo.MessageEmbed, error) {
	option := body.Data.Options[0]
	studio := activeDiscord[0].Studio
	searchString := option.Value
	repoAndCollectionRows := canvasbranch.App.Repo.QueryDB(searchString, studio.ID)
	resp := canvasbranch.App.Service.ProcessSearchDump(repoAndCollectionRows, user.ID, studio.ID, "false")
	treeMessage := ""
	repos := resp["repos"]
	for i, repo := range repos {
		if searchString == "" || strings.Contains(strings.ToLower(repo["name"].(string)), searchString) {
			canvasUrlTitle := notifications.App.Service.GenerateCanvasUrlTitle(repo["name"].(string), repo["defaultBranchID"].(uint64))
			treeMessage += fmt.Sprintf("**[%s](%s)**", repo["name"].(string), fmt.Sprintf("%s/%s/%s", configs.GetAppInfoConfig().FrontendHost, studio.Handle, url.QueryEscape(canvasUrlTitle)))
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
				treeMessage += "\n" + fmt.Sprintf("[+%d more](%s)", len(repos)-i, fmt.Sprintf("%s/%s", configs.GetAppInfoConfig().FrontendHost, studio.Handle))
			}
			break
		}
	}
	embed := &discordgo.MessageEmbed{
		Title:       "Canvases List",
		Description: treeMessage,
		Color:       0x44B244,
		Type:        "rich",
	}
	if treeMessage == "" {
		embed.Title = "The Canvas does not exist. Kindly check"
	}
	return embed, nil
}

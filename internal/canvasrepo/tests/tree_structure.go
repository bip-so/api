package main

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/cmd/api"
	"gitlab.com/phonepost/bip-be-platform/internal/canvasrepo"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/integrations"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"net/url"
	"strings"
)

func main() {
	configs.InitConfig(".env", ".")
	postgres.InitDB()
	integrations.InitDiscordGo()
	redis.InitRedis()
	api.InitAllApps()
	var collections []models.Collection
	err := postgres.GetDB().Model(models.Collection{}).Where("studio_id = ?", 383).Preload("Studio").Find(&collections).Error
	if err != nil {
		fmt.Println("Error in fetching collections", err)
	}
	CreateStudioFileTree(collections)
	//integrations.CreateChannel("958371136368959579")
}

func CreateStudioFileTree(collections []models.Collection) {
	// loop the collections
	// initiate a string and start with the collection name with no space at starting.
	// Next create another method with recursion which will recursive start adding the canvasRepos to the string
	fmt.Println(len(collections))
	for _, collection := range collections {
		var user *models.User
		postgres.GetDB().Model(models.User{}).Where("id = ?", 82).First(&user)
		embed := canvasrepo.App.Service.CreateStudioFileTree(collection.Studio, user)
		//tree := ""
		//indent := 0
		//tree += fmt.Sprintf("**%s**\n", collection.Name)
		////fmt.Println(tree)
		//canvasTree, canvasesCount := CreateCanvasRepoFileTree(collection.ID, 0, "", indent, 0)
		//diffCount := collection.ComputedRootCanvasCount - canvasesCount
		//canvasTree += "\n\n" + fmt.Sprintf("[+%d more](%s)", diffCount, fmt.Sprintf("%s/%s", configs.GetAppInfoConfig().FrontendHost, collection.Studio.Handle))
		//tree += canvasTree + "\n"
		//fmt.Println(tree)
		//embed := &discordgo.MessageEmbed{
		//	Title:       "Canvases List",
		//	Description: tree,
		//	Color:       0x44B244,
		//	Type:        "rich",
		//}
		integrations.SendDiscordEmbedToChannel("1027442589789601832", embed)
		break
	}
	//discordComponents := []interface{}{
	//	notifications.ActionRowsComponent{
	//		Type: 1,
	//		Components: []interface{}{
	//			notifications.MessageBtnComponent{
	//				Type:     2,
	//				Label:    "Check your access",
	//				Style:    2,
	//				CustomID: "checkCanvasAccess",
	//			},
	//		},
	//	},
	//}
	//fmt.Println(discordComponents)
	//msg, err := integrations.SendDiscordDMMessageToChannel("1027442589789601832", []string{}, discordComponents)
	//fmt.Println("msg", msg.GuildID, err)
}

func CreateCanvasRepoFileTree(collectionID uint64, parentCanvasRepoID uint64, canvasTree string, indent int, canvasesCount int) (string, int) {
	var subCanvas []models.CanvasRepository
	if parentCanvasRepoID == 0 {
		postgres.GetDB().Model(models.CanvasRepository{}).Where("collection_id = ? and parent_canvas_repository_id is null", collectionID).Preload("Studio").Preload("DefaultBranch").Find(&subCanvas)
	} else {
		postgres.GetDB().Model(models.CanvasRepository{}).Where("collection_id = ? and parent_canvas_repository_id = ?", collectionID, parentCanvasRepoID).Preload("Studio").Preload("DefaultBranch").Find(&subCanvas)
	}
	indent += 2
	for _, canvas := range subCanvas {
		fmt.Println("canvasName", canvas.Name, canvasesCount)
		canvasesCount++
		if canvasesCount > 20 {
			return canvasTree, canvasesCount
		}
		canvasUrlTitle := notifications.App.Service.GenerateCanvasUrlTitle(canvas.Name, *canvas.DefaultBranchID)
		canvasTree += strings.Repeat("　", indent) + fmt.Sprintf("**[%s](%s)**", canvas.Name, fmt.Sprintf("%s/%s/%s", configs.GetAppInfoConfig().FrontendHost, canvas.Studio.Handle, url.QueryEscape(canvasUrlTitle)))
		//fmt.Println(tree)
		if canvas.DefaultBranch.PublicAccess == "private" {
			canvasTree += ":lock:"
		}
		canvasTree += "\n"
		var languageCanvases []models.CanvasRepository
		postgres.GetDB().Model(models.CanvasRepository{}).Where("default_language_canvas_repo_id = ? and is_language_canvas = true and is_archived = false", canvas.ID).Order("position ASC").Find(&languageCanvases)
		for _, lcanvas := range languageCanvases {
			canvasesCount++
			if canvasesCount > 20 {
				break
			}
			lcanvasUrlTitle := notifications.App.Service.GenerateCanvasUrlTitle(lcanvas.Name, *lcanvas.DefaultBranchID)
			canvasTree += strings.Repeat("　", indent) + fmt.Sprintf("**[|अ %s](%s)**", lcanvas.Name, fmt.Sprintf("%s/%s/%s", configs.GetAppInfoConfig().FrontendHost, canvas.Studio.Handle, url.QueryEscape(lcanvasUrlTitle)))
			canvasTree += "\n"
		}
		var subCanvasTree string
		subCanvasTree, canvasesCount = CreateCanvasRepoFileTree(collectionID, canvas.ID, "", indent, canvasesCount)
		canvasTree += subCanvasTree
	}
	return canvasTree, canvasesCount
}

package notifications

import (
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

func (s notificationService) GetReelBlocksText(reel *models.Reel) string {
	contextData := map[string]interface{}{}
	err := json.Unmarshal(reel.ContextData, &contextData)
	if err != nil {
		fmt.Println("Unmarshal context data error:", err)
		return ""
	}
	selectedBlocks := map[string]map[string]BlockData{}
	err = json.Unmarshal(reel.SelectedBlocks, &selectedBlocks)
	blocksData := selectedBlocks["blocksData"]
	reelText := contextData["text"].(string)
	text := "*" + reel.CreatedByUser.Username + "*: " + reelText + "\n\n>>> "
	for _, block := range blocksData {
		if block.Type == models.BlockSimpleTableV1 {
			continue
		}
		format := "%s \n"
		if block.Type == models.BlockTypeHeading1 {
			for _, children := range block.Children {
				text = fmt.Sprintf(format+"%s", text, "**#"+children.Text+"**")
			}
		} else if block.Type == models.BlockTypeHeading2 {
			for _, children := range block.Children {
				text = fmt.Sprintf(format+"%s", text, "**##"+children.Text+"**")
			}
		} else if block.Type == models.BlockTypeHeading3 {
			for _, children := range block.Children {
				text = fmt.Sprintf(format+"%s", text, "**###"+children.Text+"**")
			}
		} else if block.Type == models.BlockTypeHeading4 {
			for _, children := range block.Children {
				text = fmt.Sprintf(format+"%s", text, "**####"+children.Text+"*")
			}
		} else if block.Type == models.BlockTypeHeading5 {
			for _, children := range block.Children {
				text = fmt.Sprintf(format+"%s", text, "**#####"+children.Text+"**")
			}
		} else if block.Type == models.BlockTypeHeading6 {
			for _, children := range block.Children {
				text = fmt.Sprintf(format+"%s", text, "**######"+children.Text+"**")
			}
		} else if block.Type == models.BlockTypeImage {
			for _, children := range block.Children {
				if children.Type == models.BlockTypeImage {
					text = fmt.Sprintf(format+"Image: %s", text, children.Url)
				}
			}
		} else if block.Type == models.BlockTypeAttachment {
			for _, children := range block.Children {
				if children.Type == models.BlockTypeAttachment {
					text = fmt.Sprintf(format+"File: %s", text, children.Url)
				}
			}
		} else if utils.Contains([]string{models.BlockTypeVideo, models.BlockTypeTweet}, block.Type) {
			for _, children := range block.Children {
				text = fmt.Sprintf(format+"Embed: %s", text, children.Url)
			}
		} else {
			for _, children := range block.Children {
				if children.Type == models.BlockTypeImage {
					text = fmt.Sprintf(format+"Image: %s", text, children.Url)
				} else if children.Type == models.BlockTypeAttachment {
					text = fmt.Sprintf(format+"File: %s", text, children.Url)
				} else if children.Type == models.BlockDrawIO {
					continue
				} else if children.Url != "" {
					text = fmt.Sprintf(format+"Embed: %s", text, children.Url)
				} else {
					text = fmt.Sprintf(format+"%s", text, children.Text)
				}
			}
		}
	}

	canvasRepo, _ := App.Repo.GetCanvasRepoByID(reel.CanvasRepositoryID)
	reelURL := configs.GetAppInfoConfig().FrontendHost
	if canvasRepo.ID != 0 {
		reelURL = App.Service.GenerateReelUUIDUrl(canvasRepo.Key, canvasRepo.Name, reel.StudioID, reel.CanvasBranchID, reel.UUID.String())
	}
	text = fmt.Sprintf("%s\n\n%s", text, reelURL)
	return text
}

func (s notificationService) GetBlockText(block *models.Block) string {
	childrenData := []BlockChildren{}
	json.Unmarshal(block.Children, &childrenData)
	var text string
	format := "%s \n"
	if block.Type == models.BlockTypeHeading1 {
		for _, children := range childrenData {
			text = fmt.Sprintf(format+"%s", text, "**#"+children.Text+"**")
		}
	} else if block.Type == models.BlockTypeHeading2 {
		for _, children := range childrenData {
			text = fmt.Sprintf(format+"%s", text, "**##"+children.Text+"**")
		}
	} else if block.Type == models.BlockTypeHeading3 {
		for _, children := range childrenData {
			text = fmt.Sprintf(format+"%s", text, "**###"+children.Text+"**")
		}
	} else if block.Type == models.BlockTypeHeading4 {
		for _, children := range childrenData {
			text = fmt.Sprintf(format+"%s", text, "**####"+children.Text+"*")
		}
	} else if block.Type == models.BlockTypeHeading5 {
		for _, children := range childrenData {
			text = fmt.Sprintf(format+"%s", text, "**#####"+children.Text+"**")
		}
	} else if block.Type == models.BlockTypeHeading6 {
		for _, children := range childrenData {
			text = fmt.Sprintf(format+"%s", text, "**######"+children.Text+"**")
		}
	} else if block.Type == models.BlockSimpleTableV1 {
		text = ""
	} else {
		for _, children := range childrenData {
			text = fmt.Sprintf(format+"%s", text, children.Text)
		}
	}
	return text
}

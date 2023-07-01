package parser2

import (
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gorm.io/datatypes"
)

type PlainTextBlockParser struct{}

func LevelMaker(level int) string {
	switch level {
	case 1:
		return Markdown_Level1_Indentation
	case 2:
		return Markdown_Level2_Indentation
	case 3:
		return Markdown_Level3_Indentation
	case 4:
		return Markdown_Level4_Indentation
	case 5:
		return Markdown_Level5_Indentation
	case 6:
		return Markdown_Level6_Indentation
	default:
		return Markdown_Level1_Indentation
	}

}

func UtilConvertByteJsonToMMap(content datatypes.JSON) []map[string]interface{} {
	var mm []map[string]interface{}
	_ = json.Unmarshal(content, &mm)
	return mm
}

type UserMentionObject struct {
	AvatarUrl             string `json:"avatarUrl"`
	CreatedByUserFullName string `json:"createdByUserFullName"`
	CreatedByUserUsername string `json:"createdByUserUsername"`
	FullName              string `json:"fullName"`
	Username              string `json:"username"`
}

func MakeMapFromInterface(b map[string]interface{}) map[string]string {
	var res = map[string]string{}
	for k, v := range b {
		if v == nil {
			continue
		}
		if reflect.TypeOf(v).String() == "string" {
			res[k] = v.(string)
		}
	}
	return res
}
func GetPreciseMapWithKey(blockMap []map[string]interface{}, key string) map[string]interface{} {
	var x map[string]interface{}
	for _, ov := range blockMap {
		for ik, value := range ov {
			if ik == "type" && value == key {
				x = ov
				return x
			}
		}
	}
	return x
}

func H1Extractor(blockMap []map[string]interface{}) string {
	content := Markdown_H1_Prefix + SimpleTextExtract(blockMap)
	return content
}

func H2Extractor(blockMap []map[string]interface{}) string {
	content := Markdown_H2_Prefix + SimpleTextExtract(blockMap)
	return content
}
func H3Extractor(blockMap []map[string]interface{}) string {
	content := Markdown_H3_Prefix + SimpleTextExtract(blockMap)
	return content
}
func H4Extractor(blockMap []map[string]interface{}) string {
	content := Markdown_H4_Prefix + SimpleTextExtract(blockMap)
	return content
}
func H5Extractor(blockMap []map[string]interface{}) string {
	content := Markdown_H5_Prefix + SimpleTextExtract(blockMap)
	return content
}
func H6Extractor(blockMap []map[string]interface{}) string {
	content := Markdown_H6_Prefix + SimpleTextExtract(blockMap)
	return content
}

func SubTextExtractor(blockMap []map[string]interface{}) string {
	content := MarkdownTilde + SimpleTextExtract(blockMap) + MarkdownTilde
	return content
}

func QuoteExtractor(blockMap []map[string]interface{}) string {
	content := Markdown_BLOCKQUOTE + SimpleTextExtract(blockMap)
	return content
}

func UserMentionTextExtractv1(blockMap []map[string]interface{}) string {
	//b1, _ := json.MarshalIndent(blockMap, "", "  ")
	//fmt.Println(string(b1))
	var content string
	for _, children := range blockMap {
		for key, innerBlock := range children {
			if key == "text" {
				content += innerBlock.(string)
			}
			if key == "mention" {
				var res = map[string]string{}
				for k, v := range innerBlock.(map[string]interface{}) {
					if reflect.TypeOf(v).String() == "string" {
						res[k] = v.(string)
					}
				}

				content += "  < Mention @" + res["username"] + "> "
			}
		}

	}
	return content
}

func UserMentionTextExtract(children map[string]interface{}) string {
	var content string
	for key, innerBlock := range children {
		if key == "text" {
			content += innerBlock.(string)
		}
		if key == "mention" {
			var res = map[string]string{}
			for k, v := range innerBlock.(map[string]interface{}) {
				if reflect.TypeOf(v).String() == "string" {
					res[k] = v.(string)
				}
			}
			username := res["username"]
			if username == "" {
				username = res["name"]
				content += "  < Mention @" + username + "> "
			} else {
				// it is a user so adding url link
				mentionURL := fmt.Sprintf("%s/@%s", configs.GetAppInfoConfig().FrontendHost, username)
				content += "  < Mention [@" + username + "](" + mentionURL + ")> "
			}
		}
	}
	return content
}

func PageMentionTextExtract(children map[string]interface{}) string {
	var content string
	for key, innerBlock := range children {
		if key == "text" {
			content += innerBlock.(string)
		}
		if key == "mention" {
			mentionData := innerBlock.(map[string]interface{})
			repoName := mentionData["repoName"].(string)
			fmt.Println(innerBlock, reflect.TypeOf(mentionData))
			studioID := mentionData["studioID"].(float64)
			fmt.Println("studio id", studioID)
			studioInstance, _ := App.Repo.GetStudioByID(uint64(studioID))
			canvasUrlTitle := notifications.App.Service.GenerateCanvasUrlTitle(mentionData["repoName"].(string), uint64(mentionData["studioID"].(float64)))
			mentionURL := fmt.Sprintf("%s/%s/%s", configs.GetAppInfoConfig().FrontendHost, url.QueryEscape(studioInstance.Handle), url.QueryEscape(canvasUrlTitle))
			content += "  [Canvas(" + repoName + ")](" + mentionURL + ") "
		}
	}
	return content
}

func SimpleTextExtract(blockMap []map[string]interface{}) string {
	//b1, _ := json.MarshalIndent(blockMap, "", "  ")
	//fmt.Println(string(b1))
	var content string
	// content += Markdown_Linebreak_Prefix
	for _, children := range blockMap {
		blockText, ok := children["text"].(string)
		if !ok {
			continue
		}
		attributeFound := false
		if children["link"] != nil && children["link"].(string) != "" {
			attributeFound = true
			content += fmt.Sprintf("[%s](%s)", blockText, children["link"].(string))
		} else {
			if children["bold"] != nil && children["italic"] != nil && children["bold"].(bool) && children["italic"].(bool) {
				attributeFound = true
				content += " " + Markdown_BOLDITALICS + strings.TrimSpace(blockText) + Markdown_BOLDITALICS + " "
			} else if children["bold"] != nil && children["bold"].(bool) {
				attributeFound = true
				content += " " + Markdown_BOLD + strings.TrimSpace(blockText) + Markdown_BOLD + " "
			} else if children["italic"] != nil && children["italic"].(bool) {
				attributeFound = true
				content += " " + Markdown_ITALICS + strings.TrimSpace(blockText) + Markdown_ITALICS + " "
			} else if children["strikethrough"] != nil && children["strikethrough"].(bool) {
				attributeFound = true
				content += " " + Markdown_STRIKETHROUGH + strings.TrimSpace(blockText) + Markdown_STRIKETHROUGH + " "
			} else if children["inlineCode"] != nil && children["inlineCode"].(bool) == true {
				content = fmt.Sprintf("`%s`", content)
			}
		}
		if !attributeFound {
			content += blockText
		}
	}
	// content += Markdown_Linebreak_Prefix
	// content += Markdown_Linebreak_Prefix
	return content
}

func ChildTextExtract(children map[string]interface{}) string {
	if children["text"] == nil {
		return ""
	}
	var content string
	blockText := children["text"].(string)
	attributeFound := false
	isLinkString := false
	if children["link"] != nil {
		_, isLinkString = children["link"].(string)
	}
	if children["link"] != nil && isLinkString && children["link"].(string) != "" {
		attributeFound = true
		content += fmt.Sprintf("[%s](%s)", blockText, children["link"].(string))
	} else {
		if children["bold"] != nil && children["italic"] != nil && children["bold"].(bool) && children["italic"].(bool) {
			attributeFound = true
			content += " " + Markdown_BOLDITALICS + strings.TrimSpace(blockText) + Markdown_BOLDITALICS + " "
		} else if children["bold"] != nil && children["strikethrough"] != nil && children["bold"].(bool) && children["strikethrough"].(bool) {
			attributeFound = true
			content += " " + "**~~" + strings.TrimSpace(blockText) + "~~**" + " "
		} else if children["bold"] != nil && children["bold"].(bool) {
			attributeFound = true
			content += " " + Markdown_BOLD + strings.TrimSpace(blockText) + Markdown_BOLD + " "
		} else if children["italic"] != nil && children["italic"].(bool) {
			attributeFound = true
			content += " " + Markdown_ITALICS + strings.TrimSpace(blockText) + Markdown_ITALICS + " "
		} else if children["strikethrough"] != nil && children["strikethrough"].(bool) {
			attributeFound = true
			content += " " + Markdown_STRIKETHROUGH + strings.TrimSpace(blockText) + Markdown_STRIKETHROUGH + " "
		} else if children["inlineCode"] != nil && children["inlineCode"].(bool) {
			attributeFound = true
			content += " " + "`" + strings.TrimSpace(blockText) + "`" + " "
		}
	}
	if !attributeFound {
		content += blockText
	}
	return content
}

func EmbedExtractor(children map[string]interface{}) string {
	if children["type"] == nil || children["url"] == nil {
		fmt.Println("children", children["type"])
		return ""
	}
	var content string
	// content += fmt.Sprintf("[%s](%s)", strings.Title(children["type"].(string)), children["url"].(string))
	content += "[" + strings.Title(children["type"].(string)) + "](" + children["url"].(string) + ")"
	return content
}

func YoutubeExtractor(blockMap []map[string]interface{}) string {
	var content string
	mainBlock := GetPreciseMapWithKey(blockMap, models.BlockTypeYoutube)
	mailBlockAsMap := MakeMapFromInterface(mainBlock)
	content += "[Youtube Link](" + mailBlockAsMap["url"] + ")"
	return content
}

func GoogleSheetExtractor(blockMap []map[string]interface{}) string {
	var content string
	mainBlock := GetPreciseMapWithKey(blockMap, models.BlockTypeGoogleSheet)
	mailBlockAsMap := MakeMapFromInterface(mainBlock)
	content += "[Googlesheet](" + mailBlockAsMap["url"] + ")"
	return content
}
func GoogleDriveExtractor(blockMap []map[string]interface{}) string {
	var content string
	mainBlock := GetPreciseMapWithKey(blockMap, models.BlockTypeGoogleDrive)
	mailBlockAsMap := MakeMapFromInterface(mainBlock)
	content += "[GoogleDrive](" + mailBlockAsMap["url"] + ")"
	return content
}
func LoomExtractor(blockMap []map[string]interface{}) string {
	var content string
	mainBlock := GetPreciseMapWithKey(blockMap, models.BlockTypeLoom)
	mailBlockAsMap := MakeMapFromInterface(mainBlock)
	content += "[Loom](" + mailBlockAsMap["url"] + ")"
	return content
}
func FigmaExtractor(blockMap []map[string]interface{}) string {
	var content string
	mainBlock := GetPreciseMapWithKey(blockMap, models.BlockTypeFigma)
	mailBlockAsMap := MakeMapFromInterface(mainBlock)
	content += "[Figma](" + mailBlockAsMap["url"] + ")"
	return content
}
func ReplitExtractor(blockMap []map[string]interface{}) string {
	var content string
	mainBlock := GetPreciseMapWithKey(blockMap, models.BlockTypeReplit)
	mailBlockAsMap := MakeMapFromInterface(mainBlock)
	content += "[Replit](" + mailBlockAsMap["url"] + ")"
	return content
}

func CodesandboxExtractor(blockMap []map[string]interface{}) string {
	var content string
	mainBlock := GetPreciseMapWithKey(blockMap, models.BlockTypeCodeSandBox)
	mailBlockAsMap := MakeMapFromInterface(mainBlock)
	content += "[codesandbox](" + mailBlockAsMap["url"] + ")"
	return content
}
func TweetExtractor(blockMap []map[string]interface{}) string {
	var content string
	mainBlock := GetPreciseMapWithKey(blockMap, models.BlockTypeTweet)
	mailBlockAsMap := MakeMapFromInterface(mainBlock)
	content += "[Tweet](" + mailBlockAsMap["url"] + ")"
	return content
}

func AttachmentExtractor(blockMap []map[string]interface{}) string {
	var content string
	mainBlock := GetPreciseMapWithKey(blockMap, models.BlockTypeAttachment)
	mailBlockAsMap := MakeMapFromInterface(mainBlock)
	attachmentUrl := strings.ReplaceAll(mailBlockAsMap["url"], " ", "+")
	content += "[Attachment](" + attachmentUrl + ")"
	return content
}
func CalloutExtractor(blockMap []map[string]interface{}) string {
	text := ""
	subText := ""
	for _, children := range blockMap {
		if children["type"] == nil && children["text"].(string) == "" {
			continue
		}
		if children["type"] == nil {
			subText += ChildTextExtract(children)
			continue
		}
		childrenType := children["type"]
		if childrenType == models.BlockTypeImage {
			text += ImageExtractor(blockMap)
		} else if childrenType == models.BlockTypeAttachment {
			text += AttachmentExtractor(blockMap)
		} else if childrenType == "userMention" {
			subText += UserMentionTextExtract(children)
		} else if childrenType == "pageMention" {
			subText += PageMentionTextExtract(children)
		} else if childrenType == "hr" {
			subText += Markdown_HR
		} else if childrenType != models.BlockTypeText {
			text += EmbedExtractor(children)
		} else {
			subText += ChildTextExtract(children)
		}
	}
	text += subText
	return Markdown_BLOCKQUOTE + text
}

func ChecklistExtractor(blockMap []map[string]interface{}, attr Attr) string {
	var content string
	var CheckListStatus string
	if attr.Checked {
		CheckListStatus = Markdown_ChecklistTrue
	} else {
		CheckListStatus = Markdown_ChecklistFalse
	}
	indentation := LevelMaker(attr.Level)
	content = indentation + CheckListStatus + SimpleTextExtract(blockMap)
	return content
}

func ULChecklistExtractor(blockMap []map[string]interface{}, attr Attr) string {
	var content string
	indentation := LevelMaker(attr.Level)
	content = indentation + Markdown_LI + SimpleTextExtract(blockMap)
	//fmt.Println(attr.Level, content)
	return content
}

func OLChecklistExtractor(blockMap []map[string]interface{}, attr Attr) string {
	var content string
	indentation := LevelMaker(attr.Level)
	content = indentation + Markdown_OL + SimpleTextExtract(blockMap)
	return content
}

func ImageExtractor(blockMap []map[string]interface{}) string {
	var content string
	mainBlock := GetPreciseMapWithKey(blockMap, models.BlockTypeImage)
	mailBlockAsMap := MakeMapFromInterface(mainBlock)
	imageUrl := strings.ReplaceAll(mailBlockAsMap["url"], " ", "+")
	content += "![Image](" + imageUrl + ")"
	return content
}

func CodeExtractor(blockMap []map[string]interface{}, attr Attr) string {
	var content string
	//b1, _ := json.MarshalIndent(blockMap, "", "  ")
	//fmt.Println(string(b1))
	//fmt.Println(attr)
	content = fmt.Sprintf("```%s\n%s\n```", attr.CodeLanguage, SimpleTextExtract(blockMap))
	// content = content + "\n"
	// content = content + "\n"
	// content = content + "\n"
	// content = content + "```" + attr.CodeLanguage
	// content = content + "\n"
	// content = "\t\t" + SimpleTextExtract(blockMap)
	// content = content + "```"
	// content = content + "\n"
	// content = content + "\n"
	return content
}

func GetMarkdownTextWithLinks(text string) {
	re := regexp.MustCompile(`((http|https)\:\/\/)?[a-zA-Z0-9\.\/\?\:@\-_=#]+\.([a-zA-Z0-9\&\.\/\?\:@\-_=#])*`)
	urlMatches := re.FindAllStringIndex(text, -1)
	for _, urlIndex := range urlMatches {
		fmt.Println(urlIndex)
		url := text[urlIndex[0]:urlIndex[1]]
		markdownUrlText := fmt.Sprintf("[%s](%s)", url, url)
		fmt.Println(markdownUrlText)
		// text[urlIndex[0]:urlIndex[1]] = markdownUrlText
	}
	index := 0
	for {
		re := regexp.MustCompile(`((http|https)\:\/\/)?[a-zA-Z0-9\.\/\?\:@\-_=#]+\.([a-zA-Z0-9\&\.\/\?\:@\-_=#])*`)
		urlMatches := re.FindAllStringIndex(text, -1)
		if len(urlMatches) == 0 {
			break
		}
		urlIndex := urlMatches[index]
		url := text[urlIndex[0]:urlIndex[1]]
		markdownUrlText := fmt.Sprintf("[%s](%s)", url, url)
		text = strings.Replace(text, url, markdownUrlText, -1)
	}
}

// This Block is Responsble for All the Fun in Lif
func ExtractText(blockMap []map[string]interface{}, MainBlockType string, ActualBlockType string, attributes Attr) string {
	var final string
	// Find out Actual Block Type
	fmt.Println("Main Block Type: ", MainBlockType)
	fmt.Println("Actual Block Type: ", ActualBlockType)
	fmt.Println(attributes)
	if MainBlockType == models.BlockTypeText && ActualBlockType == AssumedTextBlockType {
		final = final + SimpleTextExtract(blockMap)
	}
	if MainBlockType == models.BlockTypeText && ActualBlockType == "userMention" {
		// final = final + UserMentionTextExtract(blockMap)
	}
	if MainBlockType == models.BlockTypeText && ActualBlockType == models.BlockTypeYoutube {
		final = final + YoutubeExtractor(blockMap)
	}
	if MainBlockType == models.BlockTypeText && ActualBlockType == models.BlockTypeGoogleSheet {
		final = final + GoogleSheetExtractor(blockMap)
	}
	if MainBlockType == models.BlockTypeText && ActualBlockType == models.BlockTypeGoogleDrive {
		final = final + GoogleDriveExtractor(blockMap)
	}
	if MainBlockType == models.BlockTypeText && ActualBlockType == models.BlockTypeLoom {
		final = final + LoomExtractor(blockMap)
	}
	if MainBlockType == models.BlockTypeText && ActualBlockType == models.BlockTypeFigma {
		final = final + FigmaExtractor(blockMap)
	}
	if MainBlockType == models.BlockTypeText && ActualBlockType == models.BlockTypeReplit {
		final = final + ReplitExtractor(blockMap)
	}
	if MainBlockType == models.BlockTypeText && ActualBlockType == models.BlockTypeCodeSandBox {
		final = final + CodesandboxExtractor(blockMap)
	}
	if MainBlockType == models.BlockTypeText && ActualBlockType == models.BlockTypeTweet {
		final = final + TweetExtractor(blockMap)
	}
	if MainBlockType == models.BlockTypeText && ActualBlockType == models.BlockTypeAttachment {
		final = final + AttachmentExtractor(blockMap)
	}
	if MainBlockType == models.BlockTypeCallout && ActualBlockType == AssumedTextBlockType {
		final = final + CalloutExtractor(blockMap)
	}
	if MainBlockType == models.BlockTypeCheckList && ActualBlockType == AssumedTextBlockType {
		final = final + ChecklistExtractor(blockMap, attributes)
	}
	if MainBlockType == models.BlockTypeUList && ActualBlockType == AssumedTextBlockType {
		// final = final + ULChecklistExtractor(blockMap, attributes)
	}
	if MainBlockType == models.BlockTypeOList && ActualBlockType == AssumedTextBlockType {
		// final = final + OLChecklistExtractor(blockMap, attributes)
	}
	if MainBlockType == models.BlockTypeHeading1 && ActualBlockType == AssumedTextBlockType {
		final = final + H1Extractor(blockMap)
	}
	if MainBlockType == models.BlockTypeHeading2 && ActualBlockType == AssumedTextBlockType {
		final = final + H2Extractor(blockMap)
	}
	if MainBlockType == models.BlockTypeHeading3 && ActualBlockType == AssumedTextBlockType {
		final = final + H3Extractor(blockMap)
	}
	if MainBlockType == models.BlockTypeHeading4 && ActualBlockType == AssumedTextBlockType {
		final = final + H4Extractor(blockMap)
	}
	if MainBlockType == models.BlockTypeHeading5 && ActualBlockType == AssumedTextBlockType {
		final = final + H5Extractor(blockMap)
	}
	if MainBlockType == models.BlockTypeHeading6 && ActualBlockType == AssumedTextBlockType {
		final = final + H6Extractor(blockMap)
	}
	if MainBlockType == models.BlockTypeDivider && ActualBlockType == AssumedTextBlockType {
		final = final + Markdown_Linebreak_Prefix
		final = final + Markdown_HR
		final = final + Markdown_Linebreak_Prefix
	}
	if MainBlockType == models.BlockTypeText && ActualBlockType == models.BlockTypeImage {
		final = final + ImageExtractor(blockMap)
	}
	if MainBlockType == models.BlockTypeCode && ActualBlockType == AssumedTextBlockType {
		final = final + CodeExtractor(blockMap, attributes)
	}
	return final
}

func TocExtract(branchId uint64) string {
	content := "Sub-canvases list\n\n"
	branch, _ := queries.App.BranchQuery.GetBranchByID(branchId)
	canvasRepos, _ := App.Repo.GetCanvasRepos(branch.CanvasRepositoryID)
	for _, repo := range canvasRepos {
		canvasUrlTitle := notifications.App.Service.GenerateCanvasUrlTitle(repo.Name, *repo.DefaultBranchID)
		content += fmt.Sprintf("[%s](%s)\n", repo.Name, fmt.Sprintf("%s/%s/%s/%d", configs.GetAppInfoConfig().FrontendHost, branch.CanvasRepository.Studio.Handle, url.QueryEscape(canvasUrlTitle), repo.DefaultBranchID))
	}
	return content
}

func TableExtractor(blockMap []map[string]interface{}, canvasBranchID uint64) string {
	text := ""
	topRowDivider := ""
	for tableRowIndex, row := range blockMap {
		tableCellsString, _ := json.Marshal(row["children"])
		var tableCells []map[string]interface{}
		json.Unmarshal(tableCellsString, &tableCells)
		cellText := "| "
		for _, cell := range tableCells {
			var children []map[string]interface{}
			cellString, _ := json.Marshal(cell["children"])
			json.Unmarshal(cellString, &children)
			extractedText := ExtractFromChildren(children, canvasBranchID)
			fmt.Println("extracted cell TExt", extractedText)
			cellText += extractedText
			if tableRowIndex == 0 {
				topRowDivider += MarkdownTableTopRowDivider
			}
			cellText += " |"
		}
		text += cellText + "\n"
		if tableRowIndex == 0 {
			text += "|" + topRowDivider + "\n"
		}
	}
	text += "\n"
	return text
}

func ExtractFromChildren(cellChildren []map[string]interface{}, canvasBranchID uint64) string {
	text := ""
	for _, child := range cellChildren {
		var attributes Attr
		childType := child["type"].(string)
		if child["attributes"] != nil {
			attrString, _ := json.Marshal(child["attributes"])
			json.Unmarshal(attrString, &attributes)
		}
		var subCellChildren []map[string]interface{}
		subCellChildrenStr, _ := json.Marshal(child["children"])
		json.Unmarshal(subCellChildrenStr, &subCellChildren)

		if childType == models.BlockTypeOList {
			text += OLChecklistExtractor(subCellChildren, attributes)
		} else if childType == models.BlockTypeUList {
			text += ULChecklistExtractor(subCellChildren, attributes)
		} else if childType == models.BlockTypeCheckList {
			text += ChecklistExtractor(subCellChildren, attributes)
		} else {
			subText := ""
			for _, children := range subCellChildren {
				if children["type"] == nil && children["text"].(string) == "" {
					continue
				}
				if children["type"] == nil {
					subText += ChildTextExtract(children)
					continue
				}
				childrenType := children["type"]
				if childrenType == models.BlockTypeImage {
					text += ImageExtractor(subCellChildren)
				} else if childrenType == models.BlockTypeAttachment {
					text += AttachmentExtractor(subCellChildren)
				} else if childrenType == "userMention" {
					subText += UserMentionTextExtract(children)
				} else if childrenType == "pageMention" {
					subText += PageMentionTextExtract(children)
				} else if childrenType == "hr" {
					subText += Markdown_HR
				} else if childrenType == models.BlockDrawIO {
					continue
				} else if childrenType == "toc" {
					subText += TocExtract(canvasBranchID)
				} else if childrenType != models.BlockTypeText {
					text += EmbedExtractor(children)
				} else {
					subText += ChildTextExtract(children)
				}
			}
			text += subText
		}
	}
	return text
}

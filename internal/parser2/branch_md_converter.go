package parser2

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"

	"github.com/google/uuid"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/s3"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

type BlockChildren struct {
	Text string `json:"text"`
	Type string
	Url  string
}

func cleanString(s string) string {
	// Remove spaces.
	s = regexp.MustCompile("\\s+").ReplaceAllString(s, "")
	// Remove emojis.
	s = regexp.MustCompile("[^\\w\\s]").ReplaceAllString(s, "")
	return s
}

func (s parser2Service) StartExportingStudioNoEmail(studioId uint64, handle string) string {

	//zipFileName := fmt.Sprintf("export-%d-%s.zip", studioId, handle)
	cleanedHandle := cleanString(handle)
	zipFileName := fmt.Sprintf("backup-%d-%s.zip", studioId, cleanedHandle)
	fmt.Println("creating zip archive...")
	archive, err := os.Create(zipFileName)
	if err != nil {
		panic(err)
	}
	defer archive.Close()
	zipWriter := zip.NewWriter(archive)
	collectionNames := map[string]bool{}
	var collections []models.Collection
	App.Repo.db.Model(models.Collection{}).Where("studio_id = ? AND is_archived = false", studioId).Find(&collections)
	for _, collection := range collections {
		collectionName := fmt.Sprintf("%s-%s", collection.Name, uuid.New().String())
		var canvases []models.CanvasRepository
		err := App.Repo.db.Model(&models.CanvasRepository{}).Where("studio_id = ? AND collection_id = ? AND parent_canvas_repository_id is null  AND is_archived = false", studioId, collection.ID).Preload("Collection").Preload("ParentCanvasRepository").Find(&canvases).Error
		if err != nil {
			fmt.Println("Error in getting canvases", err)
			continue
		}

		for _, canvas := range canvases {
			canvasName := fmt.Sprintf("%s-%s", canvas.Name, uuid.New().String())
			if canvas.DefaultBranchID != nil {
				os.MkdirAll(collectionName, os.ModePerm)
				s.createMDFileAndSaveToZip(zipWriter, canvas, fmt.Sprintf("%s/%s", collectionName, canvasName))
			}
			if !collectionNames[canvas.Collection.Name] {
				collectionNames[canvas.Collection.Name] = true
			}
		}

		e := os.RemoveAll(collectionName)
		if e != nil {
			fmt.Println("Error in removing file", e, collectionName)
		}
	}

	// last steps to upload zip to s3 and send it via email
	zipWriter.Close()
	zipReader, _ := os.Open(zipFileName)
	zipUrl, _ := s3.UploadObjectToBucket(fmt.Sprintf("studio-exports/%d/%s", studioId, zipFileName), zipReader, true)
	// mailer := pkg.BipMailer{}
	// template := "<div>Workspace Export has been successfully completed.<br>" + fmt.Sprintf("Click <a href=%s>here</a> to download it.</div>", zipUrl)
	// err = mailer.SendEmail([]string{email}, nil, nil, "Bip Workspace Export Completed", template, template)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	fmt.Println(zipReader)
	fmt.Println(zipUrl)

	return zipUrl
}

func (s parser2Service) StartExportingStudio(studioId uint64, email string) string {

	zipFileName := fmt.Sprintf("export-%d-%s.zip", studioId, uuid.New().String())

	fmt.Println("creating zip archive...")
	archive, err := os.Create(zipFileName)
	if err != nil {
		panic(err)
	}
	defer archive.Close()
	zipWriter := zip.NewWriter(archive)
	collectionNames := map[string]bool{}
	var collections []models.Collection
	App.Repo.db.Model(models.Collection{}).Where("studio_id = ? AND is_archived = false", studioId).Find(&collections)
	for _, collection := range collections {
		collectionName := fmt.Sprintf("%s-%s", collection.Name, uuid.New().String())
		var canvases []models.CanvasRepository
		err := App.Repo.db.Model(&models.CanvasRepository{}).Where("studio_id = ? AND collection_id = ? AND parent_canvas_repository_id is null  AND is_archived = false", studioId, collection.ID).Preload("Collection").Preload("ParentCanvasRepository").Find(&canvases).Error
		if err != nil {
			fmt.Println("Error in getting canvases", err)
			continue
		}

		for _, canvas := range canvases {
			canvasName := fmt.Sprintf("%s-%s", canvas.Name, uuid.New().String())
			if canvas.DefaultBranchID != nil {
				os.MkdirAll(collectionName, os.ModePerm)
				s.createMDFileAndSaveToZip(zipWriter, canvas, fmt.Sprintf("%s/%s", collectionName, canvasName))
			}
			if !collectionNames[canvas.Collection.Name] {
				collectionNames[canvas.Collection.Name] = true
			}
		}

		e := os.RemoveAll(collectionName)
		if e != nil {
			fmt.Println("Error in removing file", e, collectionName)
		}
	}

	// last steps to upload zip to s3 and send it via email
	zipWriter.Close()
	zipReader, _ := os.Open(zipFileName)
	zipUrl, _ := s3.UploadObjectToBucket(fmt.Sprintf("exports/%s", zipFileName), zipReader, true)
	mailer := pkg.BipMailer{}
	template := "<div>Workspace Export has been successfully completed.<br>" + fmt.Sprintf("Click <a href=%s>here</a> to download it.</div>", zipUrl)
	err = mailer.SendEmail([]string{email}, nil, nil, "Bip Workspace Export Completed", template, template)
	if err != nil {
		fmt.Println(err)
	}
	return zipFileName
}

func (s parser2Service) createMDFileAndSaveToZip(zipWriter *zip.Writer, canvas models.CanvasRepository, filePath string) {
	fmt.Println("canvas name ========> ", canvas.Name, filePath)
	blocks, err := App.Repo.Get(*canvas.DefaultBranchID)
	if err != nil {
		fmt.Println("Error in getting branch blocks", err)
	}
	mdTextArray := s.GetMdTextFromBipBlocks(blocks)

	newFilePath := filePath + ".md"
	p := Parser2{}
	p.SaveToFile(newFilePath, mdTextArray)

	file, err := os.Open(newFilePath)
	if err != nil {
		fmt.Println("error in opening new file", err)
	}
	defer file.Close()

	w1, err := zipWriter.Create(newFilePath)
	if err != nil {
		fmt.Println("error in create new file path", err)
	}
	if _, err := io.Copy(w1, file); err != nil {
		fmt.Println("error in copying file to zip", err)
	}

	var subCanvases []models.CanvasRepository
	err = App.Repo.db.Model(&models.CanvasRepository{}).Where("parent_canvas_repository_id = ?  AND is_archived = false", canvas.ID).Preload("Collection").Preload("ParentCanvasRepository").Find(&subCanvases).Error
	if err != nil {
		fmt.Println("Error in flow", err)
	}
	if len(subCanvases) > 0 {
		os.MkdirAll(filePath+"/"+canvas.Name, os.ModePerm)
	}
	for _, subCanvas := range subCanvases {
		subCanvasName := fmt.Sprintf("%s-%s", subCanvas.Name, uuid.New().String())
		subFileName := fmt.Sprintf("%s/%s", filePath, subCanvasName)
		fmt.Println("Created subfolder here")
		s.createMDFileAndSaveToZip(zipWriter, subCanvas, subFileName)
	}
}

func (s parser2Service) StartMdBlockConversionForBranch(branchID uint64) {
	blocks, err := App.Repo.Get(branchID)
	if err != nil {
		fmt.Println("Error in getting branch blocks", err)
	}
	mdTextArray := s.GetMdTextFromBipBlocks(blocks)
	// for _, i := range mdTextArray {
	// 	fmt.Println(i)
	// }

	p := Parser2{}
	branch, _ := queries.App.BranchQuery.GetBranchByID(branchID)
	p.SaveToFile(fmt.Sprintf("%s-%d.md", url.QueryEscape(branch.CanvasRepository.Name), branchID), mdTextArray)
}

func (s parser2Service) GetMdTextFromBipBlocks(blocks []models.Block) []string {
	var content []string
	for _, block := range blocks {
		var blockText string
		blockText += s.ConvertBlockToMdString(block)
		blockText += "\n"
		// fmt.Println(strings.TrimSpace(blockText))
		content = append(content, blockText)
	}
	return content
}

func (s parser2Service) ConvertBlockToMdString(block models.Block) string {
	var blockChildren []map[string]interface{}
	json.Unmarshal(block.Children, &blockChildren)
	var attributes Attr
	json.Unmarshal(block.Attributes, &attributes)
	text := ""
	if block.Type == models.BlockTypeHeading1 {
		text += H1Extractor(blockChildren)
	} else if block.Type == models.BlockTypeHeading2 {
		text += H2Extractor(blockChildren)
	} else if block.Type == models.BlockTypeHeading3 {
		text += H3Extractor(blockChildren)
	} else if block.Type == models.BlockTypeHeading4 {
		text += H4Extractor(blockChildren)
	} else if block.Type == models.BlockTypeHeading5 {
		text += H5Extractor(blockChildren)
	} else if block.Type == models.BlockTypeHeading6 {
		text += H6Extractor(blockChildren)
	} else if block.Type == models.BlockTypeSubtext {
		text += SubTextExtractor(blockChildren)
	} else if block.Type == models.BlockQuote {
		text += QuoteExtractor(blockChildren)
	} else if block.Type == models.BlockTypeOList {
		text += OLChecklistExtractor(blockChildren, attributes)
	} else if block.Type == models.BlockTypeUList {
		text += ULChecklistExtractor(blockChildren, attributes)
	} else if block.Type == models.BlockTypeCheckList {
		text += ChecklistExtractor(blockChildren, attributes)
	} else if block.Type == models.BlockTypeCode {
		text += CodeExtractor(blockChildren, attributes)
	} else if block.Type == models.BlockTypeCallout {
		text += CalloutExtractor(blockChildren)
	} else if block.Type == models.BlockSimpleTableV1 {
		text += TableExtractor(blockChildren, *block.CanvasBranchID)
	} else {
		subText := ""
		for _, children := range blockChildren {
			if children["type"] == nil && children["text"].(string) == "" {
				continue
			}
			if children["type"] == nil {
				subText += ChildTextExtract(children)
				continue
			}
			childrenType := children["type"]
			if childrenType == models.BlockTypeImage {
				text += ImageExtractor(blockChildren)
			} else if childrenType == models.BlockTypeAttachment {
				text += AttachmentExtractor(blockChildren)
			} else if childrenType == "userMention" {
				subText += UserMentionTextExtract(children)
			} else if childrenType == "pageMention" {
				subText += PageMentionTextExtract(children)
			} else if childrenType == "hr" {
				subText += Markdown_HR
			} else if childrenType == models.BlockDrawIO {
				continue
			} else if childrenType == "toc" {
				subText += TocExtract(*block.CanvasBranchID)
			} else if childrenType != models.BlockTypeText {
				text += EmbedExtractor(children)
			} else {
				subText += ChildTextExtract(children)
			}
		}
		text += subText
	}
	return text
}

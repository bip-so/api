package parser2

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/phonepost/bip-be-platform/internal/canvasbranch"
	"gitlab.com/phonepost/bip-be-platform/internal/canvasrepo"
	"gitlab.com/phonepost/bip-be-platform/internal/collection"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/s3"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"gorm.io/datatypes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
)

func (s parser2Service) NotionImportZipHandler(file multipart.File, fileSize int64, authUser *models.User, studioID uint64) {
	fmt.Println("Started task from task file")
	userID := authUser.ID
	archive, err := zip.NewReader(file, fileSize)
	if err != nil {
		panic(err)
	}
	studioInstance, err := queries.App.StudioQueries.GetStudioQuery(map[string]interface{}{
		"id": studioID,
	})
	for _, f := range archive.File {
		filePath := f.Name
		if filePath[len(filePath)-3:] != ".md" && filePath[len(filePath)-5:] != ".html" {
			fmt.Println(f.Name)
			folderNameArray := strings.Split(f.Name, "/")
			//fileName := folderNameArray[len(folderNameArray)-1]
			fileReader, _ := f.Open()
			byts, err := ioutil.ReadAll(fileReader)
			if err != nil {
				fmt.Println(byts)
			}

			object := bytes.NewReader(byts)
			folderNameArray[len(folderNameArray)-1] = strings.ReplaceAll(folderNameArray[len(folderNameArray)-1], " ", "-")
			var re = regexp.MustCompile(`[^a-zA-Z0-9. ]`)
			folderNameArray[len(folderNameArray)-1] = re.ReplaceAllString(folderNameArray[len(folderNameArray)-1], `-`)
			attachmentsKey := strings.Join(folderNameArray[len(folderNameArray)-2:], "/")
			s3Path := strings.ReplaceAll(attachmentsKey, "/", "-")
			s3PathArray := strings.Split(s3Path, " ")
			s3Path = s3PathArray[len(s3PathArray)-1]
			response, err := s3.UploadObjectToBucket(fmt.Sprintf("%s/%s", "import", s3Path), object, true)
			if err != nil {
				fmt.Println("error in uploading object", err)
			}
			fmt.Println("URL uploaded", response)
		}
	}
	canvasData := map[string]map[string]string{}
	var collectionInstance *models.Collection
	var firstCanvasRepo *models.CanvasRepository
	if len(archive.File) > 0 {
		// Create a collection & collection permission and add default Administrator role to it.
		collectionsCount, _ := s.Manager.GetCount(models.COLLECTION, map[string]interface{}{"studio_id": studioID})
		if collectionsCount == 0 {
			collectionsCount = 1
		}
		collectionInstance, err = collection.App.Controller.CreateCollectionController(&collection.CollectionCreateValidator{
			Name:         "Notion Import",
			PublicAccess: models.PRIVATE,
			Position:     uint(collectionsCount + 1),
		}, userID, studioID)
		if err != nil {
			fmt.Println("Error in creating collection", err)
			return
		}
		fmt.Println("Created collection id", collectionInstance.ID)
	}
	canvasBranchIdMap := map[string]uint64{}
	for i, file := range archive.File {
		fmt.Println("Canvabranch id map", canvasBranchIdMap)
		filePath := file.Name
		if filePath[len(filePath)-3:] != ".md" && filePath[len(filePath)-5:] != ".html" && filePath[len(filePath)-4:] != ".csv" {
			continue
		}
		folderNameArray := strings.Split(filePath, "/")
		// This is a hack.
		// This condition is added to check the csv is an attachment, or it is a datatable.
		if filePath[len(filePath)-4:] == ".csv" {
			fileNameArray := strings.Split(folderNameArray[len(folderNameArray)-1], " ")
			if len(fileNameArray) == 1 {
				continue
			}
		}
		fileNameArray := strings.Split(folderNameArray[len(folderNameArray)-1], " ")
		canvasName := strings.Join(fileNameArray[:len(fileNameArray)-1], " ")
		canvasUniqueKey := fileNameArray[len(fileNameArray)-1]
		canvasNameWithID := strings.Join(fileNameArray, " ")
		//canvasID := fileNameArray[len(fileNameArray)-1]
		var canvasUniqueName string
		if canvasNameWithID[len(canvasNameWithID)-3:] == ".md" {
			//canvasKey = canvasID[:len(canvasID)-3]
			canvasUniqueName = canvasNameWithID[:len(canvasNameWithID)-3]
		} else if canvasNameWithID[len(canvasNameWithID)-5:] == ".html" {
			//canvasKey = canvasID[:len(canvasID)-5]
			canvasUniqueName = canvasNameWithID[:len(canvasNameWithID)-5]
		} else if canvasNameWithID[len(canvasNameWithID)-4:] == ".csv" {
			//canvasKey = canvasID[:len(canvasID)-5]
			canvasUniqueName = canvasNameWithID[:len(canvasNameWithID)-4]
		}
		// if file is csv file. and complete file name is present as a folder in that level then we can create a canvas.
		// We get blocks with UUID here
		blocks, err := s.GetBlocksData(file, canvasUniqueName)
		if err != nil {
			fmt.Println("Error in fetching blocks", canvasName, err, "So skipping this branch")
			continue
		}
		fmt.Println("blocks length", len(blocks))
		// create canvas, canvasBranch, add administrator role and add authUser as creator of the canvas.
		canvasKey := utils.NewNanoid()
		var parentCanvas map[string]string
		if len(folderNameArray) > 1 {
			for i, _ := range folderNameArray {
				tempFileName := folderNameArray[len(folderNameArray)-i-1]
				parentCanvas = canvasData[tempFileName]
				if parentCanvas != nil {
					break
				}
			}
		}

		canvasData[canvasUniqueName] = map[string]string{
			"canvasName": canvasName,
			"canvasKey":  canvasKey,
		}
		canvasCount := int64(0)
		var parentCanvasRepo *models.CanvasRepository
		if parentCanvas != nil {
			parentCanvasRepo, err = App.Repo.GetCanvasRepo(map[string]interface{}{"key": parentCanvas["canvasKey"]})
			if err != nil {
				fmt.Println("parent canvas not found for this key", parentCanvas["canvasKey"], parentCanvas["canvasName"])
				continue
			}
			canvasCount, _ = s.Manager.GetCount(models.CANVAS_REPO, map[string]interface{}{"parent_canvas_repository_id": parentCanvasRepo.ID})
		} else {
			canvasCount, _ = s.Manager.GetCount(models.CANVAS_REPO, map[string]interface{}{"collection_id": collectionInstance.ID})
		}
		canvasRepoBody := canvasrepo.NewCanvasRepoPost{
			CollectionID: collectionInstance.ID,
			Name:         canvasName,
			Position:     uint(canvasCount + 1),
		}
		if parentCanvasRepo != nil {
			canvasRepoBody.ParentCanvasRepositoryID = parentCanvasRepo.ID
		}
		canvasRepo, err := canvasrepo.App.Controller.CreateCanvasRepo(canvasRepoBody, userID, studioID, models.User{}, collectionInstance.PublicAccess)
		if err != nil {
			fmt.Println("error in creating canvas", canvasName, canvasKey)
		}
		canvasrepo.App.Repo.Update(map[string]interface{}{"id": canvasRepo.ID}, map[string]interface{}{"key": canvasKey, "position": uint(canvasCount + 1)})
		canvasBranchIdMap[canvasUniqueKey] = *canvasRepo.DefaultBranchID
		// Process the blocks and add to the canvas branch

		// First delete all the blocks on this branch
		App.Repo.db.Where("canvas_branch_id = ?", *canvasRepo.DefaultBranchID).Delete(&models.Block{})
		blocks = s.ProcessBlocks(blocks)
		_, err1 := canvasbranch.App.Controller.BlocksManager(*authUser, *canvasRepo.DefaultBranchID, canvasbranch.CanvasBlockPost{Blocks: blocks}, false, studioInstance)
		if err1 != nil {
			fmt.Println("Error in creating blocks", err1)
		}

		// Make canvas publish true.
		canvasbranch.App.Service.PublishCanvasBranch(*canvasRepo.DefaultBranchID, authUser, true)
		// Create blocks with the canvasBranchID created Above.
		if i == 0 {
			firstCanvasRepo = canvasRepo
		}
		fmt.Println(canvasName, "parentCanvasName======>", parentCanvas["canvasName"], "position====>", canvasCount+1)
	}
	fmt.Println("canvasBranchIdMap", canvasBranchIdMap)
	s.ProcessBranchBlocks(canvasBranchIdMap, authUser)
	// Send notification
	notifications.App.Service.PublishNewNotification(notifications.NotionImport, 0, []uint64{authUser.ID}, &studioID,
		nil, notifications.NotificationExtraData{
			CanvasRepoID:   firstCanvasRepo.ID,
			CanvasBranchID: *firstCanvasRepo.DefaultBranchID,
		}, nil, nil)
}

func (s parser2Service) GetBlocksData(file *zip.File, canvasUniqueName string) ([]models.PostBlocks, error) {
	data, _ := file.Open()
	filePath := file.Name
	var output []byte
	var blocks []models.PostBlocks
	payload := &bytes.Buffer{}
	w := multipart.NewWriter(payload)
	if filePath[len(filePath)-3:] == ".md" {
		output, _ = ioutil.ReadAll(data)
		//output = blackfriday.Run(mdDataString)
		part, err := w.CreateFormFile("file", filepath.Base(canvasUniqueName+".md"))
		if err != nil {
			fmt.Println(err)
		}
		part.Write(output)
		w.Close()
	} else if filePath[len(filePath)-5:] == ".html" {
		output, _ = ioutil.ReadAll(data)
		part, err := w.CreateFormFile("file", filepath.Base(canvasUniqueName+".html"))
		if err != nil {
			fmt.Println(err)
		}
		part.Write(output)
		w.Close()
	} else if filePath[len(filePath)-4:] == ".csv" {
		output, _ = ioutil.ReadAll(data)
		part, err := w.CreateFormFile("file", filepath.Base(canvasUniqueName+".csv"))
		if err != nil {
			fmt.Println(err)
		}
		part.Write(output)
		w.Close()
	} else {
		return nil, errors.New("file format is wrong")
	}

	req, err := http.NewRequest("POST", "http://bip-service-html-to-block.bip.so:5000/", payload)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	json.Unmarshal(body, &blocks)
	return blocks, nil
}

func (s parser2Service) ProcessBlocks(blocks []models.PostBlocks) []models.PostBlocks {
	for i, _ := range blocks {
		blocks[i].Scope = "create"
		blocks[i].Rank = int32((i + 1) * 1000)
	}
	return blocks
}

func (s parser2Service) ProcessBranchBlocks(canvasBranchIdMap map[string]uint64, user *models.User) {
	for _, canvasBranchID := range canvasBranchIdMap {
		blocks, err := queries.App.BlockQuery.GetBlocksByBranchID(canvasBranchID)
		if err != nil {
			fmt.Println("Error in getting blocks", err)
			continue
		}
		branch, _ := queries.App.BranchQuery.GetBranchByID(canvasBranchID)
		for _, block := range *blocks {
			var blockChildren []map[string]interface{}
			json.Unmarshal(block.Children, &blockChildren)
			blockMentions := []map[string]interface{}{}
			updateBlock := false
			for childIndex, child := range blockChildren {
				if child["text"] != nil && child["text"] == branch.CanvasRepository.Name {
					// delete the block
					if len(*blocks) > 1 {
						App.Repo.db.Delete(&block)
					} else {
						updateBlock = true
						child["text"] = ""
					}
					continue
				}
				if child["type"] != nil && child["type"].(string) == models.BlockTypeImportNotionMention {
					updateBlock = true
					filePath := blockChildren[childIndex]["fileName"].(string)
					delete(blockChildren[childIndex], "text")
					if !strings.Contains(filePath, "https://") && !strings.Contains(filePath, "http://") {
						// Add mention to the block.
						folderNameArray := strings.Split(filePath, "/")
						canvasName := folderNameArray[len(folderNameArray)-1]
						canvasNameArray := strings.Split(canvasName, "%20")
						canvasUniqueKey := canvasNameArray[len(canvasNameArray)-1]
						mentionBranchID := canvasBranchIdMap[canvasUniqueKey]
						fmt.Println("Canvas branch id map", canvasBranchIdMap)
						fmt.Println("canvas unique key and branchId", canvasUniqueKey, mentionBranchID)
						newBranch, _ := queries.App.BranchQuery.GetBranchByID(mentionBranchID)
						newMention := map[string]interface{}{
							"id":                    newBranch.ID,
							"key":                   newBranch.Key,
							"name":                  newBranch.Name,
							"type":                  "branch",
							"uuid":                  newBranch.UUID.String(),
							"repoID":                newBranch.CanvasRepositoryID,
							"repoKey":               newBranch.CanvasRepository.Key,
							"repoName":              strings.TrimSuffix(newBranch.CanvasRepository.Name, "\n"),
							"repoUUID":              newBranch.CanvasRepository.UUID.String(),
							"studioID":              newBranch.CanvasRepository.StudioID,
							"createdByUserID":       user.ID,
							"createdByUserFullName": user.FullName,
							"createdByUserUsername": user.Username,
						}
						blockChildren[childIndex]["type"] = "pageMention"
						blockChildren[childIndex]["mention"] = newMention
						blockChildren[childIndex]["uuid"] = uuid.New().String()
						blockChildren[childIndex]["children"] = []map[string]string{
							{
								"text": fmt.Sprintf("[Canvas](%s)", newBranch.CanvasRepository.Name),
							},
						}
						blockMentions = append(blockMentions, newMention)
					}
				} else if child["type"] != nil && child["type"].(string) == models.BlockTypeAttachment {
					updateBlock = true
					if child["override"] != nil && child["override"] == false {
						delete(blockChildren[childIndex], "override")
						continue
					}
					path := child["url"].(string)
					if strings.Contains(path, ".csv") && !strings.Contains(path, "http://") && !strings.Contains(path, "https://") {
						url := strings.ReplaceAll(path, "%20", " ")
						urlSplit := strings.Split(url, " ")
						fmt.Println("attachments Canvas branch id map", canvasBranchIdMap)
						fmt.Println("urlSplit key and branchId", urlSplit, canvasBranchIdMap[urlSplit[len(urlSplit)-1]])
						if canvasBranchIdMap[urlSplit[len(urlSplit)-1]] != 0 {
							mentionBranchID := canvasBranchIdMap[urlSplit[len(urlSplit)-1]]
							newBranch, _ := queries.App.BranchQuery.GetBranchByID(mentionBranchID)
							newMention := map[string]interface{}{
								"id":                    newBranch.ID,
								"key":                   newBranch.Key,
								"name":                  newBranch.Name,
								"type":                  "branch",
								"uuid":                  newBranch.UUID.String(),
								"repoID":                newBranch.CanvasRepositoryID,
								"repoKey":               newBranch.CanvasRepository.Key,
								"repoName":              strings.TrimSuffix(newBranch.CanvasRepository.Name, "\n"),
								"repoUUID":              newBranch.CanvasRepository.UUID.String(),
								"studioID":              newBranch.CanvasRepository.StudioID,
								"createdByUserID":       user.ID,
								"createdByUserFullName": user.FullName,
								"createdByUserUsername": user.Username,
							}
							delete(blockChildren[childIndex], "text")
							delete(blockChildren[childIndex], "attributes")
							delete(blockChildren[childIndex], "url")
							blockChildren[childIndex]["type"] = "pageMention"
							blockChildren[childIndex]["mention"] = newMention
							blockChildren[childIndex]["uuid"] = uuid.New().String()
							blockChildren[childIndex]["children"] = []map[string]string{
								{
									"text": fmt.Sprintf("[Canvas](%s)", newBranch.CanvasRepository.Name),
								},
							}
							blockMentions = append(blockMentions, newMention)
							continue
						}
						href := s.GetS3Path(path)
						blockChildren[childIndex]["url"] = href
					}
				}
			}
			if updateBlock {
				blockChildrenStr, _ := json.Marshal(blockChildren)
				blockMentionsStr, _ := json.Marshal(blockMentions)
				mentionsData := datatypes.JSON(blockMentionsStr)
				block.Mentions = &mentionsData
				block.Children = blockChildrenStr
				App.Repo.db.Save(&block)
			}
		}
	}
}

func (s parser2Service) GetS3Path(fileName string) string {
	fileName = strings.ReplaceAll(fileName, "%20", " ")
	folderNameArray := strings.Split(fileName, "/")
	folderNameArray[len(folderNameArray)-1] = strings.ReplaceAll(folderNameArray[len(folderNameArray)-1], " ", "-")
	var re = regexp.MustCompile(`[^a-zA-Z0-9. ]`)
	folderNameArray[len(folderNameArray)-1] = re.ReplaceAllString(folderNameArray[len(folderNameArray)-1], `-`)
	attachmentsKey := strings.Join(folderNameArray[len(folderNameArray)-2:], "/")
	s3Path := strings.ReplaceAll(attachmentsKey, "/", "-")
	s3PathArray := strings.Split(s3Path, " ")
	s3Path = s3PathArray[len(s3PathArray)-1]
	href := fmt.Sprintf("%s/import/%s", configs.GetAWSS3Config().CloudFrontURL, s3Path)
	return href
}

func (s parser2Service) GetAttachmentBlockChildren(href string, canvasUniqueName string) datatypes.JSON {
	canvasName := strings.Split(canvasUniqueName, " ")
	name := strings.Join(canvasName[:len(canvasName)-1], " ")
	children := []map[string]interface{}{
		{
			"text": "",
		},
		{
			"type":     "attachment",
			"override": false,
			"attributes": map[string]string{
				"fileName": name,
			},
			"children": []map[string]string{{"text": ""}},
			"url":      href,
			"uuid":     uuid.New(),
		},
		{
			"text": "",
		},
	}
	childrenStr, _ := json.Marshal(children)
	return childrenStr
}

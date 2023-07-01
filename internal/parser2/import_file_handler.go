package parser2

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/canvasbranch"
	"gitlab.com/phonepost/bip-be-platform/internal/canvasrepo"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

func (s parser2Service) NotionImportFileHandler(file multipart.FileHeader, authUser *models.User, studioID uint64) {
	blocks, err := s.GetBlocksByFile(file)
	if err != nil {
		fmt.Println("error in getting blocks by file", err)
		return
	}
	fmt.Println(blocks)
	studioInstance, _ := queries.App.StudioQueries.GetStudioQuery(map[string]interface{}{"id": studioID})
	// Get a collection
	// create a canvas
	// Add blocks
	// publish the branch
	// send a notification.
	collections, err := queries.App.CollectionQuery.GetCollections(map[string]interface{}{"studio_id": studioID, "is_archived": false})
	if len(collections) == 0 {
		fmt.Println("No collections found")
		return
	}
	canvasName := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))
	// In future if the import is for all users, then we need to get the collection by user access.
	collection := collections[len(collections)-1]
	canvasRepo, err := canvasrepo.App.Controller.CreateCanvasRepo(canvasrepo.NewCanvasRepoPost{
		CollectionID: collection.ID,
		Name:         canvasName,
		Position:     1,
	}, authUser.ID, studioID, *authUser, collection.PublicAccess)
	App.Repo.db.Where("canvas_branch_id = ?", *canvasRepo.DefaultBranchID).Delete(&models.Block{})
	blocks = s.ProcessBlocks(blocks)
	_, err1 := canvasbranch.App.Controller.BlocksManager(*authUser, *canvasRepo.DefaultBranchID, canvasbranch.CanvasBlockPost{Blocks: blocks}, false, studioInstance)
	if err1 != nil {
		fmt.Println("Error in creating blocks", err1)
	}

	// Make canvas publish true.
	canvasbranch.App.Service.PublishCanvasBranch(*canvasRepo.DefaultBranchID, authUser, true)
	// Create blocks with the canvasBranchID created Above.
	fmt.Println(canvasRepo.Name, "parentCanvasName======>", collection.Name, "position====>", canvasRepo.Position)

	// Send notification
	notifications.App.Service.PublishNewNotification(notifications.FileImport, 0, []uint64{authUser.ID}, &studioID,
		nil, notifications.NotificationExtraData{
			CanvasRepoID:   canvasRepo.ID,
			CanvasBranchID: *canvasRepo.DefaultBranchID,
		}, nil, nil)
}

func (s parser2Service) GetBlocksByFile(file multipart.FileHeader) ([]models.PostBlocks, error) {
	data, _ := file.Open()
	filePath := file.Filename
	var output []byte
	var canvasUniqueName string
	fmt.Println("filePath printing here", filePath)
	if filePath[len(filePath)-3:] == ".md" {
		output, _ = ioutil.ReadAll(data)
		//output = blackfriday.Run(mdDataString)
		canvasUniqueName = filePath[:len(filePath)-3] + ".md"
	} else if filePath[len(filePath)-5:] == ".html" {
		output, _ = ioutil.ReadAll(data)
		canvasUniqueName = filePath[:len(filePath)-5] + ".html"
	} else if filePath[len(filePath)-5:] == ".docx" {
		output, _ = ioutil.ReadAll(data)
		canvasUniqueName = filePath[:len(filePath)-5] + ".docx"
	} else {
		return nil, errors.New("file format is wrong")
	}
	payload := &bytes.Buffer{}
	w := multipart.NewWriter(payload)
	part, err := w.CreateFormFile("file", filepath.Base(canvasUniqueName))
	if err != nil {
		fmt.Println(err)
	}
	part.Write(output)
	w.Close()

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
	var blocks []models.PostBlocks
	json.Unmarshal(body, &blocks)
	return blocks, nil
}

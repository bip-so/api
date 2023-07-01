package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/russross/blackfriday/v2"
	"gitlab.com/phonepost/bip-be-platform/internal/canvasbranch"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

func main() {
	//core.InitCore(".env", ".env")
	//postgres.InitDB()
	//search.InitAlgolia()
	//stream.InitStreamClient()
	//redis.InitRedis()
	//kafka.InitKafka()
	//api.InitAllApps()
	//parser2.App.Service.NotionImportZipHandler()
	//archive, err := zip.OpenReader("Export-32312787-260f-46cf-b9fb-322b451f1900.zip")
	archive, err := zip.OpenReader("Export-2e1c67a9-f063-4275-ba0d-fb3863bbffb7.zip")
	if err != nil {
		panic(err)
	}
	defer archive.Close()
	for _, file := range archive.File {
		filePath := file.Name
		data, _ := file.Open()
		fmt.Println(filePath)
		if filePath[len(filePath)-3:] == ".md" {
			mdDataString, _ := ioutil.ReadAll(data)
			output := blackfriday.Run(mdDataString)
			fmt.Println("output here", string(output))
			payload := &bytes.Buffer{}
			w := multipart.NewWriter(payload)
			part, err := w.CreateFormFile("file", filepath.Base("boardView.html"))
			if err != nil {
				fmt.Println(err)
			}
			part.Write(output)
			w.Close()

			req, err := http.NewRequest("POST", "http://3.7.152.94:5000/", payload)
			if err != nil {
				fmt.Println(err)
			}
			req.Header.Add("Content-Type", w.FormDataContentType())

			client := &http.Client{}
			res, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
			}
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Println(err)
			}
			var blocks []canvasbranch.PostBlocks
			json.Unmarshal(body, &blocks)
			fmt.Println("blocks elnght", len(blocks))
			for _, a := range blocks {
				fmt.Println(a.UUID)
			}
			break
		} else if filePath[len(filePath)-5:] == ".html" {
			htmlDataString, _ := ioutil.ReadAll(data)
			fmt.Println(string(htmlDataString))
		}
		break
	}
	//
	//for _, f := range archive.File {
	//	filePath := f.Name
	//	if filePath[len(filePath)-3:] != ".md" && filePath[len(filePath)-5:] != ".html" {
	//		fmt.Println(f.Name)
	//		folderNameArray := strings.Split(f.Name, "/")
	//		//fileName := folderNameArray[len(folderNameArray)-1]
	//		file, _ := f.OpenRaw()
	//		attachmentsKey := strings.Join(folderNameArray[len(folderNameArray)-2:], "/")
	//		s3Path := strings.ReplaceAll(attachmentsKey, "/", "-")
	//		response, err := s3.UploadObjectToBucket(fmt.Sprintf("%s/%s", "import", s3Path), file, true)
	//		if err != nil {
	//			fmt.Println("error in uploading object", err)
	//		}
	//		fmt.Println("URL uploaded", response)
	//	}
	//}
	//canvasData := map[string]map[string]string{}
	//if len(archive.File) > 0 {
	//	// Create a collection & collection permission and add default Administrator role to it.
	//}
	//for _, file := range archive.File {
	//	filePath := file.Name
	//	if filePath[len(filePath)-3:] != ".md" && filePath[len(filePath)-5:] != ".html" {
	//		continue
	//	}
	//	folderNameArray := strings.Split(filePath, "/")
	//	var parentCanvas map[string]string
	//	if len(folderNameArray) > 1 {
	//		for i, _ := range folderNameArray {
	//			tempFileName := folderNameArray[len(folderNameArray)-i-1]
	//			parentCanvas = canvasData[tempFileName]
	//			if parentCanvas["canvasName"] != "" {
	//				break
	//			}
	//		}
	//	}
	//
	//	fileNameArray := strings.Split(folderNameArray[len(folderNameArray)-1], " ")
	//	canvasName := strings.Join(fileNameArray[:len(fileNameArray)-1], " ")
	//	canvasNameWithID := strings.Join(fileNameArray, " ")
	//	canvasID := fileNameArray[len(fileNameArray)-1]
	//	var canvasUniqueName string
	//	var canvasKey string
	//	if canvasNameWithID[len(canvasNameWithID)-3:] == ".md" {
	//		canvasKey = canvasID[:len(canvasID)-3]
	//	} else if canvasNameWithID[len(canvasNameWithID)-5:] == ".html" {
	//		canvasKey = canvasID[:len(canvasID)-5]
	//	}
	//	fmt.Println(canvasKey, canvasName)
	//	canvasData[canvasUniqueName] = map[string]string{
	//		"canvasName": canvasName,
	//		"canvasKey":  canvasKey,
	//	}
	//	//mdData, _ := file.Open()
	//	//mdDataString, _ := ioutil.ReadAll(mdData)
	//	//fmt.Println(mdDataString)
	//	// ConvertHtmlToBlocks(mdDataString)
	//	// For page mentions we need to store the blockIDs in redis and after all zip files are completed we need to update these.
	//	// We get blocks with UUID here here
	//	// create canvas, canvasBranch, add administrator role and add authUser as creator of the canvas.
	//	// Make canvas publish true.
	//	// Create blocks with the canvasBranchID created Above.
	//	//fmt.Println(canvasName, "parentCanvasName======>", parentCanvas["canvasName"])
	//}
}

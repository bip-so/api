package eventHandler

import (
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

func parseImage(remoteURL string) (imageText string, err error) {
	fmt.Println("Reading from " + remoteURL)

	remote, err := http.Get(remoteURL)
	if err != nil {
		return
	}

	defer remote.Body.Close()
	lastBin := strings.LastIndex(remoteURL, "/")
	fileName := remoteURL[lastBin+1:]

	fmt.Println("Filename is " + fileName)

	//open a file for writing
	file, err := os.Create("/tmp/" + fileName)
	if err != nil {
		return
	}

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, remote.Body)
	if err != nil {
		return
	}

	file.Close()
	fmt.Println("Image File Pulled and saved to /tmp/" + fileName)

	// //load file to read
	// buf, err := ioutil.ReadFile("/tmp/" + fileName)
	// if err != nil {
	// 	return
	// }

	// check filetype
	// if !filetype.IsImage(buf) {
	// 	fmt.Println("file is not an image\n")
	// 	return
	// }

	// fmt.Println("File is an image")

	// client := gosseract.NewClient()
	// defer client.Close()

	// client.SetImage("/tmp/" + fileName)
	// w, h := getImageDimension("/tmp/" + fileName)
	// fmt.Println("Image width is " + strconv.Itoa(h))
	// fmt.Println("Image height is " + strconv.Itoa(w))
	// imageText, err = client.Text()
	// if err != nil {
	// 	return
	// }

	// if len(imageText) >= 1 {
	// 	imageText = imageText[:len(imageText)-1]
	// }

	fmt.Println(imageText)
	fmt.Println("Image Parsed")

	return
}

func getImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("error sending message", err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Println("error sending message", err)
	}
	return image.Width, image.Height
}

// paste site handling
func parseBin(url, format string) (binText string, err error) {
	var rawURL string

	fmt.Printf("reading from %s", url)
	_, file := path.Split(url)
	rawURL = strings.Replace(format, "&filename&", file, 1)

	fmt.Println("Raw text URL is " + rawURL)

	resp, err := http.Get(rawURL)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	binText = string(body)

	fmt.Println("Contents = \n" + binText)

	return
}

// parses url contents for images and paste sites.
func parseURL(url string, parseConf Parsing) (parsedText string) {
	//Catch domains and route to the proper controllers (image, binsite parsers)
	fmt.Printf("checking for pastes and images on %s\n", url)
	// if a url ends with a / remove it. Stupid chrome adds them.
	if strings.HasSuffix(url, "/") {
		url = strings.TrimSuffix(url, "/")
	}
	if len(parseConf.Image.Sites) != 0 {
		for _, site := range parseConf.Image.Sites {
			fmt.Printf("checking paste site %s", site.URL)
			if strings.HasPrefix(url, site.URL) {
				fmt.Printf("matched on url %s", site.URL)
				_, file := path.Split(url)
				url = strings.Replace(site.Format, "&filename&", file, 1)
			}
		}
	}

	//check for image filetypes
	for _, filetype := range parseConf.Image.FileTypes {
		fmt.Println("checking if image")
		if strings.HasSuffix(url, filetype) {
			fmt.Println("found image file")
			if imageText, err := parseImage(url); err != nil {
				fmt.Printf("%s\n", err)
			} else {
				fmt.Println(imageText)
				parsedText = imageText
				return
			}
		}
	}

	// check for paste sites
	for _, paste := range parseConf.Paste.Sites {
		fmt.Println("checking if bin file")
		if strings.HasPrefix(url, paste.URL) {
			if binText, err := parseBin(url, paste.Format); err != nil {
				fmt.Printf("%s\n", err)
			} else {
				fmt.Println(binText)
				parsedText = binText
				return
			}
		}
	}

	return
}

//     __                               __
//    / /_____ __ ___    _____  _______/ /
//   /  '_/ -_) // / |/|/ / _ \/ __/ _  /
//  /_/\_\\__/\_, /|__,__/\___/_/  \_,_/
//  	     /___/

// returns response and reaction for keywords
func parseKeyword(message, botName string, channelKeywords []Keyword, parseConf Parsing) (response, reaction []string) {
	fmt.Printf("Parsing inbound chat for %s", botName)

	message = strings.ToLower(message)

	//exact match search
	fmt.Println("Testing matches")
	for _, keyWord := range channelKeywords {
		if message == keyWord.Keyword && keyWord.Exact { // if the match was an exact match
			fmt.Printf("Response is %v", keyWord.Response)
			fmt.Printf("Reaction is %v", keyWord.Reaction)
			return keyWord.Response, keyWord.Reaction
		} else if strings.Contains(message, keyWord.Keyword) && !keyWord.Exact { // if the match was just a match
			fmt.Printf("Response is %v", keyWord.Response)
			fmt.Printf("Reaction is %v", keyWord.Reaction)
			return keyWord.Response, keyWord.Reaction
		}
	}

	lastIndex := -1

	//Match on errors
	fmt.Println("Testing matches")

	for _, keyWord := range channelKeywords {
		if strings.Contains(message, keyWord.Keyword) {
			fmt.Printf("match is %s", keyWord.Keyword)
		}

		index := strings.LastIndex(message, keyWord.Keyword)
		if index > lastIndex && !keyWord.Exact {
			lastIndex = index
			response = keyWord.Response
			reaction = keyWord.Reaction
		}
	}

	return
}

//                                     __
//  _______  __ _  __ _  ___ ____  ___/ /
// / __/ _ \/  ' \/  ' \/ _ `/ _ \/ _  /
// \__/\___/_/_/_/_/_/_/\_,_/_//_/\_,_/
//

// AdminCommand commands are hard coded for now
func adminCommand(message, botName string, servCommands []Command, servKeywords []Keyword) (response, reaction []string) {
	fmt.Printf("Parsing inbound admin command for %s", botName)
	message = strings.ToLower(message)

	return
}

// ModCommand commands are hard coded for now
func modCommand(message, botName string, servCommands []Command) (response, reaction []string) {
	fmt.Printf("Parsing inbound mod command for %s", botName)
	message = strings.ToLower(message)
	return
}

// Command parses commands
func parseCommand(message, botName string, channelCommands []Command) (response, reaction []string) {
	fmt.Printf("Parsing inbound command for %s", botName)
	message = strings.ToLower(message)

	for _, command := range channelCommands {
		if command.Command == message {
			response = command.Response
			reaction = command.Reaction
		}
	}
	return
}

// general funcs
func contains(array []string, str string) bool {
	for _, value := range array {
		if value == str {
			return true
		}
	}
	return false
}

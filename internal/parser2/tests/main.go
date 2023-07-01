package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	// configs.InitConfig(".env", ".")
	// postgres.InitDB()
	// parser2.InitApp()
	// parser2.App.Service.StartMdBlockConversionForBranch(7543)
	GetMarkdownTextWithLinks("hello www.google.com I am from https://www.yahoo.com")
}

func GetMarkdownTextWithLinks(text string) {
	re := regexp.MustCompile(`((http|https)\:\/\/)?[a-zA-Z0-9\.\/\?\:@\-_=#]+\.([a-zA-Z0-9\&\.\/\?\:@\-_=#])*`)
	urlMatches := re.FindAllStringIndex(text, -1)
	for index, _ := range urlMatches {
		re := regexp.MustCompile(`((http|https)\:\/\/)?[a-zA-Z0-9\.\/\?\:@\-_=#]+\.([a-zA-Z0-9\&\.\/\?\:@\-_=#])*`)
		urlMatches := re.FindAllStringIndex(text, -1)
		urlIndex := urlMatches[index]
		url := text[urlIndex[0]:urlIndex[1]]
		markdownUrlText := fmt.Sprintf("[%s](%s)", url, url)
		text = strings.Replace(text, url, markdownUrlText, -1)
		fmt.Println(text)
	}
	fmt.Println("what happend", text)
}

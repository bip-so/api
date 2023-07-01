package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	//api := slack.New("xoxp-2976430295908-2971212875189-4119308834340-ab6d6d60ced261a1ad062bfebac1604a")

	// Add reaction test
	//err := api.AddReaction("white_check_mark", slack.ItemRef{Channel: "C02UQCN9KMJ", Timestamp: "1664349741.201169"})
	//fmt.Println(err)
	url := "https://files.slack.com/files-pri/T02UQCN8PSQ-F043T4FC571/5bbab31177f74e09acc5af3185ec0b00-download-1.png"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer xoxp-2976430295908-2971212875189-4119308834340-ab6d6d60ced261a1ad062bfebac1604a")
	req.Header.Add("Cookie", "x=c3f3984ef4d3b9453cfeaa02d371c143.1664355818")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

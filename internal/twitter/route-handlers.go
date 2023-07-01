package twitter

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"io/ioutil"
	"net/http"
)

// GetMetadata Get Twitter tweet Metadata
// @Summary 	Get Twitter tweet Metadata
// @Description
// @Tags		Twitter
// @Security 	bearerAuth
// @Param 		id 	path 	string	true "Tweet ID"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/twitter/metadata/{id} [get]
func (impl *TwitterImpl) GetMetadata(c *gin.Context) {
	tweetId := c.Param("id")
	url := "https://api.twitter.com/2/tweets/" + tweetId + "?tweet.fields=attachments,author_id,created_at," +
		"entities,geo,id,in_reply_to_user_id,lang,possibly_sensitive,referenced_tweets,source,text,withheld"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	// @todo: SR: Move to ENV
	req.Header.Add("Authorization", "Bearer AAAAAAAAAAAAAAAAAAAAABjLeQEAAAAA2HBnr8amAFKXBes8MZySJ8Spi1c%3DKGF9mc3ANXq49zIiy6NpokYMC0WHEytLQ33iz2slZ9B9MGwIEJ")

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
	var result interface{}
	json.Unmarshal(body, &result)
	response.RenderResponse(c, result)
	return
}

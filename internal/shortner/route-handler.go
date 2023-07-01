package shortner

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (r *shortRoutes) Get(c *gin.Context) {
	shortID := c.Param("shortID")
	shortInstance, errCreatURL := App.Service.Get(shortID)
	fmt.Println(shortInstance)
	if errCreatURL != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errCreatURL.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Create",
		"code":    shortInstance.ShortCode,
		"url":     shortInstance.OriginalURL,
	})
	return
}

func (r *shortRoutes) Create(c *gin.Context) {
	var body ShortPost
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	_, errGettingUser := r.GetLoggedInUser(c)
	if errGettingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errGettingUser.Error(),
		})
		return
	}
	fmt.Println(body)
	shortInstance, errCreatURL := App.Service.Create(body.OriginalURL)
	if errCreatURL != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errCreatURL.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Create",
		"url":     shortInstance.ShortCode,
	})

}

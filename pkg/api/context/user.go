package context

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

func GetAuthUser(c *gin.Context) (*models.User, error) {
	user, isOk := c.Get("currentUser")
	if !isOk || user == nil {
		return nil, errors.New("Auth User Not Found")
	}
	authUser := user.(*models.User)
	return authUser, nil
}

func GetAuthStudio(c *gin.Context) (uint64, error) {
	studio, isOk := c.Get("currentStudio")
	if !isOk || studio == nil {
		return 0, errors.New("Studio Id Not Found")
	}
	studioId := studio.(uint64)
	return studioId, nil
}

package core

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

type RouteHelper struct {
}

// GetLoggedInUser Returns User Instance
func (rh RouteHelper) GetLoggedInUser(c *gin.Context) (*models.User, error) {

	user, isOk := c.Get("currentUser")
	if !isOk || user == nil {
		return nil, errors.New("Auth User Not Found")
	}
	authUser := user.(*models.User)
	return authUser, nil
}

// GetLoggedInUserId Returns UserID as uint64
func (rh RouteHelper) GetLoggedInUserId(c *gin.Context) uint64 {
	ctxUser, _ := c.Get("currentUser")
	if ctxUser == nil {
		return 0
	}
	loggedInUser := ctxUser.(*models.User)
	if loggedInUser == nil {
		return 0
	}
	return loggedInUser.ID
}

// GetStudioId Returns Studio Id
func (rh RouteHelper) GetStudioId(c *gin.Context) (uint64, error) {

	studio, isOk := c.Get("currentStudio")
	if !isOk || studio == nil {
		return 0, errors.New("Studio Id Not Found")
	}
	studioId := studio.(uint64)
	return studioId, nil
}

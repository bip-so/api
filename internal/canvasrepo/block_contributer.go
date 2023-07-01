package canvasrepo

import (
	"encoding/json"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gorm.io/datatypes"
	"time"
)

func (s canvasRepoService) BlockContributorFirst(user models.User, branchId uint64) datatypes.JSON {
	contribution := queries.UserBlockContributor{
		Id:        user.ID,
		UUID:      user.UUID.String(),
		FullName:  user.FullName,
		Username:  user.Username,
		AvatarUrl: user.AvatarUrl,
		Timestamp: time.Now(),
		BranchID:  branchId,
	}
	j := []queries.UserBlockContributor{contribution}
	singleContrib, _ := json.Marshal(j)
	first := datatypes.JSON(singleContrib)
	return first
}
func (s canvasRepoService) BlockContributorNext(user models.User, branchId uint64) datatypes.JSON {
	contribution := queries.UserBlockContributor{
		Id:        user.ID,
		UUID:      user.UUID.String(),
		FullName:  user.FullName,
		Username:  user.Username,
		AvatarUrl: user.AvatarUrl,
		Timestamp: time.Now(),
		BranchID:  branchId,
	}
	singleContrib, _ := json.Marshal(contribution)
	first := datatypes.JSON(singleContrib)
	return first
}

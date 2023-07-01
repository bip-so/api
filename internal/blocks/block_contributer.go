package blocks

import (
	"encoding/json"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gorm.io/datatypes"
	"time"
)

func (s blockService) BlockContributorFirst(user models.User, branchId uint64) datatypes.JSON {
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
func (s blockService) BlockContributorNext(user models.User, branchId uint64) datatypes.JSON {
	// We need to check if contribution is present and just update the time
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
	next := datatypes.JSON(singleContrib)
	return next
}

func (s blockService) BuildUniqueContribs(sample []queries.UserBlockContributor) []queries.UserBlockContributor {
	var unique []queries.UserBlockContributor
uniqueLoop:
	for _, v := range sample {
		for i, u := range unique {
			if v.Id == u.Id {
				unique[i] = v
				continue uniqueLoop
			}
		}
		unique = append(unique, v)
	}
	return unique
}

// Merges the Contributiuon from two Blocks
func (s blockService) BlockContributorMerge(parent datatypes.JSON, child datatypes.JSON) datatypes.JSON {
	var p, c, pc, pc2 []queries.UserBlockContributor
	_ = json.Unmarshal(parent, &p)
	_ = json.Unmarshal(child, &c)

	pc = append(p, c...)
	pc2 = s.BuildUniqueContribs(pc)
	mergedContrib, _ := json.Marshal(pc2)
	return datatypes.JSON(mergedContrib)
}

// Merges the Reactions from two Blocks
func (s blockService) BlockReactionsMerge(parent datatypes.JSON, child datatypes.JSON) datatypes.JSON {
	var p, c, pc []models.ReactionCounter
	_ = json.Unmarshal(parent, &p)
	_ = json.Unmarshal(child, &c)
	pc = append(p, c...)
	mergedContrib, _ := json.Marshal(pc)
	return datatypes.JSON(mergedContrib)
}

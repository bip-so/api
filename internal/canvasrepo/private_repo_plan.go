package canvasrepo

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

func (s *canvasRepoRepo) StudioPlanCheck(studioID uint64) (bool, int64, int) {
	var studio models.Studio
	_ = s.db.Model(models.Studio{}).Where("id = ? and is_archived = ?", studioID, false).Preload("CreatedByUser").First(&studio).Error

	privateRepoCount := s.PrivateRepoCount(studioID)
	ENV := configs.GetConfigString("ENV")
	stripeData := utils.STRIPE_CONST_LOOKUPS[ENV]

	// will improve more.
	if studio.IsEarlyAdopter || studio.IsNonProfit || (studio.ID == studio.CreatedByUser.DefaultStudioID) {
		return false, privateRepoCount, 1000000
	} else {
		if int(privateRepoCount) > models.CanvasMaxPrivateReposAllowed && studio.StripePriceID == stripeData["LITE_PRICE"] {
			return true, privateRepoCount, models.CanvasMaxPrivateReposAllowed
		} else {
			return false, privateRepoCount, models.CanvasMaxPrivateReposAllowed
		}
	}
}

func (s *canvasRepoRepo) PrivateRepoCount(studioID uint64) int64 {
	type Result struct {
		Privates int64
	}
	query := "SELECT  ( SELECT COUNT(*) FROM canvas_repositories as CR INNER JOIN canvas_branches as CB ON CB.id = CR.default_branch_id where CR.studio_id = ? and CR.is_archived = false and cb.public_access = 'private' and CR.is_published = true) AS Privates"
	var result Result
	s.db.Raw(query, studioID).Scan(&result)
	fmt.Println("Private Repo Count", result.Privates)
	return result.Privates
}

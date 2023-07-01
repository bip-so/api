package parser2

import "gitlab.com/phonepost/bip-be-platform/internal/models"

func (s parser2Service) GetBranchBlocks(branchID uint64) *[]models.Block {
	blocks, err := App.Repo.Get(branchID)
	if err != nil {
		return nil
	}
	return &blocks
}

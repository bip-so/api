package canvasbranch

import "gitlab.com/phonepost/bip-be-platform/internal/models"

/*
This all works as the  Block we are creating is not on the MAIN BRANCH
All we are doing is moving things from ROUGH to MAIN
*/

// MoveThreadCommentsReelsReactionsMainBranchNewBlock : Will move the block threads, Reels and Reactions to the Main Branch
func (s canvasBranchService) MoveThreadCommentsReelsReactionsMainBranchNewBlock(blockID uint64, destinationBranchID *uint64) {
	App.Repo.UpdateBlockThreadCommentsBranch(blockID, destinationBranchID)
	App.Repo.UpdateReelsBranch(blockID, destinationBranchID)
	App.Repo.TransferReactionsFromNewBlock(blockID, destinationBranchID)
}

// UpdateBlockThreadCommentsBranch Find all block threads for this Block ID and Change the Branch
func (r canvasBranchRepo) UpdateBlockThreadCommentsBranch(blockID uint64, branchID *uint64) {
	r.db.Model(models.BlockThread{}).Where("start_block_id = ?", blockID).Updates(map[string]interface{}{"canvas_branch_id": &branchID})
}

// UpdateReelsBranch Move all the reels on this Rough Branch (Block)   to Main Branch (Block)
func (r canvasBranchRepo) UpdateReelsBranch(blockID uint64, branchID *uint64) {
	r.db.Model(models.Reel{}).Where("start_block_id = ?", blockID).Updates(map[string]interface{}{"canvas_branch_id": &branchID})
}

// TransferReactionsFromNewBlock We are moving  BlockReaction, BlockThreadReaction, BlockCommentReaction instances to new block
func (r canvasBranchRepo) TransferReactionsFromNewBlock(blockID uint64, branchID *uint64) {
	// We are moving the SAME BLOCK ID"S BRANCH THAT ALL
	r.db.Model(models.BlockReaction{}).Where("block_id = ?", blockID).Updates(map[string]interface{}{"canvas_branch_id": &branchID})
	r.db.Model(models.BlockThreadReaction{}).Where("block_id = ?", blockID).Updates(map[string]interface{}{"canvas_branch_id": &branchID})
	r.db.Model(models.BlockCommentReaction{}).Where("block_id = ?", blockID).Updates(map[string]interface{}{"canvas_branch_id": &branchID})
}

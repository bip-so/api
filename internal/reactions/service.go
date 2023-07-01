package reactions

import (
	"context"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

//---------------------------------- Create things ------------------------------//

func (s reactionService) CreateBlockReaction(obj CreateMentionPost, studioID uint64, userID uint64) error {
	err := App.Repo.CreateBlockReaction(obj, studioID, userID)
	if err != nil {
		return err
	}
	return nil
}
func (s reactionService) CreateBlockThreadReaction(obj CreateMentionPost, studioID uint64, userID uint64) error {
	err := App.Repo.CreateBlockThreadReaction(obj, studioID, userID)
	if err != nil {
		return err
	}
	return nil
}
func (s reactionService) CreateBlockThreadCommentReaction(obj CreateMentionPost, studioID uint64, userID uint64) error {
	err := App.Repo.CreateBlockThreadCommentReaction(obj, studioID, userID)
	if err != nil {
		return err
	}
	return nil
}
func (s reactionService) CreateReelReaction(obj CreateMentionPost, studioID uint64, userID uint64) error {
	err := App.Repo.CreateReelReaction(obj, studioID, userID)
	if err != nil {
		return err
	}
	return nil
}
func (s reactionService) CreateReelCommentReaction(obj CreateMentionPost, studioID uint64, userID uint64) error {
	err := App.Repo.CreateReelCommentReaction(obj, studioID, userID)
	if err != nil {
		return err
	}
	return nil
}

//---------------------------------- Remove things ------------------------------//

func (s reactionService) RemoveBlockReaction(obj CreateMentionPost, studioID uint64, userID uint64) error {
	err := App.Repo.RemoveBlockReaction(obj, studioID, userID)
	if err != nil {
		return err
	}
	return nil
}
func (s reactionService) RemoveBlockThreadReaction(obj CreateMentionPost, studioID uint64, userID uint64) error {
	err := App.Repo.RemoveBlockThreadReaction(obj, studioID, userID)
	if err != nil {
		return err
	}
	return nil
}
func (s reactionService) RemoveBlockThreadCommentReaction(obj CreateMentionPost, studioID uint64, userID uint64) error {
	err := App.Repo.RemoveBlockThreadCommentReaction(obj, studioID, userID)
	if err != nil {
		return err
	}
	return nil
}
func (s reactionService) RemoveReelReaction(obj CreateMentionPost, studioID uint64, userID uint64) error {
	err := App.Repo.RemoveReelReaction(obj, studioID, userID)
	if err != nil {
		return err
	}
	return nil
}
func (s reactionService) RemoveReelCommentReaction(obj CreateMentionPost, studioID uint64, userID uint64) error {
	err := App.Repo.RemoveReelCommentReaction(obj, studioID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s reactionService) InvalidateBranchBlocks(canvasBranchID uint64) {
	cachedBranchKey := models.CANVAS_BRANCH + ":" + utils.String(canvasBranchID) + ":*"
	s.cache.DeleteMatching(context.Background(), cachedBranchKey)
	cachedBranchKey = models.CANVAS_BRANCH + ":" + utils.String(canvasBranchID)
	s.cache.Delete(context.Background(), cachedBranchKey)
}

func (s *reactionService) InvalidateReelsCachingViaStudio(studioID uint64) {
	//s.cache.HDelete(context.Background(), "studio-reels:"+utils.String(studioID), "*")
	s.cache.HDeleteMatching(context.Background(), "cached-studio-reels:"+utils.String(studioID), "*")
}

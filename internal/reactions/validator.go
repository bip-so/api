package reactions

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

type CreateMentionPost struct {
	Scope          string `json:"scope" binding:"required"` // "block","block_comment","reel","reel_comment","block_thread"
	Emoji          string `json:"emoji" binding:"required"`
	CanvasBranchID uint64 `json:"canvasBranchID" binding:"required"`
	ReelID         uint64 `json:"reelID"`
	ReelCommentID  uint64 `json:"reelCommentID"`
	//BlockID        uint64 `json:"blockID"`
	BlockUUID      uuid.UUID `json:"blockUUID"`
	BlockCommentID uint64    `json:"blockCommentID"`
	BlockThreadID  uint64    `json:"blockThreadID"`
}

// Validate
func (obj CreateMentionPost) Validate() error {
	fmt.Println(obj)

	allowedScope := []string{"block", "block_thread", "block_comment", "reel", "reel_comment"}
	if !utils.SliceContainsItem(allowedScope, obj.Scope) {
		return errors.New("Please only send  \"block\",\"block_thread\", \"block_comment\", \"reel\", \"reel_comment\"")
	}

	if obj.Emoji == "" {
		return errors.New("Emoji can't be  empty.")
	}

	switch obj.Scope {
	case "block":
		if obj.BlockUUID == uuid.Nil {
			return errors.New("Block UUID Required.")
		}
		break
	case "block_thread":
		if obj.BlockUUID == uuid.Nil || obj.BlockThreadID == 0 {
			return errors.New("Block ID and Thread ID are Required.")
		}
		break
	case "block_comment":
		if obj.BlockUUID == uuid.Nil || obj.BlockThreadID == 0 || obj.BlockCommentID == 0 {
			return errors.New("Block ID Block Comment and Thread ID are Required.")
		}
		break
	case "reel":
		if obj.ReelID == 0 {
			return errors.New("Reel ID Required.")
		}
		break
	case "reel_comment":
		if obj.ReelID == 0 || obj.ReelCommentID == 0 {
			return errors.New("Reel ID and Reel Comment Required.")
		}
		break
	default:
		return errors.New("Validation Incorrect Input")
	}
	return nil
}

type ReactionCounter struct {
	Scope string `json:"scope"`
	Count string `json:"count"`
}

type ReactionCounterArray struct {
	Scope string `json:"scope"`
	Data  ReactionCounter
}

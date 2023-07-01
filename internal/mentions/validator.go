package mentions

import (
	"errors"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

type MentionPost struct {
	Users          []uint64 `json:"users"`      // List of user ID's
	Branches       []uint64 `json:"branches"`   // List of Branches ID's
	Roles          []uint64 `json:"roles"`      // List of Role ID's
	Scope          string   `json:"scope"`      // Scope
	ObjectID       uint64   `json:"objectID"`   //The id of the Scope
	ObjectUUID     string   `json:"objectUUID"` // UUID for the block
	CanvasBranchID uint64   `json:"canvasBranchId"`
}

func (obj MentionPost) Validate() error {
	fmt.Println(obj)
	allowedScope := []string{"block", "block_thread", "block_comment", "reel", "reel_comment"}
	if !utils.SliceContainsItem(allowedScope, obj.Scope) {
		return errors.New("Please only send  \"block\",\"block_thread\", \"block_comment\", \"reel\", \"reel_comment\"")
	}
	return nil
}

package canvasbranch

import (
	"errors"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"gorm.io/datatypes"
)

type CanvasBranchVisibilityPost struct {
	Visibility string `json:"visibility" binding:"required"` //"private" "view" "comment"  "edit"
}
type CanvasBranchPublishPost struct {
	Status bool `json:"status" binding:"required"`
}

type initPost struct {
	CollectionID uint64 `json:"collectionID" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Icon         string `json:"icon"`
}

// post
type newCanvasCreatePost struct {
	CollectionID uint64 `json:"collectionID" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Icon         string `json:"icon"`
	BranchId     string `json:"branchId"` // this is needed only if
}

// newCanvasBranchPost
type newCanvasBranchPost struct {
	CanvasRepoID       uint64 `json:"canvasRepoId" binding:"required"`
	CollectionID       uint64 `json:"collectionId" binding:"required"`
	ParentCanvasRepoID uint64 `json:"parentCanvasRepoId"`
	FromCanvasBranchID uint64 `json:"fromCanvasBranchId"`
}

// post for the blocks
type NewDraftBranchPost struct {
	CollectionID       uint64 `json:"collectionId"`
	CanvasRepoID       uint64 `json:"canvasRepoId"`
	ParentCanvasRepoID uint64 `json:"parentCanvasRepoId"`
}
type EmptyPost struct {
}

//type InitMergePost struct {
//	MergeRequestID  *uint64                 `json:"mergeRequestId"`
//	MergeStatus     string                  `json:"status"`
//	ChangesAccepted *map[string]interface{} `json:"changesAccepted"`
//	CommitMessage   string                  `json:"commitMessage" binding:"required"`
//}
type MergeRequestAcceptPartialPost struct {
	//MergeRequestID  uint64                  `json:"mergeRequestID"`
	MergeStatus     string                  `json:"status"`
	ChangesAccepted *map[string]interface{} `json:"changesAccepted"`
	CommitMessage   string                  `json:"commitMessage" binding:"required"`
}

// type MergeRequestMergeRequestRejectedPost struct {
// 	MergeRequestID uint64 `json:"mergeRequestID"`
// 	MergeStatus    string `json:"status"`
// }
type MergeRequestCreatePost struct {
	CommitMessage string `json:"commitMessage"`
}

type ManagePublishRequest struct {
	Accept bool `json:"accept"`
}

type InitPRPost struct {
	Message string `json:"message" binding:"required"`
}

type MergeMergePost struct {
	CommitMessage   string         `json:"commitMessage" binding:"required"`
	Status          string         `json:"status" binding:"required"`
	ChangesAccepted datatypes.JSON `json:"changesAccepted"`
}

//MERGE_REQUEST_OPEN               = "OPEN"
//MERGE_REQUEST_ACCEPTED           = "ACCEPTED"
//MERGE_REQUEST_REJECTED           = "REJECTED"
//MERGE_REQUEST_PARTIALLY_ACCEPTED = "PARTIALLY_ACCEPTED"
func (obj MergeMergePost) Validate() error {
	fmt.Println(obj)

	allowedStatus := []string{"OPEN", "ACCEPTED", "REJECTED", "PARTIALLY_ACCEPTED"}
	if !utils.SliceContainsItem(allowedStatus, obj.Status) {
		return errors.New("Please only send \"OPEN\", \"ACCEPTED\", \"REJECTED\", \"PARTIALLY_ACCEPTED\"")
	}
	// Add validation for
	if obj.Status == "PARTIALLY_ACCEPTED" {
		// This means we can't have
		if len(obj.ChangesAccepted.String()) == 0 {
			return errors.New("we have selected `PARTIALLY_ACCEPTED` ChangesAccepted is Required")
		}
	}

	return nil
}

type PlaceHolder struct {
}

// post for the blocks
type CanvasBlockPost struct {
	Blocks []models.PostBlocks `json:"blocks"`
}

type CanvasBlockAssociationPost struct {
	PermissionContext string              `json:"permissionContext"`
	Blocks            []models.PostBlocks `json:"blocks"`
}

func (obj CanvasBlockAssociationPost) Validate() error {
	allowedScope := []string{"block_thread", "reel"}
	if !utils.SliceContainsItem(allowedScope, obj.PermissionContext) {
		return errors.New("Please send `block_thread` or  `reel`")
	}
	return nil
}

type GetCanvasBranches struct {
	CanvasID       uint64 `json:"canvasId"`
	ParentCanvasID uint64 `json:"parentCanvasId"`
	CollectionID   uint64 `json:"collectionId"`
}

type SearchBranchRepos struct {
	Query string `json:"query"`
	//CanvasID       uint64 `json:"canvasId"`
	//ParentCanvasID uint64 `json:"parentCanvasId"`
	//CollectionID   uint64 `json:"collectionId"`
}

func (v SearchBranchRepos) Validate() error {
	if len(v.Query) < 3 {
		return errors.New("Min 3 chars")
	}
	return nil
}

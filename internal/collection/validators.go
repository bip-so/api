package collection

import (
	"errors"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

type CollectionCreateValidator struct {
	Name               string `json:"name" binding:"required"`
	Position           uint   `json:"position" binding:"required"`
	Icon               string `json:"icon"`
	ParentCollectionID uint64 `json:"parentCollectionID"`
	PublicAccess       string `json:"publicAccess" binding:"required"`
}

type CollectionUpdateValidator struct {
	ID           uint64 `json:"id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Icon         string `json:"icon"`
	PublicAccess string `json:"publicAccess" binding:"required"`
}

type VisibilityUpdateValidator struct {
	PublicAccess string `json:"publicAccess" binding:"required"`
}

func (obj VisibilityUpdateValidator) Validate() error {
	allowed := []string{"private", "view", "comment", "edit"}
	if !utils.SliceContainsItem(allowed, obj.PublicAccess) {
		return errors.New("Please only send allowed values. \"private\", \"view\", \"comment\", \"edit\"")
	}
	return nil
}

type CollectionMoveValidator struct {
	CollectionId uint64 `json:"collectionId" binding:"required"`
	Position     uint   `json:"position" binding:"required"`
}

package models

import (
	"errors"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

const (
	ACCESS_REQUEST_PENDING  = "PENDING"
	ACCESS_REQUEST_ACCEPTED = "ACCEPTED"
	ACCESS_REQUEST_REJECTED = "REJECTED"
)

type AccessRequest struct {
	BaseModel

	StudioID                    uint64
	CollectionID                uint64
	CanvasRepositoryID          uint64
	CanvasBranchID              uint64
	CanvasBranchPermissionGroup *string
	Message                     *string
	Status                      string `gorm:"default:'PENDING'"`
	ReviewedByUserID            *uint64
	ReviewedByUser              *User `gorm:"foreignKey:ReviewedByUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	CreatedByID   uint64
	CreatedByUser *User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Studio           *Studio           `gorm:"foreignKey:StudioID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Collection       *Collection       `gorm:"foreignKey:CollectionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CanvasRepository *CanvasRepository `gorm:"foreignKey:CanvasRepositoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CanvasBranch     *CanvasBranch     `gorm:"foreignkey:CanvasBranchID;constraint:OnDelete:CASCADE;"`
}

func (obj AccessRequest) Validate() error {
	allowedScope := []string{PGCanvasNoneSysName, PGCanvasViewSysName, PGCanvasCommentSysName, PGCanvasEditSysName, PGCanvasModerateSysName}
	if !utils.SliceContainsItem(allowedScope, *obj.CanvasBranchPermissionGroup) {
		return errors.New("Please send corrent perms for CB.")
	}
	allowedStatus := []string{ACCESS_REQUEST_PENDING, ACCESS_REQUEST_ACCEPTED, ACCESS_REQUEST_REJECTED}
	if !utils.SliceContainsItem(allowedStatus, obj.Status) {
		return errors.New("Please send corrent perms for CB.")
	}

	return nil
}

func (obj AccessRequest) NewAccessRequest(studioID uint64, collectionID uint64, repoID uint64, branchID uint64, userID uint64) *AccessRequest {
	var instance = AccessRequest{
		StudioID:           studioID,
		CollectionID:       collectionID,
		CanvasRepositoryID: repoID,
		CanvasBranchID:     branchID,
		CreatedByID:        userID,
	}
	return &instance
}

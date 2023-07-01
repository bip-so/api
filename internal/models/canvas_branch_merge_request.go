package models

import (
	"gorm.io/datatypes"
	"time"
)

const (
	MERGE_REQUEST_OPEN               = "OPEN"
	MERGE_REQUEST_ACCEPTED           = "ACCEPTED"
	MERGE_REQUEST_REJECTED           = "REJECTED"
	MERGE_REQUEST_PARTIALLY_ACCEPTED = "PARTIALLY_ACCEPTED"
)

type MergeRequest struct {
	BaseModel
	// canvas id
	CanvasRepositoryID uint64
	CanvasRepository   *CanvasRepository `gorm:"foreignKey:CanvasRepositoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// Source and dest Canvas Branch
	SourceBranchID      uint64
	SourceBranch        *CanvasBranch `gorm:"foreignkey:SourceBranchID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	DestinationBranchID uint64
	DestinationBranch   *CanvasBranch `gorm:"foreignkey:DestinationBranchID;constraint:OnDelete:CASCADE;"`

	CommitMessage  string `gorm:"type: varchar(500)"`
	ClosedByUserId *uint64
	ClosedAt       time.Time

	Status              string `gorm:"type: varchar(20)"` // Consts Above
	CommitID            string `gorm:"type: varchar(40)"`
	SourceCommitID      string `gorm:"type: varchar(40)"`
	DestinationCommitID string `gorm:"type: varchar(40)"`
	ChangesAccepted     datatypes.JSON
	CreatedByID         uint64
	UpdatedByID         uint64
	CreatedByUser       *User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser       *User `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ClosedByUser        *User `gorm:"foreignKey:ClosedByUserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	IsArchived          bool
	ArchivedAt          time.Time
	ArchivedByID        *uint64
	ArchivedByUser      *User `gorm:"foreignKey:ArchivedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (m MergeRequest) NewMergeRequest(canvasRepositoryID uint64, sourceBranchID uint64, destinationBranchID uint64, commitMessage string, status string, createdByID uint64, updatedByID uint64) *MergeRequest {
	return &MergeRequest{CanvasRepositoryID: canvasRepositoryID, SourceBranchID: sourceBranchID, DestinationBranchID: destinationBranchID, CommitMessage: commitMessage, Status: status, CreatedByID: createdByID, UpdatedByID: updatedByID}
}

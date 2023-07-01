package models

/*

Merge_request

canvas CANVAS_ID
source CanvasBranchID
destination CanvasBranchID
message

created_by_id
status : GitMergeStatus()

Post Merge Commit (bip git api )
commit_id : ()
src_commit_id
dest_commit_id

closed_by_user_id
closed_at
changed_accepted ["blockid:true",  ] // Block by block which were accepted or rejcted.

]}

*/

// Parent Canvas
//type MergeRequest struct {
//	BaseModel
//	CreatedByID   uint64
//	UpdatedByID   uint64
//	CreatedByUser *User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
//	UpdatedByUser *User `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
//
//	IsArchived     bool
//	ArchivedAt     time.Time
//	ArchivedByID   uint64
//	ArchivedByUser *User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
//	// canvas id
//	CanvasID uint64
//	// Source and dest Canvas Branch
//	SourceID      uint64
//	DestinationID uint64
//
//	CommitMessage string `gorm:"type: varchar(500)"`
//	// status -> Enum *GitMergeStatus
//
//	// @todo: Paras
//	//Post Merge Commit (bip git api )
//	//commit_id : ()
//	//src_commit_id
//	//dest_commit_id
//
//	ClosedByUserId uint64
//	ClosedAt       time.Time
//
//	Canvas      *Canvas       `gorm:"foreignKey:CanvasID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
//	Source      *CanvasBranch `gorm:"foreignKey:SourceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
//	Destination *CanvasBranch `gorm:"foreignKey:DestinationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
//	User        *User         `gorm:"foreignkey:ClosedByUserId;contraint;OnDelete:SET NULL;"`
//}

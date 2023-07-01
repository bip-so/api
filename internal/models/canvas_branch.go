package models

import "time"

// Todo: Note: Please Talk to PW or CC as there is a migration issue here.
// We need to Create a Mock of CanvasBranch using CreateTable before AutoMigrate
// https://gorm.io/docs/migration.html
// We need to comment out
// https://golang.hotexamples.com/examples/github.com.jinzhu.gorm/DB/CreateTable/golang-db-createtable-method-examples.html

const (
	CANVAS_BRANCH_NAME_MAIN             = "main"
	CANVAS_BRANCH_PUBLIC_ACCESS_PRIVATE = "private"
	REDIS_CANVAS_BRANCH_BLOCKS          = "cached-branch-blocks:"
)

type CanvasBranch struct {
	BaseModel

	CanvasRepositoryID uint64            // Canvas Branch
	CanvasRepository   *CanvasRepository `gorm:"foreignKey:CanvasRepositoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Name string `gorm:"type: varchar(256)"`

	// This is a kinds confusing variable but essentially think of this waay.
	// If By defaut this is True meaning each Branch is by default "IsDraft"

	IsDraft bool `gorm:"default:true"`
	// IsPublished

	IsMerged bool `gorm:"default:false"`

	// This one also is kep for future User
	// We can drop this is needed but this essentially which branch is "Main" Branch
	//v so may be you don't delete it
	IsDefault bool `gorm:"default:false"`

	// Moved from CanvasRepo to CanvasBranch
	PublicAccess string `gorm:"default:private"` //"Private" "View" "Comment"  "Edit" - always save in lower case.

	// Rough Branch Fields
	IsRoughBranch        bool          `gorm:"default:false"` // This is a flag will should
	RoughFromBranchID    *uint64       // The Branch from which this Rough Branch was created
	RoughFromBranch      *CanvasBranch `gorm:"foreignkey:RoughFromBranchID;constraint:OnDelete:CASCADE;"`
	RoughBranchCreatorID *uint64       // Person who started this Rough Branch
	RoughBranchCreator   *User         `gorm:"foreignkey:RoughBranchCreatorID;constraint:OnDelete:CASCADE;"`
	WasRoughBranch       bool          `gorm:"default:false"` // This is a flag which tells first it was a rough branch then created to normal branch. In old bip it was type == 'branchDraft'

	// PW Double Check
	FromBranchID        *uint64 // FK to self
	CreatedFromCommitID string
	// This makes a branch Readonly or not
	Committed                         bool
	LastSyncedAllAttributionsCommitID string // ??

	CreatedByID uint64
	UpdatedByID uint64
	// Soft delete
	IsArchived   bool
	ArchivedAt   time.Time
	ArchivedByID *uint64
	Key          string

	CreatedByUser  *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser  *User         `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ArchivedByUser *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	FromBranch     *CanvasBranch `gorm:"foreignkey:FromBranchID;constraint:OnDelete:CASCADE;"`
}

func (m *CanvasBranch) TableName() string {
	return "canvas_branches"
}

func NewCanvasBranch(name string, canvasRepositoryID uint64, createdByUserID uint64, publicAccess string) *CanvasBranch {
	return &CanvasBranch{
		Name:               name,
		CanvasRepositoryID: canvasRepositoryID,
		CreatedByID:        createdByUserID,
		UpdatedByID:        createdByUserID,
		IsDefault:          name == CANVAS_BRANCH_NAME_MAIN,
		PublicAccess:       publicAccess,
	}
}

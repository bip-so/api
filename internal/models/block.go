package models

import (
	"time"

	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

/*
Block Model

This is Model for Block and This is the primary Model we'll have
CanvasBranch has Many Block.
Each Block Will have different Data
We'll have a Block Thing

- CR
	- CB1
		- B1
		- B2
		- B3
	- CB2
		- B2
		- B3

There is a JSON Key here, which is what the FE uses to dump data
Data
	- position : uint
	- version: 2
	- type Type
	- Data
			- commnt on this block
			- data: {}
			- embeds {}
			- com

We'll discuss later
- Comments
- Reactions
- Mentions
- Reels
*/

type BlockAttribution struct {
	UserUUID  string
	Username  string
	FullName  string
	TimeStamp time.Time // TZ aware??
	AvatarUrl string
}

/*
	Create a block2/Branch2
	- UUID()
	block2/branch3
Main -> Forked Branch -> (Copy) -> Merge -> Main

	UUID:Block1 - Branch1
	UUID:Block1 - Branch2
	UUID:Block1 - Branch3
*/

/*CR
	CB1
	CB2
	CB3
		Bloc1,2,
Lifecycle of block across branches.

CRUD

Copy(toBranch, fromBranch)
Merge(toBranch, fromBranch)
	B1 - Block(UUID1), Block(UUID2),  Block(UUID3)
	B2 This branch will also be deleted before merge)
		- (same)Block(UUID1),
		- (data changed)Block(UUID2), // we delete the "parent" block and change the branchid from b2(Forked) to b2(Main)
		- (deleted)Block(UUID3), -> Delete The Block from b1
		- (new)Block(UUID4) -> change the branchid from b2(Forked) to b2(Main)

	List :  (b1, b7) (b2, b2) (b3, b5) (b4, b4) (b5, b3)
	FinalBlocks (b1, b3, b4, b5)

	B1 - 1,2,3,4,5,6,7
	B2 - 7,2,5,4,3,6,1 // Final

Delete(fromBranch)
Move()
--
Move

// Only show in branch created in and carry forward on merge. For now, no option to NOT carry forward on merge.
// Will always carry forward as long as relevant blocks / content is accepted in merge. In the future, will be optional.
*/

func (m Block) TableName() string {
	return "blocks"
}

//type BlockChildren struct {
//	Child datatypes.JSON
//}

type Block struct {
	// Boilerplate Stuff
	ID                 uint64    `gorm:"primaryKey,type:BIGSERIAL UNSIGNED NOT NULL AUTO_INCREMENT"`
	UUID               uuid.UUID `gorm:"type:uuid;index;"`
	CanvasRepositoryID uint64
	CanvasBranchID     *uint64 // Canvas Branch this belongs to _>  CB1
	ParentBlock        *uint64
	Version            uint   //Top level block version (Default will be 2) newer block.
	Type               string // Block Type

	Rank int32 `gorm:"default:0"`

	//Data               datatypes.JSON // Slate spec if Primary data bucket used by FE
	Children datatypes.JSON
	//BlockChildren   pgtype.JSONBArray

	CommentCount uint `gorm:"default:0"` // This is just a count so we can show is UI //- Everything
	ReelCount    uint `gorm:"default:0"` // Number of Reels

	CreatedByID  uint64
	UpdatedByID  uint64
	ArchivedByID uint64
	IsArchived   bool
	ArchivedAt   time.Time
	Contributors datatypes.JSON

	// Purpose of this field is to keep a reference to the orrignal block this was created from
	// Main used in cloning and merging artifacts
	ClonedFromBlockID uint64 `gorm:"default:0"`

	CanvasBranch   *CanvasBranch `gorm:"foreignkey:CanvasBranchID;constraint:OnDelete:CASCADE;"`
	CreatedByUser  *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser  *User         `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ArchivedByUser *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt      time.Time     `gorm:"autoCreateTime"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime"`

	Reactions  datatypes.JSON  // Reactions
	Attributes datatypes.JSON  // Attributes - Extra Slate Data
	Mentions   *datatypes.JSON // Mentions

}

func (m Block) NewBlock(
	uuid uuid.UUID,
	userID, canvasRepoID, canvasBranchID uint64,
	parentBlock *uint64,
	rank int32, version uint,
	blocktype string,
	blockdata, attributes, contributors datatypes.JSON,
) (*Block, error) {
	return &Block{
		UUID:               uuid,
		CanvasRepositoryID: canvasRepoID,
		CanvasBranchID:     &canvasBranchID,
		ParentBlock:        parentBlock,
		Version:            2,
		Type:               blocktype,
		Rank:               rank,
		Children:           blockdata,
		CreatedByID:        userID,
		UpdatedByID:        userID,
		Attributes:         attributes,
		Contributors:       contributors,
	}, nil
}

func MyFirstBlockJson() datatypes.JSON {
	return datatypes.JSON(`[{"text": ""}]`)
}

func MyFirstEmptyBlockJson() datatypes.JSON {
	return datatypes.JSON(`{}`)
}

// GetBlockByUUIDAndBranchID : Get Block ID from UUID and BranchID
func (m Block) GetBlockByUUIDAndBranchID(query map[string]interface{}) (*Block, error) {
	var block Block
	err := postgres.GetDB().Model(&Block{}).Where(query).First(&block).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &block, nil
}

type PostBlocks struct {
	ID           uint64         `json:"id"`
	UUID         uuid.UUID      `json:"uuid"`
	Version      uint           `json:"version"`
	Type         string         `json:"type"`
	Rank         int32          `json:"rank"`
	Children     datatypes.JSON `json:"children"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	CreatedByID  uint64         `json:"createdByID"`
	UpdatedByID  uint64         `json:"updatedByID"`
	ArchivedByID *uint64        `json:"archivedByID"`
	IsArchived   bool           `json:"isArchived"`
	ArchivedAt   time.Time      `json:"archivedAt"`
	Scope        string         `json:"scope"`
	Attributes   datatypes.JSON `json:"attributes"`
}

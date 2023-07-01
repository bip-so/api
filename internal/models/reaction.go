package models

/*
Clubing all reaction models here
follow foillowing pattern and fk to the thing

- emoji
- who
- when
models Blocks, Block Thread, Block Comment, Reel, Reel Comment
This one gets hard deleted
Comment / Block or Reel

Q: How do we query this?
*/

type BlockReaction struct {
	Emoji               string
	CanvasBranchID      *uint64
	BaseModel           // Boilerplate Stuff
	CreatedByID         uint64
	UpdatedByID         uint64
	BlockID             uint64
	ClonedBlockReaction uint64        `gorm:"default:0"`
	Block               *Block        `gorm:"foreignKey:BlockID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CanvasBranch        *CanvasBranch `gorm:"foreignkey:CanvasBranchID;constraint:OnDelete:CASCADE;"`
	CreatedByUser       *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser       *User         `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type BlockThreadReaction struct {
	Emoji                     string
	CanvasBranchID            *uint64
	BaseModel                 // Boilerplate Stuff
	CreatedByID               uint64
	UpdatedByID               uint64
	BlockID                   uint64
	BlockThreadID             uint64
	ClonedBlockThreadReaction uint64 `gorm:"default:0"`

	Block       *Block       `gorm:"foreignKey:BlockID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BlockThread *BlockThread `gorm:"foreignKey:BlockThreadID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CanvasBranch *CanvasBranch `gorm:"foreignkey:CanvasBranchID;constraint:OnDelete:CASCADE;"`

	CreatedByUser *User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type BlockCommentReaction struct {
	Emoji                      string
	CanvasBranchID             *uint64
	BaseModel                  // Boilerplate Stuff
	CreatedByID                uint64
	UpdatedByID                uint64
	BlockID                    uint64
	BlockThreadID              uint64
	BlockCommentID             uint64
	ClonedBlockCommentReaction uint64        `gorm:"default:0"`
	BlockThread                *BlockThread  `gorm:"foreignKey:BlockThreadID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Block                      *Block        `gorm:"foreignKey:BlockID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BlockComment               *BlockComment `gorm:"foreignKey:BlockCommentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CanvasBranch *CanvasBranch `gorm:"foreignkey:CanvasBranchID;constraint:OnDelete:CASCADE;"`

	CreatedByUser *User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type ReelReaction struct {
	Emoji          string
	CanvasBranchID *uint64
	BaseModel      // Boilerplate Stuff
	CreatedByID    uint64
	UpdatedByID    uint64
	ReelID         uint64
	Reel           *Reel         `gorm:"foreignKey:ReelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CanvasBranch   *CanvasBranch `gorm:"foreignkey:CanvasBranchID;constraint:OnDelete:CASCADE;"`
	CreatedByUser  *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser  *User         `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
type ReelCommentReaction struct {
	Emoji          string
	CanvasBranchID *uint64
	BaseModel      // Boilerplate Stuff
	CreatedByID    uint64
	UpdatedByID    uint64
	ReelID         uint64
	ReelCommentID  uint64
	Reel           *Reel         `gorm:"foreignKey:ReelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ReelComment    *ReelComment  `gorm:"foreignKey:ReelCommentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CanvasBranch   *CanvasBranch `gorm:"foreignkey:CanvasBranchID;constraint:OnDelete:CASCADE;"`
	CreatedByUser  *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser  *User         `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type ReactionCounter struct {
	Emoji string `json:"emoji"`
	Count string `json:"count"`
}

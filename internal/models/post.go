package models

import (
	"gorm.io/datatypes"
)

/*

Post
PostComment


Reactions
- PostReaction
- PostCommentReaction

*/

type Post struct {
	BaseModel
	CreatedByID   uint64
	UpdatedByID   uint64
	CreatedByUser *User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	StudioID      uint64
	Studio        *Studio `gorm:"foreignKey:StudioID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// specific
	Children   datatypes.JSON
	Attributes datatypes.JSON // Attributes
	//Mentions  *datatypes.JSON // Mentions
	//Reactions *datatypes.JSON // Reactions
	IsPublic     bool   `gorm:"default:false;"`        // If the is_public true then ignore the Roles
	Roles        []Role `gorm:"many2many:post_roles;"` // Post can be associated with many roles
	CommentCount uint   `gorm:"default:0;"`            // This needs to come from auto service kind of thing
}

// PostComment Table ONE post can have many PostComments
type PostComment struct {
	BaseModel
	PostID        uint64
	Post          *Post `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	IsEdited      bool
	Comment       string
	CreatedByID   uint64
	UpdatedByID   uint64
	CreatedByUser *User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	// Parent Comment
	ParentPostCommentID *uint64
	ParentPostComment   *PostComment `gorm:"foreignKey:ParentPostCommentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CommentCount        uint         `gorm:"default:0;"` // This needs to come from auto service kind of thing
}

type PostReaction struct {
	BaseModel
	PostID        uint64
	Post          *Post `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Emoji         string
	CreatedByID   uint64
	CreatedByUser *User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type PostCommentReaction struct {
	BaseModel

	PostID uint64
	Post   *Post `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	PostCommentID uint64
	PostComment   *PostComment `gorm:"foreignKey:PostCommentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Emoji         string
	CreatedByID   uint64
	CreatedByUser *User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

////////////////////////////////////////////////////////////

//type PostCommentReaction struct {
//	BaseModel
//
//	PostID uint64
//	Post   *Post `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
//
//	PostCommentID uint64
//	PostComment   *PostComment `gorm:"foreignKey:PostCommentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
//
//	Emoji         string
//	CreatedByID   uint64
//	CreatedByUser *User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
//}

//type PostCommentReplyReaction struct {
//	BaseModel
//
//	PostID uint64
//	Post   *Post `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
//
//	PostCommentID uint64
//	PostComment   *PostComment `gorm:"foreignKey:PostCommentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
//
//	PostCommentReplyID uint64
//	PostCommentReply   *PostCommentReply `gorm:"foreignKey:PostCommentReplyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
//
//	Emoji         string
//	CreatedByID   uint64
//	CreatedByUser *User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
//}

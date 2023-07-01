package post

import (
	"gorm.io/datatypes"
)

// Post
type NewPostThread struct {
	IsPublic   bool           `json:"isPublic"`
	Children   datatypes.JSON `json:"children"`
	Attributes datatypes.JSON `json:"attributes"`
	Roles      []uint64       `json:"roleIds"`
}

type UpdatePostThread struct {
	IsPublic   bool           `json:"isPublic"`
	Children   datatypes.JSON `json:"children"`
	Attributes datatypes.JSON `json:"attributes"`
	Roles      []uint64       `json:"roleIds"`
}

// Post Comment
type CreatePostCommentValidation struct {
	IsEdited            bool   `json:"isEdited"`
	Comment             string `json:"comment"`
	ParentPostCommentID uint64 `json:"parentPostCommentID"`
}

type UpdatePostCommentValidation struct {
	IsEdited bool
	Comment  string
}

// Post Reactions
type NewPostReaction struct {
	Emoji string `json:"emoji" binding:"required"`
}
type RemovePostReaction struct {
	Emoji string `json:"emoji" binding:"required"`
}

type NewPostCommentReaction struct {
	Emoji string `json:"emoji" binding:"required"`
}
type RemovePostCommentReaction struct {
	Emoji string `json:"emoji" binding:"required"`
}

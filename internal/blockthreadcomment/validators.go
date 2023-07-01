package blockThreadCommentcomment

import (
	"gorm.io/datatypes"
)

type PostBlockThreadComment struct {
	ThreadID uint64         `json:"threadId" binding:"required"`
	ParentID *uint64        `json:"parentId"`
	Position uint           `json:"position"`
	Data     datatypes.JSON `json:"data" binding:"required"`
	IsReply  bool           `json:"isReply"`
	IsEdited bool           `json:"isEdited"`
}

type PatchBlockThreadComment struct {
	ID       uint64         `json:"id" binding:"required"`
	ThreadID uint64         `json:"threadId" binding:"required"`
	ParentID *uint64        `json:"parentId" binding:"required"`
	Position uint           `json:"position" binding:"required"`
	Data     datatypes.JSON `json:"data" binding:"required"`
	IsReply  bool           `json:"isReply" binding:"required"`
}

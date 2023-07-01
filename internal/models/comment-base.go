package models

import "gorm.io/datatypes"

type CommentBase struct {
	Position uint           // May be useful may comment in same
	Data     datatypes.JSON // Slate spec if Primary data bucket used by FE {"text" : "hello"} // Always needs to have "text" which is a plain text version
	IsEdited bool
	IsReply  bool
}

package notion

import "mime/multipart"

type PostImportNotion struct {
	StudioID uint64               `json:"studioId"`
	File     multipart.FileHeader `json:"file" form:"file"`
}

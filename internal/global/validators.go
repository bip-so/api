package global

import "mime/multipart"

type ImageUpload struct {
	File   multipart.FileHeader `json:"file" form:"file"`
	Model  string               `json:"model" form:"model"`
	UUID   string               `json:"uuid" form:"uuid"`
	RepoID uint64               `json:"repoId" form:"repoId"`
}

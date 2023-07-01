package canvasrepo

type InitCanvasRepoPost struct {
	CollectionID             uint64 `json:"collectionID" binding:"required"`
	Name                     string `json:"name" binding:"required"`
	Icon                     string `json:"icon"`
	Position                 uint   `json:"position" binding:"required"`
	ParentCanvasRepositoryID uint64 `json:"parentCanvasRepositoryID"`
}

type NewCanvasRepoPost struct {
	CollectionID             uint64 `json:"collectionID" binding:"required"`
	ParentCanvasRepositoryID uint64 `json:"parentCanvasRepositoryID"`
	Name                     string `json:"name" binding:"required"`
	Icon                     string `json:"icon"`
	Position                 uint   `json:"position" binding:"required"`
}

type UpdateCanvasRepoPost struct {
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	CoverUrl string `json:"coverUrl"`
}

type MoveCanvasRepoPost struct {
	ToCollectionID             uint64 `json:"toCollectionID"`
	ToParentCanvasRepositoryID uint64 `json:"toParentCanvasRepositoryID"`
	CanvasRepoID               uint64 `json:"canvasRepoID"`
	FuturePosition             uint   `json:"futurePosition"`
}

type GetAllCanvasValidator struct {
	ParentCollectionID       uint64 `json:"parentCollectionID"`
	ParentCanvasRepositoryID uint64 `json:"parentCanvasRepositoryID"`
}

type CreateLanguageValidator struct {
	CanvasRepositoryID uint64   `json:"canvasRepositoryID"`
	Languages          []string `json:"languages"`
	AutoTranslate      bool     `json:"autoTranslate"`
}

package canvasbranchpermissions

// post
type NewCanvasBranchPermissionCreatePost struct {
	CollectionId                uint64 `json:"collectionId" binding:"required"`
	CanvasBranchId              uint64 `json:"canvasBranchId" binding:"required"`
	CanvasRepositoryID          uint64 `json:"canvasRepositoryId" binding:"required"`
	CbpParentCanvasRepositoryID uint64 `json:"parentCanvasRepositoryId"`
	PermGroup                   string `json:"permGroup" binding:"required"`
	RoleID                      uint64 `json:"roleID"`
	MemberID                    uint64 `json:"memberID"`
	UserID                      uint64 `json:"userID"`
	IsOverridden                bool   `json:"isOverridden"`
}

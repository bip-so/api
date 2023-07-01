package permissions

func (s permissionService) CalculateCollectionRolePermissions(roleID uint64) map[uint64]string {
	collectionPerms, _ := App.Repo.getCollectionPermissions(map[string]interface{}{"role_id": roleID})
	permissionList := make(map[uint64]string)
	for _, perm := range collectionPerms {
		permissionList[perm.CollectionId] = perm.PermissionGroup
	}
	return permissionList
}

func (s permissionService) CalculateCanvasRolePermissions(roleID uint64, collectionIDs []uint64) map[uint64]map[uint64]string {
	permissionList := make(map[uint64]map[uint64]string)
	canvasPerms, _ := App.Repo.GetCanvasBranchPerms(map[string]interface{}{"role_id": roleID, "collection_id": collectionIDs})

	for _, perm := range canvasPerms {
		if len(permissionList[perm.CanvasRepositoryID]) == 0 {
			permissionList[perm.CanvasRepositoryID] = make(map[uint64]string)
		}
		if perm.CanvasBranchID != nil {
			permissionList[perm.CanvasRepositoryID][*perm.CanvasBranchID] = perm.PermissionGroup
		}
	}
	return permissionList
}

func (s permissionService) CalculateSubCanvasRolePermissions(roleID uint64, parentCanvasRepoID uint64) map[uint64]map[uint64]string {
	permissionList := make(map[uint64]map[uint64]string)
	canvasPerms, _ := App.Repo.GetCanvasBranchPerms(map[string]interface{}{"role_id": roleID, "cbp_parent_canvas_repository_id": parentCanvasRepoID})

	for _, perm := range canvasPerms {
		if len(permissionList[perm.CanvasRepositoryID]) == 0 {
			permissionList[perm.CanvasRepositoryID] = make(map[uint64]string)
		}
		if perm.CanvasBranchID != nil {
			permissionList[perm.CanvasRepositoryID][*perm.CanvasBranchID] = perm.PermissionGroup
		}
	}
	return permissionList
}

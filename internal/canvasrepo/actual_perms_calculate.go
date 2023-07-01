package canvasrepo

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

func RoleBranchPermissionActualCalculator(collectionID uint64, repoID, branchID uint64, member *models.Member, studioID uint64) []RoleBranchActualPermissionsObject {
	var actual []RoleBranchActualPermissionsObject
	// This mean we didn't find anything on the branch with membership
	// Now we query on Perms table where there are no members
	branchPermissionsObjectsViaRole, _ := App.Repo.GetBranchesPermissionAll(map[string]interface{}{
		"studio_id":            studioID,
		"collection_id":        collectionID,
		"member_id":            nil,
		"canvas_repository_id": repoID,
		"canvas_branch_id":     branchID,
	})
	for _, val := range *branchPermissionsObjectsViaRole {
		var count int64
		// Looping all collectionperms and checking is the memner is found in the role.
		results := postgres.GetDB().Table("role_members").Where("role_id = ? and member_id = ?", val.Role.ID, member.ID)
		results.Count(&count)
		fmt.Println("Count", count)
		if count == 1 {
			y := RoleBranchActualPermissionsObject{
				BranchPermissionID: val.ID,
				CollectionID:       collectionID,
				RepoID:             repoID,
				BranchID:           branchID,
				IsOverRidden:       val.IsOverridden,
				RoleID:             val.Role.ID,
				Name:               val.Role.Name,
				PG:                 val.PermissionGroup,
			}
			actual = append(actual, y)
		}
	}
	return actual
}

func MemberBranchPermissionActualCalculator(collectionID uint64, repoID, branchID uint64, member *models.Member, studioID uint64) MemberBranchActualPermissionsObject {
	var actual MemberBranchActualPermissionsObject
	// Check Perms via Member ID
	branchPermissionObjectViaMember, _ := App.Repo.GetBranchPermission(map[string]interface{}{
		"studio_id":            studioID,
		"collection_id":        collectionID,
		"canvas_repository_id": repoID,
		"canvas_branch_id":     branchID,
		"member_id":            member.ID,
	})
	if branchPermissionObjectViaMember.ID != 0 {
		actual = MemberBranchActualPermissionsObject{
			BranchPermissionID: branchPermissionObjectViaMember.ID,
			CollectionID:       collectionID,
			RepoID:             repoID,
			BranchID:           branchID,
			IsOverRidden:       branchPermissionObjectViaMember.IsOverridden,
			MemberID:           member.ID,
			PG:                 branchPermissionObjectViaMember.PermissionGroup,
		}
		return actual
	}

	return actual
}

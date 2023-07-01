package collection

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

type CollectionRoleActualPermissionsObject struct {
	RoleID uint64 `json:"roleId"`
	Name   string `json:"name"`
	PG     string `json:"pg"`
}
type CollectionMemberActualPermissionsObject struct {
	MemberID uint64 `json:"memberId"`
	PG       string `json:"pg"`
}
type CollectionActualPermissionsObject struct {
	CollectionPermissionID uint64                                   `json:"collectionPermissionID"`
	CollectionID           uint64                                   `json:"collectionID"`
	ActualRole             *CollectionRoleActualPermissionsObject   `json:"actualRole"`
	ActualMember           *CollectionMemberActualPermissionsObject `json:"actualMember"`
	IsOverRidden           bool                                     `json:"isOverRidden"`
}

func PluckTheObject(cpo []CollectionActualPermissionsObject, collectionID uint64) CollectionActualPermissionsObject {
	for _, val := range cpo {
		if val.CollectionPermissionID == 0 {
			continue
		}
		if val.CollectionID == collectionID {
			return val
		}
	}
	return CollectionActualPermissionsObject{}
}

// This is slightly unoptimized
func CollectionPermissionActual(collectionID uint64, memberID uint64, studioID uint64) []CollectionActualPermissionsObject {
	var actual []CollectionActualPermissionsObject

	// Check Perms via Member ID
	collectionPermissionObjectViaMember, _ := App.Repo.GetCollectionPermission(map[string]interface{}{"studio_id": studioID, "collection_id": collectionID, "member_id": memberID})
	if collectionPermissionObjectViaMember.ID != 0 {
		x := CollectionActualPermissionsObject{
			CollectionPermissionID: collectionPermissionObjectViaMember.ID,
			CollectionID:           collectionID,
			ActualRole:             nil,
			IsOverRidden:           collectionPermissionObjectViaMember.IsOverridden,
			ActualMember: &CollectionMemberActualPermissionsObject{
				MemberID: memberID,
				PG:       collectionPermissionObjectViaMember.PermissionGroup,
			},
		}
		//temp.ID = collectionPermissionObjectViaMember.ID
		//temp.ActualRole = nil
		//temp.IsOverRidden =
		//temp.ActualMember.MemberID = memberID
		//temp.ActualMember.PG = collectionPermissionObjectViaMember.PermissionGroup
		actual = append(actual, x)
	}

	// This mean we didn't find anything on the collection with membership
	// Now we query on Perms table where there are no members
	collectionPermissionsObjectsViaRole, _ := App.Repo.GetCollectionsPermissionAll(map[string]interface{}{"studio_id": studioID, "collection_id": collectionID, "member_id": nil})
	for _, val := range *collectionPermissionsObjectsViaRole {
		var count int64
		// Looping all collectionperms and checking is the memner is found in the role.
		if val.Role == nil {
			continue
		}
		results := postgres.GetDB().Table("role_members").Where("role_id = ? and member_id = ?", val.Role.ID, memberID)
		fmt.Println("results.RowsAffected", results.RowsAffected)
		fmt.Println("Count", results.Count(&count))
		if results.RowsAffected == 1 {
			y := CollectionActualPermissionsObject{
				CollectionPermissionID: val.ID,
				CollectionID:           collectionID,
				ActualRole: &CollectionRoleActualPermissionsObject{
					RoleID: val.Role.ID,
					Name:   val.Role.Name,
					PG:     val.PermissionGroup,
				},
				IsOverRidden: val.IsOverridden,
				ActualMember: nil,
			}

			//temp2.ID = val.ID
			//temp2.ActualMember = nil
			//temp2.IsOverRidden = val.IsOverridden
			//temp2.ActualRole.PG = val.PermissionGroup
			//temp2.ActualRole.RoleID = val.Role.ID
			//temp2.ActualRole.Name = val.Role.Name
			actual = append(actual, y)
		}
	}

	return actual
}

func MemberCollectionPermissionActualCalculator(collectionID uint64, member *models.Member, studioID uint64) MemberCollectionActualPermissionsObject {
	var actual MemberCollectionActualPermissionsObject
	// Check Perms via Member ID
	collectionPermissionObjectViaMember, _ := App.Repo.GetCollectionPermission(map[string]interface{}{
		"studio_id":     studioID,
		"collection_id": collectionID,
		"member_id":     member.ID,
	})
	if collectionPermissionObjectViaMember.ID != 0 {
		actual = MemberCollectionActualPermissionsObject{
			CollectionPermissionID: collectionPermissionObjectViaMember.ID,
			CollectionID:           collectionID,
			IsOverRidden:           collectionPermissionObjectViaMember.IsOverridden,
			MemberID:               member.ID,
			PG:                     collectionPermissionObjectViaMember.PermissionGroup,
		}
		return actual
	}

	return actual
}

func RoleCollectionPermissionActualCalculator(collectionID uint64, member *models.Member, studioID uint64) []RoleCollectionActualPermissionsObject {
	var actual []RoleCollectionActualPermissionsObject
	// This mean we didn't find anything on the branch with membership
	// Now we query on Perms table where there are no members
	collectionPermissionsObjectsViaRole, _ := App.Repo.GetCollectionsPermissionAll(map[string]interface{}{
		"studio_id":     studioID,
		"collection_id": collectionID,
		"member_id":     nil,
	})
	for _, val := range *collectionPermissionsObjectsViaRole {
		var count int64
		// Looping all collectionperms and checking is the memner is found in the role.
		if val.Role == nil {
			continue
		}
		results := postgres.GetDB().Table("role_members").Where("role_id = ? and member_id = ?", val.Role.ID, member.ID)
		results.Count(&count)
		fmt.Println("Count", count)
		if count == 1 {
			y := RoleCollectionActualPermissionsObject{
				CollectionPermissionID: val.ID,
				CollectionID:           collectionID,
				IsOverRidden:           val.IsOverridden,
				RoleID:                 val.Role.ID,
				Name:                   val.Role.Name,
				PG:                     val.PermissionGroup,
			}
			actual = append(actual, y)
		}
	}
	return actual
}

package collection

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

type CollectionSerializer struct {
	Id                      uint64                                  `json:"id"`
	Name                    string                                  `json:"name"`
	Position                uint                                    `json:"position"`
	Icon                    string                                  `json:"icon"`
	StudioID                uint64                                  `json:"studioID"`
	UserID                  uint64                                  `json:"userID"`
	ParentCollectionID      *uint64                                 `json:"parentCollectionID"`
	PublicAccess            string                                  `json:"publicAccess"`
	ComputedRootCanvasCount int                                     `json:"computedRootCanvasCount"`
	ComputedAllCanvasCount  int                                     `json:"computedAllCanvasCount"`
	Type                    string                                  `json:"type"`
	Parent                  int                                     `json:"parent"`
	Permission              string                                  `json:"permission"`
	Rank                    int32                                   `json:"rank"`
	HasPublicCanvas         bool                                    `json:"hasPublicCanvas"`
	ActualPermsObject       CollectionActualPermissionsObject       `json:"actualPermsObject"`
	MemberPermsObject       MemberCollectionActualPermissionsObject `json:"actualMemberPermsObject"`
	RolePermsObject         []RoleCollectionActualPermissionsObject `json:"actualRolePermsObject"`
}

func CollectionSerializerData(collection *models.Collection) CollectionSerializer {
	return CollectionSerializer{
		Id:                      collection.ID,
		Name:                    collection.Name,
		Position:                collection.Position,
		Icon:                    collection.Icon,
		StudioID:                collection.StudioID,
		UserID:                  collection.CreatedByID,
		ParentCollectionID:      collection.ParentCollectionID,
		PublicAccess:            collection.PublicAccess,
		ComputedRootCanvasCount: collection.ComputedRootCanvasCount,
		ComputedAllCanvasCount:  collection.ComputedAllCanvasCount,
		Type:                    "COLLECTION",
		Rank:                    collection.Rank,
		HasPublicCanvas:         collection.HasPublicCanvas,
	}
}

func MultiCollectionSerializerData(collections *[]models.Collection) *[]CollectionSerializer {
	collectionsData := &[]CollectionSerializer{}

	for _, collection := range *collections {
		view := CollectionSerializer{
			Id:                      collection.ID,
			Name:                    collection.Name,
			Position:                collection.Position,
			Icon:                    collection.Icon,
			StudioID:                collection.StudioID,
			UserID:                  collection.CreatedByID,
			ParentCollectionID:      collection.ParentCollectionID,
			PublicAccess:            collection.PublicAccess,
			ComputedRootCanvasCount: collection.ComputedRootCanvasCount,
			ComputedAllCanvasCount:  collection.ComputedAllCanvasCount,
			Type:                    "COLLECTION",
			Permission:              models.AnonymousUserPerms["collection"],
			Rank:                    collection.Rank,
			HasPublicCanvas:         collection.HasPublicCanvas,
		}
		if collection.HasPublicCanvas {
			view.Permission = models.PGCollectionNoneSysName
		}
		*collectionsData = append(*collectionsData, view)
	}
	return collectionsData
}

type MemberCollectionActualPermissionsObject struct {
	CollectionPermissionID uint64 `json:"collectionPermissionID"`
	CollectionID           uint64 `json:"collectionID"`
	IsOverRidden           bool   `json:"isOverRidden"`
	MemberID               uint64 `json:"memberId"`
	PG                     string `json:"pg"`
}

type RoleCollectionActualPermissionsObject struct {
	CollectionPermissionID uint64 `json:"collectionPermissionID"`
	CollectionID           uint64 `json:"collectionID"`
	IsOverRidden           bool   `json:"isOverRidden"`
	RoleID                 uint64 `json:"roleId"`
	Name                   string `json:"name"`
	PG                     string `json:"pg"`
}

func CollectionSerializerDataMini(collection *models.Collection) *CollectionSerializer {
	if collection == nil || collection.ID == 0 {
		return nil
	}
	return &CollectionSerializer{
		Id:                      collection.ID,
		Name:                    collection.Name,
		Position:                collection.Position,
		Icon:                    collection.Icon,
		StudioID:                collection.StudioID,
		UserID:                  collection.CreatedByID,
		ParentCollectionID:      collection.ParentCollectionID,
		PublicAccess:            collection.PublicAccess,
		ComputedRootCanvasCount: collection.ComputedRootCanvasCount,
		ComputedAllCanvasCount:  collection.ComputedAllCanvasCount,
		Type:                    "COLLECTION",
		Rank:                    collection.Rank,
		HasPublicCanvas:         collection.HasPublicCanvas,
	}
}

package canvasrepo

import (
	"time"

	"github.com/gosimple/slug"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

type BranchRoleActualPermissionsObject struct {
}
type MemberBranchActualPermissionsObject struct {
	BranchPermissionID uint64 `json:"branchPermissionID"`
	CollectionID       uint64 `json:"collectionID"`
	RepoID             uint64 `json:"repoID"`
	BranchID           uint64 `json:"branchID"`
	IsOverRidden       bool   `json:"isOverRidden"`
	MemberID           uint64 `json:"memberId"`
	PG                 string `json:"pg"`
}
type RoleBranchActualPermissionsObject struct {
	BranchPermissionID uint64 `json:"branchPermissionID"`
	CollectionID       uint64 `json:"collectionID"`
	RepoID             uint64 `json:"repoID"`
	BranchID           uint64 `json:"branchID"`
	IsOverRidden       bool   `json:"isOverRidden"`
	RoleID             uint64 `json:"roleId"`
	Name               string `json:"name"`
	PG                 string `json:"pg"`
}
type CanvasBranchMiniSerializer struct {
	ID                uint64                              `json:"id"`
	Name              string                              `json:"name"`
	UUID              string                              `json:"uuid"`
	Key               string                              `json:"key"`
	PublicAccess      string                              `json:"publicAccess"`
	Permission        string                              `json:"permission"`
	HasPublishRequest bool                                `json:"hasPublishRequest"`
	CanPublish        bool                                `json:"canPublish"`
	IsPublishedBy     uint64                              `json:"isPublishedBy"`
	Slug              string                              `json:"slug"`
	MemberPermsObject MemberBranchActualPermissionsObject `json:"actualMemberPermsObject"`
	RolePermsObject   []RoleBranchActualPermissionsObject `json:"actualRolePermsObject"`
}

func checkHasPublishRequestOnDefaultBranch(branch *models.CanvasBranch) bool {
	return App.Repo.checkHasPublishRequestOnDefaultBranch(branch)
}

func SerializeCanvasBranchMini(branch *models.CanvasBranch) *CanvasBranchMiniSerializer {
	view := &CanvasBranchMiniSerializer{
		ID:                branch.ID,
		UUID:              branch.UUID.String(),
		Name:              branch.Name,
		Key:               branch.Key,
		PublicAccess:      branch.PublicAccess,
		HasPublishRequest: checkHasPublishRequestOnDefaultBranch(branch),
		Slug:              slug.Make(branch.Name),
	}

	return view
}

type CanvasRepoDefaultSerializer struct {
	ID                          uint64                     `json:"id"`
	UUID                        string                     `json:"uuid"`
	CollectionID                uint64                     `json:"collectionID"`
	Name                        string                     `json:"name"`
	Position                    uint                       `json:"position"`
	Icon                        string                     `json:"icon"`
	CoverUrl                    string                     `json:"coverUrl"`
	IsPublished                 bool                       `json:"isPublished"`
	ParentCanvasRepositoryID    *uint64                    `json:"parentCanvasRepositoryID"`
	CreatedAt                   time.Time                  `json:"createdAt"`
	UpdatedAt                   time.Time                  `json:"updatedAt"`
	CreatedByID                 uint64                     `json:"createdByID"`
	UpdatedByID                 uint64                     `json:"updatedByID"`
	Key                         string                     `json:"key"`
	DefaultBranchID             *uint64                    `json:"defaultBranchID"`
	DefaultBranch               CanvasBranchMiniSerializer `json:"defaultBranch"`
	Type                        string                     `json:"type"`
	Parent                      *uint64                    `json:"parent"`
	SubCanvasCount              int                        `json:"subCanvasCount"`
	Rank                        int32                      `json:"rank"`
	DefaultLanguageCanvasRepoID *uint64                    `json:"defaultLanguageCanvasRepoId"`
	Language                    *string                    `json:"language"`
	IsLanguageCanvas            bool                       `json:"isLanguageCanvas"`
	AutoTranslated              bool                       `json:"autoTranslated"`
	HasPublicCanvas             bool                       `json:"hasPublicCanvas"`
	MergeRequestCount           int64                      `json:"mergeRequestCount"`
	PrivateRepoCount            int64                      `json:"privateRepoCount"`
	MaxAllowedPrivateRepo       int                        `json:"maxAllowedPrivateRepo"`
	Nudge                       bool                       `json:"nudge"`
	SearchMatch                 bool                       `json:"searchMatch"`
}

func SerializeDefaultCanvasRepo(cr *models.CanvasRepository) *CanvasRepoDefaultSerializer {
	flag, count, maxCount := App.Repo.StudioPlanCheck(cr.StudioID)
	view := CanvasRepoDefaultSerializer{
		Nudge:                       flag,
		PrivateRepoCount:            count,
		MaxAllowedPrivateRepo:       maxCount,
		MergeRequestCount:           App.Repo.MergeRequestCount(cr.ID),
		ID:                          cr.ID,
		UUID:                        cr.UUID.String(),
		Key:                         cr.Key,
		CollectionID:                cr.CollectionID,
		Name:                        cr.Name,
		Position:                    cr.Position,
		Icon:                        cr.Icon,
		CoverUrl:                    cr.CoverUrl,
		IsPublished:                 cr.IsPublished,
		DefaultBranchID:             cr.DefaultBranchID,
		ParentCanvasRepositoryID:    cr.ParentCanvasRepositoryID,
		SubCanvasCount:              cr.SubCanvasCount,
		CreatedAt:                   cr.CreatedAt,
		UpdatedAt:                   cr.UpdatedAt,
		CreatedByID:                 cr.CreatedByID,
		UpdatedByID:                 cr.UpdatedByID,
		Type:                        "CANVAS",
		Rank:                        cr.Rank,
		DefaultLanguageCanvasRepoID: cr.DefaultLanguageCanvasRepoID,
		Language:                    cr.Language,
		IsLanguageCanvas:            cr.IsLanguageCanvas,
		AutoTranslated:              cr.AutoTranslated,
		HasPublicCanvas:             cr.HasPublicCanvas,
		SearchMatch:                 true,
	}
	if view.ParentCanvasRepositoryID != nil {
		view.Parent = view.ParentCanvasRepositoryID
	} else {
		view.Parent = &view.CollectionID
	}
	//fmt.Println("BRANCHHHHH")
	//fmt.Println(cr.DefaultBranch)
	//fmt.Println(*cr.DefaultBranch)

	view.DefaultBranch = *SerializeCanvasBranchMini(cr.DefaultBranch)
	return &view
}

func SerializeDefaultCanvasRepoMini(cr *models.CanvasRepository) *CanvasRepoDefaultSerializer {
	if cr == nil {
		return nil
	}
	view := CanvasRepoDefaultSerializer{
		MergeRequestCount:           App.Repo.MergeRequestCount(cr.ID),
		ID:                          cr.ID,
		UUID:                        cr.UUID.String(),
		Key:                         cr.Key,
		CollectionID:                cr.CollectionID,
		Name:                        cr.Name,
		Position:                    cr.Position,
		Icon:                        cr.Icon,
		CoverUrl:                    cr.CoverUrl,
		IsPublished:                 cr.IsPublished,
		DefaultBranchID:             cr.DefaultBranchID,
		ParentCanvasRepositoryID:    cr.ParentCanvasRepositoryID,
		SubCanvasCount:              cr.SubCanvasCount,
		CreatedAt:                   cr.CreatedAt,
		UpdatedAt:                   cr.UpdatedAt,
		CreatedByID:                 cr.CreatedByID,
		UpdatedByID:                 cr.UpdatedByID,
		Type:                        "CANVAS",
		Rank:                        cr.Rank,
		DefaultLanguageCanvasRepoID: cr.DefaultLanguageCanvasRepoID,
		Language:                    cr.Language,
		IsLanguageCanvas:            cr.IsLanguageCanvas,
		AutoTranslated:              cr.AutoTranslated,
		HasPublicCanvas:             cr.HasPublicCanvas,
	}
	if view.ParentCanvasRepositoryID != nil {
		view.Parent = view.ParentCanvasRepositoryID
	} else {
		view.Parent = &view.CollectionID
	}
	return &view
}

func MultiSerializeDefaultCanvasRepo(crs *[]models.CanvasRepository) *[]CanvasRepoDefaultSerializer {
	canvasRepoViews := &[]CanvasRepoDefaultSerializer{}
	for _, cr := range *crs {
		canvasRepoView := CanvasRepoDefaultSerializer{
			ID:                          cr.ID,
			UUID:                        cr.UUID.String(),
			Key:                         cr.Key,
			CollectionID:                cr.CollectionID,
			Name:                        cr.Name,
			Position:                    cr.Position,
			Icon:                        cr.Icon,
			IsPublished:                 cr.IsPublished,
			DefaultBranchID:             cr.DefaultBranchID,
			ParentCanvasRepositoryID:    cr.ParentCanvasRepositoryID,
			SubCanvasCount:              cr.SubCanvasCount,
			CreatedAt:                   cr.CreatedAt,
			UpdatedAt:                   cr.UpdatedAt,
			CreatedByID:                 cr.CreatedByID,
			UpdatedByID:                 cr.UpdatedByID,
			Type:                        "CANVAS",
			Rank:                        cr.Rank,
			HasPublicCanvas:             cr.HasPublicCanvas,
			IsLanguageCanvas:            cr.IsLanguageCanvas,
			Language:                    cr.Language,
			DefaultLanguageCanvasRepoID: cr.DefaultLanguageCanvasRepoID,
		}
		if canvasRepoView.ParentCanvasRepositoryID != nil {
			canvasRepoView.Parent = canvasRepoView.ParentCanvasRepositoryID
		} else {
			canvasRepoView.Parent = &canvasRepoView.CollectionID
		}
		canvasRepoView.DefaultBranch = *SerializeCanvasBranchMini(cr.DefaultBranch)
		if cr.HasPublicCanvas {
			canvasRepoView.DefaultBranch.Permission = models.PGCanvasViewMetadataSysName
		} else {
			canvasRepoView.DefaultBranch.Permission = models.AnonymousUserPerms["canvas_branch"]
		}
		*canvasRepoViews = append(*canvasRepoViews, canvasRepoView)
	}
	return canvasRepoViews
}

type DiscordMessagesData struct {
	CollectionsMap    map[string]map[string]interface{} `json:"collectionsMap"`
	CanvasesChannelID string                            `json:"canvasesChannelId"`
	MessageIDs        []string                          `json:"messageIds"`
}

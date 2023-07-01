package canvasbranch

import (
	"fmt"
	"time"

	"gitlab.com/phonepost/bip-be-platform/internal/canvasrepo"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
)

type CanvasBranchMiniSerializer struct {
	ID                                uint64    `json:"id"`
	Name                              string    `json:"name"`
	UUID                              string    `json:"uuid"`
	Key                               string    `json:"key"`
	PublicAccess                      string    `json:"publicAccess"`
	Type                              string    `json:"type"`
	CanvasRepositoryID                uint64    `json:"CanvasRepositoryId"`
	IsDraft                           bool      `json:"isDraft"`
	IsMerged                          bool      `json:"isMerged"`
	IsDefault                         bool      `json:"isDefault"`
	IsRoughBranch                     bool      `json:"isRoughBranch"`
	RoughFromBranchID                 uint64    `json:"roughFromBranchId"`
	RoughBranchCreatorID              uint64    `json:"roughBranchCreatorId"`
	FromBranchID                      uint64    `json:"fromBranchId"`
	CreatedFromCommitID               string    `json:"createdFromCommitId"`
	Committed                         bool      `json:"committed"`
	LastSyncedAllAttributionsCommitID string    `json:"lastSyncedAllAttributionsCommitId"`
	CreatedByID                       uint64    `json:"createdById"`
	UpdatedByID                       uint64    `json:"updatedById"`
	IsArchived                        bool      `json:"isArchived"`
	ArchivedAt                        time.Time `json:"archivedAt"`
	ArchivedByID                      uint64    `json:"archivedById"`
	Permission                        string    `json:"permission"`
	MergeRequestCount                 int       `json:"mergeRequestCount"`
}

type CanvasRepoDefaultSerializer struct {
	ID                       uint64                       `json:"id"`
	UUID                     string                       `json:"uuid"`
	Name                     string                       `json:"name"`
	HasBranch                bool                         `json:"hasBranch"`
	Branches                 []CanvasBranchMiniSerializer `json:"branches"`
	Position                 uint                         `json:"position"`
	Icon                     string                       `json:"icon"`
	IsPublished              bool                         `json:"isPublished"`
	ParentCanvasRepositoryID *uint64                      `json:"parentCanvasRepositoryID"`
	CreatedAt                time.Time                    `json:"createdAt"`
	UpdatedAt                time.Time                    `json:"updatedAt"`
	CreatedByID              uint64                       `json:"createdByID"`
	UpdatedByID              uint64                       `json:"updatedByID"`
	Key                      string                       `json:"key"`
	DefaultBranchID          *uint64                      `json:"defaultBranchID"`
	DefaultBranch            CanvasBranchMiniSerializer   `json:"defaultBranch"`
	Type                     string                       `json:"type"`
	Parent                   *uint64                      `json:"parent"`
	SubCanvasCount           int                          `json:"subCanvasCount"`
	Rank                     int32                        `json:"rank"`

	CollectionID                uint64    `json:"collectionId"`
	StudioID                    uint64    `json:"studioId"`
	DefaultLanguageCanvasRepoID *uint64   `json:"defaultLanguageCanvasRepoId"`
	Language                    *string   `json:"language"`
	IsLanguageCanvas            bool      `json:"isLanguageCanvas"`
	AutoTranslated              bool      `json:"autoTranslated"`
	IsArchived                  bool      `json:"isArchived"`
	ArchivedAt                  time.Time `json:"archivedAt"`
	ArchivedByID                uint64    `json:"archivedById"`
	HasPublicCanvas             bool      `json:"hasPublicCanvas"`
	MergeRequestCount           int64     `json:"mergeRequestCount"`
}

type CollectionRootNavSerializer struct {
	Id                      uint64                        `json:"id"`
	Name                    string                        `json:"name"`
	Position                uint                          `json:"position"`
	Icon                    string                        `json:"icon"`
	StudioID                uint64                        `json:"studioID"`
	PublicAccess            string                        `json:"publicAccess"`
	ComputedRootCanvasCount int                           `json:"computedRootCanvasCount"`
	ComputedAllCanvasCount  int                           `json:"computedAllCanvasCount"`
	Repos                   []CanvasRepoDefaultSerializer `json:"repos"`
	Type                    string                        `json:"type"`
	HasBranch               bool                          `json:"hasBranch"`
	Rank                    int32                         `json:"rank"`
	CreatedByID             uint64                        `json:"createdById"`
	UpdatedByID             uint64                        `json:"updatedById"`
	IsArchived              bool                          `json:"isArchived"`
	ArchivedAt              time.Time                     `json:"archivedAt"`
	ArchivedByID            uint64                        `json:"archivedById"`
	Permission              string                        `json:"permission"`
	HasPublicCanvas         bool                          `json:"hasPublicCanvas"`
}

func GetSerializedBranches(repoID uint64) []CanvasBranchMiniSerializer {
	var cb []CanvasBranchMiniSerializer
	branches, _ := App.Repo.GetBranchesSimple(map[string]interface{}{"canvas_repository_id": repoID, "is_archived": false, "is_rough_branch": false})
	for _, branch := range *branches {
		// Todo: Calc MR count

		cb = append(cb, CanvasBranchMiniSerializer{
			MergeRequestCount:  0,
			ID:                 branch.ID,
			Name:               branch.Name,
			UUID:               branch.UUID.String(),
			Key:                branch.Key,
			PublicAccess:       branch.PublicAccess,
			CanvasRepositoryID: branch.CanvasRepositoryID,
			IsDraft:            branch.IsDraft,
			IsMerged:           branch.IsMerged,
			IsDefault:          branch.IsDefault,
			IsRoughBranch:      branch.IsRoughBranch,
			//RoughFromBranchID:                 *branch.RoughFromBranchID,
			//RoughBranchCreatorID:              *branch.RoughBranchCreatorID,
			//FromBranchID:                      *branch.FromBranchID,
			CreatedFromCommitID:               branch.CreatedFromCommitID,
			Committed:                         branch.Committed,
			LastSyncedAllAttributionsCommitID: branch.LastSyncedAllAttributionsCommitID,
			CreatedByID:                       branch.CreatedByID,
			UpdatedByID:                       branch.UpdatedByID,
			IsArchived:                        branch.IsArchived,
			ArchivedAt:                        branch.ArchivedAt,
			//ArchivedByID:                      *branch.ArchivedByID,
			Type: "BRANCH",
		})
	}
	return cb
}

func RootSerialized(cb *models.CanvasBranch, userID uint64, public string) *[]CollectionRootNavSerializer {
	collectionsSerialSlice := []CollectionRootNavSerializer{}
	collectionPermissionGroup := models.PGCollectionNoneSysName
	var collectionPermissionList map[uint64]string
	var canvasRepoPermissionList map[uint64]map[uint64]string

	collections, _ := App.Repo.GetCollections(map[string]interface{}{"studio_id": cb.CanvasRepository.Collection.StudioID, "is_archived": false})
	if userID != 0 {
		collectionPermissionList, _ = permissions.App.Service.CalculateCollectionPermissions(userID, cb.CanvasRepository.Collection.StudioID)
	}

	for _, collection := range *collections {
		if public == "true" && (collection.PublicAccess == models.PRIVATE && collection.HasPublicCanvas == false) {
			continue
		}
		collectionPermissionGroup = collectionPermissionList[collection.ID]
		if collectionPermissionGroup == "" || collectionPermissionGroup == models.PGCollectionNoneSysName {
			if collection.HasPublicCanvas {
				collectionPermissionGroup = models.PGCollectionNoneSysName
			} else if collection.PublicAccess == models.EDIT {
				collectionPermissionGroup = models.PGCollectionEditSysName
			} else if collection.PublicAccess == models.COMMENT {
				collectionPermissionGroup = models.PGCollectionCommentSysName
			} else if collection.PublicAccess == models.VIEW {
				collectionPermissionGroup = models.PGCollectionViewSysName
			} else {
				continue
			}
		}

		var HasBranch bool
		HasBranch = false
		var repoSlice []CanvasRepoDefaultSerializer
		if collection.ID == cb.CanvasRepository.Collection.ID {
			HasBranch = true
			canvasRepos, _ := canvasrepo.App.Repo.GetCanvasRepos(map[string]interface{}{"collection_id": cb.CanvasRepository.Collection.ID, "is_archived": false, "is_processing": false})
			for _, repo := range *canvasRepos {
				if public == "true" && (repo.DefaultBranch.PublicAccess == models.PRIVATE && repo.HasPublicCanvas == false) {
					continue
				}
				if repo.DefaultBranchID == nil {
					continue
				}
				canvasRepoPermissionGroup := models.PGCanvasNoneSysName
				if userID != 0 {
					if repo.ParentCanvasRepositoryID != nil {
						canvasRepoPermissionList, _ = permissions.App.Service.CalculateSubCanvasRepoPermissions(userID, cb.CanvasRepository.Collection.StudioID, collection.ID, *repo.ParentCanvasRepositoryID)
						canvasRepoPermissionGroup = canvasRepoPermissionList[repo.ID][*repo.DefaultBranchID]
					} else {
						canvasRepoPermissionList, _ = permissions.App.Service.CalculateCanvasRepoPermissions(userID, cb.CanvasRepository.Collection.StudioID, collection.ID)
						canvasRepoPermissionGroup = canvasRepoPermissionList[repo.ID][*repo.DefaultBranchID]
					}
				}
				if canvasRepoPermissionGroup == "" || canvasRepoPermissionGroup == models.PGCanvasNoneSysName {
					if repo.HasPublicCanvas {
						canvasRepoPermissionGroup = models.PGCanvasViewMetadataSysName
					} else if repo.DefaultBranch.PublicAccess == models.EDIT && repo.IsPublished {
						canvasRepoPermissionGroup = models.PGCanvasEditSysName
					} else if repo.DefaultBranch.PublicAccess == models.COMMENT && repo.IsPublished {
						canvasRepoPermissionGroup = models.PGCanvasCommentSysName
					} else if repo.DefaultBranch.PublicAccess == models.VIEW && repo.IsPublished {
						canvasRepoPermissionGroup = models.PGCanvasViewSysName
					} else {
						continue
					}
				}

				var serialBranches []CanvasBranchMiniSerializer
				var isRepo bool

				isRepo = false
				if cb.CanvasRepository.ID == repo.ID {
					isRepo = true
					serialBranches = GetSerializedBranches(cb.CanvasRepository.ID)
				}
				repoSlice = append(repoSlice, CanvasRepoDefaultSerializer{
					MergeRequestCount:        App.Repo.MergeRequestCount(repo.ID),
					ID:                       repo.ID,
					UUID:                     repo.UUID.String(),
					Name:                     repo.Name,
					Position:                 repo.Position,
					Icon:                     repo.Icon,
					IsPublished:              repo.IsPublished,
					ParentCanvasRepositoryID: repo.ParentCanvasRepositoryID,
					Key:                      repo.Key,
					DefaultBranchID:          repo.DefaultBranchID,
					Rank:                     repo.Rank,
					DefaultBranch: CanvasBranchMiniSerializer{
						ID:                 repo.DefaultBranch.ID,
						Name:               repo.DefaultBranch.Name,
						UUID:               repo.DefaultBranch.UUID.String(),
						Key:                repo.DefaultBranch.Key,
						PublicAccess:       repo.DefaultBranch.PublicAccess,
						Type:               "BRANCH",
						CanvasRepositoryID: repo.DefaultBranch.CanvasRepositoryID,
						IsDraft:            repo.DefaultBranch.IsDraft,
						IsMerged:           repo.DefaultBranch.IsMerged,
						IsDefault:          repo.DefaultBranch.IsDefault,
						IsRoughBranch:      repo.DefaultBranch.IsRoughBranch,
						//RoughFromBranchID:                 *repo.DefaultBranch.RoughFromBranchID,
						//RoughBranchCreatorID:              *repo.DefaultBranch.RoughBranchCreatorID,
						//FromBranchID:                      *repo.DefaultBranch.FromBranchID,
						CreatedFromCommitID:               repo.DefaultBranch.CreatedFromCommitID,
						Committed:                         repo.DefaultBranch.Committed,
						LastSyncedAllAttributionsCommitID: repo.DefaultBranch.LastSyncedAllAttributionsCommitID,
						CreatedByID:                       repo.DefaultBranch.CreatedByID,
						UpdatedByID:                       repo.DefaultBranch.UpdatedByID,
						IsArchived:                        repo.DefaultBranch.IsArchived,
						ArchivedAt:                        repo.DefaultBranch.ArchivedAt,
						//ArchivedByID:                      *repo.DefaultBranch.ArchivedByID,
						Permission: canvasRepoPermissionGroup,
					},
					Type:           "REPO",
					SubCanvasCount: repo.SubCanvasCount,
					CreatedAt:      repo.CreatedAt,
					UpdatedAt:      repo.UpdatedAt,
					CreatedByID:    repo.CreatedByID,
					UpdatedByID:    repo.UpdatedByID,
					HasBranch:      isRepo,
					Branches:       serialBranches,

					CollectionID:                repo.CollectionID,
					StudioID:                    repo.StudioID,
					DefaultLanguageCanvasRepoID: repo.DefaultLanguageCanvasRepoID,
					Language:                    repo.Language,
					IsLanguageCanvas:            repo.IsLanguageCanvas,
					AutoTranslated:              repo.AutoTranslated,
					IsArchived:                  repo.IsArchived,
					ArchivedAt:                  repo.ArchivedAt,
					HasPublicCanvas:             repo.HasPublicCanvas,
					//	ArchivedByID:     *repo.ArchivedByID,
				})
			}
		}
		collectionRootSerializerData := CollectionRootNavSerializer{
			Id:                      collection.ID,
			Name:                    collection.Name,
			Position:                collection.Position,
			Icon:                    collection.Icon,
			StudioID:                collection.StudioID,
			PublicAccess:            collection.PublicAccess,
			ComputedRootCanvasCount: collection.ComputedRootCanvasCount,
			ComputedAllCanvasCount:  collection.ComputedAllCanvasCount,
			Rank:                    collection.Rank,
			Type:                    "COLLECTION",
			Repos:                   repoSlice,
			HasBranch:               HasBranch,
			CreatedByID:             collection.CreatedByID,
			UpdatedByID:             collection.UpdatedByID,
			IsArchived:              collection.IsArchived,
			ArchivedAt:              collection.ArchivedAt,
			Permission:              collectionPermissionGroup,
			HasPublicCanvas:         collection.HasPublicCanvas,
		}
		if collection.ArchivedByID != nil {
			collectionRootSerializerData.ArchivedByID = *collection.ArchivedByID
		}
		collectionsSerialSlice = append(collectionsSerialSlice, collectionRootSerializerData)

	}

	return &collectionsSerialSlice
}

func RootSerializedByStudioID(studioID uint64, userID uint64) *[]CollectionRootNavSerializer {
	collectionsSerialSlice := []CollectionRootNavSerializer{}
	collectionPermissionGroup := models.PGCollectionNoneSysName
	var collectionPermissionList map[uint64]string
	var canvasRepoPermissionList map[uint64]map[uint64]string

	collections, _ := App.Repo.GetCollections(map[string]interface{}{"studio_id": studioID, "is_archived": false})
	if userID != 0 {
		collectionPermissionList, _ = permissions.App.Service.CalculateCollectionPermissions(userID, studioID)
	}

	for _, collection := range *collections {
		collectionPermissionGroup = collectionPermissionList[collection.ID]
		if collectionPermissionGroup == "" || collectionPermissionGroup == models.PGCollectionNoneSysName {
			if collection.HasPublicCanvas {
				collectionPermissionGroup = models.PGCollectionViewMetadataSysName
			} else if collection.PublicAccess == models.EDIT {
				collectionPermissionGroup = models.PGCollectionEditSysName
			} else if collection.PublicAccess == models.COMMENT {
				collectionPermissionGroup = models.PGCollectionCommentSysName
			} else if collection.PublicAccess == models.VIEW {
				collectionPermissionGroup = models.PGCollectionViewSysName
			} else {
				continue
			}
		}

		var HasBranch bool
		HasBranch = false
		var repoSlice []CanvasRepoDefaultSerializer

		HasBranch = true
		canvasRepos, _ := canvasrepo.App.Repo.GetCanvasRepos(map[string]interface{}{"collection_id": collection.ID, "is_archived": false, "is_processing": false})
		fmt.Println("Collectionid and canvasrepocount", collection.ID, len(*canvasRepos))
		for _, repo := range *canvasRepos {
			if repo.DefaultBranchID == nil {
				continue
			}
			canvasRepoPermissionGroup := models.PGCanvasNoneSysName
			if userID != 0 {
				if repo.ParentCanvasRepositoryID != nil {
					canvasRepoPermissionList, _ = permissions.App.Service.CalculateSubCanvasRepoPermissions(userID, studioID, collection.ID, *repo.ParentCanvasRepositoryID)
					canvasRepoPermissionGroup = canvasRepoPermissionList[repo.ID][*repo.DefaultBranchID]
				} else {
					canvasRepoPermissionList, _ = permissions.App.Service.CalculateCanvasRepoPermissions(userID, studioID, collection.ID)
					canvasRepoPermissionGroup = canvasRepoPermissionList[repo.ID][*repo.DefaultBranchID]
				}
			}
			if canvasRepoPermissionGroup == "" || canvasRepoPermissionGroup == models.PGCanvasNoneSysName {
				if repo.HasPublicCanvas {
					canvasRepoPermissionGroup = models.PGCanvasViewMetadataSysName
				} else if repo.DefaultBranch.PublicAccess == models.EDIT && repo.IsPublished {
					canvasRepoPermissionGroup = models.PGCanvasEditSysName
				} else if repo.DefaultBranch.PublicAccess == models.COMMENT && repo.IsPublished {
					canvasRepoPermissionGroup = models.PGCanvasCommentSysName
				} else if repo.DefaultBranch.PublicAccess == models.VIEW && repo.IsPublished {
					canvasRepoPermissionGroup = models.PGCanvasViewSysName
				} else {
					continue
				}
			}
			var serialBranches []CanvasBranchMiniSerializer
			repoSlice = append(repoSlice, CanvasRepoDefaultSerializer{
				MergeRequestCount:        App.Repo.MergeRequestCount(repo.ID),
				ID:                       repo.ID,
				UUID:                     repo.UUID.String(),
				Name:                     repo.Name,
				Position:                 repo.Position,
				Icon:                     repo.Icon,
				IsPublished:              repo.IsPublished,
				ParentCanvasRepositoryID: repo.ParentCanvasRepositoryID,
				Key:                      repo.Key,
				DefaultBranchID:          repo.DefaultBranchID,
				Rank:                     repo.Rank,
				DefaultBranch: CanvasBranchMiniSerializer{
					ID:                                repo.DefaultBranch.ID,
					Name:                              repo.DefaultBranch.Name,
					UUID:                              repo.DefaultBranch.UUID.String(),
					Key:                               repo.DefaultBranch.Key,
					PublicAccess:                      repo.DefaultBranch.PublicAccess,
					Type:                              "BRANCH",
					CanvasRepositoryID:                repo.DefaultBranch.CanvasRepositoryID,
					IsDraft:                           repo.DefaultBranch.IsDraft,
					IsMerged:                          repo.DefaultBranch.IsMerged,
					IsDefault:                         repo.DefaultBranch.IsDefault,
					IsRoughBranch:                     repo.DefaultBranch.IsRoughBranch,
					CreatedFromCommitID:               repo.DefaultBranch.CreatedFromCommitID,
					Committed:                         repo.DefaultBranch.Committed,
					LastSyncedAllAttributionsCommitID: repo.DefaultBranch.LastSyncedAllAttributionsCommitID,
					CreatedByID:                       repo.DefaultBranch.CreatedByID,
					UpdatedByID:                       repo.DefaultBranch.UpdatedByID,
					IsArchived:                        repo.DefaultBranch.IsArchived,
					ArchivedAt:                        repo.DefaultBranch.ArchivedAt,
					Permission:                        canvasRepoPermissionGroup,
				},
				Type:                        "REPO",
				SubCanvasCount:              repo.SubCanvasCount,
				CreatedAt:                   repo.CreatedAt,
				UpdatedAt:                   repo.UpdatedAt,
				CreatedByID:                 repo.CreatedByID,
				UpdatedByID:                 repo.UpdatedByID,
				Branches:                    serialBranches,
				CollectionID:                repo.CollectionID,
				StudioID:                    repo.StudioID,
				DefaultLanguageCanvasRepoID: repo.DefaultLanguageCanvasRepoID,
				Language:                    repo.Language,
				IsLanguageCanvas:            repo.IsLanguageCanvas,
				AutoTranslated:              repo.AutoTranslated,
				IsArchived:                  repo.IsArchived,
				ArchivedAt:                  repo.ArchivedAt,
				HasPublicCanvas:             repo.HasPublicCanvas,
			})
		}

		collectionRootSerializerData := CollectionRootNavSerializer{
			Id:                      collection.ID,
			Name:                    collection.Name,
			Position:                collection.Position,
			Icon:                    collection.Icon,
			StudioID:                collection.StudioID,
			PublicAccess:            collection.PublicAccess,
			ComputedRootCanvasCount: collection.ComputedRootCanvasCount,
			ComputedAllCanvasCount:  collection.ComputedAllCanvasCount,
			Rank:                    collection.Rank,
			Type:                    "COLLECTION",
			Repos:                   repoSlice,
			HasBranch:               HasBranch,
			CreatedByID:             collection.CreatedByID,
			UpdatedByID:             collection.UpdatedByID,
			IsArchived:              collection.IsArchived,
			ArchivedAt:              collection.ArchivedAt,
			Permission:              collectionPermissionGroup,
			HasPublicCanvas:         collection.HasPublicCanvas,
		}
		if collection.ArchivedByID != nil {
			collectionRootSerializerData.ArchivedByID = *collection.ArchivedByID
		}
		collectionsSerialSlice = append(collectionsSerialSlice, collectionRootSerializerData)

	}
	return &collectionsSerialSlice
}

// Todo
func NodeSerialized(cb *models.CanvasBranch, userID uint64) *[]CanvasRepoDefaultSerializer {

	//collectionsSerialSlice := []CollectionRootNavSerializer{}
	//collectionPermissionGroup := models.PGCollectionNoneSysName
	//var collectionPermissionList map[uint64]string
	//var canvasRepoPermissionList map[uint64]map[uint64]string

	//collections, _ := App.Repo.GetCollections(map[string]interface{}{"studio_id": cb.CanvasRepository.Collection.StudioID, "is_archived": false})
	if userID != 0 {
		//collectionPermissionList, _ = permissions.App.Service.CalculateCollectionPermissions(userID, cb.CanvasRepository.Collection.StudioID)
	}
	//
	//for _, collection := range *collections {
	//	collectionPermissionGroup = collectionPermissionList[collection.ID]
	//	if collectionPermissionGroup == "" || collectionPermissionGroup == models.PGCollectionNoneSysName {
	//		if collection.HasPublicCanvas {
	//			collectionPermissionGroup = models.PGCollectionViewMetadataSysName
	//		} else if collection.PublicAccess == models.EDIT {
	//			collectionPermissionGroup = models.PGCollectionEditSysName
	//		} else if collection.PublicAccess == models.COMMENT {
	//			collectionPermissionGroup = models.PGCollectionCommentSysName
	//		} else if collection.PublicAccess == models.VIEW {
	//			collectionPermissionGroup = models.PGCollectionViewSysName
	//		} else {
	//			continue
	//		}
	//	}
	//
	//	var HasBranch bool
	//	HasBranch = false
	//	var repoSlice []CanvasRepoDefaultSerializer
	//	if collection.ID == cb.CanvasRepository.Collection.ID {
	//		HasBranch = true
	//		canvasRepos, _ := canvasrepo.App.Repo.GetCanvasRepos(map[string]interface{}{"collection_id": cb.CanvasRepository.Collection.ID, "is_archived": false})
	//		for _, repo := range *canvasRepos {
	//			if repo.DefaultBranchID == nil {
	//				continue
	//			}
	//			canvasRepoPermissionGroup := models.PGCanvasNoneSysName
	//			if userID != 0 {
	//				if repo.ParentCanvasRepositoryID != nil {
	//					canvasRepoPermissionList, _ = permissions.App.Service.CalculateSubCanvasRepoPermissions(userID, cb.CanvasRepository.Collection.StudioID, collection.ID, *repo.ParentCanvasRepositoryID)
	//					canvasRepoPermissionGroup = canvasRepoPermissionList[repo.ID][*repo.DefaultBranchID]
	//				} else {
	//					canvasRepoPermissionList, _ = permissions.App.Service.CalculateCanvasRepoPermissions(userID, cb.CanvasRepository.Collection.StudioID, collection.ID)
	//					canvasRepoPermissionGroup = canvasRepoPermissionList[repo.ID][*repo.DefaultBranchID]
	//				}
	//			}
	//			if canvasRepoPermissionGroup == "" || canvasRepoPermissionGroup == models.PGCanvasNoneSysName {
	//				if repo.HasPublicCanvas {
	//					canvasRepoPermissionGroup = models.PGCanvasViewMetadataSysName
	//				} else if repo.DefaultBranch.PublicAccess == models.EDIT && repo.IsPublished {
	//					canvasRepoPermissionGroup = models.PGCanvasEditSysName
	//				} else if repo.DefaultBranch.PublicAccess == models.COMMENT && repo.IsPublished {
	//					canvasRepoPermissionGroup = models.PGCanvasCommentSysName
	//				} else if repo.DefaultBranch.PublicAccess == models.VIEW && repo.IsPublished {
	//					canvasRepoPermissionGroup = models.PGCanvasViewSysName
	//				} else {
	//					continue
	//				}
	//			}
	//
	//			var serialBranches []CanvasBranchMiniSerializer
	//			var isRepo bool
	//
	//			isRepo = false
	//			if cb.CanvasRepository.ID == repo.ID {
	//				isRepo = true
	//				serialBranches = GetSerializedBranches(cb.CanvasRepository.ID)
	//			}
	//			repoSlice = append(repoSlice, CanvasRepoDefaultSerializer{
	//				MergeRequestCount:        App.Repo.MergeRequestCount(repo.ID),
	//				ID:                       repo.ID,
	//				UUID:                     repo.UUID.String(),
	//				Name:                     repo.Name,
	//				Position:                 repo.Position,
	//				Icon:                     repo.Icon,
	//				IsPublished:              repo.IsPublished,
	//				ParentCanvasRepositoryID: repo.ParentCanvasRepositoryID,
	//				Key:                      repo.Key,
	//				DefaultBranchID:          repo.DefaultBranchID,
	//				Rank:                     repo.Rank,
	//				DefaultBranch: CanvasBranchMiniSerializer{
	//					ID:                 repo.DefaultBranch.ID,
	//					Name:               repo.DefaultBranch.Name,
	//					UUID:               repo.DefaultBranch.UUID.String(),
	//					Key:                repo.DefaultBranch.Key,
	//					PublicAccess:       repo.DefaultBranch.PublicAccess,
	//					Type:               "BRANCH",
	//					CanvasRepositoryID: repo.DefaultBranch.CanvasRepositoryID,
	//					IsDraft:            repo.DefaultBranch.IsDraft,
	//					IsMerged:           repo.DefaultBranch.IsMerged,
	//					IsDefault:          repo.DefaultBranch.IsDefault,
	//					IsRoughBranch:      repo.DefaultBranch.IsRoughBranch,
	//					//RoughFromBranchID:                 *repo.DefaultBranch.RoughFromBranchID,
	//					//RoughBranchCreatorID:              *repo.DefaultBranch.RoughBranchCreatorID,
	//					//FromBranchID:                      *repo.DefaultBranch.FromBranchID,
	//					CreatedFromCommitID:               repo.DefaultBranch.CreatedFromCommitID,
	//					Committed:                         repo.DefaultBranch.Committed,
	//					LastSyncedAllAttributionsCommitID: repo.DefaultBranch.LastSyncedAllAttributionsCommitID,
	//					CreatedByID:                       repo.DefaultBranch.CreatedByID,
	//					UpdatedByID:                       repo.DefaultBranch.UpdatedByID,
	//					IsArchived:                        repo.DefaultBranch.IsArchived,
	//					ArchivedAt:                        repo.DefaultBranch.ArchivedAt,
	//					//ArchivedByID:                      *repo.DefaultBranch.ArchivedByID,
	//					Permission: canvasRepoPermissionGroup,
	//				},
	//				Type:           "REPO",
	//				SubCanvasCount: repo.SubCanvasCount,
	//				CreatedAt:      repo.CreatedAt,
	//				UpdatedAt:      repo.UpdatedAt,
	//				CreatedByID:    repo.CreatedByID,
	//				UpdatedByID:    repo.UpdatedByID,
	//				HasBranch:      isRepo,
	//				Branches:       serialBranches,
	//
	//				CollectionID:                repo.CollectionID,
	//				StudioID:                    repo.StudioID,
	//				DefaultLanguageCanvasRepoID: repo.DefaultLanguageCanvasRepoID,
	//				Language:                    repo.Language,
	//				IsLanguageCanvas:            repo.IsLanguageCanvas,
	//				AutoTranslated:              repo.AutoTranslated,
	//				IsArchived:                  repo.IsArchived,
	//				ArchivedAt:                  repo.ArchivedAt,
	//				HasPublicCanvas:             repo.HasPublicCanvas,
	//				//	ArchivedByID:     *repo.ArchivedByID,
	//			})
	//		}
	//	}
	//	collectionRootSerializerData := CollectionRootNavSerializer{
	//		Id:                      collection.ID,
	//		Name:                    collection.Name,
	//		Position:                collection.Position,
	//		Icon:                    collection.Icon,
	//		StudioID:                collection.StudioID,
	//		PublicAccess:            collection.PublicAccess,
	//		ComputedRootCanvasCount: collection.ComputedRootCanvasCount,
	//		ComputedAllCanvasCount:  collection.ComputedAllCanvasCount,
	//		Rank:                    collection.Rank,
	//		Type:                    "COLLECTION",
	//		Repos:                   repoSlice,
	//		HasBranch:               HasBranch,
	//		CreatedByID:             collection.CreatedByID,
	//		UpdatedByID:             collection.UpdatedByID,
	//		IsArchived:              collection.IsArchived,
	//		ArchivedAt:              collection.ArchivedAt,
	//		Permission:              collectionPermissionGroup,
	//		HasPublicCanvas:         collection.HasPublicCanvas,
	//	}
	//	if collection.ArchivedByID != nil {
	//		collectionRootSerializerData.ArchivedByID = *collection.ArchivedByID
	//	}
	//	collectionsSerialSlice = append(collectionsSerialSlice, collectionRootSerializerData)
	//
	//}

	return &[]CanvasRepoDefaultSerializer{}
}

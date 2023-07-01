package canvasbranch

import (
	"github.com/google/uuid"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"time"
)

// We need to return a Map of {collections, repos an d branches}

type fullRow struct {
	ID                          uint64
	UUID                        uuid.UUID
	CreatedAt                   time.Time
	UpdatedAt                   time.Time
	CollectionID                uint64
	StudioID                    uint64
	Name                        string
	Position                    uint
	Icon                        string
	PublicAccess                string
	IsPublished                 bool
	DefaultBranchID             uint64
	ParentCanvasRepositoryID    *uint64
	CreatedByID                 uint64
	UpdatedByID                 uint64
	IsArchived                  bool
	ArchivedAt                  time.Time
	ArchivedByID                *uint64
	DefaultLanguageCanvasRepoID *uint64
	Language                    string
	IsLanguageCanvas            bool
	AutoTranslated              bool
	Key                         string
	SubCanvasCount              uint
	ID2                         uint64
	UUID2                       uuid.UUID
	Name2                       string
	Position2                   uint64
	PublicAccess2               string
	Icon2                       string
	HasPublicAccess             bool
	ComputedRootCanvasCount2    int
}

func (r canvasBranchRepo) QueryDB(q string, studioID uint64) *[]fullRow {
	var finalResult []fullRow
	var onerow fullRow
	q = q + "%"
	query := "select canvas_repositories.*, " +
		"collections.id as id2, " +
		"collections.uuid as uuid2, " +
		"collections.name as name2, " +
		"collections.icon as icon2, " +
		"collections.position as position2, " +
		"collections.public_access as public_access2, " +
		"collections.computed_root_canvas_count as computed_root_canvas_count2 " +
		"from canvas_repositories left join collections " +
		"ON canvas_repositories.collection_id = collections.id " +
		"where canvas_repositories.studio_id = ? \n" +
		"and canvas_repositories.is_archived = false \n" +
		"and collections.is_archived = false \n" +
		"and ( collections.name ILIKE ?" +
		"or canvas_repositories.name ILIKE ?) order by collections.position, canvas_repositories.position;"
	rows, _ := r.db.Raw(query, studioID, "%"+q+"%", "%"+q+"%").Rows()
	defer rows.Close()
	for rows.Next() {
		r.db.ScanRows(rows, &onerow)
		finalResult = append(finalResult, onerow)
	}
	return &finalResult
}

func (s canvasBranchService) ProcessSearchDump(records *[]fullRow, userID uint64, studioID uint64, isPublic string) map[string][]map[string]interface{} {
	var CollectionsMap = []map[string]interface{}{}
	var RepoMap = []map[string]interface{}{}
	var BranchMap = []map[string]interface{}{}
	collectionIds := []uint64{}
	repoIDs := []uint64{}
	collectionCheckMap := make(map[uint64]bool)
	repoIDMap := make(map[uint64]bool)

	collectionPermList, _ := permissions.App.Service.CalculateCollectionPermissions(userID, studioID)
	repoPermList := map[uint64]map[uint64]string{}

	for _, v := range *records {
		collectionIds = append(collectionIds, v.ID2)
		repoIDs = append(repoIDs, v.ID)
		repoIDMap[v.ID] = true
	}

	for _, v := range *records {
		var permList map[uint64]map[uint64]string
		if v.ParentCanvasRepositoryID != nil {
			permList, _ = permissions.App.Service.CalculateSubCanvasRepoPermissions(userID, studioID, v.CollectionID, *v.ParentCanvasRepositoryID)
		} else {
			permList, _ = permissions.App.Service.CalculateCanvasRepoPermissions(userID, studioID, v.CollectionID)
		}
		for key, value := range permList {
			repoPermList[key] = value
		}
		hasPermission, _ := permissions.App.Service.CanUserDoThisOnBranch(userID, v.DefaultBranchID, permissiongroup.CANVAS_BRANCH_VIEW_METADATA)
		if !hasPermission {
			continue
		}
		RepoMap = append(RepoMap, map[string]interface{}{
			"id":                       v.ID,
			"uuid":                     v.UUID,
			"name":                     v.Name,
			"position":                 v.Position,
			"icon":                     v.Icon,
			"key":                      v.Key,
			"collectionID":             v.CollectionID,
			"isPublished":              v.IsPublished,
			"defaultBranchID":          v.DefaultBranchID,
			"parentCanvasRepositoryID": v.ParentCanvasRepositoryID,
			"subCanvasCount":           v.SubCanvasCount,
			"createdByID":              v.CreatedByID,
			"hasPublicAccess":          v.HasPublicAccess,
			"match":                    true,
		})
		// Logic to add parent canvas which are not matched but for the tree to build in FrontEnd
		if v.ParentCanvasRepositoryID != nil && !repoIDMap[*v.ParentCanvasRepositoryID] {
			var parentRecords []map[string]interface{}
			parentRecords, repoIDMap = s.GetParentRecords(*v.ParentCanvasRepositoryID, repoIDMap, []map[string]interface{}{})
			for _, parentRepo := range parentRecords {
				if parentRepo["parentCanvasRepositoryID"].(*uint64) != nil {
					permList, _ = permissions.App.Service.CalculateSubCanvasRepoPermissions(userID, studioID, v.CollectionID, *parentRepo["parentCanvasRepositoryID"].(*uint64))
				} else {
					permList, _ = permissions.App.Service.CalculateCanvasRepoPermissions(userID, studioID, v.CollectionID)
				}
				for key, value := range permList {
					repoPermList[key] = value
				}
				repoIDs = append(repoIDs, parentRepo["id"].(uint64))
			}
			RepoMap = append(RepoMap, parentRecords...)
		}
	}

	for _, v := range *records {
		permission := collectionPermList[v.ID2]
		if permission == "" {
			permission = models.PGCollectionNoneSysName
		}
		_, isPresent := collectionCheckMap[v.ID2]
		if isPresent {
			continue
		}
		collectionCheckMap[v.ID2] = true
		CollectionsMap = append(CollectionsMap, map[string]interface{}{
			"id":                         v.ID2,
			"uuid":                       v.UUID2,
			"name":                       v.Name2,
			"position":                   v.Position2,
			"icon":                       v.Icon2,
			"permission":                 permission,
			"computed_root_canvas_count": v.ComputedRootCanvasCount2,
		})
	}

	branches, _ := App.Repo.GetBranchesNoPreload(repoIDs)
	canvasRepoBranchMap := make(map[uint64]map[string]interface{})
	for _, v := range *branches {
		permission := repoPermList[v.CanvasRepositoryID][v.ID]
		if permission == "" {
			permission = models.PGCanvasNoneSysName
		}
		branch := map[string]interface{}{
			"id":                 v.ID,
			"uuid":               v.UUID,
			"name":               v.Name,
			"key":                v.Key,
			"isDraft":            v.IsDraft,
			"isDefault":          v.IsDefault,
			"canvasRepositoryID": v.CanvasRepositoryID,
			"isRoughBranch":      v.IsRoughBranch,
			"publicAccess":       v.PublicAccess,
			"createdByID":        v.CreatedByID,
			"permission":         permission,
		}
		if v.IsDefault {
			BranchMap = append(BranchMap, branch)
			canvasRepoBranchMap[v.CanvasRepositoryID] = branch
		}
	}

	finalRepoMap := []map[string]interface{}{}
	finalCollectionHash := map[uint64]bool{}
	for _, repo := range RepoMap {
		defaultBranch := canvasRepoBranchMap[repo["id"].(uint64)]
		if isPublic == "true" {
			if defaultBranch["publicAccess"].(string) != models.PRIVATE || repo["hasPublicAccess"].(bool) {
				repo["defaultBranch"] = canvasRepoBranchMap[repo["id"].(uint64)]
				finalCollectionHash[repo["collectionID"].(uint64)] = true
				finalRepoMap = append(finalRepoMap, repo)
			}
		} else {
			if defaultBranch["permission"].(string) == models.PGCanvasNoneSysName && defaultBranch["publicAccess"].(string) == models.PRIVATE && !repo["hasPublicAccess"].(bool) {
				continue
			}
			repo["defaultBranch"] = canvasRepoBranchMap[repo["id"].(uint64)]
			finalCollectionHash[repo["collectionID"].(uint64)] = true
			finalRepoMap = append(finalRepoMap, repo)
		}
	}

	finalCollectionMap := []map[string]interface{}{}
	for _, col := range CollectionsMap {
		if finalCollectionHash[col["id"].(uint64)] {
			finalCollectionMap = append(finalCollectionMap, col)
		}
	}

	var Responded = map[string][]map[string]interface{}{
		"collections": finalCollectionMap,
		"repos":       finalRepoMap,
		"branches":    BranchMap,
	}

	return Responded
}

func (s canvasBranchService) GetParentRecords(canvasID uint64, repoIDMap map[uint64]bool, canvasParentRecords []map[string]interface{}) ([]map[string]interface{}, map[uint64]bool) {
	canvasRepo, _ := App.Repo.GetCanvasRepoInstance(map[string]interface{}{"id": canvasID})
	if canvasRepo.ParentCanvasRepositoryID != nil && !repoIDMap[canvasRepo.ID] {
		canvasParentRecords, repoIDMap = s.GetParentRecords(*canvasRepo.ParentCanvasRepositoryID, repoIDMap, canvasParentRecords)
	}
	if repoIDMap[canvasRepo.ID] {
		return canvasParentRecords, repoIDMap
	}
	recordMap := map[string]interface{}{
		"id":                       canvasRepo.ID,
		"uuid":                     canvasRepo.UUID,
		"name":                     canvasRepo.Name,
		"position":                 canvasRepo.Position,
		"icon":                     canvasRepo.Icon,
		"key":                      canvasRepo.Key,
		"collectionID":             canvasRepo.CollectionID,
		"isPublished":              canvasRepo.IsPublished,
		"defaultBranchID":          canvasRepo.DefaultBranchID,
		"parentCanvasRepositoryID": canvasRepo.ParentCanvasRepositoryID,
		"subCanvasCount":           canvasRepo.SubCanvasCount,
		"createdByID":              canvasRepo.CreatedByID,
		"match":                    false,
		"hasPublicAccess":          canvasRepo.HasPublicCanvas,
	}
	repoIDMap[canvasRepo.ID] = true
	canvasParentRecords = append(canvasParentRecords, recordMap)
	return canvasParentRecords, repoIDMap
}

/*
var onerow rorow
	q = q + "%"
	query := "select * from canvas_repositories " +
		"left join collections " +
		"ON canvas_repositories.collection_id = collections.id " +
		"where collections.name ILIKE ? " +
		"\n " +
		"union all" +
		"\n " +
		"select * from canvas_repositories left join collections " +
		"ON canvas_repositories.collection_id = collections.id " +
		"where  canvas_repositories.studio_id = ? " +
		"and canvas_repositories.name ILIKE ? " +
		"and canvas_repositories.is_archived is false "
	rows, _ := r.db.Raw(query, q, studioID, q).Rows()
	defer rows.Close()
	for rows.Next() {
		r.db.ScanRows(rows, &onerow)
		fmt.Println("-------")
		fmt.Printf("%+v\n", rows)
		fmt.Println("-------")
		fmt.Printf("%+v\n", onerow)

		os.Exit(42)
		// do something
	}
*/

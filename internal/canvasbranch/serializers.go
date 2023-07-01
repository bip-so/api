package canvasbranch

import (
	"fmt"
	ar "gitlab.com/phonepost/bip-be-platform/internal/accessrequest"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"time"

	"github.com/gosimple/slug"

	"github.com/google/uuid"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

type BranchDefault struct {
	ID                                uint64    `json:"id"`
	UUID                              uuid.UUID `json:"uuid"`
	CanvasRepositoryID                uint64    `json:"canvasRepositoryID"`
	Name                              string    `json:"name"`
	Key                               string    `json:"key"`
	IsDraft                           bool      `json:"isDraft"`
	IsMerged                          bool      `json:"isMerged"`
	IsDefault                         bool      `json:"isDefault"`
	PublicAccess                      string    `json:"publicAccess"`
	IsRoughBranch                     bool      `json:"isRoughBranch"`
	RoughFromBranchID                 *uint64   `json:"roughFromBranchID"`
	RoughBranchCreatorID              *uint64   `json:"roughBranchCreatorID"`
	FromBranchID                      *uint64   `json:"fromBranchID"`
	CreatedFromCommitID               string    `json:"createdFromCommitID"`
	Committed                         bool      `json:"committed"`
	LastSyncedAllAttributionsCommitID string    `json:"lastSyncedAllAttributionsCommitID"`
	CreatedByID                       uint64    `json:"createdByID"`
	UpdatedByID                       uint64    `json:"updatedByID"`
	CreatedAt                         time.Time `json:"createdAt"`
	UpdatedAt                         time.Time `json:"updatedAt"`
	Slug                              string    `json:"slug"`
}

type AccessTokenSerial struct {
	InviteCode      string `json:"inviteCode"`
	PermissionGroup string `json:"permissionGroup"`
	IsActive        bool   `json:"isActive"`
	CreatedByID     uint64 `json:"createdById"`
}

type BranchMeta struct {
	ID                                uint64                             `json:"id"`
	UUID                              uuid.UUID                          `json:"uuid"`
	CanvasRepositoryID                uint64                             `json:"canvasRepositoryID"`
	Name                              string                             `json:"name"`
	Key                               string                             `json:"key"`
	IsDraft                           bool                               `json:"isDraft"`
	IsMerged                          bool                               `json:"isMerged"`
	IsDefault                         bool                               `json:"isDefault"`
	PublicAccess                      string                             `json:"publicAccess"`
	IsRoughBranch                     bool                               `json:"isRoughBranch"`
	RoughFromBranchID                 *uint64                            `json:"roughFromBranchID"`
	RoughBranchCreatorID              *uint64                            `json:"roughBranchCreatorID"`
	FromBranchID                      *uint64                            `json:"fromBranchID"`
	CreatedFromCommitID               string                             `json:"createdFromCommitID"`
	Committed                         bool                               `json:"committed"`
	LastSyncedAllAttributionsCommitID string                             `json:"lastSyncedAllAttributionsCommitID"`
	CreatedByID                       uint64                             `json:"createdByID"`
	UpdatedByID                       uint64                             `json:"updatedByID"`
	CreatedAt                         time.Time                          `json:"createdAt"`
	UpdatedAt                         time.Time                          `json:"updatedAt"`
	MergeRequest                      *MergeRequestSerializer            `json:"mergeRequest"`
	AccessRequests                    *ar.AccessRequestSerializerDefault `json:"accessRequests"`
	ContributorsList                  []GitAttributionView               `json:"contributorsList"`
	BranchAccessTokens                []AccessTokenSerial                `json:"branchAccessTokens"`
	CanvasRepoKey                     string                             `json:"canvasRepoKey"`
	Slug                              string                             `json:"slug"`
}

func SimpleBranchDiffSerializer(model *models.CanvasBranch) BranchDefault {

	view := BranchDefault{
		ID:                                model.ID,
		UUID:                              model.UUID,
		UpdatedByID:                       model.UpdatedByID,
		CreatedAt:                         model.CreatedAt,
		UpdatedAt:                         model.UpdatedAt,
		CanvasRepositoryID:                model.CanvasRepositoryID,
		Name:                              model.Name,
		Key:                               model.Key,
		IsDraft:                           model.IsDraft,
		IsMerged:                          model.IsMerged,
		IsDefault:                         model.IsDefault,
		PublicAccess:                      model.PublicAccess,
		IsRoughBranch:                     model.IsRoughBranch,
		RoughFromBranchID:                 model.RoughFromBranchID,
		RoughBranchCreatorID:              model.RoughBranchCreatorID,
		FromBranchID:                      model.FromBranchID,
		CreatedFromCommitID:               model.CreatedFromCommitID,
		Committed:                         model.Committed,
		LastSyncedAllAttributionsCommitID: model.LastSyncedAllAttributionsCommitID,
		CreatedByID:                       model.CreatedByID,
		Slug:                              slug.Make(model.Name),
	}

	return view

}

func BranchTokensSerialMaker(id uint64) []AccessTokenSerial {
	var tokens []AccessTokenSerial
	tokensInstances, _ := queries.App.BranchAccessTokenQuery.GetAllBranchTokenInstance(map[string]interface{}{"branch_id": id})
	for _, v := range *tokensInstances {
		tokens = append(tokens, AccessTokenSerial{
			InviteCode:      v.InviteCode,
			PermissionGroup: v.PermissionGroup,
			IsActive:        v.IsActive,
			CreatedByID:     v.CreatedByID,
		})
	}

	return tokens
}

func SimpleBranchDefaultSerializer(model *models.CanvasBranch) BranchMeta {
	view := BranchMeta{
		ID:                                model.ID,
		UUID:                              model.UUID,
		UpdatedByID:                       model.UpdatedByID,
		CreatedAt:                         model.CreatedAt,
		UpdatedAt:                         model.UpdatedAt,
		CanvasRepositoryID:                model.CanvasRepositoryID,
		Name:                              model.Name,
		Key:                               model.Key,
		IsDraft:                           model.IsDraft,
		IsMerged:                          model.IsMerged,
		IsDefault:                         model.IsDefault,
		PublicAccess:                      model.PublicAccess,
		IsRoughBranch:                     model.IsRoughBranch,
		RoughFromBranchID:                 model.RoughFromBranchID,
		RoughBranchCreatorID:              model.RoughBranchCreatorID,
		FromBranchID:                      model.FromBranchID,
		CreatedFromCommitID:               model.CreatedFromCommitID,
		Committed:                         model.Committed,
		LastSyncedAllAttributionsCommitID: model.LastSyncedAllAttributionsCommitID,
		CreatedByID:                       model.CreatedByID,
		BranchAccessTokens:                BranchTokensSerialMaker(model.ID),
		Slug:                              slug.Make(model.Name),
	}

	if model.CanvasRepository != nil {
		view.CanvasRepoKey = model.CanvasRepository.Key
	}

	return view

}

func BranchMetaWithMRSerializer(model *models.CanvasBranch, branchID uint64, userID uint64, attributions *[]models.Attribution) BranchMeta {
	serialized := SimpleBranchDefaultSerializer(model)
	if attributions != nil {
		serialized.ContributorsList = *GitAttributionsSerializerData(attributions)
	} else {
		serialized.ContributorsList = []GitAttributionView{}
	}

	if userID == 0 {
		serialized.MergeRequest = nil
		return serialized
	}
	mrInstance, _ := App.Repo.GetMergeRequest(map[string]interface{}{"source_branch_id": model.ID, "destination_branch_id": model.RoughFromBranchID, "status": models.MERGE_REQUEST_OPEN})
	if mrInstance == nil {
		serialized.MergeRequest = nil
		return serialized
	}
	//arInstances, _ := App.Repo.GetAllAcceesRequests(map[string]interface{}{"canvas_branch_id": model.ID})
	//if arInstances == nil {
	//	serialized.AccessRequests = nil
	//	return serialized
	//}

	serialized.MergeRequest = &MergeRequestSerializer{
		Status:        mrInstance.Status,
		ID:            mrInstance.ID,
		CommitMessage: mrInstance.CommitMessage,
		CreatedByID:   mrInstance.CreatedByID,
	}

	//for k,v  in arInstances{}
	serialized.AccessRequests = nil

	return serialized
}

type ReactedCount struct {
	Emoji   string `json:"emoji"`
	Count   int    `json:"count"`
	Reacted *bool  `json:"reacted"`
}

type GetCanvasBranchSerializer struct {
	ID                                uint64  `json:"id"`
	Name                              string  `json:"name"`
	UUID                              string  `json:"uuid"`
	Key                               string  `json:"key"`
	PublicAccess                      string  `json:"publicAccess"`
	IsDraft                           bool    `json:"isDraft"`
	IsDefault                         bool    `json:"isDefault"`
	IsMerged                          bool    `json:"isMerged"`
	IsRoughBranch                     bool    `json:"isRoughBranch"`
	RoughFromBranchID                 *uint64 `json:"roughFromBranchID"`
	RoughBranchCreatorID              *uint64 `json:"roughBranchCreatorID"`
	FromBranchID                      *uint64 `json:"fromBranchID"`
	CreatedFromCommitID               string  `json:"createdFromCommitID"`
	Committed                         bool    `json:"committed"`
	LastSyncedAllAttributionsCommitID string  `json:"lastSyncedAllAttributionsCommitID"`
	CreatedByID                       uint64  `json:"createdById"`
	UpdatedByID                       uint64  `json:"updatedById"`
	ArchivedByID                      *uint64 `json:"archivedById"`
	Permission                        string  `json:"permission"`
	Type                              string  `json:"type"`
	CanvasRepositoryID                uint64  `json:"canvasRepositoryId"`
	// Please understand this was a Choice we'll removed
	CanvasRepositoryID2 uint64 `json:"canvasRepositoryID"`
	CanvasRepoKey       string `json:"canvasRepoKey"`
	CanvasRepoName      string `json:"canvasRepoName"`
	Slug                string `json:"slug"`
	CanvasIcon          string `json:"canvasIcon"`
}

func SerializeCanvasBranch(branch *models.CanvasBranch) *GetCanvasBranchSerializer {
	view := &GetCanvasBranchSerializer{
		ID:                                branch.ID,
		UUID:                              branch.UUID.String(),
		Name:                              branch.Name,
		Key:                               branch.Key,
		PublicAccess:                      branch.PublicAccess,
		IsDraft:                           branch.IsDraft,
		IsDefault:                         branch.IsDefault,
		IsMerged:                          branch.IsMerged,
		IsRoughBranch:                     branch.IsRoughBranch,
		RoughFromBranchID:                 branch.RoughFromBranchID,
		RoughBranchCreatorID:              branch.RoughBranchCreatorID,
		FromBranchID:                      branch.FromBranchID,
		CreatedFromCommitID:               branch.CreatedFromCommitID,
		Committed:                         branch.Committed,
		LastSyncedAllAttributionsCommitID: branch.LastSyncedAllAttributionsCommitID,
		Permission:                        models.PGCanvasModerateSysName,
		CreatedByID:                       branch.CreatedByID,
		UpdatedByID:                       branch.UpdatedByID,
		Type:                              "BRANCH",
		CanvasRepositoryID:                branch.CanvasRepositoryID,
		CanvasRepositoryID2:               branch.CanvasRepositoryID,
		Slug:                              slug.Make(branch.Name),
	}
	if branch.CanvasRepository != nil {
		fmt.Println("KEY branch.CanvasRepository.Key:::: ", branch.CanvasRepository.Key)
		view.CanvasRepoKey = branch.CanvasRepository.Key
		view.CanvasRepoName = branch.CanvasRepository.Name
		view.CanvasIcon = branch.CanvasRepository.Icon
	} else {
		canvasRepo := App.Repo.CanvasRepo(branch.CanvasRepositoryID)
		view.CanvasRepoKey = canvasRepo.Key
		view.CanvasRepoName = canvasRepo.Name
		view.CanvasIcon = canvasRepo.Icon
	}
	return view
}

func MultiSerializeCanvasBranch(branches []models.CanvasBranch) *[]GetCanvasBranchSerializer {
	canvasBranches := &[]GetCanvasBranchSerializer{}
	for _, branch := range branches {
		branchSerializedData := SerializeCanvasBranch(&branch)
		*canvasBranches = append(*canvasBranches, *branchSerializedData)
	}
	return canvasBranches
}

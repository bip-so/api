package canvasbranch

import (
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/internal/shared"
	"time"

	"github.com/google/uuid"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gorm.io/datatypes"
)

type MergeRequestSerializer struct {
	Status        string `json:"status"`
	ID            uint64 `json:"id"`
	CommitMessage string `json:"commitMessage"`
	CreatedByID   uint64 `json:"createdByID"`
}
type MergeRequestDiffSerializer struct {
	Status        string `json:"status"`
	ID            uint64 `json:"id"`
	CommitMessage string `json:"commitMessage"`
}

type MergeRequestDefault struct {
	ID                  uint64                          `json:"id"`
	UUID                uuid.UUID                       `json:"uuid"`
	CommitMessage       string                          `json:"commitMessage"`
	Status              string                          `json:"status"`
	CreatedAt           time.Time                       `json:"createdAt"`
	SourceBranchID      uint64                          `json:"sourceBranchID"`
	DestinationBranchID uint64                          `json:"destinationBranchID"`
	CreatedByID         uint64                          `json:"createdByID"`
	ChangesAccepted     datatypes.JSON                  `json:"changesAccepted"`
	CreatedByUser       shared.CommonUserMiniSerializer `json:"createdByUser"`
	ClosedByUser        shared.CommonUserMiniSerializer `json:"closedByUser"`
	ClosedAt            time.Time                       `json:"closedAt"`
}
type MergeResponseObject struct {
	Branch            BranchDefault            `json:"branch"`
	CanvasRepo        MiniCanvasRepoSerializer `json:"canvasRepository"`
	MergeRequest      MergeRequestDefault      `json:"mergeRequest"`
	SourceBlocks      *[]BulkBlocks            `json:"sourceBlocks"`
	DestinationBlocks *[]BulkBlocks            `json:"destinationBlocks"`
}

type MiniCanvasRepoSerializer struct {
	ID           uint64    `json:"id"`
	UUID         string    `json:"uuid"`
	CollectionID uint64    `json:"collectionID"`
	Name         string    `json:"name"`
	Position     uint      `json:"position"`
	Icon         string    `json:"icon"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	CreatedByID  uint64    `json:"createdByID"`
	UpdatedByID  uint64    `json:"updatedByID"`
	Key          string    `json:"key"`
	Permission   string    `json:"permission"`
}

type DiffBeforeMergeResponseObject struct {
	SourceBlocks      *[]BulkBlocks `json:"sourceBlocks"`
	DestinationBlocks *[]BulkBlocks `json:"destinationBlocks"`
}

type MergeRequestDefaultSerializer struct {
	ID                  uint64                          `json:"id"`
	CommitMessage       string                          `json:"commitMessage"`
	UUID                uuid.UUID                       `json:"uuid"`
	CanvasRepositoryID  uint64                          `json:"canvasRepositoryID"`
	SourceBranchID      uint64                          `json:"sourceBranchID"`
	DestinationBranchID uint64                          `json:"destinationBranchID"`
	CreatedAt           time.Time                       `json:"createdAt"`
	UpdatedAt           time.Time                       `json:"updatedAt"`
	CreatedByID         uint64                          `json:"createdByID"`
	ClosedByUserId      uint64                          `json:"closedByUserId"`
	ClosedAt            time.Time                       `json:"closedAt"`
	Status              string                          `json:"status"`
	CommitID            string                          `json:"commitID"`
	SourceCommitID      string                          `json:"sourceCommitID"`
	DestinationCommitID string                          `json:"destinationCommitID"`
	ChangesAccepted     datatypes.JSON                  `json:"changesAccepted"`
	CreatedByUser       shared.CommonUserMiniSerializer `json:"createdByUser"`
}

type ManyMergeRequest struct {
	Data []MergeRequestDefaultSerializer `json:"data"`
}

func MergeRequestSerializerMany(instances *[]models.MergeRequest) ManyMergeRequest {
	var mergeRequests []MergeRequestDefaultSerializer
	for _, mr := range *instances {
		//MRuser, _ := user2.App.Repo.GetUser(map[string]interface{}{"id": mr.CreatedByID})
		MRuser, _ := queries.App.UserQueries.GetUserByID(mr.CreatedByID)

		mergeRequests = append(mergeRequests, MergeRequestDefaultSerializer{
			ID:                  mr.ID,
			CommitMessage:       mr.CommitMessage,
			UUID:                mr.UUID,
			CanvasRepositoryID:  mr.CanvasRepositoryID,
			SourceBranchID:      mr.SourceBranchID,
			DestinationBranchID: mr.DestinationBranchID,
			CreatedAt:           mr.CreatedAt,
			UpdatedAt:           mr.UpdatedAt,
			CreatedByID:         mr.CreatedByID,
			//ClosedByUserId:      mr.ClosedByUserId,
			ClosedAt:            mr.ClosedAt,
			Status:              mr.Status,
			CommitID:            mr.CommitID,
			SourceCommitID:      mr.SourceCommitID,
			DestinationCommitID: mr.DestinationCommitID,
			ChangesAccepted:     mr.ChangesAccepted,
			CreatedByUser: shared.CommonUserMiniSerializer{
				Id:        MRuser.ID,
				UUID:      MRuser.UUID.String(),
				FullName:  MRuser.FullName,
				Username:  MRuser.Username,
				AvatarUrl: MRuser.AvatarUrl,
			},
		})
	}
	many := ManyMergeRequest{}
	many.Data = mergeRequests
	return many
}

func SimpleMergeRequestSerializer(model *models.MergeRequest) MergeRequestDefault {
	//MRuser, _ := user2.App.Repo.GetUser(map[string]interface{}{"id": model.CreatedByID})
	MRuser, _ := queries.App.UserQueries.GetUserByID(model.CreatedByID)

	view := MergeRequestDefault{
		ID:                  model.ID,
		UUID:                model.UUID,
		CommitMessage:       model.CommitMessage,
		Status:              model.Status,
		CreatedAt:           model.CreatedAt,
		SourceBranchID:      model.SourceBranchID,
		DestinationBranchID: model.DestinationBranchID,
		CreatedByID:         model.CreatedByID,
		ChangesAccepted:     model.ChangesAccepted,
		CreatedByUser: shared.CommonUserMiniSerializer{
			Id:        MRuser.ID,
			UUID:      MRuser.UUID.String(),
			FullName:  MRuser.FullName,
			Username:  MRuser.Username,
			AvatarUrl: MRuser.AvatarUrl,
		},
		ClosedAt: model.ClosedAt,
	}
	if model.ClosedByUserId != nil {
		//closedByUser, _ := user2.App.Repo.GetUser(map[string]interface{}{"id": model.ClosedByUserId})
		closedByUser, _ := queries.App.UserQueries.GetUserByID(*model.ClosedByUserId)

		view.ClosedByUser = shared.CommonUserMiniSerializer{
			Id:        closedByUser.ID,
			UUID:      closedByUser.UUID.String(),
			FullName:  closedByUser.FullName,
			Username:  closedByUser.Username,
			AvatarUrl: closedByUser.AvatarUrl,
		}
	}
	return view
}

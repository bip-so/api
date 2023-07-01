package ar

import (
	"github.com/google/uuid"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/internal/shared"
	"time"
)

type AccessRequestSerializerDefault struct {
	ID                          uint64                          `json:"id"`
	UUID                        uuid.UUID                       `json:"uuid"`
	CreatedAt                   time.Time                       `json:"createdAt"`
	UpdatedAt                   time.Time                       `json:"updatedAt"`
	CreatedByID                 uint64                          `json:"createdByID"`
	ReviewedByUserID            uint64                          `json:"reviewedByUserID"`
	StudioID                    uint64                          `json:"studioID"`
	CollectionID                uint64                          `json:"collectionID"`
	CanvasRepositoryID          uint64                          `json:"canvasRepositoryID"`
	CanvasBranchID              uint64                          `json:"canvasBranchID"`
	CanvasBranchPermissionGroup string                          `json:"canvasBranchPermissionGroup"`
	Status                      string                          `json:"status"`
	Message                     string                          `json:"message"`
	CreatedByUser               shared.CommonUserMiniSerializer `json:"createdByUser"`
}
type ManyAccessRequests struct {
	Data []AccessRequestSerializerDefault `json:"data"`
}

func possiblyNullString(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func possiblyNullUINT64(v *uint64) uint64 {
	if v == nil {
		return 0
	}
	return *v
}

func AccessRequestSerializerMany(instances *[]models.AccessRequest) ManyAccessRequests {
	var accessRequests []AccessRequestSerializerDefault
	for _, ar := range *instances {
		// Deprecated
		//usr, _ := user2.App.Repo.GetUser(map[string]interface{}{"id": ar.CreatedByID})
		usr, _ := queries.App.UserQueries.GetUserByID(ar.CreatedByID)

		accessRequests = append(accessRequests, AccessRequestSerializerDefault{
			ID:                          ar.ID,
			UUID:                        ar.UUID,
			CreatedAt:                   ar.CreatedAt,
			UpdatedAt:                   ar.UpdatedAt,
			CreatedByID:                 ar.CreatedByID,
			StudioID:                    ar.StudioID,
			CollectionID:                ar.CollectionID,
			CanvasRepositoryID:          ar.CanvasRepositoryID,
			CanvasBranchID:              ar.CanvasBranchID,
			Status:                      ar.Status,
			Message:                     possiblyNullString(ar.Message),
			ReviewedByUserID:            possiblyNullUINT64(ar.ReviewedByUserID),
			CanvasBranchPermissionGroup: possiblyNullString(ar.CanvasBranchPermissionGroup),
			CreatedByUser: shared.CommonUserMiniSerializer{
				Id:        usr.ID,
				UUID:      usr.UUID.String(),
				FullName:  usr.FullName,
				Username:  usr.Username,
				AvatarUrl: usr.AvatarUrl,
			},
		})
	}
	Many := ManyAccessRequests{}
	Many.Data = accessRequests
	return Many
}

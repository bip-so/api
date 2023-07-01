package queries

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

// StudioMembersRequest
func (q studioMemberRequestQuery) CheckAndDeleteStudioMembersRequestIfExists(studioID uint64, userID uint64) {
	var quickcheck models.StudioMembersRequest
	_ = postgres.GetDB().Model(&models.StudioMembersRequest{}).Where(map[string]interface{}{
		"studio_id": studioID,
		"user_id":   userID,
	}).First(&quickcheck).Error
	if quickcheck.ID != 0 {
		// Delete this Instance
		_ = postgres.GetDB().Delete(quickcheck)
	}
}

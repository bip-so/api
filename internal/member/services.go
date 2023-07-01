package member

import (
	"gitlab.com/phonepost/bip-be-platform/internal/feed"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/supabase"
)

func (m memberService) GetMembersByStudio(studioId uint64, skip int) ([]models.Member, error) {
	members, err := queries.App.MemberQuery.GetMembersByStudioIDPaginated(studioId, skip)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (m memberService) LeaveStudio(userIds []uint64, studioId uint64) error {
	err := queries.App.MemberQuery.UpdateHasLeft(userIds, studioId, true)
	if err != nil {
		return nil
	}
	members, err := queries.App.MemberQuery.GetMembersByUserIDs(userIds, studioId)
	if err != nil {
		return nil
	}
	err = m.RemoveMembersToStudioMemberRole(studioId, members)
	if err != nil {
		return nil
	}
	go func() {
		for _, userID := range userIds {
			queries.App.StudioQueries.DeleteUserAssociatedStudioDataByUserID(userID)
			supabase.UpdateUserSupabase(userID, true)
			// unfollow in stream
			feed.App.Service.LeaveStudio(studioId, userID)
		}
	}()

	return err
}

func (m memberService) BanUser(userID uint64, studioId uint64, banReason string, removedByID uint64) error {
	err := queries.App.MemberQuery.BanUser(userID, studioId, banReason, removedByID)
	if err != nil {
		return nil
	}
	member, err := queries.App.MemberQuery.GetMember(map[string]interface{}{"user_id": userID, "studio_id": studioId})
	if err != nil {
		return nil
	}
	err = m.RemoveMembersToStudioMemberRole(studioId, []models.Member{*member})
	if err != nil {
		return nil
	}
	return err
}

func (m memberService) RemoveMembersToStudioMemberRole(studioID uint64, members []models.Member) error {
	studioPermission, err := queries.App.StudioPermissionQuery.GetStudioPermission(
		map[string]interface{}{"studio_id": studioID, "permission_group": "pg_studio_none"})
	if err != nil {
		return err
	}

	if studioPermission.RoleId != nil {
		role, err := queries.App.RoleQuery.GetRole(*studioPermission.RoleId)
		if err != nil {
			return err
		}
		err = queries.App.MemberQuery.RemoveMembersInRole(members, role)
		if err != nil {
			return err
		}
	}
	return nil
}

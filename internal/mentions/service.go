package mentions

import (
	"errors"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

func (s mentionsService) AddMentionManager(obj MentionPost, user *models.User, studioID uint64) ([]map[string]interface{}, error) {
	var userObjects *[]MentionedUserSerializer
	var canvasObjects *[]MentionedCanvasSerializer
	var roleObjects *[]MentionedRolesSerializer
	var errGettingObjects error

	if len(obj.Users) > 0 {
		userObjects, errGettingObjects = s.GetUserObjects(obj.Users)
		if errGettingObjects != nil {
			return []map[string]interface{}{}, errGettingObjects
		}
	}

	if len(obj.Branches) > 0 {
		canvasObjects, errGettingObjects = s.GetCanvasObjects(obj.Branches)
		if errGettingObjects != nil {
			return []map[string]interface{}{}, errGettingObjects
		}
	}

	if len(obj.Roles) > 0 {
		roleObjects, errGettingObjects = s.GetRoleObjects(obj.Roles)
		if errGettingObjects != nil {
			return []map[string]interface{}{}, errGettingObjects
		}
	}

	switch obj.Scope {
	case "block":
		data, err := s.AddMentionToBlock(obj, userObjects, canvasObjects, roleObjects, user, studioID)
		if err != nil {
			return []map[string]interface{}{}, err
		}
		return data, nil
	case "block_thread":
		data, err := s.AddMentionToBlockThread(obj, userObjects, canvasObjects, roleObjects, user, studioID)
		if err != nil {
			return []map[string]interface{}{}, err
		}
		return data, nil
	case "block_comment":
		data, err := s.AddMentionToBlockThreadComment(obj, userObjects, canvasObjects, roleObjects, user, studioID)
		if err != nil {
			return []map[string]interface{}{}, err
		}
		return data, nil
	case "reel":
		data, err := s.AddMentionToReel(obj, userObjects, canvasObjects, roleObjects, user, studioID)
		if err != nil {
			return []map[string]interface{}{}, err
		}
		return data, nil
	case "reel_comment":
		data, err := s.AddMentionToReelComment(obj, userObjects, canvasObjects, roleObjects, user, studioID)
		if err != nil {
			return []map[string]interface{}{}, err
		}
		return data, nil
	default:
		return nil, errors.New("scope validation Error")
	}
}

func (s mentionsService) GetUserObjects(users []uint64) (*[]MentionedUserSerializer, error) {
	instances, err := App.Repo.GetUserObjects(users)
	if err != nil {
		return nil, err
	}
	var serials []MentionedUserSerializer
	for _, k := range *instances {
		serials = append(serials, MentionedUserSerializer{
			Type:      "user",
			ID:        k.ID,
			UUID:      k.UUID.String(),
			FullName:  k.FullName,
			Username:  k.Username,
			AvatarUrl: k.AvatarUrl,
		})
	}

	return &serials, nil
}

func (s mentionsService) GetRoleObjects(roles []uint64) (*[]MentionedRolesSerializer, error) {
	instances, err := App.Repo.GetRoleObjects(roles)
	if err != nil {
		return nil, err
	}
	var serials []MentionedRolesSerializer
	for _, k := range *instances {
		serials = append(serials, MentionedRolesSerializer{
			Type: "role",
			ID:   k.ID,
			Name: k.Name,
			UUID: k.UUID.String(),
		})
	}

	return &serials, nil
}

func (s mentionsService) GetCanvasObjects(branches []uint64) (*[]MentionedCanvasSerializer, error) {
	instances, err := App.Repo.GetBranchObjects(branches)
	if err != nil {
		return nil, err
	}
	var serials []MentionedCanvasSerializer
	for _, k := range *instances {
		serials = append(serials, MentionedCanvasSerializer{
			Type: "branch",
			ID:   k.ID,
			Name: k.Name,
			UUID: k.UUID.String(),
			Key:  k.Key,
		})
	}

	return &serials, nil
}

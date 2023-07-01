package mentions

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

func (s mentionsService) MentionMapMaker(userObjects *[]MentionedUserSerializer, canvasObjects *[]MentionedCanvasSerializer, roleObjects *[]MentionedRolesSerializer, user *models.User, studioID uint64) []map[string]interface{} {
	var MentionJson []map[string]interface{}
	if userObjects != nil {
		for _, v := range *userObjects {
			MentionJson = append(MentionJson, map[string]interface{}{
				"studioID":              studioID,
				"type":                  v.Type,
				"id":                    v.ID,
				"uuid":                  v.UUID,
				"fullName":              v.FullName,
				"username":              v.Username,
				"avatarUrl":             v.AvatarUrl,
				"createdByUserID":       user.ID,
				"createdByUserUsername": user.Username,
				"createdByUserFullName": user.FullName,
			})
		}
	}
	if canvasObjects != nil {
		for _, v := range *canvasObjects {
			canvasBranch, _ := App.Repo.GetBranchObject(v.ID)
			canvasRepo, _ := App.Repo.GetRepo(map[string]interface{}{"id": canvasBranch.CanvasRepositoryID})
			MentionJson = append(MentionJson, map[string]interface{}{
				"studioID":              studioID,
				"type":                  v.Type,
				"id":                    v.ID,
				"uuid":                  v.UUID,
				"name":                  v.Name,
				"key":                   v.Key,
				"repoID":                canvasRepo.ID,
				"repoKey":               canvasRepo.Key,
				"repoName":              canvasRepo.Name,
				"repoUUID":              canvasRepo.UUID.String(),
				"createdByUserID":       user.ID,
				"createdByUserUsername": user.Username,
				"createdByUserFullName": user.FullName,
			})
		}
		fmt.Println(MentionJson)
	}
	if roleObjects != nil {
		for _, v := range *roleObjects {
			MentionJson = append(MentionJson, map[string]interface{}{
				"studioID":              studioID,
				"type":                  v.Type,
				"id":                    v.ID,
				"uuid":                  v.UUID,
				"name":                  v.Name,
				"createdByUserID":       user.ID,
				"createdByUserUsername": user.Username,
				"createdByUserFullName": user.FullName,
			})
		}
		fmt.Println(MentionJson)
	}
	return MentionJson
}

package studio

import (
	"fmt"
	"gorm.io/gorm"
	"strings"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/search"
)

type StudioService interface {
	NewStudio(name, handle, description, website, imageURL string, userID uint64) (*models.Studio, error)
	AddAllStudiosToAlgolia() error
	AddStudioToAlgolia(studioID uint64) error
}

func (ss studioService) NewStudio(name, handle, description, website, imageURL string, userID uint64) (*models.Studio, error) {
	// validate here
	if len(handle) < 2 {
		return nil, StudioErrorsHandleInvalid
	}
	if len(handle) > 36 {
		return nil, StudioErrorsHandleLenMax
	}
	if len(name) > 100 {
		return nil, StudioErrorsNameLenMax
	}
	if len(description) > 300 {
		return nil, StudioErrorsDescriptionLenMax
	}
	isValid := models.CheckHandleValidity(handle)
	if !isValid {
		return nil, StudioErrorsHandleInvalid
	}
	available, err := App.StudioRepo.checkHandleAvailablity(handle)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, ErrHandleUnavailable
	}
	return &models.Studio{
		DisplayName: name,
		Handle:      handle,
		Description: description,
		Website:     website,
		ImageURL:    imageURL,
		CreatedByID: userID,
		UpdatedByID: userID,
	}, nil
}

func (ss studioService) AddAllStudiosToAlgolia() error {
	studios, err := App.StudioRepo.GetAllStudios()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	studioDocs := []interface{}{}
	for _, stdio := range *studios {
		res := strings.ToLower(stdio.DisplayName)
		res1 := strings.ToLower(stdio.Handle)
		res2 := strings.ToLower(stdio.Description)
		if strings.Contains(res, "test") || strings.Contains(res1, "test") || strings.Contains(res2, "test") {
			continue
		}
		studioDocs = append(studioDocs, *StudioModelToStudioDocument(&stdio))
	}
	err = search.GetIndex(search.StudioDocumentIndexName).SaveRecords(studioDocs)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func (ss studioService) AddStudioToAlgolia(studioID uint64) error {
	stdio, err := App.StudioRepo.GetStudioByID(studioID)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if stdio.IsArchived {
		return nil
	}
	res := strings.ToLower(stdio.DisplayName)
	res1 := strings.ToLower(stdio.Handle)
	res2 := strings.ToLower(stdio.Description)
	if strings.Contains(res, "test") || strings.Contains(res1, "test") || strings.Contains(res2, "test") {
		return nil
	}
	studioDoc := StudioModelToStudioDocument(stdio)
	err = search.GetIndex(search.StudioDocumentIndexName).SaveRecord(studioDoc)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func (ss studioService) DeleteStudioFromAlgolia(studioID uint64) error {
	err := search.GetIndex(search.StudioDocumentIndexName).DeleteRecordByID(studioID)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return err
}

// Todo: Will return true if a user with this email is found in the studio
func (ss studioService) DoesUserAlreadyBelongsToStudio(email string, studioID uint64) bool {
	return false
}

func (ss studioService) StudioStats(studioID uint64) map[string]interface{} {
	//totalRepoCount := App.StudioRepo.GetStudioRepoCount(studioID)
	return map[string]interface{}{}
}

// Will Return if user is found on BIP

/// Topics

type TopicService interface {
	NewTopic(name string) (*models.Topic, error)
}

func (ts studioTopicService) NewTopic(name string) (*models.Topic, error) {
	if len(name) > 100 {
		return nil, StudioErrorsTopicLenMax
	}
	return &models.Topic{
		Name: name,
	}, nil
}

func (s studioService) GetRequestStudioList(studioID uint64) (*[]models.StudioMembersRequest, error) {
	var results *[]models.StudioMembersRequest
	err := s.db.Model(&models.StudioMembersRequest{}).
		Where("studio_id = ? and action = ?", studioID, "Pending").
		Preload("User").
		Order("created_at DESC").
		Find(&results).Error

	if err != nil {
		return nil, err
	}
	return results, nil
}

func (ss studioService) MockStudioMembershipRequestObject(studioID uint64, loggedInUserID uint64, requestingUserID uint64) *models.StudioMembersRequest {
	return &models.StudioMembersRequest{
		UserID:     requestingUserID,
		Action:     "Pending",
		StudioID:   studioID,
		ActionByID: loggedInUserID,
	}
}

func (s studioService) CreateStudioRequestInstance(instance *models.StudioMembersRequest) error {
	// Add a dedup check.
	var quickcheck models.StudioMembersRequest
	_ = s.db.Model(&models.StudioMembersRequest{}).Where(map[string]interface{}{
		"studio_id": instance.StudioID,
		"user_id":   instance.UserID,
	}).First(&quickcheck).Error
	fmt.Println("StudioMembersRequest: Found")
	if quickcheck.ID == 0 {
		err := s.db.Create(instance).Error
		return err
	} else if quickcheck.Action != "Pending" {
		s.db.Model(&models.StudioMembersRequest{}).Where(map[string]interface{}{
			"studio_id": instance.StudioID,
			"user_id":   instance.UserID,
		}).Update("action", "Pending")
	}

	return nil
}

func (s studioService) BuildStudioStats(studioID uint64) map[string]interface{} {
	/*
		SELECT
		( SELECT COUNT(*) FROM posts where studio_id = ? ) AS post_count,
		( SELECT COUNT(*) FROM reels where studio_id = ?) AS reel_count,
		( SELECT COUNT(*) FROM canvas_repositories where studio_id = ? and is_archived = false ) AS repo_count,
		( SELECT COUNT(*) FROM canvas_repositories where studio_id = ? and is_archived = false ) AS repo_count,
		( SELECT COUNT(*) FROM collections where studio_id = ?) AS collection_count
	*/

	type Result struct {
		Posts       int
		Reels       int
		Repos       int
		Privates    int
		Collections int
	}
	query := "SELECT ( SELECT COUNT(*) FROM posts where studio_id = ? ) AS Posts, ( SELECT COUNT(*) FROM reels where studio_id = ?) AS Reels, ( SELECT COUNT(*) FROM canvas_repositories where studio_id = ? and is_archived = false ) AS Repos, ( SELECT COUNT(*) FROM canvas_repositories as CR INNER JOIN canvas_branches as CB ON CB.id = CR.default_branch_id where CR.studio_id = ? and CR.is_archived = false and cb.public_access = 'private' and CR.is_published = true) AS Privates, ( SELECT COUNT(*) FROM collections where studio_id = ?) AS Collections"
	var result Result
	s.db.Raw(query, studioID, studioID, studioID, studioID, studioID).Scan(&result)
	fmt.Println(result)

	var stats = map[string]interface{}{
		"all_repos":         result.Repos,
		"private_repos":     result.Privates,
		"collections_count": result.Collections,
		"post_count":        result.Posts,
		"reel_count":        result.Reels,
	}
	return stats
}

func (s studioService) CheckIsRequested(userID, studioID uint64) bool {
	if userID == 0 {
		return false
	}
	_, err := App.StudioRepo.GetUserStudioMemberRequest(userID, studioID)
	if err == gorm.ErrRecordNotFound {
		return false
	}
	return true
}

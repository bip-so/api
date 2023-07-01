package studio

import (
	"gitlab.com/phonepost/bip-be-platform/internal/shared"
	"time"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

type StudioSerializer struct {
	ID                    uint64            `json:"id"`
	UUID                  string            `json:"uuid"`
	DisplayName           string            `json:"displayName"`
	Handle                string            `json:"handle"`
	ImageURL              string            `json:"imageUrl"`
	Description           string            `json:"description"`
	Website               string            `json:"website"`
	FollowerCount         int               `json:"followerCount"`
	IsJoined              *bool             `json:"isJoined"`
	IsRequested           bool              `json:"isRequested"`
	Topics                []TopicSerializer `json:"topics"`
	CreatedAt             time.Time         `json:"createdAt"`
	UpdatedAt             time.Time         `json:"updatedAt"`
	MembersCount          int64             `json:"membersCount"`
	Permission            string            `json:"permission"`
	CreatedByID           uint64            `json:"createdById"`
	DefaultCanvasRepoID   uint64            `json:"defaultCanvasRepoId"`
	DefaultCanvasRepoName string            `json:"defaultCanvasRepoName"`
	DefaultCanvasRepoKey  string            `json:"defaultCanvasRepoKey"`
	DefaultCanvasBranchID uint64            `json:"defaultCanvasBranchId"`
	IsEarlyAdopter        bool              `json:"isEarlyAdopter"`
	IsNonProfit           bool              `json:"isNonProfit"`
	AllowPublicMembership bool              `json:"allowPublicMembership"`
	StripeCustomerID      string            `json:"stripeCustomerID"`
	StripeProductID       string            `json:"stripeProductID"`
	StripePriceID         string            `json:"stripePriceID"`
	StripePriceUnit       int64             `json:"stripePriceUnit"`
	StripeSubscriptionsID string            `json:"stripeSubscriptionsID"`
}

type TopicSerializer struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func SerializeTopicData(topics []models.Topic) []TopicSerializer {
	var topicsData []TopicSerializer
	for _, topic := range topics {
		topicsData = append(topicsData, TopicSerializer{
			ID:   topic.ID,
			Name: topic.Name,
		})
	}
	return topicsData
}

func SerializeStudioForUser(studio *models.Studio, authUser *models.User, members *[]models.Member) *StudioSerializer {
	serialized := SerializeStudio(studio)
	serialized.IsJoined = checkIsJoined(studio, authUser, members)
	if authUser != nil {
		serialized.IsRequested = App.StudioService.CheckIsRequested(authUser.ID, studio.ID)
	}
	return serialized
}

func SerializeStudio(studio *models.Studio) *StudioSerializer {
	return &StudioSerializer{
		ID:                    studio.ID,
		UUID:                  studio.UUID.String(),
		DisplayName:           studio.DisplayName,
		Handle:                studio.Handle,
		Description:           studio.Description,
		ImageURL:              studio.ImageURL,
		Website:               studio.Website,
		FollowerCount:         studio.ComputedFollowerCount,
		Topics:                SerializeTopicData(studio.Topics),
		CreatedAt:             studio.CreatedAt,
		UpdatedAt:             studio.UpdatedAt,
		CreatedByID:           studio.CreatedByID,
		IsEarlyAdopter:        studio.IsEarlyAdopter,
		IsNonProfit:           studio.IsNonProfit,
		AllowPublicMembership: studio.AllowPublicMembership,
		StripeCustomerID:      studio.StripeCustomerID,
		StripeProductID:       studio.StripeProductID,
		StripePriceID:         studio.StripePriceID,
		StripePriceUnit:       studio.StripePriceUnit,
		StripeSubscriptionsID: studio.StripeSubscriptionsID,
	}
}

type MemberCountSerializer struct {
	Members uint64 `json:"members"`
}

func SerializeMemberCount(count uint64) *MemberCountSerializer {
	return &MemberCountSerializer{
		Members: count,
	}
}

func checkIsJoined(studio *models.Studio, authUser *models.User, studioMembers *[]models.Member) *bool {
	if studioMembers == nil {
		return nil
	}
	joined := false
	for _, member := range *studioMembers {
		if member.StudioID == studio.ID && authUser.ID == member.UserID {
			joined = true
			return &joined
		}
	}
	return &joined
}

// Single StudioMembersRequest
type StudioMembersRequestDefaultSerializer struct {
	ID       uint64                          `json:"id"`
	UserID   uint64                          `json:"userId"`
	User     shared.CommonUserMiniSerializer `json:"user"`
	Action   string                          `json:"action"`
	StudioID uint64                          `json:"studioID"`
}

func SerializeSingleStudioMembersRequest(model *models.StudioMembersRequest) *StudioMembersRequestDefaultSerializer {
	return &StudioMembersRequestDefaultSerializer{
		ID:       model.ID,
		UserID:   model.UserID,
		Action:   model.Action,
		StudioID: model.StudioID,
		User: shared.CommonUserMiniSerializer{
			Id:        model.User.ID,
			UUID:      model.User.UUID.String(),
			Username:  model.User.Username,
			FullName:  model.User.FullName,
			AvatarUrl: model.User.AvatarUrl,
		},
	}
}

func ManySerializeStudioMembersRequest(modelInstances *[]models.StudioMembersRequest) *[]StudioMembersRequestDefaultSerializer {
	results := &[]StudioMembersRequestDefaultSerializer{}

	if len(*modelInstances) == 0 {
		return results
	}

	for _, model := range *modelInstances {
		*results = append(*results, *SerializeSingleStudioMembersRequest(&model))
	}
	return results
}

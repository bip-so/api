package post

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/internal/reel"
	"gitlab.com/phonepost/bip-be-platform/internal/shared"
	"gitlab.com/phonepost/bip-be-platform/internal/studio"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"gorm.io/datatypes"
	"time"
)

type PostSerializer struct {
	ID                uint64                            `json:"id"`
	UUID              string                            `json:"uuid"`
	StudioID          uint64                            `json:"studioID"`
	Children          datatypes.JSON                    `json:"children"`
	Attributes        datatypes.JSON                    `json:"attributes"`
	CreatedByID       uint64                            `json:"createdById"`
	UpdatedByID       uint64                            `json:"updatedById"`
	CreatedAt         time.Time                         `json:"createdAt"`
	CreatedByUser     shared.CommonUserMiniSerializer   `json:"createdByUser"`
	UpdatedByUser     shared.CommonUserMiniSerializer   `json:"updatedByUser"`
	Studio            shared.CommonStudioMiniSerializer `json:"studio"`
	CommentCount      uint                              `json:"commentCount"`
	ReactionCounter   []shared.ReactedCount             `json:"reactions"`
	ReactionCopy      string                            `json:"reactionCopy"`
	IsUserFollower    bool                              `json:"isUserFollower"`
	IsStudioMember    bool                              `json:"isStudioMember"`
	IsUserStudioAdmin bool                              `json:"isUserStudioAdmin"`
}

func OnePostSerializerData(modelInstances *models.Post, loggedInUserID uint64) *PostSerializer {
	posts := &PostSerializer{}
	posts = SinglePostSerializerData(modelInstances, loggedInUserID)
	return posts
}

func SinglePostSerializerData(model *models.Post, loggedInUserID uint64) *PostSerializer {
	fmt.Println(model)

	var userFollowings *[]models.FollowUser
	//creatorByUser, _ := user2.App.Repo.GetUser(map[string]interface{}{"id": model.CreatedByID})
	//updatedByUser, _ := user2.App.Repo.GetUser(map[string]interface{}{"id": model.CreatedByID})

	creatorByUser, _ := queries.App.UserQueries.GetUserByID(model.CreatedByID)
	updatedByUser, _ := queries.App.UserQueries.GetUserByID(model.UpdatedByID)

	userFollowings, _ = reel.App.Repo.GetUserFollowings(loggedInUserID)
	view := PostSerializer{
		StudioID:     model.StudioID,
		ID:           model.ID,
		UUID:         model.UUID.String(),
		CreatedAt:    model.CreatedAt,
		CreatedByID:  model.CreatedByID,
		UpdatedByID:  model.UpdatedByID,
		Children:     model.Children,
		Attributes:   model.Attributes,
		CommentCount: model.CommentCount,
		CreatedByUser: shared.CommonUserMiniSerializer{
			Id:        creatorByUser.ID,
			UUID:      creatorByUser.UUID.String(),
			FullName:  creatorByUser.FullName,
			Username:  creatorByUser.Username,
			AvatarUrl: creatorByUser.AvatarUrl,
		},
		UpdatedByUser: shared.CommonUserMiniSerializer{
			Id:        updatedByUser.ID,
			UUID:      updatedByUser.UUID.String(),
			FullName:  updatedByUser.FullName,
			Username:  updatedByUser.Username,
			AvatarUrl: updatedByUser.AvatarUrl,
		},
		Studio: shared.CommonStudioMiniSerializer{
			ID:                    model.Studio.ID,
			UUID:                  model.Studio.UUID.String(),
			DisplayName:           model.Studio.DisplayName,
			Handle:                model.Studio.Handle,
			ImageURL:              model.Studio.ImageURL,
			CreatedByID:           model.Studio.CreatedByID,
			AllowPublicMembership: model.Studio.AllowPublicMembership,
			IsRequested:           studio.App.StudioService.CheckIsRequested(loggedInUserID, model.StudioID),
		},
		ReactionCounter: TransposePostReactions(model.ID, loggedInUserID),
		ReactionCopy:    BuildRandomReactionsString(model, loggedInUserID),

		IsStudioMember:    shared.IsUserStudioMember(loggedInUserID, model.StudioID),
		IsUserStudioAdmin: shared.IsUserStudioAdmin(loggedInUserID, model.StudioID),
	}
	if model.CreatedByID == loggedInUserID {
		view.IsUserFollower = true
	} else {
		view.IsUserFollower = shared.CheckIsUserFollowing(model.CreatedByID, userFollowings)
	}
	return &view
}

func ManyPostSerializerData(modelInstances *[]models.Post, loggedInUser uint64) *[]PostSerializer {
	posts := &[]PostSerializer{}

	if len(*modelInstances) == 0 {
		return posts
	}

	for _, model := range *modelInstances {
		*posts = append(*posts, *SerializePostWithReactions(&model, loggedInUser))
	}

	return posts
}

// Likes text - Darshana Hazarika and 62 others reacted to this post.  - BE.
func BuildRandomReactionsString(model *models.Post, loggedInUserID uint64) string {
	count := shared.CountResults(map[string]interface{}{"post_id": model.ID}, &models.PostReaction{})
	if count == 0 {
		return ""
	}
	pr := models.PostReaction{}
	if count == 1 {
		_ = postgres.GetDB().Model(&models.PostReaction{}).
			Where("post_id = ?", model.ID).
			Preload("CreatedByUser").
			Limit(1).
			First(&pr)
		return pr.CreatedByUser.FullName
	}

	pr2 := models.PostReaction{}
	// This query returns Random
	_ = postgres.GetDB().Model(&models.PostReaction{}).
		Where("post_id = ? and created_by_id != ?", model.ID, loggedInUserID).
		Order("created_at desc").
		Preload("CreatedByUser").
		Limit(1).
		First(&pr2)
	var name string
	if pr2.ID == 0 {
		name = ""
	} else {
		name = pr2.CreatedByUser.FullName
	}
	//countForString := count - 1
	//return fmt.Sprintf("%s and %d others reacted to this post", name, countForString)
	return fmt.Sprintf("%s", name)
	// and 62 others reacted to this post.
}

//func SerializePostWithReactions(model *models.Post, user *models.User) *PostSerializer {
func SerializePostWithReactions(model *models.Post, loggedInUser uint64) *PostSerializer {
	serialized := SinglePostSerializerData(model, loggedInUser)
	return serialized
}

func TransposePostReactions(postid uint64, uid uint64) []shared.ReactedCount {
	// Database call
	type ReactionData struct {
		Emoji string `json:"emoji"`
		Count int    `json:"count"`
	}
	var Reactions []ReactionData
	postgres.GetDB().Raw("SELECT emoji, COUNT(emoji) as Count  FROM post_reactions  where post_id = ? GROUP BY emoji", postid).Scan(&Reactions)
	// Prepare the Response
	var updated []shared.ReactedCount
	for _, v := range Reactions {
		// We need to query DB for Existence of the User Reaction
		var exists bool
		_ = postgres.GetDB().Raw("SELECT EXISTS (SELECT 1 from post_reactions WHERE created_by_id = ? and post_id = ? and emoji = ?)", uid, postid, v.Emoji).Scan(&exists)
		updated = append(updated, shared.ReactedCount{
			Emoji:   v.Emoji,
			Count:   v.Count,
			Reacted: exists,
		})
	}
	return updated
}

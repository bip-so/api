package user

import (
	"time"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

type UserDocument struct {
	ID          uint64    `json:"id"`
	UUID        string    `json:"uuid"`
	ObjectID    string    `json:"objectID"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	FullName    string    `json:"fullName"`
	Username    string    `json:"handle"`
	AvatarUrl   string    `json:"avatarUrl"`
	IsFollowing *bool     `json:"isFollowing"`
	Email       string    `json:"email"`
	Followers   uint64    `json:"followers"`
	Following   uint64    `json:"following"`
}

func UserModelToUserDocument(user *models.User) *UserDocument {
	return &UserDocument{
		ID:          user.ID,
		UUID:        user.UUID.String(),
		ObjectID:    user.UUID.String(),
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		FullName:    user.FullName,
		Username:    user.Username,
		AvatarUrl:   user.AvatarUrl,
		IsFollowing: nil,
		Email:       user.Email.String,
	}
}

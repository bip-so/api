package global

import (
	"gitlab.com/phonepost/bip-be-platform/internal/shared"
	"time"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

type SerializedMessage struct {
	ID          uint64                          `json:"id"`
	UUID        string                          `json:"uuid"`
	Text        string                          `json:"text"`
	Author      shared.CommonUserMiniSerializer `json:"author"`
	Attachments []string                        `json:"attachments"`
	Timestamp   time.Time                       `json:"timestamp"`
	CreatedAt   time.Time                       `json:"createdAt"`
	UpdatedAt   time.Time                       `json:"updatedAt"`
}

func SerializeMessages(messages *[]models.Message) *[]SerializedMessage {
	serializeMessages := []SerializedMessage{}
	for _, message := range *messages {
		author := &message.Author
		serializeMessages = append(serializeMessages, SerializedMessage{
			ID:   message.ID,
			UUID: message.UUID.String(),
			Text: message.Text,
			//Author:      user.UserMiniSerializerData(&message.Author),
			Author: shared.CommonUserMiniSerializer{
				Id:        author.ID,
				UUID:      author.UUID.String(),
				FullName:  author.FullName,
				Username:  author.Username,
				AvatarUrl: author.AvatarUrl,
			},
			Attachments: message.Attachments,
			Timestamp:   message.Timestamp,
			CreatedAt:   message.CreatedAt,
			UpdatedAt:   message.UpdatedAt,
		})
	}
	return &serializeMessages
}

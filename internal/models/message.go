package models

import (
	"time"

	"github.com/lib/pq"
)

const (
	MessageTypeDiscord = "discord"
	MessageTypeSlack   = "slack"
)

type Message struct {
	BaseModel
	Text           string
	Type           string
	RefID          string `gorm:"index:idx_refid_userid,unique"`
	UserID         uint64 `gorm:"index:idx_refid_userid,unique"`
	AuthorID       uint64
	Attachments    pq.StringArray `gorm:"type:text[]"`
	IsUsed         bool           `gorm:"default:false;index:idx_refid_userid,unique"`
	TeamName       string
	TeamImage      string
	TeamID         string
	NotificationID string
	Timestamp      time.Time

	User   User `gorm:"foreignkey:UserID;"`
	Author User `gorm:"foreignkey:AuthorID;"`
}

func newMessage(refID, text string, authorID uint64, userID uint64, mtype string, timestamp time.Time, attachments []string, teamID, teamName, teamImage, notificationID string) *Message {
	return &Message{
		Text:           text,
		AuthorID:       authorID,
		UserID:         userID,
		RefID:          refID,
		Type:           mtype,
		Timestamp:      timestamp,
		Attachments:    attachments,
		TeamName:       teamName,
		TeamImage:      teamImage,
		TeamID:         teamID,
		NotificationID: notificationID,
	}
}

func NewDiscordMessage(refID, text string, authorID uint64, userID uint64, timestamp time.Time, attachments []string, teamID, teamName, teamImage, notificationID string) *Message {
	return newMessage(refID, text, authorID, userID, MessageTypeDiscord, timestamp, attachments, teamID, teamName, teamImage, notificationID)
}

func NewSlackMessage(refID, text string, authorID uint64, userID uint64, timestamp time.Time, attachments []string, teamID, teamName, teamImage, notificationID string) *Message {
	return newMessage(refID, text, authorID, userID, MessageTypeSlack, timestamp, attachments, teamID, teamName, teamImage, notificationID)
}

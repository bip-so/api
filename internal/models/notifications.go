package models

import (
	"gorm.io/datatypes"
	"time"
)

/*
	Notification kafka Flow

	At every notification trigger action we create a kafka event.
	KafkaEvent will be of
	key = Entity.
	Value = Notification specific data.(Will confirm the format later)

	When kafka event is consumed then
	Notification creation process
	We get the NotificationEntity based on kafka KEY entity(for now there are 10).
	NotificationEntity contains the default values of that entity.
	And Based on event we get the activity and default text from notificationEntity.events
	After getting all the data we create the notification record in PG DB or, we can make a call to store it in MongoDB.

	Notification needs to send to app
	We update the count of the NotificationCount table

	Notification Sending to user via email/Discord
	NotificationEntity has default values to send it to email/Discord
	If default value is false we return not sending the notification
	If defaults true we need to get userSettings and check and then decide to send or not.
*/
type Notification struct {
	BaseModel
	Version uint

	// Event -> trigger Action Eg: Published
	// Activity -> trigger Action Eg: Canvas Published
	// Entity -> UserSetting Notification entity Eg: SystemNotifications
	Entity   string
	Event    string
	Activity string
	Priority string // enum types will be ["high", "medium", "low"]
	// For some reason email or discord notification fails to send we can store the high priority notifications in redis. So that we can run this separately by some scripts.

	Text string // Actual notification message.

	// This will be only ID not a foreign key. To Query the db easily.
	NotifierID uint64  // User who receives the notification
	ReactorID  *uint64 // Who reacted on notification.
	StudioID   *uint64

	ObjectId            *uint64
	ContentObject       *string
	TargetObjectId      *uint64
	TargetContentObject *string

	// to store some extra data
	/*
		{
			"notifier": {"id": "", "name": "", "avatarUrl": "", "handle": ""},
			"reactor": {"id": "", "name": "", "avatarUrl": "", "handle": ""},
			"studio": {"id": "", "DisplayName": "", "handle": ""},
			"canvasRepo": {"id": "", "name": "", "handle": "", "emoji": ""},
			"canvasBranch": {"id": "", "name": ""},
			"objectId": 12,
			"contentType": "x",
			"targetObjectId": "",
			"targetContentTypeId": "",
		}
	*/
	ExtraData datatypes.JSON // to store any denormalized data as json

	// To know the notification status. These are used for app notifications
	Read bool `gorm:"read, default:false"`
	Seen bool `gorm:"seen, default:false"`

	// IsPersonal helps with the calculating of personal notification count
	IsPersonal            bool    `gorm:"isPersonal;default:false;"`
	DiscordDmID           *string // Discord messageId to modify the message on click on notification on discord.
	SlackDmID             *string
	SlackChannelID        *string
	IsWebsiteNotification bool `gorm:"default:true;"`

	// User Related data about the notification.
	CreatedByID  uint64
	UpdatedByID  uint64
	ArchivedByID uint64
	IsArchived   bool
	ArchivedAt   time.Time
}

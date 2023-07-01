package notifications

import (
	"encoding/json"
	"time"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gorm.io/datatypes"
)

type PostNotification struct {
	Entity                string                `json:"entity"`
	Event                 string                `json:"event"`
	Activity              string                `json:"activity"`
	Priority              string                `json:"priority"`
	Text                  string                `json:"text"`
	NotifierIDs           []uint64              `json:"notifierIds"`
	ReactorID             *uint64               `json:"reactorId"`
	StudioID              *uint64               `json:"studioId"`
	ExtraData             NotificationExtraData `json:"notificationExtraData"`
	IsPersonal            bool                  `json:"isPersonal"`
	DiscordDmID           *string               `json:"discordDmId"`
	IsWebsiteNotification bool                  `json:"isWebsiteNotification"`
	CreatedByID           uint64                `json:"createdById"`
	RoleIDs               *[]uint64             `json:"roleIDs"`
	ObjectID              *uint64               `json:"objectId"`
	ContentObject         *string               `json:"contentObject"`
	TargetObjectID        *uint64               `json:"targetObjectId"`
	TargetContentObject   *string               `json:"targetContentObjectId"`
}

type NotificationExtraData struct {
	CollectionID           uint64                   `json:"collectionId"`
	CanvasRepoID           uint64                   `json:"canvasRepoId"`
	CanvasBranchID         uint64                   `json:"canvasBranchId"`
	Status                 string                   `json:"actionStatus"`
	Message                string                   `json:"message"`
	PermissionGroup        string                   `json:"permissionGroup"`
	Data                   string                   `json:"data"`
	DiscordComponents      []interface{}            `json:"discordComponents"`
	SlackComponents        []map[string]interface{} `json:"slackComponents"`
	DiscordMessage         []string                 `json:"discordMessage"`
	SlackMessage           string                   `json:"slackMessage"`
	AppUrl                 string                   `json:"appUrl"`
	ActionOnText           string                   `json:"actionOnText"`
	BlockUUID              string                   `json:"blockUUID"`
	BlockThreadUUID        string                   `json:"blockThreadUUID"`
	BlockThreadCommentUUID string                   `json:"blockThreadCommentUUID"`
	ReelUUID               string                   `json:"reelUUID"`
	ReelID                 uint64                   `json:"reelId"`
	ReelCommentUUID        string                   `json:"reelCommentUUID"`
	EmailSubject           string                   `json:"emailSubject"`
}

type MessageBtnComponent struct {
	Type     int    `json:"type"`
	Label    string `json:"label"`
	Style    int    `json:"style"`
	Url      string `json:"url,omitempty"`
	CustomID string `json:"custom_id,omitempty"`
}

type ActionRowsComponent struct {
	Type       int           `json:"type"`
	Components []interface{} `json:"components"`
}

func PostNotificationInstance(event string, createdByID uint64, notifierIDs []uint64, studioID *uint64, roleIDs *[]uint64, extraData NotificationExtraData, objectID *uint64, contentObject *string) PostNotification {
	return PostNotification{
		Event:         event,
		CreatedByID:   createdByID,
		NotifierIDs:   notifierIDs,
		StudioID:      studioID,
		RoleIDs:       roleIDs,
		ExtraData:     extraData,
		ObjectID:      objectID,
		ContentObject: contentObject,
	}
}

// StudioData for the NotificationCount Studio
type StudioData struct {
	Type        string `json:"type"`
	ID          uint64 `json:"id"`
	DisplayName string `json:"displayName"`
	Handle      string `json:"handle"`
	ImageURL    string `json:"imageUrl"`
}

//StudioNotificationCountView studio notification count view
type StudioNotificationCountView struct {
	Count  int        `json:"count"`
	Studio StudioData `json:"studio"`
}

// BuildStudio for notification count view out of studio
func BuildStudio(p *models.Studio) StudioData {
	view := StudioData{
		Type:        "studio",
		ID:          p.ID,
		DisplayName: p.DisplayName,
		Handle:      p.Handle,
		ImageURL:    p.ImageURL,
	}
	return view
}

type studioCount struct {
	Count    int
	StudioID uint64
}

type StudioCount struct {
	Count  int `json:"count"`
	Studio struct {
		Type        string `json:"type"`
		Id          int    `json:"id"`
		DisplayName string `json:"displayName"`
		Handle      string `json:"handle"`
		ImageUrl    string `json:"imageUrl"`
	} `json:"studio"`
}

type NotificationCountSerializer struct {
	ID       uint64        `json:"id"`
	UserID   uint64        `json:"userId"`
	All      int64         `json:"all"`
	Personal int64         `json:"personal"`
	Studio   []StudioCount `json:"studio"`
}

func SerializeNotificationCount(notificationCount models.NotificationCount) NotificationCountSerializer {
	view := NotificationCountSerializer{
		ID:       notificationCount.ID,
		UserID:   notificationCount.UserID,
		All:      notificationCount.All,
		Personal: notificationCount.Personal,
	}
	var studios []StudioCount
	json.Unmarshal([]byte(notificationCount.Studio), &studios)
	view.Studio = studios
	return view
}

type NotificationGetSerializer struct {
	ID                    uint64         `json:"id"`
	UUID                  string         `json:"uuid"`
	CreatedAt             time.Time      `json:"createdAt"`
	UpdatedAt             time.Time      `json:"updatedAt"`
	Entity                string         `json:"entity"`
	Event                 string         `json:"event"`
	Activity              string         `json:"activity"`
	Priority              string         `json:"priority"`
	Text                  string         `json:"text"`
	NotifierID            uint64         `json:"notifierId"`
	ReactorID             *uint64        `json:"reactorId"`
	StudioID              *uint64        `json:"studioId"`
	ObjectId              *uint64        `json:"objectId"`
	ContentObject         *string        `json:"contentObject"`
	TargetObjectId        *uint64        `json:"targetObjectId"`
	TargetContentObject   *string        `json:"targetContentObject"`
	ExtraData             datatypes.JSON `json:"extraData"`
	Read                  bool           `json:"read"`
	Seen                  bool           `json:"seen"`
	IsPersonal            bool           `json:"isPersonal"`
	DiscordDmID           *string        `json:"discordDmId"`
	IsWebsiteNotification bool           `json:"isWebsiteNotification"`
	CreatedByID           uint64         `json:"createdById"`
	UpdatedByID           uint64         `json:"updatedById"`
	ArchivedByID          uint64         `json:"archivedById"`
	IsArchived            bool           `json:"isArchived"`
	ArchivedAt            time.Time      `json:"archivedAt"`
}

func SerializeNotification(notification models.Notification) NotificationGetSerializer {
	view := NotificationGetSerializer{
		ID:                    notification.ID,
		UUID:                  notification.UUID.String(),
		CreatedAt:             notification.CreatedAt,
		UpdatedAt:             notification.UpdatedAt,
		Entity:                notification.Entity,
		Event:                 notification.Event,
		Activity:              notification.Activity,
		Priority:              notification.Priority,
		Text:                  notification.Text,
		NotifierID:            notification.NotifierID,
		StudioID:              notification.StudioID,
		ObjectId:              notification.ObjectId,
		ContentObject:         notification.ContentObject,
		TargetObjectId:        notification.TargetObjectId,
		TargetContentObject:   notification.TargetContentObject,
		ExtraData:             notification.ExtraData,
		Read:                  notification.Read,
		Seen:                  notification.Seen,
		IsPersonal:            notification.IsPersonal,
		DiscordDmID:           notification.DiscordDmID,
		IsWebsiteNotification: notification.IsWebsiteNotification,
		CreatedByID:           notification.CreatedByID,
		UpdatedByID:           notification.UpdatedByID,
		ArchivedByID:          notification.ArchivedByID,
		IsArchived:            notification.IsArchived,
		ArchivedAt:            notification.ArchivedAt,
	}
	return view
}

func MultiSerializeNotification(notifications []models.Notification) []NotificationGetSerializer {
	notificationsData := []NotificationGetSerializer{}
	for _, notification := range notifications {
		data := SerializeNotification(notification)
		notificationsData = append(notificationsData, data)
	}
	return notificationsData
}

type MentionsSerializer struct {
	Type      string `json:"type"`
	ID        uint64 `json:"id"`
	UUID      string `json:"uuid"`
	FullName  string `json:"fullName"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatarUrl"`
	Key       string `json:"key"`
	RepoID    uint64 `json:"repoID"`
	RepoKey   string `json:"repoKey"`
	RepoName  string `json:"repoName"`
	RepoUUID  string `json:"repoUUID"`
}

type NotificationEmailData struct {
	Text     string `json:"text"`
	AppUrl   string `json:"appUrl"`
	Activity string `json:"activity"`
	Subject  string `json:"subject"`
	Event    string `json:"event"`
}

type BlockChildren struct {
	Text string `json:"text"`
	Type string `json:"type"`
	Url  string `json:"url"`
}

type BlockData struct {
	ID       uint64
	UUID     string
	Type     string
	Children []struct {
		Text string
		Type string
		Url  string
	}
}

type CommentsData struct {
	Text string `json:"text"`
}

type BlockAttributes struct {
	MessageID string `json:"messageId"`
}

type SlackSocialAuthMetadata struct {
	AccessToken string `json:"accessToken"`
}

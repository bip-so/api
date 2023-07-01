package models

const (
	NOTIFICATION_TYPE_DISCORD = "discord"
	NOTIFICATION_TYPE_APP     = "app"
	NOTIFICATION_TYPE_EMAIL   = "email"
	NOTIFICATION_TYPE_SLACK   = "slack"
)

var DefaultNotificationTypes = []string{"app", "email"}

func (m *UserSettings) TableName() string {
	return "user_settings"
}

// @todo in future needed notifications specific for canvas

type UserSettings struct {
	BaseModel

	UserID                  uint64
	Type                    string
	AllComments             bool `gorm:"default:false"`
	RepliesToMe             bool `gorm:"default:false"`
	Mentions                bool `gorm:"default:false"`
	Reactions               bool `gorm:"default:false"`
	Invite                  bool `gorm:"default:false"`
	FollowedMe              bool `gorm:"default:false"`
	FollowedMyStudio        bool `gorm:"default:false"`
	PublishAndMergeRequests bool `gorm:"default:false"`
	ResponseToMyRequests    bool `gorm:"default:false"`
	SystemNotifications     bool `gorm:"default:false"`
	DarkMode                bool `gorm:"default:false"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

func NewDefaultUserSettings(userID uint64, notificationType string) *UserSettings {
	if notificationType == NOTIFICATION_TYPE_APP {
		return &UserSettings{
			UserID:                  userID,
			Type:                    notificationType,
			AllComments:             true,
			RepliesToMe:             true,
			Mentions:                true,
			Reactions:               true,
			Invite:                  true,
			FollowedMe:              true,
			FollowedMyStudio:        true,
			PublishAndMergeRequests: true,
			ResponseToMyRequests:    true,
		}
	} else if notificationType == NOTIFICATION_TYPE_DISCORD {
		return &UserSettings{
			UserID:                  userID,
			Type:                    notificationType,
			AllComments:             true,
			RepliesToMe:             true,
			Mentions:                true,
			Reactions:               false,
			Invite:                  true,
			FollowedMe:              false,
			FollowedMyStudio:        false,
			PublishAndMergeRequests: true,
			ResponseToMyRequests:    true,
		}
	} else if notificationType == NOTIFICATION_TYPE_EMAIL {
		return &UserSettings{
			UserID:                  userID,
			Type:                    notificationType,
			AllComments:             true,
			RepliesToMe:             true,
			Mentions:                true,
			Reactions:               false,
			Invite:                  true,
			FollowedMe:              false,
			FollowedMyStudio:        false,
			PublishAndMergeRequests: true,
			ResponseToMyRequests:    true,
		}
	} else if notificationType == NOTIFICATION_TYPE_SLACK {
		return &UserSettings{
			UserID:                  userID,
			Type:                    notificationType,
			AllComments:             true,
			RepliesToMe:             true,
			Mentions:                true,
			Reactions:               false,
			Invite:                  true,
			FollowedMe:              false,
			FollowedMyStudio:        false,
			PublishAndMergeRequests: true,
			ResponseToMyRequests:    true,
		}
	}
	return &UserSettings{
		UserID: userID,
		Type:   notificationType,
	}
}

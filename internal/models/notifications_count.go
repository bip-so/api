package models

type NotificationCount struct {
	BaseModel
	UserID   uint64
	All      int64  `gorm:"default:0"`
	Personal int64  `gorm:"default:0"`
	Studio   string `gorm:"type:json"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

func NewNotificationCount(userID uint64, all int64, personal int64, studioCount string) *NotificationCount {
	return &NotificationCount{
		UserID:   userID,
		All:      all,
		Personal: personal,
		Studio:   studioCount,
	}
}

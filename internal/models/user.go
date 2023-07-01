package models

import (
	"database/sql"
	"time"
)

const RedisUserOtpNS = "userotp:"
const NewUsers = "users:"
const RedisUserForgotPassword = "user-forgot-pass:"
const GhostLogin = "ghost-login:"

func (m *User) TableName() string {
	return "users"
}

type User struct {
	BaseModel

	Email     sql.NullString `gorm:"index"`
	Password  string
	LastLogin *time.Time
	Username  string `gorm:"uniqueIndex"`
	// Commented for now: Firstname and Last name/
	FullName        string
	AvatarUrl       string
	IsSuperuser     bool `gorm:"default:false"`
	IsEmpty         bool `gorm:"default:false"`
	IsSetupDone     bool `gorm:"default:false"`
	IsEmailVerified bool `gorm:"default:false"`
	HasPassword     bool `gorm:"default:true"`
	Timezone        sql.NullString
	DateJoined      time.Time `gorm:"autoCreateTime"`

	// Not adding any FK for a reason.
	DefaultStudioID uint64       `gorm:"default:0"`
	UserProfile     *UserProfile `gorm:"foreignkey:UserID;references:ID"`

	ViaDiscordStore   bool   `gorm:"default:false"`
	ClientReferenceId string `gorm:"default:na"`
}

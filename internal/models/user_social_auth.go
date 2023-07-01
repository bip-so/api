package models

import "gorm.io/datatypes"

func (m *UserSocialAuth) TableName() string {
	return "user_social_auths"
}

type UserSocialAuth struct {
	BaseModel
	UserID       uint64
	ProviderName string
	ProviderID   string
	Metadata     datatypes.JSON
	User         *User `gorm:"foreignkey:UserID;contraint;OnDelete:CASCADE;"`
}

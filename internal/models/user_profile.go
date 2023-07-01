package models

func (m *UserProfile) TableName() string {
	return "user_profiles"
}

type UserProfile struct {
	UserID     uint64
	Bio        string
	Website    *string
	TwitterUrl *string
	Location   *string
	User       *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

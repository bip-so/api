package models

func (m *UserContact) TableName() string {
	return "user_contacts"
}

//UserContact to store users contact list or user favourites
type UserContact struct {
	BaseModel

	UserID        uint64
	ContactUserID *uint64
	Email         string
	Phone         string
	Photo         string
	Name          string
	Deleted       bool

	User        User  `gorm:"foreignkey:UserID;constraint:OnDelete:CASCADE;"`
	ContactUser *User `gorm:"foreignkey:ContactUserID;constraint:OnDelete:CASCADE;"`
}

func (user User) NewUserContact(contactUserID *uint64, phone string, email string, name string, photo string) UserContact {
	return UserContact{
		UserID:        user.ID,
		ContactUserID: contactUserID,
		Email:         email,
		Phone:         phone,
		Name:          name,
		Photo:         photo,
		Deleted:       false,
	}
}

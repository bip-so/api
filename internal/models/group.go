package models

import (
	"time"
)

/* Groups are a generic way of categorizing users to apply permissions, or
some other label, to those users. A user can belong to any number of
groups. */
func (m *Group) TableName() string {
	return "groups"
}

type Group struct {
	BaseModel
	Name  string `gorm:"type: varchar(150)"`
	Users []User

	CreatedByID uint64
	UpdatedByID uint64

	CreatedByUser *User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	IsArchived     bool
	ArchivedAt     time.Time
	ArchivedByID   *uint64
	ArchivedByUser *User `gorm:"foreignKey:ArchivedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

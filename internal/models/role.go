package models

import (
	"database/sql"
	"time"
)

const SYSTEM_ROLE_MEMBER = "Member"

//const SYSTEM_ADMIN_ROLE = "BIP Admin"
const SYSTEM_ADMIN_ROLE = "Administrator"
const SYSTEM_BILLING_ROLE = "Billing"

func (m *Role) TableName() string {
	return "roles"
}

type Role struct {
	// Roles always belong to a studio
	BaseModel
	// BoilerPlateStuff
	CreatedByID *uint64
	UpdatedByID *uint64

	IsArchived   bool
	ArchivedAt   time.Time
	ArchivedByID uint64

	// This should be almost always required
	StudioID uint64 // Studio to which roles belongs
	Name     string `gorm:"type: varchar(50)"`
	Color    string `gorm:"type: varchar(16)"` // Hex value
	IsSystem bool   `gorm:"default:false"`
	Icon     string

	// Roles also have many members
	Members    []Member `gorm:"many2many:role_members;"`
	IsNonPerms bool     `gorm:"default:false"`

	DiscordRoleID sql.NullString
	SlackRoleID   sql.NullString

	Studio         *Studio `gorm:"foreignKey:StudioID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedByUser  *User   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser  *User   `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ArchivedByUser *User   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type RoleMember struct {
	MemberID uint64
	RoleID   uint64
}

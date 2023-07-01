package models

import (
	"regexp"
	"time"

	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

func (m *Studio) TableName() string {
	return "studios"
}

type Studio struct {
	BaseModel
	CreatedByID uint64
	UpdatedByID uint64

	CreatedByUser *User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	IsArchived     bool
	ArchivedAt     time.Time
	ArchivedByID   *uint64
	ArchivedByUser *User `gorm:"foreignKey:ArchivedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	DisplayName string `gorm:"type: varchar(100)"`
	Handle      string `gorm:"type: varchar(60)"`
	Description string `gorm:"type: varchar(300)"`
	Website     string
	ImageURL    string

	DiscordNotificationsEnabled bool `gorm:"default:false"`
	SlackNotificationsEnabled   bool `gorm:"default:false"`

	ComputedFollowerCount int

	// Membership Allowed
	AllowPublicMembership bool `gorm:"default:false"`

	// Studio Plans
	IsEarlyAdopter bool `gorm:"default:false"`
	IsNonProfit    bool `gorm:"default:false"`

	// Stripe IDs
	StripeCustomerID      string `gorm:"default:na"`
	StripeProductID       string `gorm:"default:na"`
	StripePriceID         string `gorm:"default:na"`
	StripeSubscriptionsID string `gorm:"default:na"`
	StripePriceUnit       int64  `gorm:"default:0"`

	Topics []Topic `gorm:"many2many:studio_topics;"`
	// Update VendorName post success -> SR
	VendorName       string `gorm:"default:na"`
	FeatureFlagHasXP bool   `gorm:"default:false"`
	// Vendors

}

//
//func (s *Studio) AfterSave(tx *gorm.DB) (err error) {
//	fmt.Println("Studio Just got Created.")
//	fmt.Println(s.DisplayName)
//	fmt.Println(s.Handle)
//	fmt.Println(s.ID)
//	noneRole := &Role{
//		Name:     "NONE",
//		StudioID: s.ID,
//		Color:    "#ffffff",
//		IsSystem: true,
//		Icon:     "",
//	}
//	err = postgres.GetDB().Create(noneRole).Error
//	fmt.Println("Error creating NONE ROLE")
//	fmt.Println(err)
//	return
//}

func CheckHandleValidity(handle string) bool {
	regex := regexp.MustCompile("^[a-zA-Z0-9]+$")
	return regex.Match([]byte(handle))
}

func (studio *Studio) TopicsDiff(newTopics []string) (addedTopics []string, removedTopics []string) {
	oldTopics := []string{}
	for _, topic := range studio.Topics {
		oldTopics = append(oldTopics, topic.Name)
	}
	processedTopics := []string{}
	for _, topic := range newTopics {
		if !utils.SliceContainsItem(oldTopics, topic) {
			addedTopics = append(addedTopics, topic)
		} else {
			processedTopics = append(processedTopics, topic)
		}
	}
	for _, topic := range oldTopics {
		if !utils.SliceContainsItem(processedTopics, topic) {
			removedTopics = append(removedTopics, topic)
		}
	}
	return
}

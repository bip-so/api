package models

const RedisFollowUserStudioNS = "followstudio:"

func (m *FollowStudio) TableName() string {
	return "follow_studios"
}

type FollowStudio struct {
	BaseModel
	StudioId   uint64 `gorm:"index:idx_studioid_followerid,unique"`
	FollowerId uint64 `gorm:"index:idx_studioid_followerid,unique"`

	Studio       *Studio `gorm:"foreignKey:StudioId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FollowerUser *User   `gorm:"foreignKey:FollowerId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

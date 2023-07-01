package models

const RedisFollowUserNS = "followuser:"

func (m *FollowUser) TableName() string {
	return "follow_users"
}

/*
 For getting people following this USER
	Query with "user_id": user.ID : This is followers list
 For getting people this USER is following
	Query with "follower_id": user.ID This will list Users this user is following.

*/

type FollowUser struct {
	BaseModel

	//UserID -> Hero
	//FollowerID -> Following
	//FollowerID is Following UserID

	UserId     uint64 `gorm:"index:idx_userid_followerid,unique"`
	FollowerId uint64 `gorm:"index:idx_userid_followerid,unique"` // user

	User         *User `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FollowerUser *User `gorm:"foreignKey:FollowerId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

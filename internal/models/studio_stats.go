package models

type StudioStats struct {
	BaseModel
	// Stats
	CollectionCount        int64
	PrivateCanvasRepoCount int64
	PublicCanvasRepoCount  int64
	StorageObjectCount     int64
	StorageObjectSizeCount int64
	MemberCount            int64
	PostCount              int64
	ReelCount              int64

	StudioID uint64
	Studio   *Studio `gorm:"foreignKey:StudioID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

const StudioStatsCollectionCount = "collections"
const StudioStatsPrivateCanvasRepoCount = "private_repo"
const StudioStatsPublicCanvasRepoCount = "public_repo"
const StudioStatsStorageObjectCount = "storage_count"
const StudioStatsStorageObjectSizeCount = "storage_size"
const StudioStatsMemberCount = "members_count"
const StudioStatsPostCount = "post_count"
const StudioStatsReelCount = "reel_count"

// @todo: Need to update Fields POST migrations
var StudioStatsCounterMap = map[string]string{
	StudioStatsCollectionCount:        "wrong",
	StudioStatsPrivateCanvasRepoCount: "wrong",
	StudioStatsPublicCanvasRepoCount:  "wrong",
	StudioStatsStorageObjectCount:     "wrong",
	StudioStatsStorageObjectSizeCount: "wrong",
	StudioStatsMemberCount:            "wrong",
	StudioStatsPostCount:              "wrong",
	StudioStatsReelCount:              "wrong",
}

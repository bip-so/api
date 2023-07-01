package canvasrepo

import (
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	bipredis "gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"time"
)

type NewQueueItem struct {
	StudioID uint64 `json:"studioID"`
	BranchID uint64 `json:"branchID"`
	RepoID   uint64 `json:"repoID"`
	UserID   uint64 `json:"userID"`
}

const AutoPublishQ = "plans-auto-publish:"

func (crc *repoAutoPublichCachingService) AddToAutoPublishQueue(branchID, repoID, studioID, userID uint64) {
	currentTime := time.Now()
	TodayDateAsKey := currentTime.Format("2006-01-02") //MM-DD-YYYY
	var data NewQueueItem
	data.BranchID = branchID
	data.StudioID = studioID
	data.UserID = userID
	data.RepoID = repoID
	repoDataJson, _ := json.Marshal(data)
	ctx := redis.GetBgContext()
	options := bipredis.Options{
		Expiration: 48 * time.Hour,
	}
	crc.cache.Set(ctx, fmt.Sprintf("%s%s:%s", AutoPublishQ, TodayDateAsKey, utils.String(branchID)), repoDataJson, &options)
}

func (crc *repoAutoPublichCachingService) ProcessCanvasBranchAccess(queueItemData string) {
	var data NewQueueItem
	json.Unmarshal([]byte(queueItemData), &data)
	_ = UpdateCanvasBranchVisibility(data.BranchID, data.UserID, "view")
	_ = UpdateCanvasLanguageBranchesVisibility(data.BranchID, data.UserID, "view")
}

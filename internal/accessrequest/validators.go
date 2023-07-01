package ar

import (
	"errors"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

type ManageAccessRequestPost struct {
	Status                      string `json:"status"`
	CanvasBranchPermissionGroup string `json:"canvasBranchPermissionGroup"`
}

func (obj ManageAccessRequestPost) Validate() error {
	allowedScope := []string{models.PGCanvasNoneSysName, models.PGCanvasViewSysName, models.PGCanvasCommentSysName, models.PGCanvasEditSysName, models.PGCanvasModerateSysName}
	if !utils.SliceContainsItem(allowedScope, obj.CanvasBranchPermissionGroup) {
		return errors.New("Please send corrent perms for CB.")
	}
	allowedStatus := []string{models.ACCESS_REQUEST_PENDING, models.ACCESS_REQUEST_ACCEPTED, models.ACCESS_REQUEST_REJECTED}
	if !utils.SliceContainsItem(allowedStatus, obj.Status) {
		return errors.New("Please send corrent perms for CB.")
	}

	return nil
}

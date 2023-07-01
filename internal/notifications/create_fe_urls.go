package notifications

import (
	"fmt"
	"net/url"
	"strings"

	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

func (s notificationService) GenerateStudioUrl(studioHandle string) string {
	return fmt.Sprintf("%s/%s", configs.GetAppInfoConfig().FrontendHost, studioHandle)
}

func (s notificationService) GenerateCollectionUrl(studioHandle string, collectionID uint64) string {
	return fmt.Sprintf("%s/%s/collection/%d", configs.GetAppInfoConfig().FrontendHost, studioHandle, collectionID)
}

func (s notificationService) GenerateStudioIntegrationSettingsUrl(studioHandle string) string {
	return fmt.Sprintf("%s/%s/about?open_settings=true&tab=3", configs.GetAppInfoConfig().FrontendHost, studioHandle)
}

func (s notificationService) GenerateStudioPendingRequestsUrl(studioHandle string) string {
	return fmt.Sprintf("%s/%s/about?open_settings=true&tab=4", configs.GetAppInfoConfig().FrontendHost, studioHandle)
}

func (s notificationService) GenerateStudioBillingUrl(studioHandle string) string {
	return fmt.Sprintf("%s/%s/about?open_settings=true", configs.GetAppInfoConfig().FrontendHost, studioHandle)
}

func (s notificationService) GenerateCanvasBranchUrl(canvasKey, canvasName string, studioId, canvasBranchId uint64) string {
	studio, _ := App.Repo.GetStudioByID(studioId)
	// return fmt.Sprintf("%s/@%s/canvas/%s/%s/%s", configs.GetAppInfoConfig().FrontendHost, studio.Handle, canvasKey, utils.String(canvasBranchId), url.QueryEscape(canvasName))
	urlTitle := s.GenerateCanvasUrlTitle(canvasName, canvasBranchId)
	return fmt.Sprintf("%s/%s/%s", configs.GetAppInfoConfig().FrontendHost, studio.Handle, url.QueryEscape(urlTitle))
}

func (s notificationService) GenerateCanvasBranchUrlByID(canvasBranchId uint64) string {
	canvasBranch, _ := App.Repo.GetCanvasBranchByIDPreload(canvasBranchId)
	studio, _ := App.Repo.GetStudioByID(canvasBranch.CanvasRepository.StudioID)
	// return fmt.Sprintf("%s/@%s/canvas/%s/%s/%s", configs.GetAppInfoConfig().FrontendHost, studio.Handle, canvasKey, utils.String(canvasBranchId), url.QueryEscape(canvasName))
	urlTitle := s.GenerateCanvasUrlTitle(canvasBranch.CanvasRepository.Name, canvasBranchId)
	return fmt.Sprintf("%s/%s/%s", configs.GetAppInfoConfig().FrontendHost, studio.Handle, url.QueryEscape(urlTitle))
}

func (s notificationService) GenerateCanvasBranchBlockUrl(canvasKey, canvasName string, studioId, canvasBranchId uint64, blockUUID string) string {
	studio, _ := App.Repo.GetStudioByID(studioId)
	// return fmt.Sprintf("%s/@%s/canvas/%s/%s/%s", configs.GetAppInfoConfig().FrontendHost, studio.Handle, canvasKey, utils.String(canvasBranchId), url.QueryEscape(canvasName))
	urlTitle := s.GenerateCanvasUrlTitle(canvasName, canvasBranchId)
	return fmt.Sprintf("%s/%s/%s?blockUUID=%s", configs.GetAppInfoConfig().FrontendHost, studio.Handle, url.QueryEscape(urlTitle), blockUUID)
}

func (s notificationService) GenerateMergeRequestUrl(canvasKey, canvasName string, studioId, canvasBranchId, mergeReqID uint64) string {
	studio, _ := App.Repo.GetStudioByID(studioId)
	// return fmt.Sprintf("%s/@%s/canvas/%s/%s/%s/merge-req/%s", configs.GetAppInfoConfig().FrontendHost, studio.Handle, canvasKey, utils.String(canvasBranchId), url.QueryEscape(canvasName), utils.String(mergeReqID))
	urlTitle := s.GenerateCanvasUrlTitle(canvasName, canvasBranchId)
	return fmt.Sprintf("%s/%s/%s/merge-req/%s", configs.GetAppInfoConfig().FrontendHost, studio.Handle, url.QueryEscape(urlTitle), utils.String(mergeReqID))
}

func (s notificationService) GenerateBlockCommentUrl(canvasKey, canvasName, blockThreadUUID string, studioId, canvasBranchId uint64) string {
	studio, _ := App.Repo.GetStudioByID(studioId)
	urlTitle := s.GenerateCanvasUrlTitle(canvasName, canvasBranchId)
	return fmt.Sprintf("%s/%s/%s?threadUUID=%s", configs.GetAppInfoConfig().FrontendHost, studio.Handle, url.QueryEscape(urlTitle), blockThreadUUID)
}

func (s notificationService) GenerateBlockReactionUrl(canvasKey, canvasName, blockUUID string, studioId, canvasBranchId uint64) string {
	studio, _ := App.Repo.GetStudioByID(studioId)
	urlTitle := s.GenerateCanvasUrlTitle(canvasName, canvasBranchId)
	return fmt.Sprintf("%s/%s/%s?reactionBlockUUID=%s", configs.GetAppInfoConfig().FrontendHost, studio.Handle, url.QueryEscape(urlTitle), blockUUID)
}

func (s notificationService) GenerateReelCommentUrl(canvasKey, canvasName string, studioId, canvasBranchId uint64, reelUUID string) string {
	studio, _ := App.Repo.GetStudioByID(studioId)
	urlTitle := s.GenerateCanvasUrlTitle(canvasName, canvasBranchId)
	return fmt.Sprintf("%s/%s/%s?reelUUID=%s", configs.GetAppInfoConfig().FrontendHost, studio.Handle, url.QueryEscape(urlTitle), reelUUID)
}

func (s notificationService) GenerateReelUrl(canvasKey, canvasName string, studioId, canvasBranchId uint64, reelUUID string) string {
	studio, _ := App.Repo.GetStudioByID(studioId)
	urlTitle := s.GenerateCanvasUrlTitle(canvasName, canvasBranchId)
	return fmt.Sprintf("%s/%s/%s?reelUUID=%s", configs.GetAppInfoConfig().FrontendHost, studio.Handle, url.QueryEscape(urlTitle), reelUUID)
}

func (s notificationService) GenerateCanvasUrlTitle(canvasName string, canvasBranchID uint64) string {
	canvasName = strings.ToLower(canvasName)
	canvasName = strings.ReplaceAll(canvasName, " ", "-")
	urlTitle := fmt.Sprintf("%s-%dc", canvasName, canvasBranchID)
	return urlTitle
}

func (s notificationService) GenerateReelUUIDUrl(canvasKey, canvasName string, studioId, canvasBranchId uint64, reelUUID string) string {
	studio, _ := App.Repo.GetStudioByID(studioId)
	urlTitle := s.GenerateCanvasUrlTitle(canvasName, canvasBranchId)
	return fmt.Sprintf("%s/%s/%s?reelUUID=%s", configs.GetAppInfoConfig().FrontendHost, studio.Handle, urlTitle, reelUUID)
}

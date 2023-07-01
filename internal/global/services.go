package global

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/canvasbranch"
	"gitlab.com/phonepost/bip-be-platform/internal/canvasrepo"
	"gitlab.com/phonepost/bip-be-platform/internal/collection"
	"gitlab.com/phonepost/bip-be-platform/internal/feed"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/internal/workflows"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
)

// PostStudioSetup : After a studio is created we create a collection, canvas repo, canvas branch, 1 Paragraph Block.
func PostStudioSetup(std *models.Studio, user *models.User) error {
	err := CreateDefaultDocsInStudio(std, user)
	if err != nil {
		fmt.Println("Error in creating default docs")
	}

	//collectionView := &collection.CollectionCreateValidator{
	//	Name:         "My new collection",
	//	Position:     2,
	//	PublicAccess: "private",
	//}
	//collectionInstance, err := collection.App.Controller.CreateCollectionController(collectionView, std.CreatedByID, std.ID)
	//if err != nil {
	//	return err
	//}
	//
	//canvasView := canvasrepo.InitCanvasRepoPost{
	//	CollectionID: collectionInstance.ID,
	//	Name:         "Default Canvas",
	//	Icon:         "📋",
	//	Position:     1,
	//}
	//_, errCreatingCanvasBranch := shared.WorkflowHelperInitCanvasRepo(shared.InitCanvasRepoPost{
	//	CollectionID:             canvasView.CollectionID,
	//	Name:                     canvasView.Name,
	//	Icon:                     canvasView.Icon,
	//	Position:                 canvasView.Position,
	//	ParentCanvasRepositoryID: canvasView.ParentCanvasRepositoryID,
	//}, collectionInstance.CreatedByID, collectionInstance.StudioID, *user)
	//if errCreatingCanvasBranch != nil {
	//	return errCreatingCanvasBranch
	//}
	//
	//collectionInstance, err = queries.App.CollectionQuery.UpdateCollection(
	//	collectionInstance.ID, map[string]interface{}{
	//		"computed_root_canvas_count": 1,
	//		"computed_all_canvas_count":  1,
	//	})

	go feed.App.Service.JoinStudio(std.ID, std.CreatedByID)
	// No more creating Customer on New studio
	//go payments.App.Service.CreateNewCustomerOnStripe(std.ID)
	return nil
}

func CreateDefaultDocsInStudio(std *models.Studio, user *models.User) error {
	// Create a collection
	// Add mod permission to the collection

	// Create repo `Getting Started`
	// Add mod permission to the canvas
	// move all the blocks from branch ID 34241
	// Move the comments

	// Create repo `About <<Studio Name>>`
	// Add mod permission to the canvas
	// move all the blocks from branch ID 34239
	// Move the comments

	// Create repo `Contributing or How to Contribute Guide `
	// Add mod permission to the canvas
	// move all the blocks from branch ID 34240
	// Move the comments

	// Update mentions on repo `Getting Started`
	// Update mentions on repo `About <<Studio Name>>`
	gettingStartedCanvasBranchID := uint64(6881)
	aboutStudioCanvasBranchID := uint64(6882)
	contributionCanvasBranchID := uint64(6883)
	if configs.GetConfigString("APP_MODE") == "production" {
		gettingStartedCanvasBranchID = uint64(34241)
		aboutStudioCanvasBranchID = uint64(34239)
		contributionCanvasBranchID = uint64(34240)
	}
	collectionView := &collection.CollectionCreateValidator{
		Name:         "INTRODUCTION",
		Position:     1,
		PublicAccess: "private",
	}
	collectionInstance, err := collection.App.Controller.CreateCollectionController(collectionView, std.CreatedByID, std.ID)
	if err != nil {
		fmt.Println("Error in creating collection", err)
		return err
	}

	gettingStartedCanvasView := canvasrepo.InitCanvasRepoPost{
		CollectionID: collectionInstance.ID,
		Name:         "Getting Started",
		Icon:         "",
		Position:     1,
	}
	gettingStartedCanvas, errCreatingCanvasBranch := workflows.WorkflowHelperInitCanvasRepo(workflows.InitCanvasRepoPost{
		CollectionID:             gettingStartedCanvasView.CollectionID,
		Name:                     gettingStartedCanvasView.Name,
		Icon:                     gettingStartedCanvasView.Icon,
		Position:                 gettingStartedCanvasView.Position,
		ParentCanvasRepositoryID: gettingStartedCanvasView.ParentCanvasRepositoryID,
	}, collectionInstance.CreatedByID, collectionInstance.StudioID, *user)
	if errCreatingCanvasBranch != nil {
		return errCreatingCanvasBranch
	}
	err = BlocksCloner(gettingStartedCanvasBranchID, *gettingStartedCanvas.DefaultBranchID, user, gettingStartedCanvas.ID)
	if err != nil {
		fmt.Println("Error in cloning blocks of getting started", err)
	}

	aboutStudioCanvasView := canvasrepo.InitCanvasRepoPost{
		CollectionID: collectionInstance.ID,
		Name:         "About <<Workspace Name>>",
		Icon:         "",
		Position:     2,
	}
	aboutStudioCanvas, errCreatingCanvasBranch := workflows.WorkflowHelperInitCanvasRepo(workflows.InitCanvasRepoPost{
		CollectionID:             aboutStudioCanvasView.CollectionID,
		Name:                     aboutStudioCanvasView.Name,
		Icon:                     aboutStudioCanvasView.Icon,
		Position:                 aboutStudioCanvasView.Position,
		ParentCanvasRepositoryID: aboutStudioCanvasView.ParentCanvasRepositoryID,
	}, collectionInstance.CreatedByID, collectionInstance.StudioID, *user)
	if errCreatingCanvasBranch != nil {
		return errCreatingCanvasBranch
	}
	err = BlocksCloner(aboutStudioCanvasBranchID, *aboutStudioCanvas.DefaultBranchID, user, aboutStudioCanvas.ID)
	if err != nil {
		fmt.Println("Error in cloning blocks of about studio canvas", err)
	}

	contributingCanvasView := canvasrepo.InitCanvasRepoPost{
		CollectionID: collectionInstance.ID,
		Name:         "Contributing or How to Contribute Guide",
		Icon:         "",
		Position:     3,
	}
	contributingCanvas, errCreatingCanvasBranch := workflows.WorkflowHelperInitCanvasRepo(workflows.InitCanvasRepoPost{
		CollectionID:             contributingCanvasView.CollectionID,
		Name:                     contributingCanvasView.Name,
		Icon:                     contributingCanvasView.Icon,
		Position:                 contributingCanvasView.Position,
		ParentCanvasRepositoryID: contributingCanvasView.ParentCanvasRepositoryID,
	}, collectionInstance.CreatedByID, collectionInstance.StudioID, *user)
	if errCreatingCanvasBranch != nil {
		return errCreatingCanvasBranch
	}
	err = BlocksCloner(contributionCanvasBranchID, *contributingCanvas.DefaultBranchID, user, contributingCanvas.ID)
	if err != nil {
		fmt.Println("Error in cloning blocks of contributions", err)
	}

	canvasBranchMap := map[uint64]uint64{
		// Prod mapping
		34241: *gettingStartedCanvas.DefaultBranchID,
		34239: *aboutStudioCanvas.DefaultBranchID,
		34240: *contributingCanvas.DefaultBranchID,

		// Stage mapping
		6881: *gettingStartedCanvas.DefaultBranchID,
		6882: *aboutStudioCanvas.DefaultBranchID,
		6883: *contributingCanvas.DefaultBranchID,
	}
	UpdateMentionsInBlocks(*gettingStartedCanvas.DefaultBranchID, canvasBranchMap, user, std)
	UpdateMentionsInBlocks(*aboutStudioCanvas.DefaultBranchID, canvasBranchMap, user, std)

	// Publishing the branches
	canvasbranch.App.Service.PublishCanvasBranch(*gettingStartedCanvas.DefaultBranchID, user, true)
	canvasbranch.App.Service.PublishCanvasBranch(*aboutStudioCanvas.DefaultBranchID, user, true)
	canvasbranch.App.Service.PublishCanvasBranch(*contributingCanvas.DefaultBranchID, user, true)

	collectionInstance, err = queries.App.CollectionQuery.UpdateCollection(
		collectionInstance.ID, map[string]interface{}{
			"computed_root_canvas_count": 3,
			"computed_all_canvas_count":  3,
		})
	return nil
}

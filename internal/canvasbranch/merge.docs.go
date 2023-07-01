package canvasbranch

/*

//- Publishing Canvas
//- Creating rough branch/branches
//- Getting a list of branches for a canvas & showing them
//- Creating a merge request
//- Showing merge requests
//- Merge request screen to view diff
//- Accepting/Rejecting Merge requests
//- Viewing/Accessing commit history


*/

/*

READ ME : THIS NEEDS TO BE CLEANED UP.
Blocks are cache between GIT and FE
Merge Two given Branches. Rough Branch/Branch from this : CanvasBranch has RoughFromBranchID
Main -> Forked Branch -> (Copy) -> Merge -> Main

- Branch1	Block1, Block2, Block3, Block4, Block5, Block6
- Branch2   Block1, Block2 (Delete), Block3 (Updated), Block4 (Moved), Block5, Block6, Block7 - (Create)

- CopyBlocks(toBranch, fromBranch) -> Generic Function
- TryMerge(toBranch, fromBranch)
	- Create Merge Request
	- Actual Merge
		- Merge Blocks
		- Delete Branch (RoughBranch)

Get all the blocks for the rough branch
- Branch2 (RoughBranch)   : Block1, , Block3 (Updated), Block4 (Moved), Block5, Block6, Block7 - (Create)
- Branch1 (ParentBranch) : Block1, Block2, Block3, Block4, Block5, Block6


- All Blocks ParentBranch : []
- RoughBranchBlocks: All Blocks RoughBranch: []
- Get List of all UUID's from ParentBranch
- for RoughBranchBlocks
	- Blocks to update (UUID found in ParentBranch and UUID RoughBranch)
	- Blocks to delete the (UUID's only on ParentBranch : Deleted Blocks)
    - Blocks to create: (UUID's not found in the ParentBranch)

Preserve Block Artifacts
	- Block Mentions
	- Block Threads
		- Block Thread Mentions
		- Block Comments
		- Block Comments Mentions
	- Reels
	- Reel Mention
		- Reel Comments
		- Reel Comment Mention

// Get Canvas Branch Instance
	// Validate: if the Branch is actually a RoughBranch (optional later)
	// Validate: Check U "CANVAS_BRANCH_CREATE_MERGE_REQUEST" On Rough for event allowing this action 400
	// Validate : CANVAS_BRANCH_MANAGE_MERGE_REQUESTS on Parent 1 Service 2
	// Validate : CANVAS_BRANCH_MANAGE_MERGE_REQUESTS on Parent 0 Service 1
	// Get users permissions on the CanvasBranch -> Merge Permissions Check
	// Check rough Branch permission "
	// Can the user perform the Merge request ?
	// NR: Check ParentBranch permission / Manage merge request
	// Service1 -> Create Merge Request Instance -> Return Response (Merge Request has been created)
	// Merge request has been Reqy we a counbt
	// After Service1 is Successfully Done -> Change the Rough Branch to be (readonly).Committed = True

	// Get the from instance RoughFromBranchID
	// Service2 -> Start the merge activity
	// Documentation Above.
	// After Service2 is Successfully Done -> Delete the branch and delete all the blocks on RoughBlocks

---- BLOCK LEVEL MERGE ----

//- Blocks to delete the (UUID's only on ParentBranch : Deleted Blocks)

	//- Blocks to create: (UUID's not found in the ParentBranch)

	// We need to create 3 Data Structures for following
	// + DELETE
	// Blocks to be deleted in the Parent Branch - QQ: Delete is fine; What happens to Threads/Comments  and Reels Here

	// + CREATE
	// Blocks to be added (created) on Parent Branch  :
	// Create We will Just Change the (CanvasBranchID)Branch on the Block so all other associations are remaining
	// Reset all the other fields including RoughBranch to Nil
	// + UPDATE
	// Blocks to be Updated on Parent Branch
	// This is similar to create, but we'll do the following
	// START MOVING THINGS FROM ROUGH -> CHILDREN
	// based on the SAME BLOCK on PARENT BRANCH
	// We'll get
	// - BlockThreads and change CanvasRepositoryID/CanvasBranchID/StartBlockID = RoughBranchBlockID : Rerun to cal CommentCount
	// - (REEEL) DO WE EVEN NEED TO CHECK REELS ON ROUGH BRANCH:
	//  Reel: We have to be careful as it requires
	//  CanvasRepositoryID uint64
	//	CanvasBranchID
	//	StartBlockID
	//	StartBlockUUID
	// BlockCommentS No Change since they are connected to BlockThreads
	// REEL COMMENT WILL COME AUTOMATICALLY
	// REACTIONS ->
	// BlockReaction

	// 3 way merge
	// Canvas Branch A
	// - Rough NR
	// Canvas Branch A
	// - Rough NR
	// Call API snapshot in final version
	// List blocks that shoud exist (Post Merge)
	// POST Merge

*/

package permissiongroup

// STUDIO_CREATE_DELETE_ROLE Studio Permisions
const STUDIO_CREATE_DELETE_ROLE = "STUDIO_CREATE_DELETE_ROLE"
const STUDIO_EDIT_STUDIO_PROFILE = "STUDIO_EDIT_STUDIO_PROFILE"
const STUDIO_ADD_REMOVE_USER_TO_ROLE = "STUDIO_ADD_REMOVE_USER_TO_ROLE"
const STUDIO_MANAGE_INTEGRATION = "STUDIO_MANAGE_INTEGRATION"
const STUDIO_CREATE_COLLECTION = "STUDIO_CREATE_COLLECTION"
const STUDIO_METADATA_UPDATE = "STUDIO_METADATA_UPDATE"
const STUDIO_CHANGE_CANVAS_COLLECTION_POSITION = "STUDIO_CHANGE_CANVAS_COLLECTION_POSITION"
const STUDIO_DELETE = "STUDIO_DELETE"
const STUDIO_MANAGE_PERMS = "STUDIO_MANAGE_PERMS"

// New studio plan perms
const STUDIO_CAN_MANAGE_BILLING = "CAN_MANAGE_BILLING"

// COLLECTION_MEMBERSHIP_MANAGE Collection Permisions
// @todo COLLECTION_MEMBERSHIP_MANAGE, COLLECTION_PUBLIC_ACCESS_CHANGE can be merged into COLLECTION_MANAGE_PERMS here
const COLLECTION_MEMBERSHIP_MANAGE = "COLLECTION_MEMBERSHIP_MANAGE"
const COLLECTION_PUBLIC_ACCESS_CHANGE = "COLLECTION_PUBLIC_ACCESS_CHANGE"
const COLLECTION_MANAGE_PERMS = "COLLECTION_MANAGE_PERMS"
const COLLECTION_DELETE = "COLLECTION_DELETE"
const COLLECTION_OVERRIDE_STUDIO_MODE_ROLE = "COLLECTION_OVERRIDE_STUDIO_MODE_ROLE"
const COLLECTION_EDIT_NAME = "COLLECTION_EDIT_NAME"
const COLLECTION_VIEW_METADATA = "COLLECTION_VIEW_METADATA"

// This is added on 4Th July. Mainly for Future not User as of now.
// Permission is added for Future: Relevent for When PR is being requested

const COLLECTION_MANAGE_PUBLISH_REQUEST = "COLLECTION_MANAGE_PUBLISH_REQUEST"

const CANVAS_BRANCH_VIEW = "CANVAS_BRANCH_VIEW"
const CANVAS_BRANCH_EDIT = "CANVAS_BRANCH_EDIT"
const CANVAS_BRANCH_EDIT_NAME = "CANVAS_BRANCH_EDIT_NAME"
const CANVAS_BRANCH_DELETE = "CANVAS_BRANCH_DELETE"
const CANVAS_BRANCH_ADD_COMMENT = "CANVAS_BRANCH_ADD_COMMENT"
const CANVAS_BRANCH_ADD_REACTION = "CANVAS_BRANCH_ADD_REACTION"
const CANVAS_BRANCH_CREATE_REEL = "CANVAS_BRANCH_CREATE_REEL"
const CANVAS_BRANCH_COMMENT_ON_REEL = "CANVAS_BRANCH_COMMENT_ON_REEL"
const CANVAS_BRANCH_REACT_TO_REEL = "CANVAS_BRANCH_REACT_TO_REEL"
const CANVAS_BRANCH_MANAGE_PERMS = "CANVAS_BRANCH_MANAGE_PERMS"

// Can a user manage a merge request  checked while merging.
// CHECK COLLECTION -> CANVAS REPO -> "MAIN"

const CANVAS_BRANCH_MANAGE_MERGE_REQUESTS = "CANVAS_BRANCH_MANAGE_MERGE_REQUESTS"

// "CANVAS_BRANCH_MANAGE_PUBLISH_REQUESTS" on Parent -> Main Pe ()
// NR -> cc can moderate -> Manage Publish Request
//
// Permission is added for Future: Relevent for When PR is being requested
const CANVAS_BRANCH_MANAGE_PUBLISH_REQUESTS = "CANVAS_BRANCH_MANAGE_PUBLISH_REQUESTS"

// Can a user create a merge request - Checked while merging.
const CANVAS_BRANCH_CREATE_MERGE_REQUEST = "CANVAS_BRANCH_CREATE_MERGE_REQUEST"
const CANVAS_BRANCH_CREATE_PUBLISH_REQUEST = "CANVAS_BRANCH_CREATE_PUBLISH_REQUEST"
const CANVAS_BRANCH_CLONE = "CANVAS_BRANCH_CLONE"
const CANVAS_BRANCH_RESOLVE_COMMENTS = "CANVAS_BRANCH_RESOLVE_COMMENTS"
const CANVAS_BRANCH_MANAGE_CONTENT = "CANVAS_BRANCH_MANAGE_CONTENT"

const CANVAS_BRANCH_VIEW_METADATA = "CANVAS_BRANCH_VIEW_METADATA"

package notifications

var REQUEST_EVENTS = []string{PublishRequested, MergeRequested, AccessRequested, AccessRequestedUpdate, PublishRequestedUpdate, MergeRequestedUpdate}
var REPLIES_EVENTS = []string{BlockComment, ReelComment, CommentReply, ReelCommentReply, BlockReact, ReelReact, BlockCommentReact,
	ReelCommentReact, BlockMention, BlockThreadMention, BlockThreadCommentMention, ReelMention, ReelCommentMention}
var PR_EVENTS = []string{PublishRequested, PublishRequestedUpdate}

const RoughBranchNameSpace = "rough-branch-notifications:"

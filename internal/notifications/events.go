package notifications

const (
	BlockComment                 = "BlockComment"
	BlockReact                   = "BlockReact"
	ReelComment                  = "ReelComment"
	ReelReact                    = "ReelReact"
	ReelCommentReply             = "ReelCommentReply"
	ReelCommentReact             = "ReelCommentReact"
	CommentReply                 = "CommentReply"
	BlockCommentReact            = "BlockCommentReact"
	CanvasMerged                 = "CanvasMerged"
	TranslatedCanvas             = "TranslatedCanvas"
	JoinedStudio                 = "JoinedStudio"
	PublishRequested             = "PublishRequested"
	PublishRequestedUpdate       = "PublishRequestedUpdate"
	MergeRequested               = "MergeRequested"
	MergeRequestedUpdate         = "MergeRequestedUpdate"
	FollowUser                   = "FollowUser"
	StudioInviteByName           = "StudioInviteByName"
	StudioInviteByGroup          = "StudioInviteByGroup"
	CollectionInviteByName       = "CollectionInviteByName"
	CollectionInviteByGroup      = "CollectionInviteByGroup"
	BlockMention                 = "BlockMention"
	BlockThreadMention           = "BlockThreadMention"
	BlockThreadCommentMention    = "BlockThreadCommentMention"
	ReelMention                  = "ReelMention"
	ReelCommentMention           = "ReelCommentMention"
	AccessRequested              = "AccessRequested"
	AccessRequestedUpdate        = "AccessRequestedUpdate"
	CanvasInviteByName           = "CanvasInviteByName"
	CanvasInviteByGroup          = "CanvasInviteByGroup"
	BipMarkMessageAdded          = "BipMarkMessageAdded"
	NotionImport                 = "NotionImport"
	FileImport                   = "FileImport"
	TranslateCanvas              = "TranslateCanvas"
	DiscordIntegrationTask       = "DiscordIntegrationTask"
	DiscordIntegrationTaskFailed = "DiscordIntegrationTaskFailed"
	CreateRequestToJoinStudio    = "CreateRequestToJoinStudio"
	RejectRequestToJoinStudio    = "RejectRequestToJoinStudio"
	AcceptRequestToJoinStudio    = "AcceptRequestToJoinStudio"
	CanvasLimitExceed            = "CanvasLimitExceed"
)

var systemNotificationEvents = map[string]NotificationEvent{
	TranslatedCanvas: {
		Activity: "Translated canvas",
		Text:     "Canvas is Translated",
		Priority: "high",
	},
	CanvasMerged: {
		Activity: "Canvas Merged",
		Text:     "@%s  merged changes to `📄 %s`",
		Priority: "high",
	},
	BipMarkMessageAdded: {
		Activity: "Bip Mark Message Added",
		Text:     "Your discord message was added to `📄 %s` by `@%s`",
		Priority: "medium",
	},
	NotionImport: {
		Activity: "Notion import Completed",
		Text:     "Notion import has been successfully completed",
		Priority: "medium",
	},
	FileImport: {
		Activity: "File import Completed",
		Text:     "File import has been successfully completed",
		Priority: "medium",
	},
	TranslateCanvas: {
		Activity: "Translate canvas Completed",
		Text:     "Canvas has been successfully translated",
		Priority: "medium",
	},
	DiscordIntegrationTask: {
		Activity: "Discord Integration Task",
		Text:     "Discord Integration Completed. All your Discord members and roles added to your workspace.",
		Priority: "medium",
	},
	DiscordIntegrationTaskFailed: {
		Activity: "Discord Integration Task Failed",
		Text:     "Discord integration failed, please retry after some time or chat with us.",
		Priority: "medium",
	},
	CreateRequestToJoinStudio: {
		Activity: "Request To Join Studio",
		Text:     "%s requesting to join `%s` workspace",
		Priority: "medium",
	},
	RejectRequestToJoinStudio: {
		Activity: "Reject Request To Join Studio",
		Text:     "%s has rejected your request to join the workspace",
		Priority: "medium",
	},
	AcceptRequestToJoinStudio: {
		Activity: "Accept Request To Join Studio",
		Text:     "%s has accepted your request to join the workspace",
		Priority: "medium",
	},
	CanvasLimitExceed: {
		Activity: "Canvas Limit Exceeded",
		Text:     "",
		Priority: "high",
	},
}

var invitesEvents = map[string]NotificationEvent{
	StudioInviteByName: {
		Activity: "Studio Invite By Name",
		Text:     "@%s invited you to workspace `%s`",
		Priority: "medium",
	},
	StudioInviteByGroup: {
		Activity: "Studio Invite By Group",
		Text:     "@%s invited you to workspace `%s`",
		Priority: "medium",
	},
	CollectionInviteByName: {
		Activity: "Collection Invite By Name",
		Text:     "@%s invited you to `📄 %s`",
		Priority: "medium",
	},
	CollectionInviteByGroup: {
		Activity: "Collection Invite By Group",
		Text:     "@%s invited you to `📄 %s`",
		Priority: "medium",
	},
	CanvasInviteByName: {
		Activity: "Canvas Invite By Name",
		Text:     "@%s invited you to `📄 %s`",
		Priority: "medium",
	},
	CanvasInviteByGroup: {
		Activity: "Canvas Invite By Group",
		Text:     "@%s invited you to `📄 %s`",
		Priority: "medium",
	},
}

var allCommentsEvents = map[string]NotificationEvent{
	BlockComment: {
		Activity: "Block Comment",
		Text:     "%s commented in `📄 %s`",
		Priority: "medium",
	},
	ReelComment: {
		Activity: "Reel Comment",
		Text:     "%s commented in reel",
		Priority: "medium",
	},
}

var repliesToMeEvents = map[string]NotificationEvent{
	CommentReply: {
		Activity: "Comment Reply",
		Text:     "%s commented in `📄 %s`",
		Priority: "medium",
	},
	ReelCommentReply: {
		Activity: "Reel Comment Reply",
		Text:     "%s replied in reel comment",
		Priority: "medium",
	},
}

var reactionsEvents = map[string]NotificationEvent{
	BlockReact: {
		Activity: "Block Reacted",
		Text:     "@%s reacted with %s in `📄 %s`",
		Priority: "medium",
	},
	ReelReact: {
		Activity: "Reel Reacted",
		Text:     "@%s reacted with %s in `📄 %s`",
		Priority: "medium",
	},
	BlockCommentReact: {
		Activity: "Comment Reacted",
		Text:     "@%s reacted with %s in `📄 %s`",
		Priority: "medium",
	},
	ReelCommentReact: {
		Activity: "Reel Comment react",
		Text:     "@%s reacted with %s in `📄 %s`",
		Priority: "medium",
	},
}

var mentionsEvents = map[string]NotificationEvent{
	BlockMention: {
		Activity: "Block Mentioned",
		Text:     "@%s mentioned you in a %s",
		Priority: "medium",
	},
	BlockThreadMention: {
		Activity: "Block Mentioned",
		Text:     "@%s mentioned you in a %s",
		Priority: "medium",
	},
	BlockThreadCommentMention: {
		Activity: "Block Mentioned",
		Text:     "@%s mentioned you in a %s",
		Priority: "medium",
	},
	ReelMention: {
		Activity: "Block Mentioned",
		Text:     "@%s mentioned you in a %s",
		Priority: "medium",
	},
	ReelCommentMention: {
		Activity: "Block Mentioned",
		Text:     "@%s mentioned you in a %s",
		Priority: "medium",
	},
}

var FollowedMyStudioEvents = map[string]NotificationEvent{
	JoinedStudio: {
		Activity: "Joined Studio",
		Text:     "@%s has joined `🎨 %s`",
		Priority: "medium",
	},
}

var FollowedMeEvents = map[string]NotificationEvent{
	FollowUser: {
		Activity: "Follow User",
		Text:     "@%s is following you",
		Priority: "medium",
	},
}

var PublishAndMergeRequestsEvents = map[string]NotificationEvent{
	PublishRequested: {
		Activity: "Publish Requested",
		Text:     "@%s is requesting to publish `📄 %s` in `🎨 %s`",
		Priority: "high",
	},
	MergeRequested: {
		Activity: "Merge Requested",
		Text:     "@%s is requesting to merge changes to `📄 %s`",
		Priority: "high",
	},
	AccessRequested: {
		Activity: "Access Requested",
		Text:     "@%s is requesting access to `📄 %s`",
		Priority: "Medium",
	},
	AccessRequestedUpdate: {
		Activity: "Access Requested Update",
		Text:     "@%s has granted **%s** in `📄 %s`",
		Priority: "Medium",
	},
}

var ResponseToMyRequestsEvents = map[string]NotificationEvent{
	PublishRequestedUpdate: {
		Activity: "Publish Request update",
		Text:     "@%s **%s** your request to publish `📄 %s`",
		Priority: "high",
	},
	MergeRequestedUpdate: {
		Activity: "Merge Request updated",
		Text:     "@%s **%s** your merge request in `📄 %s`",
		Priority: "high",
	},
}

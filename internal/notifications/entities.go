package notifications

const (
	AllComments             = "AllComments"
	RepliesToMe             = "RepliesToMe"
	Mentions                = "Mentions"
	Reactions               = "Reactions"
	Invite                  = "Invite"
	FollowedMe              = "FollowedMe"
	FollowedMyStudio        = "FollowedStudio"
	PublishAndMergeRequests = "PublishAndMergeRequests"
	ResponseToMyRequests    = "ResponseToMyRequests"
	SystemNotifications     = "SystemNotifications"
)

var AllCommentsEntity = NotificationEntity{
	Entity:     AllComments,
	App:        true,
	Email:      true,
	Discord:    true,
	IsPersonal: false,
	Events:     allCommentsEvents,
}

var RepliesToMeEntity = NotificationEntity{
	Entity:     RepliesToMe,
	App:        true,
	Email:      true,
	Discord:    true,
	IsPersonal: true,
	Events:     repliesToMeEvents,
}

var ReactionsEntity = NotificationEntity{
	Entity:     Reactions,
	App:        true,
	Email:      false,
	Discord:    false,
	IsPersonal: true,
	Events:     reactionsEvents,
}

var MentionsEntity = NotificationEntity{
	Entity:     Mentions,
	App:        true,
	Email:      true,
	Discord:    true,
	IsPersonal: true,
	Events:     mentionsEvents,
}

var InvitesEntity = NotificationEntity{
	Entity:     Invite,
	App:        true,
	Email:      true,
	Discord:    true,
	IsPersonal: true,
	Events:     invitesEvents,
}

var FollowedMeEntity = NotificationEntity{
	Entity:     FollowedMe,
	App:        true,
	Email:      false,
	Discord:    false,
	IsPersonal: true,
	Events:     FollowedMeEvents,
}

var FollowedMyStudioEntity = NotificationEntity{
	Entity:     FollowedMyStudio,
	App:        true,
	Email:      false,
	Discord:    false,
	IsPersonal: true,
	Events:     FollowedMyStudioEvents,
}

var PublishAndMergeRequestsEntity = NotificationEntity{
	Entity:     PublishAndMergeRequests,
	App:        true,
	Email:      true,
	Discord:    true,
	IsPersonal: false,
	Events:     PublishAndMergeRequestsEvents,
}

var ResponseToMyRequestsEntity = NotificationEntity{
	Entity:     ResponseToMyRequests,
	App:        true,
	Email:      true,
	Discord:    true,
	IsPersonal: true,
	Events:     ResponseToMyRequestsEvents,
}

var SystemNotificationEntity = NotificationEntity{
	Entity:     SystemNotifications,
	App:        true,
	Email:      true,
	Discord:    true,
	IsPersonal: false,
	Events:     systemNotificationEvents,
}

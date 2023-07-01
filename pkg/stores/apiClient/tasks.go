package apiClient

// Reels tasks
const (
	AddReelToAlgolia      = "tasks:AddReelToAlgolia"
	DeleteReelFromAlgolia = "tasks:DeleteReelFromAlgolia"
)

const (
	SendToIntegration     = "tasks:SendToIntegration"
	SendPostToIntegration = "tasks:SendPostToIntegration"
)

const (
	LoginEmail = "tasks:HandleLoginEmailTask"
)

const (
	DeleteMergeRequestNotifications   = "tasks:DeleteMergeRequestNotifications"
	DeletePublishRequestNotifications = "tasks:DeletePublishRequestNotifications"
	DeleteModsOnCanvas                = "tasks:DeleteModsOnCanvas"
)

const (
	NotionImportHandler = "tasks:NotionImportHandler"
)

const (
	SlackIntegrationTask    = "tasks:SlackIntegrationTask"
	SlackBipMarkAction      = "tasks:SlackBipMarkAction"
	SlackEventSubscriptions = "tasks:SlackEventSubscriptions"
	SlackSlashCommands      = "tasks:SlackSlashCommands"
)

const (
	TranslateCanvasRepositories = "tasks:TranslateCanvasRepositories"
)

const UpdateDiscordTreeMessage = "tasks:UpdateDiscordTreeMessage"
const DiscordIntegrationTask = "discord-tasks:DiscordIntegrationTask"

// cron tasks
const (
	TestTaskCronMethod         = "tasks:TestTaskCronMethod"
	TestTaskCronMethod1        = "tasks:TestTaskCronMethod1"
	CanvasBranchAccessCron     = "tasks:CanvasBranchAccessCron"
	RunFailedDiscordEventsCron = "tasks:RunFailedDiscordEventsCron"
)

const (
	RoughBranchNotificationsOnMerge = "tasks:RoughBranchNotificationsOnMerge"
)

package eventHandler

const (
	// redis event names
	RedisDiscordNamespace          = "discordevents:"
	RedisDiscordNamespaceProcessed = "discordevents-success:"
	RedisDiscordNamespaceFailed    = "discordevents-failed:"
	RedisDiscordNamespaceAll       = "discordevents:*"

	MessageCreate     = "messageCreate"
	GuildMemberAdd    = "guildMemberAdd"
	GuildMemberRemove = "guildMemberRemove"
	GuildRoleCreate   = "guildRoleCreate"
	GuildRoleDelete   = "guildRoleDelete"
	GuildRoleUpdate   = "guildRoleUpdate"
	GuildMemberUpdate = "guildMemberUpdate"
)

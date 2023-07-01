package eventHandler

type Discord struct {
	Bots []DiscordBot `json:"bots,omitempty"`
}

type DiscordBot struct {
	BotName string           `json:"bot_name,omitempty"`
	BotID   string           `json:"bot_id,omitempty"`
	Config  DiscordBotConfig `json:"config,omitempty"`
	Servers []DiscordServer  `json:"servers,omitempty"`
}

type DiscordBotConfig struct {
	Token  string        `json:"token,omitempty"`
	Game   string        `json:"game,omitempty"`
	DMResp ResponseArray `json:"dm_response,omitempty"`
}

type DiscordServer struct {
	ServerID    string              `json:"server_id,omitempty"`
	Config      DiscordServerConfig `json:"config,omitempty"`
	ChanGroups  []ChannelGroup      `json:"channel_groups,omitempty"`
	Permissions []Permission        `json:"permissions,omitempty"`
	Filters     []Filter            `json:"filters,omitempty"`
}

type DiscordServerConfig struct {
	Prefix   string          `json:"prefix,omitempty"`
	Clear    bool            `json:"clear_commands,omitempty"`
	WebHooks discordWebHooks `json:"web_hooks,omitempty"`
}

type discordWebHooks struct {
	Logs string `json:"logs,omitempty"`
}

type discordKickOnMention struct {
	Roles   []string      `json:"roles,omitempty"`
	Users   []string      `json:"users,omitempty"`
	Direct  ResponseArray `json:"dm,omitempty"`
	Channel ResponseArray `json:"channel,omitempty"`
	Kick    bool          `json:"kick,omitempty"`
}

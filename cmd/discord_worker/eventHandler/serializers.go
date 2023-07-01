package eventHandler

import "github.com/bwmarrin/discordgo"

type EventSerializer struct {
	Event interface{} `json:"event"`
	Type  string      `json:"type"`
}

type MessageCreateSerializer struct {
	Event *discordgo.MessageCreate `json:"event"`
	Type  string                   `json:"type"`
}

type GuildMemberAddSerializer struct {
	Event *discordgo.GuildMemberAdd `json:"event"`
	Type  string                    `json:"type"`
}

type GuildMemberRemoveSerializer struct {
	Event *discordgo.GuildMemberRemove `json:"event"`
	Type  string                       `json:"type"`
}

type GuildMemberUpdateSerializer struct {
	Event *discordgo.GuildMemberUpdate `json:"event"`
	Type  string                       `json:"type"`
}

type GuildRoleCreateSerializer struct {
	Event *discordgo.GuildRoleCreate `json:"event"`
	Type  string                     `json:"type"`
}

type GuildRoleUpdateSerializer struct {
	Event *discordgo.GuildRoleUpdate `json:"event"`
	Type  string                     `json:"type"`
}

type GuildRoleDeleteSerializer struct {
	Event *discordgo.GuildRoleDelete `json:"event"`
	Type  string                     `json:"type"`
}

package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// InteractionType indicates the type of an interaction event.
type InteractionType uint8

// Interaction types
const (
	InteractionPing                           InteractionType = 1
	InteractionApplicationCommand             InteractionType = 2
	InteractionMessageComponent               InteractionType = 3
	InteractionApplicationCommandAutocomplete InteractionType = 4
)

func (t InteractionType) String() string {
	switch t {
	case InteractionPing:
		return "Ping"
	case InteractionApplicationCommand:
		return "ApplicationCommand"
	case InteractionMessageComponent:
		return "MessageComponent"
	}
	return fmt.Sprintf("InteractionType(%d)", t)
}

type Option struct {
	Name    string `json:"name"`
	Type    int    `json:"type"`
	Value   string `json:"value"`
	Focused bool   `json:"focused"`
}

type ComponentType uint

type InteractionData struct {
	CustomID      string        `json:"custom_id"`
	ComponentType ComponentType `json:"component_type"`
	Values        []string      `json:"values"`
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Type          int           `json:"type"`
	Options       []Option      `json:"options"`
	Resolved      struct {
		Messages map[string]discordgo.Message `json:"messages"`
	} `json:"resolved"`
}

type Interaction struct {
	ID        string          `json:"id"`
	Type      InteractionType `json:"type"`
	Data      InteractionData `json:"data"`
	GuildID   string          `json:"guild_id"`
	ChannelID string          `json:"channel_id"`

	// The message on which interaction was used.
	// NOTE: this field is only filled when a button click triggered the interaction. Otherwise it will be nil.
	Message *discordgo.Message `json:"message"`

	// The member who invoked this interaction.
	// NOTE: this field is only filled when the slash command was invoked in a guild;
	// if it was invoked in a DM, the `User` field will be filled instead.
	// Make sure to check for `nil` before using this field.
	Member *discordgo.Member `json:"member"`
	// The user who invoked this interaction.
	// NOTE: this field is only filled when the slash command was invoked in a DM;
	// if it was invoked in a guild, the `Member` field will be filled instead.
	// Make sure to check for `nil` before using this field.
	User *discordgo.User `json:"user"`

	Token   string `json:"token"`
	Version int    `json:"version"`
}

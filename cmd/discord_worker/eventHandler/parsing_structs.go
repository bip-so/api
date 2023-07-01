package eventHandler

// generic structs
type Permission struct {
	Group       string   `json:"group,omitempty"`
	Users       []string `json:"users,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	Commands    []string `json:"commands,omitempty"`
	Blacklisted bool     `json:"blacklisted,omitempty"`
}

type Command struct {
	Command  string   `json:"command,omitempty"`
	Response []string `json:"response,omitempty"`
	Reaction []string `json:"reaction,omitempty"`
}

type Keyword struct {
	Keyword  string   `json:"keyword,omitempty"`
	Reaction []string `json:"reaction,omitempty"`
	Response []string `json:"response,omitempty"`
	Exact    bool     `json:"exact,omitempty"`
}

type Mentions struct {
	Ping    ResponseArray `json:"ping,omitempty"`
	Mention ResponseArray `json:"mention,omitempty"`
}

type Filter struct {
	Term   string   `json:"term,omitempty"`
	Reason []string `json:"reason,omitempty"`
}

type ResponseArray struct {
	Reaction []string `json:"reaction,omitempty"`
	Response []string `json:"response,omitempty"`
}

type Parsing struct {
	Image ParsingImageConfig `json:"image,omitempty"`
	Paste ParsingPasteConfig `json:"paste,omitempty"`
}

type ParsingConfig struct {
	Name   string `json:"name,omitempty"`
	URL    string `json:"url,omitempty"`
	Format string `json:"format,omitempty"`
}

type ParsingImageConfig struct {
	FileTypes []string        `json:"filetypes,omitempty"`
	Sites     []ParsingConfig `json:"sites,omitempty"`
}

type ParsingPasteConfig struct {
	Sites  []ParsingConfig `json:"sites,omitempty"`
	Ignore []ParsingConfig `json:"ignore,omitmepty"`
}

type ChannelGroup struct {
	ChannelIDs  []string             `json:"channels,omitempty"`
	Mentions    Mentions             `json:"mentions,omitempty"`
	Commands    []Command            `json:"commands,omitempty"`
	Keywords    []Keyword            `json:"keywords,omitempty"`
	Parsing     Parsing              `json:"parsing,omitempty"`
	Permissions []Permission         `json:"permissions,omitempty"`
	KOM         discordKickOnMention `json:"kick_on_mention,omitempty"`
}

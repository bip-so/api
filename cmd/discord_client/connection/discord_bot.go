package connection

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
)

var dg *discordgo.Session
var bot *string

var (
	StopDiscordServer    = make(chan string)
	DiscordServerStopped = make(chan string)
	discordLoad          = make(chan string)
	commands             = []*discordgo.ApplicationCommand{
		{
			Name:        "bip-new",
			Description: "Create a new canvas",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "title",
					Description: "Title of the new canvas to be created",
					Required:    true,
				},
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "collection",
					Description:  "Title of the collection under which new canvas is to be created",
					Required:     true,
					Autocomplete: true,
				},
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "parent-canvas",
					Description:  "Title of the parent canvas under which new canvas is to be created",
					Required:     false,
					Autocomplete: true,
				},
			},
		},
		{
			Name:        "bip-search",
			Description: "Search a canvas",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "title",
					Description:  "Title of the canvas",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
	}
)

// StartDiscordBotConnection: Inits a new session with the Discord Bot
// We are passing the Token and the Name to this function
// https://github.com/bwmarrin/discordgo
// More examples : https://github.com/bwmarrin/discordgo/tree/master/examples

func StartDiscordBotConnection(token string, botName string) {
	fmt.Printf("⇾ Starting connections for %s ", botName)
	// Initializing Discord connection
	// Create a new Discord session using the provided bot token.
	fmt.Printf("using token '%s' to auth", token)
	var err error
	dg, err = discordgo.New("Bot " + token)
	if err != nil {
		fmt.Printf("x Error creating Discord session for %s: %v", botName, err)
		return
	}

	bot = &botName
	dg.Identify.Intents = dg.Identify.Intents | discordgo.IntentsGuildMembers
	fmt.Println("Current sessions intent is ", dg.Identify.Intents)
	// Logging
	messagePrefix := "\n → Discord Event : "

	dg.AddHandler(func(dg *discordgo.Session, event *discordgo.MessageCreate) {
		if event != nil {
			fmt.Println(messagePrefix, *event.Message)
			if !event.Author.Bot {
				setDataInRedis(event, event.Message.ID, MessageCreate)
			}
		}
	})

	dg.AddHandler(func(dg *discordgo.Session, event *discordgo.GuildMemberAdd) {
		if event != nil {
			fmt.Println(messagePrefix, *event.Member)
			setDataInRedis(event, event.Member.User.ID, GuildMemberAdd)
		}

	})
	dg.AddHandler(func(dg *discordgo.Session, event *discordgo.GuildMemberRemove) {
		if event != nil {
			fmt.Println(messagePrefix, *event.Member)
			setDataInRedis(event, event.Member.User.ID, GuildMemberRemove)
		}

	})
	dg.AddHandler(func(dg *discordgo.Session, event *discordgo.GuildRoleCreate) {
		if event != nil {
			fmt.Println(messagePrefix, *event.GuildRole)
			setDataInRedis(event, event.GuildRole.Role.ID, GuildRoleCreate)
		}

	})
	dg.AddHandler(func(dg *discordgo.Session, event *discordgo.GuildRoleDelete) {
		if event != nil {
			fmt.Println(messagePrefix, *event)
			setDataInRedis(event, event.RoleID, GuildRoleDelete)
		}

	})
	dg.AddHandler(func(dg *discordgo.Session, event *discordgo.GuildRoleUpdate) {
		if event != nil {
			fmt.Println(messagePrefix, *event.GuildRole)
			setDataInRedis(event, event.GuildRole.Role.ID, GuildRoleUpdate)
		}

	})
	dg.AddHandler(func(dg *discordgo.Session, event *discordgo.GuildMemberUpdate) {
		if event != nil {
			fmt.Println(messagePrefix, *event.Member)
			setDataInRedis(event, event.Member.User.ID, GuildMemberUpdate)
		}

	})
	fmt.Printf("⇾  Discord service connected for %s", botName)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Print("x Error opening connection,", err)
		return
	}

	for _, v := range commands {
		res, err := dg.ApplicationCommandCreate(configs.GetConfigString("DISCORD_APP_ID"), "", v)
		if err != nil {
			fmt.Println("x Cannot create " + v.Name + " command: " + err.Error())
		}
		fmt.Println(res)
	}

	discordLoad <- ""

	<-StopDiscordServer

	fmt.Printf("⇾  Stop recieved on %s", botName)
	// properly send a shutdown to the discord server so the bot goes offline.
	if err := dg.Close(); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s sent close", botName)
	// return the shutdown signal
	DiscordServerStopped <- fmt.Sprintf("⇾  Closed connection for %s", botName)
}

// DiscordClient return the discord go client instance.
func DiscordClient() *discordgo.Session {
	return dg
}

func BotName() *string {
	return bot
}

package integrations

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
)

var (
	discord              *discordgo.Session
	StopDiscordServer    = make(chan string)
	DiscordServerStopped = make(chan string)
	DiscordPostsChannel  = "bip-posts"
)

func InitDiscordGo() error {
	var err error
	fmt.Println(configs.GetConfigString("DISCORD_TOKEN"))
	discord, err = discordgo.New("Bot " + configs.GetConfigString("DISCORD_TOKEN"))
	if err != nil {
		return err
	}

	return nil
}

func CreateChannel(guildID string) (*discordgo.Channel, error) {
	ch, err := discord.GuildChannelCreate(guildID, "bip stage 1", discordgo.ChannelTypeGuildText)
	if err != nil {
		fmt.Println("Error ", err)
		return nil, err
	}
	fmt.Println(ch.Name, ch.Type, ch.GuildID, ch.ID)
	return ch, err
}

func CreateBipBotCanvasChannel(guildID string, channelName string) (*discordgo.Channel, error) {
	channels, err := discord.GuildChannels(guildID)
	if err != nil {
		return nil, err
	}
	var channel *discordgo.Channel = nil
	for _, ch := range channels {
		if ch.Name == channelName {
			channel = ch
		}
	}
	if channel != nil {
		_, err = discord.ChannelDelete(channel.ID)
	}
	ch, err := discord.GuildChannelCreate(guildID, channelName, discordgo.ChannelTypeGuildText)
	if err != nil {
		fmt.Println("Error ", err)
		return nil, err
	}
	channel = ch
	return channel, nil
}

func CreateBipPostsCanvasChannel(guildID string) (*discordgo.Channel, error) {
	channels, err := discord.GuildChannels(guildID)
	if err != nil {
		return nil, err
	}
	var channel *discordgo.Channel = nil
	for _, ch := range channels {
		if ch.Name == DiscordPostsChannel {
			channel = ch
			return channel, nil
		}
	}
	ch, err := discord.GuildChannelCreate(guildID, "bip-posts", discordgo.ChannelTypeGuildText)
	if err != nil {
		fmt.Println("Error in creating guild channel", err)
		return nil, err
	}
	channel = ch
	return channel, nil
}

func GetDiscordTeamMembers(discordTeamID string) ([]*discordgo.Member, error) {
	var result []*discordgo.Member
	for {
		after := ""
		if len(result) > 0 {
			after = result[len(result)-1].User.ID
		}
		members, err := discord.GuildMembers(discordTeamID, after, 1000)
		if len(members) == 0 {
			break
		}
		if err != nil {
			return nil, err
		}
		result = append(result, members...)
	}
	return result, nil
}

func GetDiscordTeamRoles(discordTeamID string) ([]*discordgo.Role, error) {
	roles, err := discord.GuildRoles(discordTeamID)
	if err != nil {
		return []*discordgo.Role{}, err
	}
	return roles, nil
}

func GetDiscordCategories(discordTeamID string) ([]*discordgo.Channel, error) {
	channels, err := discord.GuildChannels(discordTeamID)
	if err != nil {
		return []*discordgo.Channel{}, err
	}
	return channels, nil
}

func GetDiscordTeam(discordTeamID string) (*discordgo.Guild, error) {
	guild, err := discord.Guild(discordTeamID)
	if err != nil {
		return nil, err
	}
	return guild, nil
}

func ReplyToDiscordMessage(channelID string, reference *discordgo.MessageReference, responseArray []string) (err error) {
	// if there is no response to sen just return
	if len(responseArray) == 0 {
		return errors.New("no message")
	}

	response := strings.Join(responseArray, "\n")

	// if there is an error return the error
	if _, err = discord.ChannelMessageSendReply(channelID, response, reference); err != nil {
		return err
	}

	return nil
}

type messageSend struct {
	Content    string        `json:"content,omitempty"`
	Components []interface{} `json:"components"`
	TTS        bool          `json:"tts"`
}

type MessageBtnComponent struct {
	Type  int    `json:"type"`
	Label string `json:"label"`
	Style int    `json:"style"`
	Url   string `json:"url"`
}

type ActionRowsComponent struct {
	Type       int           `json:"type"`
	Components []interface{} `json:"components"`
}

func SendDiscordUserDM(userID string, responseArray []string) (string, error) {

	response := strings.Join(responseArray, "\n")
	channel, err := discord.UserChannelCreate(userID)
	if err != nil {
		return "", err
	}

	message, err := discord.ChannelMessageSend(channel.ID, response)
	if err != nil {
		return "", err
	}
	return message.ID, nil
}

func EditDiscordUserDM(userID string, messageID string, responseArray []string) (string, error) {
	if messageID == "" {
		return "", errors.New("messageId is empty")
	}
	response := strings.Join(responseArray, "\n")
	channel, err := discord.UserChannelCreate(userID)
	if err != nil {
		return "", err
	}

	message, err := discord.ChannelMessageEdit(channel.ID, messageID, response)
	if err != nil {
		return "", err
	}
	return message.ID, nil
}

func SendDiscordUserComponentsDM(userID string, responseArray []string, components []interface{}) (discordgo.Message, error) {
	response := strings.Join(responseArray, "\n")
	channel, err := discord.UserChannelCreate(userID)
	if err != nil {
		return discordgo.Message{}, err
	}

	message := messageSend{
		Content:    response,
		Components: components,
	}

	endpoint := discordgo.EndpointChannelMessages(channel.ID)

	result, err := discord.RequestWithBucketID("POST", endpoint, message, endpoint)
	if err != nil {
		return discordgo.Message{}, err
	}
	m := discordgo.Message{}
	err = json.Unmarshal(result, &m)
	if err != nil {
		return discordgo.Message{}, err
	}
	return m, nil
}

func SendDiscordReaction(channelID string, messageID string, reactionArray []string) (err error) {

	if len(reactionArray) == 0 {
		return
	}

	for _, reaction := range reactionArray {
		if err = discord.MessageReactionAdd(channelID, messageID, reaction); err != nil {
			return
		}
	}
	return
}

func DeleteDiscordMessage(userID string, messageID string) (err error) {
	channel, err := discord.UserChannelCreate(userID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = discord.ChannelMessageDelete(channel.ID, messageID)
	if err != nil {
		fmt.Println("error in cloud delete msg====>", err)
		return
	}
	return
}

// send a dm message with an embed
func SendDiscordDMEmbed(userID string, embed *discordgo.MessageEmbed) error {
	// if there is an error sending the embed message
	channel, err := discord.UserChannelCreate(userID)
	if err != nil {
		return err
	}
	if _, err := discord.ChannelMessageSendEmbed(channel.ID, embed); err != nil {
		fmt.Printf("Embed send error")
		fmt.Println(err)
		return err
	}

	return nil
}

// bip-stage-1 0 958371136368959579 1027442589789601832
func SendDiscordEmbedToChannel(channelID string, embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	msg, err := discord.ChannelMessageSendEmbed(channelID, embed)
	return msg, err
}

func EditDiscordEmbedToChannel(channelID string, messageId string, embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	msg, err := discord.ChannelMessageEditEmbed(channelID, messageId, embed)
	return msg, err
}

func EditDiscordEmbedComplexToChannel(channelID string, messageId string, embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	text := ""
	msg, err := discord.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    channelID,
		ID:         messageId,
		Embed:      embed,
		Components: []discordgo.MessageComponent{},
		Content:    &text,
	})
	return msg, err
}

func DeleteDiscordEmbedToChannel(channelID string, messageId string) error {
	err := discord.ChannelMessageDelete(channelID, messageId)
	return err
}

func SendDiscordDMMessageToChannel(channelID string, responseArray []string, components []interface{}) (*discordgo.Message, error) {
	response := strings.Join(responseArray, "\n")
	message := messageSend{
		Content:    response,
		Components: components,
	}
	endpoint := discordgo.EndpointChannelMessages(channelID)
	result, err := discord.RequestWithBucketID("POST", endpoint, message, endpoint)
	if err != nil {
		return nil, err
	}
	m := &discordgo.Message{}
	err = json.Unmarshal(result, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func GetDiscordUserByID(discordUserID string) (*discordgo.User, error) {
	user, err := discord.User(discordUserID)
	return user, err
}

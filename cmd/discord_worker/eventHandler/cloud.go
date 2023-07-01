package eventHandler

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
)

func InitDiscordGo() error {
	var err error
	discord, err = discordgo.New("Bot " + configs.GetConfigString("DISCORD_TOKEN"))
	if err != nil {
		return err
	}

	return nil
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
func SendDiscordDMEmbed(dg *discordgo.Session, userID string, embed *discordgo.MessageEmbed) error {
	// if there is an error sending the embed message
	channel, err := discord.UserChannelCreate(userID)
	if err != nil {
		return err
	}
	if _, err := dg.ChannelMessageSendEmbed(channel.ID, embed); err != nil {
		fmt.Printf("Embed send error")
		fmt.Println(err)
		return err
	}

	return nil
}

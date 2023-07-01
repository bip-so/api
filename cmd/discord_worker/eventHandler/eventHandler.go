package eventHandler

import (
	"context"
	"errors"
	"fmt"
	user2 "gitlab.com/phonepost/bip-be-platform/internal/user"
	"gitlab.com/phonepost/bip-be-platform/lambda/connect_discord/connect"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/supabase"
	"regexp"
	"strings"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"mvdan.cc/xurls/v2"

	"github.com/bwmarrin/discordgo"
)

// This function will be called (due to AddHandler) every time a new
// message is created on any channel that the authenticated bot has access to.
func discordMessageHandler(dg *discordgo.Session, m *discordgo.MessageCreate, botName string, ctx context.Context) error {

	bot, err := dg.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
	}

	// Ignore all messages created by bots (stops the bot uprising)
	if m.Author.Bot {
		fmt.Printf("User is a bot and being ignored.")
		return nil
	}

	isBotMentioned := false
	for _, user := range m.Mentions {
		if user.Username == bot.Username {
			isBotMentioned = true
			break
		}
	}

	// if m.MessageReference == nil && isBotMentioned {

	// 	// parentMessage, err := dg.ChannelMessage(m.ChannelID, m.MessageReference.MessageID)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }

	// 	// captureMessageFlow(dg, m, botName, ctx, bot.Email, parentMessage)

	// } else
	if m.MessageReference != nil && isBotMentioned {

		parentMessage, err := dg.ChannelMessage(m.ChannelID, m.MessageReference.MessageID)
		if err != nil {
			return err
		}

		if !parentMessage.Author.Bot {
			err = captureMessageFlow(dg, m, botName, ctx, bot.Email, parentMessage)
			if err != nil {
				return err
			}
		} else {
			err = replyToReelFlow(dg, m, botName, ctx, bot, bot.Email)
			if err != nil {
				return err
			}
		}
	}
	return nil

}

func captureMessageFlow(dg *discordgo.Session, m *discordgo.MessageCreate, botName string, ctx context.Context, userEmail string, parentMessage *discordgo.Message) error {

	var err error
	users, _ := FindUsersByDiscordIDs([]string{m.Author.ID, parentMessage.Author.ID})
	var foundAuthor *models.UserSocialAuth
	var foundUser *models.UserSocialAuth
	if len(users) > 0 {
		for i, user := range users {
			if parentMessage.Author.ID == user.ProviderID {
				foundAuthor = &users[i]
			}
			if m.Author.ID == user.ProviderID {
				foundUser = &users[i]
			}
		}
	}
	if foundAuthor == nil {
		discordUserName := parentMessage.Author.Username + parentMessage.Author.Discriminator
		//nickName := parentMessage.Author.Username
		user := CreateNewUser("", "", discordUserName, m.Author.Username, parentMessage.Author.AvatarURL(""))
		discordUser := NewDiscordUser(user.ID, parentMessage.Author.ID, nil)
		err = CreateUserSocialAuth(discordUser)
		if err != nil {
			println("error while creating user", err.Error())
			return err
		}
		foundAuthor = discordUser
	}
	if foundUser == nil {
		discordUserName := m.Author.Username + m.Author.Discriminator
		//nickName := m.Author.Username
		// Discord User
		user := CreateNewUser("", "", discordUserName, m.Author.Username, m.Author.AvatarURL(""))
		// TODO: pass metadata
		discordUser := NewDiscordUser(user.ID, m.Author.ID, nil)
		err = CreateUserSocialAuth(discordUser)
		if err != nil {
			println("error while creating user", err.Error())
			return err
		}
		foundUser = discordUser
	}

	timestamp := parentMessage.Timestamp

	re := regexp.MustCompile(`\<\@.*\>`)
	messageText := re.ReplaceAllString(parentMessage.Content, "")
	attachments := []string{}
	for _, att := range parentMessage.Attachments {
		attachments = append(attachments, att.URL)
	}
	guild, _ := GetDiscordTeam(m.GuildID)
	message := models.NewDiscordMessage(parentMessage.ID, messageText, foundAuthor.UserID, foundUser.UserID, timestamp, attachments, guild.ID, guild.Name, guild.IconURL(), "")
	err = CreateMessage([]*models.Message{message})
	if err != nil {
		println("error while creating message ", err.Error())
	}
	// go SendDiscordUserDM(foundUser.DiscordID.String, response)
	embed := &discordgo.MessageEmbed{
		URL:   "https://discordapp.com/channels/" + guild.ID + "/" + parentMessage.ChannelID + "/" + parentMessage.ID,
		Title: "✅ bip Marked.Type // in canvas to select & embed bip Marks",
		Color: 0x44b244,
	}

	go SendDiscordDMEmbed(dg, foundUser.ProviderID, embed)

	// go ReplyToDiscordMessage(m.ChannelID, m.MessageReference, response)
	if err == nil {
		go SendDiscordReaction(m.ChannelID, m.ID, []string{"☑️"})
		go SendDiscordReaction(parentMessage.ChannelID, parentMessage.ID, []string{"✅"})
	}
	return nil
}

func replyToReelFlow(dg *discordgo.Session, m *discordgo.MessageCreate, botName string, ctx context.Context, bot *discordgo.User, userEmail string) error {
	// data to send to discord
	var response []string
	var reaction []string

	if m.MessageReference != nil {
		msgId := m.MessageReference.MessageID
		addDiscordMessage(ctx, map[string]interface{}{
			"id":            msgId,
			"user_id":       m.Author.ID,
			"email":         userEmail,
			"username":      m.Author.Username,
			"discriminator": m.Author.Discriminator,
			"text":          m.Content,
			"avatar":        m.Author.AvatarURL(""),
		})
	} else {
		return errors.New("failed")
	}

	// get channel information
	channel, err := dg.State.Channel(m.ChannelID)
	if err != nil {
		fmt.Println("Channel error", err)
		return err
	}

	// get all the configs
	// requires channel info we get from the channel info above
	prefix := getPrefix("discord", botName, channel.GuildID)
	channelCommands := getCommands("discord", botName, channel.GuildID, m.ChannelID)
	channelKeywords := getKeywords("discord", botName, channel.GuildID, m.ChannelID)
	channelParsing := getParsing("discord", botName, channel.GuildID, m.ChannelID)
	serverFilter := getFilter("discord", botName, channel.GuildID)

	fmt.Printf("\nprefix: %s", prefix)

	// bot level configs for log reading
	maxLogs, logResponse, logReaction, allowIP := getBotParseConfig()

	// if the channel is a DM
	if channel.Type == 1 {
		_, dmResp := getMentions("discord", botName, channel.GuildID, "DirectMessage")
		if err := sendDiscordMessage(dg, m.ChannelID, m.Author.ID, getPrefix("discord", botName, channel.GuildID), dmResp.Reaction); err != nil {
			fmt.Println(err)
			return err
		}

		if err := sendDiscordReaction(dg, m.ChannelID, m.ID, dmResp.Reaction); err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	}

	// fmt.Printf("all channels %s", getChannels("discord", botName, channel.GuildID))

	//filter logic
	fmt.Printf("filtering messages")
	if len(serverFilter) == 0 {
		fmt.Printf("no filtered terms found")
	} else {
		for _, filter := range serverFilter {
			if strings.Contains(m.Content, filter.Term) {
				fmt.Printf("message was removed for containing %s", filter.Term)
				if err := deleteDiscordMessage(dg, m.ChannelID, m.ID, ""); err != nil {
					fmt.Println(err)
					return err
				}

				if err := sendDiscordMessage(dg, m.ChannelID, m.Author.ID, prefix, filter.Reason); err != nil {
					fmt.Println(err)
					return err
				}
				return nil
			} else {
				continue
			}
		}
	}

	// if the channel isn't in a group drop the message
	fmt.Printf("checking channels")
	if !contains(getChannels("discord", botName, channel.GuildID), m.ChannelID) {
		fmt.Printf("channel not found")
		return errors.New("channel not found")
	}

	fmt.Printf("checking blacklist")

	// drop messages from blacklisted users
	for _, user := range getBlacklist("discord", botName, channel.GuildID, m.ChannelID) {
		if user == m.Author.ID {
			fmt.Printf("user %s is blacklisted username is %s", m.Author.ID, m.Author.Username)
			return errors.New("blacklisted user")
		}
	}

	fmt.Printf("checking attachments")

	// for all attachment urls
	var attachmentURLs []string
	for _, url := range m.Attachments {
		attachmentURLs = append(attachmentURLs, url.ProxyURL)
	}

	fmt.Printf("all attachments %s", attachmentURLs)
	fmt.Printf("all ignores %+v", channelParsing.Paste.Ignore)

	fmt.Printf("checking for any urls in the message")
	var allURLS []string
	for _, url := range xurls.Relaxed().FindAllString(m.Content, -1) {
		fmt.Printf("checking on %s", url)
		// if the url is an ip filter it out
		if match, err := regexp.Match("^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])", []byte(url)); err != nil {
			fmt.Println(err)
		} else if match && !allowIP {
			fmt.Printf("adding %s to the list", url)
			continue
		}

		fmt.Printf("looking for ignored domains")
		if len(channelParsing.Paste.Ignore) == 0 {
			fmt.Printf("appending %s to allURLS", url)
			allURLS = append(allURLS, url)
			fmt.Printf("no ignored domain found")
			continue
		} else {
			var ignored bool
			for _, ignoreURL := range channelParsing.Paste.Ignore {
				fmt.Printf("url should be ignored: %t", strings.HasPrefix(url, ignoreURL.URL))
				if strings.HasPrefix(url, ignoreURL.URL) {
					ignored = true
					fmt.Printf("domain %s is being ignored.", ignoreURL.URL)
					break
				}
			}
			if ignored {
			} else {
				fmt.Printf("appending %s to allURLS", url)
				allURLS = append(allURLS, url)
			}
		}
	}

	// add all urls together
	fmt.Printf("adding attachment URLS to allURLS")
	for i := 0; i < len(attachmentURLs); i++ {
		allURLS = append(allURLS, attachmentURLs[i])
	}

	// fmt.Printf(allURLS)
	fmt.Printf("checking mentions")
	if len(m.Mentions) != 0 {
		ping, mention := getMentions("discord", botName, channel.GuildID, m.ChannelID)
		if m.Mentions[0].ID == bot.ID && strings.Replace(m.Content, "<@!"+dg.State.User.ID+">", "", -1) == "" {
			fmt.Printf("bot was pinged")
			response = ping.Response
			reaction = ping.Reaction
		} else {
			for _, mentioned := range m.Mentions {
				if mentioned.ID == bot.ID {
					fmt.Printf("bot was mentioned")
					response = mention.Response
					reaction = mention.Reaction
				}
			}
		}
	} else {
		fmt.Printf("no mentions found")
		if strings.HasPrefix(m.Content, prefix) {
			// command
			response, reaction = parseCommand(strings.TrimPrefix(m.Content, prefix), botName, channelCommands)
			// if the flag for clearing commands is set and there is a response
			if getCommandClear("discord", botName, channel.GuildID) && len(response) > 0 {
				fmt.Printf("removing command message %s", m.ID)
				if err := deleteDiscordMessage(dg, m.ChannelID, m.ID, ""); err != nil {
					fmt.Println(err)
				}
			} else {

			}
		} else {
			// keyword
			response, reaction = parseKeyword(m.Content, botName, channelKeywords, channelParsing)
		}
	}

	if len(channelParsing.Image.FileTypes) == 0 && len(channelParsing.Paste.Sites) == 0 {
		fmt.Printf("no parsing configs found")
	} else {
		fmt.Printf("allURLS: %s", allURLS)
		fmt.Printf("allURLS count: %d", len(allURLS))

		// if we have too many logs ignore it.
		if len(allURLS) == 0 {
			fmt.Printf("no URLs to read")
		} else if len(allURLS) > maxLogs {
			fmt.Printf("too many logs or screenshots to try and read.")
			if err := sendDiscordMessage(dg, m.ChannelID, m.Author.ID, prefix, logResponse); err != nil {
				fmt.Println(err)
			}
			if err := sendDiscordReaction(dg, m.ChannelID, m.ID, logReaction); err != nil {
				fmt.Println(err)
			}
			return err
		} else {
			fmt.Printf("reading logs")
			if err := sendDiscordReaction(dg, m.ChannelID, m.ID, []string{"👀"}); err != nil {
				fmt.Println(err)
			}

			// get parsed content for each url/attachment
			fmt.Printf("reading all attachments and logs")
			allParsed := make(map[string]string)
			for _, url := range allURLS {
				allParsed[url] = parseURL(url, channelParsing)
			}

			//parse logs and append to current response.
			for _, url := range allURLS {
				fmt.Printf("passing %s to keyword parser", url)
				urlResponse, _ := parseKeyword(allParsed[url], botName, channelKeywords, channelParsing)
				fmt.Printf("response length = %d", len(urlResponse))
				if len(urlResponse) == 1 && urlResponse[0] == "" || len(urlResponse) == 0 {

				} else {
					response = append(response, fmt.Sprintf("I have found the following for: <%s>", url))
					for _, singleLine := range urlResponse {
						response = append(response, singleLine)
					}
				}
			}
		}
	}

	// send response to channel
	fmt.Printf("sending response %s to %s", response, m.ChannelID)
	if err := sendDiscordMessage(dg, m.ChannelID, m.Author.ID, prefix, response); err != nil {
		fmt.Println(err)
	}

	// send reaction to channel
	fmt.Printf("sending reaction %s", reaction)
	if err := sendDiscordReaction(dg, m.ChannelID, m.ID, reaction); err != nil {
		fmt.Println(err)
	}
	return nil
}

// clean up messages if configured to
func deleteDiscordMessage(dg *discordgo.Session, channelID, messageID, message string) (err error) {
	fmt.Printf("Removing message \n'%s'\n from %s", message, channelID)

	if err = dg.ChannelMessageDelete(channelID, messageID); err != nil {
		return err
	}
	return nil
}

// send message handling
func sendDiscordMessage(dg *discordgo.Session, channelID, authorID, prefix string, responseArray []string) (err error) {
	// if there is no response to sen just return
	if len(responseArray) == 0 {
		return nil
	}

	response := strings.Join(responseArray, "\n")
	response = strings.Replace(response, "&user&", "<@"+authorID+">", -1)
	response = strings.Replace(response, "&prefix&", prefix, -1)
	response = strings.Replace(response, "&react&", "", -1)

	// if there is an error return the error
	if _, err = dg.ChannelMessageSend(channelID, response); err != nil {
		return err
	}

	return nil
}

// send a reaction to a message
func sendDiscordReaction(dg *discordgo.Session, channelID string, messageID string, reactionArray []string) (err error) {
	// if there is no reaction to sen just return
	if len(reactionArray) == 0 {
		return
	}

	for _, reaction := range reactionArray {
		fmt.Printf("sending \"%s\" as a reaction to message: %s", reaction, messageID)
		// if there is an error sending a message return it
		if err = dg.MessageReactionAdd(channelID, messageID, reaction); err != nil {
			return
		}
	}
	return
}

// send a message with an embed
func sendDiscordEmbed(dg *discordgo.Session, channelID string, embed *discordgo.MessageEmbed) error {
	// if there is an error sending the embed message
	if _, err := dg.ChannelMessageSendEmbed(channelID, embed); err != nil {
		fmt.Printf("Embed send error")
		return err
	}

	return nil
}

// service handling
// start all the bots
func startDiscordsBots(token string, botName string) {
	fmt.Printf("Starting discord server connections\n")
	// range over the bots available to start
	for _, bot := range discordGlobal.Bots {
		fmt.Printf("Connecting to %s\n", bot.BotName)

		// spin up a channel to tell the bot to shutdown later
		// stopDiscord[bot.BotName] = make(chan string)

		// start the bot
		//go StartDiscordBotConnection(token, botName)
		// wait on bot being able to start.
		<-discordLoad
	}

	fmt.Printf("Discord service started\n")
	servStart <- "discord_online"
}

// when a shutdown is sent close out services properly
func stopDiscordBots() {
	fmt.Printf("stopping discord connections")
	// loop through bots and send shutdowns
	for _, bot1 := range discordGlobal.Bots {
		StopDiscordServer <- bot1.BotName
		//stopDiscordServer <- bot1.BotName
	}

	for range discordGlobal.Bots {
		botIn := <-DiscordServerStopped
		//botIn := <-discordServerStopped
		fmt.Printf("%s", botIn)
	}

	fmt.Printf("discord connections stopped")
	// return shutdown signal on channel
	servStopped <- "discord_stopped"
}

func addDiscordMessage(ctx context.Context, event map[string]interface{}) {
	id, isOk := event["id"].(string)
	if isOk {
		ext, err := GetExternalReference(id, DISCORD_INTEGRATION_TYPE)
		if err != nil {
			println("extternal reference not found...")
			return
		}
		fmt.Println(ext)

		activityId := ext.InternalID
		userID := event["user_id"].(string)
		userName := event["username"].(string)
		avatar := event["avatar"].(string)
		commentText := event["text"].(string)
		email := event["email"].(string)
		discriminator := event["discriminator"].(string)

		re := regexp.MustCompile(`\<\@.*\>`)
		commentText = re.ReplaceAllString(commentText, "")
		var newUser *models.User

		var activity *models.Reel

		extUser, err := GetExternalReference(userID, DISCORD_USER_TYPE)
		if err != nil || extUser.ID == "" {
			println("discord user not found.....")
			activity, err = GetReel(activityId)
			if err != nil {
				println("error while getting activity", err.Error())
				return
			}
			integration, err := GetDiscordStudioIntegration(activity.StudioID)
			if err != nil {
				println("error while getting integration ", err.Error())
			}
			if len(integration) < 1 {
				return
			}

			uniqueUsername := userName + discriminator
			existingUsers, err := FindUsersByDiscordIDs([]string{userID})
			if err != nil {
				return
			}
			println("userLEn", len(existingUsers), err != nil)
			if err != nil || len(existingUsers) < 1 {

				newUser = CreateNewUser(email, "", uniqueUsername, userName, avatar)

				if err != nil {
					return
				}

				if err != nil {
					println("error while creating new user: error: ", err.Error())
				}
			} else {
				println("\n\n\n===>")
				newUser = existingUsers[0].User
				println("===>>>")
			}

		} else {
			newUser, err = GetUser(map[string]interface{}{"id": extUser.ExternalID})
		}
		println("created new user", newUser.ID)

		/* TODO
		comment, _, err := newUser.NewComment(ctx, newUser.UserID,
			models.ReelFeedContentTypeName, activityId, commentText, 0, 0, "", false, []string{}, "")
		if err != nil {
			println("error while createing new comment", err.Error())
			return
		}
		newUser.CreateComment(ctx, comment, activity.PageID)
		*/
	}
}

func discordMemberAddRemoveHandler(dg *discordgo.Session, m *discordgo.Member, botName string, ctx context.Context, isAdd bool) error {
	fmt.Printf("bot is %s", botName)

	integrations, err := GetStudioIntegrationByDiscordTeamId(m.GuildID)
	if err != nil {
		println("error while getting integration ", err.Error())
		return err
	}
	if len(integrations) < 1 {
		return errors.New("no integration found")
	}

	integration := integrations[0]

	studio, err := FindStudio(integration.StudioID)
	if err != nil {
		fmt.Println("product not found")
		return err
	}

	existingUsers, err := FindUsersByDiscordIDs([]string{m.User.ID})
	if err != nil {
		println("error while getting discord user", err.Error())
		return err
	}
	discordRolesList, _ := GetDiscordTeamRoles(integration.TeamID)
	if err != nil {
		println("error while getting discord roles", err.Error())
		return err
	}

	allRoleIDs := []string{}
	for _, role := range discordRolesList {
		allRoleIDs = append(allRoleIDs, role.ID)
	}
	existingMemberGroups, _ := GetRolesByDiscordRoleIDs(allRoleIDs, studio.ID)

	var user *models.UserSocialAuth
	if len(existingUsers) == 1 {
		user = &existingUsers[0]
		InvalidateUserPermissionRedisCache(user.UserID)
	} else {
		discordUserName := m.User.Username + m.User.Discriminator
		nickName := m.User.Username
		if m.Nick != "" {
			nickName = m.Nick
		}
		newuser := connect.CreateNewUser("", "", discordUserName, nickName, m.User.AvatarURL(""))
		// TODO: pass metadata
		user = NewDiscordUser(newuser.ID, m.User.ID, nil)
		err = CreateUserSocialAuth(user)
		if err != nil {
			println("error while creating discord user", err.Error())
			return err
		}
		user2.App.Service.AddUserToAlgolia(newuser.ID)
	}
	if isAdd {
		UserJoinStudio(studio.ID, user.UserID)
	} else {
		LeaveStudio(studio.ID, user.UserID)
	}
	for _, role := range discordRolesList {
		if role.Managed || (role.Name == "@everyone" && isAdd) {
			continue
		}
		var foundMemberGroup *models.Role = nil
		for _, mGroup := range *existingMemberGroups {
			if mGroup.DiscordRoleID.String == role.ID {
				foundMemberGroup = &mGroup
				break
			}
		}
		if foundMemberGroup == nil {
			continue
		}
		if !isAdd {
			add := []uint64{}
			remove := []uint64{user.UserID}
			err = UpdateMembershipRole(foundMemberGroup.ID, add, remove)
			if err != nil {
				return err
			}
			continue
		}
		if !ExistsIn(role.ID, m.Roles) {
			continue
		}
		if isAdd && foundMemberGroup != nil {
			add := []uint64{user.UserID}
			remove := []uint64{}
			err = UpdateMembershipRole(foundMemberGroup.ID, add, remove)
			if err != nil {
				return err
			}
		}
	}
	DeleteUserAssociatedStudioDataByUserID(user.UserID)
	supabase.UpdateUserSupabase(user.UserID, true)
	fmt.Println("members add remove successfully")
	return nil

}

func discordNewRoleCreateHandler(dg *discordgo.Session, m *discordgo.GuildRoleCreate, botName string, ctx context.Context) error {
	fmt.Printf("bot is %s", botName)
	fmt.Printf("bot is %s", botName)

	integrations, err := GetStudioIntegrationByDiscordTeamId(m.GuildID)
	if err != nil {
		println("error while getting integration ", err.Error())
		return err
	}
	/* TESTING
	var integration *models.StudioIntegration
	if len(integrations) < 1 {
		//TODO remove (temp)
		integration, err = CreateStudioIntegration(m.GuildID)
		if err != nil {
			return
		}
	}
	*/

	if len(integrations) < 1 {
		return errors.New("no integrations found")
	}
	integration := integrations[0]
	studio, err := FindStudio(integration.StudioID)
	if err != nil {
		fmt.Println("product not found")
		return err
	}

	discordMemberList, err := GetDiscordTeamMembers(integration.TeamID)
	if err != nil {
		return err
	}
	discordIDs := []string{}
	for _, member := range discordMemberList {
		discordIDs = append(discordIDs, member.User.ID)
	}

	existingUsers, err := FindUsersByDiscordIDs(discordIDs)
	if err != nil {
		return err
	}
	userIDs := []uint64{}
	for _, member := range discordMemberList {
		if member.User.Bot {
			continue
		}
		for _, eUser := range existingUsers {
			if eUser.ProviderID == member.User.ID {
				userIDs = []uint64{eUser.UserID}
				break
			}
		}
	}

	role, _ := CreateNewRole(studio.ID, m.Role.Name, userIDs, m.Role.ID)
	fmt.Println("new role created  ", role)
	return err

}

func discordRoleDeleteHandler(dg *discordgo.Session, m *discordgo.GuildRoleDelete, botName string, ctx context.Context) error {
	fmt.Printf("bot is %s", botName)

	integrations, err := GetStudioIntegrationByDiscordTeamId(m.GuildID)
	if err != nil {
		println("error while getting integration ", err.Error())
		return err
	}
	if len(integrations) < 1 {
		return errors.New("no integrations found")
	}

	integration := integrations[0]
	studio, err := FindStudio(integration.StudioID)
	if err != nil {
		fmt.Println("product not found")
		return err
	}

	role, err := GetRoleByDiscordRoleID(m.RoleID, studio.ID)
	if err != nil {
		fmt.Println("role not found")
		return err
	}

	err = DeleteRole(role.ID)
	if err != nil {
		fmt.Println("err in deleting role")
		return err
	}
	fmt.Println("role deleted successfully")
	return nil
}

func discordRoleUpdateHandler(dg *discordgo.Session, m *discordgo.GuildRoleUpdate, botName string, ctx context.Context) error {
	fmt.Printf("bot is %s", botName)

	integrations, err := GetStudioIntegrationByDiscordTeamId(m.GuildID)
	if err != nil {
		println("error while getting integration ", err.Error())
		return err
	}
	/* TESTING
	var integration *models.StudioIntegration
	if len(integrations) < 1 {
		integration, err = CreateStudioIntegration(m.GuildID)
		if err != nil {
			return
		}
	}
	*/
	if len(integrations) < 1 {
		return err
	}

	integration := integrations[0]
	studio, err := FindStudio(integration.StudioID)
	if err != nil {
		fmt.Println("product not found")
		return err
	}

	role, err := GetRoleByDiscordRoleID(m.Role.ID, studio.ID)
	if err != nil {
		fmt.Println("group not found")
		return err
	}

	err = UpdateRole(role.ID, map[string]interface{}{
		"name": m.Role.Name,
	})
	if err != nil {
		return err
	}
	fmt.Println("role updated successfully")
	return nil

}

func discordMemberUpdateHandler(dg *discordgo.Session, m *discordgo.GuildMemberUpdate, botName string, ctx context.Context) error {
	fmt.Printf("bot is %s", botName)
	// updating user info
	UpdateDiscordUserOnEvent(m.User)

	integrations, err := GetStudioIntegrationByDiscordTeamId(m.GuildID)
	if err != nil {
		println("error while getting integration ", err.Error())
		return err
	}
	if len(integrations) < 1 {
		return errors.New("no integration found")
	}

	integration := integrations[0]
	studio, err := FindStudio(integration.StudioID)
	if err != nil {
		fmt.Println("product not found")
		return err
	}

	discordRolesList, _ := GetDiscordTeamRoles(integration.TeamID)
	if err != nil {
		println("error while getting discord roles", err.Error())
		return err
	}

	existingUsers, err := FindUsersByDiscordIDs([]string{m.User.ID})
	if err != nil {
		return err
	}
	if len(existingUsers) == 1 {
		user := existingUsers[0]

		roles, err := GetRolesByStudioID(studio.ID)

		if err != nil {
			println("error while getting discord roles", err.Error())
			return err
		}
		for _, role := range discordRolesList {
			if role.Managed || role.Name == "@everyone" {
				continue
			}
			var foundGroup *models.Role
			for _, mGroup := range roles {
				if role.Name == mGroup.Name {
					foundGroup = &mGroup
					break
				}
			}
			if foundGroup == nil {
				fmt.Println("did not foundGroup", role.Name)
				continue
			}

			if ExistsIn(role.ID, m.Roles) {
				err = UpdateMembershipRole(foundGroup.ID, []uint64{user.UserID}, []uint64{})
			} else if !ExistsIn(role.ID, m.Roles) {
				err = UpdateMembershipRole(foundGroup.ID, []uint64{}, []uint64{user.UserID})
			}
			if err != nil {
				return err
			}
		}
		InvalidateUserPermissionRedisCache(user.UserID)
	}
	fmt.Println("members updated successfully")
	return nil
}

func UpdateDiscordUserOnEvent(m *discordgo.User) {
	if m.ID != "" {
		// updating user avatar url
		existingUsers, _ := FindUsersByDiscordIDs([]string{m.ID})
		if len(existingUsers) == 1 {
			user := existingUsers[0]
			fmt.Println("user Avatar from discord", m.Avatar)
			err := connect.UpdateUserByID(user.UserID, map[string]interface{}{"avatar_url": m.AvatarURL("")})
			if err != nil {
				fmt.Println("error in updating user", err)
			}
			supabase.UpdateUserSupabase(user.UserID, true)
		}
	}
}

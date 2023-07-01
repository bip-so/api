package eventHandler

import (
	"fmt"
)

var (
	discordGlobal Discord
	discordLoad   = make(chan string)
	servStart     = make(chan string)
	servStopped   = make(chan string)
)

func getBlacklist(inService, botName, inServer, inChannel string) (blacklist []string) {
	perms := []Permission{}
	switch inService {
	case "discord":
		for _, bot := range discordGlobal.Bots {
			if bot.BotName == botName {
				for _, server := range bot.Servers {
					if inServer == server.ServerID {
						perms = server.Permissions
					}
				}
			}
		}
	case "irc":
		for _, group := range getChannelGroups(inService, botName, inServer, inChannel) {
			for _, channel := range group.ChannelIDs {
				if channel == inChannel {
					perms = group.Permissions
				}
			}
		}
	default:
	}

	// load users that are in blacklisted groups
	for _, perm := range perms {
		if perm.Blacklisted {
			for _, user := range perm.Users {
				blacklist = append(blacklist, user)
			}
		}
	}

	return
}

func getChannels(inService, botName, inServer string) (channels []string) {
	fmt.Printf("service: %s, bot: %s, server: %s", inService, botName, inServer)
	switch inService {
	case "discord":
		for bid := range discordGlobal.Bots {
			fmt.Printf("checking for bot: %s", discordGlobal.Bots[bid].BotName)
			if botName == discordGlobal.Bots[bid].BotName {
				fmt.Printf("matched for %s", discordGlobal.Bots[bid].BotName)
				for sid := range discordGlobal.Bots[bid].Servers {
					fmt.Printf("checking for server: %s", discordGlobal.Bots[bid].Servers[sid].ServerID)
					if inServer == discordGlobal.Bots[bid].Servers[sid].ServerID {
						fmt.Printf("matched for %s", discordGlobal.Bots[bid].Servers[sid].ServerID)
						for gid := range discordGlobal.Bots[bid].Servers[sid].ChanGroups {
							fmt.Printf("%s", discordGlobal.Bots[bid].Servers[sid].ChanGroups[gid].ChannelIDs)
							for _, channel := range discordGlobal.Bots[bid].Servers[sid].ChanGroups[gid].ChannelIDs {
								channels = append(channels, channel)
							}
						}
					}
				}
			}
		}

	default:
	}

	fmt.Printf("handing channels back with a value of %s", channels)

	return
}

func getChannelGroups(inService, botName, inServer, inChannel string) (chanGroups []ChannelGroup) {
	switch inService {
	case "discord":
		for _, bot := range discordGlobal.Bots {
			if bot.BotName == botName {
				for _, server := range bot.Servers {
					if inServer == server.ServerID {
						chanGroups = server.ChanGroups
					}
				}
			}
		}
	default:
	}

	return
}

func getCommands(inService, botName, inServer, inChannel string) (commands []Command) {
	// prep stuff for passing to the parser
	for _, group := range getChannelGroups(inService, botName, inServer, inChannel) {
		for _, channel := range group.ChannelIDs {
			if inChannel == channel {
				for _, command := range group.Commands {
					commands = append(commands, command)
				}
			}
		}
	}

	return
}

func getKeywords(inService, botName, inServer, inChannel string) (keywords []Keyword) {
	// prep stuff for passing to the parser
	for _, group := range getChannelGroups(inService, botName, inServer, inChannel) {
		for _, channel := range group.ChannelIDs {
			if inChannel == channel {
				for _, keyword := range group.Keywords {
					keywords = append(keywords, keyword)
				}
			}
		}
	}

	return
}

func getMentions(inService, botName, inServer, inChannel string) (ping, mention ResponseArray) {
	switch inService {
	case "discord":
		for _, bot := range discordGlobal.Bots {
			if bot.BotName == botName {
				for _, server := range bot.Servers {
					if inServer == server.ServerID {
						if inChannel == "DirectMessage" {
							mention = bot.Config.DMResp
						} else {
							for _, group := range server.ChanGroups {
								for _, channel := range group.ChannelIDs {
									if inChannel == channel {
										fmt.Printf("bot was mentioned on channel %s", channel)
										fmt.Printf("ping resp %s", group.Mentions.Ping)
										fmt.Printf("mention resp %s", group.Mentions.Mention)
										ping = group.Mentions.Ping
										mention = group.Mentions.Mention
										return
									}
								}
							}
						}
					}
				}
			}
		}
	default:
	}

	return
}

func getParsing(inService, botName, inServer, inChannel string) (parseConf Parsing) {
	// prep stuff for passing to the parser
	for _, group := range getChannelGroups(inService, botName, inServer, inChannel) {
		for _, channel := range group.ChannelIDs {
			if inChannel == channel {
				parseConf = group.Parsing
			}
		}
	}

	return
}

func getFilter(inService, botName, inServer string) (filters []Filter) {
	// prep stuff for passing to the parser
	switch inService {
	case "discord":
		for _, bot := range discordGlobal.Bots {
			if bot.BotName == botName {
				for _, server := range bot.Servers {
					if inServer == server.ServerID {
						filters = server.Filters
					}
				}
			}
		}
	case "irc":
	default:
	}

	return
}

func getBotParseConfig() (maxLogs int, response, reaction []string, allowIP bool) {
	return
	//return botConfig.Parsing.Max, botConfig.Parsing.Response, botConfig.Parsing.Reaction, botConfig.Parsing.AllowIP
}

func getPrefix(inService, botName, inServer string) (prefix string) {
	switch inService {
	case "discord":
		for _, bot := range discordGlobal.Bots {
			if bot.BotName == botName {
				for _, server := range bot.Servers {
					if inServer == server.ServerID {
						prefix = server.Config.Prefix
					}
				}
			}
		}
	default:
	}

	return
}

func getCommandClear(inService, botName, inServer string) (clear bool) {
	switch inService {
	case "discord":
		for _, bot := range discordGlobal.Bots {
			if bot.BotName == botName {
				for _, server := range bot.Servers {
					if inServer == server.ServerID {
						clear = server.Config.Clear
					}
				}
			}
		}
	default:
	}

	return
}

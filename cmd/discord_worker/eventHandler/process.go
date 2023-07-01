package eventHandler

import (
	"context"
	"encoding/json"
	"fmt"
)

func Process(data EventSerializer, valStr string) error {

	ctx := context.Background()
	botName := "bip"
	dg := discord

	//bot := &botName
	//dg := bot.DiscordClient()
	//botName := bot.BotName()

	switch data.Type {

	case MessageCreate:
		fmt.Println("processing message create")
		result := MessageCreateSerializer{}
		err := json.Unmarshal([]byte(valStr), &result)
		if err != nil {
			return err
		}

		// @todo commenting it out here because this is failing in `cmd/discord_worker/eventHandler/eventHandler.go`
		// at line `bot, err := dg.User("@me")` To execute all the other failed events commenting this for now.
		err = discordMessageHandler(dg, result.Event, botName, ctx)

		if err != nil {
			return err
		}
		fmt.Println("processing message create done")

	case GuildMemberAdd:
		fmt.Println("processing guild member add")
		result := GuildMemberAddSerializer{}
		err := json.Unmarshal([]byte(valStr), &result)
		if err != nil {
			return err
		}
		err = discordMemberAddRemoveHandler(dg, result.Event.Member, botName, ctx, true)
		if err != nil {
			fmt.Println("Error on discord member add", err)
			return err
		}
		fmt.Println("processing guild member add done")

	case GuildMemberRemove:
		fmt.Println("processing guild member remove")
		result := GuildMemberRemoveSerializer{}
		err := json.Unmarshal([]byte(valStr), &result)
		if err != nil {
			return err
		}
		err = discordMemberAddRemoveHandler(dg, result.Event.Member, botName, ctx, false)
		if err != nil {
			return err
		}
		fmt.Println("processing guild member remove done")

	case GuildMemberUpdate:
		fmt.Println("processing guild member update")
		result := GuildMemberUpdateSerializer{}
		err := json.Unmarshal([]byte(valStr), &result)
		if err != nil {
			return err
		}
		err = discordMemberUpdateHandler(dg, result.Event, botName, ctx)
		if err != nil {
			return err
		}
		fmt.Println("processing guild member done")

	case GuildRoleCreate:
		fmt.Println("processing guild role create")
		result := GuildRoleCreateSerializer{}
		err := json.Unmarshal([]byte(valStr), &result)
		if err != nil {
			return err
		}
		err = discordNewRoleCreateHandler(dg, result.Event, botName, ctx)
		if err != nil {
			return err
		}
		fmt.Println("processing guild role create done")

	case GuildRoleDelete:
		fmt.Println("processing message delete")
		result := GuildRoleDeleteSerializer{}
		err := json.Unmarshal([]byte(valStr), &result)
		if err != nil {
			return err
		}
		err = discordRoleDeleteHandler(dg, result.Event, botName, ctx)
		if err != nil {
			return err
		}
		fmt.Println("processing guild role delete done")

	case GuildRoleUpdate:
		fmt.Println("processing guild role update")
		result := GuildRoleUpdateSerializer{}
		err := json.Unmarshal([]byte(valStr), &result)
		if err != nil {
			return err
		}
		err = discordRoleUpdateHandler(dg, result.Event, botName, ctx)
		if err != nil {
			return err
		}
		fmt.Println("processing guild role update done")
	}
	return nil
}

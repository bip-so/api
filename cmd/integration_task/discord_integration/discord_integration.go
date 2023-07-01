package discord_integration

import (
	"errors"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/lambda/connect_discord/connect"
	"gitlab.com/phonepost/bip-be-platform/pkg/integrations"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/supabase"
)

func DiscordIntegrationTask(integration *models.StudioIntegration) error {
	studio := integration.Studio
	fmt.Println("teamid: ", integration.TeamID)
	discordMemberList, err := integrations.GetDiscordTeamMembers(integration.TeamID)
	if err != nil {
		fmt.Println("Error in getting discord members", err)
		return err
	}
	discordRolesList, err := integrations.GetDiscordTeamRoles(integration.TeamID)
	if err != nil {
		fmt.Println("Error in getting discord member roles", err)
		return err
	}
	fmt.Println("discordMemberList & roles", discordMemberList, discordRolesList)
	discordIDs := []string{}
	for _, member := range discordMemberList {
		discordIDs = append(discordIDs, member.User.ID)
	}
	fmt.Println("discordIDs", discordIDs)

	existingUsers, err := connect.FindUsersByDiscordIDs(discordIDs)
	if err != nil {
		fmt.Println("Error in find users by discordIDs", err)
		return err
	}
	fmt.Println("existingUsers", existingUsers)

	members, err := connect.GetMembers(studio.ID)
	if err != nil {
		fmt.Println("Error in get studio members", err)
		return err
	}
	fmt.Println("followers", members)

	roleUserIDMap := map[string][]uint64{}
	usersSocialAuth := []models.UserSocialAuth{}
	users := []models.User{}
	memberStudioUserIds := []uint64{}
	for _, member := range discordMemberList {
		if member.User.Bot {
			continue
		}
		var foundUser *models.UserSocialAuth = nil
		var foundMember *models.Member = nil
		// TODO: can be optimised by using map
		for _, eUser := range existingUsers {
			if eUser.ProviderID == member.User.ID {
				foundUser = &eUser
				break
			}
		}
		for _, member := range members {
			if foundUser != nil && foundUser.UserID == member.UserID {
				foundMember = &member
				break
			}
		}
		if foundUser == nil {
			discordUserName := member.User.Username + member.User.Discriminator + connect.GenerateRandomKey(4)
			nickName := member.User.Username
			if member.Nick != "" {
				nickName = member.Nick
			}
			user := connect.CreateNewUser("", "", discordUserName, nickName, member.User.AvatarURL(""))
			discordUser := connect.NewDiscordUser(user.ID, member.User.ID, nil)
			users = append(users, *user)
			usersSocialAuth = append(usersSocialAuth, *discordUser)
			foundUser = discordUser
		}
		if foundMember == nil {
			memberStudioUserIds = append(memberStudioUserIds, foundUser.UserID)
		}
		for _, role := range member.Roles {
			if roleUserIDMap[role] == nil {
				roleUserIDMap[role] = []uint64{foundUser.UserID}
			} else {
				roleUserIDMap[role] = append(roleUserIDMap[role], foundUser.UserID)
			}
			roleUserIDMap["@everyone"] = append(roleUserIDMap["@everyone"], foundUser.UserID)
		}
	}
	fmt.Println("users", users)
	fmt.Println("roleUserIDMap", roleUserIDMap)
	if len(users) > 0 {
		err = connect.CreateUsers(&users, &usersSocialAuth)
		if err != nil {
			fmt.Println("Error in create users", err)
			return err
		}
	}
	fmt.Println("memberStudioUserIds", memberStudioUserIds)
	connect.BulkJoinStudio(studio.ID, memberStudioUserIds)
	allRoleNames := []string{"Discord Members"}
	for _, role := range discordRolesList {
		allRoleNames = append(allRoleNames, role.Name)
	}
	existingMemberGroups, _ := connect.GetRolesByNames(allRoleNames, studio.ID)
	for _, role := range discordRolesList {
		if role.Managed {
			continue
		}
		userIDs := roleUserIDMap[role.ID]
		name := role.Name
		if name == "@everyone" {
			name = "Discord Members"
			userIDs = roleUserIDMap["@everyone"]
		}
		var foundMemberGroup *models.Role = nil
		for _, mGroup := range existingMemberGroups {
			if mGroup.DiscordRoleID.String == role.ID {
				foundMemberGroup = &mGroup
				break
			}
		}
		if foundMemberGroup == nil {
			group, err := connect.CreateNewRole(studio.ID, name, role.ID)
			if err != nil {
				fmt.Println(err, name, role.ID)
				continue
			}
			fmt.Println("new role cerated======>", group)
			foundMemberGroup = group
		}
		connect.UpdateMembershipRole(foundMemberGroup.ID, userIDs, []uint64{})
	}

	// Send studio file structure to the discord channel
	err = SendCanvasTree(studio, integration.TeamID, integration.ID)
	if err != nil {
		fmt.Println("Error in send Canvas Tree", err)
		return err
	}
	postsChannelID := connect.CreateBipPostsChannel(integration.TeamID, integration.ID)
	if postsChannelID == "" {
		return errors.New("error in creating bip-posts channel ID")
	}
	DeleteUserAssociatedStudioDataByUserID(integration.CreatedByID)
	supabase.UpdateUserSupabase(integration.CreatedByID, true)
	notifications.App.Service.PublishNewNotification(notifications.DiscordIntegrationTask, 0, []uint64{}, &studio.ID,
		nil, notifications.NotificationExtraData{}, nil, nil)
	return nil
}

func SendCanvasTree(studio *models.Studio, guildID string, integrationID uint64) error {
	err := connect.CreateStudioFileTree(studio.ID, guildID, integrationID)
	return err
}

func DeleteUserAssociatedStudioDataByUserID(userID uint64) (*models.UserAssociatedStudio, error) {
	var userStudios models.UserAssociatedStudio
	err := postgres.GetDB().Where(models.UserAssociatedStudio{UserID: userID}).Delete(&userStudios).Error
	return &userStudios, err
}

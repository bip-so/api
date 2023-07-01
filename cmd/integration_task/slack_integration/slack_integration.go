package slack_integration

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/lambda/slack/connect"
	"gitlab.com/phonepost/bip-be-platform/pkg/integrations"
)

func SlackIntegrationTask(integration *models.StudioIntegration) error {
	accessToken := integration.AccessKey
	teamID := integration.TeamID
	studioID := integration.StudioID

	slackMemberList, err := integrations.GetSlackTeamMembers(accessToken, teamID)
	if err != nil {
		fmt.Println("Error in getting slack members", slackMemberList)
	}
	fmt.Println("slackMemberList", slackMemberList)

	slackIDs := []string{}
	for _, slackMember := range slackMemberList.Members {
		slackIDs = append(slackIDs, slackMember.Id)
	}
	fmt.Println("slackIDs", slackIDs)

	existingUsers, err := connect.FindUsersBySlackIDs(slackIDs)
	if err != nil {
		fmt.Println("error in getting user by slackIDs", err)
	}
	fmt.Println("Existing user", existingUsers)

	members, err := connect.GetMembers(studioID)
	if err != nil {
		fmt.Println("Error in getting studio members", slackMemberList)
	}
	fmt.Println("Studio Members", members)

	usersSocialAuth := []models.UserSocialAuth{}
	users := []models.User{}
	memberStudioUserIds := []uint64{}
	for _, member := range slackMemberList.Members {
		if member.IsBot || member.Deleted || member.Id == "USLACKBOT" {
			continue
		}
		var foundUser *models.UserSocialAuth
		var foundMember *models.Member
		for _, eUser := range existingUsers {
			if eUser.ProviderID == member.Id {
				foundUser = &eUser
				break
			}
		}

		for _, studioMember := range members {
			if foundUser != nil && foundUser.UserID == studioMember.UserID {
				foundMember = &studioMember
				break
			}
		}

		if foundUser == nil {
			slackUserName := member.Profile.DisplayName + connect.GenerateRandomKey(4)
			nickName := member.Profile.DisplayName
			user := connect.CreateNewUser("", "", slackUserName, nickName, member.Profile.Avatar)
			slackUser := connect.NewSlackUser(user.ID, member.Id, nil)
			users = append(users, *user)
			usersSocialAuth = append(usersSocialAuth, *slackUser)
			foundUser = slackUser
		} else {
			queries.App.StudioQueries.DeleteUserAssociatedStudioDataByUserID(foundUser.UserID)
		}
		if foundMember == nil {
			memberStudioUserIds = append(memberStudioUserIds, foundUser.UserID)
		}
	}
	if len(users) > 0 {
		err = connect.CreateUsers(&users, &usersSocialAuth)
		if err != nil {
			fmt.Println("Error in creating social auths")
			return err
		}
	}
	connect.BulkJoinStudio(studioID, memberStudioUserIds)

	slackUserGroupResponse, err := integrations.GetSlackUserGroups(accessToken, teamID)
	if err != nil {
		fmt.Println("Error in getting slack user groups")
		return err
	}
	allRoleNames := []string{}
	for _, userGroup := range slackUserGroupResponse.UserGroups {
		allRoleNames = append(allRoleNames, userGroup.Name)
	}
	existingMemberGroups, _ := connect.GetRolesByNames(allRoleNames, studioID)
	for _, userGroup := range slackUserGroupResponse.UserGroups {
		var foundMemberGroup *models.Role
		for _, mGroup := range existingMemberGroups {
			if mGroup.SlackRoleID.String == userGroup.Id {
				foundMemberGroup = &mGroup
				break
			}
		}
		if foundMemberGroup == nil {
			group, err := connect.CreateNewRole(studioID, userGroup.Name, userGroup.Id)
			if err != nil {
				fmt.Println(err, group.Name, group.ID)
				continue
			}
			fmt.Println("new role created======>", group)
			foundMemberGroup = group
		}
		userIDs := []uint64{}
		userSocialAuths, _ := connect.FindUsersBySlackIDs(userGroup.Users)
		for _, usa := range userSocialAuths {
			userIDs = append(userIDs, usa.UserID)
		}
		connect.UpdateMembershipRole(foundMemberGroup.ID, userIDs, []uint64{})
	}
	return nil
}

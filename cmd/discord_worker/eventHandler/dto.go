package eventHandler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"log"
	"regexp"
	"strconv"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	bipStream "gitlab.com/phonepost/bip-be-platform/pkg/stores/stream"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"gorm.io/datatypes"
	"gorm.io/gorm/clause"
)

const (
	DISCORD_INTEGRATION_TYPE = "discord"
	DISCORD_USER_TYPE        = "discord_user"
	DISCORD_PROVIDER         = "discord"
)

func FindStudio(studioId uint64) (studio *models.Studio, err error) {
	err = postgres.GetDB().Model(&models.Studio{}).Where("id = ?", studioId).Find(&studio).Error

	return
}

func FindUsersByDiscordIDs(userIds []string) (users []models.UserSocialAuth, err error) {
	err = postgres.GetDB().Model(&models.UserSocialAuth{}).Where("provider_id IN ?", userIds).Preload("User").Find(&users).Error

	return
}

func GetRolesByDiscordRoleIDs(discordRoleIDs []string, studioID uint64) (*[]models.Role, error) {
	var roles []models.Role
	err := postgres.GetDB().Model(&models.Role{}).Where("discord_role_id IN ? AND studio_id = ?", discordRoleIDs, studioID).Find(&roles).Error
	return &roles, err
}

func GetStudioIntegrationByDiscordTeamId(teamID string) (integration []models.StudioIntegration, err error) {
	condition := models.StudioIntegration{TeamID: teamID}

	condition.Type = DISCORD_INTEGRATION_TYPE
	err = postgres.GetDB().Model(&models.StudioIntegration{}).Where(condition).Find(&integration).Error

	return
}

func GetDiscordStudioIntegration(studioId uint64) (integration []models.StudioIntegration, err error) {
	condition := models.StudioIntegration{StudioID: studioId}

	condition.Type = DISCORD_INTEGRATION_TYPE
	err = postgres.GetDB().Model(&models.StudioIntegration{}).Where(condition).Find(&integration).Error
	return
}

func GetRoleByDiscordRoleID(discordRoleID string, studioID uint64) (*models.Role, error) {
	var role *models.Role
	err := postgres.GetDB().Model(&models.Role{}).Where("discord_role_id = ? AND studio_id = ?", discordRoleID, studioID).Find(&role).Error
	return role, err
}

func CreateUserSocialAuth(user *models.UserSocialAuth) error {
	err := postgres.GetDB().Create(user).Error
	return err
}

func NewDiscordUser(userId uint64, providerId string, metadata datatypes.JSON) *models.UserSocialAuth {
	return &models.UserSocialAuth{
		UserID:       userId,
		ProviderName: DISCORD_PROVIDER,
		ProviderID:   providerId,
		Metadata:     metadata,
	}
}

func CreateUser(user *models.User) error {
	err := postgres.GetDB().Create(user).Error
	return err
}
func CreateNewUser(email, password, username, fullName, avatarURL string) *models.User {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	processedUsername := reg.ReplaceAllString(username, "")
	user := &models.User{
		Email: sql.NullString{
			String: email,
			Valid:  true,
		},
		Password:        password,
		Username:        processedUsername,
		FullName:        fullName,
		HasPassword:     false,
		AvatarUrl:       avatarURL,
		IsEmailVerified: true,
		IsSetupDone:     true,
		//timezone
	}

	// TODO: runs one query for every user. Can be optimised?
	_ = CreateUser(user)
	_ = postgres.GetDB().Create(models.UserProfile{UserID: user.ID}).Error
	//user2.App.Service.AddUserToAlgolia(user.ID)
	return user
}

func CreateNewRole(studioId uint64, name string, userIds []uint64, discordRoleId string) (*models.Role, error) {

	//TODO: store members based on userids
	roleObject := &models.Role{
		StudioID: studioId,
		Name:     name,
		Color:    "#ffffff",
		IsSystem: false,
		Icon:     "",
	}
	roleObject.DiscordRoleID = sql.NullString{String: discordRoleId, Valid: true}

	createdRole, err := CreateRole(roleObject)
	if err != nil {
		return nil, err
	}

	return createdRole, nil
}

func CreateRole(role *models.Role) (*models.Role, error) {
	result := postgres.GetDB().Create(role)
	if result.Error != nil {
		return nil, result.Error
	}
	return role, nil
}

func UpdateRole(ID uint64, updates map[string]interface{}) error {
	err := postgres.GetDB().Model(&models.Role{}).Where("id = ?", ID).Updates(updates).Error
	return err
}

func UpdateMembershipRole(roleId uint64, addUserIds []uint64, removeUserIds []uint64) error {
	fmt.Println("Started Update Membership role Task")
	role, err := GetRole(roleId)

	if err != nil {
		return err
	}

	studioMembers, err := GetMembersByStudioID(role.StudioID, 0)

	if err != nil {
		return err
	}

	members, err := FindMembersInRole(role)
	if err != nil {
		return err
	}

	addMembers := []models.Member{}
	for _, mID := range addUserIds {
		found := false
		for _, eMember := range members {
			if mID == eMember.UserID {
				found = true
				break
			}
		}
		if !found {
			// given ID not found in role
			// check if the member exists, if not create a new member
			// given ID can be userId, or memberId
			exists := false
			var memberID uint64
			for _, memb := range studioMembers {
				if mID == memb.UserID {
					exists = true
					memberID = memb.ID
					break
				}
			}
			if !exists {
				//create new member
				memberID = AddMemberToStudio(mID, role.StudioID)
				if memberID == 0 {
					return errors.New("couldn't add user to studio")
				}
			}
			addMembers = append(addMembers, models.Member{BaseModel: models.BaseModel{ID: memberID}})
		}
	}
	fmt.Println("Finished processing the members")
	removeMembers := []models.Member{}
	for _, mID := range removeUserIds {
		found := false
		var memberID uint64
		for _, eMember := range members {
			if mID == eMember.UserID {
				found = true
				memberID = eMember.ID
				break
			}
		}
		if found {
			removeMembers = append(removeMembers, models.Member{BaseModel: models.BaseModel{ID: memberID}})
		}
	}

	if len(addMembers) > 0 {
		err = AddMembersInRole(addMembers, role)

		if err != nil {
			return err
		}
	}

	if len(removeMembers) > 0 {
		err = RemoveMembersInRole(removeMembers, role)
	}
	return err
}

func DeleteRole(roleId uint64) error {
	roleObject := models.Role{BaseModel: models.BaseModel{ID: roleId}}
	//err := postgres.GetDB().Delete(roleObject).Error
	err := postgres.GetDB().Unscoped().Select(clause.Associations).Delete(roleObject).Error
	return err
}

func GetRolesByStudioID(studioID uint64) ([]models.Role, error) {
	var roles []models.Role
	result := postgres.GetDB().Model(&roles).Where("studio_id = ?", studioID).Preload("Members").Preload("Studio").Preload("Members.User").Order("id asc").Find(&roles)

	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

func FindMembersInRole(role *models.Role) (members []models.Member, err error) {
	err = postgres.GetDB().Model(&role).Association("Members").Find(&members)
	return members, err
}

func AddMembersInRole(addMembers []models.Member, role *models.Role) error {
	//err := postgres.GetDB().Model(&role).Association("Members").Append(addMembers)
	query := `insert into role_members (member_id, role_id) values`
	for i, member := range addMembers {
		query = fmt.Sprintf("%s (%d,%d)", query, member.ID, role.ID)
		if i < (len(addMembers) - 1) {
			query += ","
		}
	}
	println("Query====>  ", query)
	err := postgres.GetDB().Exec(query).Error
	if err != nil {
		fmt.Println("Error in creating role_members: ", err)
	}
	return err
}

func RemoveMembersInRole(removeMembers []models.Member, role *models.Role) error {
	err := postgres.GetDB().Model(&role).Association("Members").Delete(removeMembers)
	return err
}

func GetRole(roleId uint64) (*models.Role, error) {
	var role *models.Role
	postgres.GetDB().Model(&models.Role{}).Where("id = ?", roleId).Preload("Members").Preload("Studio").Preload("Members.User").First(&role)
	return role, nil
}

func NewMember(userId uint64, studioId uint64) *models.Member {
	return &models.Member{
		UserID:      userId,
		StudioID:    studioId,
		CreatedByID: userId,
		UpdatedByID: userId,
	}
}

func AddMemberToStudio(userdId uint64, studioId uint64) uint64 {
	memberObject := NewMember(userdId, studioId)

	newCreatedMemberId := CreateMember(memberObject)
	if newCreatedMemberId == 0 {
		return 0
	}
	return newCreatedMemberId
}

func CreateMember(member *models.Member) uint64 {
	//result := postgres.GetDB().Create(member)
	//return member.ID, result.Error
	if err := postgres.GetDB().Create(&member).Error; err != nil {
		return 0
	}

	return member.ID
}

func CreateMessage(messages []*models.Message) error {
	return postgres.GetDB().Create(messages).Error
}

func GetMessagesByProductID(userID string, skip int) (*[]models.Message, error) {
	var messages []models.Message
	err := postgres.GetDB().Model(&models.Message{}).Where("user_id = ? and is_used = ?", userID, false).Preload("Author").Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return &messages, nil
}

// func GetMessageByID(messageID string) (*models.Message, error) {
// 	var message models.Message
// 	err := postgres.GetDB().Model(&models.Message{}).Where(models.Message{ID: messageID}).Preload("Author").First(&message).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &message, nil
// }

func GetMessageByRefID(refID string) (*models.Message, error) {
	var message models.Message
	err := postgres.GetDB().Model(&models.Message{}).Where(models.Message{RefID: refID}).Preload("Author").First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func MarkMessageAsUsed(messageId string, isUsed bool) error {
	err := postgres.GetDB().Model(&models.Message{}).Where("id = ? and is_used = ?", messageId, false).Update("is_used", isUsed).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteMessageById(ctx context.Context, userID string, messageId string) error {
	err := postgres.GetDB().Model(&models.Message{}).Where("user_id = ? and id = ?", userID, messageId).Delete(&models.Message{}).Error
	if err != nil {
		return err
	}
	return nil
}

func GetMembersByStudioID(studioID uint64, skip int) ([]models.Member, error) {
	var members []models.Member
	result := postgres.GetDB().Model(&models.Member{}).Where("studio_id = ?", studioID).Order("id asc").Find(&members)

	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

func GetUser(query map[string]interface{}) (*models.User, error) {
	var user models.User
	err := postgres.GetDB().Model(&models.User{}).Where(query).Preload("UserProfile").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetExternalReference(externalID string, externalSource string) (
	ext models.ExternalReference, err error) {
	err = postgres.GetDB().Where(models.ExternalReference{ExternalID: externalID,
		ExternalSourceType: externalSource}).First(&ext).Error
	return
}

func UserJoinStudio(studioID uint64, userID uint64) {
	// check user is present as member before creating a member.
	var member *models.Member
	err := postgres.GetDB().Model(models.Member{}).Where("user_id = ? and studio_id = ?", userID, studioID).First(&member).Error
	if member == nil || err != nil {
		member = models.NewMember(userID, studioID)
		CreateMember(member)
		AddMembersToStudioMemberRole(studioID, []models.Member{*member})
		role, err := GetRoleByName("Discord Members", studioID)
		if err != nil {
			return
		}
		err = AddMembersInRole([]models.Member{*member}, role)
		if err != nil {
			return
		}
	} else if member.HasLeft {
		postgres.GetDB().Model(models.Member{}).Where("user_id = ? and studio_id = ?", userID, studioID).Updates(map[string]interface{}{"has_left": false}).First(&member)
		AddMembersToStudioMemberRole(studioID, []models.Member{*member})
		role, err := GetRoleByName("Discord Members", studioID)
		if err != nil {
			return
		}
		err = AddMembersInRole([]models.Member{*member}, role)
		if err != nil {
			return
		}
	}
}

func LeaveStudio(studioId uint64, userId uint64) error {
	updates := map[string]interface{}{
		"has_left": true,
	}
	err := postgres.GetDB().Model(&models.Member{}).Where("user_id = ? and studio_id = ?", userId, studioId).Updates(updates).Error
	if err != nil {
		return nil
	}
	members, err := GetMembersByUserIDs([]uint64{userId}, studioId)
	if err != nil {
		return nil
	}
	err = RemoveMembersToStudioMemberRole(studioId, members)
	if err != nil {
		return nil
	}
	studioRoles, err := GetRolesByStudioID(studioId)
	if err != nil {
		return nil
	}
	for _, studioRole := range studioRoles {
		err = UpdateMembershipRole(studioRole.ID, []uint64{}, []uint64{userId})
		if err != nil {
			continue
		}
	}

	// @todo Later move to kafka on member join we need to invalidate the user associated studios and send event to supabase
	go func() {
		DeleteUserAssociatedStudioDataByUserID(userId)
		StreamLeaveStudio(studioId, userId)
	}()

	return err
}

func StreamLeaveStudio(studioID uint64, userID uint64) {
	s := bipStream.StreamClient()
	followee, err := s.FlatFeed(bipStream.FlatFeedName, utils.String(studioID))
	if err != nil {
		fmt.Println(err)
	}
	followerTimeline, err := s.FlatFeed(bipStream.FlatTimelineName, utils.String(userID))
	if err != nil {
		fmt.Println(err)
	}
	response, err := followerTimeline.Unfollow(context.Background(), followee)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Unfollow studio response:", response)
}

func DeleteUserAssociatedStudioDataByUserID(userID uint64) (*models.UserAssociatedStudio, error) {
	var userStudios models.UserAssociatedStudio
	err := postgres.GetDB().Where(models.UserAssociatedStudio{UserID: userID}).Delete(&userStudios).Error
	return &userStudios, err
}

func RemoveMembersToStudioMemberRole(studioID uint64, members []models.Member) error {
	studioPermission, err := getStudioPermission(
		map[string]interface{}{"studio_id": studioID, "permission_group": "pg_studio_none"})
	if err != nil {
		return err
	}

	if studioPermission.RoleId != nil {
		role, err := GetRole(*studioPermission.RoleId)
		if err != nil {
			return err
		}
		err = RemoveMembersInRole(members, role)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetMembersByUserIDs(userIds []uint64, studioID uint64) ([]models.Member, error) {
	var members []models.Member
	result := postgres.GetDB().Model(&models.Member{}).Where("user_id IN ? and studio_id = ?", userIds, studioID).Find(&members)

	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

func GetRoleByName(roleName string, studioID uint64) (*models.Role, error) {
	var roles models.Role
	postgres.GetDB().Model(&models.Role{}).Where("studio_id = ? and name = ?", studioID, roleName).First(&roles)
	return &roles, nil
}

func AddMembersToStudioMemberRole(studioID uint64, members []models.Member) error {
	studioPermission, err := getStudioPermission(
		map[string]interface{}{"studio_id": studioID, "permission_group": "pg_studio_none"})
	if err != nil {
		return err
	}
	if studioPermission.RoleId != nil {
		role, err := GetRole(*studioPermission.RoleId)
		if err != nil {
			return err
		}
		err = AddMembersInRole(members, role)
		if err != nil {
			return err
		}
	}
	return nil
}

func getStudioPermission(query map[string]interface{}) (models.StudioPermission, error) {
	var studioPerms models.StudioPermission
	err := postgres.GetDB().Model(&studioPerms).Where(query).Find(&studioPerms).Error
	return studioPerms, err
}

func deleteCahceStudioFollowerCount(studioid uint64) {
	studioIDStr := strconv.FormatUint(studioid, 10)
	rc := redis.RedisClient()
	rc.Del(redis.GetBgContext(), models.RedisFollowUserStudioNS+studioIDStr)
}

func GetReel(ID uint64) (*models.Reel, error) {
	var reel *models.Reel
	err := postgres.GetDB().Model(&models.Reel{}).Where("id = ?", ID).Find(&reel).Error
	if err != nil {
		return nil, err
	}
	return reel, nil
}

func InvalidateUserPermissionRedisCache(userID uint64) {
	rc := redis.RedisClient()
	rc.Del(redis.GetBgContext(), fmt.Sprintf("%s%d", permissions.PermissionsHash, userID))
}

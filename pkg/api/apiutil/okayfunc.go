package apiutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"gorm.io/datatypes"
	"strconv"
	"time"
)

func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

func TestMailer(g *gin.Context) {
	err := pkg.TestMailer()
	if err != nil {
		g.JSON(400, gin.H{
			"message": "Email not sent check console",
		})
	}

	g.JSON(200, gin.H{
		"message": "Email sent!",
	})

}

type UT struct {
	Id        uint64    `json:"id"`
	UUID      string    `json:"uuid"`
	FullName  string    `json:"fullName"`
	Username  string    `json:"username"`
	AvatarUrl string    `json:"avatarUrl"`
	Timestamp time.Time `json:"timestamp"`
	BranchID  uint64    `json:"branchID"`
	RepoID    uint64    `json:"repoID"`
}

func Tested(g *gin.Context) {
	var xx []UT
	var x datatypes.JSON
	x = []byte(`[{"id": 12, "uuid": "61fbadd4-9aef-465f-93d6-0c9a5bd3383d", "repoID": 0, "branchID": 201, "fullName": "Martin Love", "username": "Martin_", "avatarUrl": "https://dei6r756sazdx.cloudfront.net/user/61fbadd4-9aef-465f-93d6-0c9a5bd3383d/_8Fm-LLYVgAq.png", "timestamp": "2022-07-27T06:05:32.556058645Z"}]`)
	_ = json.Unmarshal(x, &xx)
	fmt.Println(x)
	fmt.Println("========================================")
	//for _, c := range xx {
	//	fmt.Printf("%+v\n", c)
	//}

}

func Okay(g *gin.Context) {

	scheme := "http"
	if g.Request.TLS != nil {
		scheme = "https"
	}
	g.JSON(200, gin.H{
		"swagger": scheme + "://" + g.Request.Host + g.Request.URL.Path + "swagger/index.html",
		"version": pkg.Version,
	})
}

type UserBlockContributor struct {
	Id        uint64    `json:"id"`
	UUID      string    `json:"uuid"`
	FullName  string    `json:"fullName"`
	Username  string    `json:"username"`
	AvatarUrl string    `json:"avatarUrl"`
	Timestamp time.Time `json:"timestamp"`
	BranchID  uint64    `json:"branchID"`
	RepoID    uint64    `json:"repoID"`
}

func BuildUniqueContribs(blk []UserBlockContributor) []UserBlockContributor {
	var unique []UserBlockContributor
uniqueLoop:
	for _, v := range blk {
		for i, u := range unique {
			if v.Id == u.Id {
				unique[i] = v
				continue uniqueLoop
			}
		}
		unique = append(unique, v)
	}
	return unique
}

func SearchandReplaceCanvasBranchPermission(redundantMID uint64, newMasterID uint64) {
	var perms *[]models.CanvasBranchPermission
	postgres.GetDB().Model(&models.CanvasBranchPermission{}).Where("member_id = ?", redundantMID).Find(&perms)
	for _, v := range *perms {
		updates := map[string]interface{}{
			"member_id": newMasterID,
		}
		_ = postgres.GetDB().Model(&models.CanvasBranchPermission{}).Where("id = ?", v.ID).Updates(&updates)
	}
}

func SearchandReplaceCollectionPermission(redundantMID uint64, newMasterID uint64) {
	var perms *[]models.CollectionPermission
	postgres.GetDB().Model(&models.CollectionPermission{}).Where("member_id = ?", redundantMID).Find(&perms)
	for _, v := range *perms {
		updates := map[string]interface{}{
			"member_id": newMasterID,
		}
		_ = postgres.GetDB().Model(&models.CollectionPermission{}).Where("id = ?", v.ID).Updates(&updates)
	}
}

func SearchandReplaceStudioPermission(redundantMID uint64, newMasterID uint64) {
	var perms *[]models.StudioPermission
	postgres.GetDB().Model(&models.StudioPermission{}).Where("member_id = ?", redundantMID).Find(&perms)
	for _, v := range *perms {
		updates := map[string]interface{}{
			"member_id": newMasterID,
		}
		_ = postgres.GetDB().Model(&models.StudioPermission{}).Where("id = ?", v.ID).Updates(&updates)
	}
}

func AddMembersInRole(addMembers []models.Member, role *models.Role) error {
	err := postgres.GetDB().Model(&role).Association("Members").Append(addMembers)
	return err
}

func RemoveMembersInRole(removeMembers []models.Member, role *models.Role) error {
	err := postgres.GetDB().Model(&role).Association("Members").Delete(removeMembers)
	return err
}

func SearchandReplaceInRole(redundantMID uint64, newMasterID uint64) {
	var memberRemoveInatance, memberAddInstance models.Member
	postgres.GetDB().Model(&models.Member{}).Where(map[string]interface{}{"id": redundantMID}).Preload("Roles").First(&memberRemoveInatance)
	postgres.GetDB().Model(&models.Member{}).Where(map[string]interface{}{"id": newMasterID}).Preload("Roles").First(&memberAddInstance)
	rolesOnRedunrdant := memberRemoveInatance.Roles
	//rolesOnNew := memberAddInstance.Roles

	_ = postgres.GetDB().Model(&memberAddInstance).Association("Roles").Append(rolesOnRedunrdant)
	_ = postgres.GetDB().Model(&memberRemoveInatance).Association("Roles").Clear()
}

type MembershipObject struct {
	UserID            uint64
	MasterMemberID    uint64
	RedundantMemberID []uint64
}

func MemberFix(g *gin.Context) {

	StudioID := 226
	fmt.Println("Doing Studio ID", StudioID)
	var members *[]models.Member
	var skipUserList []uint64
	postgres.GetDB().Table("members").Where("studio_id = ?", StudioID).Find(&members)
	for _, v := range *members {
		if utils.SliceContainsInt(skipUserList, v.UserID) {
			fmt.Println("skipping ", v.UserID)
			continue
		}
		skipUserList = append(skipUserList, v.UserID)
		var redundantmembers *[]models.Member
		var redundantmembersID []uint64

		masterMemberID := v.ID
		postgres.GetDB().Table("members").Where("studio_id = ? and user_id = ?", StudioID, v.UserID).Find(&redundantmembers)
		if len(*redundantmembers) > 1 {
			fmt.Println("Members found : ", len(*redundantmembers))
			for _, innerLoopMember := range *redundantmembers {
				if innerLoopMember.ID == v.ID {
					continue
				} else {
					redundantmembersID = append(redundantmembersID, innerLoopMember.ID)

				}
			}
			fmt.Println(" The Master memberID ", masterMemberID)
			fmt.Println(" The Redundant memberIDd ", redundantmembersID)
			fmt.Println("Fixing")

			for _, redundantId := range redundantmembersID {
				SearchandReplaceCanvasBranchPermission(redundantId, masterMemberID)
				SearchandReplaceCollectionPermission(redundantId, masterMemberID)
				SearchandReplaceStudioPermission(redundantId, masterMemberID)
				SearchandReplaceInRole(redundantId, masterMemberID)
				fmt.Println("Deleting the redundant member now", redundantId)
				postgres.GetDB().Delete(&models.Member{}, redundantId)
			}

		}
	}
}

// BlockFixer This one fixes the cotributor key
func BlockFixer(g *gin.Context) {
	//var p []UserBlockContributor
	fmt.Println("The fuck>>>>")
	var blocks *[]models.Block
	//	layout := "2006-01-02T15:04:05Z07:00"
	//updated_at > ?
	//result := postgres.GetDB().Table("blocks").Where("updated_at > ? ", "2022-09-10 00:04:12.738103 +00:00").Find(&blocks)
	//blockIDs := []uint64{351330, 351317, 351333, 351318, 351316, 351315, 351313, 351421, 351334, 351332, 351331, 351329, 351328, 351327, 351326, 351325, 351324, 351323, 351322, 351321, 351320, 351319, 351314, 351356, 351355, 351354, 351352, 351351, 351350, 351349, 351348, 351347, 351346, 351345, 351344, 351343, 351342, 351341, 351340, 351339, 351338, 351337, 351336, 351335, 351365, 351364, 351363, 351362, 351361, 351360, 351359, 351358, 351357, 351368, 351367, 351366, 351381, 351380, 351379}
	blockIDs := []uint64{}
	result := postgres.GetDB().Table("blocks").Where("id in ?", blockIDs).Find(&blocks)

	fmt.Println(result)
	//_ = postgres.GetDB().Table("blocks").Where("updated_at > ? ", "2022-09-10 00:04:12.738103 +00:00").Find(&blocks).Error
	for i, v := range *blocks {

		fmt.Println("Processing Blocck #", i, len(*blocks))
		var p, u []UserBlockContributor
		_ = json.Unmarshal(v.Contributors, &p)
		fmt.Println("Number of Attribs ", len(p))
		//if len(p) == 3 {
		//	continue
		//}
		u = BuildUniqueContribs(p)
		fmt.Println("Fixing for --------------------------------------------", v.ID)
		mergedContrib, _ := json.Marshal(u)
		fmt.Println(u)
		LocalUserPushBlockContributors(v.ID, mergedContrib)
	}

	g.JSON(200, gin.H{
		"message": "Hello",
	})
}
func LocalUserPushBlockContributors(blockID uint64, contrib datatypes.JSON) {
	updates := map[string]interface{}{
		"contributors": contrib,
	}
	_ = postgres.GetDB().Model(&models.Block{}).Where("id = ?", blockID).Updates(&updates)
}

func SadLife(g *gin.Context) {
	if configs.GetConfigString("APP_MODE") == "production" {
		g.JSON(200, gin.H{
			"message": "abort",
		})

	}

	//Get the User ID
	userIDParam := g.Param("userid")
	userId, _ := strconv.ParseUint(userIDParam, 10, 64)

	// get all studios created by this User
	// Get All the studios associated
	var studios *[]models.Studio
	var members *[]models.Member
	_ = postgres.GetDB().Model(models.Studio{}).Where("created_by_id = ?", userId).Find(&studios).Error
	for _, studio := range *studios {
		fmt.Println(studio)
		// Delete instances of studio_topics where studio_id = studio.id
		//Delete studio_topics for this Studio
		postgres.GetDB().Model(&studio).Association("Topics").Clear()
		//Get All Members
		//Get All Roles
		var roles *[]models.Role
		//Delete the role_members (with this MemberID)
		_ = postgres.GetDB().Model(models.Member{}).Where("studio_id = ?", studio.ID).Find(&members).Error
		_ = postgres.GetDB().Model(models.Role{}).Where("studio_id = ?", studio.ID).Find(&roles).Error

		for _, member := range *members {
			postgres.GetDB().Model(&member).Association("Roles").Clear()
			postgres.GetDB().Delete(&models.Member{}, member.ID)
		}
		for _, role := range *roles {
			postgres.GetDB().Model(&role).Association("Members").Clear()
			postgres.GetDB().Delete(&models.Role{}, role.ID)
		}
		//Delete the role_members (with this MemberID)
		//Delete Studio
		postgres.GetDB().Delete(&models.Studio{}, studio.ID)
		//Delete user_associated_studios
	}
	_ = postgres.GetDB().Model(models.Member{}).Where("user_id = ?", userId).Find(&members).Error
	for _, member := range *members {
		postgres.GetDB().Model(&member).Association("Roles").Clear()
		postgres.GetDB().Delete(&models.Member{}, member.ID)
	}
	// Delete all messages
	var messages *[]models.Message
	_ = postgres.GetDB().Model(models.Message{}).Where("user_id = ?", userId).Delete(&messages).Error

	//Get All the user_social_auths Delete
	var socialAuth *[]models.UserSocialAuth
	_ = postgres.GetDB().Model(models.UserSocialAuth{}).Where("user_id = ?", userId).Delete(&socialAuth).Error

	var userAssStudios *[]models.UserAssociatedStudio
	_ = postgres.GetDB().Model(models.UserAssociatedStudio{}).Where("user_id = ?", userId).Delete(&userAssStudios).Error

	postgres.GetDB().Delete(&models.User{}, userId)

	g.JSON(200, gin.H{
		"missions": userId,
		"account":  "delete",
	})
}

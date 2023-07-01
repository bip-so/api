package core

import (
	"encoding/json"
	"fmt"
	"time"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"gorm.io/gorm"
)

type QuerySet struct {
}

/*
	GetByID This is just an example we can extend this more.

	Args:
		tableName string
		entityID uint64
	Response:
		resultString []byte
		err error

	Example:

	repo := core.QuerySet{}
	var collection models.Collection
	collectionInstance, _ := repo.GetByID("collections", 170)
	json.Unmarshal(collectionInstance, &collection)
	// Collection will be you models struct.
	fmt.Println(collection.ID, collection.StudioID, collection)
*/
func (repo QuerySet) GetByID(tableName string, entityID uint64) ([]byte, error) {
	var result map[string]interface{}
	postgres.GetDB().Table(tableName).Where("id = ?", entityID).Take(&result)
	resultString, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return resultString, nil
}

/*
	GetByParams This is just an example we can extend this more.

	Args:
		tableName string
		query map[string]interface{}
	Response:
		resultString []byte
		err error

	Example:

	repo := core.QuerySet{}
	var collections []models.Collection
	collectionInstance, _ := repo.GetByParams("collections", map[string]interface{}{"studio_id": 110})
	json.Unmarshal(collectionInstance, &collections)
	fmt.Println(cole, reflect.TypeOf(collections))
	for _, c := range collections {
		fmt.Println(c.ID)
	}
*/
func (repo QuerySet) GetByParams(tableName string, query map[string]interface{}) ([]byte, error) {
	var result []map[string]interface{}
	postgres.GetDB().Table(tableName).Where(query).Find(&result)
	resultString, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return resultString, nil
}

/*
	UpdateEntityByID Updates the record by ID and the fields mentioned

	Args:
		tableName string
		entityID uint64
		updateMap map[string]interface{}
	Response:
		err error

	Example:

	repo := core.QuerySet{}
	err := repo.UpdateEntityByID("collections", 170, map[string]interface{}{"publicAccess": "view"})
*/
func (repo QuerySet) UpdateEntityByID(tableName string, entityID uint64, updateMap map[string]interface{}) error {
	err := postgres.GetDB().Table(tableName).Where("id", entityID).Updates(updateMap).Error
	return err
}

/*
	SoftDeleteByID This is reversible change.

	Args:
		tableName string
		entityID uint64
		userID uint64
	Response:
		err error

	Example:

	repo := core.QuerySet{}
	err := repo.SoftDeleteByID("collections", 170, 56)
*/
func (repo QuerySet) SoftDeleteByID(tableName string, entityID uint64, userID uint64) error {
	result := map[string]interface{}{
		"is_archived":    true,
		"archived_at":    time.Now(),
		"archived_by_id": userID,
	}
	err := postgres.GetDB().Table(tableName).Where("id", entityID).Updates(result).Error

	return err
}

/*
	HardDeleteByID This will irreversible change

	Args:
		tableName string
		entityID uint64
	Response:
		err error

	Example:

	repo := core.QuerySet{}
	err := repo.HardDeleteByID("collections", 170)
*/
func (repo QuerySet) HardDeleteByID(tableName string, entityID uint64) error {
	result := map[string]interface{}{"id": entityID}
	err := postgres.GetDB().Table(tableName).Where("id", entityID).Delete(result).Error
	return err
}

// RecordExists:
func (repo QuerySet) RecordExists(tableName string, id uint64) bool {
	var exists bool
	postgres.GetDB().Table(tableName).Select("count(*) > 0").Where("id = ?", id).Find(&exists)
	return exists
}

/*
Function will update CommentCount On the Given model but the field should be `comment_count`
Please Note: it expects the Model Fields to be comment_count
Example: 	_ = r.Manager.CommentCountPlus(models.BLOCK_THREAD, instance.ThreadID)

*/
func (repo QuerySet) CommentCountPlus(tableName string, entityID uint64) error {
	err := postgres.GetDB().Table(tableName).Where("id", entityID).UpdateColumn("comment_count", gorm.Expr("comment_count  + ?", 1)).Error
	fmt.Println("CommentCountPlus Erer")
	fmt.Println(err)
	return err
}
func (repo QuerySet) ReelCountPlus(tableName string, entityID uint64) error {
	//result := map[string]interface{}{"id": entityID}
	err := postgres.GetDB().Table(tableName).Where("id", entityID).UpdateColumn("reel_count", gorm.Expr("reel_count  + ?", 1)).Error
	return err
}
func (repo QuerySet) ReelCountMinus(tableName string, entityID uint64) error {
	//result := map[string]interface{}{"id": entityID}
	err := postgres.GetDB().Table(tableName).Where("id", entityID).UpdateColumn("reel_count", gorm.Expr("reel_count  - ?", 1)).Error
	return err
}
func (repo QuerySet) CommentCountMinus(tableName string, entityID uint64) error {
	//result := map[string]interface{}{"id": entityID}
	err := postgres.GetDB().Table(tableName).Where("id", entityID).UpdateColumn("comment_count", gorm.Expr("comment_count  - ?", 1)).Error
	return err
}

/*
returns array of emoji name and its count
*/
func (repo QuerySet) GetEmojiCounter(tableName string, query string) ([]models.ReactionCounter, error) {
	var counter []models.ReactionCounter
	err := postgres.GetDB().Table(tableName).Select("Count(emoji) as count, emoji").Where(query).Group("emoji").Find(&counter).Error
	return counter, err
}

// Returns the count based on query
func (repo QuerySet) GetCount(tableName string, query map[string]interface{}) (int64, error) {
	var count int64
	err := postgres.GetDB().Table(tableName).Where(query).Count(&count).Error
	return count, err
}

type UserBlockContributor struct {
	Id        uint64    `json:"id"`
	UUID      string    `json:"uuid"`
	FullName  string    `json:"fullName"`
	Username  string    `json:"username"`
	AvatarUrl string    `json:"avatarUrl"`
	Timestamp time.Time `json:"timestamp"`
	BranchID  uint64    `json:"branchID"`
}

/*
This function is essentiually a Syntatic Sugar for
	update blocks
	set "field"= "field" || '{"key2": "value2"}' // valid json string
	where id = 66
*/
func (repo QuerySet) UserPushBlockContributors(blockID uint64, contrib string) {
	query := "update blocks set contributors = contributors || @contrib where id = @blockid"
	_ = postgres.GetDB().Raw(query, map[string]interface{}{"contrib": contrib, "blockid": blockID}).Error
}

func (repo *QuerySet) GetAllStudioIDsByUserID(userID uint64) []uint64 {
	var members []models.Member
	_ = postgres.GetDB().Model(&models.Member{}).Where("user_id = ? AND has_left = false AND is_removed = false", userID).Preload("Studio").Find(&members)
	//var studios []models.Studio
	var studiosArray []uint64
	for _, member := range members {
		studiosArray = append(studiosArray, member.Studio.ID)
		//studios = append(studios, *member.Studio)
	}
	return studiosArray
}

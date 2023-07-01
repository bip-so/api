package search

import (
	"fmt"
	"strconv"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
)

var algoliaClient *search.Client

func InitAlgolia() {
	configuration := search.Configuration{
		AppID:  configs.GetAlgoliaConfig().AppID,
		APIKey: configs.GetAlgoliaConfig().AdminAPIKey,
	}
	algoliaClient = search.NewClientWithConfig(configuration)
}

type AlgoliaIndex struct {
	index *search.Index
}

func GetAlgoliaClient() *search.Client {
	return algoliaClient
}

func GetIndex(indexName string) *AlgoliaIndex {
	return &AlgoliaIndex{
		index: algoliaClient.InitIndex(indexName),
	}
}

func (idx AlgoliaIndex) SaveRecord(object interface{}) error {
	opts := opt.AutoGenerateObjectIDIfNotExist(true)
	_, err := idx.index.SaveObject(object, opts)
	if err != nil {
		fmt.Println(err)
		logger.Error("error saving object in " + idx.index.GetName())
	}
	return err
}

func (idx AlgoliaIndex) SaveRecords(objects []interface{}) error {
	opts := opt.AutoGenerateObjectIDIfNotExist(true)
	_, err := idx.index.SaveObjects(objects, opts)
	if err != nil {
		logger.Error("error saving object in " + idx.index.GetName() + " : " + err.Error())
	}
	return err
}

func (idx AlgoliaIndex) DeleteRecordByID(documentID uint64) error {
	_, err := idx.index.DeleteBy(opt.Filters("id:" + strconv.FormatUint(documentID, 10)))
	if err != nil {
		logger.Error("error deleting object in " + idx.index.GetName() + " : " + err.Error())
	}
	return err
}

func (idx AlgoliaIndex) DeleteRecordByIDs(documentIDs []string) error {
	_, err := idx.index.DeleteObjects(documentIDs)
	if err != nil {
		logger.Error("error deleting object in " + idx.index.GetName() + " : " + err.Error())
	}
	return err
}

func (idx AlgoliaIndex) SetupIndex(searchableAttributes []string) {
	_, err := idx.index.SetSettings(search.Settings{
		SearchableAttributes: opt.SearchableAttributes(
			searchableAttributes...,
		),
	})
	if err != nil {
		logger.Error("error in setup index in " + idx.index.GetName() + " : " + err.Error())
	}
}

func (idx AlgoliaIndex) Search(query string, skipInt int) (search.QueryRes, error) {
	opts := []interface{}{
		opt.Offset(skipInt),
		opt.Length(configs.PAGINATION_LIMIT),
	}
	queryResult, err := idx.index.Search(query, opts...)
	if err != nil {
		logger.Error("error deleting object in " + idx.index.GetName() + " : " + err.Error())
	}
	return queryResult, err
}

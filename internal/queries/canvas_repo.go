package queries

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

func (q canvasRepoQuery) GetCanvasRepos(query map[string]interface{}) ([]models.CanvasRepository, error) {
	var repo []models.CanvasRepository
	err := postgres.GetDB().Model(&models.CanvasRepository{}).Where(query).Preload("DefaultBranch").Order("position ASC").Find(&repo).Error
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (q canvasRepoQuery) GetCanvasRepoInstance(query map[string]interface{}) (*models.CanvasRepository, error) {
	var repo *models.CanvasRepository
	err := postgres.GetDB().Model(&models.CanvasRepository{}).Where(query).First(&repo).Error
	if err != nil {

		return nil, err
	}
	return repo, nil
}

func (q canvasRepoQuery) QueryDB(search string, studioID uint64) *[]models.CanvasRepoFullRow {
	var finalResult []models.CanvasRepoFullRow
	var onerow models.CanvasRepoFullRow
	search = "%" + search + "%"
	query := "select canvas_repositories.*, " +
		"collections.id as id2, " +
		"collections.uuid as uuid2, " +
		"collections.name as name2, " +
		"collections.icon as icon2, " +
		"collections.position as position2, " +
		"collections.public_access as public_access2, " +
		"collections.computed_root_canvas_count as computed_root_canvas_count2 " +
		"from canvas_repositories left join collections " +
		"ON canvas_repositories.collection_id = collections.id " +
		"where canvas_repositories.studio_id = ? \n" +
		"and canvas_repositories.is_archived = false \n" +
		"and collections.is_archived = false \n" +
		"and ( collections.name ILIKE ?" +
		"or canvas_repositories.name ILIKE ?) order by collections.position, canvas_repositories.position;"
	rows, _ := postgres.GetDB().Raw(query, studioID, search, search).Rows()
	defer rows.Close()
	for rows.Next() {
		postgres.GetDB().ScanRows(rows, &onerow)
		finalResult = append(finalResult, onerow)
	}
	return &finalResult
}

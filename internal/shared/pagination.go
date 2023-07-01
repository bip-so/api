package shared

import (
	"gitlab.com/phonepost/bip-be-platform/pkg/api/apiutil"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"math"
)

type PaginationData struct {
	NextPage     int         `json:"nextPage"`
	PreviousPage int         `json:"previousPage"`
	CurrentPage  int         `json:"currentPage"`
	TotalPages   int         `json:"totalPages"`
	Offset       int         `json:"offset"`
	Data         interface{} `json:"data"`
}

//

func GetPaginationData(page int, query map[string]interface{}, model interface{}, data interface{}) PaginationData {
	var totalRows int64
	offset := (page - 1) * apiutil.SharedPaginationPerPage
	postgres.GetDB().Model(model).Where(query).Count(&totalRows)
	// totalPage = (int) Math.Ceiling((double) imagesFound.Length / PageSize);
	var perPageStart float64
	perPageStart = float64(totalRows) / apiutil.SharedPaginationPerPage
	totalPages := math.Ceil(perPageStart)
	// Debug

	CurrentPageCount := page
	TotalPages := int(totalPages)
	NextPageCount := page + 1

	if CurrentPageCount == int(totalPages) {
		NextPageCount = -1
	}

	return PaginationData{
		NextPage:     NextPageCount,
		PreviousPage: page - 1,
		CurrentPage:  page,
		TotalPages:   TotalPages,
		Offset:       offset,
		Data:         data,
	}
}

type SimpleListData struct {
	Data interface{} `json:"data"`
}

func GetSimpleListData(data interface{}) SimpleListData {
	return SimpleListData{
		Data: data,
	}
}

type SimpleStringData struct {
	Message string `json:"message"`
}

func GetStringMessageData(data string) SimpleStringData {
	return SimpleStringData{
		Message: data,
	}
}

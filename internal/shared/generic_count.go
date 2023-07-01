package shared

import (
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

func CountResults(query map[string]interface{}, model interface{}) int64 {
	var totalRows int64
	postgres.GetDB().Model(model).Where(query).Count(&totalRows)
	return totalRows
}

package shared

type NewGenericResponseV1 struct {
	Data interface{} `json:"data"`
	Next string      `json:"next"`
}

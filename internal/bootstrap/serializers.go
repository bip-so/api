package bootstrap

import "encoding/json"

type AssociatedStudio struct {
	ID          uint64 `json:"id"`
	UUID        string `json:"uuid"`
	Handle      string `json:"handle"`
	DisplayName string `json:"displayName"`
	ImageURL    string `json:"imageUrl"`
	CreatedByID uint64 `json:"created_by_id"`
}

func SerializeAssociatedStudios(data string) *[]AssociatedStudio {
	var associatedStudios []AssociatedStudio
	err := json.Unmarshal([]byte(data), &associatedStudios)
	if err != nil {
		return nil
	}
	return &associatedStudios
}

package models

func (m *UserAssociatedStudio) TableName() string {
	return "user_associated_studios"
}

type UserAssociatedStudio struct {
	BaseModel
	UserID      uint64 `json:"user_id"`
	StudiosData string `json:"studios"`
}

func NewUserAssociatedStudio(userID uint64, studiosData string) *UserAssociatedStudio {
	return &UserAssociatedStudio{
		UserID:      userID,
		StudiosData: studiosData,
	}
}

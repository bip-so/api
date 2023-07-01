package models

func (m *Topic) TableName() string {
	return "topics"
}

type Topic struct {
	BaseModel
	Name    string   `gorm:"type: varchar(100)"`
	Studios []Studio `gorm:"many2many:studio_topics;"`
}

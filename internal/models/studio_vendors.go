package models

type StudioVendor struct {
	BaseModel
	GuildId           string `json:"guildId"`
	GuildName         string `json:"guildName"`
	PartnerName       string `json:"partnerName"`
	IntegrationStatus string `gorm:"type: varchar(20)"`
	Handle            string `json:"handle"`
}

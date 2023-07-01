package models

import (
	"database/sql"
	"github.com/lib/pq"
	"time"
)

type BlockOg struct {
	ID                          string
	Type                        string
	Text                        string
	UserID                      string
	TweetID                     sql.NullString
	PageID                      string
	Position                    int
	URL                         sql.NullString
	CreatedAt                   time.Time
	UpdatedAt                   time.Time
	Properties                  string
	BranchID                    string
	CommentCount                int
	ReactionCountList           string
	BlockTextHighlightCountList string
	LastUpdatedAttributedUserID sql.NullString
	LastUpdatedAttributedTime   time.Time
	ReelCount                   int            `gorm:"default:0"`
	MentionedUserIDs            pq.StringArray `gorm:"type:text"`
	MentionedGroupIDs           pq.StringArray `gorm:"type:text"`
	MentionedPageIDs            pq.StringArray `gorm:"type:text"`
}

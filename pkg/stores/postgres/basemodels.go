package postgres

import (
	"time"

	"github.com/google/uuid"
)

// PG: shared model for all tables
type BaseModel struct {
	ID        uint64    `gorm:"primaryKey,type:BIGSERIAL UNSIGNED NOT NULL AUTO_INCREMENT"`
	UUID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// Commented
//type OwnershipMixin struct {
//	CreatedByID uint64 `gorm:"type:BIGSERIAL"`
//	UpdatedByID uint64 `gorm:"type:BIGSERIAL"`
//}
//
//type SoftDeleteMixin struct {
//	IsArchived   bool
//	ArchivedAt   time.Time
//	ArchivedByID int `gorm:"type:BIGSERIAL"`
//}

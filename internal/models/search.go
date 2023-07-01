package models

import (
	"github.com/google/uuid"
	"time"
)

type CanvasRepoFullRow struct {
	ID                          uint64
	UUID                        uuid.UUID
	CreatedAt                   time.Time
	UpdatedAt                   time.Time
	CollectionID                uint64
	StudioID                    uint64
	Name                        string
	Position                    uint
	Icon                        string
	PublicAccess                string
	IsPublished                 bool
	DefaultBranchID             uint64
	ParentCanvasRepositoryID    *uint64
	CreatedByID                 uint64
	UpdatedByID                 uint64
	IsArchived                  bool
	ArchivedAt                  time.Time
	ArchivedByID                *uint64
	DefaultLanguageCanvasRepoID *uint64
	Language                    string
	IsLanguageCanvas            bool
	AutoTranslated              bool
	Key                         string
	SubCanvasCount              uint
	ID2                         uint64
	UUID2                       uuid.UUID
	Name2                       string
	Position2                   uint64
	PublicAccess2               string
	Icon2                       string
	HasPublicAccess             bool
	ComputedRootCanvasCount2    int
}

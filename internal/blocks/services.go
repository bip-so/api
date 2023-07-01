package blocks

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
)

func (s blockService) Create(b *models.Block) error {
	_, err := queries.App.BlockQuery.CreateBlock(b)
	return err
}

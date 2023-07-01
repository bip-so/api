package workflows

import "gitlab.com/phonepost/bip-be-platform/internal/models"

func InvalidateCaching(event string) {
	switch event {
	case models.FECanvasBranchCaching:
		// pass call an FE API for caching
	}
}

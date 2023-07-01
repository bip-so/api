package blocks

import (
	"github.com/gin-gonic/gin"
)

func (a blockApp) RegisterRoutes(r *gin.RouterGroup) {
	App.Routes = *r.Group("block")
	{

	}
}

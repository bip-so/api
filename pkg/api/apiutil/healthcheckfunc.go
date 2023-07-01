package apiutil

import (
	"github.com/gin-gonic/gin"
)

func HealthCheck(g *gin.Context) {
	//val, _ := g.Get("X-Revision")
	//err := core.CoreHealth()
	//if err != nil {
	//	g.JSON(503, gin.H{"status": "DOWN", "reason": err.Error()})
	//	return
	//}
	//scheme := "http"
	//if g.Request.TLS != nil {
	//	scheme = "https"
	//}

	g.JSON(200, gin.H{
		"status": "OKay",
		//"version": val.(string),
		//"swagger": scheme + "://" + g.Request.Host + "/" + "swagger/index.html",
	})
}

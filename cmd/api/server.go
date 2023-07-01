package api

import (
	br "github.com/anargu/gin-brotli"
	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/urls"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
)

func RouterSetUp() *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())

	// Nitish and Chirag have tested don't change.
	AwesomeCompression := br.Options{
		WriterOptions: brotli.WriterOptions{
			Quality: 5,
			LGWin:   0,
		},
		SkipExtensions: []string{".png", ".gif", ".jpeg", ".jpg", ".mp3", ".mp4"},
	}
	router.Use(br.Brotli(AwesomeCompression))
	//router.Use(gzip.Gzip(gzip.BestCompression))
	router.Use(middlewares.Benchmark())
	// This will set request_id, studio_id, currentUser
	router.Use(middlewares.SauronMiddleware())
	// Paused
	//router.Use(middlewares.SentryMiddleware())
	return router
}

func StartHttpServer() {
	r := RouterSetUp()
	urls.Router(r)

	// Generic Page not found.
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	// Multiple Listners
	//go http.ListenAndServe(PORT, handlerA)
	//http.ListenAndServe(PORT, handlerB)
	if configs.GetConfigString("APP_MODE") == "local" {
		r.Run("localhost:9001")
	} else {
		r.Run(":9001")
	}
}

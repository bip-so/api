package discord

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
)

func ErrorHandlerMessage(c *gin.Context) {
	response.RenderCustomResponse(c, map[string]interface{}{
		"type": 4,
		"data": map[string]interface{}{
			"tts":     false,
			"content": "Error capturing the message. Please try again.",
			"embeds":  []string{},
			"allowed_mentions": map[string]interface{}{
				"parse": []string{},
			},
			"flags": 1 << 6,
		},
	})
}

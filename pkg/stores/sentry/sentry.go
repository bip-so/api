package sentry

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
)

func InitSentry() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              configs.GetConfigString(configs.SENTRY_DSN),
		Environment:      configs.GetConfigString(configs.ENV_KEY),
		Debug:            true,
		AttachStacktrace: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)
}

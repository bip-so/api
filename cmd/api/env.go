package api

import (
	"os"
)

//os.environ.setdefault("BIP_ENV", "api.settings")
func GetEnvFileName() string {
	foo := os.Getenv("APP_MODE") + ".env"
	return foo
}

package main

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/aws"
)

func main() {
	configs.InitConfig(".env", ".")
	translated, err := aws.Translate("hello change me", "en", "hi")
	if err != nil {
		fmt.Println("error", err)
	}
	fmt.Println("Translated text ", translated)
}

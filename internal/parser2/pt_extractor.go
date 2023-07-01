package parser2

import (
	"encoding/json"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

func ExtractInnerBlockData(mainBlock []map[string]interface{}) string {
	for _, value := range mainBlock {
		if _, ok := value["type"]; !ok {
			continue
		} else {
			return value["type"].(string)
		}
	}
	//fmt.Println(mainBlock)
	return AssumedTextBlockType
}

type Attr struct {
	Level        int    `json:"level"`
	Checked      bool   `json:"checked"`
	CodeLanguage string `json:"codeLanguage"`
}

func ExtractAttributes(block models.Block) Attr {
	a := Attr{}
	_ = json.Unmarshal(block.Attributes, &a)
	return a
}

package parser2

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"os"
)

const AssumedTextBlockType = "assumedText"

func (p Parser2) NewParser(branchID uint64) *Parser2 {
	blocks := App.Service.GetBranchBlocks(branchID)
	return &Parser2{
		BranchID: branchID,
		Blocks:   blocks,
	}
}

func (p Parser2) SaveToFile(filePath string, values []string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, value := range values {
		fmt.Fprintln(f, value) // print values to f, one per line
	}
	return nil
}

func GenericBlockHandler(blockType string, block models.Block) string {
	blockMap := UtilConvertByteJsonToMMap(block.Children)
	actualBlockType := ExtractInnerBlockData(blockMap)
	attributes := ExtractAttributes(block)
	return ExtractText(blockMap, blockType, actualBlockType, attributes)
}

func (p Parser2) PlainText(id uint64) {
	p2 := p.NewParser(id)
	var finalStrArray []string
	//var previousBlockType string
	finalStrArray = append(finalStrArray, "### This will be title to be added later")
	finalStrArray = append(finalStrArray, Markdown_Heading_Prefix)
	finalStrArray = append(finalStrArray, "\n\n")

	for _, v := range *p2.Blocks {
		//if previousBlockType == "text" && (v.Type == "ulist" || v.Type == "olist" || v.Type == "checklist") {
		//	finalStrArray = append(finalStrArray, "\n")
		//}

		tempFromBlock := GenericBlockHandler(v.Type, v)
		finalStrArray = append(finalStrArray, tempFromBlock)
		//previousBlockType = v.Type
	}
	//for _, v1 := range finalStrArray {
	//	fmt.Println(v1)
	//}
	p.SaveToFile("test.md", finalStrArray)

}

package models

// Talk to PW
const (
	// BlockTypeHeading1 defines heading1
	BlockTypeHeading1 = "heading1"
	// BlockTypeHeading2 defines heading2
	BlockTypeHeading2 = "heading2"
	// BlockTypeHeading3 defines heading3
	BlockTypeHeading3 = "heading3"
	// BlockTypeHeading4 defines heading4
	BlockTypeHeading4 = "heading4"
	// BlockTypeHeading5 defines heading5
	BlockTypeHeading5 = "heading5"
	// BlockTypeHeading6 defines heading6
	BlockTypeHeading6 = "heading6"
	// BlockTypeText defines text
	BlockTypeText = "text"
	// BlockTypeText defines code
	BlockTypeCode = "code"
	// BlockTypeText defines callout
	BlockTypeCallout = "callout"
	// BlockTypeTweet defines tweet
	BlockTypeTweet = "tweet"
	// BlockTypeImage defines image
	BlockTypeImage = "image"
	// BlockTypeUList defines unordered list
	BlockTypeUList = "ulist"
	// BlockTypeOList defines ordered list
	BlockTypeOList = "olist"
	// BlockTypeOList defines ordered list
	BlockTypeListItem = "listitem"
	// BlockTypeAttachment defines attachment
	BlockTypeAttachment = "attachment"
	// BlockTypeVideo defines video
	BlockTypeVideo = "video"
	// BlockTypeCheckList defines todo list
	BlockTypeCheckList = "checklist"

	// BlockContentTypeName this will be used to store
	// *** DO NOT CHANGE THIS VALUE *** Block struct name in content Type model
	BlockContentTypeName = "Block"
	// BlockTypeGoogleMap defines google map
	BlockTypeGoogleMap = "googlemap"
	// BlockTypeFigma defines figma
	BlockTypeFigma = "figma"
	//BlockTypeCodePen defines codepen
	BlockTypeCodePen = "codepen"
	//BlockTypeReplit defines replit
	BlockTypeReplit = "replit"
	//BlockTypeCodeSandBox defines codesandbox
	BlockTypeCodeSandBox = "codesandbox"
	//BlockTypeGoogleDrive defines googledrive
	BlockTypeGoogleDrive = "googledrive"
	//BlockTypeGoogleSheet defines googlesheet
	BlockTypeGoogleSheet = "googlesheet"
	//BlockTypeMiro defines miro
	BlockTypeMiro = "miro"
	//BlockTypeExcaliDraw defines excalidraw
	BlockTypeExcaliDraw = "excalidraw"
	//BlockTypeOgMetaTag defines ogmetatag
	BlockTypeOgMetaTag = "ogmetatag"
	//BlockTypeOgMetaTag defines ogmetatag
	BlockTypeYoutube = "youtube"
	//BlockTypeOgMetaTag defines ogmetatag
	BlockTypeLoom = "loom"
	//BlockTypeDivider defines ogmetatag
	BlockTypeDivider = "hr"
	//BlockTypeTempTable defines tempTable used in importing notion docs
	BlockTypeTempTable = "temptable"
	//BlockTypeTempTable defines tempTable used in importing notion docs
	BlockTypeImportNotionMention = "importNotionMention"
	BlockTypeSimpleTable         = "simpletable"
	BlockMentions                = "mentions"
	BlockTypeTOC                 = "toc"

	BlockTypeBIPMark = "bipmark"
	BlockTypeSubtext = "subtext"
	BlockDrawIO      = "drawio"

	// BlockSimpleTableV1 BlockTable
	BlockSimpleTableV1 = "simple_table_v1"
	BlockSubtext       = "subtext"
	BlockQuote         = "quote"
	//subtext
	//quote
)

var AllowedBlockTypes = []string{
	BlockTypeHeading1,
	BlockTypeHeading2,
	BlockTypeHeading3,
	BlockTypeHeading4,
	BlockTypeHeading5,
	BlockTypeHeading6,
	BlockTypeText,
	BlockTypeCode,
	BlockTypeCallout,
	BlockTypeTweet,
	BlockTypeImage,
	BlockTypeUList,
	BlockTypeOList,
	BlockTypeListItem,
	BlockTypeAttachment,
	BlockTypeVideo,
	BlockTypeCheckList,
	BlockTypeGoogleMap,
	BlockTypeFigma,
	BlockTypeCodePen,
	BlockTypeReplit,
	BlockTypeCodeSandBox,
	BlockTypeGoogleDrive,
	BlockTypeGoogleSheet,
	BlockTypeMiro,
	BlockTypeExcaliDraw,
	BlockTypeOgMetaTag,
	BlockTypeYoutube,
	BlockTypeLoom,
	BlockTypeDivider,
	BlockTypeImportNotionMention,
	BlockTypeTempTable,
	BlockTypeSimpleTable,
	BlockMentions,
	BlockTypeTOC,
	BlockTypeBIPMark,
	BlockTypeSubtext,
	BlockDrawIO,
	BlockSimpleTableV1,
	BlockSubtext,
	BlockQuote,
}

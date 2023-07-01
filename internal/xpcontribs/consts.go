package xpcontribs

var XPEventGroupMerge = "xp_merge"
var XPEventGroupPublish = "xp_publish"

var XP_EVENT_GROUP_MAPPER = map[string]map[string]int{
	XPEventGroupMerge: {
		"NEW_BLOCK":    5,
		"EDITED_BLOCK": 1,
		"DELETE_BLOCK": 0,
	},
	XPEventGroupPublish: {
		"NEW_BLOCK":    5,
		"EDITED_BLOCK": 1,
		"DELETE_BLOCK": 0,
	},
}

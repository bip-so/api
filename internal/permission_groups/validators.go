package permissiongroup

const PGSchemaVersion = 1

type PermissionsSchemaResponse struct {
	Version     int                   `json:"version"`
	Group       string                `json:"group"`
	Permissions []PermissionsTemplate `json:"permissionGroups"`
}

type PermissionsSchemaResponseArray struct {
	Type string
	Data []PermissionsSchemaResponse
}

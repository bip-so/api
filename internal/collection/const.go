package collection

type CollectionPublicAccessEnum int

const (
	Private CollectionPublicAccessEnum = iota
	View
	Comment
	Edit
)

func (cpa CollectionPublicAccessEnum) String() string {
	return []string{"private", "view", "comment", "edit"}[cpa]
}

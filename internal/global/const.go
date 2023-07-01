package global

//import "context"
//
//func BipGlobarContext() (context.Context, context.Context) {
//	var ctx context.Context
//	ctxbg := context.Background()
//	return ctx, ctxbg
//}

func GetCurrentSessionStudioId() uint64 {
	var id uint64

	return id
}

// CanvasPublicAccessEnum
type CanvasPublicAccessEnum int

const (
	Private CanvasPublicAccessEnum = iota
	View
	Comment
	Edit
)

func (cpa CanvasPublicAccessEnum) String() string {
	return []string{"private", "view", "comment", "edit"}[cpa]
}

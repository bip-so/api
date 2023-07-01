package queries

import (
	"fmt"

	"gitlab.com/phonepost/bip-be-platform/pkg/core"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/kafka"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	bipredis "gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"gorm.io/gorm"
)

type userQuery struct {
	Manager core.QuerySet
	db      *gorm.DB
}
type studioQuery struct {
}
type memberQuery struct {
}
type blockQuery struct {
	Manager core.QuerySet
}
type repoQuery struct {
	kafka *kafka.BipKafka
}
type canvasRepoQuery struct {
}

type collectionQuery struct {
}

type branchQuery struct {
}
type permsQuery struct {
	cache *bipredis.Cache
}

type studioMemberRequestQuery struct{}
type studioInviteEmailQuery struct{}
type branchInviteViaEmailQuery struct{}
type roleQuery struct {
}
type stripeQuery struct {
}
type studioVendorQuery struct{}
type studioIntegrationQuery struct{}
type studioPermissionQuery struct{}
type canvasBranchPermissionQuery struct {
}
type accessRequestQuery struct {
}

type attributionQuery struct {
}
type publishRequestQuery struct {
}
type branchAccessTokenQuery struct {
	Manager core.QuerySet
}
type messageQuery struct {
}

type QueryApp struct {
	Name                        string
	UserQueries                 userQuery
	StudioQueries               studioQuery
	StudioMemberRequestQuery    studioMemberRequestQuery
	RoleQuery                   roleQuery
	MemberQuery                 memberQuery
	BlockQuery                  blockQuery
	RepoQuery                   repoQuery
	CollectionQuery             collectionQuery
	BranchQuery                 branchQuery
	BranchInviteViaEmailQuery   branchInviteViaEmailQuery
	PermsQuery                  permsQuery
	StudioInviteEmailQuery      studioInviteEmailQuery
	StripeQuery                 stripeQuery
	StudioVendorQuery           studioVendorQuery
	StudioIntegrationQuery      studioIntegrationQuery
	StudioPermissionQuery       studioPermissionQuery
	CanvasBranchPermissionQuery canvasBranchPermissionQuery
	AccessRequestQuery          accessRequestQuery
	AttributionQuery            attributionQuery
	PublishRequestQuery         publishRequestQuery
	BranchAccessTokenQuery      branchAccessTokenQuery
	CanvasRepoQuery             canvasRepoQuery
	MessageQuery                messageQuery
}

var App QueryApp

func InitApp() {
	App.Name = "Studio"
	App.UserQueries.db = postgres.GetDB()
	App.RepoQuery.kafka = kafka.GetKafkaClient()
	App.PermsQuery.cache = bipredis.NewCache()
	fmt.Println(App.Name + " started. ")
}

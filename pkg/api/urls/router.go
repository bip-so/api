package urls

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gitlab.com/phonepost/bip-be-platform/docs"
	"gitlab.com/phonepost/bip-be-platform/internal/auth"
	"gitlab.com/phonepost/bip-be-platform/internal/blocks"
	"gitlab.com/phonepost/bip-be-platform/internal/blockthread"
	blockThreadCommentcomment "gitlab.com/phonepost/bip-be-platform/internal/blockthreadcomment"
	"gitlab.com/phonepost/bip-be-platform/internal/bootstrap"
	"gitlab.com/phonepost/bip-be-platform/internal/canvasbranch"
	"gitlab.com/phonepost/bip-be-platform/internal/canvasbranchpermissions"
	"gitlab.com/phonepost/bip-be-platform/internal/canvasrepo"
	"gitlab.com/phonepost/bip-be-platform/internal/collection"
	"gitlab.com/phonepost/bip-be-platform/internal/collectionpermissions"
	"gitlab.com/phonepost/bip-be-platform/internal/discord"
	"gitlab.com/phonepost/bip-be-platform/internal/follow"
	"gitlab.com/phonepost/bip-be-platform/internal/global"
	"gitlab.com/phonepost/bip-be-platform/internal/member"
	"gitlab.com/phonepost/bip-be-platform/internal/mentions"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/internal/parser2"
	"gitlab.com/phonepost/bip-be-platform/internal/payments"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/internal/post"
	"gitlab.com/phonepost/bip-be-platform/internal/pr"
	"gitlab.com/phonepost/bip-be-platform/internal/reactions"
	"gitlab.com/phonepost/bip-be-platform/internal/reel"
	"gitlab.com/phonepost/bip-be-platform/internal/role"
	"gitlab.com/phonepost/bip-be-platform/internal/shortner"
	slack2 "gitlab.com/phonepost/bip-be-platform/internal/slack"
	"gitlab.com/phonepost/bip-be-platform/internal/studio"
	"gitlab.com/phonepost/bip-be-platform/internal/studio_integration"
	"gitlab.com/phonepost/bip-be-platform/internal/studiopermissions"
	"gitlab.com/phonepost/bip-be-platform/internal/twitter"
	"gitlab.com/phonepost/bip-be-platform/internal/user"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/apiutil"
)

func Router(r *gin.Engine) {
	r.GET("/", apiutil.Okay)
	r.GET("/sad-life-alice/:userid", apiutil.SadLife)
	r.GET("/fixblock", apiutil.BlockFixer)
	r.GET("/memberfix", apiutil.MemberFix)

	r.GET("/plain", parser2.Mdplain)
	r.GET("/local", apiutil.Tested)
	r.GET("/test-mailer", apiutil.TestMailer)
	r.GET("/health", apiutil.HealthCheck)
	//r.GET("/exp", apiutil.Exp)

	//r.Use(middlewares.Constraint())
	/* ---------------------------  Public Swagger routes  --------------------------- */
	api := r.Group("/")
	docs.SwaggerInfo.BasePath = "/api"
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	/* ---------------------------  bip Platform routes  --------------------------- */

	api_v1 := r.Group("/api/v1")
	auth.App.RegisterRoutes(api_v1)

	//api_v1.Use(middlewares.TokenAuthorizationMiddleware())
	user.App.RegisterUser(api_v1)
	studio.App.RegisterRoutes(api_v1)
	studiopermissions.RegisterStudioPermissionRoutes(api_v1)
	global.RegisterGlobalRoutes(api_v1)
	bootstrap.RegisterBootstrapRoutes(api_v1)
	permissiongroup.RegisterPermsGroupRoutes(api_v1)
	collection.App.RegisterCollectionRoutes(api_v1)
	follow.App.RegisterRoutes(api_v1)
	role.RegisterRoleCrudRoutes(api_v1)
	collectionpermissions.RegisterCollectionPermissionRoutes(api_v1)
	canvasbranch.App.RegisterRoutes(api_v1)
	canvasrepo.App.RegisterRoutes(api_v1)
	blocks.App.RegisterRoutes(api_v1)
	member.App.RegisterRoutes(api_v1)
	blockthread.App.RegisterRoutes(api_v1)
	blockThreadCommentcomment.App.RegisterRoutes(api_v1)
	permissions.App.RegisterRoutes(api_v1)
	canvasbranchpermissions.RegisterCanvasBranchPermissionRoutes(api_v1)
	reel.App.RegisterRoutes(api_v1)
	discord.RegisterRoutes(api)
	reactions.App.RegisterRoutes(api_v1)
	notifications.App.RegisterRoutes(api_v1)
	twitter.RegisterRoutes(api_v1)
	pr.App.RegisterRoutes(api_v1)
	mentions.App.RegisterRoutes(api_v1)
	shortner.App.RegisterRoutes(api_v1)
	studio_integration.App.RegisterRoutes(api_v1)
	parser2.App.RegisterRoutes(api_v1)
	slack2.App.RegisterRoutes(api_v1)
	post.App.RegisterRoutes(api_v1)
	payments.RegisterStripeRoutes(api_v1)
}

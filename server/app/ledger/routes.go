package ledger

import (
	"baize/app/common/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup) {
	public := router.Group("/ledger/auth")
	public.POST("/register", Register)
	public.POST("/login", Login)

	private := router.Group("/ledger")
	private.Use(AuthMiddleware())
	private.GET("/sync-consent", GetSyncConsent)
	private.PUT("/sync-consent", UpdateSyncConsent)
	private.GET("/profile", GetProfile)
	private.PUT("/profile", UpdateProfile)
	private.GET("/referrals", GetReferrals)
	private.DELETE("/account", DeleteAccount)
	private.PUT("/records", UploadRecords)
	private.POST("/records", CreateRecord)
	private.GET("/records", ReadRecords)
	private.PUT("/records/:recordId", UpdateRecord)
	private.DELETE("/records/:recordId", DeleteRecord)
	private.DELETE("/records/batches/:batchId", DeleteImportBatch)
	private.GET("/import-batches", ReadImportBatches)
	private.POST("/import-batches", SaveImportBatch)
	private.PUT("/import-batches", UploadImportBatches)
	private.DELETE("/import-batches/:batchId", DeleteImportBatch)

	admin := router.Group("/ledger/admin")
	admin.Use(middlewares.JWTAuthMiddleware())
	admin.GET("/accounts", middlewares.HasPermission("ledger:account:list"), AdminAccounts)
	admin.GET("/accounts/:accountId/records", middlewares.HasPermission("ledger:account:list"), AdminAccountRecords)
	admin.DELETE("/accounts/:accountId/records/:recordId", middlewares.HasPermission("ledger:account:removeRecord"), AdminDeleteAccountRecord)
}

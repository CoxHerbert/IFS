package controller

import (
	"baize/app/common/baize/baizeContext"
	customermiddleware "baize/app/customer/middleware"
	"baize/app/customer/models"
	"baize/app/customer/service"
	freightModels "baize/app/freight/models"
	freightService "baize/app/freight/service"
	"baize/app/utils/slicesUtils"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var customerService = service.GetCustomerService()
var shipmentService = freightService.GetShipmentService()
var receiptService = freightService.GetReceiptService()

type resetPasswordBody struct {
	Password string `json:"password" binding:"required"`
}

type accountRoleBody struct {
	RoleIds []string `json:"roleIds"`
}

func CustomerList(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	customer := new(models.CustomerDQL)
	c.ShouldBind(customer)
	if !freightService.CanManageAllShipments(bzc.GetCurrentUser()) {
		customer.SalesUserId = bzc.GetCurrentUserId()
	}
	customer.SetLimit(c)
	list, count := customerService.SelectCustomerList(customer)
	bzc.SuccessListData(list, count)
}

func CustomerGetInfo(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	customerId := bzc.ParamInt64("customerId")
	if customerId == 0 {
		bzc.ParameterError()
		return
	}
	if !canAccessCustomer(bzc, customerId) {
		bzc.Waring("无权查看该客户")
		return
	}
	bzc.SuccessData(customerService.SelectCustomerById(customerId))
}

func CustomerOptions(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	if freightService.CanManageAllShipments(bzc.GetCurrentUser()) {
		bzc.SuccessData(customerService.SelectCustomerOptions(c.Query("keyword")))
		return
	}
	bzc.SuccessData(customerService.SelectCustomerOptionsBySales(c.Query("keyword"), bzc.GetCurrentUserId()))
}

func CustomerAdd(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	customer := new(models.CustomerDML)
	if err := c.ShouldBindJSON(customer); err != nil {
		zap.L().Error("参数错误", zap.Error(err))
		bzc.ParameterError()
		return
	}
	customer.CreateBy = bzc.GetCurrentUserName()
	customer.UpdateBy = bzc.GetCurrentUserName()
	if !freightService.CanManageAllShipments(bzc.GetCurrentUser()) {
		customer.SalesUserId = bzc.GetCurrentUserId()
		customer.SalesUserName = bzc.GetCurrentUserName()
	}
	customerService.InsertCustomer(customer)
	bzc.Success()
}

func CustomerEdit(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	customer := new(models.CustomerDML)
	if err := c.ShouldBindJSON(customer); err != nil {
		zap.L().Error("参数错误", zap.Error(err))
		bzc.ParameterError()
		return
	}
	if !canAccessCustomer(bzc, customer.CustomerId) {
		bzc.Waring("无权修改该客户")
		return
	}
	if !freightService.CanManageAllShipments(bzc.GetCurrentUser()) {
		customer.SalesUserId = bzc.GetCurrentUserId()
		customer.SalesUserName = bzc.GetCurrentUserName()
	}
	customer.UpdateBy = bzc.GetCurrentUserName()
	customerService.UpdateCustomer(customer)
	bzc.Success()
}

func CustomerRemove(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	var ids slicesUtils.Slices = strings.Split(c.Param("customerIds"), ",")
	customerIds := ids.StrSlicesToInt()
	for _, customerId := range customerIds {
		if !canAccessCustomer(bzc, customerId) {
			bzc.Waring("无权删除所选客户")
			return
		}
	}
	customerService.DeleteCustomerByIds(customerIds)
	bzc.Success()
}

func ContactList(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	contact := new(models.CustomerContactDQL)
	c.ShouldBind(contact)
	contact.CustomerId = bzc.ParamInt64("customerId")
	if !canAccessCustomer(bzc, contact.CustomerId) {
		bzc.Waring("无权查看该客户联系人")
		return
	}
	contact.SetLimit(c)
	list, count := customerService.SelectContactList(contact)
	bzc.SuccessListData(list, count)
}

func canAccessCustomer(bzc *baizeContext.BaiZeContext, customerId int64) bool {
	customer := customerService.SelectCustomerById(customerId)
	if customer == nil {
		return false
	}
	return freightService.CanManageAllShipments(bzc.GetCurrentUser()) || customer.SalesUserId == bzc.GetCurrentUserId()
}

func ContactGetInfo(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	contactId := bzc.ParamInt64("contactId")
	if contactId == 0 {
		bzc.ParameterError()
		return
	}
	contact := customerService.SelectContactById(contactId)
	if contact == nil || !canAccessCustomer(bzc, contact.CustomerId) {
		bzc.Waring("无权查看该联系人")
		return
	}
	bzc.SuccessData(contact)
}

func ContactAdd(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	contact := new(models.CustomerContactDML)
	if err := c.ShouldBindJSON(contact); err != nil {
		zap.L().Error("参数错误", zap.Error(err))
		bzc.ParameterError()
		return
	}
	if !canAccessCustomer(bzc, contact.CustomerId) {
		bzc.Waring("无权维护该客户联系人")
		return
	}
	contact.CreateBy = bzc.GetCurrentUserName()
	contact.UpdateBy = bzc.GetCurrentUserName()
	if err := customerService.InsertContact(contact); err != nil {
		bzc.Waring(err.Error())
		return
	}
	bzc.Success()
}

func ContactEdit(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	contact := new(models.CustomerContactDML)
	if err := c.ShouldBindJSON(contact); err != nil {
		zap.L().Error("参数错误", zap.Error(err))
		bzc.ParameterError()
		return
	}
	existingContact := customerService.SelectContactById(contact.ContactId)
	if existingContact == nil || !canAccessCustomer(bzc, existingContact.CustomerId) {
		bzc.Waring("无权修改该联系人")
		return
	}
	contact.CustomerId = existingContact.CustomerId
	contact.UpdateBy = bzc.GetCurrentUserName()
	customerService.UpdateContact(contact)
	bzc.Success()
}

func ContactRemove(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	var ids slicesUtils.Slices = strings.Split(c.Param("contactIds"), ",")
	contactIds := ids.StrSlicesToInt()
	for _, contactId := range contactIds {
		contact := customerService.SelectContactById(contactId)
		if contact == nil || !canAccessCustomer(bzc, contact.CustomerId) {
			bzc.Waring("无权删除所选联系人")
			return
		}
	}
	customerService.DeleteContactByIds(contactIds)
	bzc.Success()
}

func AccountList(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	account := new(models.CustomerAccountDQL)
	c.ShouldBind(account)
	if !freightService.CanManageAllShipments(bzc.GetCurrentUser()) {
		account.SalesUserId = bzc.GetCurrentUserId()
	}
	account.SetLimit(c)
	list, count := customerService.SelectAccountList(account)
	bzc.SuccessListData(list, count)
}

func AccountGetInfo(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	accountId := bzc.ParamInt64("accountId")
	if accountId == 0 {
		bzc.ParameterError()
		return
	}
	account := customerService.SelectAccountById(accountId)
	if account == nil || !canAccessCustomer(bzc, account.CustomerId) {
		bzc.Waring("无权查看该客户账号")
		return
	}
	data := map[string]interface{}{
		"account": account,
		"roleIds": customerService.SelectAccountRoleIds(accountId),
	}
	bzc.SuccessData(data)
}

func AccountAdd(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	account := new(models.CustomerAccountDML)
	if err := c.ShouldBindJSON(account); err != nil {
		zap.L().Error("参数错误", zap.Error(err))
		bzc.ParameterError()
		return
	}
	if !canAccessCustomer(bzc, account.CustomerId) {
		bzc.Waring("无权新增该客户账号")
		return
	}
	account.CreateBy = bzc.GetCurrentUserName()
	account.UpdateBy = bzc.GetCurrentUserName()
	if err := customerService.InsertAccount(account); err != nil {
		bzc.Waring(err.Error())
		return
	}
	bzc.Success()
}

func AccountEdit(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	account := new(models.CustomerAccountDML)
	if err := c.ShouldBindJSON(account); err != nil {
		zap.L().Error("参数错误", zap.Error(err))
		bzc.ParameterError()
		return
	}
	existingAccount := customerService.SelectAccountById(account.AccountId)
	if existingAccount == nil || !canAccessCustomer(bzc, existingAccount.CustomerId) {
		bzc.Waring("无权修改该客户账号")
		return
	}
	account.CustomerId = existingAccount.CustomerId
	account.UpdateBy = bzc.GetCurrentUserName()
	customerService.UpdateAccount(account)
	bzc.Success()
}

func AccountResetPassword(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	accountId := bzc.ParamInt64("accountId")
	body := new(resetPasswordBody)
	if accountId == 0 || c.ShouldBindJSON(body) != nil {
		bzc.ParameterError()
		return
	}
	account := customerService.SelectAccountById(accountId)
	if account == nil || !canAccessCustomer(bzc, account.CustomerId) {
		bzc.Waring("无权重置该客户账号密码")
		return
	}
	customerService.ResetAccountPassword(accountId, body.Password)
	bzc.Success()
}

func AccountRemove(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	var ids slicesUtils.Slices = strings.Split(c.Param("accountIds"), ",")
	accountIds := ids.StrSlicesToInt()
	for _, accountId := range accountIds {
		account := customerService.SelectAccountById(accountId)
		if account == nil || !canAccessCustomer(bzc, account.CustomerId) {
			bzc.Waring("无权删除所选客户账号")
			return
		}
	}
	customerService.DeleteAccountByIds(accountIds)
	bzc.Success()
}

func PortalCustomerLogin(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	login := new(models.CustomerLoginBody)
	if err := c.ShouldBindJSON(login); err != nil {
		bzc.ParameterError()
		return
	}
	result, err := customerService.Login(login)
	if err != nil {
		bzc.Waring(err.Error())
		return
	}
	bzc.SuccessData(result)
}

func PortalCustomerProfile(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	value, ok := c.Get(customermiddleware.CustomerClaimsKey)
	if !ok {
		bzc.InvalidToken()
		return
	}
	claims := value.(*service.CustomerClaims)
	bzc.SuccessData(customerService.SelectPortalProfile(claims.AccountId))
}

func PortalCustomerUpdateProfile(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	value, ok := c.Get(customermiddleware.CustomerClaimsKey)
	if !ok {
		bzc.InvalidToken()
		return
	}
	claims := value.(*service.CustomerClaims)
	body := new(models.PortalProfileUpdateBody)
	if err := c.ShouldBindJSON(body); err != nil {
		bzc.ParameterError()
		return
	}
	account, err := customerService.UpdatePortalProfile(claims.AccountId, claims.Username, body)
	if err != nil {
		bzc.Waring(err.Error())
		return
	}
	bzc.SuccessData(account)
}

func PortalCustomerRouters(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	value, ok := c.Get(customermiddleware.CustomerClaimsKey)
	if !ok {
		bzc.InvalidToken()
		return
	}
	claims := value.(*service.CustomerClaims)
	bzc.SuccessData(customerService.SelectPortalRouters(claims.AccountId))
}

func PortalCustomerUpdatePassword(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	value, ok := c.Get(customermiddleware.CustomerClaimsKey)
	if !ok {
		bzc.InvalidToken()
		return
	}
	claims := value.(*service.CustomerClaims)
	body := new(models.PortalPasswordUpdateBody)
	if err := c.ShouldBindJSON(body); err != nil {
		bzc.ParameterError()
		return
	}
	if err := customerService.UpdatePortalPassword(claims.AccountId, body.OldPassword, body.NewPassword, body.ConfirmPassword); err != nil {
		bzc.Waring(err.Error())
		return
	}
	bzc.Success()
}

func PortalShipmentList(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	value, ok := c.Get(customermiddleware.CustomerClaimsKey)
	if !ok {
		bzc.InvalidToken()
		return
	}
	claims := value.(*service.CustomerClaims)
	query := new(freightModels.ShipmentPlanDQL)
	c.ShouldBind(query)
	query.CustomerId = claims.CustomerId
	query.SetLimit(c)
	list, count := shipmentService.SelectShipmentList(query)
	bzc.SuccessListData(list, count)
}

func PortalPaymentLedgerList(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	value, ok := c.Get(customermiddleware.CustomerClaimsKey)
	if !ok {
		bzc.InvalidToken()
		return
	}
	claims := value.(*service.CustomerClaims)
	query := new(freightModels.PaymentLedgerDQL)
	c.ShouldBind(query)
	query.CustomerId = claims.CustomerId
	query.SetLimit(c)
	list, total := receiptService.SelectPaymentLedger(query)
	bzc.SuccessListData(list, total)
}

func PortalShipmentStatusDict(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	bzc.SuccessData(shipmentService.SelectShipmentStatusDict())
}

func PortalShipmentDetail(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	value, ok := c.Get(customermiddleware.CustomerClaimsKey)
	if !ok {
		bzc.InvalidToken()
		return
	}
	claims := value.(*service.CustomerClaims)
	shipmentId := bzc.ParamInt64("shipmentId")
	if shipmentId == 0 {
		bzc.ParameterError()
		return
	}
	detail := shipmentService.SelectShipmentDetail(shipmentId)
	if detail == nil || detail.Plan == nil || detail.Plan.CustomerId != claims.CustomerId {
		bzc.Waring("出货计划不存在")
		return
	}
	// 客户端只展示付款摘要，不下发后台付款明细及凭证。
	detail.Payments = nil
	bzc.SuccessData(detail)
}

func PortalMenuList(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	menu := new(models.CustomerPortalMenuDQL)
	c.ShouldBind(menu)
	bzc.SuccessData(customerService.SelectPortalMenuList(menu))
}

func PortalMenuGetInfo(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	menuId := bzc.ParamInt64("menuId")
	if menuId == 0 {
		bzc.ParameterError()
		return
	}
	bzc.SuccessData(customerService.SelectPortalMenuById(menuId))
}

func PortalMenuAdd(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	menu := new(models.CustomerPortalMenuDML)
	if err := c.ShouldBindJSON(menu); err != nil {
		bzc.ParameterError()
		return
	}
	if customerService.CheckPortalMenuNameUnique(menu) {
		bzc.Waring("客户端菜单名称已存在")
		return
	}
	menu.CreateBy = bzc.GetCurrentUserName()
	menu.UpdateBy = bzc.GetCurrentUserName()
	customerService.InsertPortalMenu(menu)
	bzc.Success()
}

func PortalMenuEdit(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	menu := new(models.CustomerPortalMenuDML)
	if err := c.ShouldBindJSON(menu); err != nil {
		bzc.ParameterError()
		return
	}
	if customerService.CheckPortalMenuNameUnique(menu) {
		bzc.Waring("客户端菜单名称已存在")
		return
	}
	menu.UpdateBy = bzc.GetCurrentUserName()
	customerService.UpdatePortalMenu(menu)
	bzc.Success()
}

func PortalMenuRemove(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	menuId := bzc.ParamInt64("menuId")
	if menuId == 0 {
		bzc.ParameterError()
		return
	}
	if customerService.HasPortalMenuChildByMenuId(menuId) {
		bzc.Waring("存在子菜单，不能删除")
		return
	}
	if customerService.CheckPortalMenuExistRole(menuId) {
		bzc.Waring("菜单已分配角色，不能删除")
		return
	}
	customerService.DeletePortalMenuById(menuId)
	bzc.Success()
}

func PortalRoleList(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	role := new(models.CustomerPortalRoleDQL)
	c.ShouldBind(role)
	role.SetLimit(c)
	list, count := customerService.SelectPortalRoleList(role)
	bzc.SuccessListData(list, count)
}

func PortalRoleGetInfo(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	roleId := bzc.ParamInt64("roleId")
	if roleId == 0 {
		bzc.ParameterError()
		return
	}
	bzc.SuccessData(customerService.SelectPortalRoleById(roleId))
}

func PortalRoleAdd(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	role := new(models.CustomerPortalRoleDML)
	if err := c.ShouldBindJSON(role); err != nil {
		bzc.ParameterError()
		return
	}
	if customerService.CheckPortalRoleNameUnique(role) {
		bzc.Waring("客户端角色名称已存在")
		return
	}
	if customerService.CheckPortalRoleKeyUnique(role) {
		bzc.Waring("客户端角色权限字符已存在")
		return
	}
	role.CreateBy = bzc.GetCurrentUserName()
	role.UpdateBy = bzc.GetCurrentUserName()
	customerService.InsertPortalRole(role)
	bzc.Success()
}

func PortalRoleEdit(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	role := new(models.CustomerPortalRoleDML)
	if err := c.ShouldBindJSON(role); err != nil {
		bzc.ParameterError()
		return
	}
	if customerService.CheckPortalRoleNameUnique(role) {
		bzc.Waring("客户端角色名称已存在")
		return
	}
	if customerService.CheckPortalRoleKeyUnique(role) {
		bzc.Waring("客户端角色权限字符已存在")
		return
	}
	role.UpdateBy = bzc.GetCurrentUserName()
	customerService.UpdatePortalRole(role)
	bzc.Success()
}

func PortalRoleChangeStatus(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	role := new(models.CustomerPortalRoleDML)
	if err := c.ShouldBindJSON(role); err != nil || role.RoleId == 0 {
		bzc.ParameterError()
		return
	}
	customerService.UpdatePortalRoleStatus(role.RoleId, role.Status, bzc.GetCurrentUserName())
	bzc.Success()
}

func PortalRoleRemove(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	roleIds := bzc.ParamInt64Array("roleIds")
	if customerService.CountAccountRoleByRoleIds(roleIds) {
		bzc.Waring("角色已分配客户账号，不能删除")
		return
	}
	customerService.DeletePortalRoleByIds(roleIds, bzc.GetCurrentUserName())
	bzc.Success()
}

func PortalRoleMenuTreeselect(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	roleId := bzc.ParamInt64("roleId")
	if roleId == 0 {
		bzc.ParameterError()
		return
	}
	bzc.SuccessData(customerService.SelectPortalRoleMenuTreeselect(roleId))
}

func PortalRoleOptions(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	bzc.SuccessData(customerService.SelectPortalRoleOptions())
}

func AccountRoleEdit(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	accountId := bzc.ParamInt64("accountId")
	body := new(accountRoleBody)
	if accountId == 0 || c.ShouldBindJSON(body) != nil {
		bzc.ParameterError()
		return
	}
	account := customerService.SelectAccountById(accountId)
	if account == nil || !canAccessCustomer(bzc, account.CustomerId) {
		bzc.Waring("无权配置该客户账号角色")
		return
	}
	customerService.UpdateAccountRoles(accountId, body.RoleIds)
	bzc.Success()
}

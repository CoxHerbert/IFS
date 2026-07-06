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

type resetPasswordBody struct {
	Password string `json:"password" binding:"required"`
}

type accountRoleBody struct {
	RoleIds []string `json:"roleIds"`
}

type portalShipmentCreateBody struct {
	OrderNo       string                          `json:"orderNo"`
	Pol           string                          `json:"pol"`
	Pod           string                          `json:"pod"`
	PlannedEtd    string                          `json:"plannedEtd"`
	PlannedEta    string                          `json:"plannedEta"`
	Remark        string                          `json:"remark"`
	PreferredType string                          `json:"preferredType"`
	CargoList     []*freightModels.CargoImportReq `json:"cargoList" binding:"required"`
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
	bzc.SuccessData(customerService.SelectCustomerById(customerId))
}

func CustomerOptions(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	bzc.SuccessData(customerService.SelectCustomerOptions(c.Query("keyword")))
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
	customer.UpdateBy = bzc.GetCurrentUserName()
	customerService.UpdateCustomer(customer)
	bzc.Success()
}

func CustomerRemove(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	var ids slicesUtils.Slices = strings.Split(c.Param("customerIds"), ",")
	customerService.DeleteCustomerByIds(ids.StrSlicesToInt())
	bzc.Success()
}

func ContactList(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	contact := new(models.CustomerContactDQL)
	c.ShouldBind(contact)
	contact.CustomerId = bzc.ParamInt64("customerId")
	contact.SetLimit(c)
	list, count := customerService.SelectContactList(contact)
	bzc.SuccessListData(list, count)
}

func ContactGetInfo(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	contactId := bzc.ParamInt64("contactId")
	if contactId == 0 {
		bzc.ParameterError()
		return
	}
	bzc.SuccessData(customerService.SelectContactById(contactId))
}

func ContactAdd(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	contact := new(models.CustomerContactDML)
	if err := c.ShouldBindJSON(contact); err != nil {
		zap.L().Error("参数错误", zap.Error(err))
		bzc.ParameterError()
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
	contact.UpdateBy = bzc.GetCurrentUserName()
	customerService.UpdateContact(contact)
	bzc.Success()
}

func ContactRemove(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	var ids slicesUtils.Slices = strings.Split(c.Param("contactIds"), ",")
	customerService.DeleteContactByIds(ids.StrSlicesToInt())
	bzc.Success()
}

func AccountList(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	account := new(models.CustomerAccountDQL)
	c.ShouldBind(account)
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
	data := map[string]interface{}{
		"account": customerService.SelectAccountById(accountId),
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
	customerService.ResetAccountPassword(accountId, body.Password)
	bzc.Success()
}

func AccountRemove(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	var ids slicesUtils.Slices = strings.Split(c.Param("accountIds"), ",")
	customerService.DeleteAccountByIds(ids.StrSlicesToInt())
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

func PortalShipmentAssistantEstimate(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	req := new(freightModels.ShipmentEstimateReq)
	if err := c.ShouldBindJSON(req); err != nil {
		bzc.ParameterError()
		return
	}
	data, err := shipmentService.EstimateShipment(req)
	if err != nil {
		bzc.Waring(err.Error())
		return
	}
	bzc.SuccessData(data)
}

func PortalShipmentCreateFromAssistant(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	value, ok := c.Get(customermiddleware.CustomerClaimsKey)
	if !ok {
		bzc.InvalidToken()
		return
	}
	claims := value.(*service.CustomerClaims)
	req := new(portalShipmentCreateBody)
	if err := c.ShouldBindJSON(req); err != nil {
		bzc.ParameterError()
		return
	}
	account := customerService.SelectAccountProfile(claims.AccountId)
	customerName := ""
	if account != nil {
		customerName = account.CustomerName
	}
	detail, err := shipmentService.ImportShipment(&freightModels.ShipmentImportReq{
		CustomerId:    claims.CustomerId,
		CustomerName:  customerName,
		OrderNo:       req.OrderNo,
		Pol:           req.Pol,
		Pod:           req.Pod,
		PlannedEtd:    req.PlannedEtd,
		PlannedEta:    req.PlannedEta,
		Remark:        req.Remark,
		PreferredType: req.PreferredType,
		CargoList:     req.CargoList,
	}, claims.Username)
	if err != nil {
		bzc.Waring(err.Error())
		return
	}
	bzc.SuccessData(detail)
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
	customerService.UpdateAccountRoles(accountId, body.RoleIds)
	bzc.Success()
}

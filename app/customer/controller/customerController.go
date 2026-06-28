package controller

import (
	"baize/app/common/baize/baizeContext"
	customermiddleware "baize/app/customer/middleware"
	"baize/app/customer/models"
	"baize/app/customer/service"
	"baize/app/utils/slicesUtils"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var customerService = service.GetCustomerService()

type resetPasswordBody struct {
	Password string `json:"password" binding:"required"`
}

func CustomerList(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	customer := new(models.CustomerDQL)
	c.ShouldBind(customer)
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
	bzc.SuccessData(customerService.SelectAccountById(accountId))
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
	bzc.SuccessData(customerService.SelectAccountProfile(claims.AccountId))
}

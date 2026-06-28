package controller

import (
	"baize/app/common/baize/baizeContext"
	"baize/app/portal/models"
	"baize/app/utils/slicesUtils"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ContactList(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	contact := new(models.ContactDQL)
	c.ShouldBind(contact)
	contact.SetLimit(c)
	list, count := contactService.SelectContactList(contact)
	bzc.SuccessListData(list, count)
}

func ContactExport(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	contact := new(models.ContactDQL)
	c.ShouldBind(contact)
	data := contactService.ExportContact(contact)
	bzc.DataPackageExcel(data)
}

func ContactGetInfo(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	contactId := bzc.ParamInt64("contactId")
	if contactId == 0 {
		zap.L().Error("参数错误")
		bzc.ParameterError()
		return
	}
	bzc.SuccessData(contactService.SelectContactById(contactId))
}

func ContactEdit(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	bzc.SetLog("官网线索", "UPDATE")
	contact := new(models.ContactDML)
	if err := c.ShouldBindJSON(contact); err != nil {
		zap.L().Error("参数错误", zap.Error(err))
		bzc.ParameterError()
		return
	}
	contact.UpdateBy = bzc.GetCurrentUserName()
	contactService.UpdateContact(contact)
	bzc.Success()
}

func ContactRemove(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	bzc.SetLog("官网线索", "DELETE")
	var ids slicesUtils.Slices = strings.Split(c.Param("contactIds"), ",")
	contactService.DeleteContactByIds(ids.StrSlicesToInt())
	bzc.Success()
}

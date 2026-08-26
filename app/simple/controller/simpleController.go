package controller

import (
	"baize/app/common/baize/baizeContext"
	customerService "baize/app/customer/service"
	freightService "baize/app/freight/service"

	"github.com/gin-gonic/gin"
)

var simpleCustomerService = customerService.GetCustomerService()

// SimpleCustomerOptions provides lightweight customer data for selects and autocomplete inputs.
func SimpleCustomerOptions(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	if freightService.CanManageAllShipments(bzc.GetCurrentUser()) {
		bzc.SuccessData(simpleCustomerService.SelectCustomerOptions(c.Query("keyword")))
		return
	}
	bzc.SuccessData(simpleCustomerService.SelectCustomerOptionsBySales(c.Query("keyword"), bzc.GetCurrentUserId()))
}

package controller

import (
	"baize/app/common/baize/baizeContext"
	"baize/app/constant/constants"
	customerServicePkg "baize/app/customer/service"
	"baize/app/freight/models"
	"baize/app/freight/service"
	"baize/app/utils/fileUploadUtils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strconv"
	"strings"
)

var receiptService = service.GetReceiptService()
var receiptCustomerService = customerServicePkg.GetCustomerService()

func ReceiptList(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	q := new(models.ReceiptDQL)
	c.ShouldBind(q)
	if !service.CanManageAllShipments(bzc.GetCurrentUser()) {
		q.SalesUserId = bzc.GetCurrentUserId()
	}
	q.SetLimit(c)
	list, total := receiptService.SelectList(q)
	bzc.SuccessListData(list, total)
}
func ReceiptGet(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id := bzc.ParamInt64("receiptId")
	vo := receiptService.SelectById(id)
	if vo == nil {
		bzc.Waring("收款单不存在")
		return
	}
	if !canAccessReceiptCustomer(bzc, vo.CustomerId) {
		bzc.Waring("无权查看该客户的收款记录")
		return
	}
	bzc.SuccessData(vo)
}

func ReceiptCreate(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	req := new(models.ReceiptCreateReq)
	amount, err := strconv.ParseFloat(strings.TrimSpace(c.PostForm("amount")), 64)
	if err != nil {
		bzc.ParameterError()
		return
	}
	req.Amount = amount
	req.Currency = c.PostForm("currency")
	req.ReceiptTime = c.PostForm("receiptTime")
	req.PaymentMethod = c.PostForm("paymentMethod")
	req.CustomerName = c.PostForm("customerName")
	req.Remark = c.PostForm("remark")
	req.CustomerId, err = strconv.ParseInt(c.PostForm("customerId"), 10, 64)
	if err != nil {
		bzc.ParameterError()
		return
	}
	if !canAccessReceiptCustomer(bzc, req.CustomerId) {
		bzc.Waring("无权为该客户登记收款")
		return
	}
	if raw := c.PostForm("allocations"); raw != "" {
		if err = json.Unmarshal([]byte(raw), &req.Allocations); err != nil {
			bzc.ParameterError()
			return
		}
	}
	voucherUrl, voucherName := "", ""
	if file, fileErr := c.FormFile("voucher"); fileErr == nil {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if file.Size > maxPaymentVoucherSize || (ext != ".pdf" && ext != ".png" && ext != ".jpg" && ext != ".jpeg") {
			bzc.Waring("收款凭证仅支持 PDF、PNG、JPG，且不能超过10MB")
			return
		}
		voucherName = filepath.Base(file.Filename)
		voucherUrl = constants.ResourcePrefix + fileUploadUtils.Upload(constants.PaymentVoucherPath, file)
	}
	vo, err := receiptService.Create(req, voucherUrl, voucherName, bzc.GetCurrentUserName())
	if err != nil {
		bzc.Waring(err.Error())
		return
	}
	bzc.SuccessData(vo)
}

func ReceiptRemove(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id := bzc.ParamInt64("receiptId")
	vo := receiptService.SelectById(id)
	if vo == nil || !canAccessReceiptCustomer(bzc, vo.CustomerId) {
		bzc.Waring("无权删除该收款记录")
		return
	}
	if err := receiptService.Delete(id); err != nil {
		bzc.Waring(err.Error())
		return
	}
	bzc.Success()
}

func PaymentLedgerList(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	q := new(models.PaymentLedgerDQL)
	c.ShouldBind(q)
	if !service.CanManageAllShipments(bzc.GetCurrentUser()) {
		q.SalesUserId = bzc.GetCurrentUserId()
	}
	q.SetLimit(c)
	list, total := receiptService.SelectPaymentLedger(q)
	bzc.SuccessListData(list, total)
}
func PaymentLedgerUpdatePayable(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id := bzc.ParamInt64("shipmentId")
	req := new(models.PayableAmountUpdateReq)
	if id == 0 || c.ShouldBindJSON(req) != nil {
		bzc.ParameterError()
		return
	}
	if !shipmentService.CanOperateShipment(id, bzc.GetCurrentUserId(), service.CanManageAllShipments(bzc.GetCurrentUser())) {
		bzc.Waring("无权维护该客户的应收金额")
		return
	}
	if err := receiptService.UpdatePayableAmount(id, req.PayableAmount); err != nil {
		bzc.Waring(err.Error())
		return
	}
	bzc.Success()
}
func ShipmentChargeList(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id := bzc.ParamInt64("shipmentId")
	if id == 0 {
		bzc.ParameterError()
		return
	}
	if !shipmentService.CanOperateShipment(id, bzc.GetCurrentUserId(), service.CanManageAllShipments(bzc.GetCurrentUser())) {
		bzc.Waring("无权查看该客户的费用明细")
		return
	}
	bzc.SuccessData(receiptService.SelectShipmentCharges(id))
}
func ShipmentChargeReplace(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id := bzc.ParamInt64("shipmentId")
	req := new(models.ShipmentChargeUpdateReq)
	if id == 0 || c.ShouldBindJSON(req) != nil {
		bzc.ParameterError()
		return
	}
	if !shipmentService.CanOperateShipment(id, bzc.GetCurrentUserId(), service.CanManageAllShipments(bzc.GetCurrentUser())) {
		bzc.Waring("无权维护该客户的费用明细")
		return
	}
	if err := receiptService.ReplaceShipmentCharges(id, req.Items, bzc.GetCurrentUserName()); err != nil {
		bzc.Waring(err.Error())
		return
	}
	bzc.Success()
}

func canAccessReceiptCustomer(bzc *baizeContext.BaiZeContext, customerId int64) bool {
	if service.CanManageAllShipments(bzc.GetCurrentUser()) {
		return true
	}
	customer := receiptCustomerService.SelectCustomerById(customerId)
	return customer != nil && customer.SalesUserId == bzc.GetCurrentUserId()
}

func PaymentDeclarationList(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	q := new(models.PaymentDeclarationDQL)
	c.ShouldBind(q)
	if !service.CanManageAllShipments(bzc.GetCurrentUser()) {
		q.SalesUserId = bzc.GetCurrentUserId()
	}
	q.SetLimit(c)
	list, total := receiptService.SelectDeclarationList(q)
	bzc.SuccessListData(list, total)
}
func PaymentDeclarationGet(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	vo := receiptService.SelectDeclarationById(bzc.ParamInt64("declarationId"))
	if vo == nil {
		bzc.Waring("付款申报不存在")
		return
	}
	if !shipmentService.CanOperateShipment(vo.ShipmentId, bzc.GetCurrentUserId(), service.CanManageAllShipments(bzc.GetCurrentUser())) {
		bzc.Waring("无权查看该客户的付款申报")
		return
	}
	bzc.SuccessData(vo)
}
func PaymentDeclarationApprove(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id := bzc.ParamInt64("declarationId")
	item := receiptService.SelectDeclarationById(id)
	if item == nil || !shipmentService.CanOperateShipment(item.ShipmentId, bzc.GetCurrentUserId(), service.CanManageAllShipments(bzc.GetCurrentUser())) {
		bzc.Waring("无权审核该客户的付款申报")
		return
	}
	req := new(models.PaymentDeclarationReviewReq)
	_ = c.ShouldBindJSON(req)
	vo, err := receiptService.ApproveDeclaration(id, bzc.GetCurrentUserName(), req.Remark)
	if err != nil {
		bzc.Waring(err.Error())
		return
	}
	bzc.SuccessData(vo)
}
func PaymentDeclarationReject(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id := bzc.ParamInt64("declarationId")
	item := receiptService.SelectDeclarationById(id)
	if item == nil || !shipmentService.CanOperateShipment(item.ShipmentId, bzc.GetCurrentUserId(), service.CanManageAllShipments(bzc.GetCurrentUser())) {
		bzc.Waring("无权审核该客户的付款申报")
		return
	}
	req := new(models.PaymentDeclarationReviewReq)
	if c.ShouldBindJSON(req) != nil {
		bzc.ParameterError()
		return
	}
	if err := receiptService.RejectDeclaration(id, bzc.GetCurrentUserName(), req.Remark); err != nil {
		bzc.Waring(err.Error())
		return
	}
	bzc.Success()
}

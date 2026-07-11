package controller

import (
	"baize/app/common/baize/baizeContext"
	"baize/app/constant/constants"
	"baize/app/freight/models"
	"baize/app/freight/service"
	"baize/app/utils/fileUploadUtils"
	"baize/app/utils/slicesUtils"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const maxPaymentVoucherSize = 10 << 20

var shipmentService = service.GetShipmentService()

func ShipmentImport(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	req := new(models.ShipmentImportReq)
	if err := c.ShouldBindJSON(req); err != nil {
		zap.L().Error("shipment import param error", zap.Error(err))
		bzc.ParameterError()
		return
	}
	detail, err := shipmentService.ImportShipment(req, bzc.GetCurrentUserName(), bzc.GetCurrentUserId())
	if err != nil {
		bzc.Waring(err.Error())
		return
	}
	bzc.SuccessData(detail)
}

func ShipmentList(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	query := new(models.ShipmentPlanDQL)
	c.ShouldBind(query)
	if !service.CanManageAllShipments(bzc.GetCurrentUser()) {
		query.SalesUserId = bzc.GetCurrentUserId()
	}
	query.SetLimit(c)
	list, count := shipmentService.SelectShipmentList(query)
	bzc.SuccessListData(list, count)
}

func ShipmentGetInfo(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	shipmentId := bzc.ParamInt64("shipmentId")
	if shipmentId == 0 {
		bzc.ParameterError()
		return
	}
	if !shipmentService.CanOperateShipment(shipmentId, bzc.GetCurrentUserId(), service.CanManageAllShipments(bzc.GetCurrentUser())) {
		bzc.Waring("无权查看该客户的出货计划")
		return
	}
	detail := shipmentService.SelectShipmentDetail(shipmentId)
	if detail == nil {
		bzc.Waring("出货计划不存在")
		return
	}
	bzc.SuccessData(detail)
}

func ShipmentUpdateStatus(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	shipmentId := bzc.ParamInt64("shipmentId")
	req := new(models.ShipmentStatusUpdateReq)
	if shipmentId == 0 || c.ShouldBindJSON(req) != nil {
		bzc.ParameterError()
		return
	}
	if err := shipmentService.UpdateShipmentStatus(shipmentId, req, bzc.GetCurrentUserName(), bzc.GetCurrentUserId(), service.CanManageAllShipments(bzc.GetCurrentUser())); err != nil {
		bzc.Waring(err.Error())
		return
	}
	bzc.Success()
}

func ShipmentBindCustomer(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	shipmentId := bzc.ParamInt64("shipmentId")
	req := new(models.ShipmentCustomerBindReq)
	if shipmentId == 0 || c.ShouldBindJSON(req) != nil {
		bzc.ParameterError()
		return
	}
	if !shipmentService.CanOperateShipment(shipmentId, bzc.GetCurrentUserId(), service.CanManageAllShipments(bzc.GetCurrentUser())) {
		bzc.Waring("无权维护该客户的出货计划")
		return
	}
	if err := shipmentService.UpdateShipmentCustomer(shipmentId, req, bzc.GetCurrentUserName()); err != nil {
		bzc.Waring(err.Error())
		return
	}
	bzc.Success()
}

func ShipmentConfirm(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	shipmentId := bzc.ParamInt64("shipmentId")
	if shipmentId == 0 {
		bzc.ParameterError()
		return
	}
	if !shipmentService.CanOperateShipment(shipmentId, bzc.GetCurrentUserId(), service.CanManageAllShipments(bzc.GetCurrentUser())) {
		bzc.Waring("无权维护该客户的出货计划")
		return
	}
	order, err := shipmentService.ConfirmShipment(shipmentId, bzc.GetCurrentUserName())
	if err != nil {
		bzc.Waring(err.Error())
		return
	}
	bzc.SuccessData(order)
}

func ShipmentShareInfo(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	shipmentId := bzc.ParamInt64("shipmentId")
	if !shipmentService.CanOperateShipment(shipmentId, bzc.GetCurrentUserId(), service.CanManageAllShipments(bzc.GetCurrentUser())) {
		bzc.Waring("无权查看该客户的出货计划")
		return
	}
	detail := shipmentService.SelectShipmentDetail(shipmentId)
	if shipmentId == 0 || detail == nil || detail.Plan == nil {
		bzc.ParameterError()
		return
	}
	token := detail.Plan.ShareToken
	bzc.SuccessData(models.ShareInfoVo{Token: token, ShareUrl: "/shipment/share/" + token})
}

func ShipmentRemove(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	var ids slicesUtils.Slices = strings.Split(c.Param("shipmentIds"), ",")
	shipmentIds := ids.StrSlicesToInt()
	canManageAll := service.CanManageAllShipments(bzc.GetCurrentUser())
	for _, shipmentId := range shipmentIds {
		if !shipmentService.CanOperateShipment(shipmentId, bzc.GetCurrentUserId(), canManageAll) {
			bzc.Waring("无权删除该客户的出货计划")
			return
		}
	}
	shipmentService.DeleteShipmentByIds(shipmentIds)
	bzc.Success()
}

func ShipmentPaymentAdd(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	shipmentId := bzc.ParamInt64("shipmentId")
	amount, err := strconv.ParseFloat(strings.TrimSpace(c.PostForm("amount")), 64)
	if shipmentId == 0 || err != nil || amount <= 0 { bzc.ParameterError(); return }
	canManageAll := service.CanManageAllShipments(bzc.GetCurrentUser())
	if !shipmentService.CanOperateShipment(shipmentId, bzc.GetCurrentUserId(), canManageAll) {
		bzc.Waring("无权维护该客户的付款记录"); return
	}
	payment := &models.ShipmentPaymentDML{
		ShipmentId: shipmentId, Amount: amount, Currency: c.PostForm("currency"), PaymentTime: c.PostForm("paymentTime"),
		PaymentMethod: c.PostForm("paymentMethod"), Remark: c.PostForm("remark"), CreateBy: bzc.GetCurrentUserName(),
	}
	file, fileErr := c.FormFile("voucher")
	if fileErr == nil {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if file.Size > maxPaymentVoucherSize || (ext != ".pdf" && ext != ".png" && ext != ".jpg" && ext != ".jpeg") {
			bzc.Waring("付款凭证仅支持 PDF、PNG、JPG，且不能超过10MB"); return
		}
		payment.VoucherName = filepath.Base(file.Filename)
		payment.VoucherUrl = constants.ResourcePrefix + fileUploadUtils.Upload(constants.PaymentVoucherPath, file)
	}
	if err := shipmentService.AddPayment(payment, bzc.GetCurrentUserId(), canManageAll); err != nil {
		bzc.Waring(err.Error()); return
	}
	bzc.SuccessData(payment)
}

func ShipmentPaymentRemove(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	shipmentId := bzc.ParamInt64("shipmentId")
	paymentId, err := strconv.ParseInt(c.Param("paymentId"), 10, 64)
	if shipmentId == 0 || err != nil || paymentId == 0 { bzc.ParameterError(); return }
	if err = shipmentService.DeletePayment(shipmentId, paymentId, bzc.GetCurrentUserName(), bzc.GetCurrentUserId(), service.CanManageAllShipments(bzc.GetCurrentUser())); err != nil {
		bzc.Waring(err.Error()); return
	}
	bzc.Success()
}

func PortalShipmentShare(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	token := strings.TrimSpace(c.Param("token"))
	if token == "" {
		bzc.ParameterError()
		return
	}
	detail := shipmentService.SelectShareDetail(token)
	if detail == nil {
		bzc.Waring("分享链接无效或已过期")
		return
	}
	bzc.SuccessData(detail)
}

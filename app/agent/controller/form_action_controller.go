package controller

import (
	"baize/app/agent/protocol"
	"baize/app/agent/request"
	"baize/app/agent/service"
	"baize/app/common/baize/baizeContext"
	"baize/app/constant/constants"
	"baize/app/utils/fileUploadUtils"
	customermiddleware "baize/app/customer/middleware"
	customerService "baize/app/customer/service"
	"encoding/json"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var formActionService = service.GetFormActionService()

func SubmitAgentForm(c *gin.Context) {
	req := new(request.FormSubmitRequest)
	if strings.HasPrefix(c.ContentType(), "multipart/form-data") {
		req.SessionID, _ = strconv.ParseInt(c.PostForm("sessionId"), 10, 64)
		req.FormCode = c.PostForm("formCode")
		if err := json.Unmarshal([]byte(c.PostForm("values")), &req.Values); err != nil { c.JSON(400, protocol.NewErrorResult("invalid form values")); return }
		if file, err := c.FormFile("voucher"); err == nil {
			ext := strings.ToLower(filepath.Ext(file.Filename))
			if file.Size > 10<<20 || (ext != ".pdf" && ext != ".png" && ext != ".jpg" && ext != ".jpeg") { c.JSON(400, protocol.NewErrorResult("付款凭证仅支持 PDF、PNG、JPG，且不能超过10MB")); return }
			req.VoucherName = filepath.Base(file.Filename)
			req.VoucherURL = constants.ResourcePrefix + fileUploadUtils.Upload(constants.PaymentVoucherPath, file)
		}
	} else if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, protocol.NewErrorResult("invalid form submit payload"))
		return
	}
	fillSubmitContext(c, req)
	c.JSON(200, formActionService.SubmitForm(req))
}

func ExecuteAgentAction(c *gin.Context) {
	req := new(request.ActionExecuteRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, protocol.NewErrorResult("invalid action payload"))
		return
	}
	c.JSON(200, formActionService.ExecuteAction(req))
}

func fillSubmitContext(c *gin.Context, req *request.FormSubmitRequest) {
	if claims := customerClaimsFromContext(c); claims != nil {
		req.Source = "customer"
		req.CustomerID = claims.CustomerId
		req.OperatorName = claims.Username
		if customer := customerService.GetCustomerService().SelectCustomerById(claims.CustomerId); customer != nil {
			req.CustomerName = customer.CustomerName
		}
		return
	}

	if claims := customerClaimsFromAuthorization(c); claims != nil {
		req.Source = "customer"
		req.CustomerID = claims.CustomerId
		req.OperatorName = claims.Username
		if customer := customerService.GetCustomerService().SelectCustomerById(claims.CustomerId); customer != nil {
			req.CustomerName = customer.CustomerName
		}
		return
	}

	req.Source = "admin"
	bzc := baizeContext.NewBaiZeContext(c)
	if bzc.GetCurrentLoginUser() != nil {
		req.OperatorID = bzc.GetCurrentUserId()
		req.OperatorName = bzc.GetCurrentUserName()
		req.Permissions = currentPermissions(bzc)
	} else {
		req.Source = "public"
	}
	if req.OperatorName == "" {
		req.OperatorName = "agent"
	}
}

func customerClaimsFromContext(c *gin.Context) *customerService.CustomerClaims {
	value, ok := c.Get(customermiddleware.CustomerClaimsKey)
	if !ok {
		return nil
	}
	claims, ok := value.(*customerService.CustomerClaims)
	if !ok {
		return nil
	}
	return claims
}

func customerClaimsFromAuthorization(c *gin.Context) *customerService.CustomerClaims {
	authHeader := c.Request.Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != constants.TokenPrefix {
		return nil
	}
	claims, err := customerService.ParseCustomerToken(parts[1])
	if err != nil {
		return nil
	}
	return claims
}

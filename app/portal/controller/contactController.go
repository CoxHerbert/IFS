package controller

import (
	"baize/app/common/baize/baizeContext"
	"baize/app/portal/models"
	"baize/app/portal/service"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var contactService = service.GetContactService()

type ContactSubmitReq struct {
	ContactName string `json:"contactName"`
	CompanyName string `json:"companyName"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Route       string `json:"route"`
	CargoInfo   string `json:"cargoInfo"`
	Message     string `json:"message"`
	Source      string `json:"source"`
}

type ContactSubmitResp struct {
	LeadNo string `json:"leadNo"`
}

func SubmitContact(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	var req ContactSubmitReq
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("官网联系我们参数解析失败", zap.Error(err))
		bzc.ParameterError()
		return
	}

	req.ContactName = strings.TrimSpace(req.ContactName)
	req.CompanyName = strings.TrimSpace(req.CompanyName)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Email = strings.TrimSpace(req.Email)
	req.Route = strings.TrimSpace(req.Route)
	req.CargoInfo = strings.TrimSpace(req.CargoInfo)
	req.Message = strings.TrimSpace(req.Message)
	req.Source = strings.TrimSpace(req.Source)

	if req.ContactName == "" || (req.Phone == "" && req.Email == "") || req.Message == "" {
		bzc.Waring("请填写联系人、电话或邮箱，以及需求说明")
		return
	}

	contact := &models.ContactDML{
		ContactName: req.ContactName,
		CompanyName: req.CompanyName,
		Phone:       req.Phone,
		Email:       req.Email,
		Route:       req.Route,
		CargoInfo:   req.CargoInfo,
		Message:     req.Message,
		Source:      req.Source,
		IpAddr:      c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
	}
	leadNo := contactService.InsertContact(contact)
	zap.L().Info("官网联系我们提交",
		zap.String("leadNo", leadNo),
		zap.String("contactName", req.ContactName),
		zap.String("companyName", req.CompanyName),
		zap.String("phone", req.Phone),
		zap.String("email", req.Email),
		zap.String("route", req.Route),
		zap.String("cargoInfo", req.CargoInfo),
		zap.String("message", req.Message),
		zap.String("source", req.Source),
		zap.String("ip", c.ClientIP()),
	)

	bzc.SuccessData(ContactSubmitResp{LeadNo: leadNo})
}

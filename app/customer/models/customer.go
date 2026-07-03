package models

import (
	"baize/app/common/baize/baizeUnix"
	"baize/app/common/commonModels"
)

type CustomerDQL struct {
	CustomerName string `form:"customerName" db:"customer_name"`
	CompanyName  string `form:"companyName" db:"company_name"`
	ContactName  string `form:"contactName" db:"contact_name"`
	Phone        string `form:"phone" db:"phone"`
	Email        string `form:"email" db:"email"`
	Status       string `form:"status" db:"status"`
	BeginTime    string `form:"beginTime" db:"begin_time"`
	EndTime      string `form:"endTime" db:"end_time"`
	commonModels.BaseEntityDQL
}

type CustomerDML struct {
	CustomerId   int64  `json:"customerId,string" db:"customer_id"`
	CustomerNo   string `json:"customerNo" db:"customer_no"`
	CustomerName string `json:"customerName" db:"customer_name"`
	CompanyName  string `json:"companyName" db:"company_name"`
	ContactName  string `json:"contactName" db:"contact_name"`
	Phone        string `json:"phone" db:"phone"`
	Email        string `json:"email" db:"email"`
	Status       string `json:"status" db:"status"`
	Remark       string `json:"remark" db:"remark"`
	CreateBy     string `json:"createBy" db:"create_by"`
	UpdateBy     string `json:"updateBy" db:"update_by"`
}

type CustomerVo struct {
	CustomerId   int64                `json:"customerId,string" db:"customer_id"`
	CustomerNo   string               `json:"customerNo" db:"customer_no"`
	CustomerName string               `json:"customerName" db:"customer_name"`
	CompanyName  string               `json:"companyName" db:"company_name"`
	ContactName  string               `json:"contactName" db:"contact_name"`
	Phone        string               `json:"phone" db:"phone"`
	Email        string               `json:"email" db:"email"`
	Status       string               `json:"status" db:"status"`
	Remark       string               `json:"remark" db:"remark"`
	CreateBy     string               `json:"createBy" db:"create_by"`
	CreateTime   *baizeUnix.BaiZeTime `json:"createTime" db:"create_time"`
	UpdateBy     string               `json:"updateBy" db:"update_by"`
	UpdateTime   *baizeUnix.BaiZeTime `json:"updateTime" db:"update_time"`
}

type CustomerOptionVo struct {
	CustomerId   int64  `json:"customerId,string" db:"customer_id"`
	CustomerNo   string `json:"customerNo" db:"customer_no"`
	CustomerName string `json:"customerName" db:"customer_name"`
	CompanyName  string `json:"companyName" db:"company_name"`
}

type CustomerContactDQL struct {
	CustomerId  int64  `form:"customerId" db:"customer_id"`
	ContactName string `form:"contactName" db:"contact_name"`
	Phone       string `form:"phone" db:"phone"`
	Status      string `form:"status" db:"status"`
	commonModels.BaseEntityDQL
}

type CustomerContactDML struct {
	ContactId   int64  `json:"contactId,string" db:"contact_id"`
	CustomerId  int64  `json:"customerId,string" db:"customer_id"`
	ContactName string `json:"contactName" db:"contact_name"`
	Position    string `json:"position" db:"position"`
	Phone       string `json:"phone" db:"phone"`
	Email       string `json:"email" db:"email"`
	Wechat      string `json:"wechat" db:"wechat"`
	IsPrimary   string `json:"isPrimary" db:"is_primary"`
	Status      string `json:"status" db:"status"`
	Remark      string `json:"remark" db:"remark"`
	CreateBy    string `json:"createBy" db:"create_by"`
	UpdateBy    string `json:"updateBy" db:"update_by"`
}

type CustomerContactVo struct {
	ContactId   int64                `json:"contactId,string" db:"contact_id"`
	CustomerId  int64                `json:"customerId,string" db:"customer_id"`
	ContactName string               `json:"contactName" db:"contact_name"`
	Position    string               `json:"position" db:"position"`
	Phone       string               `json:"phone" db:"phone"`
	Email       string               `json:"email" db:"email"`
	Wechat      string               `json:"wechat" db:"wechat"`
	IsPrimary   string               `json:"isPrimary" db:"is_primary"`
	Status      string               `json:"status" db:"status"`
	Remark      string               `json:"remark" db:"remark"`
	CreateBy    string               `json:"createBy" db:"create_by"`
	CreateTime  *baizeUnix.BaiZeTime `json:"createTime" db:"create_time"`
	UpdateBy    string               `json:"updateBy" db:"update_by"`
	UpdateTime  *baizeUnix.BaiZeTime `json:"updateTime" db:"update_time"`
}

type CustomerAccountDQL struct {
	CustomerId int64  `form:"customerId" db:"customer_id"`
	Username   string `form:"username" db:"username"`
	RealName   string `form:"realName" db:"real_name"`
	Phone      string `form:"phone" db:"phone"`
	Status     string `form:"status" db:"status"`
	BeginTime  string `form:"beginTime" db:"begin_time"`
	EndTime    string `form:"endTime" db:"end_time"`
	commonModels.BaseEntityDQL
}

type CustomerAccountDML struct {
	AccountId  int64  `json:"accountId,string" db:"account_id"`
	CustomerId int64  `json:"customerId,string" db:"customer_id"`
	Username   string `json:"username" db:"username"`
	Password   string `json:"password,omitempty" db:"password"`
	RealName   string `json:"realName" db:"real_name"`
	Phone      string `json:"phone" db:"phone"`
	Email      string `json:"email" db:"email"`
	IsMain     string `json:"isMain" db:"is_main"`
	Status     string `json:"status" db:"status"`
	Remark     string `json:"remark" db:"remark"`
	CreateBy   string `json:"createBy" db:"create_by"`
	UpdateBy   string `json:"updateBy" db:"update_by"`
}

type CustomerAccountVo struct {
	AccountId     int64                `json:"accountId,string" db:"account_id"`
	CustomerId    int64                `json:"customerId,string" db:"customer_id"`
	CustomerNo    string               `json:"customerNo" db:"customer_no"`
	CustomerName  string               `json:"customerName" db:"customer_name"`
	CompanyName   string               `json:"companyName" db:"company_name"`
	Username      string               `json:"username" db:"username"`
	RealName      string               `json:"realName" db:"real_name"`
	Phone         string               `json:"phone" db:"phone"`
	Email         string               `json:"email" db:"email"`
	IsMain        string               `json:"isMain" db:"is_main"`
	Status        string               `json:"status" db:"status"`
	RoleNames     string               `json:"roleNames" db:"role_names"`
	LastLoginTime *baizeUnix.BaiZeTime `json:"lastLoginTime" db:"last_login_time"`
	Remark        string               `json:"remark" db:"remark"`
	CreateBy      string               `json:"createBy" db:"create_by"`
	CreateTime    *baizeUnix.BaiZeTime `json:"createTime" db:"create_time"`
	UpdateBy      string               `json:"updateBy" db:"update_by"`
	UpdateTime    *baizeUnix.BaiZeTime `json:"updateTime" db:"update_time"`
}

type CustomerLoginBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Uuid     string `json:"uuid" binding:"required"`
}

type CustomerLoginResult struct {
	Token string             `json:"token"`
	User  *CustomerAccountVo `json:"user"`
}

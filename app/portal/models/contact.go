package models

import (
	"baize/app/common/baize/baizeUnix"
	"baize/app/common/commonModels"
)

type ContactDQL struct {
	ContactName string `form:"contactName" db:"contact_name"`
	CompanyName string `form:"companyName" db:"company_name"`
	Phone       string `form:"phone" db:"phone"`
	Email       string `form:"email" db:"email"`
	Route       string `form:"route" db:"route"`
	Source      string `form:"source" db:"source"`
	Status      string `form:"status" db:"status"`
	BeginTime   string `form:"beginTime" db:"begin_time"`
	EndTime     string `form:"endTime" db:"end_time"`
	commonModels.BaseEntityDQL
}

type ContactDML struct {
	ContactId   int64  `json:"contactId,string" db:"contact_id"`
	LeadNo      string `json:"leadNo" db:"lead_no"`
	ContactName string `json:"contactName" db:"contact_name"`
	CompanyName string `json:"companyName" db:"company_name"`
	Phone       string `json:"phone" db:"phone"`
	Email       string `json:"email" db:"email"`
	Route       string `json:"route" db:"route"`
	CargoInfo   string `json:"cargoInfo" db:"cargo_info"`
	Message     string `json:"message" db:"message"`
	Source      string `json:"source" db:"source"`
	Status      string `json:"status" db:"status"`
	IpAddr      string `json:"ipAddr" db:"ip_addr"`
	UserAgent   string `json:"userAgent" db:"user_agent"`
	Remark      string `json:"remark" db:"remark"`
	CreateBy    string `json:"createBy" db:"create_by"`
	UpdateBy    string `json:"updateBy" db:"update_by"`
}

type ContactVo struct {
	ContactId   int64                `json:"contactId,string" db:"contact_id"`
	LeadNo      string               `json:"leadNo" db:"lead_no"`
	ContactName string               `json:"contactName" db:"contact_name"`
	CompanyName string               `json:"companyName" db:"company_name"`
	Phone       string               `json:"phone" db:"phone"`
	Email       string               `json:"email" db:"email"`
	Route       string               `json:"route" db:"route"`
	CargoInfo   string               `json:"cargoInfo" db:"cargo_info"`
	Message     string               `json:"message" db:"message"`
	Source      string               `json:"source" db:"source"`
	Status      string               `json:"status" db:"status"`
	IpAddr      string               `json:"ipAddr" db:"ip_addr"`
	UserAgent   string               `json:"userAgent" db:"user_agent"`
	Remark      string               `json:"remark" db:"remark"`
	CreateBy    string               `json:"createBy" db:"create_by"`
	CreateTime  *baizeUnix.BaiZeTime `json:"createTime" db:"create_time"`
	UpdateBy    string               `json:"updateBy" db:"update_by"`
	UpdateTime  *baizeUnix.BaiZeTime `json:"updateTime" db:"update_time"`
}

func ContactListToRows(contacts []*ContactVo, statusLabels map[string]string) (rows [][]string) {
	rows = make([][]string, 0, len(contacts)+1)
	rows = append(rows, []string{"线索编号", "联系人", "公司名称", "电话", "邮箱", "目标航线", "货物信息", "需求说明", "来源", "状态", "提交时间"})
	for _, contact := range contacts {
		createTime := ""
		if contact.CreateTime != nil {
			createTime = contact.CreateTime.ToString()
		}
		rows = append(rows, []string{
			contact.LeadNo,
			contact.ContactName,
			contact.CompanyName,
			contact.Phone,
			contact.Email,
			contact.Route,
			contact.CargoInfo,
			contact.Message,
			contact.Source,
			contactStatusLabel(contact.Status, statusLabels),
			createTime,
		})
	}
	return
}

func contactStatusLabel(status string, statusLabels map[string]string) string {
	if label := statusLabels[status]; label != "" {
		return label
	}
	switch status {
	case "10":
		return "待跟进"
	case "20":
		return "跟进中"
	case "30":
		return "已完成"
	case "40":
		return "无效"
	default:
		return status
	}
}

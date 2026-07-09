package models

import (
	"baize/app/common/baize/baizeUnix"
	"baize/app/common/commonModels"
)

type NotificationDQL struct {
	Title    string `form:"title" db:"title"`
	BizType  string `form:"bizType" db:"biz_type"`
	ReadFlag string `form:"readFlag" db:"read_flag"`
	commonModels.BaseEntityDQL
}

type NotificationDML struct {
	NotificationId int64  `json:"notificationId,string" db:"notification_id"`
	UserId         int64  `json:"userId,string" db:"user_id"`
	Title          string `json:"title" db:"title"`
	Content        string `json:"content" db:"content"`
	BizType        string `json:"bizType" db:"biz_type"`
	BizId          int64  `json:"bizId,string" db:"biz_id"`
	ReadFlag       string `json:"readFlag" db:"read_flag"`
	CreateBy       string `json:"createBy" db:"create_by"`
}

type NotificationVo struct {
	NotificationId int64                `json:"notificationId,string" db:"notification_id"`
	UserId         int64                `json:"userId,string" db:"user_id"`
	Title          string               `json:"title" db:"title"`
	Content        string               `json:"content" db:"content"`
	BizType        string               `json:"bizType" db:"biz_type"`
	BizId          int64                `json:"bizId,string" db:"biz_id"`
	ReadFlag       string               `json:"readFlag" db:"read_flag"`
	CreateBy       string               `json:"createBy" db:"create_by"`
	CreateTime     *baizeUnix.BaiZeTime `json:"createTime" db:"create_time"`
	ReadTime       *baizeUnix.BaiZeTime `json:"readTime" db:"read_time"`
}

type UnreadCountVo struct {
	UnreadCount int64 `json:"unreadCount"`
}

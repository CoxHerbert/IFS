package models

import (
	"baize/app/common/baize/baizeUnix"
	"baize/app/common/commonModels"
)

type ArticleDQL struct {
	Title    string `form:"title" db:"title"`
	Category string `form:"category" db:"category"`
	Status   string `form:"status" db:"status"`
	Keyword  string `form:"keyword" db:"keyword"`
	commonModels.BaseEntityDQL
}

type ArticleDML struct {
	ArticleId   int64  `json:"articleId,string" db:"article_id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Slug        string `json:"slug" db:"slug"`
	Summary     string `json:"summary" db:"summary"`
	Category    string `json:"category" db:"category" binding:"required"`
	CoverUrl    string `json:"coverUrl" db:"cover_url"`
	Content     string `json:"content" db:"content"`
	Status      string `json:"status" db:"status"`
	Sort        int64  `json:"sort" db:"sort"`
	PublishTime string `json:"publishTime" db:"publish_time"`
	CreateBy    string `json:"createBy" db:"create_by"`
	UpdateBy    string `json:"updateBy" db:"update_by"`
}

type ArticleVo struct {
	ArticleId   int64                `json:"articleId,string" db:"article_id"`
	Title       string               `json:"title" db:"title"`
	Slug        string               `json:"slug" db:"slug"`
	Summary     string               `json:"summary" db:"summary"`
	Category    string               `json:"category" db:"category"`
	CoverUrl    string               `json:"coverUrl" db:"cover_url"`
	Content     string               `json:"content" db:"content"`
	Status      string               `json:"status" db:"status"`
	Sort        int64                `json:"sort" db:"sort"`
	PublishTime string               `json:"publishTime" db:"publish_time"`
	CreateBy    string               `json:"createBy" db:"create_by"`
	CreateTime  *baizeUnix.BaiZeTime `json:"createTime" db:"create_time"`
	UpdateBy    string               `json:"updateBy" db:"update_by"`
	UpdateTime  *baizeUnix.BaiZeTime `json:"updateTime" db:"update_time"`
}

package service

import (
	"baize/app/cms/dao"
	"baize/app/cms/models"
	"baize/app/utils/snowflake"
	"regexp"
	"strings"
	"time"
)

var articleServiceImpl = &articleService{articleDao: dao.GetArticleDao()}

type articleService struct {
	articleDao interface {
		SelectList(query *models.ArticleDQL) (list []*models.ArticleVo, total *int64)
		SelectById(articleId int64) *models.ArticleVo
		SelectBySlug(slug string) *models.ArticleVo
		Insert(item *models.ArticleDML)
		Update(item *models.ArticleDML)
		DeleteByIds(articleIds []int64)
	}
}

func GetArticleService() *articleService {
	return articleServiceImpl
}

func (service *articleService) SelectList(query *models.ArticleDQL) (list []*models.ArticleVo, total *int64) {
	return service.articleDao.SelectList(query)
}

func (service *articleService) SelectPublishedList(query *models.ArticleDQL) (list []*models.ArticleVo, total *int64) {
	query.Status = "0"
	return service.articleDao.SelectList(query)
}

func (service *articleService) SelectById(articleId int64) *models.ArticleVo {
	return service.articleDao.SelectById(articleId)
}

func (service *articleService) SelectBySlug(slug string) *models.ArticleVo {
	return service.articleDao.SelectBySlug(slug)
}

func (service *articleService) Insert(item *models.ArticleDML, username string) {
	item.ArticleId = snowflake.GenID()
	normalizeArticle(item)
	item.CreateBy = username
	item.UpdateBy = username
	service.articleDao.Insert(item)
}

func (service *articleService) Update(item *models.ArticleDML, username string) {
	normalizeArticle(item)
	item.UpdateBy = username
	service.articleDao.Update(item)
}

func (service *articleService) DeleteByIds(articleIds []int64) {
	service.articleDao.DeleteByIds(articleIds)
}

func normalizeArticle(item *models.ArticleDML) {
	item.Title = strings.TrimSpace(item.Title)
	item.Category = strings.TrimSpace(item.Category)
	item.Slug = normalizeSlug(item.Slug)
	if item.Slug == "" {
		item.Slug = normalizeSlug(item.Title)
	}
	if item.Slug == "" {
		item.Slug = strings.ToLower(time.Now().Format("20060102150405"))
	}
	if item.Status == "" {
		item.Status = "0"
	}
	if item.PublishTime == "" && item.Status == "0" {
		item.PublishTime = time.Now().Format("2006-01-02 15:04:05")
	}
}

func normalizeSlug(value string) string {
	text := strings.ToLower(strings.TrimSpace(value))
	text = regexp.MustCompile(`[^a-z0-9\u4e00-\u9fa5]+`).ReplaceAllString(text, "-")
	return strings.Trim(text, "-")
}

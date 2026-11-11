package controller

import (
	"baize/app/cms/models"
	"baize/app/cms/service"
	"baize/app/common/baize/baizeContext"
	"baize/app/constant/constants"
	"baize/app/utils/fileUploadUtils"
	"baize/app/utils/slicesUtils"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var articleService = service.GetArticleService()

const maxArticleImageSize = 5 << 20

func PortalArticleList(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	query := new(models.ArticleDQL)
	c.ShouldBind(query)
	query.SetLimit(c)
	list, total := articleService.SelectPublishedList(query)
	bzc.SuccessListData(list, total)
}

func PortalArticleDetail(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	item := articleService.SelectBySlug(strings.TrimSpace(c.Param("slug")))
	if item == nil {
		bzc.Waring("文章不存在或未发布")
		return
	}
	bzc.SuccessData(item)
}

func ArticleList(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	query := new(models.ArticleDQL)
	c.ShouldBind(query)
	query.SetLimit(c)
	list, total := articleService.SelectList(query)
	bzc.SuccessListData(list, total)
}

func ArticleGetInfo(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	articleId := bzc.ParamInt64("articleId")
	if articleId == 0 {
		bzc.ParameterError()
		return
	}
	bzc.SuccessData(articleService.SelectById(articleId))
}

func ArticleAdd(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	item := new(models.ArticleDML)
	if err := c.ShouldBindJSON(item); err != nil {
		zap.L().Error("article add param error", zap.Error(err))
		bzc.ParameterError()
		return
	}
	articleService.Insert(item, bzc.GetCurrentUserName())
	bzc.Success()
}

func ArticleEdit(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	item := new(models.ArticleDML)
	if err := c.ShouldBindJSON(item); err != nil || item.ArticleId == 0 {
		zap.L().Error("article edit param error", zap.Error(err))
		bzc.ParameterError()
		return
	}
	articleService.Update(item, bzc.GetCurrentUserName())
	bzc.Success()
}

func ArticleRemove(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	var ids slicesUtils.Slices = strings.Split(c.Param("articleIds"), ",")
	articleService.DeleteByIds(ids.StrSlicesToInt())
	bzc.Success()
}

func ArticleUploadImage(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	file, err := c.FormFile("file")
	if err != nil {
		bzc.ParameterError()
		return
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if file.Size > maxArticleImageSize || !isAllowedArticleImage(file, ext) {
		bzc.Waring("图片仅支持 PNG、JPG、JPEG、GIF、WEBP，且不能超过5MB")
		return
	}
	url := constants.ResourcePrefix + fileUploadUtils.Upload(constants.CmsArticleImagePath, file)
	bzc.SuccessData(gin.H{"url": url})
}

func isAllowedArticleImage(fileHeader *multipart.FileHeader, ext string) bool {
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".gif" && ext != ".webp" {
		return false
	}
	file, err := fileHeader.Open()
	if err != nil {
		return false
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false
	}
	head := buffer[:n]
	contentType := http.DetectContentType(head)
	if contentType == "image/png" && ext == ".png" {
		return true
	}
	if contentType == "image/jpeg" && (ext == ".jpg" || ext == ".jpeg") {
		return true
	}
	if contentType == "image/gif" && ext == ".gif" {
		return true
	}
	return ext == ".webp" && isWebPHeader(head)
}

func isWebPHeader(head []byte) bool {
	return len(head) >= 12 &&
		string(head[0:4]) == "RIFF" &&
		string(head[8:12]) == "WEBP"
}

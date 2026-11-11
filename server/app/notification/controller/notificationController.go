package controller

import (
	"baize/app/common/baize/baizeContext"
	"baize/app/notification/models"
	"baize/app/notification/service"
	"baize/app/utils/slicesUtils"
	"strings"

	"github.com/gin-gonic/gin"
)

var notificationService = service.GetNotificationService()

func NotificationList(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	query := new(models.NotificationDQL)
	c.ShouldBind(query)
	query.SetLimit(c)
	list, total := notificationService.List(bzc.GetCurrentUserId(), query)
	bzc.SuccessListData(list, total)
}

func NotificationUnreadCount(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	bzc.SuccessData(models.UnreadCountVo{UnreadCount: notificationService.CountUnread(bzc.GetCurrentUserId())})
}

func NotificationRead(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	notificationId := bzc.ParamInt64("notificationId")
	if notificationId == 0 {
		bzc.ParameterError()
		return
	}
	notificationService.MarkRead(notificationId, bzc.GetCurrentUserId())
	bzc.Success()
}

func NotificationReadAll(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	notificationService.MarkAllRead(bzc.GetCurrentUserId())
	bzc.Success()
}

func NotificationRemove(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	var ids slicesUtils.Slices = strings.Split(c.Param("notificationIds"), ",")
	notificationService.DeleteByIds(ids.StrSlicesToInt(), bzc.GetCurrentUserId())
	bzc.Success()
}

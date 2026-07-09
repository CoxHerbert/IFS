package service

import (
	freightModels "baize/app/freight/models"
	"baize/app/notification/dao"
	notificationModels "baize/app/notification/models"
	"baize/app/utils/snowflake"
	"fmt"

	"go.uber.org/zap"
)

var notificationServiceImpl *notificationService

func init() {
	notificationServiceImpl = &notificationService{notificationDao: dao.GetNotificationDao()}
}

type notificationService struct {
	notificationDao interface {
		Insert(item *notificationModels.NotificationDML)
		SelectList(userId int64, query *notificationModels.NotificationDQL) (list []*notificationModels.NotificationVo, total *int64)
		CountUnread(userId int64) int64
		MarkRead(notificationId int64, userId int64) bool
		MarkAllRead(userId int64) int64
		DeleteByIds(notificationIds []int64, userId int64) int64
	}
}

func GetNotificationService() *notificationService {
	return notificationServiceImpl
}

func (service *notificationService) List(userId int64, query *notificationModels.NotificationDQL) (list []*notificationModels.NotificationVo, total *int64) {
	return service.notificationDao.SelectList(userId, query)
}

func (service *notificationService) CountUnread(userId int64) int64 {
	return service.notificationDao.CountUnread(userId)
}

func (service *notificationService) MarkRead(notificationId int64, userId int64) bool {
	return service.notificationDao.MarkRead(notificationId, userId)
}

func (service *notificationService) MarkAllRead(userId int64) int64 {
	return service.notificationDao.MarkAllRead(userId)
}

func (service *notificationService) DeleteByIds(notificationIds []int64, userId int64) int64 {
	return service.notificationDao.DeleteByIds(notificationIds, userId)
}

func (service *notificationService) NotifyShipmentCreated(plan *freightModels.ShipmentPlanDML, operatorName string, operatorUserId int64) {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Warn("create shipment notification failed", zap.Any("error", err), zap.Int64("shipmentId", plan.ShipmentId))
		}
	}()

	recipientUserId := plan.SalesUserId
	if recipientUserId == 0 {
		recipientUserId = operatorUserId
	}
	if recipientUserId == 0 {
		return
	}

	customerName := plan.CustomerName
	if customerName == "" {
		customerName = "未分配客户"
	}
	content := fmt.Sprintf("出货计划 %s 已创建，客户：%s，航线：%s -> %s。", plan.ShipmentNo, customerName, defaultText(plan.Pol), defaultText(plan.Pod))
	service.notificationDao.Insert(&notificationModels.NotificationDML{
		NotificationId: snowflake.GenID(),
		UserId:         recipientUserId,
		Title:          "出货计划已创建",
		Content:        content,
		BizType:        "shipment",
		BizId:          plan.ShipmentId,
		ReadFlag:       "0",
		CreateBy:       operatorName,
	})
}

func defaultText(value string) string {
	if value == "" {
		return "-"
	}
	return value
}

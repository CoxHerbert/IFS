# 消息通知

## 目标

消息通知模块用于把关键业务事件主动推送给后台用户。当前第一批落地场景是“出货计划创建后自动生成通知”。

当前能力：

- 出货计划创建成功后自动写入通知
- 后台顶部 Header 显示未读数和最近通知
- 后台“消息通知管理”页面支持筛选、标记已读、全部已读、删除

## 触发规则

当前触发点：

- 出货计划创建成功后触发通知

触发位置：

- `app/freight/service/shipmentService.go`

触发方法：

- `NotifyShipmentCreated(plan, operatorName, operatorUserId)`

## 接收人规则

优先级如下：

1. 优先发送给出货计划绑定客户的业务员 `sales_user_id`
2. 如果没有绑定业务员，则退回发送给当前后台操作人
3. 如果两者都为空，则本次不生成通知

## 数据表

- `sys_notification`

重点字段：

- `notification_id`
- `user_id`
- `title`
- `content`
- `biz_type`
- `biz_id`
- `read_flag`
- `read_time`

## SQL

- 统一入口：`sql/ifs_init.sql`
- 业务合并脚本：`sql/ifs_business.sql`

## 后端接口

- `GET /system/notification/list`
- `GET /system/notification/unread-count`
- `PUT /system/notification/:notificationId/read`
- `PUT /system/notification/read-all`
- `DELETE /system/notification/:notificationIds`

权限标识：

- `system:notification:list`
- `system:notification:edit`
- `system:notification:remove`

## 前端入口

顶部通知入口：

- `baize-ui/src/layout/components/Navbar.vue`

后台管理页：

- `baize-ui/src/views/system/notification/index.vue`

菜单位置：

- 系统管理 / 消息通知

## 当前实现方式

当前不是 WebSocket 实时推送，而是“后端落库 + 前端轮询展示”的实现方式：

- Header 进入页面后轮询未读数
- 打开通知面板时拉取最近通知
- 点击通知后调用已读接口

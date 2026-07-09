# 消息通知

## 目标

消息通知模块用于把关键业务事件主动推送给后台用户，当前第一批落地场景是“出货计划发起后自动生成通知”。

当前能力包括：

- 出货计划创建成功后自动写入通知
- 后台顶部 Header 显示未读数和最近通知
- 后台“消息通知管理”页支持筛选、标记已读、全部已读、删除

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

通知表：

- `sys_notification`

字段关注点：

- `notification_id`：通知主键
- `user_id`：接收人
- `title`：标题
- `content`：内容
- `biz_type`：业务类型，当前为 `shipment`
- `biz_id`：业务主键，当前对应 `shipment_id`
- `read_flag`：是否已读，`0` 未读，`1` 已读
- `read_time`：已读时间

相关 SQL：

- `sql/shipment_notification.sql`

## 后端接口

通知接口路由：

- `GET /system/notification/list`
- `GET /system/notification/unread-count`
- `PUT /system/notification/:notificationId/read`
- `PUT /system/notification/read-all`
- `DELETE /system/notification/:notificationIds`

权限标识：

- `system:notification:list`
- `system:notification:edit`
- `system:notification:remove`

后端代码位置：

- `app/notification/models`
- `app/notification/dao`
- `app/notification/service`
- `app/notification/controller`
- `app/routes/systemRouter/sysNotificationRouter.go`

## 前端入口

顶部通知入口：

- `baize-ui/src/layout/components/Navbar.vue`

后台管理页：

- 路由组件：`baize-ui/src/views/system/notification/index.vue`
- 菜单位置：系统管理 / 消息通知

接口封装：

- `baize-ui/src/api/system/notification.ts`

## 菜单与权限

后台消息通知管理菜单脚本：

- `sql/system_notification_menu.sql`

菜单作用：

- 新增“系统管理 / 消息通知”
- 分配查询、编辑、删除权限点

## 当前实现方式

当前不是 WebSocket 实时推送，而是“后端落库 + 前端轮询展示”的实现方式。

- Header 进入页面后轮询未读数
- 打开通知面板时拉取最近通知
- 点击通知后调用已读接口

这种方式实现成本低，适合当前单业务场景。后续如果通知类型增多、对实时性要求更高，可以再升级为 WebSocket 或 SSE。

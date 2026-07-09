# 出货计划与出货查询

## 目标

出货计划模块维护正式出货计划、货物明细、柜型/LCL 方案和客户查询数据。

Agent 生成的分析结果只有在提交“保存为正式出货计划”表单后，才会写入正式出货计划表。

## 数据表

- `freight_shipment_plan`
- `freight_shipment_cargo`
- `freight_container_plan`
- `freight_shipment_order`

SQL：`sql/freight_shipment.sql`

## 客户归属

| 来源 | 规则 |
| --- | --- |
| 客户端 Agent 发起 | 自动绑定当前客户端登录账号对应的客户。 |
| 后台 Agent 发起 | 先创建未绑定计划，`customer_id = 0`。 |
| 后台人工维护 | 可在出货计划列表中选择客户进行绑定。 |

绑定接口：

- `PUT /freight/shipment/:shipmentId/customer`

## LCL 与柜型规则

小体积优先 LCL，避免 2-3 CBM 这类货物被错误建议整柜。

当前规则：

- `< 15 CBM`：LCL 拼箱
- `<= 28 CBM`：`1×20GP`
- `<= 58 CBM`：`1×40GP`
- `<= 68 CBM`：`1×40HQ`
- `> 68 CBM`：`ceil(totalCBM / 68) × 40HQ`

## 客户端查询

客户只能查看当前登录账号绑定客户的数据。

接口：

- `GET /portal/customer/shipments`
- `GET /portal/customer/shipment/:shipmentId`
- `GET /portal/shipment/share/:token`

## 后台管理

菜单：货代业务 / 出货计划

主要接口：

- `GET /freight/shipment/list`
- `GET /freight/shipment/:shipmentId`
- `POST /freight/shipment`
- `PUT /freight/shipment`
- `DELETE /freight/shipment/:shipmentIds`
- `PUT /freight/shipment/:shipmentId/customer`

## 消息通知联动

出货计划创建成功后，系统会自动生成一条后台消息通知。

规则：

- 优先通知该客户绑定的业务员
- 如果没有绑定业务员，则通知当前后台操作人
- 通知内容会带出货计划编号、客户名和航线

详细说明见 [消息通知](notification-center.md)。

## 关键文件

后端：

- `app/freight`
- `app/routes/freightRoutes`

前端：

- `baize-ui/src/views/freight/shipment/index.vue`
- `portal-ui/src/views/workspace/WorkspaceShipmentTrackingView.vue`

Agent 分析和保存表单见 [IFS 本地 Agent](ifs-agent.md)。

# 出货计划与客户分享页功能文档

## 1. 功能目标

本功能面向国际货运代理业务，目标不是只做一张出货单，而是把客户最关心的“货在哪里、下一步是什么、是否有人在跟进”透明化。

核心链路：

1. 后台导入客户出货清单。
2. 系统根据货物体积和重量生成智能出货计划，并推荐需要的货柜。
3. 运营确认计划后生成出货单。
4. 运营维护出货状态。
5. 系统生成免登录分享链接。
6. 客户通过分享页查看单票货物状态、货物明细、柜型建议和关键节点。

## 2. 业务边界

当前版本解决的是“计划、排柜、状态透明、分享查询”。

已包含：

- 出货清单录入/导入接口。
- 货物明细保存。
- 按总体积、总重量推荐柜型。
- 出货计划列表、详情、状态维护。
- 出货计划确认并生成出货单。
- 分享 token 查询，不要求客户登录。
- 客户视角状态时间线。

暂未包含，但预留扩展空间：

- Excel 文件解析导入。
- 多柜拆分装载明细。
- 每个状态节点的操作日志。
- 附件/照片/回单上传。
- 订舱单、报关单、提单等正式单证。
- 船司、码头、海外仓 API 对接。
- 微信/邮件/短信自动通知。

## 3. 状态字典

字典类型：`freight_shipment_status`

状态值使用数字字符串，主流程按 10 递增，异常状态放到 900。这样方便排序、插入中间状态，也方便前端时间线展示。

| 状态值 | 状态名称 | 客户理解 |
| --- | --- | --- |
| 10 | 计划已创建 | 货代已建立这票货的出货计划 |
| 20 | 出货计划已确认 | 计划已确认，可以按计划推进 |
| 30 | 等待客户发货 | 等待客户送货或安排提货 |
| 40 | 已提货/已送仓 | 货物已经离开发货地 |
| 50 | 仓库已收货 | 仓库已接收货物 |
| 60 | 已入仓/码头进仓 | 仓库或码头已完成进仓 |
| 70 | 订舱处理中 | 正在安排舱位 |
| 80 | 舱位已确认 | 船司或渠道已确认舱位 |
| 90 | 报关资料已收齐 | 报关资料已齐备 |
| 100 | 报关已放行 | 出口报关已放行 |
| 110 | 已装柜 | 已完成装柜 |
| 120 | 已进港/码头放行 | 货柜已进港或码头放行 |
| 130 | 船舶已开航 | 船已经开航 |
| 140 | 目的港已到港 | 船已到目的港 |
| 150 | 目的港清关中 | 目的港正在清关 |
| 160 | 目的港已清关 | 目的港清关完成 |
| 170 | 已派送/已签收 | 派送中或客户已签收 |
| 900 | 异常处理中 | 发生异常，运营正在处理 |

实现注意：

- 后端状态流排序必须按数字比较，不能按字符串比较，否则 `100` 会错误地排在 `90` 前面。
- `900` 是异常态，不参与正常进度递增；只有当前状态为 `900` 时才显示为当前异常节点。
- 后续如果需要更细状态，可以插入 `75`、`85` 或新增 `180`，不要复用已有状态值。

## 4. 数据模型

SQL 文件：`sql/freight_shipment.sql`

### 4.1 freight_shipment_plan

出货计划主表，一票货一条记录。

关键字段：

- `shipment_id`：出货计划 ID。
- `shipment_no`：出货计划编号，例如 `SP20260630xxxxxx`。
- `order_no`：客户订单号或参考号。
- `customer_id` / `customer_name`：客户信息。
- `pol`：起运港。
- `pod`：目的港。
- `planned_etd` / `planned_eta`：计划开船/到港。
- `actual_etd` / `actual_eta`：实际开船/到港。
- `status`：出货状态，来自 `freight_shipment_status`。
- `total_weight`：总重量 KG。
- `total_volume`：总体积 CBM。
- `total_cartons`：总箱数。
- `share_token`：免登录分享令牌。

### 4.2 freight_shipment_cargo

货物明细表，一票货可以有多条货物明细。

关键字段：

- `cargo_name`：货名。
- `sku`：SKU、唛头或货号。
- `quantity`：件数。
- `cartons`：箱数。
- `weight_kg`：重量。
- `volume_cbm`：体积。
- `length_cm` / `width_cm` / `height_cm`：尺寸。

### 4.3 freight_container_plan

智能货柜计划表，保存系统推荐结果。

关键字段：

- `container_type`：柜型，例如 `20GP`、`40GP`、`40HQ`。
- `quantity`：柜量。
- `max_volume` / `max_weight`：推荐柜型可承载能力。
- `used_volume` / `used_weight`：当前货物占用。
- `load_rate`：装载率。
- `remark`：推荐说明。

### 4.4 freight_shipment_order

出货单表。运营确认出货计划后生成。

关键字段：

- `order_id`：出货单 ID。
- `shipment_id`：关联出货计划。
- `order_no`：出货单号。
- `status`：生成时的出货状态。

## 5. 智能货柜推荐逻辑

当前版本采用保守规则：

1. 先汇总所有货物的 `total_volume` 和 `total_weight`。
2. 如果用户指定偏好柜型，则按偏好柜型计算柜量。
3. 如果未指定偏好柜型，则优先选择能一次装下当前货量的最小柜型。
4. 如果单柜装不下，则按体积和重量瓶颈计算所需柜量。
5. 装载率取 `体积占用率` 和 `重量占用率` 的较大值。

默认柜型能力：

| 柜型 | 参考体积 | 参考限重 |
| --- | --- | --- |
| 20GP | 28 CBM | 21700 KG |
| 40GP | 58 CBM | 26500 KG |
| 40HQ | 68 CBM | 26500 KG |

后续可增强：

- 按长宽高做装箱算法。
- 支持混柜、多柜、多目的港。
- 对超长、超重、危险品、冷链货做特殊规则。
- 输出“为什么推荐这个柜型”的更详细解释。

## 6. 后端接口

### 6.1 后台接口

接口前缀：`/freight/shipment`

| 方法 | 路径 | 权限 | 用途 |
| --- | --- | --- | --- |
| GET | `/list` | `freight:shipment:list` | 查询出货计划列表 |
| POST | `/import` | `freight:shipment:import` | 导入出货清单并生成计划 |
| GET | `/:shipmentId` | `freight:shipment:query` | 查看计划详情 |
| PUT | `/:shipmentId/status` | `freight:shipment:edit` | 更新客户可见状态 |
| POST | `/:shipmentId/confirm` | `freight:shipment:confirm` | 确认计划并生成出货单 |
| GET | `/:shipmentId/share` | `freight:shipment:share` | 获取免登录分享链接 |
| DELETE | `/:shipmentIds` | `freight:shipment:remove` | 删除出货计划 |

### 6.2 门户公开接口

接口前缀：`/portal/shipment`

| 方法 | 路径 | 登录 | 用途 |
| --- | --- | --- | --- |
| GET | `/share/:token` | 不需要 | 通过分享 token 查看单票货状态 |

## 7. 前端页面

### 7.1 后台管理页面

文件：`baize-ui/src/views/freight/shipment/index.vue`

主要能力：

- 按计划编号、客户、目的港、状态筛选。
- 导入出货清单。
- 录入客户、起运港、目的港、偏好柜型。
- 录入多行货物明细。
- 查看智能货柜建议。
- 查看货物明细。
- 更新状态。
- 生成出货单。
- 复制客户分享链接。

### 7.2 客户分享页面

文件：`portal-ui/src/views/ShipmentShareView.vue`

路由：

`/shipment/share/:token`

页面展示：

- 出货计划编号。
- 当前状态。
- 起运港到目的港。
- 状态时间线。
- 客户参考号。
- 出货单号。
- 箱数、体积、重量。
- 计划/实际 ETD、ETA。
- 推荐货柜。
- 货物明细。

设计原则：

- 不要求客户登录，减少沟通摩擦。
- 展示客户能理解的状态，不暴露内部复杂字段。
- 用明确时间线增强信任感。

## 8. 权限与菜单

菜单脚本在 `sql/freight_shipment.sql`。

后台权限点：

- `freight:shipment:list`
- `freight:shipment:query`
- `freight:shipment:import`
- `freight:shipment:edit`
- `freight:shipment:confirm`
- `freight:shipment:share`
- `freight:shipment:remove`

菜单位置：

- 货代业务
- 出货计划

生产环境执行 SQL 前注意：

- 当前 SQL 包含 `DROP TABLE IF EXISTS`，适合新库或测试库。
- 已有生产数据时不要直接执行，应改成增量迁移脚本。
- `sys_menu`、`sys_dict_type`、`sys_dict_data` 的 ID 如与现有库冲突，需要调整。

## 9. 当前实现文件

后端：

- `app/freight/models/shipment.go`
- `app/freight/dao/shipmentDao.go`
- `app/freight/service/shipmentService.go`
- `app/freight/controller/shipmentController.go`
- `app/routes/freightRoutes/shipmentRouter.go`
- `app/routes/routes.go`

后台前端：

- `baize-ui/src/api/freight/shipment.js`
- `baize-ui/src/views/freight/shipment/index.vue`

客户门户：

- `portal-ui/src/api/customer.ts`
- `portal-ui/src/router/index.ts`
- `portal-ui/src/views/ShipmentShareView.vue`

数据库：

- `sql/freight_shipment.sql`

## 10. 后续完善建议

### 10.1 导入能力

建议下一步做 Excel 模板导入：

- 提供标准模板下载。
- 支持批量解析货名、SKU、箱数、重量、体积、尺寸。
- 导入前做数据预览和错误行提示。
- 导入后再生成计划，避免脏数据直接入库。

### 10.2 状态轨迹表

当前状态只保存在主表。后续建议新增 `freight_shipment_status_log`：

- `log_id`
- `shipment_id`
- `from_status`
- `to_status`
- `status_label`
- `remark`
- `operator`
- `operate_time`
- `customer_visible`

这样客户分享页可以展示“更新时间”和“节点说明”，运营也能追溯是谁改了状态。

### 10.3 分享链接安全

当前分享 token 是随机字符串。后续可增强：

- 支持分享链接过期时间。
- 支持重新生成 token。
- 支持关闭分享。
- 支持查看访问次数和最后访问时间。

### 10.4 客户协同

可继续加入：

- 客户确认出货计划。
- 客户上传报关资料。
- 客户留言。
- 客户对异常状态发起询问。
- 运营回复记录。

### 10.5 异常处理

`900 异常处理中` 只是总异常态。后续可以拆分异常原因：

- 资料缺失。
- 仓库差异。
- 海关查验。
- 船期延误。
- 目的港滞箱。
- 派送异常。

建议异常原因使用单独字段或单独字典，不要把所有异常都塞进主流程状态。

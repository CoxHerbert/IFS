# 出货计划与客户分享页说明

## 功能概览

该功能面向国际货代业务，核心目标是把出货计划、状态进度和客户查询入口连起来。

当前链路：

1. 后台录入或导入出货清单。
2. 系统根据货物体积和重量生成出货计划，并推荐货柜。
3. 运营确认计划后生成出货单。
4. 运营持续维护出货状态。
5. 系统生成免登录分享链接。
6. 客户通过分享页查看状态、货物明细和货柜建议。

补充：当执行“确认计划”时，如果当前状态仍小于 `20`，系统会先把状态推进到 `20`，再生成出货单。

## 数据与脚本

初始化脚本：`sql/freight_shipment.sql`

脚本已合并以下内容：

- 出货计划主表 `freight_shipment_plan`
- 货物明细表 `freight_shipment_cargo`
- 货柜计划表 `freight_container_plan`
- 出货单表 `freight_shipment_order`
- 状态字典 `freight_shipment_status`
- 货柜类型字典 `freight_container_type`
- 后台“货代业务 -> 出货计划”菜单及按钮权限

## 状态字典

字典类型：`freight_shipment_status`

主要状态：

- `10` 计划已创建
- `20` 出货计划已确认
- `30` 等待客户发货
- `40` 已提货/已送仓
- `50` 仓库已收货
- `60` 已入仓/码头进仓
- `70` 订舱处理中
- `80` 舱位已确认
- `90` 报关资料已收齐
- `100` 报关已放行
- `110` 已装柜
- `120` 已进港/码头放行
- `130` 船舶已开船
- `140` 目的港已到港
- `150` 目的港清关中
- `160` 目的港已清关
- `170` 已派送/已签收
- `900` 异常处理中

实现注意：状态排序必须按数值比较，不能按字符串比较。

## 接口

后台接口前缀：`/freight/shipment`

- `GET /list`：查询出货计划列表
- `POST /import`：导入出货清单并生成计划
- `GET /:shipmentId`：查看计划详情
- `PUT /:shipmentId/status`：更新客户可见状态
- `POST /:shipmentId/confirm`：确认计划并生成出货单
- `GET /:shipmentId/share`：获取分享链接
- `DELETE /:shipmentIds`：删除出货计划

门户公开接口前缀：`/portal/shipment`

- `GET /share/:token`：通过分享 token 查看出货详情

## 前端页面

后台页面：`baize-ui/src/views/freight/shipment/index.vue`

客户端分享页：`portal-ui/src/views/portal/PortalShipmentShareView.vue`

分享链接生成说明：

- 优先使用环境变量 `VITE_PORTAL_BASE_URL`
- 未配置时回退为当前访问域名

## 当前实现文件

后端：

- `app/freight/models/shipment.go`
- `app/freight/dao/shipmentDao.go`
- `app/freight/service/shipmentService.go`
- `app/freight/controller/shipmentController.go`
- `app/routes/freightRoutes/shipmentRouter.go`

前端：

- `baize-ui/src/api/freight/shipment.js`
- `baize-ui/src/views/freight/shipment/index.vue`
- `portal-ui/src/api/portal/shipment.ts`
- `portal-ui/src/views/portal/PortalShipmentShareView.vue`

数据库：

- `sql/freight_shipment.sql`

## 注意事项

`sql/freight_shipment.sql` 是初始化脚本，包含 `DROP TABLE IF EXISTS`，适合新库、测试库和重建场景；线上已有数据时不要直接执行，应改为增量迁移脚本。

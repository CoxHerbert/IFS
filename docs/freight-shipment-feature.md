# 出货计划与客户端出货助手说明

## 1. 功能概览

当前出货相关能力分成两条链路：

### 后台出货计划链路

1. 后台录入或导入出货明细
2. 系统汇总重量、体积、箱数
3. 系统推荐整柜方案
4. 运营确认计划后生成出货单
5. 持续维护出货状态
6. 生成分享链接给客户查询

### 客户端出货助手链路

1. 客户导入 Excel
2. 客户在前端表格里修正货物明细
3. 系统自动补算体积
4. 系统返回整柜推荐和散货建议
5. 客户把整理后的数据再交给运营或用于后续下单

客户端助手当前只做测算，不直接落库。

## 2. 数据与脚本

初始化脚本：

- `sql/freight_shipment.sql`
- `sql/customer_management.sql`
- `sql/customer_workspace_permission_migration.sql`

其中：

- `freight_shipment.sql` 负责正式出货计划表结构和后台菜单
- `customer_management.sql` 负责客户端菜单、角色和默认工作台菜单
- `customer_workspace_permission_migration.sql` 负责旧客户端权限表迁移

## 3. 后台接口

后台出货接口前缀：

- `/freight/shipment`

接口列表：

- `GET /list` 查询出货计划
- `POST /import` 导入货物并生成出货计划
- `GET /:shipmentId` 查询详情
- `PUT /:shipmentId/status` 更新状态
- `POST /:shipmentId/confirm` 确认计划并生成出货单
- `GET /:shipmentId/share` 获取分享链接
- `DELETE /:shipmentIds` 删除计划

## 4. 门户与客户端接口

公开分享接口：

- `GET /portal/shipment/share/:token`

客户端工作台接口：

- `POST /portal/customer/login`
- `GET /portal/customer/profile`
- `GET /portal/customer/routers`
- `POST /portal/customer/shipment-assistant/estimate`

其中 `/portal/customer/shipment-assistant/estimate` 为本轮新增，用于客户端 Excel 测算。

## 5. 柜型与测算规则

当前默认柜型容量：

- `20GP`：`28 CBM` / `21700 KG`
- `40GP`：`58 CBM` / `26500 KG`
- `40HQ`：`68 CBM` / `26500 KG`

测算逻辑：

1. 逐行读取货物
2. 如果 `volumeCbm` 为空，且存在 `lengthCm + widthCm + heightCm + cartons`，则自动换算体积
3. 汇总总数量、总箱数、总重量、总体积
4. 按体积和重量的较大占用比推荐柜数
5. 当体积较小或客户手动选择 `LCL` 时，同时给出散货建议

## 6. 前端页面

后台页面：

- `baize-ui/src/views/freight/shipment/index.vue`

门户分享页：

- `portal-ui/src/views/portal/PortalShipmentShareView.vue`

客户端工作台页面：

- `portal-ui/src/views/workspace/WorkspaceShipmentTrackingView.vue`
- `portal-ui/src/views/workspace/WorkspaceShipmentAssistantView.vue`

客户端助手使用：

- `vxe-table`
- `xlsx`

## 7. 关键文件

后端：

- `app/freight/models/shipment.go`
- `app/freight/service/shipmentService.go`
- `app/freight/controller/shipmentController.go`
- `app/customer/controller/customerController.go`
- `app/routes/customerRoutes/customerRouter.go`

前端：

- `portal-ui/src/api/portal/shipment.ts`
- `portal-ui/src/api/workspace/shipmentAssistant.ts`
- `portal-ui/src/views/workspace/WorkspaceShipmentTrackingView.vue`
- `portal-ui/src/views/workspace/WorkspaceShipmentAssistantView.vue`

数据库：

- `sql/freight_shipment.sql`
- `sql/customer_management.sql`
- `sql/customer_workspace_permission_migration.sql`

## 8. 注意事项

- 客户端工作台页面是否展示，取决于后台菜单配置和角色分配
- 新增客户端菜单时，前端必须同步补 `workspaceComponentMap`
- `customer_management.sql` 中默认已加入“智能出货助手”菜单
- `customer_workspace_permission_migration.sql` 已补迁移兼容
- `freight_shipment.sql` 仍是初始化脚本，线上有数据时不要直接全量执行

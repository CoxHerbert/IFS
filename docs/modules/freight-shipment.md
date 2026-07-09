# 出货计划与出货查询

## 目标

出货计划模块用于维护正式出货计划、货物明细、柜型或 LCL 方案，以及客户查询数据。

Agent 生成的分析结果只有在提交“保存为正式出货计划”后，才会写入正式出货计划表。

## 数据表

- `freight_shipment_plan`
- `freight_shipment_cargo`
- `freight_container_plan`
- `freight_shipment_order`

## SQL

- 统一入口：`sql/ifs_init.sql`
- 业务合并脚本：`sql/ifs_business.sql`

出货状态和柜型字典已经合并到 `sql/ifs_business.sql` 中统一维护：

- `freight_shipment_status`
- `freight_container_type`

## 客户归属

| 来源 | 规则 |
| --- | --- |
| 客户端 Agent 发起 | 自动绑定当前客户端登录账号对应的客户 |
| 后台 Agent 发起 | 先创建未绑定计划，`customer_id = 0` |
| 后台人工维护 | 可在出货计划列表中选择客户进行绑定 |

绑定接口：

- `PUT /freight/shipment/:shipmentId/customer`

## LCL 与柜型规则

小体积优先走 LCL，避免 2-3 CBM 这类货物被错误建议整柜。

当前规则：

- `< 15 CBM`：LCL 拼箱
- `<= 28 CBM`：`1 x 20GP`
- `<= 58 CBM`：`1 x 40GP`
- `<= 68 CBM`：`1 x 40HQ`
- `> 68 CBM`：`ceil(totalCBM / 68) x 40HQ`

## 客户端查询

接口：

- `GET /portal/customer/shipments`
- `GET /portal/customer/shipment/:shipmentId`
- `GET /portal/shipment/share/:token`

分享页说明：

- 支持免登录访问
- 移动端优先展示
- 当前已压缩状态区域布局，减小头部占用
- 付款状态和付款金额放在头部摘要区域

## 后台管理

菜单：

- 货代业务 / 出货计划

主要接口：

- `GET /freight/shipment/list`
- `GET /freight/shipment/:shipmentId`
- `POST /freight/shipment`
- `PUT /freight/shipment`
- `DELETE /freight/shipment/:shipmentIds`
- `PUT /freight/shipment/:shipmentId/customer`

## 后台页面说明

后台出货计划页面使用 `vxe` 组件体系维护，当前实现包括：

- 查询区使用 `vxe-form`
- 工具栏使用 `vxe-toolbar`
- 列表使用 `vxe-table`
- 分页使用 `vxe-pager`
- 导入、详情、状态维护、绑定客户、分享链接均使用 `vxe-modal`

页面能力：

- 按计划编号、客户名称、目的港、状态筛选
- 导入出货清单并生成正式计划
- 维护客户可见状态
- 绑定客户
- 生成出货单
- 删除出货计划
- 生成分享链接

分享链接交互：

- 点击“分享”后自动复制链接
- 弹窗内展示完整链接
- 可再次点击“复制链接”重复复制

## 状态字典

后台列表、状态维护弹窗、门户查询页统一使用字典：

- 字典类型：`freight_shipment_status`

当前字典值由 `sql/ifs_business.sql` 初始化，前后端不再写死状态名称。

## 关键文件

后端：

- `app/freight`
- `app/routes/freightRoutes`
- `app/freight/service/shipmentService.go`

后台前端：

- `baize-ui/src/views/freight/shipment/index.vue`

门户前端：

- `portal-ui/src/views/workspace/WorkspaceShipmentTrackingView.vue`

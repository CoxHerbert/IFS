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

## 2026-07 规则更新

### 客户绑定

- 发起出货计划时必须绑定客户，不允许创建未绑定客户的计划。
- 后台导入、后台 Agent、客户端工作台创建计划时都必须带客户信息。
- 后台修改绑定客户后，需要同步刷新计划上的业务员归属快照。

### 付款信息

- `freight_shipment_plan` 已增加付款状态和付款金额字段。
- 付款状态、付款金额由后台管理系统录入和维护。
- 客户端工作台、客户端详情页、免登录分享页只能查看付款状态和付款金额，不能修改。
- 当前付款状态固定为 `UNPAID`、`PARTIAL`、`PAID`。
- 收款作为独立功能维护，通过核销明细关联出货计划。
- 一笔收款默认核销一个出货计划，也可拆分核销到多个出货计划。
- 每笔付款包含金额、币种、付款时间、付款方式、备注，以及可选付款凭证。
- 付款凭证支持 PDF、PNG、JPG，单个文件最大 10MB。
- 新增或删除付款记录时，计划头上的付款金额会自动按明细合计，供列表快速展示。
- 后台“付款申报”菜单用于审核客户提交的凭证；审核通过后事务性生成正式收款并完成默认一对一全额核销，驳回时必须填写原因。

### 页面影响

- 后台出货计划列表、导入弹窗、详情弹窗支持展示付款状态和付款金额。
- 客户端详情页与分享页展示付款状态和付款金额摘要。
- 客户端“根据报告生成出货计划”不再提供付款录入入口。

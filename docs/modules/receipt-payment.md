# 收款与付款申报

## 目标

这一组能力用于把“正式收款”和“客户付款申报”拆开管理，避免客户端直接写正式收款，也避免管理端把核销动作和申报动作混在一起。

## 功能分类

### 1. 独立收款

- 后台存在独立的收款单管理页面
- 一张收款单可以关联多个出货计划
- 默认录入方式仍然按一对一处理，必要时可扩展为一对多核销
- 收款单支持上传付款凭证

### 2. 付款申报

- 客户端提交付款申报
- 申报可附带付款凭证文件
- 申报进入后台审核后，才生成正式收款单并完成核销
- 驳回时必须填写原因

### 3. 客户端查询

- 客户端可查看自己的收款核销摘要
- 客户端可查看自己的付款申报状态
- 客户端只读，不直接维护正式收款数据

## 关键数据表

- `freight_receipt`
- `freight_receipt_allocation`
- `freight_payment_declaration`

## 关键接口

### 后台收款

- `GET /freight/receipt/list`
- `GET /freight/receipt/:receiptId`
- `POST /freight/receipt`
- `DELETE /freight/receipt/:receiptId`

### 后台付款申报

- `GET /freight/payment-declaration/list`
- `GET /freight/payment-declaration/:declarationId`
- `POST /freight/payment-declaration/:declarationId/approve`
- `POST /freight/payment-declaration/:declarationId/reject`

### 客户端查询

- `GET /portal/customer/receipts`
- `GET /portal/customer/payment-declarations`

## 关联页面

- `baize-ui/src/views/freight/receipt/index.vue`
- `baize-ui/src/views/freight/paymentDeclaration/index.vue`
- `portal-ui/src/views/workspace/WorkspaceShipmentTrackingView.vue`
- `portal-ui/src/views/workspace/WorkspaceAgentChatView.vue`

## 关联说明

- 付款凭证支持 `PDF`、`PNG`、`JPG`
- 单文件最大 `10MB`
- 一张收款单允许分摊到多个出货计划
- 付款申报审核通过后，会生成正式收款单并执行默认核销

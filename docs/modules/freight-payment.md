# 收款与付款管理

## 目标

收款与付款能力从出货计划中独立出来，避免把收款单、付款申报、凭证上传和核销逻辑继续塞进出货计划页面里。

## 业务边界

- 收款单是独立业务功能
- 一张收款单默认对应一条核销明细，也支持一对多分摊
- 付款申报由客户侧提交，后台审核后生成正式收款记录
- 付款和收款都支持凭证上传
- 核销动作必须通过后台管理端完成

## 页面与入口

- 后台收款管理：`/freight/receipt`
- 后台付款申报：`/freight/payment-declaration`
- 客户端付款申报入口：客户工作台或 Agent 表单

## 数据结构

- `freight_receipt`
- `freight_receipt_allocation`
- `freight_payment_declaration`

## 关键规则

- 收款单可关联多个出货计划核销明细
- 默认创建场景按一对一展示，必要时可拆分为多条核销
- 收款凭证支持 PDF、PNG、JPG
- 核销金额不得超过收款单金额
- 付款申报审核通过后，再生成正式收款和核销记录

## 相关文件

- `app/freight/models/receipt.go`
- `app/freight/service/receiptService.go`
- `app/freight/controller/receiptController.go`
- `app/freight/models/paymentDeclaration.go`
- `app/freight/service/paymentDeclarationService.go`
- `baize-ui/src/views/freight/receipt/index.vue`
- `baize-ui/src/views/freight/paymentDeclaration/index.vue`

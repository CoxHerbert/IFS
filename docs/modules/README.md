# 模块文档索引

本目录存放 IFS 已落地模块的说明文档，重点记录业务目标、页面入口、关键接口、数据表和权限边界。

## 站点与客户入口

- [门户站点](portal-site.md)  
  说明官网首页、新闻资讯、服务页、联系页和悬浮 Agent 入口。
- [客户工作台](customer-workspace.md)  
  说明客户登录、动态菜单、客户侧出货查询、客户侧 Agent 对话等能力。

## 内容与系统管理

- [CMS 管理](cms-management.md)  
  说明后台新闻资讯管理、富文本编辑、图片上传、门户展示和安全边界。
- [系统基础管理](system-management.md)  
  说明后台菜单、权限、图标选择与基础展示规则。

## 业务执行模块

- [出货计划与出货查询](freight-shipment.md)  
  说明后台出货计划、客户归属、LCL/柜型规则、客户查询与后台管理接口。
- [收款与付款管理](freight-payment.md)  
  说明收款单、付款申报、凭证上传和核销边界。
- [收款与付款](receipt-payment.md)  
  说明收款付款的补充业务规则。
- [消息通知](notification-center.md)  
  说明出货计划触发通知、通知中心、后台消息通知管理和相关 SQL。

## 智能能力

- [IFS Agent](ifs-agent.md)  
  说明 Agent 页面入口、会话能力、本地技能和 Block Protocol 结构。

## 需求与扩展方案

- [智能出货工具需求说明](shipment-tool-requirements.md)  
  说明智能出货工具的完整业务需求、规则设计和扩展方向。

## 推荐阅读顺序

1. 先看 [门户站点](portal-site.md)、[CMS 管理](cms-management.md) 和 [客户工作台](customer-workspace.md)。
2. 涉及后台菜单、图标和权限时，看 [系统基础管理](system-management.md)。
3. 业务执行优先看 [出货计划与出货查询](freight-shipment.md)、[收款与付款管理](freight-payment.md) 和 [消息通知](notification-center.md)。
4. 涉及智能助手时，看 [IFS Agent](ifs-agent.md)。
5. 做新能力设计时，再看 [智能出货工具需求说明](shipment-tool-requirements.md)。

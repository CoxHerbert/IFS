# 模块文档索引

本目录存放 IFS 已落地模块的说明文档，重点描述业务目标、入口页面、关键接口、核心表结构和权限边界。

## 一、站点与客户入口

- [门户站点](portal-site.md)
  说明官网首页、服务页、联系页、新增页面和浮动 Agent 入口。
- [客户工作台](customer-workspace.md)
  说明客户登录、动态菜单、客户侧出货查询、客户侧 Agent 对话等能力。

## 二、系统与基础管理

- [系统基础管理](system-management.md)
  说明后台菜单、权限、图标选择与基础展示规则。

## 三、业务执行模块

- [出货计划与出货查询](freight-shipment.md)
  说明后台出货计划、客户归属、LCL/柜型规则、客户查询与后台管理接口。
- [收款与付款管理](freight-payment.md)
  说明收款单、付款申报、凭证上传和核销边界。
- [消息通知](notification-center.md)
  说明出货计划触发通知、通知中心、后台消息通知管理和相关 SQL。

## 四、智能能力

- [IFS Agent](ifs-agent.md)
  说明 Agent 页面入口、会话能力、本地技能和 Block Protocol 结构。

## 五、需求与扩展方案

- [智能出货工具需求说明](shipment-tool-requirements.md)
  说明智能出货工具的完整业务需求、规则设计和扩展方向。

## 推荐阅读顺序

1. 先看 [门户站点](portal-site.md) 和 [客户工作台](customer-workspace.md)。
2. 涉及后台菜单、图标和权限时，先看 [系统基础管理](system-management.md)。
3. 再看 [出货计划与出货查询](freight-shipment.md)、[收款与付款管理](freight-payment.md) 与 [消息通知](notification-center.md)。
4. 如果涉及智能助手，再看 [IFS Agent](ifs-agent.md)。
5. 如果要做新能力设计，再看 [智能出货工具需求说明](shipment-tool-requirements.md)。

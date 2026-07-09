# 客户工作台

## 目标

客户工作台用于客户登录后查看自己的业务数据，并按后台配置动态加载客户侧菜单。

当前能力：

- 客户账号独立登录
- 客户端菜单、角色、账号角色由后台管理
- 客户只能查看绑定到自己客户主体的出货计划
- 客户可使用智能出货助手和 Agent 对话

## 入口

| 场景 | 路径 |
| --- | --- |
| 客户工作台 | `/customer` |
| 客户登录页 | `/customer-login` |
| 出货查询 | `/customer/shipment` |
| 智能出货助手 | `/customer/shipment-assistant` |
| Agent 对话 | `/customer/agent-chat` |

## 动态菜单

后台菜单表：

- `customer_workspace_menu`

前端组件映射：

- `portal-ui/src/router/workspace-runtime.ts`

当前组件编码：

| 菜单 | component |
| --- | --- |
| 工作台 | `workspace/dashboard` |
| 账号资料 | `workspace/account-profile` |
| 出货查询 | `workspace/shipment-tracking` |
| 智能出货助手 | `workspace/shipment-assistant` |
| Agent 对话 | `workspace/agent-chat` |

## 出货查询

客户只能查看当前登录账号绑定客户的数据。

接口：

- `GET /portal/customer/shipments`
- `GET /portal/customer/shipment/:shipmentId`
- `GET /portal/shipment/share/:token`

## Agent 对话

组件：

- `portal-ui/src/views/workspace/WorkspaceAgentChatView.vue`

权限标识：

- `portal:agentChat:view`

客户登录后，前端会把客户 `token` 放入 `Authorization` 请求头。后端按客户工作台账号隔离会话，避免客户之间互相看到同一批对话记录。

## SQL

- 新环境统一入口：`sql/ifs_init.sql`
- 业务合并脚本：`sql/ifs_business.sql`

执行顺序见：

- [初始化与升级](../operations/setup.md)

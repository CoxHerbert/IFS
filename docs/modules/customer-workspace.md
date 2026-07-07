# 客户与客户端工作台

## 目标

客户工作台用于客户登录后查看自己的业务数据，并按后台配置加载客户端菜单。

当前能力：

- 客户账号独立登录
- 客户端菜单、角色、账号角色由后台管理
- 客户可查看绑定到自己客户主体的出货计划
- 客户可使用智能出货助手和 Agent 对话

## 入口

| 场景 | 路径 |
| --- | --- |
| 客户工作台 | `/customer` |
| 客户登录页 | `/customer-login` |
| 出货查询 | `/customer/shipment` |
| 智能出货助手 | `/customer/shipment-assistant` |
| Agent 对话 | `/customer/agent-chat` |

## 客户端动态菜单

后端菜单表：`customer_workspace_menu`

前端组件映射：`portal-ui/src/router/workspace-runtime.ts`

当前组件编码：

| 菜单 | component |
| --- | --- |
| 工作台 | `workspace/dashboard` |
| 账号资料 | `workspace/account-profile` |
| 出货查询 | `workspace/shipment-tracking` |
| 智能出货助手 | `workspace/shipment-assistant` |
| Agent 对话 | `workspace/agent-chat` |

## 客户端 Agent 对话

组件：`portal-ui/src/views/workspace/WorkspaceAgentChatView.vue`

权限标识：`portal:agentChat:view`

客户端登录后，前端会把客户 `token` 放入 `Authorization` 请求头。后端按客户工作台账号隔离会话，避免客户之间、客户与后台用户之间看到同一批对话记录。

当前前端交互补充：

- 点击对话标题直接进入内联编辑
- 不再展示单独的重命名按钮
- 对话发起区已统一为深色卡片样式，只保留上传与发送

Agent 详细协议和接口见 [IFS 本地 Agent](ifs-agent.md)。

## 客户端 Header

组件：`portal-ui/src/layouts/workspace/components/WorkspaceHeader.vue`

当前 Header 规则：

- 保留折叠菜单、刷新、面包屑、全屏、账号菜单
- 新增设置图标入口
- 主题切换（浅色 / 深色）已从账号下拉迁移到设置下拉中
- 账号下拉只保留返回门户、账号资料和退出登录

## 出货查询

客户只能查看当前登录账号绑定客户的数据。

接口：

- `GET /portal/customer/shipments`
- `GET /portal/customer/shipment/:shipmentId`
- `GET /portal/shipment/share/:token`

正式出货计划规则见 [出货计划与出货查询](freight-shipment.md)。

## SQL

- 新环境：`sql/customer_management.sql`
- 旧客户权限表迁移：`sql/customer_workspace_permission_migration.sql`

执行顺序见 [SQL 初始化与升级](../operations/sql.md)。

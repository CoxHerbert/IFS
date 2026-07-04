# SQL 初始化与升级

## 当前保留脚本

| 脚本 | 用途 |
| --- | --- |
| `sql/baize2022-01-08.sql` | 基础系统表、基础菜单、角色、字典。 |
| `sql/customer_management.sql` | 客户、联系人、客户账号、客户工作台菜单与权限。 |
| `sql/freight_shipment.sql` | 出货计划、货物明细、柜型/LCL、出货状态和后台出货菜单。 |
| `sql/ifs_agent.sql` | Agent 会话、消息、记忆、动态表单和后台 Agent 菜单。 |
| `sql/customer_workspace_permission_migration.sql` | 仅用于旧版 `customer_portal_*` 表迁移到 `customer_workspace_*`，并补齐客户端菜单权限。 |

## 全新环境执行顺序

1. `sql/baize2022-01-08.sql`
2. `sql/customer_management.sql`
3. `sql/freight_shipment.sql`
4. `sql/ifs_agent.sql`

## 已有环境升级原则

- 不要在生产库直接重跑带 `DROP TABLE` 的全量脚本。
- 客户端旧权限表需要迁移时，执行 `sql/customer_workspace_permission_migration.sql`。
- Agent 功能统一执行或对照 `sql/ifs_agent.sql` 做增量。
- 出货计划和客户模块已有数据时，应按现场差异写增量脚本。

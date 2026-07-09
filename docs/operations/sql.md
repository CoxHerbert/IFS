# SQL 初始化与升级

## 当前保留脚本

| 脚本 | 用途 |
| --- | --- |
| `sql/baize2022-01-08.sql` | 基础系统表、基础菜单、角色、字典 |
| `sql/customer_management.sql` | 客户、联系人、客户账号、客户工作台菜单与权限 |
| `sql/freight_shipment.sql` | 出货计划、货物明细、柜型/LCL、出货状态和后台出货菜单 |
| `sql/ifs_agent.sql` | Agent 会话、消息、记忆、动态表单和后台 Agent 菜单 |
| `sql/shipment_notification.sql` | 消息通知表 `sys_notification` |
| `sql/system_notification_menu.sql` | 后台“系统管理 / 消息通知”菜单及权限点 |
| `sql/customer_workspace_permission_migration.sql` | 旧版客户权限表迁移到 `customer_workspace_*` |

## 全新环境执行顺序

1. `sql/baize2022-01-08.sql`
2. `sql/customer_management.sql`
3. `sql/freight_shipment.sql`
4. `sql/ifs_agent.sql`
5. `sql/shipment_notification.sql`
6. `sql/system_notification_menu.sql`

## 已有环境升级建议

### 1. 消息通知能力

如果环境中还没有消息通知能力，至少执行：

1. `sql/shipment_notification.sql`
2. `sql/system_notification_menu.sql`

作用：

- 创建通知表
- 增加后台消息通知菜单
- 增加消息通知相关权限点

### 2. 客户工作台旧权限迁移

如果仍在使用旧版 `customer_portal_*` 表结构，再执行：

- `sql/customer_workspace_permission_migration.sql`

### 3. 出货与 Agent 已有环境

如果现场已经存在出货计划和 Agent 数据：

- 不要直接重跑全量脚本
- 按现场差异拆成增量 SQL
- 先核对表、索引、菜单、权限是否已存在

## 升级原则

1. 不要在生产库直接重跑带 `DROP TABLE` 的全量脚本。
2. 新增表结构时，优先提供独立增量脚本。
3. 新增菜单和权限时，同时补菜单 SQL 和文档说明。
4. 业务上线后，文档里的 SQL 清单必须同步更新。

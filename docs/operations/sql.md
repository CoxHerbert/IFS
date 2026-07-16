# SQL 说明

SQL 初始化、执行顺序和维护建议已经统一整理到：

- [初始化与升级](setup.md)

当前推荐规则：

1. 新环境执行 `sql/ifs_init.sql`
2. 业务模块统一维护到 `sql/ifs_business.sql`
3. 历史拆分迁移脚本不再作为默认执行入口
4. 如果测试环境允许重建表，优先整理合并后的全量脚本，而不是继续叠加零散迁移文件
## SQL Consolidation

Current SQL entrypoints:
- `sql/baize2022-01-08.sql`: base admin framework schema and seed data.
- `sql/ifs_business.sql`: all IFS business modules, including portal, customer, freight, CMS, notification and Agent.
- `sql/ifs_init.sql`: unified entrypoint for new environments.

Scattered dated module scripts have been merged into `sql/ifs_business.sql`.

Merged modules include:
- Freight receipt and payment declaration.
- CMS article management.
- Agent runtime configuration and `Agent 配置` menu.

Do not add new dated module SQL files by default. Add new module DDL, menu data and permission data into the corresponding section of `sql/ifs_business.sql`.

## Business Module Sections

`sql/ifs_business.sql` is the only IFS business SQL file. It is organized by module:

- Portal: public website contact and article read-side support.
- Customer: customer profiles, customer accounts, workspace menus and roles.
- Freight: shipment plans, receipts and payment declarations.
- Agent: chat tables, form submissions, runtime config, Agent menus and permissions.
- Notification: notification table and backend notification menu.
- CMS: article table, CMS menus and `cms:article:*` permissions.

Deleted scattered scripts:
- `sql/20260710_freight_receipt.sql`
- `sql/20260711_cms_article.sql`
- `sql/20260716_agent_runtime_config.sql`

# IFS 文档总览

当前文档按以下几类维护，避免业务说明、初始化 SQL 和自动生成接口文件混在一起。

## 1. 模块文档

用于说明当前已经落地的业务能力、入口页面、核心接口和关键数据表。

- [模块索引](modules/README.md)
- [门户站点](modules/portal-site.md)
- [客户工作台](modules/customer-workspace.md)
- [系统基础管理](modules/system-management.md)
- [出货计划与出货查询](modules/freight-shipment.md)
- [消息通知](modules/notification-center.md)
- [IFS Agent](modules/ifs-agent.md)

## 2. 需求与方案

用于记录仍在设计中的规则、约束和后续扩展方向。

- [智能出货工具需求说明](modules/shipment-tool-requirements.md)

## 3. 运维与部署

用于记录初始化、部署、升级和 SQL 执行说明。

- [初始化与升级](operations/setup.md)
- [SQL 说明](operations/sql.md)

当前 SQL 统一入口为：

1. `sql/ifs_init.sql`
2. `sql/ifs_business.sql`

如果直接执行统一入口脚本，则使用：

- `sql/ifs_init.sql`

## 4. 自动生成接口文件

以下文件由 Swagger 工具生成，仅用于接口查看：

- `docs.go`
- `swagger.json`
- `swagger.yaml`

## 推荐阅读顺序

1. 先看 [模块索引](modules/README.md)。
2. 再按业务看门户、客户工作台、系统管理、出货、通知和 Agent 模块文档。
3. 涉及数据库初始化或升级时，查看 [初始化与升级](operations/setup.md)。

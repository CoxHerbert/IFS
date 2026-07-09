# IFS 文档总览

当前文档按“产品模块、需求方案、运维部署、接口生成文件”四类维护，避免业务说明、SQL 操作和 Swagger 文件混在一起。

## 1. 产品模块

适合了解系统当前已经落地的页面、能力、接口和数据表。

- [模块文档索引](modules/README.md)
- [门户站点](modules/portal-site.md)
- [客户工作台](modules/customer-workspace.md)
- [出货计划与出货查询](modules/freight-shipment.md)
- [消息通知](modules/notification-center.md)
- [IFS Agent](modules/ifs-agent.md)

## 2. 需求与方案

适合查看需求背景、规则约束和后续扩展方向。

- [智能出货工具需求说明](modules/shipment-tool-requirements.md)

## 3. 运维与部署

适合初始化环境、执行 SQL、排查部署差异和升级脚本。

- [运维文档索引](operations/README.md)
- [SQL 初始化与升级顺序](operations/sql.md)

## 4. 接口生成文件

以下文件由 Swagger 工具生成，只用于接口查看，不作为业务主文档维护。

- `docs.go`
- `swagger.json`
- `swagger.yaml`

## 推荐阅读顺序

1. 先看 [模块文档索引](modules/README.md)，建立系统全貌。
2. 再看 [门户站点](modules/portal-site.md)、[客户工作台](modules/customer-workspace.md)、[出货计划与出货查询](modules/freight-shipment.md)、[消息通知](modules/notification-center.md)。
3. 如果涉及 Agent 能力，再看 [IFS Agent](modules/ifs-agent.md)。
4. 如果涉及规则设计和后续扩展，再看 [智能出货工具需求说明](modules/shipment-tool-requirements.md)。
5. 部署或升级数据库时，最后看 [SQL 初始化与升级顺序](operations/sql.md)。

# IFS 文档分类说明

为避免业务文档、需求方案、运维脚本说明和 Swagger 生成文件混杂，当前文档按下面几类维护。

## 1. 模块文档

模块文档描述已经上线的功能，重点关注：

- 页面入口
- 核心接口
- 关键数据表
- 菜单与权限

入口：

- [模块索引](modules/README.md)

当前模块：

- [客户工作台](modules/customer-workspace.md)
- [门户站点](modules/portal-site.md)
- [出货计划与出货查询](modules/freight-shipment.md)
- [消息通知](modules/notification-center.md)
- [IFS Agent](modules/ifs-agent.md)

## 2. 需求与方案文档

用于记录仍在设计中的规则和扩展方向，不等同于当前已上线功能。

- [智能出货工具需求说明](modules/shipment-tool-requirements.md)

## 3. 运维与部署文档

用于记录初始化、部署、升级和 SQL 执行说明。

- [初始化与升级](operations/setup.md)

当前 SQL 结构已经收敛为：

- 统一入口：`sql/ifs_init.sql`
- 业务合并脚本：`sql/ifs_business.sql`
- 基础脚本：`sql/baize2022-01-08.sql`

## 4. 自动生成接口文件

以下文件为工具生成，不作为主文档维护：

- `docs.go`
- `swagger.json`
- `swagger.yaml`

## 维护原则

1. 新功能上线时，优先补对应模块文档。
2. 仍在设计阶段的内容写入需求与方案文档。
3. 数据库、菜单、初始化、升级相关内容统一写入运维文档。
4. SQL 入口或执行方式调整后，需要同步更新运维文档和相关模块文档。

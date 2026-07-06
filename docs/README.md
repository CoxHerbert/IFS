# IFS 文档总览

本文档目录按“模块说明、需求设计、运维说明、接口生成文件”分类维护，避免同一主题分散在多个文件重复描述。

## 1. 模块文档

适合了解系统当前已经落地的业务能力、入口页面、数据表和接口范围。

- [模块文档索引](modules/README.md)
- [客户工作台](modules/customer-workspace.md)
- [出货计划与出货查询](modules/freight-shipment.md)
- [IFS Agent](modules/ifs-agent.md)

## 2. 需求与方案文档

适合查看需求背景、规则约束和后续扩展方向。

- [智能出货工具需求说明](modules/shipment-tool-requirements.md)

## 3. 运维文档

适合初始化环境、执行 SQL、排查部署差异。

- [运维文档索引](operations/README.md)
- [SQL 初始化与升级顺序](operations/sql.md)

## 4. 自动生成接口文档

以下文件由 Swagger 工具生成，用于接口查看，不作为业务规则主文档维护。

- `docs.go`
- `swagger.json`
- `swagger.yaml`

## 5. 推荐阅读顺序

新同学或新环境接手时，建议按下面顺序阅读：

1. 先看 [模块文档索引](modules/README.md)，建立系统全貌。
2. 再看 [客户工作台](modules/customer-workspace.md)、[出货计划与出货查询](modules/freight-shipment.md)、[IFS Agent](modules/ifs-agent.md) 这三份核心模块文档。
3. 如果涉及智能出货规划或后续扩展，再看 [智能出货工具需求说明](modules/shipment-tool-requirements.md)。
4. 部署或升级数据库时，最后看 [SQL 初始化与升级顺序](operations/sql.md)。

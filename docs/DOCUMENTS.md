# IFS 文档分类说明

为避免业务说明、需求方案、运维步骤和 Swagger 生成文件混在一起，当前文档按下面分类维护。

## 1. 模块文档

描述已经落地的业务模块，重点关注页面入口、核心接口、数据表和权限范围。

- [模块文档索引](modules/README.md)
- [客户工作台](modules/customer-workspace.md)
- [出货计划与出货查询](modules/freight-shipment.md)
- [IFS Agent](modules/ifs-agent.md)

## 2. 需求与方案文档

描述需求背景、规则约束、扩展方向，不直接等同于当前已实现功能。

- [智能出货工具需求说明](modules/shipment-tool-requirements.md)

## 3. 运维文档

描述初始化、升级、部署和脚本执行顺序。

- [运维文档索引](operations/README.md)
- [SQL 初始化与升级顺序](operations/sql.md)

## 4. 接口生成文件

由工具自动生成，只用于接口查看，不作为业务规则主文档。

- `docs.go`
- `swagger.json`
- `swagger.yaml`

## 5. 维护规则

1. 新增业务能力时，优先补对应模块文档。
2. 如果是尚未完全落地的规则，写入需求与方案文档。
3. 如果是数据库、部署、升级相关内容，写入运维文档。
4. 不要把业务规则直接写进 Swagger 生成文件。

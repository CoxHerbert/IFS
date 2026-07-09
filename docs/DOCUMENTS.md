# IFS 文档分类说明

为避免业务说明、需求方案、运维步骤和 Swagger 生成文件混在一起，当前文档按下面几类整理。

## 1. 产品模块文档

描述已经落地的模块，重点关注：

- 页面入口
- 核心接口
- 关键数据表
- 权限与菜单

入口：

- [模块文档索引](modules/README.md)

当前模块：

- [客户工作台](modules/customer-workspace.md)
- [门户站点](modules/portal-site.md)
- [出货计划与出货查询](modules/freight-shipment.md)
- [消息通知](modules/notification-center.md)
- [IFS Agent](modules/ifs-agent.md)

## 2. 需求与方案文档

描述需求背景、规则约束和扩展方向，不等同于当前已经上线的功能。

当前文档：

- [智能出货工具需求说明](modules/shipment-tool-requirements.md)

## 3. 运维与部署文档

描述初始化、部署、升级和脚本执行顺序。

入口：

- [运维文档索引](operations/README.md)
- [SQL 初始化与升级顺序](operations/sql.md)

## 4. 接口生成文件

由工具自动生成，只用于接口查看：

- `docs.go`
- `swagger.json`
- `swagger.yaml`

## 维护规则

1. 新增上线能力时，优先补对应模块文档。
2. 仍在设计中的规则，写入需求与方案文档。
3. 数据库、菜单、部署、初始化相关内容，写入运维文档。
4. 不要把业务规则直接写进 Swagger 生成文件。

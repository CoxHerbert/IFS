# 初始化与升级

## 统一入口

- SQL 统一入口：`sql/ifs_init.sql`
- 业务合并脚本：`sql/ifs_business.sql`
- 运维说明入口：`docs/operations/setup.md`

`sql/ifs_init.sql` 当前只串联两份脚本：

1. `sql/baize2022-01-08.sql`
2. `sql/ifs_business.sql`

## 脚本清单

| 脚本 | 用途 |
| --- | --- |
| `sql/baize2022-01-08.sql` | 基础系统表、基础菜单、角色、字典 |
| `sql/ifs_business.sql` | 业务模块合并脚本，包含客户、出货、Agent、通知 |

## 新环境

优先直接执行：

- `sql/ifs_init.sql`

如果当前 SQL 客户端不支持 `SOURCE`，则手工按下面顺序执行：

1. `sql/baize2022-01-08.sql`
2. `sql/ifs_business.sql`

## 已有环境升级

不要在生产库直接重跑包含 `DROP TABLE` 的全量脚本。

当前仓库默认按“可重建环境”维护：

- 新环境或测试环境：直接执行 `sql/ifs_init.sql`
- 已存在正式业务数据的环境：按现场差异单独整理增量 SQL，不要直接重跑 `sql/ifs_business.sql`

## 维护原则

1. 新增上线模块时，优先把表结构、菜单、权限脚本合并进 `sql/ifs_business.sql`
2. `sql/ifs_init.sql` 保持为统一入口，不再继续堆叠过多 `SOURCE` 明细
3. 如果脚本存在固定主键冲突风险，优先改成动态 ID 或 `NOT EXISTS` 风格
4. SQL 清单和执行建议变更后，需要同步更新本页文档

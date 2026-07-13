# 初始化与升级

## 统一入口

- SQL 统一入口：`sql/ifs_init.sql`
- 业务合并脚本：`sql/ifs_business.sql`
- 运维说明入口：`docs/operations/setup.md`

`sql/ifs_init.sql` 当前串联：

1. `sql/baize2022-01-08.sql`
2. `sql/ifs_business.sql`

## 脚本清单

| 脚本 | 用途 |
| --- | --- |
| `sql/baize2022-01-08.sql` | 基础系统表、基础菜单、角色、字典 |
| `sql/ifs_business.sql` | 业务模块合并脚本，包含客户、出货、Agent、通知等 |
| `sql/20260711_cms_article.sql` | CMS 新闻资讯模块，包含 `cms_article` 表、CMS 菜单和 `cms:article:*` 权限 |

## 新环境

优先执行：

```sql
SOURCE sql/ifs_init.sql;
```

如果当前 SQL 客户端不支持 `SOURCE`，则手工按下面顺序执行：

1. `sql/baize2022-01-08.sql`
2. `sql/ifs_business.sql`
3. `sql/20260711_cms_article.sql`

## 已有环境升级

不要在生产库直接重跑包含 `DROP TABLE` 的全量脚本。

当前仓库默认按“可重建环境”维护：

- 新环境或测试环境：直接执行 `sql/ifs_init.sql`，再执行增量脚本。
- 已存在正式业务数据的环境：按现场差异单独整理增量 SQL，不要直接重跑 `sql/ifs_business.sql`。

CMS 模块升级：

```sql
SOURCE sql/20260711_cms_article.sql;
```

该脚本使用 `IF NOT EXISTS` 和 `UPDATE` 兼容已执行过旧菜单版本的环境。

## 静态资源

CMS 富文本图片上传后会保存为 `/profile/cms/article/...`。

部署要求：

- 后端需要暴露 `constants.ResourcePrefix` 对应的静态目录。
- 开发环境中，`baize-ui` 和 `portal-ui` 的 Vite 配置需要代理 `/profile` 到后端。
- 生产环境中，Nginx 或网关需要暴露 `/profile` 静态资源。

安全建议：

- 禁止 `/profile` 目录列表。
- 设置 `X-Content-Type-Options: nosniff`。
- 对 `/profile/cms/article/` 限制只返回图片 MIME。
- 不要把敏感文件放入 `/profile`。

## 维护原则

1. 新增上线模块时，优先把表结构、菜单、权限脚本合并进业务 SQL 或提供明确增量脚本。
2. `sql/ifs_init.sql` 保持为统一入口，不继续堆叠过多 `SOURCE` 明细。
3. 如果脚本存在固定主键冲突风险，优先改成 `NOT EXISTS` 或可重复执行风格。
4. SQL 清单和执行建议变更后，需要同步更新本文档。

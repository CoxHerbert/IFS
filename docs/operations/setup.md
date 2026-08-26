# 初始化与升级

## 新环境初始化

新数据库执行：

```sql
SOURCE sql/ifs_init.sql;
```

`ifs_init.sql` 按顺序加载：

1. `sql/baize2022-01-08.sql`
2. `sql/ifs_business.sql`

如果数据库客户端不支持 `SOURCE`，按上述顺序手工执行两个文件。

> `ifs_business.sql` 包含 `DROP TABLE` 等全量重建语句，只能用于新环境或明确允许重建的测试环境。

## 已有环境升级

已有业务数据的数据库只执行：

```sql
SOURCE sql/ifs_upgrade.sql;
```

升级脚本按业务顺序处理：

1. 出货计划：单件重量/体积、币种、贸易条款、运输范围、地址、船名航次、柜型容量。
2. 收付款：应付金额、费用明细、独立一级菜单、客户端付款明细菜单及角色授权。
3. Agent：独立一级菜单、历史菜单 ID 冲突修复、管理员及已有角色授权补齐。

脚本使用字段存在性判断、`CREATE TABLE IF NOT EXISTS`、`INSERT IGNORE` 和定向 `UPDATE`，可以重复执行。执行前仍建议备份生产数据库。

## 发布检查

执行升级后检查：

```sql
SHOW COLUMNS FROM freight_shipment_plan;
SHOW COLUMNS FROM freight_shipment_cargo;
SELECT menu_id, menu_name, parent_id, path, component FROM sys_menu WHERE menu_id IN (142,143,144,145,146,147,148);
SELECT dict_value, remark FROM sys_dict_data WHERE dict_type = 'freight_container_type' ORDER BY dict_sort;
```

随后重启后端，使用 admin 账号重新登录，确认：

- “收付款管理”和“Agent 管理”作为独立一级菜单出现。
- 出货计划可维护贸易条款、运输范围、船名航次、币种及货物尺寸重量。
- 客户端工作台显示“付款明细”。

## 静态资源

CMS 富文本图片保存于 `/profile/cms/article/`。生产环境应禁止目录列表、设置 `X-Content-Type-Options: nosniff`，并限制该目录仅返回允许的图片 MIME 类型。

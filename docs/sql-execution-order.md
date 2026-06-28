# SQL 脚本执行顺序

本文档用于说明 `sql/` 目录下脚本的推荐执行顺序。执行前请先备份数据库，尤其是带有 `DROP TABLE`、`DELETE FROM sys_menu` 的脚本。

## 全新数据库初始化

全新库没有业务数据时，按下面顺序执行：

1. `sql/baize2022-01-08.sql`

   初始化后台管理系统基础表、基础菜单、角色、字典等数据。

2. `sql/portal_contact.sql`

   创建官网线索表 `portal_contact`，并写入官网线索相关菜单。

3. `sql/portal_contact_status_dict.sql`

   写入官网线索状态字典 `portal_contact_status`。

4. `sql/customer_management.sql`

   创建客户主体表 `customer`、客户联系人表 `customer_contact`、客户账号表 `customer_account`，并写入“客户资料”“客户账号”后台菜单。

## 已有数据库升级

已有库通常已经执行过基础脚本，不建议重新执行 `baize2022-01-08.sql`。

推荐顺序：

1. 确认基础后台表已存在

   至少需要存在这些系统表：

   - `sys_menu`
   - `sys_role_menu`
   - `sys_dict_type`
   - `sys_dict_data`

2. 如还没有官网线索表，执行：

   `sql/portal_contact.sql`

   注意：该脚本包含 `DROP TABLE IF EXISTS portal_contact`，如果已有线索数据，不要直接执行。

3. 如已有旧版官网线索表，只需要升级状态字段，执行：

   `sql/portal_contact_status_migration.sql`

4. 执行官网线索字典：

   `sql/portal_contact_status_dict.sql`

   该脚本会先删除并重建 `portal_contact_status` 字典。

5. 如需要单独重建官网线索菜单，执行：

   `sql/portal_contact_menu.sql`

   注意：该脚本会删除相关菜单和角色菜单关系后重建。

6. 执行客户管理脚本：

   `sql/customer_management.sql`

   注意：该脚本会 `DROP TABLE IF EXISTS customer_account`、`DROP TABLE IF EXISTS customer_contact` 和 `DROP TABLE IF EXISTS customer`。如果已有客户数据，不能直接执行，需要改成增量迁移。

## 当前客户管理功能依赖

客户管理功能依赖以下数据库对象：

- `customer`
- `customer_contact`
- `customer_account`
- `sys_menu`
- `sys_role_menu`

门户客户登录依赖：

- `customer_account.username`
- `customer_account.password`
- `customer_account.status = '0'`

后台菜单权限依赖：

- `customer:customer:list`
- `customer:customer:query`
- `customer:customer:add`
- `customer:customer:edit`
- `customer:customer:remove`
- `customer:account:list`
- `customer:account:query`
- `customer:account:add`
- `customer:account:edit`
- `customer:account:remove`
- `customer:account:resetPwd`

## 建议

生产环境不要直接执行带 `DROP TABLE` 的脚本。客户管理上线后，如果后续已有客户数据，应新增类似 `customer_management_migration.sql` 的增量脚本，只做 `ALTER TABLE`、补菜单、补索引，不再重建表。

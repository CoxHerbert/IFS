# SQL 执行顺序

`sql/` 目录已收敛为按功能划分的主脚本，默认只需要关注 3 个初始化脚本和 1 个迁移脚本。

## 当前保留脚本

1. `sql/baize2022-01-08.sql`
   基础系统表、基础菜单、角色、字典。

2. `sql/customer_management.sql`
   包含以下内容：
   - 官网线索表 `portal_contact`
   - 官网线索状态字典 `portal_contact_status`
   - 客户表 `customer`
   - 客户联系人表 `customer_contact`
   - 客户账号表 `customer_account`
   - 客户端权限表 `customer_workspace_*`
   - 后台“客户管理”目录及其下全部菜单、按钮权限

3. `sql/freight_shipment.sql`
   包含以下内容：
   - 出货计划相关 4 张表
   - 出货状态字典 `freight_shipment_status`
   - 货柜类型字典 `freight_container_type`
   - 后台“货代业务 -> 出货计划”菜单及按钮权限

4. `sql/customer_workspace_permission_migration.sql`
   仅用于把旧版 `customer_portal_*` 表迁移为 `customer_workspace_*` 表。

## 全新数据库初始化

如果是全新环境且没有业务数据，建议按以下顺序执行：

1. `sql/baize2022-01-08.sql`
2. `sql/customer_management.sql`
3. `sql/freight_shipment.sql`

## 已有数据库升级

已有环境不要重新执行带 `DROP TABLE` 的全量脚本。

推荐处理方式：

1. 如果系统库还没初始化，先执行 `sql/baize2022-01-08.sql`
2. 如果现场仍是旧版客户端权限表，先执行 `sql/customer_workspace_permission_migration.sql`
3. 其他结构变更请新增增量迁移脚本，不要直接重跑 `customer_management.sql` 或 `freight_shipment.sql`

## 当前功能依赖

客户与客户端功能依赖：

- `portal_contact`
- `customer`
- `customer_contact`
- `customer_account`
- `customer_workspace_menu`
- `customer_workspace_role`
- `customer_workspace_role_menu`
- `customer_workspace_account_role`
- `sys_menu`
- `sys_role_menu`
- `sys_dict_type`
- `sys_dict_data`

出货计划功能依赖：

- `freight_shipment_plan`
- `freight_shipment_cargo`
- `freight_container_plan`
- `freight_shipment_order`
- `sys_menu`
- `sys_role_menu`
- `sys_dict_type`
- `sys_dict_data`

## 建议

生产环境不要直接执行带 `DROP TABLE` 的初始化脚本。
如果线上已有数据，统一走增量迁移脚本；初始化脚本只用于新库、测试库或一次性重建场景。

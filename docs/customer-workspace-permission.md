# 客户端菜单与角色说明

## 功能范围

后台管理系统新增两块客户端权限配置：

- 客户端菜单管理
- 客户端角色管理

同时在“客户账号”页面支持给客户账号分配客户端角色。客户端登录后不再写死侧边菜单，而是读取后台配置动态加载。

## 数据表

初始化脚本：`sql/customer_management.sql`

其中包含以下客户端权限表：

- `customer_workspace_menu`
- `customer_workspace_role`
- `customer_workspace_role_menu`
- `customer_workspace_account_role`

迁移脚本：`sql/customer_workspace_permission_migration.sql`

适用于已经上线过旧版 `customer_portal_*` 表的环境。

## 后台入口

执行初始化脚本后，会在“客户管理”下新增：

- `客户端菜单`
- `客户端角色`

已有“客户账号”页面会增加角色分配能力。

## 客户端加载链路

1. 客户账号通过 `/portal/customer/login` 登录。
2. 登录后客户端调用 `/portal/customer/profile` 获取当前账号、角色、权限。
3. 客户端调用 `/portal/customer/routers` 获取当前账号可见菜单。
4. `portal-ui` 根据返回菜单动态注册路由并渲染侧边栏。

## 当前支持的客户端页面组件

当前后台菜单组件标识仅支持以下值：

- `workspace/dashboard`
- `workspace/account-profile`
- `workspace/shipment-tracking`

对应页面文件：

- `portal-ui/src/views/workspace/WorkspaceDashboardView.vue`
- `portal-ui/src/views/workspace/WorkspaceAccountProfileView.vue`
- `portal-ui/src/views/workspace/WorkspaceShipmentTrackingView.vue`

`/portal/customer/routers` 已支持树状菜单，父级目录可以继续挂子菜单。

## 默认初始化数据

初始化脚本会写入：

- 3 个客户端菜单：工作台、账号资料、出货查询
- 1 个默认客户端角色：`基础客户端角色`

注意：脚本不会自动把角色分配给客户账号，需要在后台“客户账号”里手工分配。

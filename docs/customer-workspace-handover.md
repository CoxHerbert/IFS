# 客户端工作台改造交接文档

## 1. 目标概述

本轮改造的核心目标是把“门户”和“客户端工作台”拆开，并让客户端工作台根据后台配置动态加载菜单和路由。

同时，后台管理系统补齐客户端菜单、客户端角色、客户账号角色分配能力，前后端、数据库脚本、文档同步收口。

---

## 2. 已完成的业务改动

### 2.1 客户端工作台路由拆分

`portal-ui` 的路由已按职责拆分：

- `src/router/modules/portal.ts`
  - 门户公开页面
  - 客户登录页
- `src/router/modules/workspace.ts`
  - 客户端工作台基础壳路由 `/customer`
  - 保留静态兜底子路由：
    - `/customer/workspace`
    - `/customer/account`
    - `/customer/shipment`
- `src/router/workspace-runtime.ts`
  - 根据 `/portal/customer/routers` 返回结果动态注册工作台子路由
  - 支持树状菜单递归注册
- `src/router/index.ts`
  - 只负责组装路由和前置守卫

当前设计：

- 门户路由与工作台路由已解耦
- 工作台业务页面路由以接口返回为主
- 静态兜底路由只用于避免刷新白屏

### 2.2 客户端动态菜单与动态路由

后端已提供：

- `/portal/customer/profile`
- `/portal/customer/routers`

前端工作方式：

1. 登录成功后保存 token
2. 进入 `/customer` 时先校验 token
3. 优先从本地缓存恢复工作台路由
4. 再调用 `/portal/customer/routers` 获取最新路由并覆盖
5. 工作台侧边栏按返回菜单渲染

当前支持的组件标识：

- `workspace/dashboard`
- `workspace/account-profile`
- `workspace/shipment-tracking`

### 2.3 客户端登录态持久化

已完成持久化：

- token：`localStorage`
- profile：`localStorage`
- workspace routes：`localStorage`

对应文件：

- `portal-ui/src/api/workspace/auth.ts`

相关 key：

- `portal_customer_token`
- `portal_customer_profile`
- `portal_customer_routes`

之前刷新白屏的原因就是只有 token 做了持久化，`profile` 和动态路由缓存没有完整恢复。现在已修复。

### 2.4 工作台布局拆分

`WorkspaceShellLayout.vue` 已拆分，避免一个文件同时承载模板、菜单逻辑、账号逻辑、样式。

当前结构：

- `portal-ui/src/layouts/workspace/WorkspaceShellLayout.vue`
- `portal-ui/src/layouts/workspace/components/WorkspaceHeader.vue`
- `portal-ui/src/layouts/workspace/components/WorkspaceSidebar.vue`
- `portal-ui/src/layouts/workspace/useWorkspaceMenu.ts`

### 2.5 工作台样式调整

已完成：

- 面包屑和折叠按钮放同一栏
- 右上角账号入口收紧，不再占高
- “返回门户”已移入账号下拉菜单
- 下拉入口只展示账号，不再展示客户名称
- 左侧菜单折叠后不显示标题文字
- 客户端首页内容区改为更简单的白底黑字风格

涉及文件：

- `portal-ui/src/layouts/workspace/components/WorkspaceHeader.vue`
- `portal-ui/src/layouts/workspace/components/WorkspaceSidebar.vue`
- `portal-ui/src/views/workspace/WorkspaceDashboardView.vue`

### 2.6 后台管理系统补齐客户端权限配置

后台管理系统已新增客户端权限配置入口：

- 客户端菜单管理
- 客户端角色管理
- 客户账号角色分配

前端页面：

- `baize-ui/src/views/customer/portalMenu/index.vue`
- `baize-ui/src/views/customer/portalRole/index.vue`
- `baize-ui/src/views/customer/account/index.vue`

说明：

- `portalMenu` 页面乱码已修复
- `portalRole` 页面仍建议继续检查是否存在乱码或样式不统一问题

### 2.7 后端客户端权限链路

后端已补齐：

- 客户账号登录
- 当前客户账号 profile
- 当前客户账号 routers
- 客户端菜单 CRUD
- 客户端角色 CRUD
- 账号角色分配

核心文件：

- `app/customer/controller/customerController.go`
- `app/customer/service/customerService.go`
- `app/customer/dao/customerDao.go`
- `app/customer/models/customer.go`
- `app/customer/models/customerPortal.go`
- `app/routes/customerRoutes/customerRouter.go`

### 2.8 `/portal/customer/profile` 已修问题

#### 问题 1：token 过期太短

原因：

- 客户端 token 过期时间之前按“秒”处理
- 系统配置 `token.expireTime` 实际语义是“分钟”

修复：

- `app/customer/service/customerService.go`
  - 客户端 token 过期时间改为按分钟处理

#### 问题 2：全局 token 时长太短

修复：

- `config.yaml`
  - `token.expireTime` 从 `30` 调整为 `720`

当前效果：

- 管理端 token：12 小时
- 客户端 token：12 小时

#### 问题 3：`/portal/customer/profile` 报 500

实际原因：

- MySQL 8 下 `DISTINCT + ORDER BY 非 SELECT 字段` 报错
- 触发位置：`SelectPortalPermissionsByAccountId`

修复：

- `app/customer/dao/customerDao.go`
  - 原 SQL：`order by m.menu_id`
  - 改为：`order by m.perms`

另外还补了 `BaiZeTime` 的数据库扫描支持：

- `app/common/baize/baizeUnix/UnixTime.go`

这不是本次 500 的直接根因，但属于真实缺口，已顺手补齐。

---

## 3. 数据库脚本整理结果

`sql/` 目录已收敛到以下脚本：

- `sql/baize2022-01-08.sql`
- `sql/customer_management.sql`
- `sql/freight_shipment.sql`
- `sql/customer_workspace_permission_migration.sql`

### 3.1 `customer_management.sql`

已合并内容：

- 官网线索表 `portal_contact`
- 官网线索状态字典 `portal_contact_status`
- 客户表 `customer`
- 客户联系人表 `customer_contact`
- 客户账号表 `customer_account`
- 客户端权限表：
  - `customer_workspace_menu`
  - `customer_workspace_role`
  - `customer_workspace_role_menu`
  - `customer_workspace_account_role`
- 后台“客户管理”目录及相关菜单/按钮权限

### 3.2 `freight_shipment.sql`

已合并内容：

- 出货计划相关 4 张表
- 出货状态字典
- 货柜类型字典
- “货代业务 -> 出货计划”菜单及按钮权限

### 3.3 已删除的旧脚本

已删除：

- `sql/portal_contact.sql`
- `sql/portal_contact_menu.sql`
- `sql/portal_contact_status_dict.sql`
- `sql/portal_contact_status_migration.sql`
- `sql/customer_workspace_permission.sql`
- `sql/freight_shipment_menu.sql`

---

## 4. 当前客户端页面与命名

已按“门户 / 工作台”重构：

### 门户页面

- `portal-ui/src/views/portal/PortalHomeView.vue`
- `portal-ui/src/views/portal/PortalNewsView.vue`
- `portal-ui/src/views/portal/PortalServiceView.vue`
- `portal-ui/src/views/portal/PortalAboutView.vue`
- `portal-ui/src/views/portal/PortalContactView.vue`
- `portal-ui/src/views/portal/PortalShipmentShareView.vue`

### 工作台页面

- `portal-ui/src/views/workspace/WorkspaceLoginView.vue`
- `portal-ui/src/views/workspace/WorkspaceDashboardView.vue`
- `portal-ui/src/views/workspace/WorkspaceAccountProfileView.vue`
- `portal-ui/src/views/workspace/WorkspaceShipmentTrackingView.vue`

### API

- `portal-ui/src/api/portal/contact.ts`
- `portal-ui/src/api/portal/shipment.ts`
- `portal-ui/src/api/workspace/auth.ts`

---

## 5. 已更新文档

已同步更新：

- `docs/sql-execution-order.md`
- `docs/customer-workspace-permission.md`
- `docs/freight-shipment-feature.md`

---

## 6. 当前已验证结果

已验证：

- `portal-ui` `npm run type-check` 通过
- `portal-ui` `npm run build` 通过

未完全验证：

- `baize-ui` 未完整做一次构建收口
- Go 后端未在当前环境完成完整编译校验（本机工具链受限时需要注意）

---

## 7. 下次可继续做的事项

建议优先级：

1. 继续清理 `baize-ui` 中剩余乱码页面
   - 特别是 `portalRole` 和客户管理相关页面

2. 把工作台会话逻辑继续抽离
   - 可新增 `useWorkspaceSession.ts`
   - 把 profile 拉取、logout、缓存恢复从布局层再拆出去

3. 扩展工作台支持的动态页面组件映射
   - 当前只支持 3 个组件标识
   - 若要增加更多页面，需要同步：
     - 后端菜单配置
     - 前端 `workspaceComponentMap`

4. 为线上环境补增量迁移脚本
   - 尤其是客户管理和出货计划，不建议直接执行带 `DROP TABLE` 的初始化脚本

5. 继续统一客户端视觉风格
   - 登录页
   - 工作台头部
   - 指标卡与按钮风格

---

## 8. 关键注意事项

- `/portal/customer/routers` 已支持树状菜单，但前端组件映射仍是白名单机制
- `/customer/workspace` 白屏问题已经通过“静态兜底 + 动态恢复 + 本地缓存”解决
- token 时长现在是全局 12 小时，修改 `config.yaml` 后需要重启服务
- 如果 `/portal/customer/profile` 再报 500，优先看：
  - 角色菜单关联数据是否完整
  - SQL 是否又引入了 MySQL 8 的 `DISTINCT/ORDER BY` 兼容问题


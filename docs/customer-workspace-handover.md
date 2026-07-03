# 客户端工作台改造交接

## 1. 当前目标

本轮改造把门户和客户端工作台彻底拆开，并让客户端工作台根据后台菜单配置动态加载。

同时补齐了三部分能力：

- 后台管理端可维护客户端菜单、客户端角色、客户账号角色分配
- 客户端工作台支持动态路由、本地持久化和刷新恢复
- 新增“智能出货助手”，用于 Excel 整理、整柜测算和散货体积评估

## 2. 路由与登录

### 门户和工作台拆分

`portal-ui` 路由已经拆分为：

- `src/router/modules/portal.ts`
- `src/router/modules/workspace.ts`
- `src/router/workspace-runtime.ts`
- `src/router/index.ts`

其中：

- 门户页面走 `portal`
- 客户端工作台走 `/customer`
- `/customer` 下的业务页优先由接口动态注册
- 本地保留静态兜底路由，避免刷新白屏

### 客户端持久化

客户端登录态已做本地缓存，位置：

- `portal-ui/src/api/workspace/auth.ts`

使用的 key：

- `portal_customer_token`
- `portal_customer_profile`
- `portal_customer_routes`

当前刷新恢复链路：

1. 先校验本地 token
2. 先恢复本地 profile 和 routes
3. 再调用 `/portal/customer/routers` 覆盖本地缓存
4. 若访问 `/customer`，自动跳到第一个可用业务页

## 3. 工作台页面结构

### 当前页面

- `portal-ui/src/views/workspace/WorkspaceLoginView.vue`
- `portal-ui/src/views/workspace/WorkspaceDashboardView.vue`
- `portal-ui/src/views/workspace/WorkspaceAccountProfileView.vue`
- `portal-ui/src/views/workspace/WorkspaceShipmentTrackingView.vue`
- `portal-ui/src/views/workspace/WorkspaceShipmentAssistantView.vue`

### 当前支持的动态组件标识

- `workspace/dashboard`
- `workspace/account-profile`
- `workspace/shipment-tracking`
- `workspace/shipment-assistant`

如果后台菜单要新增客户端页面，必须同步补前端组件映射，否则接口返回了菜单，前端也不会注册对应路由。

## 4. 智能出货助手

### 业务目标

这个页面是给客户自己先整理出货 Excel 用的，降低反复手工统计工作。

当前支持：

- 导入 `.xlsx/.xls/.csv`
- 按表格在线编辑货物明细
- 自动识别常见中英文表头
- 如果没填体积，但填了长宽高和箱数，会自动换算 CBM
- 计算总数量、总箱数、总重量、总体积
- 推荐整柜柜型和柜数
- 给出散货拼箱建议

### 前端实现

依赖：

- `vxe-table`
- `xe-utils`
- `xlsx`

接入文件：

- `portal-ui/src/api/workspace/shipmentAssistant.ts`
- `portal-ui/src/views/workspace/WorkspaceShipmentAssistantView.vue`
- `portal-ui/src/main.ts`

### 后端实现

新增客户端受保护接口：

- `POST /portal/customer/shipment-assistant/estimate`

相关文件：

- `app/freight/models/shipment.go`
- `app/freight/service/shipmentService.go`
- `app/customer/controller/customerController.go`
- `app/routes/customerRoutes/customerRouter.go`

实现方式：

- 复用出货模块已有柜型能力
- 把“货物标准化 + 汇总 + 柜型推荐”抽成纯计算方法
- 客户端测算不会落库，只返回计算结果

## 5. 后台菜单和权限

后台管理端已支持客户端配置：

- 客户端菜单管理
- 客户端角色管理
- 客户账号角色配置

相关页面：

- `baize-ui/src/views/customer/portalMenu/index.vue`
- `baize-ui/src/views/customer/portalRole/index.vue`
- `baize-ui/src/views/customer/account/index.vue`

客户端工作台的菜单由后台 `customer_workspace_menu` 和角色关联决定。

本轮新增了客户端菜单项：

- `智能出货助手`
- `path = shipment-assistant`
- `component = workspace/shipment-assistant`
- `perms = portal:shipmentAssistant:view`

## 6. SQL 现状

当前主要脚本：

- `sql/baize2022-01-08.sql`
- `sql/customer_management.sql`
- `sql/freight_shipment.sql`
- `sql/customer_workspace_permission_migration.sql`

本轮已同步更新：

- `sql/customer_management.sql`
- `sql/customer_workspace_permission_migration.sql`

新增默认客户端菜单 `20004`：

- `智能出货助手`

默认基础客户端角色 `20001` 已自动关联该菜单。

## 7. 样式与布局

工作台布局已拆分：

- `portal-ui/src/layouts/workspace/WorkspaceShellLayout.vue`
- `portal-ui/src/layouts/workspace/components/WorkspaceHeader.vue`
- `portal-ui/src/layouts/workspace/components/WorkspaceSidebar.vue`
- `portal-ui/src/layouts/workspace/useWorkspaceMenu.ts`

当前视觉规则：

- 客户端右上角只展示账号
- 客户名称不再放在右上角
- 返回门户放进下拉菜单
- 菜单折叠后不展示标题文字
- 页面内容区以白底黑字为主

## 8. 验证结果

已验证：

- `portal-ui` `npm run type-check`
- `portal-ui` `npm run build`

说明：

- 构建通过
- 由于引入 `vxe-table` 与 `xlsx`，打包体积上升，已通过路由懒加载把助手页依赖拆出首屏
- 本机没有 `go` 命令，未完成后端编译验证

## 9. 后续建议

建议下一步继续做三件事：

1. 把“智能出货助手”的测算结果进一步和正式出货计划创建打通
2. 在后台菜单管理里补一份更完整的客户端菜单模板
3. 继续清理 `baize-ui` 历史乱码页面，尤其是客户端角色管理相关页

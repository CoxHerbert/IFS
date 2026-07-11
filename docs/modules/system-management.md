# 系统基础管理

## 目标

系统基础管理用于维护后台的菜单、权限与基础展示能力。当前重点是菜单图标选择的统一，避免旧 SVG 图标与 Iconify 图标混用后出现列表混乱、样式不一致和维护成本上升。

## 菜单图标规则

- 后台菜单图标选择页面：`/system/menu`
- 菜单编辑时统一使用新的 `IconPicker`
- 图标名称统一以 `mdi:*` 形式保存
- 展示层通过 `app-icon` 自动兼容：
  - `mdi:*` 走 Iconify
  - 旧 SVG 名称保留兼容显示

## 当前行为

- 菜单列表中图标列同时显示图标和名称，方便识别
- 图标选择器采用统一网格布局，支持搜索与高亮选中状态
- 系统菜单页不再使用旧的 `IconSelect` 弹层
- 客户端菜单管理页也复用同一套选择器

## 相关页面

- `baize-ui/src/views/system/menu/index.vue`
- `baize-ui/src/views/customer/portalMenu/index.vue`
- `baize-ui/src/components/IconPicker/index.vue`
- `baize-ui/src/components/AppIcon/index.vue`

## 相关说明

- 图标注册表位于 `baize-ui/src/components/AppIcon/iconRegistry.ts`
- 后台与门户统一使用离线打包的 Iconify MDI 图标
- 旧 SVG 图标仅作为历史兼容，不再作为新的默认选择来源

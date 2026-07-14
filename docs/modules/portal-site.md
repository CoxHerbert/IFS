# 门户站点

## 定位

门户站点是 IFS 的公开获客入口，承接公司展示、新闻资讯、服务能力、联系表单和悬浮智能助手。

当前定位：

- 工具型门户，不是纯宣传型官网
- 首页聚焦服务说明、新闻资讯和联系入口
- 智能助手以右下角悬浮方式存在，不占用表头导航
- 新闻资讯内容由后台 `CMS管理 -> 新闻资讯` 维护

## 页面入口

| 页面 | 路径 |
| --- | --- |
| 门户首页 | `/` |
| 新闻资讯 | `/news` |
| 新闻详情 | `/news/:slug` |
| 服务能力 | `/service` |
| 关于我们 | `/about` |
| 联系我们 | `/contact` |
| 客户中心 | `/workspace` |

说明：

- 旧的 `/shipment-agent` 独立门户模块已移除。
- 首页不再单独保留“为什么联系”“快速询价”等重复模块，相关动作收敛到“联系我们”。

## 新闻资讯

门户新闻资讯页从 CMS 公开接口读取已发布文章。

前台文件：

- `portal-ui/src/views/portal/PortalNewsView.vue`
- `portal-ui/src/api/portal/article.ts`

后台维护：

- `CMS管理 -> 新闻资讯`
- 详见 [CMS 管理](cms-management.md)

公开接口：

- `GET /portal/articles`
- `GET /portal/articles/:slug`

展示规则：

- 只展示 `status = 0` 的文章
- 支持分类、搜索、头条和详情页
- 富文本正文使用 `v-html` 渲染前必须经过 DOMPurify 清洗
- 图片资源通过 `/profile/cms/article/...` 访问

## 悬浮智能助手

组件：

- `portal-ui/src/layouts/portal/components/PortalFloatingAgent.vue`

规则：

- 门户所有公开页右下角展示悬浮助手入口
- 点击后展开聊天面板
- 支持消息发送
- 支持 Excel / CSV 文件上传分析
- 门户表头不保留“智能助手”导航项

## 页面加载态

组件：

- `portal-ui/src/layouts/portal/PortalSiteLayout.vue`

规则：

- 首次进入门户页面显示加载态
- 切换门户路由时显示短暂加载态
- 只作用于门户布局，不影响客户工作台

## 表头与导航

组件：

- `portal-ui/src/layouts/portal/components/PortalHeader.vue`

当前策略：

- 保留首页、新闻资讯、服务能力、关于我们、联系我们等主导航
- 客户中心作为登录/业务入口
- 获取报价、联系咨询等动作统一导向“联系我们”
- 智能助手由悬浮入口承接

## 开发代理

门户开发环境需要代理以下路径到后端：

- `/portal`
- `/agent-api`
- `/profile`

其中 `/profile` 用于展示 CMS 富文本图片。

配置文件：

- `portal-ui/vite.config.ts`

## 关键文件

- `portal-ui/src/views/portal/PortalHomeView.vue`
- `portal-ui/src/views/portal/PortalNewsView.vue`
- `portal-ui/src/views/portal/PortalContactView.vue`
- `portal-ui/src/layouts/portal/PortalSiteLayout.vue`
- `portal-ui/src/layouts/portal/components/PortalHeader.vue`
- `portal-ui/src/layouts/portal/components/PortalFloatingAgent.vue`

## 后续建议

- 为新闻资讯补充封面图字段在前台列表中的展示
- 为 CMS 文章增加草稿预览能力
- 为 `/profile` 静态资源增加生产环境安全响应头

## 前端接口前缀

门户工程 `portal-ui` 当前使用以下前缀：

- 门户业务接口：`VITE_PORTAL_API_PREFIX=/portal-api`
- Agent 接口：`VITE_AGENT_API_PREFIX=/agent-api`
- 静态资源：固定 `/profile`

开发代理约定：
- `/portal-api` rewrite 到后端 `/portal`
- `/agent-api` 直接代理到后端同名路由
- `/profile` 直接代理到后端静态资源目录

相关文件：
- `portal-ui/.env`
- `portal-ui/vite.config.ts`
- `portal-ui/src/utils/portal-api.ts`
- `portal-ui/src/utils/agent-api.ts`

# 门户站点

> 图标体系：门户与后台统一使用离线打包的 Iconify MDI 图标，菜单字段保存 `mdi:*` 名称；后台保留旧 SVG 图标名兼容能力。

## 目标

门户站点不再只承担公司展示作用，而是作为 `IFS` 的公开获客入口，优先承接工具使用、报价咨询和智能问答。

当前定位：

- 工具型门户，而不是纯宣传型官网
- 首页优先给出工具入口，再承接服务介绍和线索转化
- 智能助手以悬浮方式存在，不再占用表头导航

## 页面入口

| 页面 | 路径 |
| --- | --- |
| 门户首页 | `/` |
| 航线资讯 | `/news` |
| 服务能力 | `/service` |
| 出货分析工具 | `/shipment-agent` |
| 关于我们 | `/about` |
| 联系我们 | `/contact` |

## 首页信息架构

组件：`portal-ui/src/views/portal/PortalHomeView.vue`

首页当前按以下顺序组织：

1. Hero 首屏
2. 工具导航
3. 热门工具
4. 核心服务
5. 热门航线与快速询价
6. FAQ 与案例

其中“工具导航”和“热门工具”是本轮新增的核心模块。

## 工具优先入口

首页已改为“以工具为入口”的结构，工具导航模块按主题聚合真实入口，而不是只展示说明文案。

当前工具类目：

- 测算类
- 规划类
- 查询类
- AI 助手类

当前热门工具：

- 出货计划分析
- 装柜与拼箱判断
- 智能物流问答
- 快速获取报价

入口动作规则：

- 进入出货分析类工具时跳转 `/shipment-agent`
- 进入报价或资料收集类动作时跳转 `/contact`
- 进入问答类动作时直接打开悬浮智能助手

## 悬浮智能助手

组件：`portal-ui/src/layouts/portal/components/PortalFloatingAgent.vue`

门户智能助手已改为右下角悬浮入口 + 悬浮面板，不再作为表头菜单项。

当前规则：

- 门户所有公开页右下角展示悬浮助手入口
- 点击后展开聊天面板
- 支持消息发送
- 支持 Excel / CSV 文件上传分析
- 门户表头不再保留“智能助手”导航项

说明：

- 原公开页 `/agent` 路由仍保留，便于兼容旧入口
- 但门户主交互已切换为悬浮助手

## 页面加载态

组件：`portal-ui/src/layouts/portal/PortalSiteLayout.vue`

门户布局层已加入统一的页面加载遮罩。

当前规则：

- 首次进入门户页面显示“加载中”
- 切换门户路由时显示短暂加载态
- 只作用于门户布局，不影响客户工作台

## 表头与导航

组件：`portal-ui/src/layouts/portal/components/PortalHeader.vue`

当前表头策略：

- 保留首页、资讯、服务、出货分析、关于我们、联系我们等主导航
- 移除“智能助手”表头入口
- 智能助手由悬浮入口承接

## 关键文件

- `portal-ui/src/views/portal/PortalHomeView.vue`
- `portal-ui/src/layouts/portal/PortalSiteLayout.vue`
- `portal-ui/src/layouts/portal/components/PortalHeader.vue`
- `portal-ui/src/layouts/portal/components/PortalFloatingAgent.vue`
- `portal-ui/src/views/portal/ShipmentAgent.vue`
- `portal-ui/src/views/portal/PortalContactView.vue`

## 后续建议

- 为工具导航增加搜索框和热门关键词
- 为首页增加“全部工具”独立页面
- 将更多测算类工具从说明页落为独立可执行工具
- 将工具使用结果与询价表单、智能助手进一步串联

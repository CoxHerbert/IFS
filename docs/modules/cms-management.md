# CMS 管理

## 定位

CMS 是后台管理系统中的独立内容模块，用来维护门户新闻资讯内容。门户前台只负责读取已发布文章并展示，不提供公开写入能力。

当前菜单：

- 后台：`CMS管理 -> 新闻资讯`
- 门户：`/news`、`/news/:slug`

## 后台能力

后台新闻资讯支持：

- 文章列表查询
- 按标题、栏目、状态筛选
- 新增文章
- 修改文章
- 删除文章
- 预览门户文章
- 富文本正文编辑
- 富文本图片上传

后台页面文件：

- `baize-ui/src/views/cms/article/index.vue`
- `baize-ui/src/api/cms/article.ts`
- `baize-ui/src/components/RichTextEditor/index.vue`

## 前台展示

门户新闻资讯页从公开接口读取已发布文章。

门户文件：

- `portal-ui/src/views/portal/PortalNewsView.vue`
- `portal-ui/src/api/portal/article.ts`

公开接口：

- `GET /portal/articles`
- `GET /portal/articles/:slug`

## 后台接口

后台接口需要登录和权限校验：

| 接口 | 方法 | 权限 |
| --- | --- | --- |
| `/cms/article/list` | GET | `cms:article:list` |
| `/cms/article/:articleId` | GET | `cms:article:query` |
| `/cms/article` | POST | `cms:article:add` |
| `/cms/article` | PUT | `cms:article:edit` |
| `/cms/article/:articleIds` | DELETE | `cms:article:remove` |
| `/cms/article/upload-image` | POST | `cms:article:list` |

后端文件：

- `app/cms/models/article.go`
- `app/cms/dao/articleDao.go`
- `app/cms/service/articleService.go`
- `app/cms/controller/articleController.go`
- `app/routes/cmsRoutes/articleRouter.go`
- `app/routes/portalRoutes/articleRouter.go`

## 数据表

核心表：`cms_article`

主要字段：

- `article_id`：文章 ID
- `title`：标题
- `slug`：门户访问标识
- `summary`：摘要
- `category`：栏目
- `cover_url`：封面图
- `content`：富文本正文 HTML
- `status`：状态，`0` 已发布，`1` 草稿
- `sort`：排序
- `publish_time`：发布时间

初始化/升级脚本：

- `sql/ifs_business.sql`

该脚本会创建 `cms_article` 表，并新增独立后台菜单：

- `CMS管理`
- `新闻资讯`
- `cms:article:*` 权限按钮

## 富文本与图片

后台正文编辑器使用 wangEditor：

- 支持标题、字体、列表、链接、表格等常见富文本能力
- 支持上传图片并插入正文
- 图片上传接口为 `/cms/article/upload-image`
- 图片保存到后端静态资源目录：`/profile/cms/article/...`
- 图片访问 URL 会写入富文本 HTML 中

开发环境需要代理 `/profile` 到后端服务，否则编辑器和门户都无法回显上传图片：

- `baize-ui/vite.config.ts`
- `portal-ui/vite.config.ts`

生产环境需要由网关或 Nginx 暴露 `/profile` 静态资源。

## 安全边界

CMS 富文本和图片上传属于高风险入口，当前防护规则如下：

- 文章管理接口需要后台登录
- 图片上传接口需要 `cms:article:list` 权限
- 图片大小限制为 5MB
- 仅允许 `PNG / JPG / JPEG / GIF / WEBP`
- 后端同时校验文件扩展名和真实文件头
- 不允许上传 SVG、HTML、JS、PDF 等可执行或可嵌入脚本的文件
- 门户展示富文本前使用 DOMPurify 清洗 HTML

生产环境建议：

- `/profile` 禁止目录列表
- 设置 `X-Content-Type-Options: nosniff`
- 对 `/profile/cms/article/` 限制只返回图片 MIME
- 如使用对象存储，配置只读公开桶或 CDN，只允许后台服务端写入
- 定期清理未被文章引用的历史图片

## 注意事项

- 门户只展示 `status = 0` 的文章
- `slug` 可手动填写；为空时后端会根据标题生成
- 富文本正文以 HTML 保存，前台必须经过清洗后再渲染
- 不要在富文本中粘贴不可信的第三方脚本或 iframe

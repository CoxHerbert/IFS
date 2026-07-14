# 前端运行配置

## 目标

本文档记录后台管理系统 `baize-ui` 与门户工程 `portal-ui` 的前端运行时前缀、开发代理和静态资源约定。

## 前缀分类

| 类型 | 后台管理系统 | 门户工程 | 后端实际路由 |
| --- | --- | --- | --- |
| 后台业务 API | `VITE_APP_BASE_API=/admin-api` | 不使用 | 去掉 `/admin-api` 后转发到后端业务路由 |
| 门户业务 API | 不使用 | `VITE_PORTAL_API_PREFIX=/portal-api` | `/portal` |
| Agent API | `VITE_AGENT_API_PREFIX=/agent-api` | `VITE_AGENT_API_PREFIX=/agent-api` | `/agent-api` |
| 静态资源 | 固定 `/profile` | 固定 `/profile` | `/profile` |

## 后台管理系统

配置文件：
- `baize-ui/.env.development`
- `baize-ui/.env.staging`
- `baize-ui/.env.production`
- `baize-ui/vite.config.ts`

约定：
- 后台普通业务接口统一走 `VITE_APP_BASE_API`，当前为 `/admin-api`。
- 后台 Agent 接口统一走 `VITE_AGENT_API_PREFIX`，当前为 `/agent-api`。
- 上传头像、CMS 图片、付款凭证等静态资源统一通过 `/profile` 访问，不再单独配置 env。
- `baize-ui/src/utils/agent-request.ts` 负责后台 Agent 请求前缀。
- `baize-ui/src/utils/resource-url.ts` 负责 `/profile` 静态资源 URL 拼接。

## 门户工程

配置文件：
- `portal-ui/.env`
- `portal-ui/vite.config.ts`

约定：
- 门户业务接口统一走 `VITE_PORTAL_API_PREFIX`，当前为 `/portal-api`。
- 门户开发代理会将 `/portal-api` rewrite 到后端 `/portal`。
- 门户 Agent 接口统一走 `VITE_AGENT_API_PREFIX`，当前为 `/agent-api`。
- `/agent-api` 已是后端真实 Agent 路由，开发代理不再 rewrite 到 `/api`。
- `/profile` 固定用于 CMS 富文本图片、上传凭证等静态资源。
- `portal-ui/src/utils/portal-api.ts` 负责门户业务 API URL 拼接。
- `portal-ui/src/utils/agent-api.ts` 负责门户 Agent API URL 拼接。

## 迁移记录

本次调整完成了以下统一：
- 门户原 `/portal` 前端请求前缀改为 env 配置，并默认使用 `/portal-api`。
- 项目自有 Agent 接口由 `/api` 统一改为 `/agent-api`。
- 前端开发代理不再把 `/agent-api` rewrite 到 `/api`。
- `/profile` 静态资源前缀不再使用 env，后台和门户均固定代理该路径。
- 后端下发的 Agent `SubmitAPI` 统一改为 `/agent-api/agent/form/submit`。

## 不应改动的路径

- 前端 `@/api/...` 是源码目录别名，不是 HTTP 前缀。
- `app/agent/service/ollama_service.go` 中的 `BaseURL + "/api/chat"` 是 Ollama 服务自身接口路径，不属于本项目 API 前缀。

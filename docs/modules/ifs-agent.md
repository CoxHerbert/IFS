# IFS Agent

## 目标

IFS Agent 提供本地智能对话、出货文件分析和后台业务辅助能力。后端统一返回 IFS Block Protocol，前端按 `blocks` 动态渲染。

当前默认模型：

- Ollama `qwen2.5:7b`

## 入口

| 端 | 页面 | 接口前缀 |
| --- | --- | --- |
| 门户公开页 | 右下角悬浮助手 | `/agent-api/chat` |
| 客户工作台 | `/workspace/agent` | `/agent-api/chat` |
| 后台管理 | `货代业务 / Agent 对话` | `/agent/chat` |

说明：

- 门户旧的 `/shipment-agent` 独立出货分析模块已移除。
- 门户侧统一通过右下角悬浮助手承接问答和文件分析。
- 后台业务写账能力仅在后台 Agent 页面开放。

## 数据表

- `chat_session`
- `chat_message`
- `chat_memory`
- `agent_form_submission`

## 对话能力

- 多轮会话持久化
- 会话标题编辑
- 会话删除
- 普通问题调用 Ollama
- 命中本地规则时优先返回确定性结果
- 文件上传支持 `.xlsx` 和 `.csv`，不支持旧版 `.xls`
- 管理端 Agent 支持按当前用户数据权限查询客户、出货计划和收款单
- 管理端 Agent 支持通过自然语言生成新增收款确认表单
- 管理端 Agent 支持通过自然语言生成追加核销确认表单
- 收款及核销属于写账操作，必须由用户确认提交，不会从普通对话直接写库
- 客户工作台 Agent 支持查询当前客户自己的出货计划与收款核销摘要
- 客户工作台 Agent 支持提交带付款凭证的付款申报，申报状态默认为 `PENDING`
- 客户付款申报不会直接生成正式收款或修改出货计划付款状态

## IFS Block Protocol

统一返回：

```json
{
  "version": "1.0",
  "type": "agent_result",
  "title": "结果标题",
  "summary": "结果摘要",
  "blocks": []
}
```

当前支持 block：

- `summary`
- `metrics`
- `table`
- `file`
- `markdown`
- `error`
- `form`
- `action`
- `link`

## 主要接口

公开 / 客户端：

- `POST /agent-api/chat/session`
- `GET /agent-api/chat/sessions`
- `GET /agent-api/chat/session/:sessionId/messages`
- `PUT /agent-api/chat/session/:sessionId/title`
- `DELETE /agent-api/chat/session/:sessionId`
- `POST /agent-api/chat/send`
- `POST /agent-api/chat/session/:sessionId/shipment-analyze`
- `POST /agent-api/agent/form/submit`
- `POST /agent-api/agent/action/execute`

后台：

- `/agent/chat/*`
- `/agent/chat/session/:sessionId/shipment-analyze`
- `/agent/chat/form/submit`

## 关键文件

后端：

- `app/agent`
- `app/routes/agentRoutes`

门户前端：

- `portal-ui/src/layouts/portal/components/PortalFloatingAgent.vue`
- `portal-ui/src/views/workspace/WorkspaceAgentChatView.vue`
- `portal-ui/src/components/agent-renderer`

后台前端：

- `baize-ui/src/views/agent/chat/index.vue`

## 路由前缀调整

本次 Agent 接口前缀已统一为 `/agent-api`。

后端真实路由：
- `/agent-api/chat`
- `/agent-api/agent`
- `/agent-api/shipment`

前端使用方式：
- `portal-ui` 通过 `VITE_AGENT_API_PREFIX=/agent-api` 和 `src/utils/agent-api.ts` 拼接 Agent 请求。
- `baize-ui` 通过 `VITE_AGENT_API_PREFIX=/agent-api` 和 `src/utils/agent-request.ts` 拼接 Agent 请求。
- 开发代理不再将 `/agent-api` rewrite 为 `/api`。

后端下发的表单提交地址统一为：
- `/agent-api/agent/form/submit`

保留说明：
- 前端 `@/api/...` 是源码目录别名，不是 HTTP 路由前缀。
- Ollama 客户端调用的 `/api/chat` 是 Ollama 服务自身协议路径，不属于 IFS Agent 路由。

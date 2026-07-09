# IFS Agent

## 目标

IFS Agent 提供本地智能对话和出货文件分析能力，统一返回 IFS Block Protocol，由前端按 `blocks` 动态渲染。

当前模型：

- Ollama `qwen2.5:7b`

## 入口

| 端 | 页面 | 接口前缀 |
| --- | --- | --- |
| 门户悬浮助手 / 公开页 | `/agent` | `/api/chat` |
| 门户出货分析 | `/shipment-agent` | `/api/shipment/analyze` |
| 客户工作台 | `/customer/agent-chat` | `/api/chat` |
| 后台管理 | 货代业务 / Agent 对话 | `/agent/chat` |

## 数据表

- `chat_session`
- `chat_message`
- `chat_memory`
- `agent_form_submission`

## SQL

- 统一入口：`sql/ifs_init.sql`
- 业务合并脚本：`sql/ifs_business.sql`

## 对话能力

- 多轮会话持久化
- 会话标题编辑
- 会话删除
- 普通问题调用 Ollama
- 命中本地规则时优先返回确定性结果
- 文件上传支持 `.xlsx` 和 `.csv`，不支持旧版 `.xls`

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

## 主要接口

公开 / 客户端：

- `POST /api/chat/session`
- `GET /api/chat/sessions`
- `GET /api/chat/session/:sessionId/messages`
- `PUT /api/chat/session/:sessionId/title`
- `DELETE /api/chat/session/:sessionId`
- `POST /api/chat/send`
- `POST /api/chat/session/:sessionId/shipment-analyze`
- `POST /api/agent/form/submit`
- `POST /api/agent/action/execute`

后台：

- `/agent/chat/*`
- `/agent/chat/session/:sessionId/shipment-analyze`

## 关键文件

后端：

- `app/agent`
- `app/routes/agentRoutes`

前端：

- `portal-ui/src/views/portal/ChatAgent.vue`
- `portal-ui/src/layouts/portal/components/PortalFloatingAgent.vue`
- `portal-ui/src/views/workspace/WorkspaceAgentChatView.vue`
- `portal-ui/src/components/agent-renderer`
- `baize-ui/src/views/agent/chat/index.vue`

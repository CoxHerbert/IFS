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

## 本地模型配置

Agent 默认通过 Ollama 调用本地模型，配置位于 `config.yaml`：

```yaml
agent:
  ollama:
    base_url: "http://localhost:11434"
    default_model: "qwen2.5:7b"
    timeout: 90
    models:
      - label: "Qwen 2.5 7B"
        value: "qwen2.5:7b"
        description: "默认模型，适合日常货运问答和出货分析。"
        default: true
```

配置说明：
- `base_url`：Ollama 服务地址。
- `default_model`：未显式选择模型时使用的默认模型。
- `timeout`：调用本地模型的超时时间，单位为秒。
- `models`：前端模型下拉列表，`value` 必须是 Ollama 中已拉取的模型名。

运行前需要先准备模型：

```bash
ollama pull qwen2.5:7b
ollama serve
```

代码入口：
- `app/agent/service/ollama_service.go` 负责读取配置并调用 Ollama `/api/chat`。
- `app/agent/service/chat_service.go` 负责读取默认模型和模型列表。

注意：`/api/chat` 是 Ollama 自身接口路径，不属于 IFS 的 `/agent-api` 路由前缀。
## Agent 可视化配置

后台管理系统新增 `货代业务 / Agent 配置` 页面，用于维护 Agent 调用本地模型的运行参数。

配置项：
- `Ollama Base URL`：后台服务访问 Ollama 的 HTTP 地址。
- `默认模型`：未显式选择模型时使用的模型。
- `超时时间`：后台调用 Ollama 的请求超时时间，单位秒。
- `可选模型`：前端模型下拉列表，`value` 必须是 Ollama 已拉取的模型名。

后端配置接口：
- `GET /agent/config/ollama`
- `PUT /agent/config/ollama`
- `POST /agent/config/ollama/test`

配置持久化：
- 表结构脚本：`sql/ifs_business.sql`
- 配置表：`agent_runtime_config`
- 数据库配置优先，`config.yaml` 作为兜底。

阿里云服务器请求本地 Ollama 时，`Base URL` 必须填写阿里云后端能访问到的地址，例如：
- VPN / 内网互通地址：`http://10.0.0.12:11434`
- 内网穿透地址：`https://your-tunnel.example.com`
- 反向代理地址：`https://ai.example.com/ollama`

保存前可以点击“测试连接”，该测试由后端服务器发起到 `Base URL + /api/tags`，能验证阿里云服务器是否真正连得上本地 Ollama。
## 本次新增功能归类

### 后台可视化配置

- 页面入口：`货代业务 / Agent 配置`
- 前端页面：`baize-ui/src/views/agent/config/index.vue`
- 前端 API：`baize-ui/src/api/agent/config.ts`

### 后端运行时配置

- 配置接口：`GET /agent/config/ollama`
- 保存接口：`PUT /agent/config/ollama`
- 连接测试：`POST /agent/config/ollama/test`
- 读取优先级：数据库配置优先，`config.yaml` 兜底。

### 数据库与菜单脚本

- SQL：`sql/ifs_business.sql`
- 配置表：`agent_runtime_config`
- 菜单：`Agent 配置`
- 权限：`ifs:agent:config`
- 默认授权：`role_id = 1`

### 运维部署

- 运维文档：`docs/operations/agent-runtime-config.md`
- 新环境入口：`sql/ifs_init.sql`
- 已有环境升级：按现场差异从 `sql/ifs_business.sql` 提取 Agent 配置段执行。

阿里云后端访问本地 Ollama 时，`Ollama Base URL` 必须填写后端服务器可访问的地址，并建议通过 VPN、内网穿透或带白名单的反向代理暴露，不建议裸露 `11434` 到公网。

# IFS 本地 Agent

## 目标

IFS Agent 提供本地智能对话和出货文件分析能力，统一返回 IFS Block Protocol，由前端按 `blocks` 动态渲染。

模型：Ollama `qwen2.5:7b`

## 入口

| 端 | 页面 | 接口前缀 |
| --- | --- | --- |
| 门户悬浮助手 / 兼容公开页 | `右下角悬浮助手` / `/agent` | `/api/chat` |
| 门户出货分析 | `/shipment-agent` | `/api/shipment/analyze` |
| 客户端工作台 | `/customer/agent-chat` | `/api/chat` |
| 后台管理 | 货代业务 / Agent 对话 | `/agent/chat` |

客户端工作台会携带客户 `token`；后台管理会使用后台登录用户；公开门户未登录时走公共会话。

## 数据表

- `chat_session`
- `chat_message`
- `chat_memory`
- `agent_form_submission`

SQL 统一维护在 `sql/ifs_agent.sql`。

## 对话能力

- 多轮会话持久化
- 会话标题编辑
- 会话删除
- 普通问题调用 Ollama
- 命中本地规则时优先返回确定性结果
- 文件上传支持 `.xlsx` 和 `.csv`，不支持旧版 `.xls`

## 三端交互差异

### 门户

- 智能助手主入口已改为门户右下角悬浮面板。
- 门户表头不再展示“智能助手”菜单项。
- 门户对话输入区已统一为深色卡片样式，仅保留真实可用的上传和发送操作。

### 客户端工作台

- 对话标题点击后直接切换为输入框做内联重命名。
- 不再展示单独的“重命名”按钮。
- Header 中主题切换已移入设置图标下拉菜单。

### 后台管理

- Agent 对话页同样改为点击标题内联重命名。
- 不再展示重命名按钮。
- 后台展示文案统一以 `IFS` 命名，不再使用 `白泽` 作为可见系统名称。

## 本地技能

| 技能 | 说明 |
| --- | --- |
| 尺寸解析 | 支持 `100*200*150`、`100×200×150`、`100 x 200 x 150`、`100cm*200cm*150cm` |
| CBM 计算 | `长 × 宽 × 高 ÷ 1000000 × 箱数` |
| 柜型/LCL 判断 | 小体积优先 LCL，柜型规则见出货计划模块 |
| Excel/CSV 解析 | 服务端读取表格、识别表头、标准化货物明细 |
| 保存正式计划 | 通过 `form` block 提交后写入正式出货计划 |

正式出货计划归属规则见 [出货计划与出货查询](freight-shipment.md)。

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

前端不得写死业务字段，只按 `block.type` 渲染。

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
- `portal-ui/src/layouts/workspace/components/WorkspaceHeader.vue`
- `portal-ui/src/components/agent-renderer`
- `baize-ui/src/views/agent/chat/index.vue`

## 后续

- 保存成功后返回正式出货计划详情链接
- 文件解析结果增加置信度、错误行提示和确认页
- 增加附件历史记录

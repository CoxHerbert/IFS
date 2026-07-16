# Agent Runtime Configuration

This document records the runtime configuration added for IFS Agent local model access.

## Scope

The feature allows administrators to configure Agent model runtime settings from the backend management UI instead of editing `config.yaml` manually.

Backend menu:
- `货代业务 / Agent 配置`

Backend page:
- `baize-ui/src/views/agent/config/index.vue`

Backend APIs:
- `GET /agent/config/ollama`
- `PUT /agent/config/ollama`
- `POST /agent/config/ollama/test`

## Config Items

- `Ollama Base URL`: HTTP address that the backend server can use to access Ollama.
- `默认模型`: model used when the request does not specify a model.
- `超时时间`: Ollama request timeout in seconds.
- `可选模型`: model options shown in the frontend model selector.

The connection test calls `Base URL + /api/tags` from the backend server. This verifies whether the deployed backend can really reach the local model service.

## Persistence

Database table:
- `agent_runtime_config`

Migration script:
- `sql/ifs_business.sql`

The runtime loading order is:
1. Database configuration from `agent_runtime_config`.
2. `config.yaml` as fallback.
3. Built-in defaults as final fallback.

## Menu And Permission SQL

The migration script also creates backend menu records:

- `146`: `Agent 配置`
- `1200`: `Agent 配置权限`

It also grants these menus to `role_id = 1` when that role exists:

- `140`: `货代业务`
- `146`: `Agent 配置`
- `1200`: `Agent 配置权限`

For new environments, `sql/ifs_init.sql` runs `sql/ifs_business.sql`.

For existing environments, extract and run the Agent section from `sql/ifs_business.sql` according to the deployed database state.

## Alibaba Cloud Accessing Local Ollama

If the backend is deployed on Alibaba Cloud and Ollama runs on a local machine, `Ollama Base URL` must be an address reachable from the Alibaba Cloud backend.

Recommended options:
- VPN / private network, such as Tailscale, ZeroTier, or WireGuard.
- Tunnel service, such as frp, ngrok, or Cloudflare Tunnel.
- Reverse proxy with IP whitelist.

Do not expose Ollama `11434` directly to the public internet without authentication or IP restrictions.

Examples:
- `http://10.0.0.12:11434`
- `https://your-tunnel.example.com`
- `https://ai.example.com/ollama`

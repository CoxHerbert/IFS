# IFS 前后端打包部署手册

开发环境打包 → 上传服务器 → 发布验证 → 快速回滚

| 项目 | 部署信息 |
| --- | --- |
| 后端 | Go 可执行文件 /data/app/ifs-api/ifs-api |
| 后端源码 | `server/`（Go 模块根目录） |
| 后端配置 | /data/app/ifs-api/config.yaml |
| 后端服务 | ifs-api.service，监听 8080 |
| 管理端 | /data/www/ifs-admin |
| 门户端 | /data/www/ifs-portal |
| 记账 H5 | /data/www/ifs-ledger |
| 服务器 | Linux + systemd + Nginx |
| 版本 | v1.0 · 2026-07-14 |

服务器 `/data` 目录结构（以服务器实际目录为准）：

```text
/data/
├── www/    # 前端站点及前端发布包上传目录
├── app/    # 后端目录，后端程序位于 /data/app/ifs-api/ifs-api
└── logs/   # 日志目录
```

> 注意：`/data/app/ifs-api/` 是后端目录，目录内的 `ifs-api` 是后端可执行文件；`/data/www/` 下没有后端程序。

> 适用范围：用于日常前端、后端版本发布。首次配置 DNS、HTTPS 和 Nginx 请参考《IFS 系统部署与运维手册》。

仓库提供可直接调整后部署的 Nginx 配置：`deploy/nginx/ifs.conf`。前端打包脚本：`deploy/package-frontend.ps1`。

# 1. 发布流程总览

| 阶段 | 前端 | 后端 |
| --- | --- | --- |
| 打包 | npm/pnpm build 生成 dist | 交叉编译 Linux 可执行文件 |
| 上传 | 上传 ZIP 到 /data/www | 上传 ifs-api.new 到 /data/app/ifs-api |
| 备份 | 备份当前站点目录 | 备份当前 ifs-api |
| 发布 | 解压并替换静态文件 | 替换文件并重启 systemd |
| 验证 | 页面、资源、F12 Network | 状态、端口、日志、接口 |
| 回滚 | 恢复站点备份 | 恢复 ifs-api.bak 并重启 |

> 发布原则：先上传临时文件，再备份当前版本，验证无误后替换。不要直接覆盖正在使用的后端二进制。

# 2. 发布前准备

- 确认本地代码已拉取目标分支，依赖安装成功，构建无报错。

- 确认前端生产 API 地址正确：管理端使用 /admin-api，门户端使用 /portal-api。

- 确认 config.yaml 与服务器环境匹配；一般发布程序时不覆盖服务器 config.yaml。

- 确认 SSH 可登录服务器，并预留足够磁盘空间。

```bash
# 服务器检查
df -h
systemctl status ifs-api nginx --no-pager
nginx -t
```

# 3. 前端打包

## 3.1 确认包管理器

优先使用项目已有锁文件对应的工具，不要混用，否则可能造成依赖版本变化。

| 锁文件 | 工具 | 安装与构建 |
| --- | --- | --- |
| pnpm-lock.yaml | pnpm | pnpm install --frozen-lockfile；pnpm build |
| package-lock.json | npm | npm ci；npm run build |
| yarn.lock | yarn | yarn install --frozen-lockfile；yarn build |

本仓库前端都提交了 `pnpm-lock.yaml`，推荐在仓库根目录用脚本构建并生成上传 ZIP：

```powershell
# 管理端、门户端和账本 H5
powershell -ExecutionPolicy Bypass -File .\deploy\package-frontend.ps1

# 仅构建其中一个
powershell -ExecutionPolicy Bypass -File .\deploy\package-frontend.ps1 -Target admin
powershell -ExecutionPolicy Bypass -File .\deploy\package-frontend.ps1 -Target portal
powershell -ExecutionPolicy Bypass -File .\deploy\package-frontend.ps1 -Target ledger
```

ZIP 生成在对应前端目录（`ifs-admin/ifs-admin.zip`、`ifs-portal/ifs-portal.zip`、`ifs-ledger/ifs-ledger.zip`），ZIP 根目录直接包含 `index.html` 和静态资源。该脚本执行锁文件安装、生产构建并校验产物。管理端执行 `build:prod`，门户端执行 `build`（含 TypeScript 检查），账本执行 `build:h5`。

## 3.2 管理端打包

```powershell
# PowerShell：进入管理端项目
cd D:\workspace\IFS\ifs-admin

# pnpm 项目
pnpm install --frozen-lockfile
pnpm build

# 或 npm 项目
npm ci
npm run build
```

构建产物通常位于 dist 目录。检查 dist 下是否直接包含 index.html 和 assets。

## 3.3 门户端打包

```powershell
# PowerShell：进入门户端项目
cd D:\workspace\IFS\ifs-portal
pnpm install --frozen-lockfile
pnpm build
```

## 3.4 压缩构建产物

若使用上述脚本，可跳过手工压缩步骤。

在 PowerShell 中进入 dist 目录后压缩目录内的内容，避免服务器解压后多出一层 dist。

```powershell
# 管理端示例
cd D:\workspace\IFS\ifs-admin\dist
Compress-Archive -Path * -DestinationPath ..\ifs-admin.zip -Force

# 门户端示例
cd D:\workspace\IFS\ifs-portal\dist
Compress-Archive -Path * -DestinationPath ..\ifs-portal.zip -Force
```

> 正确结构：ZIP 解压后应直接看到 index.html、assets/，而不是 dist/index.html。

# 4. 前端上传与发布

## 4.1 上传

```powershell
# 本地 PowerShell
scp D:\workspace\IFS\ifs-admin\ifs-admin.zip root@服务器公网IP:/data/www/
scp D:\workspace\IFS\ifs-portal\ifs-portal.zip root@服务器公网IP:/data/www/
```

若 /data/www 不存在，先在服务器创建：

```bash
mkdir -p /data/www /data/backup
```

## 4.2 发布管理端

```bash
# 服务器执行
stamp=$(date +%Y%m%d_%H%M%S)
cp -a /data/www/ifs-admin /data/backup/ifs-admin_$stamp

mkdir -p /tmp/ifs-admin-release
unzip -oq /data/www/ifs-admin.zip -d /tmp/ifs-admin-release
test -f /tmp/ifs-admin-release/index.html

find /data/www/ifs-admin -mindepth 1 -maxdepth 1 -exec rm -rf -- {} +
cp -a /tmp/ifs-admin-release/. /data/www/ifs-admin/
ls -lah /data/www/ifs-admin
```

## 4.3 发布门户端

```bash
stamp=$(date +%Y%m%d_%H%M%S)
cp -a /data/www/ifs-portal /data/backup/ifs-portal_$stamp

mkdir -p /tmp/ifs-portal-release
unzip -oq /data/www/ifs-portal.zip -d /tmp/ifs-portal-release
test -f /tmp/ifs-portal-release/index.html

find /data/www/ifs-portal -mindepth 1 -maxdepth 1 -exec rm -rf -- {} +
cp -a /tmp/ifs-portal-release/. /data/www/ifs-portal/
ls -lah /data/www/ifs-portal
```

> Nginx：只替换静态文件一般不需要重启 Nginx。如果同时改了 Nginx 配置，则执行 nginx -t && systemctl reload nginx。

# 5. Go 后端打包

## 5.1 确认 Go 环境

```bash
go version
go env GOPATH GOOS GOARCH
```

## 5.2 Windows 编译 Linux 版本

可在仓库根目录运行脚本，一次完成依赖下载、Linux 交叉编译和 ZIP 打包：

```powershell
powershell -ExecutionPolicy Bypass -File .\deploy\package-backend.ps1
```

默认生成 `deploy/packages/ifs-api-linux-amd64.zip`，ZIP 内含 Linux 可执行文件 `ifs-api`。ARM64 服务器可指定 `-Architecture arm64`。打包目录不含 `config.yaml`，发布时保留服务器现有配置。

也可以按以下步骤手工编译：

在 `server/` Go 模块根目录执行。PowerShell 设置环境变量使用 $env: 写法。

```powershell
cd D:\workspace\IFS\server

$env:CGO_ENABLED="0"
$env:GOOS="linux"
$env:GOARCH="amd64"
go mod download
go build -trimpath -ldflags="-s -w" -o ifs-api ./

# 检查文件
Get-Item .\ifs-api
```

> CGO 项目：如果项目依赖 SQLite、Oracle 驱动或其他 CGO 库，CGO_ENABLED=0 可能无法编译或运行，需要在 Linux 环境中构建。

## 5.3 Linux 本机构建（可选）

```bash
cd /path/to/IFS/server
go mod download
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -trimpath -ldflags='-s -w' -o ifs-api ./
```

# 6. 后端上传与发布

## 6.1 上传临时文件

```powershell
# 本地 PowerShell
scp D:\workspace\IFS\server\ifs-api root@服务器公网IP:/data/app/ifs-api/ifs-api.new
```

## 6.2 替换并重启

```bash
# 服务器执行
chmod +x /data/app/ifs-api/ifs-api.new

# 先确认新文件存在
ls -lh /data/app/ifs-api/ifs-api.new

# 备份当前版本
stamp=$(date +%Y%m%d_%H%M%S)
cp /data/app/ifs-api/ifs-api /data/backup/ifs-api_$stamp

# 替换并重启
mv /data/app/ifs-api/ifs-api.new /data/app/ifs-api/ifs-api
chmod +x /data/app/ifs-api/ifs-api
systemctl restart ifs-api

# 验证
systemctl status ifs-api --no-pager
ss -lntp | grep ifs-api
curl -i http://127.0.0.1:8080/portal/version
journalctl -u ifs-api -n 100 --no-pager
```

> 配置文件：默认保留服务器现有 /data/app/ifs-api/config.yaml。只有新增配置项时才单独备份、修改并核对。

# 7. 发布后验证

| 检查项 | 操作 | 通过标准 |
| --- | --- | --- |
| 后端状态 | systemctl status ifs-api | active (running) |
| 后端端口 | ss -lntp \| grep ifs-api | 监听 8080 |
| 后端接口 | curl -i http://127.0.0.1:8080/portal/version | 返回 HTTP 200 和版本信息 |
| 后端日志 | journalctl -u ifs-api -n 100 | 无 panic/连接失败 |
| 管理端 | 打开 admin.baozenan.online | 页面和资源正常 |
| 门户端 | 打开 baozenan.online | 页面和资源正常 |
| API | F12 → Network | 核心请求成功，无 404/502/CORS |
| 业务 | 登录并操作关键流程 | 功能符合本次发布 |

```bash
curl -i http://127.0.0.1:8080/portal/version
curl -ik https://api.baozenan.online/admin-api/对应外部接口
curl -Ik https://admin.baozenan.online
curl -Ik https://baozenan.online
```

> 404 判断：该 Go 项目对不存在的路径可能返回 HTTP 200 和 {"msg":"404"}。必须使用真实接口路径及正确的 GET/POST 方法验证。

# 8. 回滚

## 8.1 后端回滚

先从 /data/backup 找到最近一次备份，再恢复。

```bash
ls -lht /data/backup/ifs-api_*
cp /data/backup/ifs-api_YYYYMMDD_HHMMSS /data/app/ifs-api/ifs-api
chmod +x /data/app/ifs-api/ifs-api
systemctl restart ifs-api
systemctl status ifs-api --no-pager
```

## 8.2 前端回滚

```bash
ls -ldht /data/backup/ifs-admin_*
find /data/www/ifs-admin -mindepth 1 -maxdepth 1 -exec rm -rf -- {} +
cp -a /data/backup/ifs-admin_YYYYMMDD_HHMMSS/. /data/www/ifs-admin/

ls -ldht /data/backup/ifs-portal_*
find /data/www/ifs-portal -mindepth 1 -maxdepth 1 -exec rm -rf -- {} +
cp -a /data/backup/ifs-portal_YYYYMMDD_HHMMSS/. /data/www/ifs-portal/
```

# 9. 常见问题

| 现象 | 原因 | 处理 |
| --- | --- | --- |
| exec format error | 上传了 Windows 程序或架构不一致 | 以 GOOS=linux、GOARCH=amd64 重编译 |
| permission denied | 可执行权限丢失 | chmod +x ifs-api |
| 服务启动后退出 | config.yaml、数据库、Redis 或端口异常 | 查看 journalctl 日志 |
| 前端白屏 | 资源 base 路径或上传目录错误 | 检查 index.html、assets 和控制台 |
| 前端仍是旧版本 | 浏览器/CDN 缓存 | 强制刷新并核对资源文件名 |
| API 502 | 后端未启动或 Nginx 端口错误 | 检查 8080 和 proxy_pass |
| API 404 | 接口路径/方法或前缀处理错误 | F12 取真实请求并与 8080 直连对比 |

# 10. 发布检查清单

☐ 已确认发布分支与提交版本

☐ 前端生产 API 地址正确且构建成功

☐ 前端 ZIP 内直接包含 index.html 和 assets/

☐ 后端已编译为 Linux amd64 可执行文件

☐ 服务器旧版本已备份

☐ 服务器 config.yaml 未被误覆盖

☐ ifs-api 为 active (running)，8080 正常监听

☐ 管理端、门户端页面加载正常

☐ F12 Network 无 404、502、CORS 错误

☐ 登录和关键业务流程验证通过

☐ 已保留本次发布包和回滚版本

# 附录 A：Nginx 初次部署

`deploy/nginx/ifs.conf` 提供管理端、门户端、记账 H5 和可选 API 域名的 server 配置，包含 HTTPS、SPA history 回退、API 前缀转发及 `/profile` 静态资源转发。记账 H5 域名为 `ledger.baozenan.online`，站点目录为 `/data/www/ifs-ledger`。部署前按实际环境修改 `server_name`、证书路径和站点目录；如果没有单独的 API 域名，可删除最后一组 `api.baozenan.online` server 块。前端域名的证书文件需预先签发并放到配置指定路径。

```bash
# 首次部署前，创建目录并上传配置文件到 /etc/nginx/conf.d/ifs.conf
mkdir -p /data/www/ifs-admin /data/www/ifs-portal /data/www/ifs-ledger
nginx -t
systemctl enable --now nginx
systemctl reload nginx
```

配置将 `/admin-api/` 转发为后端路径（剥离 `/admin-api`），将 `/portal-api/` 转发为 `/portal/`，将账本站点的 `/ledger-api/` 转发为后端实际 `/ledger/` 路由，并将 `/profile/` 原样转发到 8080。Nginx 安装若未提供 `/etc/nginx/proxy_params`，请创建包含 `Host`、`X-Real-IP`、`X-Forwarded-For`、`X-Forwarded-Proto` 的代理头配置，或将现有发行版等价配置路径填入 `include`。每次改动后都先执行 `nginx -t`，通过后再 reload。

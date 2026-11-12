# 轻记 · 个人记账

基于 uni-app Vue 3 + Vite，使用 Wot UI V2 2.4.0（`@wot-ui/ui`），目标支持 H5、微信小程序和支付宝小程序。

底部导航使用共享自定义 tabbar 切换账本、明细、导入和我的；H5 使用 history 路由，地址不包含 `#`。部署 H5 时需配置站点将深层路径回退到 `index.html`。

## 开发

```bash
pnpm install
pnpm dev:h5
pnpm dev:mp-weixin
pnpm dev:mp-alipay
```

## 账单导入

首版提供微信、支付宝账单文件选择入口，并将交易模型和本地持久化层独立在 `src/utils/ledger.ts`。账单需由用户从支付平台导出后手动选择；不同导出版本的 CSV/Excel 表头和编码需要按真实样本适配。微信小程序限制直接访问本地路径，后续应通过文件 API 读取内容，在本机完成解析、预览和去重后再保存。

当前导入页只完成文件选择交互，交易行解析/确认流程还未接入。

## 账号与数据

- 登录、注册均为可选项，不影响未登录时的本地记账和账单查看。
- 登录和注册分别提供独立页面，均可返回本地使用。
- “我的”页可编辑昵称、手机号和邮箱；个人资料保存在本机。免登录用户首次打开“我的”页时生成设备 UID，登录用户使用服务端账号 ID 派生的稳定 UID。
- 所有记录默认只保存在设备本地。
- 登录后云同步仍默认关闭；用户在“我的与数据”中明确开启后才允许上传。
服务端 API 已实现，执行 `server/sql/ledger.sql` 后可使用 `/ledger/auth/register`、`/ledger/auth/login`、`/ledger/sync-consent` 和 `/ledger/records`。H5 开发服务器默认通过 Vite 代理到 `http://127.0.0.1:8080`；真机/小程序构建需将 `VITE_LEDGER_API_BASE` 配成已备案并配置合法域名的 HTTPS API 地址。

后台“记账管理”菜单在应用 `server/sql/ifs_upgrade.sql` 后提供。后台管理员可查询注册账号和用户主动同步到云端的账单，进入单个用户查看及删除云端记录；仅保存在用户设备的本地数据不会上传，也无法从后台查看。

| Method | Path | 用途 |
| --- | --- | --- |
| POST | `/ledger/auth/register` | 注册（`loginId`、`password`），返回 bearer token |
| POST | `/ledger/auth/login` | 登录并返回同步授权状态 |
| GET | `/ledger/sync-consent` | 查询云同步授权 |
| PUT | `/ledger/sync-consent` | 显式开启/关闭云同步；关闭会删除云端账单 |
| PUT | `/ledger/records` | 在授权开启时批量上传/覆盖账单 |
| GET | `/ledger/records` | 在授权开启时回显该账号云端账单 |

登录后接口使用 `Authorization: Bearer <token>`。账单按账户隔离，密码使用 bcrypt 哈希保存。关闭云同步会清空该账号云端账单，本机数据继续保留。

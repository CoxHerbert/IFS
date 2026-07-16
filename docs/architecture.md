# IFS Architecture

## Application Surfaces

IFS is organized by user-facing surface first, then by business module.

Current surfaces:
- `门户`: public website and public lead/contact entry.
- `客户端`: customer workspace after customer login.
- `后台管理`: internal operation console for staff.

Planned surface:
- `移动端`: future mobile app or mobile H5 client. It should reuse backend modules and permissions instead of creating duplicated business logic.

## Backend Modules

Backend code is grouped under `app/` by module:

| Module | Directory | Scope |
| --- | --- | --- |
| Portal | `app/portal` | Public website contact and portal-facing data |
| Customer | `app/customer` | Customer records, customer accounts, customer workspace auth and menus |
| Admin/System | `app/system`, `app/monitor`, `app/quartz`, `app/genTable` | Backend management foundation, roles, menus, users, monitoring and tools |
| Freight | `app/freight` | Shipment plans, receipts, payment declarations |
| CMS | `app/cms` | News/article content management |
| Notification | `app/notification` | Backend notification center |
| Agent | `app/agent` | Agent chat, local skills, Ollama runtime config, form/action protocol |
| Common | `app/common`, `app/utils`, `app/constant`, `app/setting` | Shared infrastructure |

Routes are grouped under `app/routes/*Routes` with the same module boundary.

## SQL Modules

SQL is kept intentionally small at the file level:

| Script | Purpose |
| --- | --- |
| `sql/baize2022-01-08.sql` | Base system tables and original admin framework data |
| `sql/ifs_business.sql` | IFS business modules: portal, customer, freight, CMS, notification, Agent |
| `sql/ifs_init.sql` | Unified entrypoint for new environments |

Do not add scattered dated module scripts for regular feature work. Merge module DDL, menus and permissions into `sql/ifs_business.sql` under a clear module section.

Current `sql/ifs_business.sql` module sections:
- Portal website and contact leads.
- Customer workspace and customer management.
- Freight shipment, receipt and payment declaration.
- Agent chat, form/action protocol and runtime configuration.
- Notification center.
- CMS article management.

## Module Rules

- Portal and customer workspace are separate surfaces. Customer workspace still belongs to the `customer` backend module.
- Customer management is its own module, not part of freight or Agent.
- Agent is an independent backend module. It may call customer/freight services, but it should not own their data tables.
- Backend management is a surface and a foundation module. Business modules expose their management APIs through backend routes.
- Future mobile endpoints should reuse existing module services and add mobile-specific route groups only when the response contract differs.

## Future Mobile Surface

When mobile is added, keep it as a new surface, not a new copy of business modules.

Recommended shape:
- Routes: `app/routes/mobileRoutes`
- Controllers: reuse existing module services where possible.
- Auth: add mobile token/session handling only if it differs from customer/admin auth.
- SQL: add mobile-specific tables only for mobile device/session/push needs.

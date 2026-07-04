-- ----------------------------
-- Migrate customer portal permissions to workspace permissions
-- ----------------------------

SET @portal_menu_exists = (
  SELECT COUNT(*)
  FROM information_schema.TABLES
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'customer_portal_menu'
);
SET @portal_role_exists = (
  SELECT COUNT(*)
  FROM information_schema.TABLES
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'customer_portal_role'
);
SET @portal_role_menu_exists = (
  SELECT COUNT(*)
  FROM information_schema.TABLES
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'customer_portal_role_menu'
);
SET @portal_account_role_exists = (
  SELECT COUNT(*)
  FROM information_schema.TABLES
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'customer_portal_account_role'
);

SET FOREIGN_KEY_CHECKS = 0;

SET @sql = IF(@portal_role_menu_exists > 0, 'RENAME TABLE `customer_portal_role_menu` TO `customer_workspace_role_menu`', 'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql = IF(@portal_account_role_exists > 0, 'RENAME TABLE `customer_portal_account_role` TO `customer_workspace_account_role`', 'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql = IF(@portal_role_exists > 0, 'RENAME TABLE `customer_portal_role` TO `customer_workspace_role`', 'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql = IF(@portal_menu_exists > 0, 'RENAME TABLE `customer_portal_menu` TO `customer_workspace_menu`', 'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET FOREIGN_KEY_CHECKS = 1;

SET @workspace_menu_has_is_cache = (
  SELECT COUNT(*)
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'customer_workspace_menu' AND COLUMN_NAME = 'is_cache'
);
SET @sql = IF(@workspace_menu_has_is_cache = 0, 'ALTER TABLE `customer_workspace_menu` ADD COLUMN `is_cache` char(1) NOT NULL DEFAULT ''0'' COMMENT ''是否缓存（0缓存 1不缓存）'' AFTER `component`', 'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

UPDATE `customer_workspace_menu`
SET `component` = CASE `component`
  WHEN 'customer/workspace' THEN 'workspace/dashboard'
  WHEN 'customer/account' THEN 'workspace/account-profile'
  WHEN 'customer/shipment' THEN 'workspace/shipment-tracking'
  WHEN 'customer/shipment-assistant' THEN 'workspace/shipment-assistant'
  WHEN 'customer/agent-chat' THEN 'workspace/agent-chat'
  ELSE `component`
END,
`menu_name` = CASE `menu_id`
  WHEN 20001 THEN '工作台'
  WHEN 20002 THEN '账号资料'
  WHEN 20003 THEN '出货查询'
  WHEN 20004 THEN '智能出货助手'
  WHEN 20005 THEN 'Agent 对话'
  ELSE `menu_name`
END,
`remark` = CASE `menu_id`
  WHEN 20001 THEN '客户端工作台'
  WHEN 20002 THEN '客户端账号资料'
  WHEN 20003 THEN '客户端出货查询'
  WHEN 20004 THEN '客户端智能出货助手'
  WHEN 20005 THEN '客户端 Agent 对话'
  ELSE `remark`
END,
`is_cache` = IFNULL(NULLIF(`is_cache`, ''), '0');

INSERT INTO `customer_workspace_menu` (`menu_id`, `parent_id`, `menu_name`, `order_num`, `path`, `component`, `is_cache`, `menu_type`, `visible`, `status`, `perms`, `icon`, `remark`, `create_by`, `create_time`, `update_by`, `update_time`)
SELECT 20004, 0, '智能出货助手', '4', 'shipment-assistant', 'workspace/shipment-assistant', '0', 'C', '0', '0', 'portal:shipmentAssistant:view', 'CalculatorOutlined', '客户端智能出货助手', 'admin', now(), 'admin', now()
FROM dual
WHERE NOT EXISTS (SELECT 1 FROM `customer_workspace_menu` WHERE `menu_id` = 20004);

INSERT INTO `customer_workspace_menu` (`menu_id`, `parent_id`, `menu_name`, `order_num`, `path`, `component`, `is_cache`, `menu_type`, `visible`, `status`, `perms`, `icon`, `remark`, `create_by`, `create_time`, `update_by`, `update_time`)
SELECT 20005, 0, 'Agent 对话', '5', 'agent-chat', 'workspace/agent-chat', '0', 'C', '0', '0', 'portal:agentChat:view', 'MessageOutlined', '客户端 Agent 对话', 'admin', now(), 'admin', now()
FROM dual
WHERE NOT EXISTS (SELECT 1 FROM `customer_workspace_menu` WHERE `menu_id` = 20005);

INSERT INTO `customer_workspace_role_menu` (`role_id`, `menu_id`)
SELECT 20001, 20004
FROM dual
WHERE EXISTS (SELECT 1 FROM `customer_workspace_role` WHERE `role_id` = 20001)
  AND EXISTS (SELECT 1 FROM `customer_workspace_menu` WHERE `menu_id` = 20004)
  AND NOT EXISTS (SELECT 1 FROM `customer_workspace_role_menu` WHERE `role_id` = 20001 AND `menu_id` = 20004);

INSERT INTO `customer_workspace_role_menu` (`role_id`, `menu_id`)
SELECT 20001, 20005
FROM dual
WHERE EXISTS (SELECT 1 FROM `customer_workspace_role` WHERE `role_id` = 20001)
  AND EXISTS (SELECT 1 FROM `customer_workspace_menu` WHERE `menu_id` = 20005)
  AND NOT EXISTS (SELECT 1 FROM `customer_workspace_role_menu` WHERE `role_id` = 20001 AND `menu_id` = 20005);

INSERT INTO `customer_workspace_role_menu` (`role_id`, `menu_id`)
SELECT r.`role_id`, 20005
FROM `customer_workspace_role` r
WHERE r.`status` = '0'
  AND r.`del_flag` = '0'
  AND EXISTS (SELECT 1 FROM `customer_workspace_menu` WHERE `menu_id` = 20005)
  AND NOT EXISTS (
    SELECT 1
    FROM `customer_workspace_role_menu` rm
    WHERE rm.`role_id` = r.`role_id`
      AND rm.`menu_id` = 20005
  );

DELETE FROM `sys_role_menu`
WHERE `menu_id` IN (134, 135, 1255, 1256, 1257, 1258, 1259, 1260, 1261, 1262, 1263, 1264);

DELETE FROM `sys_menu`
WHERE `menu_id` IN (134, 135, 1255, 1256, 1257, 1258, 1259, 1260, 1261, 1262, 1263, 1264);

INSERT INTO `sys_menu` VALUES
(134, '客户端菜单', 130, 4, 'portal-menu', 'customer/portalMenu/index', 1, 0, 'C', '0', '0', 'customer:portalMenu:list', 'tree-table', 'admin', now(), '', NULL, '客户端菜单管理'),
(135, '客户端角色', 130, 5, 'portal-role', 'customer/portalRole/index', 1, 0, 'C', '0', '0', 'customer:portalRole:list', 'peoples', 'admin', now(), '', NULL, '客户端角色管理');

INSERT INTO `sys_menu` VALUES
(1255, '客户端菜单查询', 134, 1, '#', '', 1, 0, 'F', '0', '0', 'customer:portalMenu:query', '#', 'admin', now(), '', NULL, ''),
(1256, '客户端菜单新增', 134, 2, '#', '', 1, 0, 'F', '0', '0', 'customer:portalMenu:add', '#', 'admin', now(), '', NULL, ''),
(1257, '客户端菜单修改', 134, 3, '#', '', 1, 0, 'F', '0', '0', 'customer:portalMenu:edit', '#', 'admin', now(), '', NULL, ''),
(1258, '客户端菜单删除', 134, 4, '#', '', 1, 0, 'F', '0', '0', 'customer:portalMenu:remove', '#', 'admin', now(), '', NULL, '');

INSERT INTO `sys_menu` VALUES
(1259, '客户端角色查询', 135, 1, '#', '', 1, 0, 'F', '0', '0', 'customer:portalRole:query', '#', 'admin', now(), '', NULL, ''),
(1260, '客户端角色新增', 135, 2, '#', '', 1, 0, 'F', '0', '0', 'customer:portalRole:add', '#', 'admin', now(), '', NULL, ''),
(1261, '客户端角色修改', 135, 3, '#', '', 1, 0, 'F', '0', '0', 'customer:portalRole:edit', '#', 'admin', now(), '', NULL, ''),
(1262, '客户端角色删除', 135, 4, '#', '', 1, 0, 'F', '0', '0', 'customer:portalRole:remove', '#', 'admin', now(), '', NULL, ''),
(1263, '客户端角色分配', 135, 5, '#', '', 1, 0, 'F', '0', '0', 'customer:portalRole:assign', '#', 'admin', now(), '', NULL, ''),
(1264, '客户账号角色配置', 133, 6, '#', '', 1, 0, 'F', '0', '0', 'customer:account:edit', '#', 'admin', now(), '', NULL, '');

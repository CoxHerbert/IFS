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

UPDATE `customer_workspace_menu`
SET `component` = CASE `component`
  WHEN 'customer/workspace' THEN 'workspace/dashboard'
  WHEN 'customer/account' THEN 'workspace/account-profile'
  WHEN 'customer/shipment' THEN 'workspace/shipment-tracking'
  ELSE `component`
END,
`menu_name` = CASE `menu_id`
  WHEN 20001 THEN '工作台'
  WHEN 20002 THEN '账号资料'
  WHEN 20003 THEN '出货查询'
  ELSE `menu_name`
END,
`remark` = CASE `menu_id`
  WHEN 20001 THEN '客户端工作台'
  WHEN 20002 THEN '客户端账号资料'
  WHEN 20003 THEN '客户端出货查询'
  ELSE `remark`
END;

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

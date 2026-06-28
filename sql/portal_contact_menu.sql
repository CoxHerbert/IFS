-- ----------------------------
-- Menu data for portal_contact
-- 后台菜单：客户管理 -> 官网线索
-- 权限标识：
--   portal:contact:list
--   portal:contact:query
--   portal:contact:edit
--   portal:contact:remove
--   portal:contact:export
-- ----------------------------

-- 删除老菜单：系统管理 -> 官网线索，以及本脚本可能重复创建的新菜单。
DELETE FROM `sys_role_menu`
WHERE `menu_id` IN (
  SELECT `menu_id` FROM `sys_menu`
  WHERE `perms` IN (
    'portal:contact:list',
    'portal:contact:query',
    'portal:contact:edit',
    'portal:contact:remove',
    'portal:contact:export'
  )
)
OR `menu_id` IN (120, 130, 131, 1120, 1121, 1122, 1123, 1130, 1131, 1132, 1133);

DELETE FROM `sys_menu`
WHERE `perms` IN (
  'portal:contact:list',
  'portal:contact:query',
  'portal:contact:edit',
  'portal:contact:remove',
  'portal:contact:export'
)
OR `menu_id` IN (120, 130, 131, 1120, 1121, 1122, 1123, 1130, 1131, 1132, 1133);

-- 新建顶级目录：客户管理
INSERT INTO `sys_menu` VALUES
(130, '客户管理', 0, 5, 'customer', NULL, 1, 0, 'M', '0', '0', '', 'peoples', 'admin', now(), '', NULL, '客户管理目录');

-- 新建菜单：客户管理 -> 官网线索
INSERT INTO `sys_menu` VALUES
(131, '官网线索', 130, 1, 'contact', 'customer/contact/index', 1, 0, 'C', '0', '0', 'portal:contact:list', 'message', 'admin', now(), '', NULL, '官网联系我们线索管理菜单');

INSERT INTO `sys_menu` VALUES
(1130, '官网线索查询', 131, 1, '#', '', 1, 0, 'F', '0', '0', 'portal:contact:query', '#', 'admin', now(), '', NULL, ''),
(1131, '官网线索修改', 131, 2, '#', '', 1, 0, 'F', '0', '0', 'portal:contact:edit', '#', 'admin', now(), '', NULL, ''),
(1132, '官网线索删除', 131, 3, '#', '', 1, 0, 'F', '0', '0', 'portal:contact:remove', '#', 'admin', now(), '', NULL, ''),
(1133, '官网线索导出', 131, 4, '#', '', 1, 0, 'F', '0', '0', 'portal:contact:export', '#', 'admin', now(), '', NULL, '');

-- 如需给普通角色（role_id = 2）授权，取消下方注释。
-- INSERT INTO `sys_role_menu` VALUES (2, 130);
-- INSERT INTO `sys_role_menu` VALUES (2, 131);
-- INSERT INTO `sys_role_menu` VALUES (2, 1130);
-- INSERT INTO `sys_role_menu` VALUES (2, 1131);
-- INSERT INTO `sys_role_menu` VALUES (2, 1132);
-- INSERT INTO `sys_role_menu` VALUES (2, 1133);

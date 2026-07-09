DELETE FROM `sys_role_menu`
WHERE `menu_id` IN (142, 1166, 1167, 1168)
OR `menu_id` IN (
  SELECT `menu_id` FROM `sys_menu`
  WHERE `perms` IN (
    'system:notification:list',
    'system:notification:edit',
    'system:notification:remove'
  )
);

DELETE FROM `sys_menu`
WHERE `menu_id` IN (142, 1166, 1167, 1168)
OR `perms` IN (
  'system:notification:list',
  'system:notification:edit',
  'system:notification:remove'
);

INSERT INTO `sys_menu` VALUES
(142, '消息通知', 1, 8, 'notification', 'system/notification/index', 1, 0, 'C', '0', '0', 'system:notification:list', 'message', 'admin', now(), '', NULL, '消息通知管理菜单');

INSERT INTO `sys_menu` VALUES
(1166, '消息通知查询', 142, 1, '#', '', 1, 0, 'F', '0', '0', 'system:notification:list', '#', 'admin', now(), '', NULL, ''),
(1167, '消息通知编辑', 142, 2, '#', '', 1, 0, 'F', '0', '0', 'system:notification:edit', '#', 'admin', now(), '', NULL, ''),
(1168, '消息通知删除', 142, 3, '#', '', 1, 0, 'F', '0', '0', 'system:notification:remove', '#', 'admin', now(), '', NULL, '');

-- ----------------------------
-- Customer management, portal contact and client workspace
-- ----------------------------

DROP TABLE IF EXISTS `customer_workspace_account_role`;
DROP TABLE IF EXISTS `customer_workspace_role_menu`;
DROP TABLE IF EXISTS `customer_workspace_role`;
DROP TABLE IF EXISTS `customer_workspace_menu`;
DROP TABLE IF EXISTS `customer_account`;
DROP TABLE IF EXISTS `customer_contact`;
DROP TABLE IF EXISTS `customer`;
DROP TABLE IF EXISTS `portal_contact`;

CREATE TABLE `portal_contact` (
  `contact_id` bigint NOT NULL COMMENT '官网线索ID',
  `lead_no` varchar(64) NOT NULL COMMENT '线索编号',
  `contact_name` varchar(64) NOT NULL COMMENT '联系人',
  `company_name` varchar(128) DEFAULT '' COMMENT '公司名称',
  `phone` varchar(64) DEFAULT '' COMMENT '联系电话',
  `email` varchar(128) DEFAULT '' COMMENT '邮箱',
  `route` varchar(255) DEFAULT '' COMMENT '目标航线',
  `cargo_info` varchar(500) DEFAULT '' COMMENT '货物信息',
  `message` text NOT NULL COMMENT '需求说明',
  `source` varchar(64) DEFAULT 'portal-contact' COMMENT '来源',
  `status` char(2) DEFAULT '10' COMMENT '状态（10待跟进 20跟进中 30已完成 40无效）',
  `ip_addr` varchar(128) DEFAULT '' COMMENT '提交IP',
  `user_agent` varchar(255) DEFAULT '' COMMENT '浏览器UA',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `create_by` varchar(64) DEFAULT 'portal' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`contact_id`) USING BTREE,
  UNIQUE KEY `uk_portal_contact_lead_no` (`lead_no`) USING BTREE,
  KEY `idx_portal_contact_phone` (`phone`) USING BTREE,
  KEY `idx_portal_contact_status` (`status`) USING BTREE,
  KEY `idx_portal_contact_create_time` (`create_time`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='官网联系线索表';

CREATE TABLE `customer` (
  `customer_id` bigint NOT NULL COMMENT '客户ID',
  `customer_no` varchar(64) NOT NULL COMMENT '客户编号',
  `customer_name` varchar(128) NOT NULL COMMENT '客户名称',
  `company_name` varchar(128) DEFAULT '' COMMENT '公司名称',
  `contact_name` varchar(64) DEFAULT '' COMMENT '默认联系人',
  `phone` varchar(64) DEFAULT '' COMMENT '联系电话',
  `email` varchar(128) DEFAULT '' COMMENT '邮箱',
  `sales_user_id` bigint DEFAULT 0 COMMENT '负责业务员用户ID',
  `sales_user_name` varchar(64) DEFAULT '' COMMENT '负责业务员名称',
  `status` char(1) DEFAULT '0' COMMENT '状态（0正常 1停用）',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`customer_id`) USING BTREE,
  UNIQUE KEY `uk_customer_no` (`customer_no`) USING BTREE,
  KEY `idx_customer_name` (`customer_name`) USING BTREE,
  KEY `idx_customer_phone` (`phone`) USING BTREE,
  KEY `idx_customer_sales_user` (`sales_user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='客户主体表';

CREATE TABLE `customer_contact` (
  `contact_id` bigint NOT NULL COMMENT '客户联系人ID',
  `customer_id` bigint NOT NULL COMMENT '客户ID',
  `contact_name` varchar(64) NOT NULL COMMENT '联系人姓名',
  `position` varchar(64) DEFAULT '' COMMENT '职务',
  `phone` varchar(64) DEFAULT '' COMMENT '联系电话',
  `email` varchar(128) DEFAULT '' COMMENT '邮箱',
  `wechat` varchar(64) DEFAULT '' COMMENT '微信',
  `is_primary` char(1) DEFAULT '0' COMMENT '是否主要联系人（0否 1是）',
  `status` char(1) DEFAULT '0' COMMENT '状态（0正常 1停用）',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`contact_id`) USING BTREE,
  KEY `idx_customer_contact_customer_id` (`customer_id`) USING BTREE,
  KEY `idx_customer_contact_phone` (`phone`) USING BTREE,
  CONSTRAINT `fk_customer_contact_customer` FOREIGN KEY (`customer_id`) REFERENCES `customer` (`customer_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='客户联系人表';

CREATE TABLE `customer_account` (
  `account_id` bigint NOT NULL COMMENT '客户账号ID',
  `customer_id` bigint NOT NULL COMMENT '客户ID',
  `username` varchar(64) NOT NULL COMMENT '登录账号',
  `password` varchar(255) NOT NULL COMMENT '登录密码',
  `real_name` varchar(64) DEFAULT '' COMMENT '账号姓名',
  `phone` varchar(64) DEFAULT '' COMMENT '手机号',
  `email` varchar(128) DEFAULT '' COMMENT '邮箱',
  `is_main` char(1) DEFAULT '0' COMMENT '是否主账号（0否 1是）',
  `status` char(1) DEFAULT '0' COMMENT '状态（0正常 1停用）',
  `last_login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`account_id`) USING BTREE,
  UNIQUE KEY `uk_customer_account_username` (`username`) USING BTREE,
  KEY `idx_customer_account_customer_id` (`customer_id`) USING BTREE,
  CONSTRAINT `fk_customer_account_customer` FOREIGN KEY (`customer_id`) REFERENCES `customer` (`customer_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='客户登录账号表';

CREATE TABLE `customer_workspace_menu` (
  `menu_id` bigint NOT NULL COMMENT '客户端菜单ID',
  `parent_id` bigint NOT NULL DEFAULT '0' COMMENT '父菜单ID',
  `menu_name` varchar(64) NOT NULL COMMENT '菜单名称',
  `order_num` varchar(8) NOT NULL DEFAULT '0' COMMENT '显示顺序',
  `path` varchar(128) NOT NULL DEFAULT '' COMMENT '路由地址',
  `component` varchar(128) NOT NULL DEFAULT '' COMMENT '组件标识',
  `is_cache` char(1) NOT NULL DEFAULT '0' COMMENT '是否缓存（0缓存 1不缓存）',
  `menu_type` char(1) NOT NULL DEFAULT 'C' COMMENT '菜单类型（M目录 C菜单 F按钮）',
  `visible` char(1) NOT NULL DEFAULT '0' COMMENT '显示状态（0显示 1隐藏）',
  `status` char(1) NOT NULL DEFAULT '0' COMMENT '菜单状态（0正常 1停用）',
  `perms` varchar(128) NOT NULL DEFAULT '' COMMENT '权限标识',
  `icon` varchar(64) NOT NULL DEFAULT '' COMMENT '图标',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`menu_id`) USING BTREE,
  KEY `idx_customer_workspace_menu_parent` (`parent_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='客户端菜单表';

CREATE TABLE `customer_workspace_role` (
  `role_id` bigint NOT NULL COMMENT '客户端角色ID',
  `role_name` varchar(64) NOT NULL COMMENT '角色名称',
  `role_key` varchar(128) NOT NULL COMMENT '权限字符',
  `role_sort` int NOT NULL DEFAULT '0' COMMENT '显示排序',
  `status` char(1) NOT NULL DEFAULT '0' COMMENT '状态（0正常 1停用）',
  `del_flag` char(1) NOT NULL DEFAULT '0' COMMENT '删除标志（0存在 2删除）',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`role_id`) USING BTREE,
  UNIQUE KEY `uk_customer_workspace_role_name` (`role_name`) USING BTREE,
  UNIQUE KEY `uk_customer_workspace_role_key` (`role_key`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='客户端角色表';

CREATE TABLE `customer_workspace_role_menu` (
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `menu_id` bigint NOT NULL COMMENT '菜单ID',
  PRIMARY KEY (`role_id`, `menu_id`) USING BTREE,
  KEY `idx_customer_workspace_role_menu_menu` (`menu_id`) USING BTREE,
  CONSTRAINT `fk_customer_workspace_role_menu_role` FOREIGN KEY (`role_id`) REFERENCES `customer_workspace_role` (`role_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_customer_workspace_role_menu_menu` FOREIGN KEY (`menu_id`) REFERENCES `customer_workspace_menu` (`menu_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='客户端角色菜单关联表';

CREATE TABLE `customer_workspace_account_role` (
  `account_id` bigint NOT NULL COMMENT '客户账号ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  PRIMARY KEY (`account_id`, `role_id`) USING BTREE,
  KEY `idx_customer_workspace_account_role_role` (`role_id`) USING BTREE,
  CONSTRAINT `fk_customer_workspace_account_role_account` FOREIGN KEY (`account_id`) REFERENCES `customer_account` (`account_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_customer_workspace_account_role_role` FOREIGN KEY (`role_id`) REFERENCES `customer_workspace_role` (`role_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='客户端账号角色关联表';

DELETE FROM `sys_dict_data` WHERE `dict_type` = 'portal_contact_status';
DELETE FROM `sys_dict_type` WHERE `dict_type` = 'portal_contact_status';

SET @dict_id := (SELECT IFNULL(MAX(`dict_id`), 0) + 1 FROM `sys_dict_type`);
INSERT INTO `sys_dict_type`
(`dict_id`, `dict_name`, `dict_type`, `status`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
VALUES
(@dict_id, '官网线索跟进状态', 'portal_contact_status', '0', 'admin', now(), '', NULL, '官网联系线索跟进状态');

SET @dict_code := (SELECT IFNULL(MAX(`dict_code`), 0) FROM `sys_dict_data`);
INSERT INTO `sys_dict_data`
(`dict_code`, `dict_sort`, `dict_label`, `dict_value`, `dict_type`, `css_class`, `list_class`, `is_default`, `status`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
VALUES
(@dict_code := @dict_code + 1, 1, '待跟进', '10', 'portal_contact_status', '', 'info', 'Y', '0', 'admin', now(), '', NULL, '待跟进'),
(@dict_code := @dict_code + 1, 2, '跟进中', '20', 'portal_contact_status', '', 'warning', 'N', '0', 'admin', now(), '', NULL, '跟进中'),
(@dict_code := @dict_code + 1, 3, '已完成', '30', 'portal_contact_status', '', 'success', 'N', '0', 'admin', now(), '', NULL, '已完成'),
(@dict_code := @dict_code + 1, 4, '无效', '40', 'portal_contact_status', '', 'danger', 'N', '0', 'admin', now(), '', NULL, '无效');

DELETE FROM `customer_workspace_role_menu` WHERE `role_id` = 20001;
DELETE FROM `customer_workspace_account_role` WHERE `role_id` = 20001;
DELETE FROM `customer_workspace_role` WHERE `role_id` = 20001;
DELETE FROM `customer_workspace_menu` WHERE `menu_id` IN (20001, 20002, 20003, 20004, 20005);

INSERT INTO `customer_workspace_menu` (`menu_id`, `parent_id`, `menu_name`, `order_num`, `path`, `component`, `is_cache`, `menu_type`, `visible`, `status`, `perms`, `icon`, `remark`, `create_by`, `create_time`, `update_by`, `update_time`) VALUES
(20001, 0, '工作台', '1', 'workspace', 'workspace/dashboard', '0', 'C', '0', '0', 'portal:workspace:view', 'AppstoreOutlined', '客户端工作台', 'admin', now(), 'admin', now()),
(20002, 0, '账号资料', '2', 'account', 'workspace/account-profile', '0', 'C', '0', '0', 'portal:account:view', 'ProfileOutlined', '客户端账号资料', 'admin', now(), 'admin', now()),
(20003, 0, '出货查询', '3', 'shipment', 'workspace/shipment-tracking', '0', 'C', '0', '0', 'portal:shipment:view', 'RadarChartOutlined', '客户端出货查询', 'admin', now(), 'admin', now()),
(20004, 0, '智能出货助手', '4', 'shipment-assistant', 'workspace/shipment-assistant', '0', 'C', '0', '0', 'portal:shipmentAssistant:view', 'CalculatorOutlined', '客户端智能出货助手', 'admin', now(), 'admin', now()),
(20005, 0, 'Agent 对话', '5', 'agent-chat', 'workspace/agent-chat', '0', 'C', '0', '0', 'portal:agentChat:view', 'MessageOutlined', '客户端 Agent 对话', 'admin', now(), 'admin', now());

INSERT INTO `customer_workspace_role` (`role_id`, `role_name`, `role_key`, `role_sort`, `status`, `del_flag`, `remark`, `create_by`, `create_time`, `update_by`, `update_time`) VALUES
(20001, '基础客户端角色', 'portal:base', 1, '0', '0', '默认客户端角色', 'admin', now(), 'admin', now());

INSERT INTO `customer_workspace_role_menu` (`role_id`, `menu_id`) VALUES
(20001, 20001),
(20001, 20002),
(20001, 20003),
(20001, 20004),
(20001, 20005);

DELETE FROM `sys_role_menu`
WHERE `menu_id` IN (130, 131, 132, 133, 134, 135, 1130, 1131, 1132, 1133, 1140, 1141, 1142, 1143, 1150, 1151, 1152, 1153, 1154, 1255, 1256, 1257, 1258, 1259, 1260, 1261, 1262, 1263, 1264)
OR `menu_id` IN (
  SELECT `menu_id` FROM `sys_menu`
  WHERE `perms` IN (
    'portal:contact:list', 'portal:contact:query', 'portal:contact:edit', 'portal:contact:remove', 'portal:contact:export',
    'customer:customer:list', 'customer:customer:query', 'customer:customer:add', 'customer:customer:edit', 'customer:customer:remove',
    'customer:account:list', 'customer:account:query', 'customer:account:add', 'customer:account:edit', 'customer:account:remove', 'customer:account:resetPwd',
    'customer:portalMenu:list', 'customer:portalMenu:query', 'customer:portalMenu:add', 'customer:portalMenu:edit', 'customer:portalMenu:remove',
    'customer:portalRole:list', 'customer:portalRole:query', 'customer:portalRole:add', 'customer:portalRole:edit', 'customer:portalRole:remove', 'customer:portalRole:assign'
  )
);

DELETE FROM `sys_menu`
WHERE `menu_id` IN (130, 131, 132, 133, 134, 135, 1130, 1131, 1132, 1133, 1140, 1141, 1142, 1143, 1150, 1151, 1152, 1153, 1154, 1255, 1256, 1257, 1258, 1259, 1260, 1261, 1262, 1263, 1264)
OR `perms` IN (
  'portal:contact:list', 'portal:contact:query', 'portal:contact:edit', 'portal:contact:remove', 'portal:contact:export',
  'customer:customer:list', 'customer:customer:query', 'customer:customer:add', 'customer:customer:edit', 'customer:customer:remove',
  'customer:account:list', 'customer:account:query', 'customer:account:add', 'customer:account:edit', 'customer:account:remove', 'customer:account:resetPwd',
  'customer:portalMenu:list', 'customer:portalMenu:query', 'customer:portalMenu:add', 'customer:portalMenu:edit', 'customer:portalMenu:remove',
  'customer:portalRole:list', 'customer:portalRole:query', 'customer:portalRole:add', 'customer:portalRole:edit', 'customer:portalRole:remove', 'customer:portalRole:assign'
);

INSERT INTO `sys_menu` VALUES
(130, '客户管理', 0, 5, 'customer', NULL, 1, 0, 'M', '0', '0', '', 'peoples', 'admin', now(), '', NULL, '客户管理目录'),
(131, '官网线索', 130, 1, 'contact', 'customer/contact/index', 1, 0, 'C', '0', '0', 'portal:contact:list', 'message', 'admin', now(), '', NULL, '官网线索管理菜单'),
(132, '客户资料', 130, 2, 'customer-list', 'customer/customer/index', 1, 0, 'C', '0', '0', 'customer:customer:list', 'peoples', 'admin', now(), '', NULL, '客户资料管理菜单'),
(133, '客户账号', 130, 3, 'customer-account', 'customer/account/index', 1, 0, 'C', '0', '0', 'customer:account:list', 'user', 'admin', now(), '', NULL, '客户账号管理菜单'),
(134, '客户端菜单', 130, 4, 'portal-menu', 'customer/portalMenu/index', 1, 0, 'C', '0', '0', 'customer:portalMenu:list', 'tree-table', 'admin', now(), '', NULL, '客户端菜单管理'),
(135, '客户端角色', 130, 5, 'portal-role', 'customer/portalRole/index', 1, 0, 'C', '0', '0', 'customer:portalRole:list', 'peoples', 'admin', now(), '', NULL, '客户端角色管理');

INSERT INTO `sys_menu` VALUES
(1130, '官网线索查询', 131, 1, '#', '', 1, 0, 'F', '0', '0', 'portal:contact:query', '#', 'admin', now(), '', NULL, ''),
(1131, '官网线索修改', 131, 2, '#', '', 1, 0, 'F', '0', '0', 'portal:contact:edit', '#', 'admin', now(), '', NULL, ''),
(1132, '官网线索删除', 131, 3, '#', '', 1, 0, 'F', '0', '0', 'portal:contact:remove', '#', 'admin', now(), '', NULL, ''),
(1133, '官网线索导出', 131, 4, '#', '', 1, 0, 'F', '0', '0', 'portal:contact:export', '#', 'admin', now(), '', NULL, ''),
(1140, '客户资料查询', 132, 1, '#', '', 1, 0, 'F', '0', '0', 'customer:customer:query', '#', 'admin', now(), '', NULL, ''),
(1141, '客户资料新增', 132, 2, '#', '', 1, 0, 'F', '0', '0', 'customer:customer:add', '#', 'admin', now(), '', NULL, ''),
(1142, '客户资料修改', 132, 3, '#', '', 1, 0, 'F', '0', '0', 'customer:customer:edit', '#', 'admin', now(), '', NULL, ''),
(1143, '客户资料删除', 132, 4, '#', '', 1, 0, 'F', '0', '0', 'customer:customer:remove', '#', 'admin', now(), '', NULL, ''),
(1150, '客户账号查询', 133, 1, '#', '', 1, 0, 'F', '0', '0', 'customer:account:query', '#', 'admin', now(), '', NULL, ''),
(1151, '客户账号新增', 133, 2, '#', '', 1, 0, 'F', '0', '0', 'customer:account:add', '#', 'admin', now(), '', NULL, ''),
(1152, '客户账号修改', 133, 3, '#', '', 1, 0, 'F', '0', '0', 'customer:account:edit', '#', 'admin', now(), '', NULL, ''),
(1153, '客户账号删除', 133, 4, '#', '', 1, 0, 'F', '0', '0', 'customer:account:remove', '#', 'admin', now(), '', NULL, ''),
(1154, '客户账号重置密码', 133, 5, '#', '', 1, 0, 'F', '0', '0', 'customer:account:resetPwd', '#', 'admin', now(), '', NULL, ''),
(1255, '客户端菜单查询', 134, 1, '#', '', 1, 0, 'F', '0', '0', 'customer:portalMenu:query', '#', 'admin', now(), '', NULL, ''),
(1256, '客户端菜单新增', 134, 2, '#', '', 1, 0, 'F', '0', '0', 'customer:portalMenu:add', '#', 'admin', now(), '', NULL, ''),
(1257, '客户端菜单修改', 134, 3, '#', '', 1, 0, 'F', '0', '0', 'customer:portalMenu:edit', '#', 'admin', now(), '', NULL, ''),
(1258, '客户端菜单删除', 134, 4, '#', '', 1, 0, 'F', '0', '0', 'customer:portalMenu:remove', '#', 'admin', now(), '', NULL, ''),
(1259, '客户端角色查询', 135, 1, '#', '', 1, 0, 'F', '0', '0', 'customer:portalRole:query', '#', 'admin', now(), '', NULL, ''),
(1260, '客户端角色新增', 135, 2, '#', '', 1, 0, 'F', '0', '0', 'customer:portalRole:add', '#', 'admin', now(), '', NULL, ''),
(1261, '客户端角色修改', 135, 3, '#', '', 1, 0, 'F', '0', '0', 'customer:portalRole:edit', '#', 'admin', now(), '', NULL, ''),
(1262, '客户端角色删除', 135, 4, '#', '', 1, 0, 'F', '0', '0', 'customer:portalRole:remove', '#', 'admin', now(), '', NULL, ''),
(1263, '客户端角色分配', 135, 5, '#', '', 1, 0, 'F', '0', '0', 'customer:portalRole:assign', '#', 'admin', now(), '', NULL, ''),
(1264, '客户账号角色配置', 133, 6, '#', '', 1, 0, 'F', '0', '0', 'customer:account:edit', '#', 'admin', now(), '', NULL, '');

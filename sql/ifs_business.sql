-- IFS business modules merged SQL.
-- Includes customer workspace, freight shipment, agent, and notification modules.

-- -----------------------------------------------------------------------------
-- Customer workspace
-- -----------------------------------------------------------------------------

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
(20001, 0, '工作台', '1', 'workspace', 'workspace/dashboard', '0', 'C', '0', '0', 'portal:workspace:view', 'mdi:view-dashboard-outline', '客户端工作台', 'admin', now(), 'admin', now()),
(20002, 0, '账号资料', '2', 'account', 'workspace/account-profile', '0', 'C', '0', '0', 'portal:account:view', 'mdi:account-outline', '客户端账号资料', 'admin', now(), 'admin', now()),
(20003, 0, '出货查询', '3', 'shipment', 'workspace/shipment-tracking', '0', 'C', '0', '0', 'portal:shipment:view', 'mdi:radar', '客户端出货查询', 'admin', now(), 'admin', now()),
(20004, 0, '智能出货助手', '4', 'shipment-assistant', 'workspace/shipment-assistant', '0', 'C', '0', '0', 'portal:shipmentAssistant:view', 'mdi:calculator-variant-outline', '客户端智能出货助手', 'admin', now(), 'admin', now()),
(20005, 0, 'Agent 对话', '5', 'agent-chat', 'workspace/agent-chat', '0', 'C', '0', '0', 'portal:agentChat:view', 'mdi:message-text-outline', '客户端 Agent 对话', 'admin', now(), 'admin', now());

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


-- -----------------------------------------------------------------------------
-- Freight shipment
-- -----------------------------------------------------------------------------

-- ----------------------------
-- Freight shipment planning
-- ----------------------------

DROP TABLE IF EXISTS `freight_shipment_order`;
DROP TABLE IF EXISTS `freight_receipt_allocation`;
DROP TABLE IF EXISTS `freight_receipt`;
DROP TABLE IF EXISTS `freight_payment_declaration`;
DROP TABLE IF EXISTS `freight_shipment_payment`;
DROP TABLE IF EXISTS `freight_container_plan`;
DROP TABLE IF EXISTS `freight_shipment_cargo`;
DROP TABLE IF EXISTS `freight_shipment_plan`;

CREATE TABLE `freight_shipment_plan` (
  `shipment_id` bigint NOT NULL COMMENT '出货计划ID',
  `shipment_no` varchar(64) NOT NULL COMMENT '出货计划编号',
  `order_no` varchar(64) DEFAULT '' COMMENT '客户订单号或参考号',
  `customer_id` bigint NOT NULL COMMENT '客户ID',
  `customer_name` varchar(128) DEFAULT '' COMMENT '客户名称',
  `sales_user_id` bigint DEFAULT 0 COMMENT '负责业务员用户ID快照',
  `sales_user_name` varchar(64) DEFAULT '' COMMENT '负责业务员名称快照',
  `pol` varchar(128) DEFAULT '' COMMENT '起运港',
  `pod` varchar(128) DEFAULT '' COMMENT '目的港',
  `planned_etd` varchar(32) DEFAULT '' COMMENT '计划开船日期',
  `planned_eta` varchar(32) DEFAULT '' COMMENT '计划到港日期',
  `actual_etd` varchar(32) DEFAULT '' COMMENT '实际开船日期',
  `actual_eta` varchar(32) DEFAULT '' COMMENT '实际到港日期',
  `status` varchar(8) DEFAULT '10' COMMENT '状态字典 freight_shipment_status',
  `payment_status` varchar(16) DEFAULT 'UNPAID' COMMENT '付款状态 UNPAID/PARTIAL/PAID',
  `payment_amount` decimal(12,2) DEFAULT 0.00 COMMENT '付款金额',
  `total_weight` decimal(12,2) DEFAULT 0.00 COMMENT '总重量KG',
  `total_volume` decimal(12,2) DEFAULT 0.00 COMMENT '总体积CBM',
  `total_cartons` int DEFAULT 0 COMMENT '总箱数',
  `share_token` varchar(64) NOT NULL COMMENT '免登录分享令牌',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`shipment_id`) USING BTREE,
  UNIQUE KEY `uk_freight_shipment_no` (`shipment_no`) USING BTREE,
  UNIQUE KEY `uk_freight_share_token` (`share_token`) USING BTREE,
  KEY `idx_freight_customer_id` (`customer_id`) USING BTREE,
  KEY `idx_freight_sales_user` (`sales_user_id`) USING BTREE,
  KEY `idx_freight_status` (`status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='出货计划表';

CREATE TABLE `freight_shipment_cargo` (
  `cargo_id` bigint NOT NULL COMMENT '货物明细ID',
  `shipment_id` bigint NOT NULL COMMENT '出货计划ID',
  `sku` varchar(64) DEFAULT '' COMMENT 'SKU/唛头/货号',
  `cargo_name` varchar(128) NOT NULL COMMENT '货物名称',
  `package_type` varchar(32) DEFAULT '' COMMENT '包装类型',
  `quantity` int DEFAULT 0 COMMENT '件数',
  `cartons` int DEFAULT 0 COMMENT '箱数',
  `weight_kg` decimal(12,2) DEFAULT 0.00 COMMENT '重量KG',
  `volume_cbm` decimal(12,2) DEFAULT 0.00 COMMENT '体积CBM',
  `length_cm` decimal(12,2) DEFAULT 0.00 COMMENT '长CM',
  `width_cm` decimal(12,2) DEFAULT 0.00 COMMENT '宽CM',
  `height_cm` decimal(12,2) DEFAULT 0.00 COMMENT '高CM',
  PRIMARY KEY (`cargo_id`) USING BTREE,
  KEY `idx_freight_cargo_shipment_id` (`shipment_id`) USING BTREE,
  CONSTRAINT `fk_freight_cargo_shipment` FOREIGN KEY (`shipment_id`) REFERENCES `freight_shipment_plan` (`shipment_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='出货货物明细表';

CREATE TABLE `freight_shipment_payment` (
  `payment_id` bigint NOT NULL COMMENT '付款记录ID',
  `shipment_id` bigint NOT NULL COMMENT '出货计划ID',
  `amount` decimal(12,2) NOT NULL COMMENT '付款金额',
  `currency` varchar(8) DEFAULT 'CNY' COMMENT '币种',
  `payment_time` datetime NOT NULL COMMENT '付款时间',
  `payment_method` varchar(32) DEFAULT '' COMMENT '付款方式',
  `voucher_url` varchar(500) DEFAULT '' COMMENT '付款凭证地址',
  `voucher_name` varchar(255) DEFAULT '' COMMENT '付款凭证原文件名',
  `remark` varchar(500) DEFAULT '' COMMENT '付款备注',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`payment_id`) USING BTREE,
  KEY `idx_freight_payment_shipment_id` (`shipment_id`) USING BTREE,
  CONSTRAINT `fk_freight_payment_shipment` FOREIGN KEY (`shipment_id`) REFERENCES `freight_shipment_plan` (`shipment_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='出货计划付款记录表';

CREATE TABLE `freight_receipt` (
  `receipt_id` bigint NOT NULL, `receipt_no` varchar(64) NOT NULL, `customer_id` bigint NOT NULL, `customer_name` varchar(128) DEFAULT '',
  `amount` decimal(12,2) NOT NULL, `currency` varchar(8) DEFAULT 'CNY', `receipt_time` datetime NOT NULL,
  `payment_method` varchar(32) DEFAULT '', `status` varchar(16) DEFAULT 'UNALLOCATED',
  `voucher_url` varchar(500) DEFAULT '', `voucher_name` varchar(255) DEFAULT '', `remark` varchar(500) DEFAULT '',
  `create_by` varchar(64) DEFAULT '', `create_time` datetime DEFAULT NULL, `update_by` varchar(64) DEFAULT '', `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`receipt_id`), UNIQUE KEY `uk_freight_receipt_no` (`receipt_no`), KEY `idx_freight_receipt_customer` (`customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='独立收款单';

CREATE TABLE `freight_receipt_allocation` (
  `allocation_id` bigint NOT NULL, `receipt_id` bigint NOT NULL, `shipment_id` bigint NOT NULL, `allocated_amount` decimal(12,2) NOT NULL,
  PRIMARY KEY (`allocation_id`), UNIQUE KEY `uk_receipt_shipment` (`receipt_id`,`shipment_id`), KEY `idx_allocation_shipment` (`shipment_id`),
  CONSTRAINT `fk_allocation_receipt` FOREIGN KEY (`receipt_id`) REFERENCES `freight_receipt` (`receipt_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_allocation_shipment` FOREIGN KEY (`shipment_id`) REFERENCES `freight_shipment_plan` (`shipment_id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='收款核销明细';

CREATE TABLE `freight_payment_declaration` (
  `declaration_id` bigint NOT NULL, `declaration_no` varchar(64) NOT NULL, `customer_id` bigint NOT NULL, `customer_name` varchar(128) DEFAULT '',
  `shipment_id` bigint NOT NULL, `shipment_no` varchar(64) DEFAULT '', `amount` decimal(12,2) NOT NULL, `currency` varchar(8) DEFAULT 'CNY',
  `payment_time` varchar(32) DEFAULT '', `voucher_url` varchar(500) NOT NULL, `voucher_name` varchar(255) DEFAULT '',
  `status` varchar(16) DEFAULT 'PENDING', `remark` varchar(500) DEFAULT '', `review_by` varchar(64) DEFAULT '', `review_time` datetime DEFAULT NULL, `review_remark` varchar(500) DEFAULT '', `create_by` varchar(64) DEFAULT '', `create_time` datetime DEFAULT NULL,
  `update_by` varchar(64) DEFAULT '', `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`declaration_id`), UNIQUE KEY `uk_payment_declaration_no` (`declaration_no`), KEY `idx_declaration_customer` (`customer_id`), KEY `idx_declaration_shipment` (`shipment_id`),
  CONSTRAINT `fk_declaration_shipment` FOREIGN KEY (`shipment_id`) REFERENCES `freight_shipment_plan` (`shipment_id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='客户付款申报';

CREATE TABLE `freight_container_plan` (
  `container_plan_id` bigint NOT NULL COMMENT '货柜计划ID',
  `shipment_id` bigint NOT NULL COMMENT '出货计划ID',
  `container_type` varchar(16) NOT NULL COMMENT '柜型',
  `quantity` int DEFAULT 1 COMMENT '柜量',
  `max_volume` decimal(12,2) DEFAULT 0.00 COMMENT '可装体积CBM',
  `max_weight` decimal(12,2) DEFAULT 0.00 COMMENT '可装重量KG',
  `used_volume` decimal(12,2) DEFAULT 0.00 COMMENT '已用体积CBM',
  `used_weight` decimal(12,2) DEFAULT 0.00 COMMENT '已用重量KG',
  `load_rate` decimal(8,2) DEFAULT 0.00 COMMENT '装载率',
  `remark` varchar(255) DEFAULT '' COMMENT '推荐说明',
  PRIMARY KEY (`container_plan_id`) USING BTREE,
  KEY `idx_freight_container_shipment_id` (`shipment_id`) USING BTREE,
  CONSTRAINT `fk_freight_container_shipment` FOREIGN KEY (`shipment_id`) REFERENCES `freight_shipment_plan` (`shipment_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='智能货柜计划表';

CREATE TABLE `freight_shipment_order` (
  `order_id` bigint NOT NULL COMMENT '出货单ID',
  `shipment_id` bigint NOT NULL COMMENT '出货计划ID',
  `order_no` varchar(64) NOT NULL COMMENT '出货单号',
  `status` varchar(8) DEFAULT '10' COMMENT '状态',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`order_id`) USING BTREE,
  UNIQUE KEY `uk_freight_order_no` (`order_no`) USING BTREE,
  UNIQUE KEY `uk_freight_order_shipment` (`shipment_id`) USING BTREE,
  CONSTRAINT `fk_freight_order_shipment` FOREIGN KEY (`shipment_id`) REFERENCES `freight_shipment_plan` (`shipment_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='出货单表';

DELETE FROM `sys_dict_data` WHERE `dict_type` IN ('freight_shipment_status', 'freight_container_type');
DELETE FROM `sys_dict_type` WHERE `dict_type` IN ('freight_shipment_status', 'freight_container_type');

SET @dict_id := (SELECT IFNULL(MAX(`dict_id`), 0) + 1 FROM `sys_dict_type`);
INSERT INTO `sys_dict_type`
(`dict_id`, `dict_name`, `dict_type`, `status`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
VALUES
(@dict_id, '出货状态', 'freight_shipment_status', '0', 'admin', now(), '', NULL, '客户端分享页展示状态');

SET @dict_code := (SELECT IFNULL(MAX(`dict_code`), 0) FROM `sys_dict_data`);
INSERT INTO `sys_dict_data`
(`dict_code`, `dict_sort`, `dict_label`, `dict_value`, `dict_type`, `css_class`, `list_class`, `is_default`, `status`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
VALUES
(@dict_code := @dict_code + 1, 1, '计划已创建', '10', 'freight_shipment_status', '', 'primary', 'N', '0', 'admin', now(), '', NULL, '内部已建立出货计划'),
(@dict_code := @dict_code + 1, 2, '出货计划已确认', '20', 'freight_shipment_status', '', 'primary', 'N', '0', 'admin', now(), '', NULL, '客户或运营确认计划'),
(@dict_code := @dict_code + 1, 3, '等待客户发货', '30', 'freight_shipment_status', '', 'info', 'N', '0', 'admin', now(), '', NULL, '等待客户送货或安排提货'),
(@dict_code := @dict_code + 1, 4, '已提货/已送仓', '40', 'freight_shipment_status', '', 'warning', 'N', '0', 'admin', now(), '', NULL, '货物已离开发货地'),
(@dict_code := @dict_code + 1, 5, '仓库已收货', '50', 'freight_shipment_status', '', 'warning', 'N', '0', 'admin', now(), '', NULL, '仓库完成收货'),
(@dict_code := @dict_code + 1, 6, '已入仓/码头进仓', '60', 'freight_shipment_status', '', 'warning', 'N', '0', 'admin', now(), '', NULL, '仓库或码头完成进仓'),
(@dict_code := @dict_code + 1, 7, '订舱处理中', '70', 'freight_shipment_status', '', 'warning', 'N', '0', 'admin', now(), '', NULL, '正在向船司或渠道订舱'),
(@dict_code := @dict_code + 1, 8, '舱位已确认', '80', 'freight_shipment_status', '', 'primary', 'N', '0', 'admin', now(), '', NULL, '订舱确认'),
(@dict_code := @dict_code + 1, 9, '报关资料已收齐', '90', 'freight_shipment_status', '', 'primary', 'N', '0', 'admin', now(), '', NULL, '报关文件齐备'),
(@dict_code := @dict_code + 1, 10, '报关已放行', '100', 'freight_shipment_status', '', 'success', 'N', '0', 'admin', now(), '', NULL, '出口报关放行'),
(@dict_code := @dict_code + 1, 11, '已装柜', '110', 'freight_shipment_status', '', 'success', 'N', '0', 'admin', now(), '', NULL, '完成装柜'),
(@dict_code := @dict_code + 1, 12, '已进港/码头放行', '120', 'freight_shipment_status', '', 'success', 'N', '0', 'admin', now(), '', NULL, '已进港或码头放行'),
(@dict_code := @dict_code + 1, 13, '船舶已开船', '130', 'freight_shipment_status', '', 'success', 'N', '0', 'admin', now(), '', NULL, '船已开船'),
(@dict_code := @dict_code + 1, 14, '目的港已到港', '140', 'freight_shipment_status', '', 'success', 'N', '0', 'admin', now(), '', NULL, '船已到目的港'),
(@dict_code := @dict_code + 1, 15, '目的港清关中', '150', 'freight_shipment_status', '', 'warning', 'N', '0', 'admin', now(), '', NULL, '目的港清关处理中'),
(@dict_code := @dict_code + 1, 16, '目的港已清关', '160', 'freight_shipment_status', '', 'success', 'N', '0', 'admin', now(), '', NULL, '目的港清关完成'),
(@dict_code := @dict_code + 1, 17, '已派送/已签收', '170', 'freight_shipment_status', '', 'success', 'N', '0', 'admin', now(), '', NULL, '派送中或客户已签收'),
(@dict_code := @dict_code + 1, 99, '异常处理中', '900', 'freight_shipment_status', '', 'danger', 'N', '0', 'admin', now(), '', NULL, '状态异常，运营正在处理');

SET @dict_id := (SELECT IFNULL(MAX(`dict_id`), 0) + 1 FROM `sys_dict_type`);
INSERT INTO `sys_dict_type`
(`dict_id`, `dict_name`, `dict_type`, `status`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
VALUES
(@dict_id, '货柜类型', 'freight_container_type', '0', 'admin', now(), '', NULL, '智能出货计划柜型');

SET @dict_code := (SELECT IFNULL(MAX(`dict_code`), 0) FROM `sys_dict_data`);
INSERT INTO `sys_dict_data`
(`dict_code`, `dict_sort`, `dict_label`, `dict_value`, `dict_type`, `css_class`, `list_class`, `is_default`, `status`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
VALUES
(@dict_code := @dict_code + 1, 1, '20GP', '20GP', 'freight_container_type', '', 'default', 'N', '0', 'admin', now(), '', NULL, '约28CBM/21700KG'),
(@dict_code := @dict_code + 1, 2, '40GP', '40GP', 'freight_container_type', '', 'default', 'N', '0', 'admin', now(), '', NULL, '约58CBM/26500KG'),
(@dict_code := @dict_code + 1, 3, '40HQ', '40HQ', 'freight_container_type', '', 'default', 'Y', '0', 'admin', now(), '', NULL, '约68CBM/26500KG'),
(@dict_code := @dict_code + 1, 4, 'LCL 拼箱', 'LCL', 'freight_container_type', '', 'info', 'N', '0', 'admin', now(), '', NULL, '小票货优先按散货拼箱评估');

DELETE FROM `sys_role_menu`
WHERE `menu_id` IN (141, 1160, 1161, 1162, 1163, 1164, 1165)
OR `menu_id` IN (
  SELECT `menu_id` FROM `sys_menu`
  WHERE `perms` IN (
    'freight:shipment:list',
    'freight:shipment:query',
    'freight:shipment:import',
    'freight:shipment:edit',
    'freight:shipment:confirm',
    'freight:shipment:share',
    'freight:shipment:remove'
  )
);

DELETE FROM `sys_menu`
WHERE `menu_id` IN (141, 1160, 1161, 1162, 1163, 1164, 1165)
OR `perms` IN (
  'freight:shipment:list',
  'freight:shipment:query',
  'freight:shipment:import',
  'freight:shipment:edit',
  'freight:shipment:confirm',
  'freight:shipment:share',
  'freight:shipment:remove'
);

INSERT INTO `sys_menu` (
  `menu_id`, `menu_name`, `parent_id`, `order_num`, `path`, `component`,
  `is_frame`, `is_cache`, `menu_type`, `visible`, `status`, `perms`, `icon`,
  `create_by`, `create_time`, `update_by`, `update_time`, `remark`
)
SELECT
  140, '货代业务', 0, 4, 'freight', NULL,
  1, 0, 'M', '0', '0', '', 'mdi:earth',
  'admin', NOW(), '', NULL, '国际货代业务菜单'
FROM dual
WHERE NOT EXISTS (
  SELECT 1 FROM `sys_menu` WHERE `menu_id` = 140
);

INSERT INTO `sys_menu` VALUES
(141, '出货计划', 140, 1, 'shipment', 'freight/shipment/index', 1, 0, 'C', '0', '0', 'freight:shipment:list', 'mdi:truck-cargo-container', 'admin', now(), '', NULL, '出货计划管理菜单');

INSERT INTO `sys_menu` VALUES
(1160, '出货计划查询', 141, 1, '#', '', 1, 0, 'F', '0', '0', 'freight:shipment:query', '#', 'admin', now(), '', NULL, ''),
(1161, '出货清单导入', 141, 2, '#', '', 1, 0, 'F', '0', '0', 'freight:shipment:import', '#', 'admin', now(), '', NULL, ''),
(1162, '出货状态维护', 141, 3, '#', '', 1, 0, 'F', '0', '0', 'freight:shipment:edit', '#', 'admin', now(), '', NULL, ''),
(1163, '生成出货单', 141, 4, '#', '', 1, 0, 'F', '0', '0', 'freight:shipment:confirm', '#', 'admin', now(), '', NULL, ''),
(1164, '分享出货进度', 141, 5, '#', '', 1, 0, 'F', '0', '0', 'freight:shipment:share', '#', 'admin', now(), '', NULL, ''),
(1165, '出货计划删除', 141, 6, '#', '', 1, 0, 'F', '0', '0', 'freight:shipment:remove', '#', 'admin', now(), '', NULL, '');


-- -----------------------------------------------------------------------------
-- IFS agent
-- -----------------------------------------------------------------------------
INSERT IGNORE INTO `sys_menu` VALUES
(143, '收款管理', 140, 2, 'receipt', 'freight/receipt/index', 1, 0, 'C', '0', '0', 'freight:receipt:list', 'mdi:cash-multiple', 'admin', now(), '', NULL, '独立收款与出货计划核销'),
(1180, '收款查询', 143, 1, '#', '', 1, 0, 'F', '0', '0', 'freight:receipt:query', '#', 'admin', now(), '', NULL, ''),
(1181, '收款新增', 143, 2, '#', '', 1, 0, 'F', '0', '0', 'freight:receipt:add', '#', 'admin', now(), '', NULL, ''),
(1182, '收款删除', 143, 3, '#', '', 1, 0, 'F', '0', '0', 'freight:receipt:remove', '#', 'admin', now(), '', NULL, '');

INSERT IGNORE INTO `sys_menu` VALUES
(144, '付款申报', 140, 3, 'payment-declaration', 'freight/paymentDeclaration/index', 1, 0, 'C', '0', '0', 'freight:declaration:list', 'mdi:receipt-text-check-outline', 'admin', now(), '', NULL, '客户付款申报审核'),
(1190, '申报查询', 144, 1, '#', '', 1, 0, 'F', '0', '0', 'freight:declaration:query', '#', 'admin', now(), '', NULL, ''),
(1191, '申报审核', 144, 2, '#', '', 1, 0, 'F', '0', '0', 'freight:declaration:review', '#', 'admin', now(), '', NULL, '');

CREATE TABLE IF NOT EXISTS `chat_session` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint DEFAULT 0 COMMENT '用户ID；免登录场景可为空',
  `title` varchar(255) DEFAULT '' COMMENT '会话标题',
  `model_name` varchar(100) DEFAULT 'qwen2.5:7b' COMMENT '模型名称',
  `summary` text NULL COMMENT '会话摘要',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_chat_session_user_id` (`user_id`),
  KEY `idx_chat_session_updated_at` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='IFS Agent 会话表';

CREATE TABLE IF NOT EXISTS `chat_message` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `session_id` bigint NOT NULL COMMENT '会话ID',
  `role` varchar(20) NOT NULL COMMENT 'user/assistant/system',
  `content` longtext NOT NULL COMMENT '消息文本',
  `block_result` json NULL COMMENT 'IFS Block Protocol 返回结果',
  `model_name` varchar(100) DEFAULT '' COMMENT '模型名称',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_chat_message_session_id` (`session_id`),
  KEY `idx_chat_message_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='IFS Agent 消息表';

CREATE TABLE IF NOT EXISTS `chat_memory` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint DEFAULT 0 COMMENT '用户ID',
  `memory_key` varchar(100) NOT NULL COMMENT '记忆键',
  `memory_value` text NOT NULL COMMENT '记忆内容',
  `memory_type` varchar(50) DEFAULT 'profile' COMMENT '记忆类型',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_chat_memory_user_key` (`user_id`, `memory_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='IFS Agent 记忆表';

CREATE TABLE IF NOT EXISTS `agent_form_submission` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `session_id` bigint NOT NULL COMMENT '会话ID',
  `form_code` varchar(100) NOT NULL COMMENT '表单编码',
  `form_values` json NOT NULL COMMENT '表单提交值',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_agent_form_submission_session_id` (`session_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='IFS Agent 动态表单提交表';

INSERT INTO `sys_menu` (
  `menu_id`, `menu_name`, `parent_id`, `order_num`, `path`, `component`,
  `is_frame`, `is_cache`, `menu_type`, `visible`, `status`, `perms`, `icon`,
  `create_by`, `create_time`, `update_by`, `update_time`, `remark`
)
SELECT
  140, '货代业务', 0, 4, 'freight', NULL,
  1, 0, 'M', '0', '0', '', 'mdi:earth',
  'admin', NOW(), '', NULL, '国际货代业务菜单'
FROM dual
WHERE NOT EXISTS (
  SELECT 1 FROM `sys_menu` WHERE `menu_id` = 140
);

INSERT INTO `sys_menu` (
  `menu_id`, `menu_name`, `parent_id`, `order_num`, `path`, `component`,
  `is_frame`, `is_cache`, `menu_type`, `visible`, `status`, `perms`, `icon`,
  `create_by`, `create_time`, `update_by`, `update_time`, `remark`
)
SELECT
  142, 'Agent 对话', 140, 2, 'agent-chat', 'agent/chat/index',
  1, 0, 'C', '0', '0', 'ifs:agent:chat', 'message',
  'admin', NOW(), '', NULL, 'IFS Agent 对话管理菜单'
FROM dual
WHERE NOT EXISTS (
  SELECT 1 FROM `sys_menu` WHERE `menu_id` = 142
);

INSERT INTO `sys_menu` (
  `menu_id`, `menu_name`, `parent_id`, `order_num`, `path`, `component`,
  `is_frame`, `is_cache`, `menu_type`, `visible`, `status`, `perms`, `icon`,
  `create_by`, `create_time`, `update_by`, `update_time`, `remark`
)
SELECT
  1170, 'Agent 对话权限', 142, 1, '#', '',
  1, 0, 'F', '0', '0', 'ifs:agent:chat', '#',
  'admin', NOW(), '', NULL, ''
FROM dual
WHERE NOT EXISTS (
  SELECT 1 FROM `sys_menu` WHERE `menu_id` = 1170
);

INSERT INTO `sys_role_menu` (`role_id`, `menu_id`)
SELECT 1, m.`menu_id`
FROM `sys_menu` m
WHERE m.`menu_id` IN (140, 142, 1170)
  AND EXISTS (SELECT 1 FROM `sys_role` r WHERE r.`role_id` = 1)
  AND NOT EXISTS (
    SELECT 1 FROM `sys_role_menu` rm WHERE rm.`role_id` = 1 AND rm.`menu_id` = m.`menu_id`
  );


-- -----------------------------------------------------------------------------
-- Shipment notification
-- -----------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS `sys_notification` (
  `notification_id` bigint NOT NULL COMMENT '通知ID',
  `user_id` bigint NOT NULL COMMENT '接收用户ID',
  `title` varchar(120) NOT NULL COMMENT '标题',
  `content` varchar(500) NOT NULL COMMENT '内容',
  `biz_type` varchar(50) DEFAULT '' COMMENT '业务类型',
  `biz_id` bigint DEFAULT 0 COMMENT '业务ID',
  `read_flag` char(1) NOT NULL DEFAULT '0' COMMENT '是否已读（0未读 1已读）',
  `read_time` datetime DEFAULT NULL COMMENT '已读时间',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`notification_id`),
  KEY `idx_notification_user_read` (`user_id`, `read_flag`, `create_time`),
  KEY `idx_notification_biz` (`biz_type`, `biz_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户消息通知表';


-- -----------------------------------------------------------------------------
-- System notification menu
-- -----------------------------------------------------------------------------

DELETE FROM `sys_role_menu`
WHERE `menu_id` IN (143, 1171, 1172, 1173)
OR `menu_id` IN (
  SELECT `menu_id` FROM `sys_menu`
  WHERE `perms` IN (
    'system:notification:list',
    'system:notification:edit',
    'system:notification:remove'
  )
);

DELETE FROM `sys_menu`
WHERE `menu_id` IN (143, 1171, 1172, 1173)
OR `perms` IN (
  'system:notification:list',
  'system:notification:edit',
  'system:notification:remove'
);

INSERT INTO `sys_menu` VALUES
(143, '消息通知', 1, 8, 'notification', 'system/notification/index', 1, 0, 'C', '0', '0', 'system:notification:list', 'message', 'admin', now(), '', NULL, '消息通知管理菜单');

INSERT INTO `sys_menu` VALUES
(1171, '消息通知查询', 143, 1, '#', '', 1, 0, 'F', '0', '0', 'system:notification:list', '#', 'admin', now(), '', NULL, ''),
(1172, '消息通知编辑', 143, 2, '#', '', 1, 0, 'F', '0', '0', 'system:notification:edit', '#', 'admin', now(), '', NULL, ''),
(1173, '消息通知删除', 143, 3, '#', '', 1, 0, 'F', '0', '0', 'system:notification:remove', '#', 'admin', now(), '', NULL, '');



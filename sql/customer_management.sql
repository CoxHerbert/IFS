-- ----------------------------
-- Customer management tables
-- ----------------------------
DROP TABLE IF EXISTS `customer_account`;
DROP TABLE IF EXISTS `customer_contact`;
DROP TABLE IF EXISTS `customer`;

CREATE TABLE `customer` (
  `customer_id` bigint NOT NULL COMMENT '客户ID',
  `customer_no` varchar(64) NOT NULL COMMENT '客户编号',
  `customer_name` varchar(128) NOT NULL COMMENT '客户名称',
  `company_name` varchar(128) DEFAULT '' COMMENT '公司名称',
  `contact_name` varchar(64) DEFAULT '' COMMENT '默认联系人',
  `phone` varchar(64) DEFAULT '' COMMENT '联系电话',
  `email` varchar(128) DEFAULT '' COMMENT '邮箱',
  `status` char(1) DEFAULT '0' COMMENT '状态（0正常 1停用）',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`customer_id`) USING BTREE,
  UNIQUE KEY `uk_customer_no` (`customer_no`) USING BTREE,
  KEY `idx_customer_name` (`customer_name`) USING BTREE,
  KEY `idx_customer_phone` (`phone`) USING BTREE
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

-- ----------------------------
-- Backend menu data
-- Adjust menu_id values if they conflict in your database.
-- ----------------------------
INSERT INTO `sys_menu` VALUES (132, '客户资料', 130, 2, 'customer-list', 'customer/customer/index', 1, 0, 'C', '0', '0', 'customer:customer:list', 'peoples', 'admin', now(), '', NULL, '客户资料管理菜单');
INSERT INTO `sys_menu` VALUES (133, '客户账号', 130, 3, 'customer-account', 'customer/account/index', 1, 0, 'C', '0', '0', 'customer:account:list', 'user', 'admin', now(), '', NULL, '客户账号管理菜单');

INSERT INTO `sys_menu` VALUES (1140, '客户资料查询', 132, 1, '#', '', 1, 0, 'F', '0', '0', 'customer:customer:query', '#', 'admin', now(), '', NULL, '');
INSERT INTO `sys_menu` VALUES (1141, '客户资料新增', 132, 2, '#', '', 1, 0, 'F', '0', '0', 'customer:customer:add', '#', 'admin', now(), '', NULL, '');
INSERT INTO `sys_menu` VALUES (1142, '客户资料修改', 132, 3, '#', '', 1, 0, 'F', '0', '0', 'customer:customer:edit', '#', 'admin', now(), '', NULL, '');
INSERT INTO `sys_menu` VALUES (1143, '客户资料删除', 132, 4, '#', '', 1, 0, 'F', '0', '0', 'customer:customer:remove', '#', 'admin', now(), '', NULL, '');

INSERT INTO `sys_menu` VALUES (1150, '客户账号查询', 133, 1, '#', '', 1, 0, 'F', '0', '0', 'customer:account:query', '#', 'admin', now(), '', NULL, '');
INSERT INTO `sys_menu` VALUES (1151, '客户账号新增', 133, 2, '#', '', 1, 0, 'F', '0', '0', 'customer:account:add', '#', 'admin', now(), '', NULL, '');
INSERT INTO `sys_menu` VALUES (1152, '客户账号修改', 133, 3, '#', '', 1, 0, 'F', '0', '0', 'customer:account:edit', '#', 'admin', now(), '', NULL, '');
INSERT INTO `sys_menu` VALUES (1153, '客户账号删除', 133, 4, '#', '', 1, 0, 'F', '0', '0', 'customer:account:remove', '#', 'admin', now(), '', NULL, '');
INSERT INTO `sys_menu` VALUES (1154, '客户账号重置密码', 133, 5, '#', '', 1, 0, 'F', '0', '0', 'customer:account:resetPwd', '#', 'admin', now(), '', NULL, '');

-- ----------------------------
-- Table structure for portal_contact
-- ----------------------------
DROP TABLE IF EXISTS `portal_contact`;
CREATE TABLE `portal_contact` (
  `contact_id` bigint NOT NULL COMMENT '联系线索ID',
  `lead_no` varchar(64) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL COMMENT '线索编号',
  `contact_name` varchar(64) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL COMMENT '联系人',
  `company_name` varchar(128) CHARACTER SET utf8 COLLATE utf8_unicode_ci NULL DEFAULT '' COMMENT '公司名称',
  `phone` varchar(64) CHARACTER SET utf8 COLLATE utf8_unicode_ci NULL DEFAULT '' COMMENT '联系电话',
  `email` varchar(128) CHARACTER SET utf8 COLLATE utf8_unicode_ci NULL DEFAULT '' COMMENT '邮箱',
  `route` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NULL DEFAULT '' COMMENT '目标航线',
  `cargo_info` varchar(500) CHARACTER SET utf8 COLLATE utf8_unicode_ci NULL DEFAULT '' COMMENT '货物信息',
  `message` text CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL COMMENT '需求说明',
  `source` varchar(64) CHARACTER SET utf8 COLLATE utf8_unicode_ci NULL DEFAULT 'portal-contact' COMMENT '来源',
  `status` char(2) CHARACTER SET utf8 COLLATE utf8_unicode_ci NULL DEFAULT '10' COMMENT '状态（10待跟进 20跟进中 30已完成 40无效）',
  `ip_addr` varchar(128) CHARACTER SET utf8 COLLATE utf8_unicode_ci NULL DEFAULT '' COMMENT '提交IP',
  `user_agent` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NULL DEFAULT '' COMMENT '浏览器UA',
  `remark` varchar(500) CHARACTER SET utf8 COLLATE utf8_unicode_ci NULL DEFAULT '' COMMENT '备注',
  `create_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_unicode_ci NULL DEFAULT 'portal' COMMENT '创建者',
  `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_unicode_ci NULL DEFAULT '' COMMENT '更新者',
  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`contact_id`) USING BTREE,
  UNIQUE INDEX `uk_portal_contact_lead_no` (`lead_no`) USING BTREE,
  INDEX `idx_portal_contact_phone` (`phone`) USING BTREE,
  INDEX `idx_portal_contact_status` (`status`) USING BTREE,
  INDEX `idx_portal_contact_create_time` (`create_time`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_unicode_ci COMMENT = '官网联系我们线索表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Menu data for portal_contact
-- 后台菜单：客户管理 -> 官网线索；如菜单ID冲突，请按现场 sys_menu 最大ID调整。
-- ----------------------------
INSERT INTO `sys_menu` VALUES (130, '客户管理', 0, 5, 'customer', NULL, 1, 0, 'M', '0', '0', '', 'peoples', 'admin', now(), '', NULL, '客户管理目录');
INSERT INTO `sys_menu` VALUES (131, '官网线索', 130, 1, 'contact', 'customer/contact/index', 1, 0, 'C', '0', '0', 'portal:contact:list', 'message', 'admin', now(), '', NULL, '官网联系我们线索管理菜单');
INSERT INTO `sys_menu` VALUES (1130, '官网线索查询', 131, 1, '#', '', 1, 0, 'F', '0', '0', 'portal:contact:query', '#', 'admin', now(), '', NULL, '');
INSERT INTO `sys_menu` VALUES (1131, '官网线索修改', 131, 2, '#', '', 1, 0, 'F', '0', '0', 'portal:contact:edit', '#', 'admin', now(), '', NULL, '');
INSERT INTO `sys_menu` VALUES (1132, '官网线索删除', 131, 3, '#', '', 1, 0, 'F', '0', '0', 'portal:contact:remove', '#', 'admin', now(), '', NULL, '');
INSERT INTO `sys_menu` VALUES (1133, '官网线索导出', 131, 4, '#', '', 1, 0, 'F', '0', '0', 'portal:contact:export', '#', 'admin', now(), '', NULL, '');

-- 如需给普通角色授权，可执行以下语句。
-- INSERT INTO `sys_role_menu` VALUES (2, 130);
-- INSERT INTO `sys_role_menu` VALUES (2, 131);
-- INSERT INTO `sys_role_menu` VALUES (2, 1130);
-- INSERT INTO `sys_role_menu` VALUES (2, 1131);
-- INSERT INTO `sys_role_menu` VALUES (2, 1132);
-- INSERT INTO `sys_role_menu` VALUES (2, 1133);

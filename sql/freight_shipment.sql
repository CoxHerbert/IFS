-- ----------------------------
-- Freight shipment planning
-- ----------------------------

DROP TABLE IF EXISTS `freight_shipment_order`;
DROP TABLE IF EXISTS `freight_container_plan`;
DROP TABLE IF EXISTS `freight_shipment_cargo`;
DROP TABLE IF EXISTS `freight_shipment_plan`;

CREATE TABLE `freight_shipment_plan` (
  `shipment_id` bigint NOT NULL COMMENT '出货计划ID',
  `shipment_no` varchar(64) NOT NULL COMMENT '出货计划编号',
  `order_no` varchar(64) DEFAULT '' COMMENT '客户订单号或参考号',
  `customer_id` bigint NOT NULL COMMENT '客户ID',
  `customer_name` varchar(128) DEFAULT '' COMMENT '客户名称',
  `pol` varchar(128) DEFAULT '' COMMENT '起运港',
  `pod` varchar(128) DEFAULT '' COMMENT '目的港',
  `planned_etd` varchar(32) DEFAULT '' COMMENT '计划开船日期',
  `planned_eta` varchar(32) DEFAULT '' COMMENT '计划到港日期',
  `actual_etd` varchar(32) DEFAULT '' COMMENT '实际开船日期',
  `actual_eta` varchar(32) DEFAULT '' COMMENT '实际到港日期',
  `status` varchar(8) DEFAULT '10' COMMENT '状态字典 freight_shipment_status',
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

INSERT INTO `sys_dict_type` VALUES (200, '出货状态', 'freight_shipment_status', '0', 'admin', now(), '', NULL, '客户端分享页展示状态');
INSERT INTO `sys_dict_data` VALUES (2000, 1, '计划已创建', '10', 'freight_shipment_status', '', 'primary', 'N', '0', 'admin', now(), '', NULL, '内部已建立出货计划');
INSERT INTO `sys_dict_data` VALUES (2001, 2, '出货计划已确认', '20', 'freight_shipment_status', '', 'primary', 'N', '0', 'admin', now(), '', NULL, '客户或运营确认计划');
INSERT INTO `sys_dict_data` VALUES (2002, 3, '等待客户发货', '30', 'freight_shipment_status', '', 'info', 'N', '0', 'admin', now(), '', NULL, '等待客户送货或安排提货');
INSERT INTO `sys_dict_data` VALUES (2003, 4, '已提货/已送仓', '40', 'freight_shipment_status', '', 'warning', 'N', '0', 'admin', now(), '', NULL, '货物已离开发货地');
INSERT INTO `sys_dict_data` VALUES (2004, 5, '仓库已收货', '50', 'freight_shipment_status', '', 'warning', 'N', '0', 'admin', now(), '', NULL, '仓库完成收货');
INSERT INTO `sys_dict_data` VALUES (2005, 6, '已入仓/码头进仓', '60', 'freight_shipment_status', '', 'warning', 'N', '0', 'admin', now(), '', NULL, '仓库或码头完成进仓');
INSERT INTO `sys_dict_data` VALUES (2006, 7, '订舱处理中', '70', 'freight_shipment_status', '', 'warning', 'N', '0', 'admin', now(), '', NULL, '正在向船司或渠道订舱');
INSERT INTO `sys_dict_data` VALUES (2007, 8, '舱位已确认', '80', 'freight_shipment_status', '', 'primary', 'N', '0', 'admin', now(), '', NULL, '订舱确认');
INSERT INTO `sys_dict_data` VALUES (2008, 9, '报关资料已收齐', '90', 'freight_shipment_status', '', 'primary', 'N', '0', 'admin', now(), '', NULL, '报关文件齐备');
INSERT INTO `sys_dict_data` VALUES (2009, 10, '报关已放行', '100', 'freight_shipment_status', '', 'success', 'N', '0', 'admin', now(), '', NULL, '出口报关放行');
INSERT INTO `sys_dict_data` VALUES (2013, 11, '已装柜', '110', 'freight_shipment_status', '', 'success', 'N', '0', 'admin', now(), '', NULL, '完成装柜');
INSERT INTO `sys_dict_data` VALUES (2014, 12, '已进港/码头放行', '120', 'freight_shipment_status', '', 'success', 'N', '0', 'admin', now(), '', NULL, '已进港或码头放行');
INSERT INTO `sys_dict_data` VALUES (2015, 13, '船舶已开船', '130', 'freight_shipment_status', '', 'success', 'N', '0', 'admin', now(), '', NULL, '船已开船');
INSERT INTO `sys_dict_data` VALUES (2016, 14, '目的港已到港', '140', 'freight_shipment_status', '', 'success', 'N', '0', 'admin', now(), '', NULL, '船已到目的港');
INSERT INTO `sys_dict_data` VALUES (2017, 15, '目的港清关中', '150', 'freight_shipment_status', '', 'warning', 'N', '0', 'admin', now(), '', NULL, '目的港清关处理中');
INSERT INTO `sys_dict_data` VALUES (2018, 16, '目的港已清关', '160', 'freight_shipment_status', '', 'success', 'N', '0', 'admin', now(), '', NULL, '目的港清关完成');
INSERT INTO `sys_dict_data` VALUES (2019, 17, '已派送/已签收', '170', 'freight_shipment_status', '', 'success', 'N', '0', 'admin', now(), '', NULL, '派送中或客户已签收');
INSERT INTO `sys_dict_data` VALUES (2020, 99, '异常处理中', '900', 'freight_shipment_status', '', 'danger', 'N', '0', 'admin', now(), '', NULL, '状态异常，运营正在处理');

INSERT INTO `sys_dict_type` VALUES (201, '货柜类型', 'freight_container_type', '0', 'admin', now(), '', NULL, '智能出货计划柜型');
INSERT INTO `sys_dict_data` VALUES (2010, 1, '20GP', '20GP', 'freight_container_type', '', 'default', 'N', '0', 'admin', now(), '', NULL, '约28CBM/21700KG');
INSERT INTO `sys_dict_data` VALUES (2011, 2, '40GP', '40GP', 'freight_container_type', '', 'default', 'N', '0', 'admin', now(), '', NULL, '约58CBM/26500KG');
INSERT INTO `sys_dict_data` VALUES (2012, 3, '40HQ', '40HQ', 'freight_container_type', '', 'default', 'Y', '0', 'admin', now(), '', NULL, '约68CBM/26500KG');
INSERT INTO `sys_dict_data` VALUES (2021, 4, 'LCL 拼箱', 'LCL', 'freight_container_type', '', 'info', 'N', '0', 'admin', now(), '', NULL, '小票货优先按散货拼箱评估');

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
  1, 0, 'M', '0', '0', '', 'international',
  'admin', NOW(), '', NULL, '国际货代业务菜单'
FROM dual
WHERE NOT EXISTS (
  SELECT 1 FROM `sys_menu` WHERE `menu_id` = 140
);

INSERT INTO `sys_menu` VALUES
(141, '出货计划', 140, 1, 'shipment', 'freight/shipment/index', 1, 0, 'C', '0', '0', 'freight:shipment:list', 'list', 'admin', now(), '', NULL, '出货计划管理菜单');

INSERT INTO `sys_menu` VALUES
(1160, '出货计划查询', 141, 1, '#', '', 1, 0, 'F', '0', '0', 'freight:shipment:query', '#', 'admin', now(), '', NULL, ''),
(1161, '出货清单导入', 141, 2, '#', '', 1, 0, 'F', '0', '0', 'freight:shipment:import', '#', 'admin', now(), '', NULL, ''),
(1162, '出货状态维护', 141, 3, '#', '', 1, 0, 'F', '0', '0', 'freight:shipment:edit', '#', 'admin', now(), '', NULL, ''),
(1163, '生成出货单', 141, 4, '#', '', 1, 0, 'F', '0', '0', 'freight:shipment:confirm', '#', 'admin', now(), '', NULL, ''),
(1164, '分享出货进度', 141, 5, '#', '', 1, 0, 'F', '0', '0', 'freight:shipment:share', '#', 'admin', now(), '', NULL, ''),
(1165, '出货计划删除', 141, 6, '#', '', 1, 0, 'F', '0', '0', 'freight:shipment:remove', '#', 'admin', now(), '', NULL, '');

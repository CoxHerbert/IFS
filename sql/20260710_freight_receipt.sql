-- 独立收款与一对多核销增量脚本
CREATE TABLE IF NOT EXISTS `freight_receipt` (
  `receipt_id` bigint NOT NULL, `receipt_no` varchar(64) NOT NULL, `customer_id` bigint NOT NULL, `customer_name` varchar(128) DEFAULT '',
  `amount` decimal(12,2) NOT NULL, `currency` varchar(8) DEFAULT 'CNY', `receipt_time` datetime NOT NULL,
  `payment_method` varchar(32) DEFAULT '', `status` varchar(16) DEFAULT 'UNALLOCATED',
  `voucher_url` varchar(500) DEFAULT '', `voucher_name` varchar(255) DEFAULT '', `remark` varchar(500) DEFAULT '',
  `create_by` varchar(64) DEFAULT '', `create_time` datetime DEFAULT NULL, `update_by` varchar(64) DEFAULT '', `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`receipt_id`), UNIQUE KEY `uk_freight_receipt_no` (`receipt_no`), KEY `idx_freight_receipt_customer` (`customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='独立收款单';
CREATE TABLE IF NOT EXISTS `freight_receipt_allocation` (
  `allocation_id` bigint NOT NULL, `receipt_id` bigint NOT NULL, `shipment_id` bigint NOT NULL, `allocated_amount` decimal(12,2) NOT NULL,
  PRIMARY KEY (`allocation_id`), UNIQUE KEY `uk_receipt_shipment` (`receipt_id`,`shipment_id`), KEY `idx_allocation_shipment` (`shipment_id`),
  CONSTRAINT `fk_allocation_receipt` FOREIGN KEY (`receipt_id`) REFERENCES `freight_receipt` (`receipt_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_allocation_shipment` FOREIGN KEY (`shipment_id`) REFERENCES `freight_shipment_plan` (`shipment_id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='收款核销明细';
INSERT IGNORE INTO `sys_menu` VALUES
(143, '收款管理', 140, 2, 'receipt', 'freight/receipt/index', 1, 0, 'C', '0', '0', 'freight:receipt:list', 'mdi:cash-multiple', 'admin', now(), '', NULL, '独立收款与出货计划核销'),
(1180, '收款查询', 143, 1, '#', '', 1, 0, 'F', '0', '0', 'freight:receipt:query', '#', 'admin', now(), '', NULL, ''),
(1181, '收款新增', 143, 2, '#', '', 1, 0, 'F', '0', '0', 'freight:receipt:add', '#', 'admin', now(), '', NULL, ''),
(1182, '收款删除', 143, 3, '#', '', 1, 0, 'F', '0', '0', 'freight:receipt:remove', '#', 'admin', now(), '', NULL, '');

CREATE TABLE IF NOT EXISTS `freight_payment_declaration` (
  `declaration_id` bigint NOT NULL, `declaration_no` varchar(64) NOT NULL, `customer_id` bigint NOT NULL, `customer_name` varchar(128) DEFAULT '',
  `shipment_id` bigint NOT NULL, `shipment_no` varchar(64) DEFAULT '', `amount` decimal(12,2) NOT NULL, `currency` varchar(8) DEFAULT 'CNY',
  `payment_time` varchar(32) DEFAULT '', `voucher_url` varchar(500) NOT NULL, `voucher_name` varchar(255) DEFAULT '',
  `status` varchar(16) DEFAULT 'PENDING', `remark` varchar(500) DEFAULT '', `review_by` varchar(64) DEFAULT '', `review_time` datetime DEFAULT NULL, `review_remark` varchar(500) DEFAULT '', `create_by` varchar(64) DEFAULT '', `create_time` datetime DEFAULT NULL,
  `update_by` varchar(64) DEFAULT '', `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`declaration_id`), UNIQUE KEY `uk_payment_declaration_no` (`declaration_no`), KEY `idx_declaration_customer` (`customer_id`), KEY `idx_declaration_shipment` (`shipment_id`),
  CONSTRAINT `fk_declaration_shipment` FOREIGN KEY (`shipment_id`) REFERENCES `freight_shipment_plan` (`shipment_id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='客户付款申报';
ALTER TABLE `freight_payment_declaration` ADD COLUMN IF NOT EXISTS `review_by` varchar(64) DEFAULT '' AFTER `remark`;
ALTER TABLE `freight_payment_declaration` ADD COLUMN IF NOT EXISTS `review_time` datetime DEFAULT NULL AFTER `review_by`;
ALTER TABLE `freight_payment_declaration` ADD COLUMN IF NOT EXISTS `review_remark` varchar(500) DEFAULT '' AFTER `review_time`;
INSERT IGNORE INTO `sys_menu` VALUES
(144, '付款申报', 140, 3, 'payment-declaration', 'freight/paymentDeclaration/index', 1, 0, 'C', '0', '0', 'freight:declaration:list', 'mdi:receipt-text-check-outline', 'admin', now(), '', NULL, '客户付款申报审核'),
(1190, '申报查询', 144, 1, '#', '', 1, 0, 'F', '0', '0', 'freight:declaration:query', '#', 'admin', now(), '', NULL, ''),
(1191, '申报审核', 144, 2, '#', '', 1, 0, 'F', '0', '0', 'freight:declaration:review', '#', 'admin', now(), '', NULL, '');
UPDATE `sys_menu` SET `icon`='mdi:earth' WHERE `menu_id`=140;
UPDATE `sys_menu` SET `icon`='mdi:truck-cargo-container' WHERE `menu_id`=141;
UPDATE `sys_menu` SET `icon`='mdi:cash-multiple' WHERE `menu_id`=143;
UPDATE `sys_menu` SET `icon`='mdi:receipt-text-check-outline' WHERE `menu_id`=144;
UPDATE `customer_workspace_menu` SET `icon`='mdi:view-dashboard-outline' WHERE `menu_id`=20001;
UPDATE `customer_workspace_menu` SET `icon`='mdi:account-outline' WHERE `menu_id`=20002;
UPDATE `customer_workspace_menu` SET `icon`='mdi:radar' WHERE `menu_id`=20003;
UPDATE `customer_workspace_menu` SET `icon`='mdi:calculator-variant-outline' WHERE `menu_id`=20004;
UPDATE `customer_workspace_menu` SET `icon`='mdi:message-text-outline' WHERE `menu_id`=20005;

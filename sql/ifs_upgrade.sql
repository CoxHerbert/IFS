-- IFS existing-environment upgrade entrypoint.
-- Safe to execute repeatedly on MySQL 5.7+/8.0.
-- New environments should execute ifs_init.sql instead.

-- =============================================================================
-- Freight shipment: cargo dimensions, currency, trade and transport fields
-- =============================================================================
SET @ddl = IF((SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='freight_shipment_cargo' AND COLUMN_NAME='unit_weight_kg')=0, 'ALTER TABLE `freight_shipment_cargo` ADD COLUMN `unit_weight_kg` decimal(12,4) NOT NULL DEFAULT 0.0000 COMMENT ''单个重量KG'' AFTER `height_cm`', 'SELECT 1'); PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl = IF((SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='freight_shipment_cargo' AND COLUMN_NAME='unit_volume_cbm')=0, 'ALTER TABLE `freight_shipment_cargo` ADD COLUMN `unit_volume_cbm` decimal(12,6) NOT NULL DEFAULT 0.000000 COMMENT ''单个体积CBM'' AFTER `unit_weight_kg`', 'SELECT 1'); PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl = IF((SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='freight_shipment_plan' AND COLUMN_NAME='currency')=0, 'ALTER TABLE `freight_shipment_plan` ADD COLUMN `currency` varchar(8) NOT NULL DEFAULT ''CNY'' COMMENT ''结算币种'' AFTER `payment_amount`', 'SELECT 1'); PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl = IF((SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='freight_shipment_plan' AND COLUMN_NAME='trade_term')=0, 'ALTER TABLE `freight_shipment_plan` ADD COLUMN `trade_term` varchar(16) NOT NULL DEFAULT '''' COMMENT ''贸易条款'' AFTER `currency`', 'SELECT 1'); PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl = IF((SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='freight_shipment_plan' AND COLUMN_NAME='delivery_type')=0, 'ALTER TABLE `freight_shipment_plan` ADD COLUMN `delivery_type` varchar(32) NOT NULL DEFAULT '''' COMMENT ''运输范围'' AFTER `trade_term`', 'SELECT 1'); PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl = IF((SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='freight_shipment_plan' AND COLUMN_NAME='pickup_address')=0, 'ALTER TABLE `freight_shipment_plan` ADD COLUMN `pickup_address` varchar(500) NOT NULL DEFAULT '''' COMMENT ''提货地址'' AFTER `delivery_type`', 'SELECT 1'); PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl = IF((SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='freight_shipment_plan' AND COLUMN_NAME='delivery_address')=0, 'ALTER TABLE `freight_shipment_plan` ADD COLUMN `delivery_address` varchar(500) NOT NULL DEFAULT '''' COMMENT ''送货地址'' AFTER `pickup_address`', 'SELECT 1'); PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl = IF((SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='freight_shipment_plan' AND COLUMN_NAME='handover_location')=0, 'ALTER TABLE `freight_shipment_plan` ADD COLUMN `handover_location` varchar(255) NOT NULL DEFAULT '''' COMMENT ''约定交货地点'' AFTER `delivery_address`', 'SELECT 1'); PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl = IF((SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='freight_shipment_plan' AND COLUMN_NAME='clearance_party')=0, 'ALTER TABLE `freight_shipment_plan` ADD COLUMN `clearance_party` varchar(32) NOT NULL DEFAULT '''' COMMENT ''目的国清关方'' AFTER `handover_location`', 'SELECT 1'); PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl = IF((SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='freight_shipment_plan' AND COLUMN_NAME='duty_payer')=0, 'ALTER TABLE `freight_shipment_plan` ADD COLUMN `duty_payer` varchar(32) NOT NULL DEFAULT '''' COMMENT ''关税税费承担方'' AFTER `clearance_party`', 'SELECT 1'); PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl = IF((SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='freight_shipment_plan' AND COLUMN_NAME='vessel_name')=0, 'ALTER TABLE `freight_shipment_plan` ADD COLUMN `vessel_name` varchar(128) NOT NULL DEFAULT '''' COMMENT ''船名'' AFTER `duty_payer`', 'SELECT 1'); PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @ddl = IF((SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='freight_shipment_plan' AND COLUMN_NAME='voyage_no')=0, 'ALTER TABLE `freight_shipment_plan` ADD COLUMN `voyage_no` varchar(64) NOT NULL DEFAULT '''' COMMENT ''航次'' AFTER `vessel_name`', 'SELECT 1'); PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;

UPDATE `sys_dict_data` SET `remark`='约28CBM/21700KG' WHERE `dict_type`='freight_container_type' AND `dict_value`='20GP';
UPDATE `sys_dict_data` SET `remark`='约58CBM/26500KG' WHERE `dict_type`='freight_container_type' AND `dict_value`='40GP';
UPDATE `sys_dict_data` SET `remark`='约68CBM/26500KG' WHERE `dict_type`='freight_container_type' AND `dict_value`='40HQ';
UPDATE `sys_dict_data` SET `remark`='约15CBM/3000KG' WHERE `dict_type`='freight_container_type' AND `dict_value`='LCL';

-- =============================================================================
-- Finance: payable details, independent receipt/payment menu and portal menu
-- =============================================================================
SET @ddl = IF((SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='freight_shipment_plan' AND COLUMN_NAME='payable_amount')=0, 'ALTER TABLE `freight_shipment_plan` ADD COLUMN `payable_amount` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT ''客户应付金额'' AFTER `payment_amount`', 'SELECT 1'); PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;

CREATE TABLE IF NOT EXISTS `freight_shipment_charge` (
  `charge_id` bigint NOT NULL, `shipment_id` bigint NOT NULL, `fee_name` varchar(128) NOT NULL,
  `amount` decimal(12,2) NOT NULL, `currency` varchar(8) NOT NULL DEFAULT 'CNY', `remark` varchar(500) DEFAULT '',
  `create_by` varchar(64) DEFAULT '', `create_time` datetime DEFAULT NULL,
  PRIMARY KEY (`charge_id`), KEY `idx_charge_shipment` (`shipment_id`),
  CONSTRAINT `fk_charge_shipment` FOREIGN KEY (`shipment_id`) REFERENCES `freight_shipment_plan` (`shipment_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='出货计划应付费用明细';

INSERT IGNORE INTO `sys_menu` VALUES
(148,'收付款管理',0,5,'finance',NULL,1,0,'M','0','0','','mdi:cash-register','admin',now(),'',NULL,'收款、核销及客户账款独立模块'),
(145,'收付款台账',148,2,'payment','freight/payment/index',1,0,'C','0','0','freight:payment:list','mdi:cash-multiple','admin',now(),'',NULL,'客户应付、实收及未付台账'),
(1270,'收付款查询',145,1,'#','',1,0,'F','0','0','freight:payment:list','#','admin',now(),'',NULL,''),
(1271,'应付金额维护',145,2,'#','',1,0,'F','0','0','freight:payment:edit','#','admin',now(),'',NULL,'');
UPDATE `sys_menu` SET `menu_name`='收付款管理',`parent_id`=0,`order_num`=5,`path`='finance',`component`=NULL,`menu_type`='M',`visible`='0',`status`='0',`perms`='',`icon`='mdi:cash-register' WHERE `menu_id`=148;
UPDATE `sys_menu` SET `menu_name`='收款管理',`parent_id`=148,`order_num`=1,`path`='receipt',`component`='freight/receipt/index' WHERE `menu_id`=143;
UPDATE `sys_menu` SET `menu_name`='收付款台账',`parent_id`=148,`order_num`=2,`path`='payment',`component`='freight/payment/index' WHERE `menu_id`=145;
UPDATE `sys_menu` SET `menu_name`='付款申报',`parent_id`=148,`order_num`=3,`path`='payment-declaration',`component`='freight/paymentDeclaration/index' WHERE `menu_id`=144;
INSERT IGNORE INTO `sys_role_menu` (`role_id`,`menu_id`) SELECT DISTINCT `role_id`,148 FROM `sys_role_menu` WHERE `menu_id` IN (143,144,145,1180,1181,1182,1190,1191,1270,1271);
INSERT IGNORE INTO `customer_workspace_menu` (`menu_id`,`parent_id`,`menu_name`,`order_num`,`path`,`component`,`is_cache`,`menu_type`,`visible`,`status`,`perms`,`icon`,`remark`,`create_by`,`create_time`,`update_by`,`update_time`) VALUES (20006,0,'付款明细','6','payments','PaymentDetails/index','0','C','0','0','portal:payment:view','mdi:cash-multiple','当前客户付款明细','admin',now(),'admin',now());
INSERT IGNORE INTO `customer_workspace_role_menu` (`role_id`,`menu_id`) VALUES (20001,20006);

-- =============================================================================
-- Agent: independent root menu and role-menu repair
-- =============================================================================
DELETE rm FROM `sys_role_menu` rm JOIN `sys_menu` m ON m.`menu_id`=rm.`menu_id` WHERE m.`menu_id` IN (1200,1201) AND m.`perms` LIKE 'freight:payment:%';
DELETE FROM `sys_menu` WHERE `menu_id` IN (1200,1201) AND `perms` LIKE 'freight:payment:%';
INSERT IGNORE INTO `sys_menu` VALUES
(147,'Agent 管理',0,6,'agent',NULL,1,0,'M','0','0','','mdi:robot-outline','admin',now(),'',NULL,'Agent 独立功能模块'),
(142,'Agent 对话',147,1,'chat','agent/chat/index',1,0,'C','0','0','ifs:agent:chat','message','admin',now(),'',NULL,'IFS Agent 对话管理菜单'),
(146,'Agent 配置',147,2,'config','agent/config/index',1,0,'C','0','0','ifs:agent:config','setting','admin',now(),'',NULL,'Agent 本地模型与运行参数配置'),
(1170,'Agent 对话权限',142,1,'#','',1,0,'F','0','0','ifs:agent:chat','#','admin',now(),'',NULL,''),
(1200,'Agent 配置权限',146,1,'#','',1,0,'F','0','0','ifs:agent:config','#','admin',now(),'',NULL,'');
UPDATE `sys_menu` SET `menu_name`='Agent 管理',`parent_id`=0,`order_num`=6,`path`='agent',`component`=NULL,`menu_type`='M',`visible`='0',`status`='0',`perms`='',`icon`='mdi:robot-outline' WHERE `menu_id`=147;
UPDATE `sys_menu` SET `menu_name`='Agent 对话',`parent_id`=147,`order_num`=1,`path`='chat',`component`='agent/chat/index',`menu_type`='C',`visible`='0',`status`='0',`perms`='ifs:agent:chat' WHERE `menu_id`=142;
UPDATE `sys_menu` SET `menu_name`='Agent 配置',`parent_id`=147,`order_num`=2,`path`='config',`component`='agent/config/index',`menu_type`='C',`visible`='0',`status`='0',`perms`='ifs:agent:config' WHERE `menu_id`=146;
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) SELECT DISTINCT rm.`role_id`,147 FROM `sys_role_menu` rm WHERE rm.`menu_id` IN (142,146) AND NOT EXISTS (SELECT 1 FROM `sys_role_menu` x WHERE x.`role_id`=rm.`role_id` AND x.`menu_id`=147);
INSERT INTO `sys_role_menu` (`role_id`,`menu_id`) SELECT DISTINCT ur.`role_id`,m.`menu_id` FROM `sys_user` u JOIN `sys_user_role` ur ON ur.`user_id`=u.`user_id` JOIN `sys_role` r ON r.`role_id`=ur.`role_id` AND r.`status`='0' JOIN `sys_menu` m ON m.`menu_id` IN (147,142,146,1170,1200) WHERE u.`user_name`='admin' AND NOT EXISTS (SELECT 1 FROM `sys_role_menu` x WHERE x.`role_id`=ur.`role_id` AND x.`menu_id`=m.`menu_id`);

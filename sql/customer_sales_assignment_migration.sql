-- ----------------------------
-- Customer salesperson assignment and shipment ownership
-- ----------------------------

ALTER TABLE `customer`
  ADD COLUMN `sales_user_id` bigint DEFAULT 0 COMMENT '负责业务员用户ID' AFTER `email`,
  ADD COLUMN `sales_user_name` varchar(64) DEFAULT '' COMMENT '负责业务员名称' AFTER `sales_user_id`;

ALTER TABLE `customer`
  ADD KEY `idx_customer_sales_user` (`sales_user_id`) USING BTREE;

ALTER TABLE `freight_shipment_plan`
  ADD COLUMN `sales_user_id` bigint DEFAULT 0 COMMENT '负责业务员用户ID快照' AFTER `customer_name`,
  ADD COLUMN `sales_user_name` varchar(64) DEFAULT '' COMMENT '负责业务员名称快照' AFTER `sales_user_id`;

ALTER TABLE `freight_shipment_plan`
  ADD KEY `idx_freight_sales_user` (`sales_user_id`) USING BTREE;

UPDATE `freight_shipment_plan` p
LEFT JOIN `customer` c ON p.`customer_id` = c.`customer_id`
SET p.`sales_user_id` = IFNULL(c.`sales_user_id`, 0),
    p.`sales_user_name` = IFNULL(c.`sales_user_name`, '')
WHERE IFNULL(p.`sales_user_id`, 0) = 0;

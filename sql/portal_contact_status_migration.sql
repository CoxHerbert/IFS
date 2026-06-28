-- ----------------------------
-- Migrate portal_contact status values
-- 10 待跟进, 20 跟进中, 30 已完成, 40 无效
-- ----------------------------

ALTER TABLE `portal_contact`
  MODIFY COLUMN `status` char(2) CHARACTER SET utf8 COLLATE utf8_unicode_ci NULL DEFAULT '10'
  COMMENT '状态（10待跟进 20跟进中 30已完成 40无效）';

UPDATE `portal_contact`
SET `status` = CASE `status`
  WHEN '0' THEN '10'
  WHEN '1' THEN '20'
  WHEN '2' THEN '30'
  WHEN '3' THEN '40'
  ELSE `status`
END
WHERE `status` IN ('0', '1', '2', '3');

-- ----------------------------
-- Dict data for portal_contact status
-- 字典类型：portal_contact_status
-- 10 待跟进, 20 跟进中, 30 已完成, 40 无效
-- ----------------------------

DELETE FROM `sys_dict_data` WHERE `dict_type` = 'portal_contact_status';
DELETE FROM `sys_dict_type` WHERE `dict_type` = 'portal_contact_status';

SET @dict_id := (SELECT IFNULL(MAX(`dict_id`), 0) + 1 FROM `sys_dict_type`);
INSERT INTO `sys_dict_type`
(`dict_id`, `dict_name`, `dict_type`, `status`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
VALUES
(@dict_id, '官网线索跟进状态', 'portal_contact_status', '0', 'admin', now(), '', NULL, '官网联系我们线索跟进状态');

SET @dict_code := (SELECT IFNULL(MAX(`dict_code`), 0) FROM `sys_dict_data`);
INSERT INTO `sys_dict_data`
(`dict_code`, `dict_sort`, `dict_label`, `dict_value`, `dict_type`, `css_class`, `list_class`, `is_default`, `status`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
VALUES
(@dict_code := @dict_code + 1, 1, '待跟进', '10', 'portal_contact_status', '', 'info', 'Y', '0', 'admin', now(), '', NULL, '待跟进'),
(@dict_code := @dict_code + 1, 2, '跟进中', '20', 'portal_contact_status', '', 'warning', 'N', '0', 'admin', now(), '', NULL, '跟进中'),
(@dict_code := @dict_code + 1, 3, '已完成', '30', 'portal_contact_status', '', 'success', 'N', '0', 'admin', now(), '', NULL, '已完成'),
(@dict_code := @dict_code + 1, 4, '无效', '40', 'portal_contact_status', '', 'danger', 'N', '0', 'admin', now(), '', NULL, '无效');

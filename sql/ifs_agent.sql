-- IFS Agent initialization.
-- Includes Agent chat persistence, form submissions, and backend Agent menu.
-- Safe for repeated execution on an initialized Baize database.

CREATE TABLE IF NOT EXISTS `chat_session` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint DEFAULT 0 COMMENT '用户ID；免登录场景为0',
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
  1, 0, 'M', '0', '0', '', 'international',
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

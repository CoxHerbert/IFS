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

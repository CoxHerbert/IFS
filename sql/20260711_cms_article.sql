CREATE TABLE IF NOT EXISTS `cms_article` (
  `article_id` bigint NOT NULL COMMENT '文章ID',
  `title` varchar(200) NOT NULL COMMENT '标题',
  `slug` varchar(220) NOT NULL COMMENT '访问标识',
  `summary` varchar(500) DEFAULT '' COMMENT '摘要',
  `category` varchar(64) NOT NULL DEFAULT '' COMMENT '栏目',
  `cover_url` varchar(500) DEFAULT '' COMMENT '封面图',
  `content` text COMMENT '正文',
  `status` char(1) NOT NULL DEFAULT '0' COMMENT '状态（0发布 1草稿）',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序',
  `publish_time` datetime DEFAULT NULL COMMENT '发布时间',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`article_id`),
  UNIQUE KEY `uk_cms_article_slug` (`slug`),
  KEY `idx_cms_article_status` (`status`, `publish_time`),
  KEY `idx_cms_article_category` (`category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='CMS文章表';

INSERT INTO `sys_menu` (
  `menu_id`, `menu_name`, `parent_id`, `order_num`, `path`, `component`, `is_frame`, `is_cache`, `menu_type`,
  `visible`, `status`, `perms`, `icon`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`
)
SELECT 144, 'CMS管理', 0, 9, 'cms', '', 1, 0, 'M',
       '0', '0', '', 'message', 'admin', NOW(), '', NULL, 'CMS管理目录'
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `menu_id` = 144);

UPDATE `sys_menu`
SET `menu_name` = 'CMS管理',
    `parent_id` = 0,
    `path` = 'cms',
    `component` = '',
    `menu_type` = 'M',
    `perms` = '',
    `icon` = 'message',
    `remark` = 'CMS管理目录'
WHERE `menu_id` = 144;

INSERT INTO `sys_menu` (
  `menu_id`, `menu_name`, `parent_id`, `order_num`, `path`, `component`, `is_frame`, `is_cache`, `menu_type`,
  `visible`, `status`, `perms`, `icon`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`
)
SELECT 145, '新闻资讯', 144, 1, 'article', 'cms/article/index', 1, 0, 'C',
       '0', '0', 'cms:article:list', 'message', 'admin', NOW(), '', NULL, 'CMS新闻资讯管理菜单'
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `menu_id` = 145);

UPDATE `sys_menu`
SET `menu_name` = '新闻资讯',
    `parent_id` = 144,
    `path` = 'article',
    `component` = 'cms/article/index',
    `menu_type` = 'C',
    `perms` = 'cms:article:list',
    `icon` = 'message',
    `remark` = 'CMS新闻资讯管理菜单'
WHERE `menu_id` = 145;

INSERT INTO `sys_menu` (
  `menu_id`, `menu_name`, `parent_id`, `order_num`, `path`, `component`, `is_frame`, `is_cache`, `menu_type`,
  `visible`, `status`, `perms`, `icon`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`
)
SELECT 1174, '新闻资讯查询', 145, 1, '#', '', 1, 0, 'F', '0', '0', 'cms:article:list', '#', 'admin', NOW(), '', NULL, ''
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `menu_id` = 1174);

UPDATE `sys_menu` SET `menu_name` = '新闻资讯查询', `parent_id` = 145, `perms` = 'cms:article:list' WHERE `menu_id` = 1174;

INSERT INTO `sys_menu` (
  `menu_id`, `menu_name`, `parent_id`, `order_num`, `path`, `component`, `is_frame`, `is_cache`, `menu_type`,
  `visible`, `status`, `perms`, `icon`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`
)
SELECT 1175, '新闻资讯详情', 145, 2, '#', '', 1, 0, 'F', '0', '0', 'cms:article:query', '#', 'admin', NOW(), '', NULL, ''
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `menu_id` = 1175);

UPDATE `sys_menu` SET `menu_name` = '新闻资讯详情', `parent_id` = 145, `perms` = 'cms:article:query' WHERE `menu_id` = 1175;

INSERT INTO `sys_menu` (
  `menu_id`, `menu_name`, `parent_id`, `order_num`, `path`, `component`, `is_frame`, `is_cache`, `menu_type`,
  `visible`, `status`, `perms`, `icon`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`
)
SELECT 1176, '新闻资讯新增', 145, 3, '#', '', 1, 0, 'F', '0', '0', 'cms:article:add', '#', 'admin', NOW(), '', NULL, ''
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `menu_id` = 1176);

UPDATE `sys_menu` SET `menu_name` = '新闻资讯新增', `parent_id` = 145, `perms` = 'cms:article:add' WHERE `menu_id` = 1176;

INSERT INTO `sys_menu` (
  `menu_id`, `menu_name`, `parent_id`, `order_num`, `path`, `component`, `is_frame`, `is_cache`, `menu_type`,
  `visible`, `status`, `perms`, `icon`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`
)
SELECT 1177, '新闻资讯修改', 145, 4, '#', '', 1, 0, 'F', '0', '0', 'cms:article:edit', '#', 'admin', NOW(), '', NULL, ''
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `menu_id` = 1177);

UPDATE `sys_menu` SET `menu_name` = '新闻资讯修改', `parent_id` = 145, `perms` = 'cms:article:edit' WHERE `menu_id` = 1177;

INSERT INTO `sys_menu` (
  `menu_id`, `menu_name`, `parent_id`, `order_num`, `path`, `component`, `is_frame`, `is_cache`, `menu_type`,
  `visible`, `status`, `perms`, `icon`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`
)
SELECT 1178, '新闻资讯删除', 145, 5, '#', '', 1, 0, 'F', '0', '0', 'cms:article:remove', '#', 'admin', NOW(), '', NULL, ''
WHERE NOT EXISTS (SELECT 1 FROM `sys_menu` WHERE `menu_id` = 1178);

UPDATE `sys_menu` SET `menu_name` = '新闻资讯删除', `parent_id` = 145, `perms` = 'cms:article:remove' WHERE `menu_id` = 1178;

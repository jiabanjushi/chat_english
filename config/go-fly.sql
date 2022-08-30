DROP TABLE IF EXISTS `user`|
CREATE TABLE `user` (
 `id` int(11) NOT NULL AUTO_INCREMENT,
 `pid` int(11) unsigned NOT NULL DEFAULT 0,
 `name` varchar(125) NOT NULL DEFAULT '',
 `password` varchar(50) NOT NULL DEFAULT '',
 `nickname` varchar(50) NOT NULL DEFAULT '',
 `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `expired_at` timestamp NULL DEFAULT NULL,
 `updated_at` timestamp NULL DEFAULT NULL,
 `deleted_at` timestamp NULL DEFAULT NULL,
 `avator` varchar(100) NOT NULL DEFAULT '',
 `rec_num` int(10) unsigned NOT NULL DEFAULT 0,
 `online_status` tinyint(4) NOT NULL DEFAULT 1,
 `status` tinyint(4) NOT NULL DEFAULT '0',
 `agent_num` int(11) unsigned NOT NULL DEFAULT 10,
 PRIMARY KEY (`id`),
 UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8|
INSERT INTO `user` (`id`, `name`, `password`, `nickname`, `created_at`, `updated_at`, `expired_at`, `avator`, `pid`) VALUES
(1, 'caonima888', '585d234cb6a583ad1bd1fd65c75b7219', '菜地超管', '2020-07-02 14:36:46', '2020-07-05 08:46:57', '2030-07-04 09:32:20', '/static/images/user/2.jpg', 0)|

DROP TABLE IF EXISTS `visitor`|
CREATE TABLE `visitor` (
 `id` int(11) NOT NULL AUTO_INCREMENT,
 `name` varchar(550) NOT NULL DEFAULT '',
 `real_name` varchar(550) NOT NULL DEFAULT '',
 `avator` varchar(500) NOT NULL DEFAULT '',
 `source_ip` varchar(50) NOT NULL DEFAULT '',
 `to_id` varchar(50) NOT NULL DEFAULT '',
 `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `updated_at` timestamp NULL DEFAULT NULL,
 `deleted_at` timestamp NULL DEFAULT NULL,
 `visitor_id` varchar(100) NOT NULL DEFAULT '',
 `status` tinyint(4) NOT NULL DEFAULT '0',
 `refer` varchar(500) NOT NULL DEFAULT '',
 `city` varchar(100) NOT NULL DEFAULT '',
 `client_ip` varchar(100) NOT NULL,
 `extra` varchar(2048) NOT NULL DEFAULT '',
 `ent_id` int(11) unsigned NOT NULL DEFAULT 0,
 PRIMARY KEY (`id`),
 UNIQUE KEY `visitor_id` (`visitor_id`),
 KEY `to_id` (`to_id`),
 KEY `idx_update` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4|
DROP TABLE IF EXISTS `visitor_ext`|
CREATE TABLE `visitor_ext` (
 `id` int(11) NOT NULL AUTO_INCREMENT,
 `visitor_id` varchar(100) NOT NULL DEFAULT '',
 `ent_id` int(11) unsigned NOT NULL DEFAULT 0,
 `ua` varchar(500) NOT NULL DEFAULT '',
 `title` varchar(256) NOT NULL DEFAULT '',
 `url` varchar(500) NOT NULL DEFAULT '',
 `server_ip` varchar(128) NOT NULL DEFAULT '',
 `client_ip` varchar(128) NOT NULL DEFAULT '',
 `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
 PRIMARY KEY (`id`),
 KEY `visitor_id` (`visitor_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8|
DROP TABLE IF EXISTS `visitor_attr`|
CREATE TABLE `visitor_attr` (
 `id` int(11) NOT NULL AUTO_INCREMENT,
 `visitor_id` varchar(100) NOT NULL DEFAULT '',
 `real_name` varchar(100) NOT NULL DEFAULT '',
 `tel` varchar(50) NOT NULL DEFAULT '',
 `email` varchar(50) NOT NULL DEFAULT '',
 `qq` varchar(100) NOT NULL DEFAULT '',
 `wechat` varchar(100) NOT NULL DEFAULT '',
 `comment` varchar(500) NOT NULL DEFAULT '',
 `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `ent_id` int(11) unsigned NOT NULL DEFAULT 0,
 PRIMARY KEY (`id`),
 KEY `visitor_id` (`visitor_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8|
DROP TABLE IF EXISTS `message`|
CREATE TABLE `message` (
 `id` int(11) NOT NULL AUTO_INCREMENT,
 `kefu_id` varchar(100) NOT NULL DEFAULT '',
 `visitor_id` varchar(100) NOT NULL DEFAULT '',
 `content` varchar(2048) NOT NULL DEFAULT '',
 `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `updated_at` timestamp NULL DEFAULT NULL,
 `deleted_at` timestamp NULL DEFAULT NULL,
 `mes_type` enum('kefu','visitor') NOT NULL DEFAULT 'visitor',
 `status` enum('read','unread') NOT NULL DEFAULT 'unread',
 `ent_id` int(11) unsigned NOT NULL DEFAULT 0,
 PRIMARY KEY (`id`),
 KEY `kefu_id` (`kefu_id`),
 KEY `visitor_id` (`visitor_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4|

DROP TABLE IF EXISTS `user_role`|
CREATE TABLE `user_role` (
 `id` int(11) NOT  NULL AUTO_INCREMENT,
 `user_id` int(11) NOT NULL DEFAULT '0',
 `role_id` int(11) NOT NULL DEFAULT '0',
 PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8|
INSERT INTO `user_role` (`id`, `user_id`, `role_id`) VALUES
(1, 1, 1),
(2, 2, 2),
(3, 3, 3),
(4, 4, 3)|

DROP TABLE IF EXISTS `role`|
CREATE TABLE `role` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '',
  `method` varchar(100) NOT NULL DEFAULT '',
  `path` varchar(2048) NOT NULL DEFAULT '',
  `is_super` tinyint(3) unsigned NOT NULL DEFAULT '0',
   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8|
INSERT INTO `role` (`id`, `name`, `method`, `path`, `is_super`) VALUES
(3, '普通坐席', 'GET', 'GET:/kefuinfo', '0'),
(2, '普通商户', '*', '*', '0'),
(1, '超级管理员', '*', '*', '1')|

DROP TABLE IF EXISTS `welcome`|
CREATE TABLE `welcome` (
 `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
 `user_id` varchar(100) NOT NULL DEFAULT '',
 `keyword` varchar(100) NOT NULL DEFAULT '',
 `content` varchar(500)  COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
 `is_default` tinyint(3) unsigned NOT NULL DEFAULT '0',
 `delay_second` int(10) unsigned NOT NULL DEFAULT '3',
 `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
 PRIMARY KEY (`id`),
 KEY `user_id` (`user_id`),
 KEY `keyword` (`keyword`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8|
INSERT INTO `welcome` (`id`, `user_id`, `content`, `is_default`, `ctime`, `keyword`) VALUES
(NULL, 'kefu2', '你好~', 1, '2020-08-24 02:57:49','welcome')|
INSERT INTO `welcome` (`id`, `user_id`, `content`, `is_default`, `ctime`, `keyword`) VALUES
(NULL, 'kefu2', '您好，这里是GOFLY客服多商户版，有需求可立即联系客服', 0, '2020-08-24 02:57:49','welcome')|

DROP TABLE IF EXISTS `ipblack`|
CREATE TABLE `ipblack` (
 `id` int(11) NOT NULL AUTO_INCREMENT,
 `ip` varchar(100) NOT NULL DEFAULT '',
 `name` varchar(500) NOT NULL DEFAULT '',
 `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `kefu_id` varchar(100) NOT NULL DEFAULT '',
 PRIMARY KEY (`id`),
 UNIQUE KEY `ip` (`ip`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8|
DROP TABLE IF EXISTS `config`|
CREATE TABLE `config` (
 `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
 `conf_name` varchar(255) NOT NULL DEFAULT '',
 `conf_key` varchar(255) NOT NULL DEFAULT '',
 `conf_value` varchar(2550) NOT NULL DEFAULT '',
 PRIMARY KEY (`id`),
 UNIQUE KEY `conf_key` (`conf_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, '是否允许上传附件', 'SendAttachment', 'true')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, '发送通知邮件(SMTP地址)', 'NoticeEmailSmtp', '')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, '发送通知邮件(邮箱)', 'NoticeEmailAddress', '')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, '发送通知邮件(密码)', 'NoticeEmailPassword', '')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, 'App个推(Token)', 'GetuiToken', '')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, 'App个推(AppID)', 'GetuiAppID', '')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, 'App个推(AppKey)', 'GetuiAppKey', '')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, 'App个推(AppSecret)', 'GetuiAppSecret', '')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, 'App个推(AppMasterSecret)', 'GetuiMasterSecret', '')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, '落地域名跳转( 例: 留空不跳转 )', 'LandHost', '')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, '客服系统公告', 'SystemNotice', '')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, '客服系统标题', 'SystemTitle', 'GOFLY在线客服系统')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, '客服系统登录页标题', 'SystemLoginTitle', 'GOFLY在线客服登录')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, '客服系统关键字', 'SystemKeywords', 'GOFLY在线客服系统')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, '客服系统描述', 'SystemDesc', 'GOFLY在线客服系统')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, '客服首页跳转地址( 例: http://域名/login )', 'IndexJumpUrl', '')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, '版权信息文案', 'CopyrightTxt', 'GOFLY在线客服版权所有 © 2020-2022')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, '版权链接地址', 'CopyrightUrl', '')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, '官方客服JS配置', 'SystemKefu', '')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, '是否允许注册(1是允许，2是不允许，默认允许)', 'SystemRegister', '')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, '系统版本号', 'SystemVersion', '0.6.1')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, '系统版本名称', 'SystemVersionName', '商务运营版')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, '是否显示客服名称（off是不显示）', 'ShowKefuName', '')|
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (NULL, '微信模板消息备注字段(remark)', 'WechatTemplateRemark', '')|
DROP TABLE IF EXISTS `ent_config`|
CREATE TABLE `ent_config` (
 `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
 `conf_name` varchar(255) NOT NULL DEFAULT '',
 `conf_key` varchar(255) NOT NULL DEFAULT '',
 `conf_value` varchar(255) NOT NULL DEFAULT '',
 `ent_id` varchar(255) NOT NULL DEFAULT '',
 PRIMARY KEY (`id`),
 KEY `ent_id` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8|
INSERT INTO `ent_config` (`id`, `conf_name`, `conf_key`, `ent_id`) VALUES (NULL, '发送通知邮件(SMTP地址)', 'NoticeEmailSmtp', '2')|
INSERT INTO `ent_config` (`id`, `conf_name`, `conf_key`, `ent_id`) VALUES (NULL, '发送通知邮件(邮箱)', 'NoticeEmailAddress', '2')|
INSERT INTO `ent_config` (`id`, `conf_name`, `conf_key`, `ent_id`) VALUES (NULL, '发送通知邮件(密码)', 'NoticeEmailPassword', '2')|
DROP TABLE IF EXISTS `about`|
CREATE TABLE `about` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `title_cn` varchar(255) NOT NULL DEFAULT '',
  `title_en` varchar(255) NOT NULL DEFAULT '',
  `keywords_cn` varchar(255) NOT NULL DEFAULT '',
  `keywords_en` varchar(255) NOT NULL DEFAULT '',
  `desc_cn` varchar(1024) NOT NULL DEFAULT '',
  `desc_en` varchar(1024) NOT NULL DEFAULT '',
  `css_js` text NOT NULL,
  `html_cn` text NOT NULL,
  `html_en` text NOT NULL,
  `page` varchar(50) NOT NULL DEFAULT '',
PRIMARY KEY (`id`),
UNIQUE KEY `page` (`page`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8|
INSERT INTO `about` (`id`, `title_cn`, `title_en`, `keywords_cn`, `keywords_en`, `desc_cn`, `desc_en`, `css_js`, `html_cn`, `html_en`, `page`) VALUES
(NULL, 'GOFLY客服系统-商务版演示页',
'Customer Live Chat GOFLY-PRO-demo',
'GOFLY，GO-FLY',
'GOFLY，GO-FLY',
'一款开箱即用的在线客服系统',
'a free customer live chat',
'<style>body{color: #333;padding-left: 40px;}h1{font-size: 6em;}h2{font-size: 3em;font-weight: normal;}a{color: #333;}</style>',
'<h1>:)</h1><h2>你好 <a href="">GOFLY-PRO</a> 在线客服系统 !</h2><h3><a href="/login">后台</a>&nbsp;<a href="/index_en">English</a>&nbsp;<a href="/index_cn">中文</a></h3>',
'<h1>:)</h1><h2>HELLO <a href="">GOFLY-PRO</a> LIVE CHAT !</h2><h3><a href="/login">Administrator</a>&nbsp;<a href="/index_en">English</a>&nbsp;<a href="/index_cn">中文</a></h3>',
 'index')|
DROP TABLE IF EXISTS `reply_group`|
CREATE TABLE `reply_group` (
 `id` int(11) NOT NULL AUTO_INCREMENT,
 `group_name` varchar(50) NOT NULL DEFAULT '',
 `user_id` varchar(50) NOT NULL DEFAULT '',
 PRIMARY KEY (`id`),
 KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8|
DROP TABLE IF EXISTS `reply_item`|
CREATE TABLE `reply_item` (
 `id` int(11) NOT NULL AUTO_INCREMENT,
 `content` varchar(1024) NOT NULL DEFAULT '',
 `group_id` int(11) NOT NULL DEFAULT '0',
 `user_id` varchar(50) NOT NULL DEFAULT '',
 `item_name` varchar(50) NOT NULL DEFAULT '',
 PRIMARY KEY (`id`),
 KEY `user_id` (`user_id`),
 KEY `group_id` (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8|
DROP TABLE IF EXISTS `article_cate`|
CREATE TABLE `article_cate` (
 `id` int(11) NOT NULL AUTO_INCREMENT,
 `cat_name` varchar(50) NOT NULL DEFAULT '',
 `user_id` varchar(50) NOT NULL DEFAULT '',
 `ent_id` varchar(50) NOT NULL DEFAULT '',
 `is_top` tinyint NOT NULL DEFAULT '0',
 PRIMARY KEY (`id`),
 KEY `ent_id` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8|
DROP TABLE IF EXISTS `article`|
CREATE TABLE `article` (
 `id` int(11) NOT NULL AUTO_INCREMENT,
 `title` varchar(50) NOT NULL DEFAULT '',
 `content` text,
 `cat_id` int(11) NOT NULL DEFAULT '0',
 `user_id` varchar(50) NOT NULL DEFAULT '',
 `ent_id` varchar(50) NOT NULL DEFAULT '',
 PRIMARY KEY (`id`),
 KEY `ent_id` (`ent_id`),
 KEY `cat_id` (`cat_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8|
DROP TABLE IF EXISTS `land_page`|
CREATE TABLE `land_page` (
  `id` int(11) NOT NULL,
  `title` varchar(125) NOT NULL DEFAULT '',
  `keyword` varchar(255) NOT NULL DEFAULT '',
  `content` text NOT NULL,
  `language` varchar(50) NOT NULL DEFAULT '',
  `page_id` varchar(50) NOT NULL DEFAULT '',
   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci|
DROP TABLE IF EXISTS `language`|
CREATE TABLE `language` (
  `id` int(11) NOT NULL,
  `country` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `short_key` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci|
INSERT INTO `language` (`id`, `country`, `short_key`) VALUES (1, '中文简体', 'zh-cn')|
INSERT INTO `language` (`id`, `country`, `short_key`) VALUES (2, '正體中文', 'zh-tw')|
INSERT INTO `language` (`id`, `country`, `short_key`) VALUES (3, 'English', 'en_us')|
INSERT INTO `language` (`id`, `country`, `short_key`) VALUES (4, '日本語', 'ja_jp')|
DROP TABLE IF EXISTS `user_client`|
CREATE TABLE `user_client` (
 `id` int(11) NOT NULL AUTO_INCREMENT,
 `kefu` varchar(100) NOT NULL DEFAULT '',
 `client_id` varchar(100) NOT NULL DEFAULT '',
 `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
 UNIQUE KEY `idx_user` (`kefu`,`client_id`),
 PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8|
DROP TABLE IF EXISTS `tag`|
CREATE TABLE `tag` (
 `id` int(11) NOT NULL AUTO_INCREMENT,
 `name` varchar(100) NOT NULL DEFAULT '',
 `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `kefu` varchar(100) NOT NULL DEFAULT '',
 `ent_id` int(11) NOT NULL DEFAULT 0,
 PRIMARY KEY (`id`),
 KEY `name` (`name`),
 KEY `kefu` (`kefu`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4|
DROP TABLE IF EXISTS `visitor_tag`|
CREATE TABLE `visitor_tag` (
 `id` int(11) NOT NULL AUTO_INCREMENT,
 `visitor_id` varchar(100) NOT NULL DEFAULT '',
 `tag_id` int(11) NOT NULL DEFAULT 0,
 `ent_id` int(11) NOT NULL DEFAULT 0,
 `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `kefu` varchar(100) NOT NULL DEFAULT '',
 PRIMARY KEY (`id`),
 KEY `visitor_id` (`visitor_id`),
 KEY `tag_id` (`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4|
DROP TABLE IF EXISTS `ip_auth`|
CREATE TABLE `ip_auth` (
 `id` int(11) NOT NULL AUTO_INCREMENT,
 `content` varchar(100) NOT NULL DEFAULT '',
 `ip_address` varchar(100) NOT NULL DEFAULT '',
 `expire_time` varchar(50) NOT NULL DEFAULT '',
 `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `status` tinyint NOT NULL DEFAULT '1',
 PRIMARY KEY (`id`),
 KEY `ip_address` (`ip_address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8|
DROP TABLE IF EXISTS `oauth`|
CREATE TABLE `oauth` (
 `id` int(11) NOT NULL AUTO_INCREMENT,
 `user_id` varchar(100) NOT NULL DEFAULT '',
 `oauth_id` varchar(100) NOT NULL DEFAULT '',
 `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `status` tinyint NOT NULL DEFAULT '1',
 PRIMARY KEY (`id`),
 KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8|
DROP TABLE IF EXISTS `new`|
CREATE TABLE `new` (
 `id` int(11) NOT NULL AUTO_INCREMENT,
 `tag` varchar(100) NOT NULL DEFAULT '',
 `title` varchar(100) NOT NULL DEFAULT '',
 `content` text,
 `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `status` tinyint NOT NULL DEFAULT '1',
 PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4|
DROP TABLE IF EXISTS `visitor_black`|
CREATE TABLE `visitor_black` (
 `id` int(11) NOT NULL AUTO_INCREMENT,
 `visitor_id` varchar(100) NOT NULL DEFAULT '',
 `name` varchar(500) NOT NULL DEFAULT '',
 `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `ent_id` varchar(100) NOT NULL DEFAULT '',
 `kefu_name` varchar(100) NOT NULL DEFAULT '',
 PRIMARY KEY (`id`),
 UNIQUE KEY `visitor_id` (`visitor_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
-- 垂直分表：表结构平分，数据不变，每一个表中属性通过外键相关联
-- 水平分表：表结构不变，数据平分
-- 同一hash值的文件在不同用户的收藏夹里可以有不同的名字

-- 文件表
CREATE TABLE `file` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `file_hash` char(40) NOT NULL DEFAULT '' COMMENT '文件hash',
  `filename` char(40) NOT NULL DEFAULT '' COMMENT '文件名',
  `file_size` bigint(20) DEFAULT '0' COMMENT '文件大小',
  `file_local_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT '本机文件存储位置',
  `file_remote_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT '云端文件存储位置',
  `create_time` datetime default NOW() COMMENT '创建日期',
  `update_time` datetime default NOW() on update current_timestamp() COMMENT '更新日期',
  `status` int(11) NOT NULL DEFAULT '1' COMMENT '状态(0已删除/1正常/2禁用)',
  `ext1` int(11) DEFAULT '0' COMMENT '备用字段1',
  `ext2` text COMMENT '备用字段2',
  PRIMARY KEY (`id`),
  KEY `idx_file_hash` (`file_hash`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 用户表
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `server_addr` varchar(256) NOT NULL DEFAULT '' COMMENT '用户对应服务器的地址',
  `user_pwd` varchar(256) NOT NULL DEFAULT '' COMMENT '用户encoded密码',
  `user_email` varchar(64) DEFAULT '' COMMENT '邮箱',
  `user_phone` varchar(128) DEFAULT '' COMMENT '手机号',
  `email_validated` tinyint(1) DEFAULT 0 COMMENT '邮箱是否已验证',
  `phone_validated` tinyint(1) DEFAULT 0 COMMENT '手机号是否已验证',
  `sign_up_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '注册日期',
  `latest_update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最近活跃时间戳',
  `profile` text COMMENT '用户属性',
  `status` int(11) NOT NULL DEFAULT '1' COMMENT '账户状态(0已删除/1正常/2禁用)',
  PRIMARY KEY (`id`),
  KEY `idx_username` (`username`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 用户文件表
CREATE TABLE `user_file` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(64) NOT NULL,
  `filename` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
  `file_hash` varchar(64) NOT NULL DEFAULT '' COMMENT '文件hash',
  `upload_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  `latest_update_time` datetime DEFAULT CURRENT_TIMESTAMP
          ON UPDATE CURRENT_TIMESTAMP COMMENT '最近修改时间',
  `status` int(11) NOT NULL DEFAULT '1' COMMENT '文件状态(0已删除/1正常/2禁用)',
  PRIMARY KEY (`id`),
  KEY `idx_filename` (`filename`),
  KEY `idx_status` (`status`),
  KEY `idx_user_id` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 用户登录时间表
CREATE TABLE `user_login` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `login_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '用户登录时间',
  `token` varchar(100) NOT NULL DEFAULT '' COMMENT '用户登录token',
  `server_addr` varchar(64) NOT NULL DEFAULT '' COMMENT '用户登录所绑定的服务器地址',
  PRIMARY KEY (`id`),
  KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 文件与服务器对应表
CREATE TABLE `file_server` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `file_hash` varchar(64) NOT NULL DEFAULT '' COMMENT '文件名',
  `server_addr` varchar(256) NOT NULL DEFAULT '' COMMENT '服务器的地址',
  `upload_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  `status` int(11) NOT NULL DEFAULT '1' COMMENT '账户状态(0已删除/1正常/2禁用)',
  PRIMARY KEY (`id`),
  KEY `idx_file_hash` (`file_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
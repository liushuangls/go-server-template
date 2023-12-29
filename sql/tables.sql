CREATE TABLE `users` (
 `id` bigint NOT NULL AUTO_INCREMENT,
 `create_time` datetime NOT NULL,
 `update_time` datetime NOT NULL,
 `delete_time` datetime DEFAULT NULL,
 `nickname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
 `register_type` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
 `register_ip` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
 `register_region` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
 `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
 `email_verified` tinyint(1) NOT NULL DEFAULT '0',
 `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
 `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
 `profile` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
 PRIMARY KEY (`id`),
 KEY `user_create_time` (`create_time`),
 KEY `user_delete_time` (`delete_time`),
 KEY `user_email` (`email`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

CREATE TABLE `user_oauth` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL,
  `delete_time` datetime DEFAULT NULL,
  `platform` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `open_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `union_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `user_id` bigint NOT NULL,
  PRIMARY KEY (`id`),
  KEY `useroauth_create_time` (`create_time`),
  KEY `useroauth_delete_time` (`delete_time`),
  KEY `useroauth_user_id` (`user_id`),
  KEY `useroauth_open_id` (`open_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

CREATE TABLE `server_logs` (
   `id` bigint NOT NULL AUTO_INCREMENT,
   `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
   `user_id` int DEFAULT '0',
   `ip` varchar(128) default '',
   `path` varchar(256) default '',
   `query` text,
   `body` text,
   `method` varchar(32) default '',
   `level` varchar(32) DEFAULT 'error',
   `from` varchar(32) DEFAULT 'api',
   `err_msg` varchar(4096) NOT NULL,
   `resp_err_msg` varchar(4096) default '',
   `extra` json DEFAULT NULL,
   `code` int DEFAULT '0',
   PRIMARY KEY (`id`),
   KEY `create_time` (`create_time`),
   KEY `user_id` (`user_id`),
   KEY `level` (`level`),
   KEY `code` (`code`),
   Key `path` (`path`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

create table app_configs (
     id bigint not null auto_increment primary key,
     `create_time` datetime NOT NULL,
     `update_time` datetime NOT NULL,
     `delete_time` datetime DEFAULT NULL,
     `key` varchar(64) not null,
     `value` text,
     value_type varchar(32) not null default 'string',
     `type` varchar(32) not null default 'client',
     `app_name` varchar(64) default '',
     app_version_gte varchar(64) default '',
     app_version_lte varchar(64) default '',
     unique key_app (`key`, `app_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

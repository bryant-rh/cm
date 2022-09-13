DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(100) NOT NULL DEFAULT '' unique,
  `password` varchar(100) NOT NULL DEFAULT '',
    `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
 
DROP TABLE IF EXISTS `project`;
CREATE TABLE `project` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `project_id` varchar(50)  NOT NULL DEFAULT '0' unique,
  `project_name` varchar(225) NOT NULL DEFAULT '',
  `describe` varchar(512) NOT NULL DEFAULT '',
    `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
 
DROP TABLE IF EXISTS `cluster`;
CREATE TABLE `cluster` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `cluster_id` varchar(50) NOT NULL DEFAULT '0' unique,
  `cluster_name` varchar(225) NOT NULL DEFAULT '',
    `describe` varchar(512) NOT NULL DEFAULT '',
    `labels` varchar(512) NOT NULL DEFAULT '',
    `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
 
DROP TABLE IF EXISTS `serviceaccount`;
CREATE TABLE `serviceaccount` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `sa_id` varchar(50) NOT NULL DEFAULT '0' unique,
  `sa_name` varchar(225) NOT NULL DEFAULT '',
    `sa_token` varchar(2048) NOT NULL DEFAULT '',
    `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
 
DROP TABLE IF EXISTS `namespace`;
CREATE TABLE `namespace` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ns_id` varchar(50) NOT NULL DEFAULT '0' unique,
  `ns_name` varchar(225) NOT NULL DEFAULT '',
    `cluster_id` varchar(50) NOT NULL DEFAULT '0',
    `sa_id` varchar(50) NOT NULL DEFAULT '0',
    `sa_name` varchar(225) NOT NULL DEFAULT '',
    `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
 
 
DROP TABLE IF EXISTS `label`;
CREATE TABLE `label` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `label_id` varchar(50) NOT NULL DEFAULT '0' unique,
  `key` varchar(225) NOT NULL DEFAULT '',
    `value` varchar(512) NOT NULL DEFAULT '',
    `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
 
DROP TABLE IF EXISTS `cluster_bind`;
CREATE TABLE `cluster_bind`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `project_id` varchar(50) NOT NULL DEFAULT '0',
    `cluster_id` varchar(50) NOT NULL DEFAULT '0',
    `project_name` varchar(225) NOT NULL DEFAULT '',
    `cluster_name` varchar(225) NOT NULL DEFAULT '',
    `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
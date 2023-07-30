/*
 Navicat Premium Data Transfer

 Source Server         : 51.81.192.209
 Source Server Type    : MySQL
 Source Server Version : 50739
 Source Host           : 51.81.192.209:3306
 Source Schema         : jihuoma

 Target Server Type    : MySQL
 Target Server Version : 50739
 File Encoding         : 65001

 Date: 25/10/2022 23:07:22
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for kpm_read_user
-- ----------------------------
DROP TABLE IF EXISTS `kpm_read_user`;
CREATE TABLE `kpm_read_user` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `username` varchar(30) NOT NULL,
  `password` varchar(64) NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '0停用1启用',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of kpm_read_user
-- ----------------------------
BEGIN;
INSERT INTO `kpm_read_user` VALUES (1, 'admin', 'admin888', 1);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;

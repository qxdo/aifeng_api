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

 Date: 25/10/2022 23:07:07
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for kpm_read_price
-- ----------------------------
DROP TABLE IF EXISTS `kpm_read_price`;
CREATE TABLE `kpm_read_price` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `baidan` decimal(10,2) NOT NULL,
  `yedan` decimal(10,2) NOT NULL,
  `uid` int(10) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of kpm_read_price
-- ----------------------------
BEGIN;
INSERT INTO `kpm_read_price` VALUES (4, 20.00, 12.00, 1);
INSERT INTO `kpm_read_price` VALUES (5, 1.00, 12.00, 1);
INSERT INTO `kpm_read_price` VALUES (6, 1.00, 12.00, 1);
INSERT INTO `kpm_read_price` VALUES (7, 1.00, 12.00, 1);
INSERT INTO `kpm_read_price` VALUES (8, 20.00, 12.00, 1);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;

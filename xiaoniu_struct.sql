/*
 Navicat Premium Data Transfer

 Source Server         : xiaoniu_server
 Source Server Type    : MySQL
 Source Server Version : 80031
 Source Host           : 49.235.113.82:3306
 Source Schema         : xiaoniu

 Target Server Type    : MySQL
 Target Server Version : 80031
 File Encoding         : 65001

 Date: 30/10/2022 16:10:08
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for kpm_read
-- ----------------------------
DROP TABLE IF EXISTS `kpm_read`;
CREATE TABLE `kpm_read`  (
  `id` int(0) NOT NULL AUTO_INCREMENT,
  `uid` int(0) NOT NULL,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '标题',
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '类型',
  `price` decimal(10, 2) NOT NULL,
  `demand_count` int(0) NOT NULL DEFAULT 0 COMMENT '需求刷单量',
  `before_count` int(0) NOT NULL DEFAULT 0 COMMENT '刷之前数量',
  `all_count` int(0) NOT NULL DEFAULT 0 COMMENT '刷完之后总数',
  `suc_count` int(0) NULL DEFAULT 0 COMMENT '成功条数',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '0全部1启用0停用2进行中3无效4完成',
  `add_time` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `end_time` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '完成时间',
  `is_first` tinyint(1) NOT NULL DEFAULT 0 COMMENT '0不用1用',
  `start` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 226 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for kpm_read_price
-- ----------------------------
DROP TABLE IF EXISTS `kpm_read_price`;
CREATE TABLE `kpm_read_price`  (
  `id` int(0) NOT NULL AUTO_INCREMENT,
  `baidan` decimal(10, 2) NOT NULL,
  `yedan` decimal(10, 2) NOT NULL,
  `uid` int(0) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for kpm_read_user
-- ----------------------------
DROP TABLE IF EXISTS `kpm_read_user`;
CREATE TABLE `kpm_read_user`  (
  `id` int(0) NOT NULL AUTO_INCREMENT,
  `username` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `password` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '0停用1启用',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for proxy_sleep
-- ----------------------------
DROP TABLE IF EXISTS `proxy_sleep`;
CREATE TABLE `proxy_sleep`  (
  `id` int(0) NOT NULL AUTO_INCREMENT,
  `guid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `proxy` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `count` int(0) NULL DEFAULT NULL,
  `time` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4158 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;

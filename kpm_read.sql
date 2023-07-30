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

 Date: 25/10/2022 23:06:51
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for kpm_read
-- ----------------------------
DROP TABLE IF EXISTS `kpm_read`;
CREATE TABLE `kpm_read` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `uid` int(10) NOT NULL,
  `title` varchar(255) DEFAULT NULL COMMENT '标题',
  `link` varchar(255) NOT NULL,
  `type` varchar(32) NOT NULL COMMENT '类型',
  `price` decimal(10,2) NOT NULL,
  `demand_count` int(10) NOT NULL DEFAULT '0' COMMENT '需求刷单量',
  `befor_count` int(10) NOT NULL DEFAULT '0' COMMENT '刷之前数量',
  `allcount` int(10) NOT NULL DEFAULT '0' COMMENT '刷完之后总数',
  `suc_count` int(10) NOT NULL DEFAULT '0' COMMENT '成功条数',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '0全部1启用0停用2进行中3无效4完成',
  `addtime_string` varchar(255) DEFAULT NULL,
  `endtime_string` varchar(255) DEFAULT NULL COMMENT '完成时间',
  `is_first` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0不用1用',
  `start_string` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=84 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of kpm_read
-- ----------------------------
BEGIN;
INSERT INTO `kpm_read` VALUES (76, 1, NULL, 'https://mp.weixin.qq.com/s/-r3a9Vk0uDTDSCsG5f6WVA', '白单', 20.00, 20, 132, 173, 20, 4, '2022-10-25 16:36:35', '2022-10-25 16:37:48', 0, '2022-10-25 16:37:27');
INSERT INTO `kpm_read` VALUES (77, 1, NULL, 'https://mp.weixin.qq.com/s/_J-0f78RL5cmF0HKlLTKcA', '白单', 20.00, 100, 1396, 1523, 127, 4, '2022-10-25 16:51:09', '2022-10-25 16:51:48', 0, '2022-10-25 16:51:10');
INSERT INTO `kpm_read` VALUES (78, 1, NULL, 'https://mp.weixin.qq.com/s/_J-0f78RL5cmF0HKlLTKcA', '白单', 20.00, 100, 1523, 1588, 65, 4, '2022-10-25 16:57:45', '2022-10-25 16:58:12', 0, '2022-10-25 16:57:46');
INSERT INTO `kpm_read` VALUES (79, 1, NULL, 'https://mp.weixin.qq.com/s/_J-0f78RL5cmF0HKlLTKcA', '白单', 20.00, 100, 1583, 1588, 3, 4, '2022-10-25 16:58:08', '2022-10-25 16:58:25', 0, '2022-10-25 16:58:09');
INSERT INTO `kpm_read` VALUES (80, 1, NULL, 'https://mp.weixin.qq.com/s/Y0WYwP059P2NKOLEvrJ7Bg', '白单', 20.00, 1, 865, 866, 1, 4, '2022-10-25 17:16:10', '2022-10-25 17:16:14', 0, '2022-10-25 17:16:10');
INSERT INTO `kpm_read` VALUES (81, 1, NULL, 'https://mp.weixin.qq.com/s/_vPg5cretF99XH1bE39P2w', '白单', 20.00, 10, 2054, 2074, 20, 4, '2022-10-25 17:59:38', '2022-10-25 18:00:01', 0, '2022-10-25 17:59:39');
INSERT INTO `kpm_read` VALUES (82, 1, NULL, 'https://mp.weixin.qq.com/s/_vPg5cretF99XH1bE39P2w', '白单', 20.00, 1, 2102, 2103, 1, 4, '2022-10-25 19:02:20', '2022-10-25 19:02:24', 0, '2022-10-25 19:02:21');
INSERT INTO `kpm_read` VALUES (83, 1, NULL, 'https://mp.weixin.qq.com/s/Y0WYwP059P2NKOLEvrJ7Bg', '白单', 20.00, 1, 957, 958, 1, 4, '2022-10-25 19:13:32', '2022-10-25 19:13:36', 0, '2022-10-25 19:13:32');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;

/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80012
 Source Host           : localhost:3306
 Source Schema         : gstep

 Target Server Type    : MySQL
 Target Server Version : 80012
 File Encoding         : 65001

 Date: 16/03/2025 20:50:26
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for hr_department
-- ----------------------------
DROP TABLE IF EXISTS `hr_department`;
CREATE TABLE `hr_department`  (
  `id` int(11) NOT NULL,
  `name` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `parent_id` int(11) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of hr_department
-- ----------------------------
INSERT INTO `hr_department` VALUES (1, 'gstep公司', 0);
INSERT INTO `hr_department` VALUES (2, '信息技术部', 1);
INSERT INTO `hr_department` VALUES (3, '网络室', 2);
INSERT INTO `hr_department` VALUES (4, '开发室', 2);
INSERT INTO `hr_department` VALUES (5, '人力资源部', 1);

-- ----------------------------
-- Table structure for hr_user
-- ----------------------------
DROP TABLE IF EXISTS `hr_user`;
CREATE TABLE `hr_user`  (
  `id` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '用户id',
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '姓名',
  `position` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '0' COMMENT '职位',
  `position_code` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '职位编码',
  `is_leader` tinyint(4) NOT NULL COMMENT '是否是部门负责人',
  `department_id` int(11) NOT NULL DEFAULT 0 COMMENT '部门id',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of hr_user
-- ----------------------------
INSERT INTO `hr_user` VALUES ('001', '刘协', '董事长', '01', 1, 1);
INSERT INTO `hr_user` VALUES ('101', '刘备', '主任', '02', 1, 3);
INSERT INTO `hr_user` VALUES ('102', '关羽', '网络员', '03', 0, 3);
INSERT INTO `hr_user` VALUES ('103', '张飞', '网络员', '03', 0, 3);
INSERT INTO `hr_user` VALUES ('201', '孙权', '主任', '02', 1, 4);
INSERT INTO `hr_user` VALUES ('202', '周瑜', '程序员', '04', 0, 4);
INSERT INTO `hr_user` VALUES ('301', '曹操', '主任', '02', 1, 5);
INSERT INTO `hr_user` VALUES ('302', '郭嘉', '人事员', '05', 0, 5);

-- ----------------------------
-- Table structure for process
-- ----------------------------
DROP TABLE IF EXISTS `process`;
CREATE TABLE `process`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `template_id` int(11) NULL DEFAULT NULL,
  `start_user_id` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `state` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '0' COMMENT '状态:started,finished',
  `finished_at` datetime(0) NULL DEFAULT NULL,
  `created_at` datetime(0) NULL DEFAULT NULL,
  `updated_at` datetime(0) NULL DEFAULT NULL,
  `deleted_at` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_id`(`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 97 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for task
-- ----------------------------
DROP TABLE IF EXISTS `task`;
CREATE TABLE `task`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `process_id` int(11) NULL DEFAULT NULL,
  `step_id` int(11) NULL DEFAULT NULL,
  `title` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `state` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '状态:unstart,started,pass,refuse,withdraw',
  `category` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `audit_method` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `form` varchar(6000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `candidates` varchar(4000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `created_at` datetime(0) NULL DEFAULT NULL,
  `updated_at` datetime(0) NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(0),
  `deleted_at` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 169 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for task_assignee
-- ----------------------------
DROP TABLE IF EXISTS `task_assignee`;
CREATE TABLE `task_assignee`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `process_id` int(11) NULL DEFAULT NULL,
  `step_id` int(11) NULL DEFAULT NULL,
  `task_id` int(11) NULL DEFAULT NULL,
  `user_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `state` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '状态:started,pass,refuse',
  `submit_index` int(11) NULL DEFAULT NULL,
  `form` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `created_at` datetime(0) NULL DEFAULT NULL,
  `updated_at` datetime(0) NULL DEFAULT NULL,
  `deleted_at` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 226 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for template
-- ----------------------------
DROP TABLE IF EXISTS `template`;
CREATE TABLE `template`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `mould_id` int(11) NULL DEFAULT NULL,
  `title` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `version` int(5) NULL DEFAULT NULL,
  `root_step` varchar(10000) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `fields` varchar(4000) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `created_at` datetime(0) NULL DEFAULT NULL,
  `updated_at` datetime(0) NULL DEFAULT NULL,
  `deleted_at` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_id`(`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of template
-- ----------------------------
INSERT INTO `template` VALUES (63, 2, '测试', 1, '{\"id\":1,\"title\":\"申请\",\"category\":\"start\",\"form\":null,\"candidates\":[{\"category\":\"field\",\"title\":\"申请人\",\"value\":\"applicant\"}],\"expression\":\"\",\"auditMethod\":\"or\",\"branchSteps\":[],\"nextStep\":{\"id\":5,\"title\":\"抄送申请人\",\"category\":\"notify\",\"form\":{},\"candidates\":[{\"category\":\"field\",\"title\":\"申请人\",\"value\":\"applicant\"}],\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":[],\"nextStep\":{\"id\":2,\"title\":\"科室、车间/作业区负责人确认内容\",\"category\":\"audit\",\"form\":null,\"candidates\":[{\"category\":\"field\",\"title\":\"科室、车间/作业区负责人\",\"value\":\"workshop_supervisor_id\"}],\"expression\":\"\",\"auditMethod\":\"or\",\"branchSteps\":[],\"nextStep\":{\"id\":6,\"title\":\"抄送科室、车间/作业区负责人\",\"category\":\"notify\",\"form\":{},\"candidates\":[{\"category\":\"field\",\"title\":\"科室、车间/作业区负责人\",\"value\":\"workshop_supervisor_id\"}],\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":[],\"nextStep\":{\"id\":3,\"title\":\"科室主任/建造经理/总建造师确认内容\",\"category\":\"audit\",\"form\":null,\"candidates\":[{\"category\":\"field\",\"title\":\"科室主任/建造经理/总建造师\",\"value\":\"constructor_id\"}],\"expression\":\"\",\"auditMethod\":\"or\",\"branchSteps\":[],\"nextStep\":{\"id\":7,\"title\":\"抄送科室主任/建造经理/总建造师\",\"category\":\"notify\",\"form\":{},\"candidates\":[{\"category\":\"field\",\"title\":\"科室主任/建造经理/总建造师\",\"value\":\"constructor_id\"}],\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":[],\"nextStep\":{\"id\":4,\"title\":\"安全管理人员确认\",\"category\":\"audit\",\"form\":{},\"candidates\":[{\"category\":\"field\",\"title\":\"安全管理人员\",\"value\":\"safety_manager_id\"}],\"expression\":\"\",\"auditMethod\":\"or\",\"branchSteps\":[],\"nextStep\":{\"id\":8,\"title\":\"抄送安全管理人员\",\"category\":\"notify\",\"form\":{},\"candidates\":[{\"category\":\"field\",\"title\":\"安全管理人员\",\"value\":\"safety_manager_id\"}],\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":[],\"nextStep\":{\"id\":9999,\"title\":\"结束\",\"category\":\"end\",\"form\":{},\"candidates\":[],\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":[],\"nextStep\":{\"id\":0,\"title\":\"\",\"category\":\"\",\"form\":null,\"candidates\":null,\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":null,\"nextStep\":null}}}}}}}}}}', '[{\"title\":\"申请人\",\"name\":\"applicant\"},{\"title\":\"科室、车间/作业区负责人\",\"name\":\"workshop_supervisor_id\"},{\"title\":\"安全管理人员\",\"name\":\"safety_manager_id\"},{\"title\":\"科室主任/建造经理/总建造师\",\"name\":\"constructor_id\"}]', '2025-03-12 09:02:10', '2025-03-16 20:35:25', NULL);

-- ----------------------------
-- View structure for department
-- ----------------------------
DROP VIEW IF EXISTS `department`;
CREATE ALGORITHM = UNDEFINED DEFINER = `root`@`localhost` SQL SECURITY DEFINER VIEW `department` AS with recursive `ct` as (select `a`.`id` AS `id`,`a`.`name` AS `name`,`a`.`parent_id` AS `parent_id` from `hr_department` `a` where (`a`.`id` = 1) union all select `b`.`id` AS `id`,`b`.`name` AS `name`,`b`.`parent_id` AS `parent_id` from (`hr_department` `b` join `ct` on((`ct`.`id` = `b`.`parent_id`)))) select concat(`ct`.`id`,'') AS `id`,`ct`.`name` AS `name`,concat(`ct`.`parent_id`,'') AS `parent_id` from `ct`;

-- ----------------------------
-- View structure for position
-- ----------------------------
DROP VIEW IF EXISTS `position`;
CREATE ALGORITHM = UNDEFINED DEFINER = `root`@`localhost` SQL SECURITY DEFINER VIEW `position` AS select distinct `hr_user`.`position` AS `title`,`hr_user`.`position_code` AS `code` from `hr_user` where ((length(ifnull(`hr_user`.`position_code`,'')) > 0) and (length(ifnull(`hr_user`.`position`,'')) > 0));

-- ----------------------------
-- View structure for user
-- ----------------------------
DROP VIEW IF EXISTS `user`;
CREATE ALGORITHM = UNDEFINED DEFINER = `root`@`localhost` SQL SECURITY DEFINER VIEW `user` AS select `hr_user`.`id` AS `id`,`hr_user`.`name` AS `name`,`hr_user`.`position` AS `position_title`,`hr_user`.`position_code` AS `position_code`,`hr_user`.`is_leader` AS `is_leader`,`hr_user`.`department_id` AS `department_id` from `hr_user`;

SET FOREIGN_KEY_CHECKS = 1;

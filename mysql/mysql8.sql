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

 Date: 22/03/2025 23:53:49
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for executor
-- ----------------------------
DROP TABLE IF EXISTS `executor`;
CREATE TABLE `executor`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `process_id` int(11) NULL DEFAULT NULL,
  `step_id` int(11) NULL DEFAULT NULL,
  `task_id` int(11) NULL DEFAULT NULL,
  `user_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `state` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '状态:unstart,started,pass,refuse',
  `submit_index` int(11) NULL DEFAULT NULL,
  `form` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `memo` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '备注',
  `created_at` datetime(0) NULL DEFAULT NULL,
  `updated_at` datetime(0) NULL DEFAULT NULL,
  `deleted_at` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 883 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of executor
-- ----------------------------
INSERT INTO `executor` VALUES (872, 458, 1, 2640, '曹操', 'pass', 1, '{\"applicant\":\"曹操\",\"receiverDeptLeaderId\":[\"张飞\",\"赵云\"],\"receiverId\":\"关羽\",\"safeDeptAuditorId\":\"刘备\"}', '备注1', '2025-03-22 07:53:44', '2025-03-22 07:53:44', NULL);
INSERT INTO `executor` VALUES (873, 458, 5, 2641, '曹操', 'pass', 1, '{\"applicant\":\"曹操\",\"receiverDeptLeaderId\":[\"张飞\",\"赵云\"],\"receiverId\":\"关羽\",\"safeDeptAuditorId\":\"刘备\"}', '', '2025-03-22 07:53:44', '2025-03-22 07:53:44', NULL);
INSERT INTO `executor` VALUES (874, 458, 2, 2642, '关羽', 'pass', 1, '{\"applicant\":\"曹操\",\"receiverDeptLeaderId\":[\"张飞\",\"赵云\"],\"receiverId\":\"关羽\",\"safeDeptAuditorId\":\"刘备\"}', '', '2025-03-22 07:53:44', '2025-03-22 07:57:30', NULL);
INSERT INTO `executor` VALUES (879, 458, 6, 2646, '关羽', 'pass', 1, '{\"applicant\":\"曹操\",\"receiverDeptLeaderId\":[\"马超\",\"魏延\"],\"receiverId\":\"关羽\",\"safeDeptAuditorId\":\"刘备\"}', '', '2025-03-22 07:57:30', '2025-03-22 07:57:30', NULL);
INSERT INTO `executor` VALUES (880, 458, 3, 2647, '马超', 'pass', 1, '{\"applicant\":\"曹操\",\"receiverDeptLeaderId\":[\"马超\",\"魏延\"],\"receiverId\":\"关羽\",\"safeDeptAuditorId\":\"刘备\"}', '', '2025-03-22 07:57:30', '2025-03-22 08:03:51', NULL);
INSERT INTO `executor` VALUES (881, 458, 3, 2647, '魏延', 'started', 2, '{\"applicant\":\"曹操\",\"receiverDeptLeaderId\":[\"马超\",\"魏延\"],\"receiverId\":\"关羽\",\"safeDeptAuditorId\":\"刘备\"}', '', '2025-03-22 07:57:30', '2025-03-22 07:57:30', NULL);
INSERT INTO `executor` VALUES (883, 458, 4, 2649, '刘备', 'started', 1, '{\"applicant\":\"曹操\",\"receiverDeptLeaderId\":[\"马超\",\"魏延\"],\"receiverId\":\"关羽\",\"safeDeptAuditorId\":\"刘备\"}', '', '2025-03-22 08:03:51', '2025-03-22 08:03:51', NULL);

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
INSERT INTO `hr_user` VALUES ('关羽', '关羽', '网络员', '03', 0, 3);
INSERT INTO `hr_user` VALUES ('刘协', '刘协', '董事长', '01', 1, 1);
INSERT INTO `hr_user` VALUES ('刘备', '刘备', '主任', '02', 1, 3);
INSERT INTO `hr_user` VALUES ('周瑜', '周瑜', '程序员', '04', 0, 4);
INSERT INTO `hr_user` VALUES ('孙权', '孙权', '主任', '02', 1, 4);
INSERT INTO `hr_user` VALUES ('张飞', '张飞', '网络员', '03', 0, 3);
INSERT INTO `hr_user` VALUES ('曹操', '曹操', '主任', '02', 1, 5);
INSERT INTO `hr_user` VALUES ('赵云', '赵云', '程序员', '02', 0, 3);
INSERT INTO `hr_user` VALUES ('郭嘉', '郭嘉', '人事员', '05', 0, 5);
INSERT INTO `hr_user` VALUES ('马超', '马超', '程序员', '02', 0, 3);
INSERT INTO `hr_user` VALUES ('魏延', '魏延', '程序员', '02', 0, 3);

-- ----------------------------
-- Table structure for mould
-- ----------------------------
DROP TABLE IF EXISTS `mould`;
CREATE TABLE `mould`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `created_at` datetime(0) NULL DEFAULT NULL,
  `updated_at` datetime(0) NULL DEFAULT NULL,
  `deleted_at` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_id`(`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 67 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of mould
-- ----------------------------
INSERT INTO `mould` VALUES (2, '有限空间作业审批单', '2025-03-20 09:42:11', '2025-03-20 09:42:11', NULL);
INSERT INTO `mould` VALUES (3, '访客预约单', '2025-03-20 09:42:11', '2025-03-20 09:42:11', NULL);
INSERT INTO `mould` VALUES (4, '船舶建造动火确认单', '2025-03-20 09:42:11', '2025-03-21 09:24:04', NULL);

-- ----------------------------
-- Table structure for process
-- ----------------------------
DROP TABLE IF EXISTS `process`;
CREATE TABLE `process`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `template_id` int(11) NULL DEFAULT NULL,
  `start_user_id` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `state` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '0' COMMENT '状态:started,finish_pass',
  `finished_at` datetime(0) NULL DEFAULT NULL,
  `created_at` datetime(0) NULL DEFAULT NULL,
  `updated_at` datetime(0) NULL DEFAULT NULL,
  `deleted_at` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_id`(`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 458 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of process
-- ----------------------------
INSERT INTO `process` VALUES (458, 64, '曹操', 'started', NULL, '2025-03-22 07:53:44', '2025-03-22 07:53:44', NULL);

-- ----------------------------
-- Table structure for task
-- ----------------------------
DROP TABLE IF EXISTS `task`;
CREATE TABLE `task`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `process_id` int(11) NULL DEFAULT NULL,
  `step_id` int(11) NULL DEFAULT NULL,
  `title` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `category` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `state` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '状态:unstart,started,pass,refuse',
  `audit_method` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `form` varchar(6000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `created_at` datetime(0) NULL DEFAULT NULL,
  `updated_at` datetime(0) NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(0),
  `deleted_at` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2649 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of task
-- ----------------------------
INSERT INTO `task` VALUES (2640, 458, 1, '访客预约申请', 'start', 'pass', 'or', '{\"applicant\":\"曹操\",\"receiverDeptLeaderId\":[\"张飞\",\"赵云\"],\"receiverId\":\"关羽\",\"safeDeptAuditorId\":\"刘备\"}', '2025-03-22 07:53:44', '2025-03-22 07:53:44', NULL);
INSERT INTO `task` VALUES (2641, 458, 5, '抄送预约人', 'notify', 'pass', '', '{\"applicant\":\"曹操\",\"receiverDeptLeaderId\":[\"张飞\",\"赵云\"],\"receiverId\":\"关羽\",\"safeDeptAuditorId\":\"刘备\"}', '2025-03-22 07:53:44', '2025-03-22 07:53:44', NULL);
INSERT INTO `task` VALUES (2642, 458, 2, '接待人审核', 'audit', 'pass', 'or', '{\"applicant\":\"曹操\",\"receiverDeptLeaderId\":[\"马超\",\"魏延\"],\"receiverId\":\"关羽\",\"safeDeptAuditorId\":\"刘备\"}', '2025-03-22 07:53:44', '2025-03-22 07:57:30', NULL);
INSERT INTO `task` VALUES (2646, 458, 6, '抄送接待人', 'notify', 'pass', '', '{\"applicant\":\"曹操\",\"receiverDeptLeaderId\":[\"马超\",\"魏延\"],\"receiverId\":\"关羽\",\"safeDeptAuditorId\":\"刘备\"}', '2025-03-22 07:57:30', '2025-03-22 07:57:30', NULL);
INSERT INTO `task` VALUES (2647, 458, 3, '接待部门领导审核', 'audit', 'pass', 'or', '{\"applicant\":\"曹操\",\"receiverDeptLeaderId\":[\"马超\",\"魏延\"],\"receiverId\":\"关羽\",\"safeDeptAuditorId\":\"刘备\"}', '2025-03-22 07:57:30', '2025-03-22 08:03:51', NULL);
INSERT INTO `task` VALUES (2649, 458, 4, '安保部审核', 'audit', 'started', 'or', '{\"applicant\":\"曹操\",\"receiverDeptLeaderId\":[\"马超\",\"魏延\"],\"receiverId\":\"关羽\",\"safeDeptAuditorId\":\"刘备\"}', '2025-03-22 08:03:51', '2025-03-22 08:03:51', NULL);

-- ----------------------------
-- Table structure for template
-- ----------------------------
DROP TABLE IF EXISTS `template`;
CREATE TABLE `template`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `mould_id` int(11) NULL DEFAULT NULL,
  `title` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `version` int(11) NULL DEFAULT NULL,
  `root_step` varchar(10000) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `fields` varchar(4000) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `state` varchar(40) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'edit' COMMENT '状态: 发布 release  草稿 draft',
  `created_at` datetime(0) NULL DEFAULT NULL,
  `updated_at` datetime(0) NULL DEFAULT NULL,
  `deleted_at` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_id`(`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 67 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of template
-- ----------------------------
INSERT INTO `template` VALUES (1, 1, '默认', 1, '{\"id\":1,\"title\":\"访客预约申请\",\"category\":\"start\",\"auth\":{},\"expression\":\"\",\"auditMethod\":\"or\",\"branchSteps\":[],\"nextStep\":{\"id\":9999,\"title\":\"结束\",\"category\":\"end\",\"auth\":null,\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":[],\"nextStep\":{\"id\":0,\"title\":\"\",\"category\":\"\",\"auth\":null,\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":null,\"nextStep\":null}}}', '[]', 'release', '2025-03-22 15:21:15', '2025-03-22 15:43:57', NULL);
INSERT INTO `template` VALUES (63, 2, '有限空间作业审批单', 1, '{\"id\":1,\"title\":\"申请\",\"category\":\"start\",\"form\":null,\"candidates\":[{\"category\":\"field\",\"title\":\"申请人\",\"value\":\"applicant\"}],\"expression\":\"\",\"auditMethod\":\"or\",\"branchSteps\":[],\"nextStep\":{\"id\":5,\"title\":\"抄送申请人\",\"category\":\"notify\",\"form\":{},\"candidates\":[{\"category\":\"field\",\"title\":\"申请人\",\"value\":\"applicant\"}],\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":[],\"nextStep\":{\"id\":2,\"title\":\"科室、车间/作业区负责人确认内容\",\"category\":\"audit\",\"form\":null,\"candidates\":[{\"category\":\"field\",\"title\":\"科室、车间/作业区负责人\",\"value\":\"workshop_supervisor_id\"}],\"expression\":\"\",\"auditMethod\":\"or\",\"branchSteps\":[],\"nextStep\":{\"id\":6,\"title\":\"抄送科室、车间/作业区负责人\",\"category\":\"notify\",\"form\":{},\"candidates\":[{\"category\":\"field\",\"title\":\"科室、车间/作业区负责人\",\"value\":\"workshop_supervisor_id\"}],\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":[],\"nextStep\":{\"id\":3,\"title\":\"科室主任/建造经理/总建造师确认内容\",\"category\":\"audit\",\"form\":null,\"candidates\":[{\"category\":\"field\",\"title\":\"科室主任/建造经理/总建造师\",\"value\":\"constructor_id\"}],\"expression\":\"\",\"auditMethod\":\"or\",\"branchSteps\":[],\"nextStep\":{\"id\":7,\"title\":\"抄送科室主任/建造经理/总建造师\",\"category\":\"notify\",\"form\":{},\"candidates\":[{\"category\":\"field\",\"title\":\"科室主任/建造经理/总建造师\",\"value\":\"constructor_id\"}],\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":[],\"nextStep\":{\"id\":4,\"title\":\"安全管理人员确认\",\"category\":\"audit\",\"form\":{},\"candidates\":[{\"category\":\"field\",\"title\":\"安全管理人员\",\"value\":\"safety_manager_id\"}],\"expression\":\"\",\"auditMethod\":\"or\",\"branchSteps\":[],\"nextStep\":{\"id\":8,\"title\":\"抄送安全管理人员\",\"category\":\"notify\",\"form\":{},\"candidates\":[{\"category\":\"field\",\"title\":\"安全管理人员\",\"value\":\"safety_manager_id\"}],\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":[],\"nextStep\":{\"id\":9999,\"title\":\"结束\",\"category\":\"end\",\"form\":{},\"candidates\":[],\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":[],\"nextStep\":{\"id\":0,\"title\":\"\",\"category\":\"\",\"form\":null,\"candidates\":null,\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":null,\"nextStep\":null}}}}}}}}}}', '[{\"title\":\"申请人\",\"name\":\"applicant\"},{\"title\":\"科室、车间/作业区负责人\",\"name\":\"workshop_supervisor_id\"},{\"title\":\"安全管理人员\",\"name\":\"safety_manager_id\"},{\"title\":\"科室主任/建造经理/总建造师\",\"name\":\"constructor_id\"}]', 'release', '2025-03-12 09:02:10', '2025-03-16 20:35:25', NULL);
INSERT INTO `template` VALUES (64, 3, '访客预约单', 1, '{\"id\":1,\"title\":\"访客预约申请\",\"category\":\"start\",\"candidates\":[],\"auth\":{\"applicant\":{\"title\":\"申请人\",\"value\":\"editable\"},\"receiverDeptLeaderId\":{\"title\":\"接待单位领导1\",\"value\":\"\"},\"receiverId\":{\"title\":\"接待人\",\"value\":\"editable\"},\"safeDeptAuditorId\":{\"title\":\"安保部审核人1\",\"value\":\"\"}},\"expression\":\"\",\"auditMethod\":\"or\",\"branchSteps\":[],\"nextStep\":{\"id\":5,\"title\":\"抄送预约人\",\"category\":\"notify\",\"candidates\":[],\"auth\":{\"applicant\":{\"title\":\"申请人\",\"value\":\"editable\"},\"receiverDeptLeaderId\":{\"title\":\"接待单位领导1\",\"value\":\"\"},\"receiverId\":{\"title\":\"接待人\",\"value\":\"editable\"},\"safeDeptAuditorId\":{\"title\":\"安保部审核人1\",\"value\":\"\"}},\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":[],\"nextStep\":{\"id\":2,\"title\":\"接待人审核\",\"category\":\"audit\",\"candidates\":[{\"category\":\"field\",\"title\":\"接待人\",\"value\":\"receiverId\"}],\"auth\":{\"applicant\":{\"title\":\"申请人\",\"value\":\"editable\"},\"receiverDeptLeaderId\":{\"title\":\"接待单位领导1\",\"value\":\"editable\"},\"receiverId\":{\"title\":\"接待人\",\"value\":\"editable\"},\"safeDeptAuditorId\":{\"title\":\"安保部审核人1\",\"value\":\"readonly\"}},\"expression\":\"\",\"auditMethod\":\"or\",\"branchSteps\":[],\"nextStep\":{\"id\":6,\"title\":\"抄送接待人\",\"category\":\"notify\",\"candidates\":[],\"auth\":{\"applicant\":{\"title\":\"申请人\",\"value\":\"editable\"},\"receiverDeptLeaderId\":{\"title\":\"接待单位领导1\",\"value\":\"\"},\"receiverId\":{\"title\":\"接待人\",\"value\":\"editable\"},\"safeDeptAuditorId\":{\"title\":\"安保部审核人1\",\"value\":\"\"}},\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":[],\"nextStep\":{\"id\":3,\"title\":\"接待部门领导审核\",\"category\":\"audit\",\"candidates\":[{\"category\":\"field\",\"title\":\"接待单位领导\",\"value\":\"receiverDeptLeaderId\"}],\"auth\":{\"applicant\":{\"title\":\"申请人\",\"value\":\"editable\"},\"receiverDeptLeaderId\":{\"title\":\"接待单位领导1\",\"value\":\"editable\"},\"receiverId\":{\"title\":\"接待人\",\"value\":\"editable\"},\"safeDeptAuditorId\":{\"title\":\"安保部审核人1\",\"value\":\"editable\"}},\"expression\":\"\",\"auditMethod\":\"or\",\"branchSteps\":[],\"nextStep\":{\"id\":4,\"title\":\"安保部审核\",\"category\":\"audit\",\"candidates\":[{\"category\":\"field\",\"title\":\"安保部审核人\",\"value\":\"safeDeptAuditorId\"}],\"auth\":{\"applicant\":{\"title\":\"申请人\",\"value\":\"editable\"},\"receiverDeptLeaderId\":{\"title\":\"接待单位领导1\",\"value\":\"editable\"},\"receiverId\":{\"title\":\"接待人\",\"value\":\"editable\"},\"safeDeptAuditorId\":{\"title\":\"安保部审核人1\",\"value\":\"editable\"}},\"expression\":\"\",\"auditMethod\":\"or\",\"branchSteps\":[],\"nextStep\":{\"id\":9999,\"title\":\"结束\",\"category\":\"end\",\"candidates\":[],\"auth\":null,\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":[],\"nextStep\":{\"id\":0,\"title\":\"\",\"category\":\"\",\"candidates\":null,\"auth\":null,\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":null,\"nextStep\":null}}}}}}}}', '[{\"title\":\"申请人\",\"name\":\"applicant\"},{\"title\":\"接待人\",\"name\":\"receiverId\"},{\"title\":\"接待单位领导\",\"name\":\"receiverDeptLeaderId\"},{\"title\":\"安保部审核人\",\"name\":\"safeDeptAuditorId\"}]', 'release', '2025-03-19 10:08:53', '2025-03-22 11:00:56', NULL);
INSERT INTO `template` VALUES (66, 4, '船舶建造动火确认单', 1, '{\"id\":1,\"title\":\"申请\",\"category\":\"start\",\"form\":null,\"candidates\":[],\"expression\":\"\",\"auditMethod\":\"or\",\"branchSteps\":[],\"nextStep\":{\"id\":3,\"title\":\"船舶建造师确认内容\",\"category\":\"audit\",\"form\":{},\"candidates\":[{\"category\":\"field\",\"title\":\"船舶建造师\",\"value\":\"shipbuilderId\"}],\"expression\":\"\",\"auditMethod\":\"or\",\"branchSteps\":[],\"nextStep\":{\"id\":2,\"title\":\"安全员确认内容\",\"category\":\"audit\",\"form\":null,\"candidates\":[{\"category\":\"field\",\"title\":\"监火人员\",\"value\":\"supervisorId\"}],\"expression\":\"\",\"auditMethod\":\"or\",\"branchSteps\":[],\"nextStep\":{\"id\":4,\"title\":\"明火施工后现场检查的确认\",\"category\":\"audit\",\"form\":{},\"candidates\":[{\"category\":\"field\",\"title\":\"施工后监火人员\",\"value\":\"postSupervisorId\"}],\"expression\":\"\",\"auditMethod\":\"or\",\"branchSteps\":[],\"nextStep\":{\"id\":9999,\"title\":\"结束\",\"category\":\"end\",\"form\":{},\"candidates\":[],\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":[],\"nextStep\":{\"id\":0,\"title\":\"\",\"category\":\"\",\"form\":null,\"candidates\":null,\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":null,\"nextStep\":null}}}}}}', '[{\"title\":\"船舶建造师\",\"name\":\"shipbuilderId\"},{\"title\":\"监火人员\",\"name\":\"supervisorId\"},{\"title\":\"施工后监火人员\",\"name\":\"postSupervisorId\"},{\"title\":\"申请人\",\"name\":\"creator\"}]', 'release', '2025-03-20 09:42:11', '2025-03-21 09:24:04', NULL);
INSERT INTO `template` VALUES (68, 3, '访客预约单', 2, '{\"id\":1,\"title\":\"访客预约申请\",\"category\":\"start\",\"candidates\":[],\"auth\":{\"applicant\":{\"title\":\"申请人\",\"value\":\"editable\"},\"receiverDeptLeaderId\":{\"title\":\"接待单位领导1\",\"value\":\"\"},\"receiverId\":{\"title\":\"接待人\",\"value\":\"editable\"},\"safeDeptAuditorId\":{\"title\":\"安保部审核人1\",\"value\":\"\"}},\"expression\":\"\",\"auditMethod\":\"or\",\"branchSteps\":[],\"nextStep\":{\"id\":5,\"title\":\"抄送预约人\",\"category\":\"notify\",\"candidates\":[],\"auth\":{\"applicant\":{\"title\":\"申请人\",\"value\":\"editable\"},\"receiverDeptLeaderId\":{\"title\":\"接待单位领导1\",\"value\":\"\"},\"receiverId\":{\"title\":\"接待人\",\"value\":\"editable\"},\"safeDeptAuditorId\":{\"title\":\"安保部审核人1\",\"value\":\"\"}},\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":[],\"nextStep\":{\"id\":2,\"title\":\"接待人审核\",\"category\":\"audit\",\"candidates\":[{\"category\":\"field\",\"title\":\"接待人\",\"value\":\"receiverId\"}],\"auth\":{\"applicant\":{\"title\":\"申请人\",\"value\":\"editable\"},\"receiverDeptLeaderId\":{\"title\":\"接待单位领导1\",\"value\":\"editable\"},\"receiverId\":{\"title\":\"接待人\",\"value\":\"editable\"},\"safeDeptAuditorId\":{\"title\":\"安保部审核人1\",\"value\":\"readonly\"}},\"expression\":\"\",\"auditMethod\":\"or\",\"branchSteps\":[],\"nextStep\":{\"id\":6,\"title\":\"抄送接待人\",\"category\":\"notify\",\"candidates\":[],\"auth\":{\"applicant\":{\"title\":\"申请人\",\"value\":\"editable\"},\"receiverDeptLeaderId\":{\"title\":\"接待单位领导1\",\"value\":\"\"},\"receiverId\":{\"title\":\"接待人\",\"value\":\"editable\"},\"safeDeptAuditorId\":{\"title\":\"安保部审核人1\",\"value\":\"\"}},\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":[],\"nextStep\":{\"id\":3,\"title\":\"接待部门领导审核\",\"category\":\"audit\",\"candidates\":[{\"category\":\"field\",\"title\":\"接待单位领导\",\"value\":\"receiverDeptLeaderId\"}],\"auth\":{\"applicant\":{\"title\":\"申请人\",\"value\":\"editable\"},\"receiverDeptLeaderId\":{\"title\":\"接待单位领导1\",\"value\":\"editable\"},\"receiverId\":{\"title\":\"接待人\",\"value\":\"editable\"},\"safeDeptAuditorId\":{\"title\":\"安保部审核人1\",\"value\":\"editable\"}},\"expression\":\"\",\"auditMethod\":\"or\",\"branchSteps\":[],\"nextStep\":{\"id\":4,\"title\":\"安保部审核\",\"category\":\"audit\",\"candidates\":[{\"category\":\"field\",\"title\":\"安保部审核人\",\"value\":\"safeDeptAuditorId\"}],\"auth\":{\"applicant\":{\"title\":\"申请人\",\"value\":\"editable\"},\"receiverDeptLeaderId\":{\"title\":\"接待单位领导1\",\"value\":\"editable\"},\"receiverId\":{\"title\":\"接待人\",\"value\":\"editable\"},\"safeDeptAuditorId\":{\"title\":\"安保部审核人1\",\"value\":\"editable\"}},\"expression\":\"\",\"auditMethod\":\"or\",\"branchSteps\":[],\"nextStep\":{\"id\":9999,\"title\":\"结束\",\"category\":\"end\",\"candidates\":[],\"auth\":null,\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":[],\"nextStep\":{\"id\":0,\"title\":\"\",\"category\":\"\",\"candidates\":null,\"auth\":null,\"expression\":\"\",\"auditMethod\":\"\",\"branchSteps\":null,\"nextStep\":null}}}}}}}}', '[{\"title\":\"申请人\",\"name\":\"applicant\"},{\"title\":\"接待人\",\"name\":\"receiverId\"},{\"title\":\"接待单位领导\",\"name\":\"receiverDeptLeaderId\"},{\"title\":\"安保部审核人\",\"name\":\"safeDeptAuditorId\"}]', 'draft', '2025-03-19 10:08:53', '2025-03-22 23:19:25', NULL);

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

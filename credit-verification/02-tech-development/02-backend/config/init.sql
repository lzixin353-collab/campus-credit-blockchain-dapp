CREATE DATABASE IF NOT EXISTS credit_dapp DEFAULT CHARSET utf8mb4;
USE credit_dapp;

-- 用户表（关联链上地址）
CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    address VARCHAR(64) NOT NULL COMMENT '以太坊地址',
    role VARCHAR(32) NOT NULL COMMENT '角色：student/teacher/admin/super_admin',
    name VARCHAR(64) COMMENT '姓名',
    school VARCHAR(128) COMMENT '学校',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY idx_address (address)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 学分表（链上数据同步）
CREATE TABLE IF NOT EXISTS credits (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    student_address VARCHAR(64) NOT NULL COMMENT '学生地址',
    teacher_address VARCHAR(64) NOT NULL COMMENT '录入教师地址',
    course_name VARCHAR(128) NOT NULL COMMENT '课程名',
    score DECIMAL(5,2) NOT NULL COMMENT '分数',
    status VARCHAR(32) NOT NULL COMMENT '状态：pending/approved/rejected',
    tx_hash VARCHAR(66) COMMENT '链上交易哈希',
    audit_admin VARCHAR(64) COMMENT '审核管理员地址',
    audit_time DATETIME COMMENT '审核时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_student (student_address),
    KEY idx_teacher (teacher_address),
    KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 1. 用户表（存储登录账号、密码、基础信息）
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '登录账号（唯一）',
  `password` varchar(100) NOT NULL COMMENT '加密密码（bcrypt）', -- 关键字段：password
  `address` varchar(64) DEFAULT NULL COMMENT '以太坊地址（关联合约角色）',
  `role` varchar(20) NOT NULL DEFAULT 'student' COMMENT '本地角色（teacher/admin/student）',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  UNIQUE KEY `idx_address` (`address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户基础信息表';

-- 2. 接口访问日志表（可选，用于调试）
CREATE TABLE `access_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned DEFAULT NULL,
  `path` varchar(100) NOT NULL,
  `method` varchar(10) NOT NULL,
  `ip` varchar(32) DEFAULT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='接口访问日志';
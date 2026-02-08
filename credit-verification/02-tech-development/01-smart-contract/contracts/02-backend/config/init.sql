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
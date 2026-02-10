
-- 初始化DB数据

-- 用户表
CREATE TABLE IF NOT EXISTS `users` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `username` VARCHAR(50) NOT NULL UNIQUE COMMENT '登录用户名',
    `password` VARCHAR(255) NOT NULL COMMENT '加密密码',
    `nickname` VARCHAR(100) DEFAULT '' COMMENT '用户昵称',
    `email` VARCHAR(100) DEFAULT '' COMMENT '邮箱',
    `phone` VARCHAR(20) DEFAULT '' COMMENT '手机号',
    `avatar` VARCHAR(255) DEFAULT '' COMMENT '头像URL',
    `status` TINYINT DEFAULT 1 COMMENT '状态: 0=禁用, 1=正常',
    `last_login_time` DATETIME NULL COMMENT '最后登录时间',
    `last_login_ip` VARCHAR(50) DEFAULT '' COMMENT '最后登录IP',
    `enterprise_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '企业ID',
    `role_ids` VARCHAR(255) DEFAULT '' COMMENT '角色ID列表，逗号分隔',
    -- 标准审计字段（所有表必须包含）
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NULL COMMENT '修改时间',
    `deleted_at` DATETIME NULL COMMENT '软删除时间',
    `create_by` BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '创建人ID',
    `create_by_name` VARCHAR(20) DEFAULT '系统' NOT NULL COMMENT '创建人姓名',
    `update_by` BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '修改人ID',
    `update_by_name` VARCHAR(20) DEFAULT '系统' NOT NULL COMMENT '修改人姓名',
    INDEX `idx_username` (`username`),
    INDEX `idx_enterprise` (`enterprise_id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_created_at` (`created_at`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

CREATE TABLE IF NOT EXISTS `token_blacklist` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `jti` VARCHAR(100) NOT NULL UNIQUE COMMENT 'Token唯一标识',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `expires_at` DATETIME NOT NULL COMMENT 'Token过期时间',
    -- 标准审计字段（所有表必须包含）
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NULL COMMENT '修改时间',
    `deleted_at` DATETIME NULL COMMENT '软删除时间',
    `create_by` BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '创建人ID',
    `create_by_name` VARCHAR(20) DEFAULT '系统' NOT NULL COMMENT '创建人姓名',
    `update_by` BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '修改人ID',
    `update_by_name` VARCHAR(20) DEFAULT '系统' NOT NULL COMMENT '修改人姓名',
    INDEX `idx_jti` (`jti`),
    INDEX `idx_expires` (`expires_at`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Token黑名单';


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

-- 模型分类表（树形结构）
CREATE TABLE IF NOT EXISTS `model_categories` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `parent_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '父分类ID，0为顶级分类',
    `name` VARCHAR(100) NOT NULL COMMENT '分类名称',
    `code` VARCHAR(50) NOT NULL UNIQUE COMMENT '分类编码',
    `description` VARCHAR(500) DEFAULT '' COMMENT '分类描述',
    `icon` VARCHAR(255) DEFAULT '' COMMENT '分类图标URL',
    `sort_order` INT DEFAULT 0 COMMENT '排序顺序',
    `level` INT DEFAULT 1 COMMENT '层级深度',
    `path` VARCHAR(500) DEFAULT '' COMMENT '分类路径，如：/1/2/3/',
    `status` TINYINT DEFAULT 1 COMMENT '状态: 0=禁用, 1=启用',
    -- 标准审计字段（所有表必须包含）
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NULL COMMENT '修改时间',
    `deleted_at` DATETIME NULL COMMENT '软删除时间',
    `create_by` BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '创建人ID',
    `create_by_name` VARCHAR(20) DEFAULT '系统' NOT NULL COMMENT '创建人姓名',
    `update_by` BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '修改人ID',
    `update_by_name` VARCHAR(20) DEFAULT '系统' NOT NULL COMMENT '修改人姓名',
    INDEX `idx_parent_id` (`parent_id`),
    INDEX `idx_code` (`code`),
    INDEX `idx_status` (`status`),
    INDEX `idx_sort_order` (`sort_order`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='思维模型分类表';

-- 思维模型表
CREATE TABLE IF NOT EXISTS `thinking_models` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `name` VARCHAR(100) NOT NULL COMMENT '模型名称',
    `code` VARCHAR(50) NOT NULL UNIQUE COMMENT '模型编码',
    `description` VARCHAR(2000) DEFAULT '' COMMENT '模型描述',
    `category_id` BIGINT UNSIGNED NOT NULL COMMENT '所属分类ID',
    `price` DECIMAL(10, 2) DEFAULT 0.00 COMMENT '价格，0表示免费',
    `content` JSON NOT NULL COMMENT '模型内容JSON（字段定义、结构等）',
    `cover_image` VARCHAR(255) DEFAULT '' COMMENT '封面图片URL',
    `icon` VARCHAR(255) DEFAULT '' COMMENT '模型图标',
    `difficulty` TINYINT DEFAULT 1 COMMENT '难度: 1=入门, 2=进阶, 3=高级',
    `estimated_time` INT DEFAULT 30 COMMENT '预计完成时间（分钟）',
    `usage_count` INT DEFAULT 0 COMMENT '使用次数统计',
    `status` TINYINT DEFAULT 0 COMMENT '状态: 0=草稿, 1=已发布, 2=已下架',
    `publish_time` DATETIME NULL COMMENT '发布时间',
    `version` VARCHAR(20) DEFAULT '1.0.0' COMMENT '版本号',
    `author_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '作者ID',
    `author_name` VARCHAR(100) DEFAULT '' COMMENT '作者名称',
    -- 标准审计字段（所有表必须包含）
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NULL COMMENT '修改时间',
    `deleted_at` DATETIME NULL COMMENT '软删除时间',
    `create_by` BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '创建人ID',
    `create_by_name` VARCHAR(20) DEFAULT '系统' NOT NULL COMMENT '创建人姓名',
    `update_by` BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '修改人ID',
    `update_by_name` VARCHAR(20) DEFAULT '系统' NOT NULL COMMENT '修改人姓名',
    INDEX `idx_category_id` (`category_id`),
    INDEX `idx_code` (`code`),
    INDEX `idx_status` (`status`),
    INDEX `idx_difficulty` (`difficulty`),
    INDEX `idx_price` (`price`),
    INDEX `idx_author_id` (`author_id`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='思维模型表';

-- 课题表
CREATE TABLE IF NOT EXISTS `topics` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `title` VARCHAR(200) NOT NULL COMMENT '课题标题',
    `description` VARCHAR(2000) DEFAULT '' COMMENT '课题描述',
    `status` TINYINT DEFAULT 0 COMMENT '状态: 0=进行中, 1=已完成, 2=已归档',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '所属用户ID',
    `model_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '当前使用的思维模型ID',
    `priority` TINYINT DEFAULT 1 COMMENT '优先级: 1=低, 2=中, 3=高',
    `tags` VARCHAR(500) DEFAULT '' COMMENT '标签，逗号分隔',
    `deadline` DATETIME NULL COMMENT '截止日期',
    `complete_time` DATETIME NULL COMMENT '完成时间',
    -- 标准审计字段（所有表必须包含）
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NULL COMMENT '修改时间',
    `deleted_at` DATETIME NULL COMMENT '软删除时间',
    `create_by` BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '创建人ID',
    `create_by_name` VARCHAR(20) DEFAULT '系统' NOT NULL COMMENT '创建人姓名',
    `update_by` BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '修改人ID',
    `update_by_name` VARCHAR(20) DEFAULT '系统' NOT NULL COMMENT '修改人姓名',
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_model_id` (`model_id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_priority` (`priority`),
    INDEX `idx_deadline` (`deadline`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课题表';

-- 课题分析记录表
CREATE TABLE IF NOT EXISTS `topic_analyses` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `topic_id` BIGINT UNSIGNED NOT NULL COMMENT '所属课题ID',
    `model_id` BIGINT UNSIGNED NOT NULL COMMENT '使用的思维模型ID',
    `model_name` VARCHAR(100) DEFAULT '' COMMENT '模型名称（快照）',
    `content` JSON NOT NULL COMMENT '用户填写的分析内容JSON',
    `ai_analysis` TEXT COMMENT 'AI分析结果',
    `ai_suggestions` TEXT COMMENT 'AI建议',
    `version` INT DEFAULT 1 COMMENT '版本号，同一课题同一模型的多次分析',
    `is_current` TINYINT DEFAULT 1 COMMENT '是否为当前版本: 0=否, 1=是',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    -- 标准审计字段（所有表必须包含）
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NULL COMMENT '修改时间',
    `deleted_at` DATETIME NULL COMMENT '软删除时间',
    `create_by` BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '创建人ID',
    `create_by_name` VARCHAR(20) DEFAULT '系统' NOT NULL COMMENT '创建人姓名',
    `update_by` BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '修改人ID',
    `update_by_name` VARCHAR(20) DEFAULT '系统' NOT NULL COMMENT '修改人姓名',
    INDEX `idx_topic_id` (`topic_id`),
    INDEX `idx_model_id` (`model_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_is_current` (`is_current`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课题分析记录表';

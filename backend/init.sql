
-- 初始化DB数据
CREATE TABLE `domain` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'Primary Key',
  `name` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `create_by` int unsigned NOT NULL DEFAULT '0' COMMENT '创建人id',
  `update_by` int unsigned NOT NULL DEFAULT '0' COMMENT '修改人id',
  `create_by_name` varchar(20) NOT NULL DEFAULT '' COMMENT '创建人',
  `update_by_name` varchar(20) NOT NULL DEFAULT '' COMMENT '更新人',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=223 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='领域表'

CREATE TABLE `model` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'Primary Key',
  `name` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `create_by` int unsigned NOT NULL DEFAULT '0' COMMENT '创建人id',
  `update_by` int unsigned NOT NULL DEFAULT '0' COMMENT '修改人id',
  `create_by_name` varchar(20) NOT NULL DEFAULT '' COMMENT '创建人',
  `update_by_name` varchar(20) NOT NULL DEFAULT '' COMMENT '更新人',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='模型表'

CREATE TABLE `api` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'Primary Key',
  `name` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `create_by` int unsigned NOT NULL DEFAULT '0' COMMENT '创建人id',
  `update_by` int unsigned NOT NULL DEFAULT '0' COMMENT '修改人id',
  `create_by_name` varchar(20) NOT NULL DEFAULT '' COMMENT '创建人',
  `update_by_name` varchar(20) NOT NULL DEFAULT '' COMMENT '更新人',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='接口表'
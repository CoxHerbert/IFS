-- Personal ledger accounts and optional cloud sync.
-- Apply this migration before deploying the ledger API.
CREATE TABLE IF NOT EXISTS `ledger_account` (
  `account_id` BIGINT NOT NULL AUTO_INCREMENT,
  `login_id` VARCHAR(120) NOT NULL,
  `nickname` VARCHAR(120) NOT NULL DEFAULT '',
  `avatar_key` VARCHAR(32) NOT NULL DEFAULT '',
  `inviter_uid` VARCHAR(32) NOT NULL DEFAULT '',
  `avatar_frame_key` VARCHAR(32) NOT NULL DEFAULT '',
  `password_hash` VARCHAR(100) NOT NULL,
  `cloud_sync_enabled` TINYINT(1) NOT NULL DEFAULT 1,
  `consent_updated_at` DATETIME NULL,
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`account_id`),
  UNIQUE KEY `uk_ledger_account_login_id` (`login_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `ledger_record` (
  `account_id` BIGINT NOT NULL,
  `record_id` VARCHAR(128) NOT NULL,
  `record_date` VARCHAR(40) NOT NULL,
  `merchant` VARCHAR(255) NOT NULL,
  `category` VARCHAR(100) NOT NULL,
  `amount` DECIMAL(18,2) NOT NULL,
  `record_type` VARCHAR(16) NOT NULL,
  `source` VARCHAR(16) NOT NULL,
  `import_batch_id` VARCHAR(64) NULL,
  `import_file_name` VARCHAR(255) NULL,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`account_id`, `record_id`),
  KEY `idx_ledger_record_date` (`account_id`, `record_date`),
  KEY `idx_ledger_record_batch` (`account_id`, `import_batch_id`),
  CONSTRAINT `fk_ledger_record_account` FOREIGN KEY (`account_id`) REFERENCES `ledger_account` (`account_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `ledger_import_batch` (
  `account_id` BIGINT NOT NULL,
  `batch_id` VARCHAR(64) NOT NULL,
  `source` VARCHAR(16) NOT NULL,
  `file_name` VARCHAR(255) NOT NULL,
  `imported_at` VARCHAR(40) NOT NULL,
  `imported_count` INT UNSIGNED NOT NULL DEFAULT 0,
  `duplicate_count` INT UNSIGNED NOT NULL DEFAULT 0,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`account_id`, `batch_id`),
  KEY `idx_ledger_import_batch_time` (`account_id`, `imported_at`),
  CONSTRAINT `fk_ledger_import_batch_account` FOREIGN KEY (`account_id`) REFERENCES `ledger_account` (`account_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

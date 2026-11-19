-- Apply once to existing databases to enable cloud persistence for import batches.
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

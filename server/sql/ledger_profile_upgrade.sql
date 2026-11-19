-- Safe to run repeatedly against an existing ledger database.
-- Adds profile columns only when each column is missing.

SET @ifs_ledger_has_nickname = (
  SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'ledger_account' AND COLUMN_NAME = 'nickname'
);
SET @ifs_ledger_sql = IF(
  @ifs_ledger_has_nickname = 0,
  'ALTER TABLE `ledger_account` ADD COLUMN `nickname` VARCHAR(120) NOT NULL DEFAULT '''' AFTER `login_id`',
  'SELECT ''nickname already exists'' AS migration_status'
);
PREPARE ifs_ledger_migration_stmt FROM @ifs_ledger_sql;
EXECUTE ifs_ledger_migration_stmt;
DEALLOCATE PREPARE ifs_ledger_migration_stmt;

SET @ifs_ledger_has_avatar_key = (
  SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'ledger_account' AND COLUMN_NAME = 'avatar_key'
);
SET @ifs_ledger_sql = IF(
  @ifs_ledger_has_avatar_key = 0,
  'ALTER TABLE `ledger_account` ADD COLUMN `avatar_key` VARCHAR(32) NOT NULL DEFAULT '''' AFTER `nickname`',
  'SELECT ''avatar_key already exists'' AS migration_status'
);
PREPARE ifs_ledger_migration_stmt FROM @ifs_ledger_sql;
EXECUTE ifs_ledger_migration_stmt;
DEALLOCATE PREPARE ifs_ledger_migration_stmt;

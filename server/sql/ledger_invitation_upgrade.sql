-- Safe to run repeatedly after ledger_profile_upgrade.sql.
-- Adds invitation and avatar-frame columns only when they are missing.

SET @ifs_ledger_has_inviter_uid = (
  SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'ledger_account' AND COLUMN_NAME = 'inviter_uid'
);
SET @ifs_ledger_sql = IF(
  @ifs_ledger_has_inviter_uid = 0,
  'ALTER TABLE `ledger_account` ADD COLUMN `inviter_uid` VARCHAR(32) NOT NULL DEFAULT '''' AFTER `avatar_key`',
  'SELECT ''inviter_uid already exists'' AS migration_status'
);
PREPARE ifs_ledger_migration_stmt FROM @ifs_ledger_sql;
EXECUTE ifs_ledger_migration_stmt;
DEALLOCATE PREPARE ifs_ledger_migration_stmt;

SET @ifs_ledger_has_avatar_frame_key = (
  SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'ledger_account' AND COLUMN_NAME = 'avatar_frame_key'
);
SET @ifs_ledger_sql = IF(
  @ifs_ledger_has_avatar_frame_key = 0,
  'ALTER TABLE `ledger_account` ADD COLUMN `avatar_frame_key` VARCHAR(32) NOT NULL DEFAULT '''' AFTER `inviter_uid`',
  'SELECT ''avatar_frame_key already exists'' AS migration_status'
);
PREPARE ifs_ledger_migration_stmt FROM @ifs_ledger_sql;
EXECUTE ifs_ledger_migration_stmt;
DEALLOCATE PREPARE ifs_ledger_migration_stmt;

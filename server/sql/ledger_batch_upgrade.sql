-- Run once on databases where ledger_record was created before import batch fields.
ALTER TABLE `ledger_record`
  ADD COLUMN `import_batch_id` VARCHAR(64) NULL AFTER `source`,
  ADD COLUMN `import_file_name` VARCHAR(255) NULL AFTER `import_batch_id`,
  ADD KEY `idx_ledger_record_batch` (`account_id`, `import_batch_id`);

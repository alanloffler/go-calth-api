DROP INDEX IF EXISTS idx_role_value;

ALTER TABLE roles
    DROP COLUMN IF EXISTS value,
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS deleted_at;

-- Remove two_fa_enabled column from users table
ALTER TABLE users DROP COLUMN IF EXISTS two_fa_enabled;

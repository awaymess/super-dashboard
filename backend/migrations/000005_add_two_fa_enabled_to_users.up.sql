-- Add two_fa_enabled column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS two_fa_enabled BOOLEAN DEFAULT FALSE;

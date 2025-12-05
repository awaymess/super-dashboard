-- Create two_factor_auths table for 2FA TOTP secrets
CREATE TABLE IF NOT EXISTS two_factor_auths (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    secret VARCHAR(64) NOT NULL,
    backup_codes TEXT, -- JSON array of backup codes
    verified BOOLEAN DEFAULT FALSE,
    enabled_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for user lookup
CREATE INDEX IF NOT EXISTS idx_two_factor_auths_user_id ON two_factor_auths(user_id);

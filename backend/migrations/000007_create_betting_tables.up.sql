-- Create bets table
CREATE TABLE IF NOT EXISTS bets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    market VARCHAR(50) NOT NULL,
    selection VARCHAR(50) NOT NULL,
    odds DECIMAL(10, 2) NOT NULL,
    stake DECIMAL(20, 2) NOT NULL,
    potential_return DECIMAL(20, 2) NOT NULL,
    bookmaker VARCHAR(100),
    status VARCHAR(20) DEFAULT 'pending',
    result VARCHAR(20),
    profit DECIMAL(20, 2),
    closing_odds DECIMAL(10, 2),
    value_percent DECIMAL(10, 2),
    settled_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_bets_user_id ON bets(user_id);
CREATE INDEX IF NOT EXISTS idx_bets_match_id ON bets(match_id);
CREATE INDEX IF NOT EXISTS idx_bets_status ON bets(status);

-- Create bankroll_history table
CREATE TABLE IF NOT EXISTS bankroll_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    balance DECIMAL(20, 2) NOT NULL,
    change DECIMAL(20, 2),
    reason VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_bankroll_history_user_id ON bankroll_history(user_id);
CREATE INDEX IF NOT EXISTS idx_bankroll_history_created_at ON bankroll_history(created_at DESC);

-- Create value_bets table
CREATE TABLE IF NOT EXISTS value_bets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    market VARCHAR(50) NOT NULL,
    selection VARCHAR(50) NOT NULL,
    bookmaker VARCHAR(100) NOT NULL,
    bookmaker_odds DECIMAL(10, 2) NOT NULL,
    true_probability DECIMAL(10, 6) NOT NULL,
    implied_probability DECIMAL(10, 6) NOT NULL,
    value_percent DECIMAL(10, 2) NOT NULL,
    kelly_stake DECIMAL(10, 2),
    confidence DECIMAL(10, 6),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_value_bets_match_id ON value_bets(match_id);
CREATE INDEX IF NOT EXISTS idx_value_bets_created_at ON value_bets(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_value_bets_value_percent ON value_bets(value_percent DESC);

-- Create stock_news table
CREATE TABLE IF NOT EXISTS stock_news (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stock_id UUID REFERENCES stocks(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    content TEXT,
    source VARCHAR(255),
    url TEXT,
    sentiment DECIMAL(3, 2) DEFAULT 0,
    published_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_stock_news_stock_id ON stock_news(stock_id);
CREATE INDEX IF NOT EXISTS idx_stock_news_published_at ON stock_news(published_at DESC);
CREATE INDEX IF NOT EXISTS idx_stock_news_sentiment ON stock_news(sentiment);

-- Create fair_values table
CREATE TABLE IF NOT EXISTS fair_values (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stock_id UUID NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
    dcf_value DECIMAL(20, 2),
    pe_value DECIMAL(20, 2),
    pbv_value DECIMAL(20, 2),
    graham_value DECIMAL(20, 2),
    buffett_value DECIMAL(20, 2),
    weighted_avg DECIMAL(20, 2) NOT NULL,
    current_price DECIMAL(20, 2) NOT NULL,
    margin_of_safety DECIMAL(10, 2),
    upside_percent DECIMAL(10, 2),
    recommendation VARCHAR(50),
    calculated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_fair_values_stock_id ON fair_values(stock_id);
CREATE INDEX IF NOT EXISTS idx_fair_values_calculated_at ON fair_values(calculated_at DESC);

-- Create trade_journal table
CREATE TABLE IF NOT EXISTS trade_journal (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    trade_id UUID REFERENCES trades(id) ON DELETE CASCADE,
    bet_id UUID REFERENCES bets(id) ON DELETE CASCADE,
    entry_reason TEXT,
    exit_reason TEXT,
    emotions VARCHAR(255),
    lessons_learned TEXT,
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    tags VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_trade_journal_user_id ON trade_journal(user_id);
CREATE INDEX IF NOT EXISTS idx_trade_journal_trade_id ON trade_journal(trade_id);
CREATE INDEX IF NOT EXISTS idx_trade_journal_bet_id ON trade_journal(bet_id);

-- Create goals table
CREATE TABLE IF NOT EXISTS goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    target_amount DECIMAL(20, 2) NOT NULL,
    current_amount DECIMAL(20, 2) DEFAULT 0,
    target_date TIMESTAMP,
    category VARCHAR(50),
    status VARCHAR(20) DEFAULT 'active',
    achieved_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_goals_user_id ON goals(user_id);
CREATE INDEX IF NOT EXISTS idx_goals_status ON goals(status);

-- Create settings table
CREATE TABLE IF NOT EXISTS settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    initial_bankroll DECIMAL(20, 2) DEFAULT 1000,
    current_bankroll DECIMAL(20, 2) DEFAULT 1000,
    kelly_factor DECIMAL(3, 2) DEFAULT 0.5,
    risk_level VARCHAR(20) DEFAULT 'moderate',
    default_bookmaker VARCHAR(100),
    value_bet_threshold DECIMAL(10, 2) DEFAULT 5,
    max_daily_bets INTEGER DEFAULT 10,
    max_stake_per_bet DECIMAL(20, 2),
    preferred_leagues TEXT,
    notify_email BOOLEAN DEFAULT TRUE,
    notify_telegram BOOLEAN DEFAULT FALSE,
    notify_line BOOLEAN DEFAULT FALSE,
    notify_discord BOOLEAN DEFAULT FALSE,
    telegram_chat_id VARCHAR(255),
    line_token VARCHAR(255),
    discord_webhook TEXT,
    theme VARCHAR(20) DEFAULT 'dark',
    language VARCHAR(10) DEFAULT 'en',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_settings_user_id ON settings(user_id);

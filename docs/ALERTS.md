# Alert System - Quick Reference Guide

## Overview

The Super Dashboard alert system monitors various conditions and sends notifications through multiple channels.

## Alert Types

### 1. Stock Price Alert
Monitor stock price movements.

```go
alert := &model.Alert{
    Type:        model.AlertTypeStockPrice,
    Symbol:      "AAPL",
    Condition:   model.AlertConditionAbove,
    TargetValue: 150.00,
}
```

**Use Cases:**
- Alert me when AAPL goes above $150
- Alert me when TSLA goes below $200
- Alert me when GOOGL crosses $140

---

### 2. Stock Volume Alert
Detect unusual trading volume.

```go
alert := &model.Alert{
    Type:        model.AlertTypeStockVolume,
    Symbol:      "AAPL",
    Condition:   model.AlertConditionAbove,
    TargetValue: 100000000, // 100M shares
}
```

**Use Cases:**
- Alert when volume exceeds 2x average
- Detect accumulation/distribution
- Spot breakout opportunities

---

### 3. Odds Change Alert
Monitor betting odds movements.

```go
alert := &model.Alert{
    Type:        model.AlertTypeOddsChange,
    Symbol:      "match_123:1X2:home", // match:market:outcome
    Condition:   model.AlertConditionPercentDown,
    TargetValue: 10.0, // 10% drop
}
```

**Use Cases:**
- Sharp money detection
- Line movement tracking
- Steam move alerts

---

### 4. Value Bet Alert
Automatic value bet notifications.

```go
alert := &model.Alert{
    Type:        model.AlertTypeValueBet,
    Symbol:      "match_123",
    Condition:   model.AlertConditionAbove,
    TargetValue: 10.0, // 10% value
}
```

**Use Cases:**
- Get notified of 10%+ value bets
- Track value in specific leagues
- Monitor bookmaker inefficiencies

---

### 5. Technical Alert
Technical indicator-based alerts.

```go
alert := &model.Alert{
    Type:        model.AlertTypeTechnical,
    Symbol:      "AAPL:RSI", // symbol:indicator
    Condition:   model.AlertConditionBelow,
    TargetValue: 30.0, // Oversold
}
```

**Use Cases:**
- RSI oversold/overbought
- MACD crossovers
- Moving average crosses
- Bollinger Band breakouts

---

### 6. News Alert
Sentiment-based news alerts.

```go
alert := &model.Alert{
    Type:        model.AlertTypeNews,
    Symbol:      "AAPL",
    Condition:   model.AlertConditionBelow,
    TargetValue: -0.5, // Negative sentiment
}
```

**Use Cases:**
- Detect negative news early
- Track sentiment changes
- Monitor company announcements

---

### 7. Dividend Alert
Dividend-related notifications.

```go
alert := &model.Alert{
    Type:       model.AlertTypeDividend,
    Symbol:     "AAPL",
    Condition:  model.AlertConditionEquals,
    Message:    "Dividend announcement",
}
```

**Use Cases:**
- Ex-dividend date reminders
- Dividend increase announcements
- Payout schedule tracking

---

### 8. Earnings Alert
Earnings report notifications.

```go
alert := &model.Alert{
    Type:       model.AlertTypeEarnings,
    Symbol:     "AAPL",
    Condition:  model.AlertConditionEquals,
    Message:    "Earnings report today",
}
```

**Use Cases:**
- Earnings calendar reminders
- Earnings surprise alerts
- Pre/post earnings price action

---

### 9. Match Start Alert
Match starting soon notifications.

```go
alert := &model.Alert{
    Type:        model.AlertTypeMatchStart,
    Symbol:      "match_123",
    TargetValue: 30, // 30 minutes before
}
```

**Use Cases:**
- Last-minute bet reminders
- Live betting preparation
- Match tracking

---

## Alert Conditions

### Above
Value exceeds target.

```go
Condition:   model.AlertConditionAbove,
TargetValue: 150.00,
// Triggers when current value > 150.00
```

**Examples:**
- Price > $150
- Volume > 100M
- RSI > 70

---

### Below
Value falls below target.

```go
Condition:   model.AlertConditionBelow,
TargetValue: 100.00,
// Triggers when current value < 100.00
```

**Examples:**
- Price < $100
- RSI < 30
- Odds < 2.00

---

### Equals
Value equals target (with small tolerance).

```go
Condition:   model.AlertConditionEquals,
TargetValue: 100.00,
// Triggers when |current - target| < 0.01
```

**Examples:**
- Price reaches exactly $100
- Event-based alerts
- Milestone tracking

---

### Percent Up
Percentage increase from last value.

```go
Condition:   model.AlertConditionPercentUp,
TargetValue: 5.0, // 5% increase
// Triggers when ((current - previous) / previous) * 100 >= 5.0
```

**Examples:**
- Stock up 5% today
- Volume spike 50%
- Rapid price movement

---

### Percent Down
Percentage decrease from last value.

```go
Condition:   model.AlertConditionPercentDown,
TargetValue: 5.0, // 5% decrease
// Triggers when ((previous - current) / previous) * 100 >= 5.0
```

**Examples:**
- Stock down 5% today
- Stop loss triggers
- Drawdown alerts

---

### Crosses
Value crosses target line (up or down).

```go
Condition:   model.AlertConditionCrosses,
TargetValue: 100.00,
// Triggers when previous < 100 && current >= 100
//         OR previous > 100 && current <= 100
```

**Examples:**
- MA crossovers
- Support/resistance breaks
- Threshold crossing

---

## Notification Channels

### In-App Notifications
Always enabled. Appears in notification center.

```go
// Automatic for all alerts
```

**Features:**
- Real-time updates
- Notification badge
- Click to view details
- Mark as read/unread

---

### Email Notifications
Send email when alert triggers.

```go
alert.NotifyEmail = true
```

**Email Format:**
```
Subject: Alert Triggered: AAPL

Hi John,

Your alert "AAPL above $150" has been triggered.

Current price: $152.45 (+1.63%)
Time: 2025-12-05 14:30:00

View in Dashboard: https://dashboard.app/alerts/123

Best,
Super Dashboard Team
```

---

### Telegram Notifications
Send Telegram message via bot.

```go
alert.NotifyTelegram = true
// Requires user to set telegram_chat_id in settings
```

**Message Format:**
```
🔔 Alert Triggered

📊 AAPL
💰 $152.45 (+1.63%)
🎯 Target: $150.00

Condition: Above
Time: 14:30:00

View: https://t.me/superdashboard_bot
```

---

### LINE Notifications
Send LINE Notify message.

```go
alert.NotifyLINE = true
// Requires user to set line_token in settings
```

---

### Discord Notifications
Send Discord webhook message.

```go
alert.NotifyDiscord = true
// Requires user to set discord_webhook in settings
```

**Webhook Embed:**
```json
{
  "embeds": [{
    "title": "Alert Triggered: AAPL",
    "description": "Price above $150.00",
    "color": 3066993,
    "fields": [
      {"name": "Current Price", "value": "$152.45", "inline": true},
      {"name": "Change", "value": "+1.63%", "inline": true},
      {"name": "Target", "value": "$150.00", "inline": true}
    ],
    "timestamp": "2025-12-05T14:30:00.000Z"
  }]
}
```

---

## API Examples

### Create Alert

**Endpoint:** `POST /api/v1/alerts`

**Request:**
```json
{
  "type": "stock_price",
  "symbol": "AAPL",
  "condition": "above",
  "target_value": 150.00,
  "message": "AAPL reached my target!",
  "notify_email": true,
  "notify_telegram": true,
  "active": true
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "type": "stock_price",
    "symbol": "AAPL",
    "condition": "above",
    "target_value": 150.00,
    "current_value": 145.23,
    "active": true,
    "created_at": "2025-12-05T14:00:00Z"
  }
}
```

---

### List User Alerts

**Endpoint:** `GET /api/v1/alerts`

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "type": "stock_price",
      "symbol": "AAPL",
      "condition": "above",
      "target_value": 150.00,
      "current_value": 145.23,
      "active": true,
      "last_triggered": null,
      "trigger_count": 0
    },
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "type": "stock_volume",
      "symbol": "TSLA",
      "condition": "above",
      "target_value": 100000000,
      "current_value": 85000000,
      "active": true,
      "last_triggered": "2025-12-04T10:30:00Z",
      "trigger_count": 3
    }
  ],
  "meta": {
    "total": 2,
    "active": 2
  }
}
```

---

### Update Alert

**Endpoint:** `PATCH /api/v1/alerts/:id`

**Request:**
```json
{
  "target_value": 155.00,
  "active": true
}
```

---

### Delete Alert

**Endpoint:** `DELETE /api/v1/alerts/:id`

**Response:**
```json
{
  "success": true,
  "message": "Alert deleted successfully"
}
```

---

### Get Alert History

**Endpoint:** `GET /api/v1/alerts/:id/history`

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "triggered_at": "2025-12-05T14:30:00Z",
      "value": 152.45,
      "notified_via": ["email", "telegram"]
    },
    {
      "triggered_at": "2025-12-04T10:15:00Z",
      "value": 151.20,
      "notified_via": ["email"]
    }
  ]
}
```

---

## Best Practices

### 1. Set Realistic Targets
Don't set alerts that will trigger constantly.

❌ Bad: Alert when price changes by 0.1%
✅ Good: Alert when price changes by 5%

---

### 2. Use Multiple Conditions
Combine alerts for better signals.

**Example:**
- Alert 1: RSI < 30 (oversold)
- Alert 2: Price > 20 MA (trend confirmation)
- Alert 3: Volume > 2x average (unusual activity)

---

### 3. Avoid Alert Fatigue
Too many alerts = you'll ignore them all.

**Recommended:**
- Max 10-15 active alerts per user
- Use percentage-based conditions for volatile assets
- Set minimum trigger intervals

---

### 4. Test Alerts
Create test alerts with easy-to-trigger conditions first.

```go
// Test alert
alert := &model.Alert{
    Symbol:      "AAPL",
    Condition:   model.AlertConditionAbove,
    TargetValue: 1.00, // Will trigger immediately
    NotifyEmail: false, // Don't spam yourself
}
```

---

### 5. Use Descriptive Messages
Custom messages help you remember why you set the alert.

❌ Bad: `message: ""`
✅ Good: `message: "Entry point for long position based on support level"`

---

### 6. Disable Alerts After Trigger
For one-time events, disable after first trigger.

```go
// In frontend, after receiving notification:
PATCH /api/v1/alerts/:id
{ "active": false }
```

---

## Troubleshooting

### Alert Not Triggering

**Check:**
1. Is alert active? (`active = true`)
2. Is current value actually meeting condition?
3. Check alert_checker worker logs
4. Verify data is being updated (stock_sync, odds_sync)

**Debug:**
```bash
# Check alert in database
psql -d super_dashboard -c "SELECT * FROM alerts WHERE id = 'alert-id';"

# Check current stock price
psql -d super_dashboard -c "SELECT * FROM stock_prices WHERE stock_id IN (SELECT id FROM stocks WHERE symbol = 'AAPL') ORDER BY timestamp DESC LIMIT 1;"

# Check worker logs
tail -f /var/log/super-dashboard/workers.log | grep alert_checker
```

---

### Not Receiving Notifications

**Check:**
1. Notification channels enabled? (`notify_email = true`)
2. User settings configured? (email, telegram_chat_id, etc.)
3. Check notification service logs
4. Verify notification was created in database

**Debug:**
```bash
# Check notification was created
psql -d super_dashboard -c "SELECT * FROM notifications WHERE user_id = 'user-id' ORDER BY created_at DESC LIMIT 5;"

# Check user notification settings
psql -d super_dashboard -c "SELECT notify_email, notify_telegram, telegram_chat_id FROM settings WHERE user_id = 'user-id';"
```

---

### Too Many Alerts

**Solutions:**
1. Increase target value threshold
2. Add minimum trigger interval
3. Use `crosses` condition instead of `above`/`below`
4. Temporarily disable alerts

---

## Performance Notes

- Alert checking runs every 30 seconds
- Typical processing time: < 1 second for 1000 alerts
- Notifications are sent asynchronously (don't block worker)
- Database indexes optimize alert queries
- Only active alerts are checked

---

## Security Considerations

- Users can only view/edit their own alerts
- API endpoints require authentication
- Notification tokens encrypted in database
- Rate limiting on alert creation (max 50 per user)
- WebSocket events only sent to alert owner

---

**Last Updated:** December 5, 2025
**Version:** 1.0

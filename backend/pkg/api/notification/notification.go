package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// EmailProvider defines email sending interface.
type EmailProvider interface {
	SendEmail(ctx context.Context, to []string, subject, body string) error
}

// SendGridClient implements SendGrid email provider.
type SendGridClient struct {
	apiKey     string
	fromEmail  string
	fromName   string
	httpClient *http.Client
}

// NewSendGridClient creates a new SendGrid client.
func NewSendGridClient(apiKey, fromEmail, fromName string) *SendGridClient {
	return &SendGridClient{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		fromName:  fromName,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendEmail sends an email via SendGrid.
func (c *SendGridClient) SendEmail(ctx context.Context, to []string, subject, body string) error {
	payload := map[string]interface{}{
		"personalizations": []map[string]interface{}{
			{
				"to": func() []map[string]string {
					recipients := make([]map[string]string, len(to))
					for i, email := range to {
						recipients[i] = map[string]string{"email": email}
					}
					return recipients
				}(),
			},
		},
		"from": map[string]string{
			"email": c.fromEmail,
			"name":  c.fromName,
		},
		"subject": subject,
		"content": []map[string]string{
			{
				"type":  "text/html",
				"value": body,
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("sendgrid error: status %d", resp.StatusCode)
	}

	return nil
}

// TelegramClient implements Telegram Bot API.
type TelegramClient struct {
	botToken   string
	httpClient *http.Client
}

// NewTelegramClient creates a new Telegram client.
func NewTelegramClient(botToken string) *TelegramClient {
	return &TelegramClient{
		botToken: botToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendMessage sends a message to a Telegram chat.
func (c *TelegramClient) SendMessage(ctx context.Context, chatID string, message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.botToken)

	payload := map[string]interface{}{
		"chat_id":    chatID,
		"text":       message,
		"parse_mode": "HTML",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("telegram error: status %d", resp.StatusCode)
	}

	return nil
}

// SendPhoto sends a photo to a Telegram chat.
func (c *TelegramClient) SendPhoto(ctx context.Context, chatID, photoURL, caption string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", c.botToken)

	payload := map[string]interface{}{
		"chat_id": chatID,
		"photo":   photoURL,
		"caption": caption,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("telegram error: status %d", resp.StatusCode)
	}

	return nil
}

// LINEClient implements LINE Messaging API.
type LINEClient struct {
	channelToken string
	httpClient   *http.Client
}

// NewLINEClient creates a new LINE Messaging client.
func NewLINEClient(channelToken string) *LINEClient {
	return &LINEClient{
		channelToken: channelToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// PushMessage sends a push message to a LINE user.
func (c *LINEClient) PushMessage(ctx context.Context, userID string, message string) error {
	url := "https://api.line.me/v2/bot/message/push"

	payload := map[string]interface{}{
		"to": userID,
		"messages": []map[string]string{
			{
				"type": "text",
				"text": message,
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.channelToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("LINE error: status %d", resp.StatusCode)
	}

	return nil
}

// PushFlexMessage sends a Flex Message (rich message) to LINE.
func (c *LINEClient) PushFlexMessage(ctx context.Context, userID string, altText string, flexContent map[string]interface{}) error {
	url := "https://api.line.me/v2/bot/message/push"

	payload := map[string]interface{}{
		"to": userID,
		"messages": []map[string]interface{}{
			{
				"type":     "flex",
				"altText":  altText,
				"contents": flexContent,
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.channelToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("LINE error: status %d", resp.StatusCode)
	}

	return nil
}

// DiscordClient implements Discord Webhook.
type DiscordClient struct {
	webhookURL string
	httpClient *http.Client
}

// NewDiscordClient creates a new Discord webhook client.
func NewDiscordClient(webhookURL string) *DiscordClient {
	return &DiscordClient{
		webhookURL: webhookURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendMessage sends a message to Discord channel.
func (c *DiscordClient) SendMessage(ctx context.Context, content string) error {
	payload := map[string]interface{}{
		"content": content,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("discord error: status %d", resp.StatusCode)
	}

	return nil
}

// SendEmbed sends a rich embed message to Discord.
func (c *DiscordClient) SendEmbed(ctx context.Context, embed DiscordEmbed) error {
	payload := map[string]interface{}{
		"embeds": []DiscordEmbed{embed},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("discord error: status %d", resp.StatusCode)
	}

	return nil
}

// DiscordEmbed represents a Discord embed message.
type DiscordEmbed struct {
	Title       string              `json:"title,omitempty"`
	Description string              `json:"description,omitempty"`
	URL         string              `json:"url,omitempty"`
	Color       int                 `json:"color,omitempty"`
	Fields      []DiscordEmbedField `json:"fields,omitempty"`
	Thumbnail   *DiscordEmbedImage  `json:"thumbnail,omitempty"`
	Image       *DiscordEmbedImage  `json:"image,omitempty"`
	Footer      *DiscordEmbedFooter `json:"footer,omitempty"`
	Timestamp   string              `json:"timestamp,omitempty"`
}

// DiscordEmbedField represents a field in embed.
type DiscordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// DiscordEmbedImage represents an image in embed.
type DiscordEmbedImage struct {
	URL string `json:"url"`
}

// DiscordEmbedFooter represents footer in embed.
type DiscordEmbedFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url,omitempty"`
}

// NotificationManager manages all notification channels.
type NotificationManager struct {
	email    EmailProvider
	telegram *TelegramClient
	line     *LINEClient
	discord  *DiscordClient
}

// NewNotificationManager creates a new notification manager.
func NewNotificationManager(email EmailProvider, telegram *TelegramClient, line *LINEClient, discord *DiscordClient) *NotificationManager {
	return &NotificationManager{
		email:    email,
		telegram: telegram,
		line:     line,
		discord:  discord,
	}
}

// NotifyAll sends notification to all enabled channels.
func (m *NotificationManager) NotifyAll(ctx context.Context, notification Notification) error {
	var lastErr error

	// Email
	if m.email != nil && len(notification.EmailRecipients) > 0 {
		if err := m.email.SendEmail(ctx, notification.EmailRecipients, notification.Subject, notification.Body); err != nil {
			lastErr = fmt.Errorf("email notification failed: %w", err)
		}
	}

	// Telegram
	if m.telegram != nil && notification.TelegramChatID != "" {
		if err := m.telegram.SendMessage(ctx, notification.TelegramChatID, notification.Message); err != nil {
			lastErr = fmt.Errorf("telegram notification failed: %w", err)
		}
	}

	// LINE
	if m.line != nil && notification.LINEUserID != "" {
		if err := m.line.PushMessage(ctx, notification.LINEUserID, notification.Message); err != nil {
			lastErr = fmt.Errorf("LINE notification failed: %w", err)
		}
	}

	// Discord
	if m.discord != nil {
		if err := m.discord.SendMessage(ctx, notification.Message); err != nil {
			lastErr = fmt.Errorf("discord notification failed: %w", err)
		}
	}

	return lastErr
}

// Notification represents a multi-channel notification.
type Notification struct {
	Subject          string   `json:"subject"`
	Message          string   `json:"message"`
	Body             string   `json:"body"` // HTML body for email
	EmailRecipients  []string `json:"emailRecipients,omitempty"`
	TelegramChatID   string   `json:"telegramChatId,omitempty"`
	LINEUserID       string   `json:"lineUserId,omitempty"`
}

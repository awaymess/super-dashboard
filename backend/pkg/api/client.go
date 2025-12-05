package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a base HTTP client for API requests.
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	headers    map[string]string
	rateLimit  *RateLimiter
}

// ClientConfig holds configuration for API client.
type ClientConfig struct {
	BaseURL        string
	APIKey         string
	Timeout        time.Duration
	RateLimitRPS   int // Requests per second
	CustomHeaders  map[string]string
}

// NewClient creates a new API client.
func NewClient(config ClientConfig) *Client {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	client := &Client{
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		baseURL: config.BaseURL,
		apiKey:  config.APIKey,
		headers: config.CustomHeaders,
	}

	if config.RateLimitRPS > 0 {
		client.rateLimit = NewRateLimiter(config.RateLimitRPS)
	}

	return client
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, endpoint string, params map[string]string) (*http.Response, error) {
	return c.request(ctx, http.MethodGet, endpoint, params, nil)
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, endpoint string, body interface{}) (*http.Response, error) {
	return c.request(ctx, http.MethodPost, endpoint, nil, body)
}

// request performs an HTTP request with rate limiting.
func (c *Client) request(ctx context.Context, method, endpoint string, params map[string]string, body interface{}) (*http.Response, error) {
	// Apply rate limiting
	if c.rateLimit != nil {
		if err := c.rateLimit.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limit wait: %w", err)
		}
	}

	// Build URL
	url := c.baseURL + endpoint
	if len(params) > 0 {
		url += "?"
		first := true
		for key, value := range params {
			if !first {
				url += "&"
			}
			url += fmt.Sprintf("%s=%s", key, value)
			first = false
		}
	}

	// Create request
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		reqBody = io.NopCloser(io.Reader(jsonData))
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}

	// Check status code
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}

// DecodeResponse decodes JSON response into target.
func DecodeResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	
	return nil
}

// RateLimiter implements token bucket rate limiting.
type RateLimiter struct {
	ticker   *time.Ticker
	tokens   chan struct{}
	maxTokens int
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	rl := &RateLimiter{
		ticker:    time.NewTicker(time.Second / time.Duration(requestsPerSecond)),
		tokens:    make(chan struct{}, requestsPerSecond),
		maxTokens: requestsPerSecond,
	}

	// Fill initial tokens
	for i := 0; i < requestsPerSecond; i++ {
		rl.tokens <- struct{}{}
	}

	// Refill tokens
	go func() {
		for range rl.ticker.C {
			select {
			case rl.tokens <- struct{}{}:
			default:
			}
		}
	}()

	return rl
}

// Wait blocks until a token is available.
func (rl *RateLimiter) Wait(ctx context.Context) error {
	select {
	case <-rl.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Stop stops the rate limiter.
func (rl *RateLimiter) Stop() {
	rl.ticker.Stop()
	close(rl.tokens)
}

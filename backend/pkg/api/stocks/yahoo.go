package stocks

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"super-dashboard/backend/pkg/api"
)

// YahooFinanceClient implements Yahoo Finance API client.
type YahooFinanceClient struct {
	client *api.Client
}

// NewYahooFinanceClient creates a new Yahoo Finance API client.
func NewYahooFinanceClient() *YahooFinanceClient {
	config := api.ClientConfig{
		BaseURL:      "https://query1.finance.yahoo.com",
		Timeout:      30 * time.Second,
		RateLimitRPS: 10, // Conservative rate limit
	}

	return &YahooFinanceClient{
		client: api.NewClient(config),
	}
}

// YahooQuote represents Yahoo Finance quote.
type YahooQuote struct {
	Symbol                 string    `json:"symbol"`
	RegularMarketPrice     float64   `json:"regularMarketPrice"`
	RegularMarketChange    float64   `json:"regularMarketChange"`
	RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
	RegularMarketVolume    int64     `json:"regularMarketVolume"`
	RegularMarketOpen      float64   `json:"regularMarketOpen"`
	RegularMarketDayHigh   float64   `json:"regularMarketDayHigh"`
	RegularMarketDayLow    float64   `json:"regularMarketDayLow"`
	RegularMarketPreviousClose float64 `json:"regularMarketPreviousClose"`
	MarketCap              float64   `json:"marketCap"`
	FiftyTwoWeekLow        float64   `json:"fiftyTwoWeekLow"`
	FiftyTwoWeekHigh       float64   `json:"fiftyTwoWeekHigh"`
	TrailingPE             float64   `json:"trailingPE"`
	ForwardPE              float64   `json:"forwardPE"`
	DividendYield          float64   `json:"dividendYield"`
	EPS                    float64   `json:"epsTrailingTwelveMonths"`
	Timestamp              time.Time `json:"timestamp"`
}

// YahooChart represents historical chart data.
type YahooChart struct {
	Symbol     string       `json:"symbol"`
	Timestamps []time.Time  `json:"timestamps"`
	Quotes     []YahooOHLCV `json:"quotes"`
}

// YahooOHLCV represents OHLCV data point.
type YahooOHLCV struct {
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`
}

// GetQuote retrieves real-time quote.
func (c *YahooFinanceClient) GetQuote(ctx context.Context, symbol string) (*YahooQuote, error) {
	params := map[string]string{
		"symbols": symbol,
	}

	resp, err := c.client.Get(ctx, "/v7/finance/quote", params)
	if err != nil {
		return nil, fmt.Errorf("get quote: %w", err)
	}

	var result struct {
		QuoteResponse struct {
			Result []struct {
				Symbol                      string  `json:"symbol"`
				RegularMarketPrice          float64 `json:"regularMarketPrice"`
				RegularMarketChange         float64 `json:"regularMarketChange"`
				RegularMarketChangePercent  float64 `json:"regularMarketChangePercent"`
				RegularMarketVolume         int64   `json:"regularMarketVolume"`
				RegularMarketOpen           float64 `json:"regularMarketOpen"`
				RegularMarketDayHigh        float64 `json:"regularMarketDayHigh"`
				RegularMarketDayLow         float64 `json:"regularMarketDayLow"`
				RegularMarketPreviousClose  float64 `json:"regularMarketPreviousClose"`
				MarketCap                   float64 `json:"marketCap"`
				FiftyTwoWeekLow             float64 `json:"fiftyTwoWeekLow"`
				FiftyTwoWeekHigh            float64 `json:"fiftyTwoWeekHigh"`
				TrailingPE                  float64 `json:"trailingPE"`
				ForwardPE                   float64 `json:"forwardPE"`
				DividendYield               float64 `json:"dividendYield"`
				EpsTrailingTwelveMonths     float64 `json:"epsTrailingTwelveMonths"`
				RegularMarketTime           int64   `json:"regularMarketTime"`
			} `json:"result"`
		} `json:"quoteResponse"`
	}

	if err := api.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.QuoteResponse.Result) == 0 {
		return nil, fmt.Errorf("no quote data for symbol: %s", symbol)
	}

	r := result.QuoteResponse.Result[0]
	quote := &YahooQuote{
		Symbol:                     r.Symbol,
		RegularMarketPrice:         r.RegularMarketPrice,
		RegularMarketChange:        r.RegularMarketChange,
		RegularMarketChangePercent: r.RegularMarketChangePercent,
		RegularMarketVolume:        r.RegularMarketVolume,
		RegularMarketOpen:          r.RegularMarketOpen,
		RegularMarketDayHigh:       r.RegularMarketDayHigh,
		RegularMarketDayLow:        r.RegularMarketDayLow,
		RegularMarketPreviousClose: r.RegularMarketPreviousClose,
		MarketCap:                  r.MarketCap,
		FiftyTwoWeekLow:            r.FiftyTwoWeekLow,
		FiftyTwoWeekHigh:           r.FiftyTwoWeekHigh,
		TrailingPE:                 r.TrailingPE,
		ForwardPE:                  r.ForwardPE,
		DividendYield:              r.DividendYield,
		EPS:                        r.EpsTrailingTwelveMonths,
		Timestamp:                  time.Unix(r.RegularMarketTime, 0),
	}

	return quote, nil
}

// GetChart retrieves historical chart data.
func (c *YahooFinanceClient) GetChart(ctx context.Context, symbol string, interval string, rangeStr string) (*YahooChart, error) {
	// interval: 1m, 5m, 15m, 1h, 1d, 1wk, 1mo
	// range: 1d, 5d, 1mo, 3mo, 6mo, 1y, 2y, 5y, 10y, ytd, max
	
	params := map[string]string{
		"interval": interval,
		"range":    rangeStr,
	}

	endpoint := fmt.Sprintf("/v8/finance/chart/%s", symbol)
	resp, err := c.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, fmt.Errorf("get chart: %w", err)
	}

	var result struct {
		Chart struct {
			Result []struct {
				Meta struct {
					Symbol string `json:"symbol"`
				} `json:"meta"`
				Timestamp  []int64 `json:"timestamp"`
				Indicators struct {
					Quote []struct {
						Open   []float64 `json:"open"`
						High   []float64 `json:"high"`
						Low    []float64 `json:"low"`
						Close  []float64 `json:"close"`
						Volume []int64   `json:"volume"`
					} `json:"quote"`
				} `json:"indicators"`
			} `json:"result"`
		} `json:"chart"`
	}

	if err := api.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Chart.Result) == 0 {
		return nil, fmt.Errorf("no chart data for symbol: %s", symbol)
	}

	r := result.Chart.Result[0]
	chart := &YahooChart{
		Symbol:     r.Meta.Symbol,
		Timestamps: make([]time.Time, len(r.Timestamp)),
		Quotes:     make([]YahooOHLCV, len(r.Timestamp)),
	}

	// Convert timestamps
	for i, ts := range r.Timestamp {
		chart.Timestamps[i] = time.Unix(ts, 0)
	}

	// Extract OHLCV data
	if len(r.Indicators.Quote) > 0 {
		q := r.Indicators.Quote[0]
		for i := range r.Timestamp {
			if i < len(q.Open) {
				chart.Quotes[i] = YahooOHLCV{
					Open:   q.Open[i],
					High:   q.High[i],
					Low:    q.Low[i],
					Close:  q.Close[i],
					Volume: q.Volume[i],
				}
			}
		}
	}

	return chart, nil
}

// GetHistoricalCSV retrieves historical data in CSV format (alternative method).
func (c *YahooFinanceClient) GetHistoricalCSV(ctx context.Context, symbol string, startDate, endDate time.Time) ([]PricePoint, error) {
	params := map[string]string{
		"period1": fmt.Sprintf("%d", startDate.Unix()),
		"period2": fmt.Sprintf("%d", endDate.Unix()),
		"interval": "1d",
		"events": "history",
	}

	endpoint := fmt.Sprintf("/v7/finance/download/%s", symbol)
	resp, err := c.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, fmt.Errorf("get historical CSV: %w", err)
	}
	defer resp.Body.Close()

	// Parse CSV
	reader := csv.NewReader(resp.Body)
	
	// Skip header
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("read CSV header: %w", err)
	}

	var prices []PricePoint

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read CSV row: %w", err)
		}

		if len(record) < 6 {
			continue
		}

		date, _ := time.Parse("2006-01-02", record[0])
		open, _ := strconv.ParseFloat(record[1], 64)
		high, _ := strconv.ParseFloat(record[2], 64)
		low, _ := strconv.ParseFloat(record[3], 64)
		close, _ := strconv.ParseFloat(record[4], 64)
		volume, _ := strconv.ParseInt(record[6], 10, 64)

		prices = append(prices, PricePoint{
			Date:   date,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: volume,
		})
	}

	return prices, nil
}

// SearchSymbol searches for stock symbols.
func (c *YahooFinanceClient) SearchSymbol(ctx context.Context, query string) ([]SearchResult, error) {
	params := map[string]string{
		"q": query,
	}

	resp, err := c.client.Get(ctx, "/v1/finance/search", params)
	if err != nil {
		return nil, fmt.Errorf("search symbol: %w", err)
	}

	var result struct {
		Quotes []struct {
			Symbol    string `json:"symbol"`
			ShortName string `json:"shortname"`
			LongName  string `json:"longname"`
			Exchange  string `json:"exchange"`
			QuoteType string `json:"quoteType"`
		} `json:"quotes"`
	}

	if err := api.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	results := make([]SearchResult, len(result.Quotes))
	for i, q := range result.Quotes {
		name := q.LongName
		if name == "" {
			name = q.ShortName
		}
		
		results[i] = SearchResult{
			Symbol:   q.Symbol,
			Name:     name,
			Exchange: q.Exchange,
			Type:     q.QuoteType,
		}
	}

	return results, nil
}

// SearchResult represents a search result.
type SearchResult struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Exchange string `json:"exchange"`
	Type     string `json:"type"`
}

// GetMultipleQuotes retrieves quotes for multiple symbols.
func (c *YahooFinanceClient) GetMultipleQuotes(ctx context.Context, symbols []string) ([]YahooQuote, error) {
	symbolsStr := strings.Join(symbols, ",")
	
	params := map[string]string{
		"symbols": symbolsStr,
	}

	resp, err := c.client.Get(ctx, "/v7/finance/quote", params)
	if err != nil {
		return nil, fmt.Errorf("get multiple quotes: %w", err)
	}

	var result struct {
		QuoteResponse struct {
			Result []struct {
				Symbol                      string  `json:"symbol"`
				RegularMarketPrice          float64 `json:"regularMarketPrice"`
				RegularMarketChange         float64 `json:"regularMarketChange"`
				RegularMarketChangePercent  float64 `json:"regularMarketChangePercent"`
				RegularMarketVolume         int64   `json:"regularMarketVolume"`
				RegularMarketOpen           float64 `json:"regularMarketOpen"`
				RegularMarketDayHigh        float64 `json:"regularMarketDayHigh"`
				RegularMarketDayLow         float64 `json:"regularMarketDayLow"`
				RegularMarketPreviousClose  float64 `json:"regularMarketPreviousClose"`
				MarketCap                   float64 `json:"marketCap"`
				RegularMarketTime           int64   `json:"regularMarketTime"`
			} `json:"result"`
		} `json:"quoteResponse"`
	}

	if err := api.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	quotes := make([]YahooQuote, len(result.QuoteResponse.Result))
	for i, r := range result.QuoteResponse.Result {
		quotes[i] = YahooQuote{
			Symbol:                     r.Symbol,
			RegularMarketPrice:         r.RegularMarketPrice,
			RegularMarketChange:        r.RegularMarketChange,
			RegularMarketChangePercent: r.RegularMarketChangePercent,
			RegularMarketVolume:        r.RegularMarketVolume,
			RegularMarketOpen:          r.RegularMarketOpen,
			RegularMarketDayHigh:       r.RegularMarketDayHigh,
			RegularMarketDayLow:        r.RegularMarketDayLow,
			RegularMarketPreviousClose: r.RegularMarketPreviousClose,
			MarketCap:                  r.MarketCap,
			Timestamp:                  time.Unix(r.RegularMarketTime, 0),
		}
	}

	return quotes, nil
}

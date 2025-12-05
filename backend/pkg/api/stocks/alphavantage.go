package stocks

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"super-dashboard/backend/pkg/api"
)

// AlphaVantageClient implements Alpha Vantage API client.
type AlphaVantageClient struct {
	client *api.Client
	apiKey string
}

// NewAlphaVantageClient creates a new Alpha Vantage API client.
func NewAlphaVantageClient(apiKey string) *AlphaVantageClient {
	config := api.ClientConfig{
		BaseURL:      "https://www.alphavantage.co/query",
		APIKey:       apiKey,
		Timeout:      30 * time.Second,
		RateLimitRPS: 1, // Free tier: 5 calls per minute = ~1 per 12 seconds
	}

	return &AlphaVantageClient{
		client: api.NewClient(config),
		apiKey: apiKey,
	}
}

// Quote represents a stock quote.
type Quote struct {
	Symbol           string    `json:"symbol"`
	Price            float64   `json:"price"`
	Change           float64   `json:"change"`
	ChangePercent    float64   `json:"changePercent"`
	Volume           int64     `json:"volume"`
	Open             float64   `json:"open"`
	High             float64   `json:"high"`
	Low              float64   `json:"low"`
	PreviousClose    float64   `json:"previousClose"`
	Timestamp        time.Time `json:"timestamp"`
}

// TimeSeriesDaily represents daily price data.
type TimeSeriesDaily struct {
	Symbol     string          `json:"symbol"`
	Interval   string          `json:"interval"`
	TimeSeries []PricePoint    `json:"timeSeries"`
	LastRefreshed time.Time    `json:"lastRefreshed"`
}

// PricePoint represents a single price data point.
type PricePoint struct {
	Date      time.Time `json:"date"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    int64     `json:"volume"`
}

// CompanyOverview represents fundamental data.
type CompanyOverview struct {
	Symbol                string  `json:"Symbol"`
	Name                  string  `json:"Name"`
	Description           string  `json:"Description"`
	Exchange              string  `json:"Exchange"`
	Currency              string  `json:"Currency"`
	Country               string  `json:"Country"`
	Sector                string  `json:"Sector"`
	Industry              string  `json:"Industry"`
	MarketCapitalization  float64 `json:"MarketCapitalization,string"`
	PERatio               float64 `json:"PERatio,string"`
	PEGRatio              float64 `json:"PEGRatio,string"`
	BookValue             float64 `json:"BookValue,string"`
	DividendPerShare      float64 `json:"DividendPerShare,string"`
	DividendYield         float64 `json:"DividendYield,string"`
	EPS                   float64 `json:"EPS,string"`
	RevenuePerShareTTM    float64 `json:"RevenuePerShareTTM,string"`
	ProfitMargin          float64 `json:"ProfitMargin,string"`
	OperatingMarginTTM    float64 `json:"OperatingMarginTTM,string"`
	ReturnOnAssetsTTM     float64 `json:"ReturnOnAssetsTTM,string"`
	ReturnOnEquityTTM     float64 `json:"ReturnOnEquityTTM,string"`
	RevenueTTM            float64 `json:"RevenueTTM,string"`
	GrossProfitTTM        float64 `json:"GrossProfitTTM,string"`
	DilutedEPSTTM         float64 `json:"DilutedEPSTTM,string"`
	QuarterlyEarningsGrowthYOY float64 `json:"QuarterlyEarningsGrowthYOY,string"`
	QuarterlyRevenueGrowthYOY  float64 `json:"QuarterlyRevenueGrowthYOY,string"`
	AnalystTargetPrice    float64 `json:"AnalystTargetPrice,string"`
	TrailingPE            float64 `json:"TrailingPE,string"`
	ForwardPE             float64 `json:"ForwardPE,string"`
	PriceToSalesRatioTTM  float64 `json:"PriceToSalesRatioTTM,string"`
	PriceToBookRatio      float64 `json:"PriceToBookRatio,string"`
	EVToRevenue           float64 `json:"EVToRevenue,string"`
	EVToEBITDA            float64 `json:"EVToEBITDA,string"`
	Beta                  float64 `json:"Beta,string"`
	High52Week            float64 `json:"52WeekHigh,string"`
	Low52Week             float64 `json:"52WeekLow,string"`
	MovingAverage50Day    float64 `json:"50DayMovingAverage,string"`
	MovingAverage200Day   float64 `json:"200DayMovingAverage,string"`
}

// GetQuote retrieves real-time quote for a symbol.
func (c *AlphaVantageClient) GetQuote(ctx context.Context, symbol string) (*Quote, error) {
	params := map[string]string{
		"function": "GLOBAL_QUOTE",
		"symbol":   symbol,
		"apikey":   c.apiKey,
	}

	resp, err := c.client.Get(ctx, "", params)
	if err != nil {
		return nil, fmt.Errorf("get quote: %w", err)
	}

	var result struct {
		GlobalQuote struct {
			Symbol           string `json:"01. symbol"`
			Price            string `json:"05. price"`
			Change           string `json:"09. change"`
			ChangePercent    string `json:"10. change percent"`
			Volume           string `json:"06. volume"`
			Open             string `json:"02. open"`
			High             string `json:"03. high"`
			Low              string `json:"04. low"`
			PreviousClose    string `json:"08. previous close"`
			LatestTradingDay string `json:"07. latest trading day"`
		} `json:"Global Quote"`
	}

	if err := api.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	gq := result.GlobalQuote
	quote := &Quote{
		Symbol: gq.Symbol,
	}

	// Parse numeric values
	quote.Price, _ = strconv.ParseFloat(gq.Price, 64)
	quote.Change, _ = strconv.ParseFloat(gq.Change, 64)
	quote.Volume, _ = strconv.ParseInt(gq.Volume, 10, 64)
	quote.Open, _ = strconv.ParseFloat(gq.Open, 64)
	quote.High, _ = strconv.ParseFloat(gq.High, 64)
	quote.Low, _ = strconv.ParseFloat(gq.Low, 64)
	quote.PreviousClose, _ = strconv.ParseFloat(gq.PreviousClose, 64)
	
	// Parse change percent (remove %)
	changePercentStr := gq.ChangePercent
	if len(changePercentStr) > 0 && changePercentStr[len(changePercentStr)-1] == '%' {
		changePercentStr = changePercentStr[:len(changePercentStr)-1]
	}
	quote.ChangePercent, _ = strconv.ParseFloat(changePercentStr, 64)

	// Parse date
	quote.Timestamp, _ = time.Parse("2006-01-02", gq.LatestTradingDay)

	return quote, nil
}

// GetDailyTimeSeries retrieves daily time series data.
func (c *AlphaVantageClient) GetDailyTimeSeries(ctx context.Context, symbol string, fullOutput bool) (*TimeSeriesDaily, error) {
	outputSize := "compact" // Last 100 data points
	if fullOutput {
		outputSize = "full" // Full history (20+ years)
	}

	params := map[string]string{
		"function":   "TIME_SERIES_DAILY",
		"symbol":     symbol,
		"outputsize": outputSize,
		"apikey":     c.apiKey,
	}

	resp, err := c.client.Get(ctx, "", params)
	if err != nil {
		return nil, fmt.Errorf("get daily time series: %w", err)
	}

	var result struct {
		MetaData struct {
			Symbol        string `json:"2. Symbol"`
			LastRefreshed string `json:"3. Last Refreshed"`
		} `json:"Meta Data"`
		TimeSeriesDaily map[string]struct {
			Open   string `json:"1. open"`
			High   string `json:"2. high"`
			Low    string `json:"3. low"`
			Close  string `json:"4. close"`
			Volume string `json:"5. volume"`
		} `json:"Time Series (Daily)"`
	}

	if err := api.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	ts := &TimeSeriesDaily{
		Symbol:   result.MetaData.Symbol,
		Interval: "daily",
		TimeSeries: make([]PricePoint, 0, len(result.TimeSeriesDaily)),
	}

	ts.LastRefreshed, _ = time.Parse("2006-01-02", result.MetaData.LastRefreshed)

	// Convert map to slice
	for dateStr, data := range result.TimeSeriesDaily {
		date, _ := time.Parse("2006-01-02", dateStr)
		
		point := PricePoint{
			Date: date,
		}
		
		point.Open, _ = strconv.ParseFloat(data.Open, 64)
		point.High, _ = strconv.ParseFloat(data.High, 64)
		point.Low, _ = strconv.ParseFloat(data.Low, 64)
		point.Close, _ = strconv.ParseFloat(data.Close, 64)
		point.Volume, _ = strconv.ParseInt(data.Volume, 10, 64)
		
		ts.TimeSeries = append(ts.TimeSeries, point)
	}

	return ts, nil
}

// GetCompanyOverview retrieves fundamental data and financial ratios.
func (c *AlphaVantageClient) GetCompanyOverview(ctx context.Context, symbol string) (*CompanyOverview, error) {
	params := map[string]string{
		"function": "OVERVIEW",
		"symbol":   symbol,
		"apikey":   c.apiKey,
	}

	resp, err := c.client.Get(ctx, "", params)
	if err != nil {
		return nil, fmt.Errorf("get company overview: %w", err)
	}

	var overview CompanyOverview
	if err := api.DecodeResponse(resp, &overview); err != nil {
		return nil, err
	}

	return &overview, nil
}

// TechnicalIndicator represents technical indicator data.
type TechnicalIndicator struct {
	Symbol    string                 `json:"symbol"`
	Indicator string                 `json:"indicator"`
	Interval  string                 `json:"interval"`
	Data      []IndicatorDataPoint   `json:"data"`
}

// IndicatorDataPoint represents a single indicator data point.
type IndicatorDataPoint struct {
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
}

// GetSMA retrieves Simple Moving Average.
func (c *AlphaVantageClient) GetSMA(ctx context.Context, symbol string, interval string, timePeriod int) (*TechnicalIndicator, error) {
	params := map[string]string{
		"function":    "SMA",
		"symbol":      symbol,
		"interval":    interval,
		"time_period": fmt.Sprintf("%d", timePeriod),
		"series_type": "close",
		"apikey":      c.apiKey,
	}

	resp, err := c.client.Get(ctx, "", params)
	if err != nil {
		return nil, fmt.Errorf("get SMA: %w", err)
	}

	var result struct {
		MetaData struct {
			Symbol string `json:"1: Symbol"`
		} `json:"Meta Data"`
		TechnicalAnalysis map[string]struct {
			SMA string `json:"SMA"`
		} `json:"Technical Analysis: SMA"`
	}

	if err := api.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	indicator := &TechnicalIndicator{
		Symbol:    result.MetaData.Symbol,
		Indicator: "SMA",
		Interval:  interval,
		Data:      make([]IndicatorDataPoint, 0, len(result.TechnicalAnalysis)),
	}

	for dateStr, data := range result.TechnicalAnalysis {
		date, _ := time.Parse("2006-01-02", dateStr)
		value, _ := strconv.ParseFloat(data.SMA, 64)
		
		indicator.Data = append(indicator.Data, IndicatorDataPoint{
			Date:  date,
			Value: value,
		})
	}

	return indicator, nil
}

// GetRSI retrieves Relative Strength Index.
func (c *AlphaVantageClient) GetRSI(ctx context.Context, symbol string, interval string, timePeriod int) (*TechnicalIndicator, error) {
	params := map[string]string{
		"function":    "RSI",
		"symbol":      symbol,
		"interval":    interval,
		"time_period": fmt.Sprintf("%d", timePeriod),
		"series_type": "close",
		"apikey":      c.apiKey,
	}

	resp, err := c.client.Get(ctx, "", params)
	if err != nil {
		return nil, fmt.Errorf("get RSI: %w", err)
	}

	var result struct {
		MetaData struct {
			Symbol string `json:"1: Symbol"`
		} `json:"Meta Data"`
		TechnicalAnalysis map[string]struct {
			RSI string `json:"RSI"`
		} `json:"Technical Analysis: RSI"`
	}

	if err := api.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	indicator := &TechnicalIndicator{
		Symbol:    result.MetaData.Symbol,
		Indicator: "RSI",
		Interval:  interval,
		Data:      make([]IndicatorDataPoint, 0, len(result.TechnicalAnalysis)),
	}

	for dateStr, data := range result.TechnicalAnalysis {
		date, _ := time.Parse("2006-01-02", dateStr)
		value, _ := strconv.ParseFloat(data.RSI, 64)
		
		indicator.Data = append(indicator.Data, IndicatorDataPoint{
			Date:  date,
			Value: value,
		})
	}

	return indicator, nil
}

package odds

import (
	"context"
	"fmt"
	"time"

	"super-dashboard/backend/pkg/api"
)

// BetfairClient implements Betfair Exchange API client.
type BetfairClient struct {
	client    *api.Client
	appKey    string
	sessionToken string
}

// NewBetfairClient creates a new Betfair API client.
func NewBetfairClient(appKey, sessionToken string) *BetfairClient {
	config := api.ClientConfig{
		BaseURL:      "https://api.betfair.com/exchange/betting/json-rpc/v1",
		Timeout:      30 * time.Second,
		RateLimitRPS: 5, // Betfair allows 5 requests/second
		CustomHeaders: map[string]string{
			"X-Application": appKey,
			"X-Authentication": sessionToken,
		},
	}

	return &BetfairClient{
		client:       api.NewClient(config),
		appKey:       appKey,
		sessionToken: sessionToken,
	}
}

// EventType represents a sport type in Betfair.
type EventType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Competition represents a competition/league.
type Competition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Event represents a betting event/match.
type Event struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CountryCode string    `json:"countryCode"`
	Timezone    string    `json:"timezone"`
	OpenDate    time.Time `json:"openDate"`
}

// MarketCatalogue represents available betting markets.
type MarketCatalogue struct {
	MarketID     string       `json:"marketId"`
	MarketName   string       `json:"marketName"`
	MarketType   string       `json:"marketType"`
	Competition  Competition  `json:"competition"`
	Event        Event        `json:"event"`
	Runners      []Runner     `json:"runners"`
	TotalMatched float64      `json:"totalMatched"`
}

// Runner represents a selection in a market.
type Runner struct {
	SelectionID int64  `json:"selectionId"`
	RunnerName  string `json:"runnerName"`
	Handicap    float64 `json:"handicap,omitempty"`
	SortPriority int    `json:"sortPriority"`
}

// MarketBook represents current odds and liquidity.
type MarketBook struct {
	MarketID    string         `json:"marketId"`
	IsMarketDataDelayed bool   `json:"isMarketDataDelayed"`
	Status      string         `json:"status"`
	Runners     []RunnerBook   `json:"runners"`
	TotalMatched float64       `json:"totalMatched"`
}

// RunnerBook represents current odds for a runner.
type RunnerBook struct {
	SelectionID   int64           `json:"selectionId"`
	Status        string          `json:"status"`
	LastPriceTraded float64       `json:"lastPriceTraded,omitempty"`
	TotalMatched  float64         `json:"totalMatched"`
	ExchangePrices ExchangePrices `json:"ex"`
}

// ExchangePrices represents available odds.
type ExchangePrices struct {
	AvailableToBack []PriceSize `json:"availableToBack"`
	AvailableToLay  []PriceSize `json:"availableToLay"`
}

// PriceSize represents odds and available stake.
type PriceSize struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

// GetEventTypes retrieves all sport types.
func (c *BetfairClient) GetEventTypes(ctx context.Context) ([]EventType, error) {
	body := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "SportsAPING/v1.0/listEventTypes",
		"params": map[string]interface{}{
			"filter": map[string]interface{}{},
		},
		"id": 1,
	}

	resp, err := c.client.Post(ctx, "", body)
	if err != nil {
		return nil, fmt.Errorf("get event types: %w", err)
	}

	var result struct {
		Result []struct {
			EventType EventType `json:"eventType"`
		} `json:"result"`
	}

	if err := api.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	eventTypes := make([]EventType, len(result.Result))
	for i, item := range result.Result {
		eventTypes[i] = item.EventType
	}

	return eventTypes, nil
}

// GetCompetitions retrieves competitions for an event type.
func (c *BetfairClient) GetCompetitions(ctx context.Context, eventTypeID string) ([]Competition, error) {
	body := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "SportsAPING/v1.0/listCompetitions",
		"params": map[string]interface{}{
			"filter": map[string]interface{}{
				"eventTypeIds": []string{eventTypeID},
			},
		},
		"id": 1,
	}

	resp, err := c.client.Post(ctx, "", body)
	if err != nil {
		return nil, fmt.Errorf("get competitions: %w", err)
	}

	var result struct {
		Result []struct {
			Competition Competition `json:"competition"`
		} `json:"result"`
	}

	if err := api.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	competitions := make([]Competition, len(result.Result))
	for i, item := range result.Result {
		competitions[i] = item.Competition
	}

	return competitions, nil
}

// GetMarkets retrieves available betting markets.
func (c *BetfairClient) GetMarkets(ctx context.Context, eventTypeID, competitionID string) ([]MarketCatalogue, error) {
	filter := map[string]interface{}{
		"eventTypeIds": []string{eventTypeID},
	}

	if competitionID != "" {
		filter["competitionIds"] = []string{competitionID}
	}

	body := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "SportsAPING/v1.0/listMarketCatalogue",
		"params": map[string]interface{}{
			"filter": filter,
			"marketProjection": []string{
				"COMPETITION",
				"EVENT",
				"RUNNER_METADATA",
				"MARKET_START_TIME",
			},
			"maxResults": 100,
		},
		"id": 1,
	}

	resp, err := c.client.Post(ctx, "", body)
	if err != nil {
		return nil, fmt.Errorf("get markets: %w", err)
	}

	var result struct {
		Result []MarketCatalogue `json:"result"`
	}

	if err := api.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Result, nil
}

// GetMarketOdds retrieves current odds for markets.
func (c *BetfairClient) GetMarketOdds(ctx context.Context, marketIDs []string) ([]MarketBook, error) {
	body := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "SportsAPING/v1.0/listMarketBook",
		"params": map[string]interface{}{
			"marketIds": marketIDs,
			"priceProjection": map[string]interface{}{
				"priceData": []string{"EX_BEST_OFFERS"},
			},
		},
		"id": 1,
	}

	resp, err := c.client.Post(ctx, "", body)
	if err != nil {
		return nil, fmt.Errorf("get market odds: %w", err)
	}

	var result struct {
		Result []MarketBook `json:"result"`
	}

	if err := api.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Result, nil
}

package odds

import (
	"context"
	"fmt"
	"time"

	"super-dashboard/backend/pkg/api"
)

// PinnacleClient implements Pinnacle API client.
type PinnacleClient struct {
	client *api.Client
}

// NewPinnacleClient creates a new Pinnacle API client.
func NewPinnacleClient(apiKey string) *PinnacleClient {
	config := api.ClientConfig{
		BaseURL:       "https://api.pinnacle.com/v1",
		APIKey:        apiKey,
		Timeout:       30 * time.Second,
		RateLimitRPS:  10, // Pinnacle allows 10 requests/second
		CustomHeaders: map[string]string{
			"X-API-Key": apiKey,
		},
	}

	return &PinnacleClient{
		client: api.NewClient(config),
	}
}

// Sport represents a sport in Pinnacle.
type Sport struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// League represents a league in Pinnacle.
type League struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	SportID int    `json:"sportId"`
}

// Match represents a match/event.
type Match struct {
	ID          int64     `json:"id"`
	SportID     int       `json:"sportId"`
	LeagueID    int       `json:"leagueId"`
	HomeTeam    string    `json:"home"`
	AwayTeam    string    `json:"away"`
	StartTime   time.Time `json:"starts"`
	Live        bool      `json:"live"`
	HomeScore   *int      `json:"homeScore,omitempty"`
	AwayScore   *int      `json:"awayScore,omitempty"`
}

// Odds represents betting odds.
type Odds struct {
	MatchID     int64              `json:"matchId"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	Moneyline   *MoneylineOdds     `json:"moneyline,omitempty"`
	Spread      *SpreadOdds        `json:"spread,omitempty"`
	Total       *TotalOdds         `json:"total,omitempty"`
}

// MoneylineOdds represents 1X2 or moneyline odds.
type MoneylineOdds struct {
	Home float64 `json:"home"`
	Draw float64 `json:"draw,omitempty"`
	Away float64 `json:"away"`
}

// SpreadOdds represents handicap/spread odds.
type SpreadOdds struct {
	HomeSpread float64 `json:"homeSpread"`
	HomeOdds   float64 `json:"homeOdds"`
	AwaySpread float64 `json:"awaySpread"`
	AwayOdds   float64 `json:"awayOdds"`
}

// TotalOdds represents over/under odds.
type TotalOdds struct {
	Points    float64 `json:"points"`
	OverOdds  float64 `json:"overOdds"`
	UnderOdds float64 `json:"underOdds"`
}

// GetSports retrieves all available sports.
func (c *PinnacleClient) GetSports(ctx context.Context) ([]Sport, error) {
	resp, err := c.client.Get(ctx, "/sports", nil)
	if err != nil {
		return nil, fmt.Errorf("get sports: %w", err)
	}

	var result struct {
		Sports []Sport `json:"sports"`
	}

	if err := api.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Sports, nil
}

// GetLeagues retrieves leagues for a sport.
func (c *PinnacleClient) GetLeagues(ctx context.Context, sportID int) ([]League, error) {
	params := map[string]string{
		"sportId": fmt.Sprintf("%d", sportID),
	}

	resp, err := c.client.Get(ctx, "/leagues", params)
	if err != nil {
		return nil, fmt.Errorf("get leagues: %w", err)
	}

	var result struct {
		Leagues []League `json:"leagues"`
	}

	if err := api.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Leagues, nil
}

// GetMatches retrieves upcoming matches.
func (c *PinnacleClient) GetMatches(ctx context.Context, sportID, leagueID int) ([]Match, error) {
	params := map[string]string{
		"sportId": fmt.Sprintf("%d", sportID),
	}

	if leagueID > 0 {
		params["leagueId"] = fmt.Sprintf("%d", leagueID)
	}

	resp, err := c.client.Get(ctx, "/fixtures", params)
	if err != nil {
		return nil, fmt.Errorf("get matches: %w", err)
	}

	var result struct {
		Fixtures []Match `json:"fixtures"`
	}

	if err := api.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Fixtures, nil
}

// GetOdds retrieves odds for matches.
func (c *PinnacleClient) GetOdds(ctx context.Context, sportID, leagueID int, oddsFormat string) ([]Odds, error) {
	if oddsFormat == "" {
		oddsFormat = "DECIMAL"
	}

	params := map[string]string{
		"sportId":    fmt.Sprintf("%d", sportID),
		"oddsFormat": oddsFormat,
	}

	if leagueID > 0 {
		params["leagueId"] = fmt.Sprintf("%d", leagueID)
	}

	resp, err := c.client.Get(ctx, "/odds", params)
	if err != nil {
		return nil, fmt.Errorf("get odds: %w", err)
	}

	var result struct {
		Odds []Odds `json:"odds"`
	}

	if err := api.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Odds, nil
}

// GetLiveMatches retrieves live matches with current scores.
func (c *PinnacleClient) GetLiveMatches(ctx context.Context, sportID int) ([]Match, error) {
	params := map[string]string{
		"sportId": fmt.Sprintf("%d", sportID),
		"live":    "true",
	}

	resp, err := c.client.Get(ctx, "/fixtures", params)
	if err != nil {
		return nil, fmt.Errorf("get live matches: %w", err)
	}

	var result struct {
		Fixtures []Match `json:"fixtures"`
	}

	if err := api.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Fixtures, nil
}

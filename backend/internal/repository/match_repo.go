package repository

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/awaymess/super-dashboard/backend/internal/model"
)

// ErrNotFound is returned when a requested resource is not found.
var ErrNotFound = errors.New("resource not found")

// MatchMockData represents the structure of the mock matches JSON file.
type MatchMockData struct {
	Teams   []TeamJSON  `json:"teams"`
	Matches []MatchJSON `json:"matches"`
	Odds    []OddsJSON  `json:"odds"`
}

// TeamJSON represents a team in the mock JSON format.
type TeamJSON struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Country string  `json:"country"`
	Elo     float64 `json:"elo"`
}

// MatchJSON represents a match in the mock JSON format.
type MatchJSON struct {
	ID         string `json:"id"`
	League     string `json:"league"`
	HomeTeamID string `json:"home_team_id"`
	AwayTeamID string `json:"away_team_id"`
	StartTime  string `json:"start_time"`
	Status     string `json:"status"`
	Venue      string `json:"venue"`
}

// OddsJSON represents odds in the mock JSON format.
type OddsJSON struct {
	ID        string  `json:"id"`
	MatchID   string  `json:"match_id"`
	Bookmaker string  `json:"bookmaker"`
	Market    string  `json:"market"`
	Outcome   string  `json:"outcome"`
	Price     float64 `json:"price"`
}

// MatchRepository defines the interface for match data operations.
type MatchRepository interface {
	GetAll() ([]model.Match, error)
	GetByID(id string) (*model.Match, error)
	GetOddsByMatchID(matchID string) ([]model.Odds, error)
}

// mockMatchRepository implements MatchRepository using mock JSON data.
type mockMatchRepository struct {
	teams   map[string]model.Team
	matches map[string]model.Match
	odds    map[string][]model.Odds
}

// NewMockMatchRepository creates a new mock match repository from a JSON file.
func NewMockMatchRepository(filePath string) (MatchRepository, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var mockData MatchMockData
	if err := json.Unmarshal(data, &mockData); err != nil {
		return nil, err
	}

	repo := &mockMatchRepository{
		teams:   make(map[string]model.Team),
		matches: make(map[string]model.Match),
		odds:    make(map[string][]model.Odds),
	}

	// Parse teams
	for _, t := range mockData.Teams {
		repo.teams[t.ID] = model.Team{
			ID:      stringToUUID(t.ID),
			Name:    t.Name,
			Country: t.Country,
			Elo:     t.Elo,
		}
	}

	// Parse matches
	for _, m := range mockData.Matches {
		startTime, _ := time.Parse(time.RFC3339, m.StartTime)
		homeTeam := repo.teams[m.HomeTeamID]
		awayTeam := repo.teams[m.AwayTeamID]

		repo.matches[m.ID] = model.Match{
			ID:         stringToUUID(m.ID),
			League:     m.League,
			HomeTeamID: homeTeam.ID,
			HomeTeam:   homeTeam,
			AwayTeamID: awayTeam.ID,
			AwayTeam:   awayTeam,
			StartTime:  startTime,
			Status:     m.Status,
			Venue:      m.Venue,
		}
	}

	// Parse odds
	for _, o := range mockData.Odds {
		matchOdds := model.Odds{
			ID:        stringToUUID(o.ID),
			MatchID:   stringToUUID(o.MatchID),
			Bookmaker: o.Bookmaker,
			Market:    o.Market,
			Outcome:   o.Outcome,
			Price:     o.Price,
		}
		repo.odds[o.MatchID] = append(repo.odds[o.MatchID], matchOdds)
	}

	return repo, nil
}

// GetAll returns all matches.
func (r *mockMatchRepository) GetAll() ([]model.Match, error) {
	matches := make([]model.Match, 0, len(r.matches))
	for _, m := range r.matches {
		matches = append(matches, m)
	}
	return matches, nil
}

// GetByID returns a match by ID.
func (r *mockMatchRepository) GetByID(id string) (*model.Match, error) {
	match, ok := r.matches[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &match, nil
}

// GetOddsByMatchID returns odds for a specific match.
func (r *mockMatchRepository) GetOddsByMatchID(matchID string) ([]model.Odds, error) {
	return r.odds[matchID], nil
}

// stringToUUID converts a string ID to a deterministic UUID.
func stringToUUID(id string) uuid.UUID {
	// Use UUID v5 (SHA-1 hash) with a namespace to generate deterministic UUIDs
	namespace := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8") // DNS namespace
	return uuid.NewSHA1(namespace, []byte(id))
}

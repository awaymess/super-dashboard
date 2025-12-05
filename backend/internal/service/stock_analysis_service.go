package service

import (
	"context"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"super-dashboard/backend/internal/model"
	"super-dashboard/backend/internal/repository"
)

// StockAnalysisService handles stock analysis and valuation.
type StockAnalysisService struct {
	fairValueRepo *repository.FairValueRepository
	stockRepo     *repository.StockRepository
	newsRepo      *repository.StockNewsRepository
	logger        zerolog.Logger
}

// NewStockAnalysisService creates a new StockAnalysisService.
func NewStockAnalysisService(
	fairValueRepo *repository.FairValueRepository,
	stockRepo *repository.StockRepository,
	newsRepo *repository.StockNewsRepository,
	logger zerolog.Logger,
) *StockAnalysisService {
	return &StockAnalysisService{
		fairValueRepo: fairValueRepo,
		stockRepo:     stockRepo,
		newsRepo:      newsRepo,
		logger:        logger.With().Str("service", "stock_analysis").Logger(),
	}
}

// CalculateDCFValue calculates Discounted Cash Flow fair value.
func (s *StockAnalysisService) CalculateDCFValue(ctx context.Context, symbol string, freeCashFlow, growthRate, discountRate float64, years int) (float64, error) {
	stock, err := s.stockRepo.GetBySymbol(ctx, symbol)
	if err != nil {
		return 0, fmt.Errorf("stock not found: %w", err)
	}

	// Calculate present value of future cash flows
	pv := 0.0
	currentCF := freeCashFlow

	for i := 1; i <= years; i++ {
		currentCF *= (1 + growthRate)
		discountFactor := math.Pow(1+discountRate, float64(i))
		pv += currentCF / discountFactor
	}

	// Terminal value (perpetuity growth model)
	terminalGrowth := 0.03 // 3% terminal growth rate
	terminalValue := (currentCF * (1 + terminalGrowth)) / (discountRate - terminalGrowth)
	pv += terminalValue / math.Pow(1+discountRate, float64(years))

	// Fair value per share
	sharesOutstanding := stock.MarketCap / stock.CurrentPrice
	fairValue := pv / sharesOutstanding

	// Save to database
	upside := ((fairValue - stock.CurrentPrice) / stock.CurrentPrice) * 100
	rating := s.getRating(upside)

	fv := &model.FairValue{
		Symbol:         symbol,
		Method:         "DCF",
		FairValue:      fairValue,
		CurrentPrice:   stock.CurrentPrice,
		UpsidePercent:  upside,
		Rating:         rating,
		Confidence:     s.calculateConfidence(stock, "DCF"),
	}

	if err := s.fairValueRepo.CreateFairValue(ctx, fv); err != nil {
		s.logger.Error().Err(err).Msg("Failed to save DCF valuation")
	}

	return fairValue, nil
}

// CalculateGrahamValue calculates Benjamin Graham intrinsic value.
func (s *StockAnalysisService) CalculateGrahamValue(ctx context.Context, symbol string, eps, bookValue float64) (float64, error) {
	stock, err := s.stockRepo.GetBySymbol(ctx, symbol)
	if err != nil {
		return 0, fmt.Errorf("stock not found: %w", err)
	}

	// Graham Number = sqrt(22.5 × EPS × Book Value)
	if eps <= 0 || bookValue <= 0 {
		return 0, fmt.Errorf("invalid inputs: EPS and Book Value must be positive")
	}

	fairValue := math.Sqrt(22.5 * eps * bookValue)

	// Save to database
	upside := ((fairValue - stock.CurrentPrice) / stock.CurrentPrice) * 100
	rating := s.getRating(upside)

	fv := &model.FairValue{
		Symbol:         symbol,
		Method:         "Graham",
		FairValue:      fairValue,
		CurrentPrice:   stock.CurrentPrice,
		UpsidePercent:  upside,
		Rating:         rating,
		Confidence:     s.calculateConfidence(stock, "Graham"),
	}

	if err := s.fairValueRepo.CreateFairValue(ctx, fv); err != nil {
		s.logger.Error().Err(err).Msg("Failed to save Graham valuation")
	}

	return fairValue, nil
}

// CalculatePEValue calculates fair value based on P/E ratio.
func (s *StockAnalysisService) CalculatePEValue(ctx context.Context, symbol string, eps, industryPE float64) (float64, error) {
	stock, err := s.stockRepo.GetBySymbol(ctx, symbol)
	if err != nil {
		return 0, fmt.Errorf("stock not found: %w", err)
	}

	fairValue := eps * industryPE

	// Save to database
	upside := ((fairValue - stock.CurrentPrice) / stock.CurrentPrice) * 100
	rating := s.getRating(upside)

	fv := &model.FairValue{
		Symbol:         symbol,
		Method:         "P/E",
		FairValue:      fairValue,
		CurrentPrice:   stock.CurrentPrice,
		UpsidePercent:  upside,
		Rating:         rating,
		Confidence:     s.calculateConfidence(stock, "PE"),
	}

	if err := s.fairValueRepo.CreateFairValue(ctx, fv); err != nil {
		s.logger.Error().Err(err).Msg("Failed to save P/E valuation")
	}

	return fairValue, nil
}

// getRating determines investment rating based on upside percentage.
func (s *StockAnalysisService) getRating(upside float64) string {
	switch {
	case upside >= 30:
		return "Strong Buy"
	case upside >= 15:
		return "Buy"
	case upside >= -5:
		return "Hold"
	case upside >= -15:
		return "Sell"
	default:
		return "Strong Sell"
	}
}

// calculateConfidence calculates confidence score for valuation.
func (s *StockAnalysisService) calculateConfidence(stock *model.Stock, method string) float64 {
	confidence := 50.0 // Base confidence

	// Adjust based on volume
	if stock.Volume > 1000000 {
		confidence += 10
	}

	// Adjust based on market cap (larger = more reliable data)
	if stock.MarketCap > 10000000000 { // >10B
		confidence += 15
	} else if stock.MarketCap > 1000000000 { // >1B
		confidence += 10
	}

	// Method-specific adjustments
	switch method {
	case "DCF":
		confidence += 5 // DCF is comprehensive
	case "Graham":
		confidence += 10 // Graham is conservative
	}

	if confidence > 95 {
		confidence = 95
	}

	return confidence
}

// GetLatestFairValue retrieves the latest fair value for a stock.
func (s *StockAnalysisService) GetLatestFairValue(ctx context.Context, symbol string) (*model.FairValue, error) {
	fv, err := s.fairValueRepo.GetLatestFairValue(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get fair value: %w", err)
	}

	return fv, nil
}

// GetUndervaluedStocks retrieves undervalued stocks.
func (s *StockAnalysisService) GetUndervaluedStocks(ctx context.Context, minUpside float64) ([]model.FairValue, error) {
	stocks, err := s.fairValueRepo.GetUndervaluedStocks(ctx, minUpside)
	if err != nil {
		return nil, fmt.Errorf("failed to get undervalued stocks: %w", err)
	}

	return stocks, nil
}

// GetStockWithSentiment retrieves stock analysis with sentiment data.
func (s *StockAnalysisService) GetStockWithSentiment(ctx context.Context, symbol string) (map[string]interface{}, error) {
	// Get stock
	stock, err := s.stockRepo.GetBySymbol(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("stock not found: %w", err)
	}

	// Get fair value
	fairValue, _ := s.fairValueRepo.GetLatestFairValue(ctx, symbol)

	// Get sentiment stats
	sentimentStats, _ := s.newsRepo.GetSentimentStats(ctx, symbol, 30)

	// Get recent news
	news, _ := s.newsRepo.GetNewsByStock(ctx, symbol, 5)

	return map[string]interface{}{
		"stock":      stock,
		"fair_value": fairValue,
		"sentiment":  sentimentStats,
		"news":       news,
	}, nil
}

// GetStockScreener screens stocks based on criteria.
func (s *StockAnalysisService) GetStockScreener(ctx context.Context, criteria map[string]interface{}) ([]model.Stock, error) {
	// This would implement complex filtering logic
	// For now, return all stocks
	stocks, err := s.stockRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get stocks: %w", err)
	}

	// Apply filters
	filtered := make([]model.Stock, 0)
	for _, stock := range stocks {
		if s.matchesCriteria(stock, criteria) {
			filtered = append(filtered, stock)
		}
	}

	return filtered, nil
}

// matchesCriteria checks if a stock matches screening criteria.
func (s *StockAnalysisService) matchesCriteria(stock model.Stock, criteria map[string]interface{}) bool {
	if minPrice, ok := criteria["min_price"].(float64); ok {
		if stock.CurrentPrice < minPrice {
			return false
		}
	}

	if maxPrice, ok := criteria["max_price"].(float64); ok {
		if stock.CurrentPrice > maxPrice {
			return false
		}
	}

	if minChange, ok := criteria["min_change"].(float64); ok {
		if stock.ChangePercent < minChange {
			return false
		}
	}

	if minVolume, ok := criteria["min_volume"].(int64); ok {
		if stock.Volume < minVolume {
			return false
		}
	}

	return true
}

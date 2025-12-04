package repository

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/superdashboard/backend/internal/model"
	"gorm.io/gorm"
)

func TestInMemoryPortfolioRepository(t *testing.T) {
	repo := NewInMemoryPortfolioRepository()

	t.Run("Create and GetByID", func(t *testing.T) {
		portfolio := &model.Portfolio{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			Name:        "Test Portfolio",
			CashBalance: 100000,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err := repo.Create(portfolio)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		got, err := repo.GetByID(portfolio.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}

		if got.ID != portfolio.ID {
			t.Errorf("GetByID() ID = %v, want %v", got.ID, portfolio.ID)
		}
	})

	t.Run("GetByID not found", func(t *testing.T) {
		_, err := repo.GetByID(uuid.New())
		if err != gorm.ErrRecordNotFound {
			t.Errorf("GetByID() error = %v, want %v", err, gorm.ErrRecordNotFound)
		}
	})

	t.Run("GetByUserID", func(t *testing.T) {
		userID := uuid.New()
		portfolio1 := &model.Portfolio{
			ID:          uuid.New(),
			UserID:      userID,
			Name:        "Portfolio 1",
			CashBalance: 50000,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		portfolio2 := &model.Portfolio{
			ID:          uuid.New(),
			UserID:      userID,
			Name:        "Portfolio 2",
			CashBalance: 75000,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		repo.Create(portfolio1)
		repo.Create(portfolio2)

		portfolios, err := repo.GetByUserID(userID)
		if err != nil {
			t.Fatalf("GetByUserID() error = %v", err)
		}

		if len(portfolios) < 2 {
			t.Errorf("GetByUserID() returned %d portfolios, want at least 2", len(portfolios))
		}
	})

	t.Run("Update", func(t *testing.T) {
		portfolio := &model.Portfolio{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			Name:        "Original Name",
			CashBalance: 100000,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		repo.Create(portfolio)

		portfolio.Name = "Updated Name"
		err := repo.Update(portfolio)
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		got, _ := repo.GetByID(portfolio.ID)
		if got.Name != "Updated Name" {
			t.Errorf("Update() name = %v, want 'Updated Name'", got.Name)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		portfolio := &model.Portfolio{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			Name:        "To Delete",
			CashBalance: 100000,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		repo.Create(portfolio)
		err := repo.Delete(portfolio.ID)
		if err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		_, err = repo.GetByID(portfolio.ID)
		if err != gorm.ErrRecordNotFound {
			t.Errorf("GetByID() after delete error = %v, want %v", err, gorm.ErrRecordNotFound)
		}
	})

	t.Run("List", func(t *testing.T) {
		portfolios, err := repo.List()
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if portfolios == nil {
			t.Error("List() returned nil")
		}
	})
}

func TestInMemoryPositionRepository(t *testing.T) {
	repo := NewInMemoryPositionRepository()

	t.Run("Create and GetByID", func(t *testing.T) {
		position := &model.Position{
			ID:          uuid.New(),
			PortfolioID: uuid.New(),
			Symbol:      "AAPL",
			Quantity:    100,
			AvgCost:     150.00,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err := repo.Create(position)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		got, err := repo.GetByID(position.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}

		if got.Symbol != "AAPL" {
			t.Errorf("GetByID() symbol = %v, want AAPL", got.Symbol)
		}
	})

	t.Run("GetByPortfolioID", func(t *testing.T) {
		portfolioID := uuid.New()
		position1 := &model.Position{
			ID:          uuid.New(),
			PortfolioID: portfolioID,
			Symbol:      "MSFT",
			Quantity:    50,
			AvgCost:     300.00,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		position2 := &model.Position{
			ID:          uuid.New(),
			PortfolioID: portfolioID,
			Symbol:      "GOOGL",
			Quantity:    25,
			AvgCost:     140.00,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		repo.Create(position1)
		repo.Create(position2)

		positions, err := repo.GetByPortfolioID(portfolioID)
		if err != nil {
			t.Fatalf("GetByPortfolioID() error = %v", err)
		}

		if len(positions) != 2 {
			t.Errorf("GetByPortfolioID() returned %d positions, want 2", len(positions))
		}
	})

	t.Run("GetByPortfolioAndSymbol", func(t *testing.T) {
		portfolioID := uuid.New()
		position := &model.Position{
			ID:          uuid.New(),
			PortfolioID: portfolioID,
			Symbol:      "NVDA",
			Quantity:    75,
			AvgCost:     450.00,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		repo.Create(position)

		got, err := repo.GetByPortfolioAndSymbol(portfolioID, "NVDA")
		if err != nil {
			t.Fatalf("GetByPortfolioAndSymbol() error = %v", err)
		}

		if got.Quantity != 75 {
			t.Errorf("GetByPortfolioAndSymbol() quantity = %v, want 75", got.Quantity)
		}
	})

	t.Run("GetByPortfolioAndSymbol not found", func(t *testing.T) {
		_, err := repo.GetByPortfolioAndSymbol(uuid.New(), "UNKNOWN")
		if err != gorm.ErrRecordNotFound {
			t.Errorf("GetByPortfolioAndSymbol() error = %v, want %v", err, gorm.ErrRecordNotFound)
		}
	})
}

func TestInMemoryOrderRepository(t *testing.T) {
	repo := NewInMemoryOrderRepository()

	t.Run("Create and GetByID", func(t *testing.T) {
		now := time.Now()
		order := &model.Order{
			ID:          uuid.New(),
			PortfolioID: uuid.New(),
			Symbol:      "AAPL",
			Side:        model.OrderSideBuy,
			OrderType:   model.OrderTypeMarket,
			Quantity:    10,
			Price:       150.00,
			Status:      model.OrderStatusFilled,
			FilledAt:    &now,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		err := repo.Create(order)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		got, err := repo.GetByID(order.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}

		if got.Status != model.OrderStatusFilled {
			t.Errorf("GetByID() status = %v, want %v", got.Status, model.OrderStatusFilled)
		}
	})

	t.Run("GetByPortfolioID", func(t *testing.T) {
		portfolioID := uuid.New()
		now := time.Now()
		order1 := &model.Order{
			ID:          uuid.New(),
			PortfolioID: portfolioID,
			Symbol:      "MSFT",
			Side:        model.OrderSideBuy,
			OrderType:   model.OrderTypeMarket,
			Quantity:    5,
			Price:       300.00,
			Status:      model.OrderStatusFilled,
			FilledAt:    &now,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		repo.Create(order1)

		orders, err := repo.GetByPortfolioID(portfolioID)
		if err != nil {
			t.Fatalf("GetByPortfolioID() error = %v", err)
		}

		if len(orders) != 1 {
			t.Errorf("GetByPortfolioID() returned %d orders, want 1", len(orders))
		}
	})
}

func TestInMemoryTradeRepository(t *testing.T) {
	repo := NewInMemoryTradeRepository()

	t.Run("Create and GetByID", func(t *testing.T) {
		trade := &model.Trade{
			ID:          uuid.New(),
			PortfolioID: uuid.New(),
			OrderID:     uuid.New(),
			Symbol:      "AAPL",
			Side:        model.OrderSideBuy,
			Quantity:    10,
			Price:       150.00,
			Total:       1500.00,
			ExecutedAt:  time.Now(),
		}

		err := repo.Create(trade)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		got, err := repo.GetByID(trade.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}

		if got.Total != 1500.00 {
			t.Errorf("GetByID() total = %v, want 1500.00", got.Total)
		}
	})

	t.Run("GetByPortfolioID", func(t *testing.T) {
		portfolioID := uuid.New()
		trade := &model.Trade{
			ID:          uuid.New(),
			PortfolioID: portfolioID,
			OrderID:     uuid.New(),
			Symbol:      "MSFT",
			Side:        model.OrderSideBuy,
			Quantity:    5,
			Price:       300.00,
			Total:       1500.00,
			ExecutedAt:  time.Now(),
		}

		repo.Create(trade)

		trades, err := repo.GetByPortfolioID(portfolioID)
		if err != nil {
			t.Fatalf("GetByPortfolioID() error = %v", err)
		}

		if len(trades) != 1 {
			t.Errorf("GetByPortfolioID() returned %d trades, want 1", len(trades))
		}
	})

	t.Run("GetByOrderID", func(t *testing.T) {
		orderID := uuid.New()
		trade := &model.Trade{
			ID:          uuid.New(),
			PortfolioID: uuid.New(),
			OrderID:     orderID,
			Symbol:      "GOOGL",
			Side:        model.OrderSideSell,
			Quantity:    20,
			Price:       140.00,
			Total:       2800.00,
			ExecutedAt:  time.Now(),
		}

		repo.Create(trade)

		trades, err := repo.GetByOrderID(orderID)
		if err != nil {
			t.Fatalf("GetByOrderID() error = %v", err)
		}

		if len(trades) != 1 {
			t.Errorf("GetByOrderID() returned %d trades, want 1", len(trades))
		}
	})
}

func TestSeedDefaultPortfolio(t *testing.T) {
	portfolioRepo := NewInMemoryPortfolioRepository()
	positionRepo := NewInMemoryPositionRepository()

	portfolio, err := SeedDefaultPortfolio(portfolioRepo, positionRepo)
	if err != nil {
		t.Fatalf("SeedDefaultPortfolio() error = %v", err)
	}

	if portfolio == nil {
		t.Fatal("SeedDefaultPortfolio() returned nil portfolio")
	}

	if portfolio.Name != "Default Paper Portfolio" {
		t.Errorf("SeedDefaultPortfolio() name = %v, want 'Default Paper Portfolio'", portfolio.Name)
	}

	if portfolio.CashBalance != 100000 {
		t.Errorf("SeedDefaultPortfolio() cash_balance = %v, want 100000", portfolio.CashBalance)
	}

	positions, _ := positionRepo.GetByPortfolioID(portfolio.ID)
	if len(positions) != 2 {
		t.Errorf("SeedDefaultPortfolio() created %d positions, want 2", len(positions))
	}
}

export interface Portfolio {
  id: string;
  name: string;
  initialBalance: number;
  currentBalance: number;
  totalValue: number;
  cashBalance: number;
  positions: Position[];
  transactions: Transaction[];
  performance: PerformanceMetrics;
  createdAt: string;
  updatedAt: string;
}

export interface Position {
  id: string;
  symbol: string;
  name: string;
  quantity: number;
  avgCost: number;
  currentPrice: number;
  marketValue: number;
  unrealizedPL: number;
  unrealizedPLPercent: number;
  dayChange: number;
  dayChangePercent: number;
  weight: number;
  openedAt: string;
}

export interface Transaction {
  id: string;
  symbol: string;
  type: 'buy' | 'sell';
  quantity: number;
  price: number;
  total: number;
  fees: number;
  executedAt: string;
  notes?: string;
}

export interface TradeOrder {
  symbol: string;
  type: 'buy' | 'sell';
  orderType: 'market' | 'limit' | 'stop' | 'stop_limit';
  quantity: number;
  limitPrice?: number;
  stopPrice?: number;
}

export interface PerformanceMetrics {
  totalReturn: number;
  totalReturnPercent: number;
  dayReturn: number;
  dayReturnPercent: number;
  weekReturn: number;
  weekReturnPercent: number;
  monthReturn: number;
  monthReturnPercent: number;
  yearReturn: number;
  yearReturnPercent: number;
  sharpeRatio: number;
  sortinoRatio: number;
  maxDrawdown: number;
  winRate: number;
  avgWin: number;
  avgLoss: number;
  profitFactor: number;
}

export interface JournalEntry {
  id: string;
  transactionId: string;
  symbol: string;
  type: 'buy' | 'sell';
  quantity: number;
  price: number;
  reasoning: string;
  emotions: string[];
  lessons?: string;
  rating: number;
  createdAt: string;
}

export interface BacktestConfig {
  symbol: string;
  startDate: string;
  endDate: string;
  initialCapital: number;
  strategy: BacktestStrategy;
}

export interface BacktestStrategy {
  name: string;
  type: 'sma_crossover' | 'rsi' | 'macd' | 'custom';
  params: Record<string, number>;
}

export interface BacktestResult {
  id: string;
  config: BacktestConfig;
  metrics: PerformanceMetrics;
  trades: BacktestTrade[];
  equityCurve: { date: string; value: number }[];
  completedAt: string;
}

export interface BacktestTrade {
  entryDate: string;
  exitDate: string;
  entryPrice: number;
  exitPrice: number;
  quantity: number;
  profit: number;
  profitPercent: number;
  type: 'buy' | 'sell';
}

export interface LeaderboardEntry {
  rank: number;
  username: string;
  avatar?: string;
  totalReturn: number;
  totalReturnPercent: number;
  winRate: number;
  totalTrades: number;
  sharpeRatio: number;
  badge?: 'gold' | 'silver' | 'bronze';
}

export const APP_NAME = 'Super Dashboard';
export const APP_DESCRIPTION = 'Integrated Sports Betting Analytics & Stock Monitoring Platform';

export const ROUTES = {
  HOME: '/',
  LOGIN: '/auth/login',
  REGISTER: '/auth/register',
  FORGOT_PASSWORD: '/auth/forgot-password',
  DASHBOARD: '/dashboard',
  BETTING: '/dashboard/betting',
  VALUE_BETS: '/dashboard/betting/value-bets',
  BETTING_HISTORY: '/dashboard/betting/history',
  STOCKS: '/dashboard/stocks',
  WATCHLIST: '/dashboard/stocks/watchlist',
  SCREENER: '/dashboard/stocks/screener',
  PAPER_TRADING: '/dashboard/paper-trading',
  TRADE: '/dashboard/paper-trading/trade',
  JOURNAL: '/dashboard/paper-trading/journal',
  BACKTEST: '/dashboard/paper-trading/backtest',
  ANALYTICS: '/dashboard/analytics',
  GOALS: '/dashboard/analytics/goals',
  REPORTS: '/dashboard/analytics/reports',
  SETTINGS: '/dashboard/settings',
  SECURITY: '/dashboard/settings/security',
  NOTIFICATIONS: '/dashboard/settings/notifications',
} as const;

export const KEYBOARD_SHORTCUTS = {
  COMMAND_PALETTE: { key: 'k', ctrl: true, description: 'Open Command Palette' },
  BETTING: { key: 'b', description: 'Go to Betting' },
  STOCKS: { key: 's', description: 'Go to Stocks' },
  DASHBOARD: { key: 'd', description: 'Go to Dashboard' },
  PAPER_TRADING: { key: 'p', description: 'Go to Paper Trading' },
  ANALYTICS: { key: 'a', description: 'Go to Analytics' },
  SEARCH: { key: '/', description: 'Focus Search' },
  ESCAPE: { key: 'Escape', description: 'Close Modal/Palette' },
} as const;

export const LEAGUES = [
  { id: 'epl', name: 'Premier League', country: 'England' },
  { id: 'laliga', name: 'La Liga', country: 'Spain' },
  { id: 'bundesliga', name: 'Bundesliga', country: 'Germany' },
  { id: 'seriea', name: 'Serie A', country: 'Italy' },
  { id: 'ligue1', name: 'Ligue 1', country: 'France' },
  { id: 'ucl', name: 'Champions League', country: 'Europe' },
  { id: 'uel', name: 'Europa League', country: 'Europe' },
  { id: 'tpl', name: 'Thai Premier League', country: 'Thailand' },
] as const;

export const SECTORS = [
  'Technology',
  'Healthcare',
  'Finance',
  'Consumer Discretionary',
  'Consumer Staples',
  'Energy',
  'Materials',
  'Industrials',
  'Utilities',
  'Real Estate',
  'Communication Services',
] as const;

export const EXCHANGES = [
  { id: 'NYSE', name: 'New York Stock Exchange' },
  { id: 'NASDAQ', name: 'NASDAQ' },
  { id: 'SET', name: 'Stock Exchange of Thailand' },
] as const;

export const BET_TYPES = [
  { id: '1x2_home', name: 'Home Win (1)' },
  { id: '1x2_draw', name: 'Draw (X)' },
  { id: '1x2_away', name: 'Away Win (2)' },
  { id: 'over25', name: 'Over 2.5 Goals' },
  { id: 'under25', name: 'Under 2.5 Goals' },
  { id: 'btts_yes', name: 'Both Teams to Score - Yes' },
  { id: 'btts_no', name: 'Both Teams to Score - No' },
  { id: 'double_home', name: 'Home or Draw (1X)' },
  { id: 'double_away', name: 'Away or Draw (X2)' },
] as const;

export const EMOTIONS = [
  'Confident',
  'Anxious',
  'FOMO',
  'Greedy',
  'Fearful',
  'Calm',
  'Excited',
  'Frustrated',
  'Impatient',
  'Disciplined',
] as const;

export const COLORS = {
  primary: '#3b82f6',
  secondary: '#8b5cf6',
  success: '#10b981',
  warning: '#f59e0b',
  danger: '#ef4444',
  background: '#0a0a0f',
  surface: '#12121a',
} as const;

export const CHART_COLORS = [
  '#3b82f6',
  '#8b5cf6',
  '#10b981',
  '#f59e0b',
  '#ef4444',
  '#06b6d4',
  '#ec4899',
  '#84cc16',
  '#f97316',
  '#6366f1',
] as const;

export const API_ENDPOINTS = {
  AUTH: '/api/auth',
  BETTING: '/api/betting',
  STOCKS: '/api/stocks',
  PAPER_TRADING: '/api/paper-trading',
  ANALYTICS: '/api/analytics',
} as const;

export const PAGINATION = {
  DEFAULT_PAGE: 1,
  DEFAULT_LIMIT: 20,
  MAX_LIMIT: 100,
} as const;

export const DATE_FORMATS = {
  SHORT: 'MMM d',
  MEDIUM: 'MMM d, yyyy',
  LONG: 'MMMM d, yyyy',
  TIME: 'HH:mm',
  DATETIME: 'MMM d, yyyy HH:mm',
  ISO: "yyyy-MM-dd'T'HH:mm:ss",
} as const;

export const CURRENCY = {
  USD: { symbol: '$', code: 'USD', name: 'US Dollar' },
  THB: { symbol: '฿', code: 'THB', name: 'Thai Baht' },
  EUR: { symbol: '€', code: 'EUR', name: 'Euro' },
  GBP: { symbol: '£', code: 'GBP', name: 'British Pound' },
} as const;

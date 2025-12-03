export interface User {
  id: string;
  email: string;
  name: string;
  avatar?: string;
  role: 'user' | 'admin';
  preferences: UserPreferences;
  createdAt: string;
  updatedAt: string;
}

export interface UserPreferences {
  theme: 'dark' | 'light' | 'system';
  language: 'en' | 'th';
  currency: string;
  timezone: string;
  notifications: NotificationSettings;
}

export interface NotificationSettings {
  email: boolean;
  push: boolean;
  valueBets: boolean;
  priceAlerts: boolean;
  newsAlerts: boolean;
  portfolioUpdates: boolean;
}

export interface Goal {
  id: string;
  title: string;
  description?: string;
  type: 'profit' | 'roi' | 'win_rate' | 'custom';
  target: number;
  current: number;
  unit: string;
  deadline?: string;
  status: 'active' | 'completed' | 'failed';
  createdAt: string;
}

export interface Report {
  id: string;
  type: 'daily' | 'weekly' | 'monthly' | 'yearly';
  period: string;
  metrics: ReportMetrics;
  generatedAt: string;
}

export interface ReportMetrics {
  bettingProfit: number;
  bettingROI: number;
  tradingProfit: number;
  tradingROI: number;
  totalProfit: number;
  goalsCompleted: number;
  goalsTotal: number;
}

export interface Notification {
  id: string;
  type: 'info' | 'success' | 'warning' | 'error';
  title: string;
  message: string;
  read: boolean;
  createdAt: string;
  link?: string;
}

export interface KeyboardShortcut {
  key: string;
  ctrl?: boolean;
  alt?: boolean;
  shift?: boolean;
  description: string;
  action: () => void;
}

export interface CommandPaletteItem {
  id: string;
  title: string;
  description?: string;
  icon?: React.ReactNode;
  category: string;
  action: () => void;
  keywords?: string[];
}

export type Theme = 'dark' | 'light' | 'system';

export type Language = 'en' | 'th';

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

export interface ApiError {
  code: string;
  message: string;
  details?: Record<string, string[]>;
}

export interface SelectOption {
  value: string;
  label: string;
  disabled?: boolean;
}

export interface TabItem {
  id: string;
  label: string;
  icon?: React.ReactNode;
  disabled?: boolean;
}

export interface BreadcrumbItem {
  label: string;
  href?: string;
}

export interface ChartDataPoint {
  date: string;
  value: number;
  label?: string;
}

export interface HeatmapCell {
  x: string;
  y: string;
  value: number;
}

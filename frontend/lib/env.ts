/**
 * Environment configuration for the Super Dashboard frontend.
 * This module provides access to environment variables with proper typing
 * and default values.
 */

/**
 * Whether to use mock data instead of real API calls.
 * Set NEXT_PUBLIC_USE_MOCK_DATA=true in environment to enable mock data.
 * Defaults to true in development, false in production.
 */
export const USE_MOCK_DATA: boolean =
  process.env.NEXT_PUBLIC_USE_MOCK_DATA === 'true' ||
  (process.env.NEXT_PUBLIC_USE_MOCK_DATA === undefined && process.env.NODE_ENV === 'development');

/**
 * Backend API base URL.
 */
export const API_BASE_URL: string =
  process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080';

/**
 * Environment name (development, production, test).
 */
export const ENV: string = process.env.NODE_ENV || 'development';

/**
 * Check if we're in development mode.
 */
export const IS_DEV: boolean = ENV === 'development';

/**
 * Check if we're in production mode.
 */
export const IS_PROD: boolean = ENV === 'production';

/**
 * Environment configuration object for convenience.
 */
export const envConfig = {
  useMockData: USE_MOCK_DATA,
  apiBaseUrl: API_BASE_URL,
  env: ENV,
  isDev: IS_DEV,
  isProd: IS_PROD,
} as const;

export default envConfig;

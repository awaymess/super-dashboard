/**
 * Stock Valuation Calculations
 * Includes P/E, P/BV, Buffett Intrinsic Value, and Margin of Safety
 */

export interface PEValuationInputs {
  eps: number;
  fairPE: number;
  currentPrice: number;
  industryPE?: number;
  historicalPE?: number;
}

export interface PBVValuationInputs {
  bookValuePerShare: number;
  fairPBV: number;
  currentPrice: number;
  roe?: number;
}

export interface BuffettInputs {
  ownerEarnings: number; // Net income + depreciation - capex
  growthRate: number;
  discountRate: number;
  years: number;
  sharesOutstanding: number;
}

export interface MarginOfSafetyResult {
  intrinsicValue: number;
  currentPrice: number;
  marginOfSafety: number;
  upside: number;
  rating: 'undervalued' | 'fairValue' | 'overvalued';
  recommendation: 'strong_buy' | 'buy' | 'hold' | 'sell' | 'strong_sell';
}

/**
 * Calculate fair price using P/E valuation
 * @param inputs - P/E valuation inputs
 * @returns Fair price based on P/E
 */
export function calculatePEValuation(inputs: PEValuationInputs): {
  fairPrice: number;
  marginOfSafety: number;
  isPriceAttractive: boolean;
} {
  const { eps, fairPE, currentPrice, industryPE, historicalPE } = inputs;
  
  // Use provided fair P/E or calculate from industry/historical average
  let adjustedPE = fairPE;
  if (industryPE && historicalPE) {
    adjustedPE = (fairPE + industryPE + historicalPE) / 3;
  }
  
  const fairPrice = eps * adjustedPE;
  const marginOfSafety = fairPrice > 0 
    ? ((fairPrice - currentPrice) / fairPrice) * 100 
    : 0;
  
  return {
    fairPrice,
    marginOfSafety,
    isPriceAttractive: marginOfSafety > 20,
  };
}

/**
 * Calculate fair price using P/BV valuation (best for banks, real estate)
 * @param inputs - P/BV valuation inputs
 * @returns Fair price based on P/BV
 */
export function calculatePBVValuation(inputs: PBVValuationInputs): {
  fairPrice: number;
  marginOfSafety: number;
  isPriceAttractive: boolean;
  justifiedPBV?: number;
} {
  const { bookValuePerShare, fairPBV, currentPrice, roe } = inputs;
  
  const fairPrice = bookValuePerShare * fairPBV;
  const marginOfSafety = fairPrice > 0 
    ? ((fairPrice - currentPrice) / fairPrice) * 100 
    : 0;
  
  // Calculate justified P/BV if ROE is provided
  // Justified P/BV = (ROE - g) / (r - g), simplified to ROE / cost of equity
  // Using 10% cost of equity as a typical equity risk premium for stable companies
  // This is based on historical equity returns and should be adjusted for market conditions
  const COST_OF_EQUITY = 10; // 10% - typical cost of equity for stable companies
  const justifiedPBV = roe ? roe / COST_OF_EQUITY : undefined;
  
  return {
    fairPrice,
    marginOfSafety,
    isPriceAttractive: marginOfSafety > 20,
    justifiedPBV,
  };
}

/**
 * Calculate intrinsic value using Buffett's Owner Earnings method
 * @param inputs - Buffett valuation inputs
 * @returns Intrinsic value calculation
 */
export function calculateBuffettValue(inputs: BuffettInputs): {
  intrinsicValue: number;
  perShareValue: number;
  presentValueOfEarnings: number;
  terminalValue: number;
} {
  const { ownerEarnings, growthRate, discountRate, years, sharesOutstanding } = inputs;
  
  let totalPV = 0;
  let currentEarnings = ownerEarnings;
  
  // Calculate present value of owner earnings for projection period
  for (let year = 1; year <= years; year++) {
    currentEarnings *= (1 + growthRate / 100);
    const discountFactor = Math.pow(1 + discountRate / 100, year);
    totalPV += currentEarnings / discountFactor;
  }
  
  // Terminal value using perpetuity growth model
  // Using 2.5% terminal growth as a conservative estimate - this represents:
  // - Long-term GDP growth rate (typically 2-3%)
  // - Should never exceed the discount rate
  // - Represents sustainable growth into perpetuity
  const TERMINAL_GROWTH_RATE = 2.5;
  const terminalEarnings = currentEarnings * (1 + TERMINAL_GROWTH_RATE / 100);
  const terminalValue = terminalEarnings / ((discountRate - TERMINAL_GROWTH_RATE) / 100);
  const pvTerminal = terminalValue / Math.pow(1 + discountRate / 100, years);
  
  const intrinsicValue = totalPV + pvTerminal;
  
  return {
    intrinsicValue,
    perShareValue: intrinsicValue / sharesOutstanding,
    presentValueOfEarnings: totalPV,
    terminalValue: pvTerminal,
  };
}

/**
 * Calculate comprehensive margin of safety
 * @param intrinsicValue - Calculated intrinsic value
 * @param currentPrice - Current market price
 * @returns Margin of safety analysis
 */
export function calculateMarginOfSafety(
  intrinsicValue: number,
  currentPrice: number
): MarginOfSafetyResult {
  const marginOfSafety = intrinsicValue > 0 
    ? ((intrinsicValue - currentPrice) / intrinsicValue) * 100 
    : 0;
  
  const upside = currentPrice > 0 
    ? ((intrinsicValue - currentPrice) / currentPrice) * 100 
    : 0;
  
  let rating: 'undervalued' | 'fairValue' | 'overvalued';
  let recommendation: 'strong_buy' | 'buy' | 'hold' | 'sell' | 'strong_sell';
  
  if (marginOfSafety >= 50) {
    rating = 'undervalued';
    recommendation = 'strong_buy';
  } else if (marginOfSafety >= 30) {
    rating = 'undervalued';
    recommendation = 'buy';
  } else if (marginOfSafety >= -10) {
    rating = 'fairValue';
    recommendation = 'hold';
  } else if (marginOfSafety >= -30) {
    rating = 'overvalued';
    recommendation = 'sell';
  } else {
    rating = 'overvalued';
    recommendation = 'strong_sell';
  }
  
  return {
    intrinsicValue,
    currentPrice,
    marginOfSafety,
    upside,
    rating,
    recommendation,
  };
}

/**
 * Calculate DuPont Analysis (ROE breakdown)
 * ROE = Net Profit Margin × Asset Turnover × Equity Multiplier
 */
export function calculateDuPontAnalysis(
  netIncome: number,
  revenue: number,
  totalAssets: number,
  totalEquity: number
): {
  roe: number;
  netProfitMargin: number;
  assetTurnover: number;
  equityMultiplier: number;
  decomposition: string;
} {
  const netProfitMargin = revenue > 0 ? (netIncome / revenue) * 100 : 0;
  const assetTurnover = totalAssets > 0 ? revenue / totalAssets : 0;
  const equityMultiplier = totalEquity > 0 ? totalAssets / totalEquity : 0;
  const roe = (netProfitMargin / 100) * assetTurnover * equityMultiplier * 100;
  
  return {
    roe,
    netProfitMargin,
    assetTurnover,
    equityMultiplier,
    decomposition: `ROE (${roe.toFixed(2)}%) = Margin (${netProfitMargin.toFixed(2)}%) × Turnover (${assetTurnover.toFixed(2)}) × Leverage (${equityMultiplier.toFixed(2)})`,
  };
}

/**
 * Composite valuation using multiple methods
 */
export function calculateCompositeValuation(
  dcfValue: number,
  grahamValue: number,
  peValue: number,
  pbvValue: number,
  currentPrice: number,
  weights: { dcf: number; graham: number; pe: number; pbv: number } = { dcf: 0.3, graham: 0.25, pe: 0.25, pbv: 0.2 }
): MarginOfSafetyResult {
  const compositeValue = 
    dcfValue * weights.dcf +
    grahamValue * weights.graham +
    peValue * weights.pe +
    pbvValue * weights.pbv;
  
  return calculateMarginOfSafety(compositeValue, currentPrice);
}

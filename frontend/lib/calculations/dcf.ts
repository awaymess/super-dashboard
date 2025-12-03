export interface DCFInputs {
  freeCashFlow: number;
  growthRate: number;
  terminalGrowthRate: number;
  discountRate: number;
  years: number;
  sharesOutstanding: number;
}

export interface DCFResult {
  intrinsicValue: number;
  perShareValue: number;
  presentValueOfCashFlows: number;
  terminalValue: number;
  presentValueOfTerminal: number;
  projectedCashFlows: { year: number; cashFlow: number; presentValue: number }[];
}

export function calculateDCF(inputs: DCFInputs): DCFResult {
  const {
    freeCashFlow,
    growthRate,
    terminalGrowthRate,
    discountRate,
    years,
    sharesOutstanding,
  } = inputs;

  const projectedCashFlows: { year: number; cashFlow: number; presentValue: number }[] = [];
  let totalPVCashFlows = 0;
  let currentCashFlow = freeCashFlow;

  for (let year = 1; year <= years; year++) {
    currentCashFlow *= (1 + growthRate / 100);
    const discountFactor = Math.pow(1 + discountRate / 100, year);
    const presentValue = currentCashFlow / discountFactor;
    
    projectedCashFlows.push({
      year,
      cashFlow: currentCashFlow,
      presentValue,
    });
    
    totalPVCashFlows += presentValue;
  }

  const finalCashFlow = currentCashFlow * (1 + terminalGrowthRate / 100);
  const terminalValue = finalCashFlow / ((discountRate - terminalGrowthRate) / 100);
  const pvTerminalValue = terminalValue / Math.pow(1 + discountRate / 100, years);

  const intrinsicValue = totalPVCashFlows + pvTerminalValue;
  const perShareValue = intrinsicValue / sharesOutstanding;

  return {
    intrinsicValue,
    perShareValue,
    presentValueOfCashFlows: totalPVCashFlows,
    terminalValue,
    presentValueOfTerminal: pvTerminalValue,
    projectedCashFlows,
  };
}

export function calculateWACC(
  equityWeight: number,
  debtWeight: number,
  costOfEquity: number,
  costOfDebt: number,
  taxRate: number
): number {
  return (
    (equityWeight / 100) * (costOfEquity / 100) +
    (debtWeight / 100) * (costOfDebt / 100) * (1 - taxRate / 100)
  ) * 100;
}

export function calculateCostOfEquity(
  riskFreeRate: number,
  beta: number,
  marketReturn: number
): number {
  return riskFreeRate + beta * (marketReturn - riskFreeRate);
}

export function estimateGrowthRate(
  roe: number,
  retentionRatio: number
): number {
  return (roe / 100) * (retentionRatio / 100) * 100;
}

export function reverseDCF(
  currentPrice: number,
  sharesOutstanding: number,
  freeCashFlow: number,
  discountRate: number,
  terminalGrowthRate: number,
  years: number
): number {
  const targetValue = currentPrice * sharesOutstanding;
  
  let lowGrowth = 0;
  let highGrowth = 50;
  let midGrowth = 25;
  const tolerance = 0.01;
  let iterations = 0;
  const maxIterations = 100;

  while (iterations < maxIterations) {
    const result = calculateDCF({
      freeCashFlow,
      growthRate: midGrowth,
      terminalGrowthRate,
      discountRate,
      years,
      sharesOutstanding,
    });

    const diff = result.intrinsicValue - targetValue;
    
    if (Math.abs(diff) < targetValue * tolerance) {
      return midGrowth;
    }

    if (diff > 0) {
      highGrowth = midGrowth;
    } else {
      lowGrowth = midGrowth;
    }

    midGrowth = (lowGrowth + highGrowth) / 2;
    iterations++;
  }

  return midGrowth;
}

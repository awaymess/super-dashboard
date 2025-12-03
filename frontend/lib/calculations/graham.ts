export interface GrahamInputs {
  eps: number;
  bookValuePerShare: number;
  currentPrice: number;
  growthRate?: number;
  aaa_yield?: number;
}

export interface GrahamResult {
  grahamNumber: number;
  modifiedGrahamValue: number;
  peLimit: number;
  marginOfSafety: number;
  rating: 'undervalued' | 'fairValue' | 'overvalued';
  analysis: {
    peRatio: number;
    pbRatio: number;
    peg?: number;
    isDefensive: boolean;
    isEnterprising: boolean;
  };
}

export function calculateGrahamNumber(eps: number, bookValue: number): number {
  if (eps <= 0 || bookValue <= 0) return 0;
  return Math.sqrt(22.5 * eps * bookValue);
}

export function calculateModifiedGrahamValue(
  eps: number,
  growthRate: number = 0,
  aaaYield: number = 4.4
): number {
  if (eps <= 0) return 0;
  
  const baseMultiplier = 8.5;
  const growthMultiplier = 2;
  const expectedReturn = 4.4;
  
  const value = (eps * (baseMultiplier + growthMultiplier * growthRate) * expectedReturn) / aaaYield;
  return Math.max(0, value);
}

export function calculateGrahamAnalysis(inputs: GrahamInputs): GrahamResult {
  const { eps, bookValuePerShare, currentPrice, growthRate = 0, aaa_yield = 4.4 } = inputs;

  const grahamNumber = calculateGrahamNumber(eps, bookValuePerShare);
  const modifiedGrahamValue = calculateModifiedGrahamValue(eps, growthRate, aaa_yield);
  
  const peRatio = eps > 0 ? currentPrice / eps : 0;
  const pbRatio = bookValuePerShare > 0 ? currentPrice / bookValuePerShare : 0;
  const peg = growthRate > 0 && eps > 0 ? peRatio / growthRate : undefined;

  const peLimit = 15;
  const pbLimit = 1.5;
  const combinedLimit = 22.5;

  const isDefensive = 
    peRatio > 0 && peRatio <= peLimit &&
    pbRatio > 0 && pbRatio <= pbLimit &&
    (peRatio * pbRatio) <= combinedLimit;

  const isEnterprising = 
    peRatio > 0 && peRatio <= 20 &&
    pbRatio > 0 && pbRatio <= 2;

  const intrinsicValue = Math.max(grahamNumber, modifiedGrahamValue);
  const marginOfSafety = intrinsicValue > 0 
    ? ((intrinsicValue - currentPrice) / intrinsicValue) * 100 
    : 0;

  let rating: 'undervalued' | 'fairValue' | 'overvalued';
  if (marginOfSafety >= 30) {
    rating = 'undervalued';
  } else if (marginOfSafety >= -10) {
    rating = 'fairValue';
  } else {
    rating = 'overvalued';
  }

  return {
    grahamNumber,
    modifiedGrahamValue,
    peLimit: combinedLimit,
    marginOfSafety,
    rating,
    analysis: {
      peRatio,
      pbRatio,
      peg,
      isDefensive,
      isEnterprising,
    },
  };
}

export function screenDefensiveStocks(
  stocks: Array<{ symbol: string; eps: number; bookValue: number; price: number }>
): string[] {
  return stocks
    .filter(stock => {
      const pe = stock.eps > 0 ? stock.price / stock.eps : Infinity;
      const pb = stock.bookValue > 0 ? stock.price / stock.bookValue : Infinity;
      return pe <= 15 && pb <= 1.5 && (pe * pb) <= 22.5;
    })
    .map(stock => stock.symbol);
}

export function calculateNCAV(
  currentAssets: number,
  totalLiabilities: number,
  sharesOutstanding: number
): { ncav: number; ncavPerShare: number } {
  const ncav = currentAssets - totalLiabilities;
  const ncavPerShare = sharesOutstanding > 0 ? ncav / sharesOutstanding : 0;
  return { ncav, ncavPerShare };
}

export function isNetNet(
  currentAssets: number,
  totalLiabilities: number,
  marketCap: number
): boolean {
  const ncav = currentAssets - totalLiabilities;
  return marketCap < ncav * 0.67;
}

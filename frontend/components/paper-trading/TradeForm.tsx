'use client';

import { useState } from 'react';
import { GlassCard, GlassInput, GlassButton, GlassSelect, Tabs } from '@/components/ui';

interface TradeFormProps {
  symbol?: string;
  currentPrice?: number;
  onSubmit: (trade: {
    symbol: string;
    side: 'buy' | 'sell';
    orderType: 'market' | 'limit';
    quantity: number;
    price?: number;
  }) => void;
}

export function TradeForm({ symbol = '', currentPrice = 0, onSubmit }: TradeFormProps) {
  const [side, setSide] = useState<'buy' | 'sell'>('buy');
  const [orderType, setOrderType] = useState<'market' | 'limit'>('market');
  const [tickerSymbol, setTickerSymbol] = useState(symbol);
  const [quantity, setQuantity] = useState('1');
  const [limitPrice, setLimitPrice] = useState(currentPrice.toString());

  const quantityNum = parseInt(quantity) || 0;
  const priceNum = orderType === 'market' ? currentPrice : parseFloat(limitPrice) || 0;
  const total = quantityNum * priceNum;

  const handleSubmit = () => {
    onSubmit({
      symbol: tickerSymbol,
      side,
      orderType,
      quantity: quantityNum,
      price: orderType === 'limit' ? priceNum : undefined,
    });
  };

  return (
    <GlassCard>
      <h3 className="font-bold text-white mb-4">Place Order</h3>

      <Tabs
        tabs={[
          { id: 'buy', label: 'Buy' },
          { id: 'sell', label: 'Sell' },
        ]}
        activeTab={side}
        onChange={(id) => setSide(id as 'buy' | 'sell')}
        className="mb-4"
      />

      <div className="space-y-4">
        <div>
          <label className="text-sm text-gray-400 mb-2 block">Symbol</label>
          <GlassInput
            value={tickerSymbol}
            onChange={(e) => setTickerSymbol(e.target.value.toUpperCase())}
            placeholder="AAPL"
          />
        </div>

        <div>
          <label className="text-sm text-gray-400 mb-2 block">Order Type</label>
          <div className="grid grid-cols-2 gap-2">
            <GlassButton
              variant={orderType === 'market' ? 'primary' : 'ghost'}
              size="sm"
              onClick={() => setOrderType('market')}
            >
              Market
            </GlassButton>
            <GlassButton
              variant={orderType === 'limit' ? 'primary' : 'ghost'}
              size="sm"
              onClick={() => setOrderType('limit')}
            >
              Limit
            </GlassButton>
          </div>
        </div>

        <div>
          <label className="text-sm text-gray-400 mb-2 block">Quantity</label>
          <GlassInput
            type="number"
            value={quantity}
            onChange={(e) => setQuantity(e.target.value)}
            min={1}
          />
        </div>

        {orderType === 'limit' && (
          <div>
            <label className="text-sm text-gray-400 mb-2 block">Limit Price</label>
            <GlassInput
              type="number"
              value={limitPrice}
              onChange={(e) => setLimitPrice(e.target.value)}
              step={0.01}
            />
          </div>
        )}

        <div className="p-4 bg-white/5 rounded-lg">
          <div className="flex justify-between mb-2">
            <span className="text-gray-400">Est. Price</span>
            <span className="text-white">${priceNum.toFixed(2)}</span>
          </div>
          <div className="flex justify-between mb-2">
            <span className="text-gray-400">Quantity</span>
            <span className="text-white">{quantityNum}</span>
          </div>
          <div className="flex justify-between pt-2 border-t border-white/10">
            <span className="text-gray-300 font-medium">Total</span>
            <span className="text-white font-bold">${total.toFixed(2)}</span>
          </div>
        </div>

        <GlassButton
          variant={side === 'buy' ? 'success' : 'danger'}
          className="w-full"
          onClick={handleSubmit}
          disabled={!tickerSymbol || quantityNum <= 0}
        >
          {side === 'buy' ? 'Buy' : 'Sell'} {tickerSymbol || 'Stock'}
        </GlassButton>
      </div>
    </GlassCard>
  );
}

'use client';

import { Clock, ExternalLink } from 'lucide-react';
import { GlassCard, Badge } from '@/components/ui';
import type { StockNews } from '@/types/stocks';

interface NewsCardProps {
  news: StockNews;
  onClick?: () => void;
}

export function NewsCard({ news, onClick }: NewsCardProps) {
  const sentimentColors = {
    positive: 'success',
    negative: 'danger',
    neutral: 'default',
  } as const;

  return (
    <GlassCard className="cursor-pointer hover:border-primary/50 transition-colors" onClick={onClick}>
      <div className="flex items-start gap-4">
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-2">
            <Badge variant={sentimentColors[news.sentiment]}>{news.sentiment}</Badge>
            <span className="text-xs text-gray-400">{news.source}</span>
          </div>
          <h4 className="font-semibold text-white mb-2 line-clamp-2">{news.title}</h4>
          <p className="text-sm text-gray-400 line-clamp-2 mb-3">{news.summary}</p>
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-1 text-xs text-gray-500">
              <Clock className="w-3 h-3" />
              {news.publishedAt}
            </div>
            <div className="flex items-center gap-2">
              {news.symbols?.map(symbol => (
                <span key={symbol} className="text-xs text-primary bg-primary/10 px-2 py-0.5 rounded">
                  {symbol}
                </span>
              ))}
            </div>
          </div>
        </div>
      </div>
    </GlassCard>
  );
}

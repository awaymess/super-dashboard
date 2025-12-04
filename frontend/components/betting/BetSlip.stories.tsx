import type { Meta, StoryObj } from '@storybook/react';
import { useState } from 'react';
import { action } from '@storybook/addon-actions';
import { BetSlip } from './BetSlip';
import type { Bet, Match, Team } from '@/types/betting';

const sampleTeam = (overrides: Partial<Team> = {}): Team => ({
  id: 'team-a',
  name: 'Arsenal',
  shortName: 'ARS',
  form: ['W', 'W', 'D', 'L', 'W'],
  homeForm: ['W', 'D', 'W'],
  awayForm: ['L', 'W', 'W'],
  goalsScoredAvg: 2.1,
  goalsConcededAvg: 0.9,
  cleanSheetPct: 45,
  ...overrides,
});

const match: Match = {
  id: 'match-1',
  homeTeam: sampleTeam(),
  awayTeam: sampleTeam({ id: 'team-b', name: 'Chelsea', shortName: 'CHE' }),
  league: 'Premier League',
  leagueCountry: 'England',
  date: '2025-12-12',
  time: '20:00',
  status: 'scheduled',
  odds: {
    home: 1.95,
    draw: 3.6,
    away: 4.1,
    over25: 1.8,
    under25: 2.05,
    btts: 1.75,
    bttsNo: 2.12,
    homeDouble: 1.34,
    awayDouble: 1.92,
  },
};

const baseBets: Bet[] = [
  {
    id: 'bet-1',
    matchId: match.id,
    match,
    betType: 'Home Win',
    odds: 1.95,
    stake: 25,
    potentialWin: 48.75,
    status: 'pending',
    placedAt: '2025-12-01T18:02:00Z',
  },
  {
    id: 'bet-2',
    matchId: match.id,
    match,
    betType: 'Over 2.5 Goals',
    odds: 1.80,
    stake: 15,
    potentialWin: 27,
    status: 'pending',
    placedAt: '2025-12-01T18:05:00Z',
  },
];

const meta = {
  title: 'Betting/BetSlip',
  component: BetSlip,
  parameters: {
    layout: 'centered',
    backgrounds: { default: 'dark' },
  },
  args: {
    bets: baseBets,
    onRemoveBet: action('remove-bet'),
    onPlaceBet: action('place-bet'),
    onClear: action('clear-slip'),
  },
} satisfies Meta<typeof BetSlip>;

export default meta;
type Story = StoryObj<typeof meta>;

function BetSlipPreview(args: React.ComponentProps<typeof BetSlip>) {
  const [bets, setBets] = useState(args.bets);

  return (
    <div className="min-h-screen bg-gradient-to-br from-[#05060b] via-[#0b1120] to-[#05060b] flex items-start justify-center p-10">
      <BetSlip
        {...args}
        bets={bets}
        onRemoveBet={id => {
          setBets(prev => prev.filter(bet => bet.id !== id));
          args.onRemoveBet?.(id);
        }}
        onClear={() => {
          setBets([]);
          args.onClear?.();
        }}
        onPlaceBet={stake => {
          args.onPlaceBet?.(stake);
        }}
      />
    </div>
  );
}

export const Default: Story = {
  render: args => <BetSlipPreview {...args} />,
};

export const EmptySlip: Story = {
  args: {
    bets: [],
  },
  render: args => <BetSlipPreview {...args} />,
};

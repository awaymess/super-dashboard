import type { Meta, StoryObj } from '@storybook/react';
import { action } from '@storybook/addon-actions';
import { ValueBetCard } from './ValueBetCard';
import type { Match, Team, ValueBet } from '@/types/betting';

const team = (overrides: Partial<Team> = {}): Team => ({
  id: overrides.id ?? 'team-x',
  name: overrides.name ?? 'Arsenal',
  shortName: overrides.shortName ?? 'ARS',
  form: overrides.form ?? ['W', 'W', 'D', 'L', 'W'],
  goalsScoredAvg: overrides.goalsScoredAvg ?? 2.0,
  goalsConcededAvg: overrides.goalsConcededAvg ?? 0.8,
  cleanSheetPct: overrides.cleanSheetPct ?? 48,
  homeForm: overrides.homeForm ?? ['W', 'W', 'D'],
  awayForm: overrides.awayForm ?? ['L', 'W', 'W'],
});

const match: Match = {
  id: 'match-vb-1',
  homeTeam: team(),
  awayTeam: team({ id: 'team-y', name: 'Tottenham', shortName: 'TOT' }),
  league: 'Premier League',
  leagueCountry: 'England',
  date: '2025-12-15',
  time: '21:00',
  status: 'scheduled',
  odds: {
    home: 2.35,
    draw: 3.4,
    away: 3.0,
    over25: 1.95,
    under25: 1.92,
    btts: 1.7,
    bttsNo: 2.2,
    homeDouble: 1.52,
    awayDouble: 1.55,
  },
};

const valueBet: ValueBet = {
  matchId: match.id,
  match,
  betType: 'Home Win',
  bookmakerOdds: 2.35,
  fairOdds: 2.10,
  value: 11.9,
  confidence: 82,
  kellyStake: 4.7,
  expectedValue: 0.23,
};

const meta = {
  title: 'Betting/ValueBetCard',
  component: ValueBetCard,
  parameters: {
    layout: 'centered',
  },
  args: {
    valueBet,
    onClick: action('open-value-bet'),
  },
} satisfies Meta<typeof ValueBetCard>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const LowerEdge: Story = {
  args: {
    valueBet: {
      ...valueBet,
      value: 4.1,
      bookmakerOdds: 3.4,
      fairOdds: 3.2,
      kellyStake: 2.2,
      confidence: 64,
    },
  },
};

import type { Meta, StoryObj } from '@storybook/react';
import { action } from '@storybook/addon-actions';
import { MatchCard } from './MatchCard';
import type { Match, Team } from '@/types/betting';

const team = (overrides: Partial<Team> = {}): Team => ({
  id: `team-${overrides.shortName ?? 'x'}`,
  name: overrides.name ?? 'Team X',
  shortName: overrides.shortName ?? 'TX',
  form: overrides.form ?? ['W', 'D', 'W', 'L', 'W'],
  goalsScoredAvg: overrides.goalsScoredAvg ?? 1.9,
  goalsConcededAvg: overrides.goalsConcededAvg ?? 1.1,
  cleanSheetPct: overrides.cleanSheetPct ?? 38,
  homeForm: overrides.homeForm ?? ['W', 'W', 'D'],
  awayForm: overrides.awayForm ?? ['L', 'W', 'D'],
});

const baseMatch: Match = {
  id: 'match-001',
  homeTeam: team({ name: 'Arsenal', shortName: 'ARS' }),
  awayTeam: team({ name: 'Chelsea', shortName: 'CHE' }),
  league: 'Premier League',
  leagueCountry: 'England',
  date: '2025-12-12',
  time: '20:00',
  status: 'scheduled',
  odds: {
    home: 1.9,
    draw: 3.5,
    away: 4.2,
    over25: 1.82,
    under25: 2.04,
    btts: 1.72,
    bttsNo: 2.2,
    homeDouble: 1.33,
    awayDouble: 1.95,
  },
};

const meta = {
  title: 'Betting/MatchCard',
  component: MatchCard,
  parameters: {
    layout: 'centered',
  },
  args: {
    match: baseMatch,
    onClick: action('select-match'),
  },
} satisfies Meta<typeof MatchCard>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Upcoming: Story = {};

export const Finished: Story = {
  args: {
    match: {
      ...baseMatch,
      status: 'finished',
      homeScore: 2,
      awayScore: 1,
      time: 'FT',
    },
  },
};

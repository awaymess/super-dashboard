import type { Meta, StoryObj } from '@storybook/react';
import { H2HHistory } from './H2HHistory';

const sampleMatches = [
  { date: '2025-10-12', homeTeam: 'Arsenal', awayTeam: 'Chelsea', homeScore: 2, awayScore: 0, competition: 'Premier League' },
  { date: '2025-04-03', homeTeam: 'Chelsea', awayTeam: 'Arsenal', homeScore: 1, awayScore: 1, competition: 'Premier League' },
  { date: '2024-12-29', homeTeam: 'Arsenal', awayTeam: 'Chelsea', homeScore: 3, awayScore: 2, competition: 'Carabao Cup' },
  { date: '2024-05-08', homeTeam: 'Chelsea', awayTeam: 'Arsenal', homeScore: 0, awayScore: 1, competition: 'Premier League' },
  { date: '2023-11-18', homeTeam: 'Arsenal', awayTeam: 'Chelsea', homeScore: 2, awayScore: 2, competition: 'Premier League' },
];

const meta = {
  title: 'Betting/H2HHistory',
  component: H2HHistory,
  parameters: {
    layout: 'padded',
  },
  args: {
    homeTeam: 'Arsenal',
    awayTeam: 'Chelsea',
    matches: sampleMatches,
  },
} satisfies Meta<typeof H2HHistory>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const AwayDominance: Story = {
  args: {
    matches: sampleMatches.map(match => ({ ...match, homeScore: match.awayScore + 1 })),
  },
};

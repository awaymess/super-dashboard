import type { Meta, StoryObj } from '@storybook/react';
import { TeamForm } from './TeamForm';

const meta = {
  title: 'Betting/TeamForm',
  component: TeamForm,
  parameters: {
    layout: 'padded',
  },
  args: {
    teamName: 'Arsenal',
    form: ['W', 'D', 'W', 'W', 'L', 'W'],
    stats: {
      goalsScored: 14,
      goalsConceded: 5,
      cleanSheets: 3,
      failedToScore: 1,
    },
  },
} satisfies Meta<typeof TeamForm>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const PoorForm: Story = {
  args: {
    teamName: 'Everton',
    form: ['L', 'L', 'D', 'L', 'W'],
    stats: {
      goalsScored: 4,
      goalsConceded: 13,
      cleanSheets: 0,
      failedToScore: 3,
    },
  },
};

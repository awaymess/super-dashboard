import type { Meta, StoryObj } from '@storybook/react';
import { useState } from 'react';
import { action } from '@storybook/addon-actions';
import { Notifications } from './Notifications';

const meta = {
  title: 'Common/Notifications',
  component: Notifications,
  parameters: {
    layout: 'fullscreen',
  },
  args: {
    notifications: [],
    onDismiss: action('dismiss-notification'),
  },
} satisfies Meta<typeof Notifications>;

export default meta;
type Story = StoryObj<typeof meta>;

const demoNotifications = [
  {
    id: 'success-1',
    type: 'success' as const,
    title: 'Trade Executed',
    message: 'Bought 50 shares of AAPL at $188.42',
  },
  {
    id: 'warning-1',
    type: 'warning' as const,
    title: 'High Volatility Alert',
    message: 'BTC/USD volatility exceeded 5% in the last hour',
  },
  {
    id: 'error-1',
    type: 'error' as const,
    title: 'Order Failed',
    message: 'Insufficient margin on account',
  },
  {
    id: 'info-1',
    type: 'info' as const,
    title: 'System Update',
    message: 'Paper trading P&L recalculated at 02:00 UTC',
  },
];

function NotificationsPreview(args: React.ComponentProps<typeof Notifications>) {
  const [items, setItems] = useState(demoNotifications);

  const handleDismiss = (id: string) => {
    setItems(prev => prev.filter(notification => notification.id !== id));
    args.onDismiss?.(id);
  };

  const handleAddNotification = () => {
    const templates = [
      { type: 'success', title: 'Value Bet Locked', message: 'Chelsea vs Arsenal @ 2.10' },
      { type: 'warning', title: 'Portfolio Drift', message: 'Equities > 70% allocation' },
      { type: 'error', title: 'API Rate Limit', message: 'Retrying submission in 30s' },
      { type: 'info', title: 'Insight Ready', message: 'New analyst note for NVDA' },
    ] as const;
    const next = templates[Math.floor(Math.random() * templates.length)];
    const id = `${next.type}-${Date.now()}`;
    setItems(prev => [...prev, { id, ...next }]);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-[#05060b] via-[#0b1120] to-[#05060b] text-white">
      <div className="absolute top-10 left-10 space-y-3">
        <button
          type="button"
          onClick={handleAddNotification}
          className="px-4 py-2 rounded-xl bg-white/10 hover:bg-white/20 transition"
        >
          Push random notification
        </button>
        <button
          type="button"
          onClick={() => setItems(demoNotifications)}
          className="px-4 py-2 rounded-xl bg-white/5 hover:bg-white/15 transition text-white/80"
        >
          Reset list
        </button>
      </div>
      <Notifications {...args} notifications={items} onDismiss={handleDismiss} />
    </div>
  );
}

export const Default: Story = {
  render: args => <NotificationsPreview {...args} />,
};

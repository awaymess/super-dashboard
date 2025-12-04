import type { Meta, StoryObj } from '@storybook/react';
import { Table } from './Table';
import { Badge } from './Badge';

const meta: Meta<typeof Table> = {
  title: 'UI/Table',
  component: Table,
  parameters: {
    layout: 'padded',
  },
  tags: ['autodocs'],
};

export default meta;
type Story = StoryObj<typeof Table>;

interface Stock {
  id: string;
  symbol: string;
  name: string;
  price: number;
  change: number;
  changePercent: number;
  volume: number;
}

interface User {
  id: string;
  name: string;
  email: string;
  role: 'admin' | 'user' | 'editor';
  status: 'active' | 'inactive' | 'pending';
}

const sampleStocks: Stock[] = [
  { id: '1', symbol: 'AAPL', name: 'Apple Inc.', price: 178.72, change: 2.45, changePercent: 1.39, volume: 52341000 },
  { id: '2', symbol: 'GOOGL', name: 'Alphabet Inc.', price: 141.80, change: -1.23, changePercent: -0.86, volume: 21456000 },
  { id: '3', symbol: 'MSFT', name: 'Microsoft Corp.', price: 378.91, change: 5.67, changePercent: 1.52, volume: 18234000 },
  { id: '4', symbol: 'AMZN', name: 'Amazon.com Inc.', price: 178.25, change: -0.89, changePercent: -0.50, volume: 34567000 },
  { id: '5', symbol: 'TSLA', name: 'Tesla Inc.', price: 248.50, change: 8.90, changePercent: 3.71, volume: 98765000 },
];

const sampleUsers: User[] = [
  { id: '1', name: 'John Doe', email: 'john@example.com', role: 'admin', status: 'active' },
  { id: '2', name: 'Jane Smith', email: 'jane@example.com', role: 'editor', status: 'active' },
  { id: '3', name: 'Bob Johnson', email: 'bob@example.com', role: 'user', status: 'pending' },
  { id: '4', name: 'Alice Brown', email: 'alice@example.com', role: 'user', status: 'inactive' },
  { id: '5', name: 'Charlie Wilson', email: 'charlie@example.com', role: 'editor', status: 'active' },
];

export const Default: Story = {
  render: () => (
    <Table<Stock>
      columns={[
        { key: 'symbol', header: 'Symbol', sortable: true },
        { key: 'name', header: 'Name' },
        { key: 'price', header: 'Price', sortable: true, render: (item) => `$${item.price.toFixed(2)}` },
        { 
          key: 'change', 
          header: 'Change',
          render: (item) => (
            <span className={item.change >= 0 ? 'text-green-400' : 'text-red-400'}>
              {item.change >= 0 ? '+' : ''}{item.change.toFixed(2)} ({item.changePercent.toFixed(2)}%)
            </span>
          ),
        },
        { key: 'volume', header: 'Volume', render: (item) => item.volume.toLocaleString() },
      ]}
      data={sampleStocks}
      keyExtractor={(item) => item.id}
    />
  ),
};

export const WithRowClick: Story = {
  render: () => (
    <Table<Stock>
      columns={[
        { key: 'symbol', header: 'Symbol' },
        { key: 'name', header: 'Name' },
        { key: 'price', header: 'Price', render: (item) => `$${item.price.toFixed(2)}` },
      ]}
      data={sampleStocks}
      keyExtractor={(item) => item.id}
      onRowClick={(item) => alert(`Clicked on ${item.symbol}`)}
    />
  ),
};

export const Loading: Story = {
  render: () => (
    <Table<Stock>
      columns={[
        { key: 'symbol', header: 'Symbol' },
        { key: 'name', header: 'Name' },
        { key: 'price', header: 'Price' },
      ]}
      data={[]}
      keyExtractor={(item) => item.id}
      isLoading={true}
    />
  ),
};

export const Empty: Story = {
  render: () => (
    <Table<Stock>
      columns={[
        { key: 'symbol', header: 'Symbol' },
        { key: 'name', header: 'Name' },
        { key: 'price', header: 'Price' },
      ]}
      data={[]}
      keyExtractor={(item) => item.id}
      emptyMessage="No stocks found. Try adjusting your filters."
    />
  ),
};

export const WithBadges: Story = {
  render: () => (
    <Table<User>
      columns={[
        { key: 'name', header: 'Name', sortable: true },
        { key: 'email', header: 'Email' },
        { 
          key: 'role', 
          header: 'Role',
          render: (item) => (
            <Badge variant={item.role === 'admin' ? 'danger' : item.role === 'editor' ? 'warning' : 'default'}>
              {item.role}
            </Badge>
          ),
        },
        { 
          key: 'status', 
          header: 'Status',
          render: (item) => (
            <Badge variant={item.status === 'active' ? 'success' : item.status === 'pending' ? 'warning' : 'danger'}>
              {item.status}
            </Badge>
          ),
        },
      ]}
      data={sampleUsers}
      keyExtractor={(item) => item.id}
    />
  ),
};

export const Sortable: Story = {
  render: () => (
    <Table<Stock>
      columns={[
        { key: 'symbol', header: 'Symbol', sortable: true },
        { key: 'name', header: 'Name' },
        { key: 'price', header: 'Price', sortable: true, render: (item) => `$${item.price.toFixed(2)}` },
      ]}
      data={sampleStocks}
      keyExtractor={(item) => item.id}
      sortColumn="price"
      sortDirection="desc"
      onSort={(column) => alert(`Sorting by ${column}`)}
    />
  ),
};

export const CustomStyling: Story = {
  render: () => (
    <Table<Stock>
      columns={[
        { key: 'symbol', header: 'Symbol' },
        { key: 'name', header: 'Name' },
        { key: 'price', header: 'Price', render: (item) => `$${item.price.toFixed(2)}` },
      ]}
      data={sampleStocks}
      keyExtractor={(item) => item.id}
      className="border-primary/30"
    />
  ),
};

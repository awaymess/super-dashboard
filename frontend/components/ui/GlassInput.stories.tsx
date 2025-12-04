import type { Meta, StoryObj } from '@storybook/react';
import { Search, Eye, EyeOff, Mail, Lock, User } from 'lucide-react';
import { useState } from 'react';
import { GlassInput } from './GlassInput';

const meta: Meta<typeof GlassInput> = {
  title: 'UI/GlassInput',
  component: GlassInput,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
  decorators: [
    (Story) => (
      <div style={{ width: '320px' }}>
        <Story />
      </div>
    ),
  ],
  argTypes: {
    size: {
      control: 'select',
      options: ['sm', 'md', 'lg'],
    },
    iconPosition: {
      control: 'select',
      options: ['left', 'right'],
    },
    disabled: {
      control: 'boolean',
    },
  },
};

export default meta;
type Story = StoryObj<typeof GlassInput>;

export const Default: Story = {
  args: {
    placeholder: 'Enter your text...',
  },
};

export const WithLabel: Story = {
  args: {
    label: 'Email Address',
    placeholder: 'you@example.com',
    type: 'email',
  },
};

export const WithHint: Story = {
  args: {
    label: 'Username',
    placeholder: 'Choose a username',
    hint: 'Must be at least 3 characters long',
  },
};

export const WithError: Story = {
  args: {
    label: 'Password',
    placeholder: 'Enter password',
    type: 'password',
    error: 'Password must be at least 8 characters',
  },
};

export const WithIconLeft: Story = {
  args: {
    placeholder: 'Search...',
    icon: <Search className="w-4 h-4" />,
    iconPosition: 'left',
  },
};

export const WithIconRight: Story = {
  args: {
    label: 'Email',
    placeholder: 'you@example.com',
    icon: <Mail className="w-4 h-4" />,
    iconPosition: 'right',
  },
};

export const Small: Story = {
  args: {
    placeholder: 'Small input',
    size: 'sm',
  },
};

export const Large: Story = {
  args: {
    placeholder: 'Large input',
    size: 'lg',
  },
};

export const Disabled: Story = {
  args: {
    placeholder: 'Disabled input',
    disabled: true,
    value: 'Cannot edit this',
  },
};

const PasswordInputComponent = () => {
  const [showPassword, setShowPassword] = useState(false);
  
  return (
    <GlassInput
      label="Password"
      placeholder="Enter password"
      type={showPassword ? 'text' : 'password'}
      icon={showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
      iconPosition="right"
      onIconClick={() => setShowPassword(!showPassword)}
    />
  );
};

export const PasswordWithToggle: Story = {
  render: () => <PasswordInputComponent />,
};

export const LoginForm: Story = {
  render: () => (
    <div className="space-y-4">
      <GlassInput
        label="Email"
        placeholder="you@example.com"
        type="email"
        icon={<Mail className="w-4 h-4" />}
        iconPosition="left"
      />
      <GlassInput
        label="Password"
        placeholder="Enter password"
        type="password"
        icon={<Lock className="w-4 h-4" />}
        iconPosition="left"
      />
    </div>
  ),
};

export const AllSizes: Story = {
  render: () => (
    <div className="space-y-4">
      <GlassInput placeholder="Small input" size="sm" />
      <GlassInput placeholder="Medium input" size="md" />
      <GlassInput placeholder="Large input" size="lg" />
    </div>
  ),
};

import type { Metadata } from 'next';
import './globals.css';
import '../styles/animations.css';
import { StoreProvider } from '@/store/provider';

export const metadata: Metadata = {
  title: 'Super Dashboard - Sports Betting & Stock Analytics',
  description: 'Integrated Sports Betting Analytics & Stock Monitoring Platform',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="dark">
      <body className="font-sans bg-background text-white antialiased">
        <StoreProvider>
          {children}
        </StoreProvider>
      </body>
    </html>
  );
}

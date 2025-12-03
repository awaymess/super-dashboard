'use client';

import { Inter } from 'next/font/google';
import './globals.css';
import '../styles/animations.css';
import { ReduxProvider } from '@/store/provider';

const inter = Inter({ subsets: ['latin'] });

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="dark">
      <body className={`${inter.className} bg-background text-foreground antialiased`}>
        <ReduxProvider>
          <div className="min-h-screen">
            {children}
          </div>
        </ReduxProvider>
      </body>
    </html>
  );
}

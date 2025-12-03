import createMiddleware from 'next-intl/middleware';
import { NextRequest, NextResponse } from 'next/server';
import { locales, defaultLocale } from './i18n/config';

const intlMiddleware = createMiddleware({
  locales,
  defaultLocale,
  localePrefix: 'as-needed',
});

export default function middleware(req: NextRequest) {
  const { pathname } = req.nextUrl;

  // Redirect root path to dashboard
  if (pathname === '/') {
    const url = req.nextUrl.clone();
    url.pathname = '/dashboard';
    return NextResponse.redirect(url);
  }

  // Otherwise, run next-intl middleware
  return intlMiddleware(req);
}

export const config = {
  // Run middleware only on the root and locale-prefixed paths
  matcher: ['/', '/(en|th)/:path*'],
};

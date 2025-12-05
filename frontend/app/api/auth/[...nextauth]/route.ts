/**
 * NextAuth.js API Route Handler
 * 
 * This is a scaffold for NextAuth authentication.
 * Configure providers and callbacks as needed.
 * 
 * TODO: Implement proper auth handlers with your OAuth providers.
 * TODO: Add session callbacks for custom session data.
 * TODO: Add JWT callbacks for token management.
 */

import { NextRequest, NextResponse } from 'next/server';

/**
 * Scaffold handler for NextAuth API route.
 * 
 * To enable NextAuth:
 * 1. Install next-auth: npm install next-auth
 * 2. Uncomment the NextAuth configuration below
 * 3. Configure your OAuth providers in .env.local
 * 
 * Example configuration:
 * 
 * import NextAuth from 'next-auth';
 * import GoogleProvider from 'next-auth/providers/google';
 * import GitHubProvider from 'next-auth/providers/github';
 * import CredentialsProvider from 'next-auth/providers/credentials';
 * 
 * const handler = NextAuth({
 *   providers: [
 *     // Credentials provider for email/password login
 *     CredentialsProvider({
 *       name: 'Credentials',
 *       credentials: {
 *         email: { label: 'Email', type: 'email' },
 *         password: { label: 'Password', type: 'password' },
 *       },
 *       async authorize(credentials) {
 *         // TODO: Implement actual authentication logic
 *         // Call your backend API to verify credentials
 *         // const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/login`, {
 *         //   method: 'POST',
 *         //   body: JSON.stringify(credentials),
 *         //   headers: { 'Content-Type': 'application/json' },
 *         // });
 *         // const user = await res.json();
 *         // if (res.ok && user) return user;
 *         // return null;
 *         
 *         // Placeholder: return mock user for development
 *         if (credentials?.email && credentials?.password) {
 *           return {
 *             id: '1',
 *             email: credentials.email,
 *             name: 'Demo User',
 *           };
 *         }
 *         return null;
 *       },
 *     }),
 *     
 *     // Google OAuth provider
 *     // Uncomment and configure in .env.local:
 *     // GOOGLE_CLIENT_ID=your-client-id
 *     // GOOGLE_CLIENT_SECRET=your-client-secret
 *     // GoogleProvider({
 *     //   clientId: process.env.GOOGLE_CLIENT_ID!,
 *     //   clientSecret: process.env.GOOGLE_CLIENT_SECRET!,
 *     // }),
 *     
 *     // GitHub OAuth provider
 *     // Uncomment and configure in .env.local:
 *     // GITHUB_CLIENT_ID=your-client-id
 *     // GITHUB_CLIENT_SECRET=your-client-secret
 *     // GitHubProvider({
 *     //   clientId: process.env.GITHUB_CLIENT_ID!,
 *     //   clientSecret: process.env.GITHUB_CLIENT_SECRET!,
 *     // }),
 *   ],
 *   
 *   pages: {
 *     signIn: '/login',
 *     signOut: '/login',
 *     error: '/login',
 *   },
 *   
 *   callbacks: {
 *     async jwt({ token, user }) {
 *       // Add custom claims to the JWT token
 *       if (user) {
 *         token.id = user.id;
 *       }
 *       return token;
 *     },
 *     async session({ session, token }) {
 *       // Add custom properties to the session
 *       if (session.user) {
 *         (session.user as any).id = token.id;
 *       }
 *       return session;
 *     },
 *   },
 *   
 *   session: {
 *     strategy: 'jwt',
 *   },
 * });
 * 
 * export { handler as GET, handler as POST };
 */

// Placeholder handler until NextAuth is configured
// This returns a helpful message indicating NextAuth needs to be set up
export async function GET(request: NextRequest) {
  return NextResponse.json(
    {
      message: 'NextAuth not configured',
      instructions: 'Install next-auth and configure providers in this file',
    },
    { status: 501 }
  );
}

export async function POST(request: NextRequest) {
  return NextResponse.json(
    {
      message: 'NextAuth not configured',
      instructions: 'Install next-auth and configure providers in this file',
    },
    { status: 501 }
  );
}

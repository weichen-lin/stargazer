import GitHubProvider from 'next-auth/providers/github'
import type { AuthOptions } from 'next-auth'
import { Neo4jAdapter, conn } from '@/actions/adapter'

export const options: AuthOptions = {
  providers: [
    GitHubProvider({
      clientId: process.env.GITHUB_CLIENT_ID as string,
      clientSecret: process.env.GITHUB_CLIENT_SECRET as string,
    }),
  ],
  session: {
    strategy: 'jwt',
  },
  jwt: {
    secret: process.env.NEXTAUTH_SECRET,
  },
  pages: {
    signIn: '/',
  },
  secret: process.env.NEXTAUTH_SECRET,
  adapter: Neo4jAdapter(conn) as any,
}

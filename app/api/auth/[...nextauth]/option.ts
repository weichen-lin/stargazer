import GitHubProvider from 'next-auth/providers/github'
import type { AuthOptions } from 'next-auth'
import { driver, auth } from 'neo4j-driver'
import { Neo4jAdapter } from '@/actions/adapter'

const conn = driver(process.env.NEO4J_URL as string, auth.basic('neo4j', process.env.NEO4J_PASSWORD as string))

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

import { GetUser } from '@/actions'
import RepositoryClient from '@/client/repository-client'

export const dynamic = 'force-dynamic' // defaults to auto
export async function GET() {
  const { email } = await GetUser()
  const client = new RepositoryClient(email)

  const data = await client.getLanguageDistribution()

  try {
    return Response.json(data)
  } catch (error) {
    return new Response('error', {
      status: 400,
    })
  }
}
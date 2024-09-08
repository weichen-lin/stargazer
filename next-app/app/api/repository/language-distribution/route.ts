import { GetUser } from '@/actions'
import RepositoryClient from '@/client/repository-client'

export const dynamic = 'force-dynamic' // defaults to auto
export async function GET() {
  const { email } = await GetUser()
  const client = new RepositoryClient(email)

  const data = await client.getLanguageDistribution()
  console.log({ data })

  try {
    return new Response('OK', {
      headers: {
        'Content-Type': 'application/json',
      },
    })
  } catch (error) {
    return new Response('error', {
      status: 400,
    })
  }
}

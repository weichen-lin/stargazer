import { GetUser } from '@/actions'
import { RepositoryClient } from '@/client/repository'

export const dynamic = 'force-dynamic'
export async function GET() {
  const { email } = await GetUser()

  try {
    const client = new RepositoryClient(email)

    await client.syncRepository()

    return Response.json({ status: 'ok' })
  } catch (error) {
    return new Response('error', {
      status: 400,
    })
  }
}

import { GetUser } from '@/actions'
import { RepositoryClient } from '@/client/repository'
import { isAxiosError } from 'axios'

export const dynamic = 'force-dynamic'
export async function GET() {
  const { email } = await GetUser()

  try {
    const client = new RepositoryClient(email)

    await client.syncRepository()

    return Response.json({ status: 'ok' })
  } catch (error) {
    if (isAxiosError(error)) {
      return new Response(JSON.stringify(error.response?.data), {
        status: error.response?.status,
      })
    }

    return new Response('unknown error', {
      status: 400,
    })
  }
}

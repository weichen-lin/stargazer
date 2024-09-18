import { GetUser } from '@/actions'
import { RepositoryClient } from '@/client/repository'
import { type NextRequest } from 'next/server'

export const dynamic = 'force-dynamic' // defaults to auto
export async function GET(req: NextRequest) {
  const { email } = await GetUser()

  const params = req.nextUrl.searchParams

  try {
    const client = new RepositoryClient(email)

    const data = await client.fullTextSearch(params.get('query') ?? '')
    return Response.json(data)
  } catch (error) {
    return new Response('error', {
      status: 400,
    })
  }
}

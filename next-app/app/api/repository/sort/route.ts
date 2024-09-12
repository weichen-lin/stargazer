import { GetUser } from '@/actions'
import { RepositoryClient, SortKeySchema } from '@/client/repository'
import { type NextRequest } from 'next/server'

export const dynamic = 'force-dynamic' // defaults to auto
export async function GET(req: NextRequest) {
  const { email } = await GetUser()

  const params = req.nextUrl.searchParams

  try {
    const key = SortKeySchema.parse(params.get('key'))
    const client = new RepositoryClient(email)

    const data = await client.getReposWithSortKey(key)
    return Response.json(data.data)
  } catch (error) {
    return new Response('error', {
      status: 400,
    })
  }
}

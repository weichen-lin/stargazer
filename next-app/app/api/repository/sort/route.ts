import { RepositoryClient, SortKeySchema } from '@/client/repository'
import { type NextRequest } from 'next/server'

export const dynamic = 'force-dynamic' // defaults to auto
export async function GET(req: NextRequest) {
  const params = req.nextUrl.searchParams

  const key = SortKeySchema.parse(params.get('key'))
  const client = new RepositoryClient()

  const { status_code, data } = await client.getReposWithSortKey(key)

  return Response.json(data, { status: status_code })
}

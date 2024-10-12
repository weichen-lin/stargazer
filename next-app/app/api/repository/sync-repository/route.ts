import { RepositoryClient } from '@/client/repository'

export const dynamic = 'force-dynamic'
export async function GET() {
  const client = new RepositoryClient()

  const { data, status_code } = await client.syncRepository()

  return Response.json(data, { status: status_code })
}

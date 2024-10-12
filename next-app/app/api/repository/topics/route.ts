import { RepositoryClient } from '@/client/repository'

export const dynamic = 'force-dynamic' // defaults to auto
export async function GET() {
  const client = new RepositoryClient()

  const { data, status_code } = await client.getTopics()

  return Response.json(data, { status: status_code })
}

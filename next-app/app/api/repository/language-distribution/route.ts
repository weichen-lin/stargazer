import { RepositoryClient } from '@/client/repository'

export const dynamic = 'force-dynamic' // defaults to auto
export async function GET() {
  const client = new RepositoryClient()

  const { status_code, data } = await client.getLanguageDistribution()

  return Response.json(data, { status: status_code })
}

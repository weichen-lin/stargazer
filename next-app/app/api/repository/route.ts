import { RepositoryClient } from '@/client/repository'
import { type NextRequest } from 'next/server'

export const dynamic = 'force-dynamic' // defaults to auto
export async function GET(req: NextRequest) {
  const params = req.nextUrl.searchParams
  const languages = params.get('languages') ?? ''
  const page = params.get('page') ?? '1'
  const limit = params.get('limit') ?? '20'

  const client = new RepositoryClient()

  const { status_code, data } = await client.getReposWithLanguage(languages, page, limit)

  return Response.json(data, { status: status_code })
}

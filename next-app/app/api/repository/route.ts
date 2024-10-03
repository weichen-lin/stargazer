import { GetUser } from '@/actions'
import { RepositoryClient, SortKeySchema } from '@/client/repository'
import { type NextRequest } from 'next/server'

export const dynamic = 'force-dynamic' // defaults to auto
export async function GET(req: NextRequest) {
  const { email } = await GetUser()

  const params = req.nextUrl.searchParams

  try {
    const languages = params.get('languages') ?? ''
    const page = params.get('page') ?? '1'
    const limit = params.get('limit') ?? '20'

    const client = new RepositoryClient(email)

    const data = await client.getReposWithLanguage(languages, page, limit)
    return Response.json(data)
  } catch (error) {
    return new Response('error', {
      status: 400,
    })
  }
}

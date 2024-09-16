import { GetUser } from '@/actions'
import trendClient from '@/client/trends'
import { DateRange } from '@/client/trends/type'
import { type NextRequest } from 'next/server'

export const dynamic = 'force-dynamic' // defaults to auto
export async function GET(req: NextRequest) {
  await GetUser()

  const params = req.nextUrl.searchParams

  const since = (params.get('since') ?? 'daily') as DateRange
  const language = params.get('language')

  const data = await trendClient.getTrendRepos({ since, language })

  try {
    return Response.json(data)
  } catch (error) {
    return new Response('error', {
      status: 400,
    })
  }
}

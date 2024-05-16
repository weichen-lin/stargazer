import { getUserRepos } from '@/actions/neo4j'
import { GetUser } from '@/actions'

export const dynamic = 'force-dynamic' // defaults to auto
export async function POST(request: Request) {
  const { email } = await GetUser()
  const body = await request.json()

  try {
    const { total, repos } = await getUserRepos({ email, ...body })

    return new Response(JSON.stringify({ total, repos }), {
      headers: {
        'Content-Type': 'application/json',
      },
    })
  } catch (error) {
    return new Response('error', {
      status: 400,
    })
  }
}

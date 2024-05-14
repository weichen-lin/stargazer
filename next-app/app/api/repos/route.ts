import { getServerSession } from 'next-auth'
import { getUserRepos } from '@/actions/neo4j'

export const dynamic = 'force-dynamic' // defaults to auto
export async function POST(request: Request) {
  const body = await request.json()
  const session = await getServerSession()

  const user = session?.user
  if (!user) {
    return Response.redirect('/')
  }

  try {
    const { total, repos } = await getUserRepos({ email: user.email, ...body })

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

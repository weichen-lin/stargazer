import { GetUser } from '@/actions'
import { CollectionClient } from '@/client/collection'
import { type NextRequest } from 'next/server'

export const dynamic = 'force-dynamic' // defaults to auto
export async function POST(req: NextRequest, { params }: { params: { id: string } }) {
  const { email } = await GetUser()
  const { id } = params

  const data = await req.json()

  const searchParams = req.nextUrl.searchParams

  try {
    const page = searchParams.get('page') ?? '1'
    const limit = searchParams.get('limit') ?? '20'

    const client = new CollectionClient(email)

    const data = await client.getReposCollections(id, page, limit)
    return Response.json(data)
  } catch (error) {
    return new Response('error', {
      status: 400,
    })
  }
}

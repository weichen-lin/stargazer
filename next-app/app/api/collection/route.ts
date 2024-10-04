import { GetUser } from '@/actions'
import { type NextRequest } from 'next/server'
import { CollectionClient } from '@/client/collection'

export const dynamic = 'force-dynamic' // defaults to auto
export async function GET(req: NextRequest) {
  const { email } = await GetUser()

  const params = req.nextUrl.searchParams

  try {
    const page = params.get('page') ?? '1'
    const limit = params.get('limit') ?? '20'

    const client = new CollectionClient(email)

    const data = await client.getCollections(page, limit)
    return Response.json(data)
  } catch (error) {
    return new Response('error', {
      status: 400,
    })
  }
}

export async function POST(req: NextRequest) {
  const { email } = await GetUser()
  const data = await req.json()

  try {
    const client = new CollectionClient(email)

    const res = await client.createCollection(data?.name)
    return Response.json(res)
  } catch (error) {
    return new Response('error', {
      status: 400,
    })
  }
}

export async function DELETE(req: NextRequest) {
  const { email } = await GetUser()
  const data = await req.json()

  try {
    const client = new CollectionClient(email)

    const res = await client.deleteCollection(data?.id)
    return Response.json(res)
  } catch (error) {
    return new Response('error', {
      status: 400,
    })
  }
}

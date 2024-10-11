import { GetUser } from '@/actions'
import { type NextRequest } from 'next/server'
import { CollectionClient } from '@/client/collection'

export const dynamic = 'force-dynamic'
export async function PATCH(req: NextRequest, { params }: { params: { id: string } }) {
  const { email } = await GetUser()
  const { id } = params
  const data = await req.json()

  try {
    const client = new CollectionClient(email)

    const res = await client.updateCollection(id, data)
    return Response.json(res)
  } catch (error) {
    return new Response('error', {
      status: 400,
    })
  }
}

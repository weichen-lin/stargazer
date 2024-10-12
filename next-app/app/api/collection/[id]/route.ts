import { type NextRequest } from 'next/server'
import { CollectionClient } from '@/client/collection'

export const dynamic = 'force-dynamic'
export async function PATCH(req: NextRequest, { params }: { params: { id: string } }) {
  const { id } = params
  const req_json = await req.json()

  const client = new CollectionClient()

  const { data, status_code } = await client.updateCollection(id, req_json)
  return Response.json(data, { status: status_code })
}

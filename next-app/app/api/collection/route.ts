import { type NextRequest } from 'next/server'
import { CollectionClient } from '@/client/collection'

export const dynamic = 'force-dynamic' // defaults to auto
export async function GET(req: NextRequest) {
  const params = req.nextUrl.searchParams
  const page = params.get('page') ?? '1'
  const limit = params.get('limit') ?? '20'

  const client = new CollectionClient()

  const { data, status_code } = await client.getCollections(page, limit)
  return Response.json(data, { status: status_code })
}

export async function POST(req: NextRequest) {
  const req_json = await req.json()

  const client = new CollectionClient()

  const { data, status_code } = await client.createCollection(req_json?.name)
  return Response.json(data, { status: status_code })
}

export async function DELETE(req: NextRequest) {
  const req_json = await req.json()

  const client = new CollectionClient()

  const { data, status_code } = await client.deleteCollection(req_json?.id)
  return Response.json(data, { status: status_code })
}

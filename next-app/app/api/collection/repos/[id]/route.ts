import { CollectionClient } from '@/client/collection'
import { type NextRequest } from 'next/server'

export const dynamic = 'force-dynamic' // defaults to auto
export async function GET(req: NextRequest, { params }: { params: { id: string } }) {
  const { id } = params

  const searchParams = req.nextUrl.searchParams

  const page = searchParams.get('page') ?? '1'
  const limit = searchParams.get('limit') ?? '20'

  const client = new CollectionClient()

  const { data, status_code } = await client.getReposCollections(id, page, limit)
  return Response.json(data, { status: status_code })
}

export async function POST(req: NextRequest, { params }: { params: { id: string } }) {
  const { id } = params

  const req_json = await req.json()

  const client = new CollectionClient()

  const { status_code, data } = await client.addRepoToCollection(id, req_json?.repo_ids)
  return Response.json(data, { status: status_code })
}

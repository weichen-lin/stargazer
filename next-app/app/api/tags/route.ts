import type { NextRequest } from 'next/server'
import TagClient, { ITagPayload } from '@/client/tag'

export async function POST(req: NextRequest) {
  const req_json = await req.json()

  const client = new TagClient()
  const { data, status_code } = await client.createTag(req_json)

  return Response.json(data, { status: status_code })
}

export async function DELETE(req: NextRequest) {
  const req_json = await req.json()

  const client = new TagClient()
  const { data, status_code } = await client.deleteTag(req_json)
  return Response.json(data, { status: status_code })
}

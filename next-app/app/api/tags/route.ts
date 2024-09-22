import type { NextRequest } from 'next/server'
import { GetUser } from '@/actions'
import TagClient, { ITagPayload } from '@/client/tag'

export async function POST(req: NextRequest) {
  const { email } = await GetUser()
  const data = await req.json()

  try {
    const client = new TagClient(email)
    await client.createTag(data)

    return new Response(JSON.stringify('ok'))
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
    const client = new TagClient(email)
    await client.deleteTag(data)

    return new Response(JSON.stringify('ok'))
  } catch (error) {
    return new Response('error', {
      status: 400,
    })
  }
}

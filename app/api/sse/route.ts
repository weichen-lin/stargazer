import { NextRequest } from 'next/server'
import { getServerSession } from 'next-auth'
import { getUserStarsRelationRepos, updateStarRelation } from '@/actions/neo4j'

export const runtime = 'nodejs'
// This is required to enable streaming
export const dynamic = 'force-dynamic'

export const maxDuration = 300

export async function GET(request: NextRequest) {
  // check referer is from the same origin
  const origin = request.headers.get('host')
  const referer = request.headers.get('referer')
  const session = await getServerSession()

  let isClose: boolean = false

  if (!session?.user) {
    return new Response('Forbidden', { status: 401 })
  }

  const { email, name } = session.user
  if (!email || !name) {
    return new Response('Forbidden', { status: 403 })
  }

  if (!referer || !origin) {
    return new Response('Forbidden', { status: 403 })
  }

  const { hostname: refererHost } = new URL(referer)
  const { hostname: originHost } = new URL(`http://${origin}`)

  if (originHost !== refererHost) {
    return new Response('Forbidden', { status: 403 })
  }

  let responseStream = new TransformStream()
  const writer = responseStream.writable.getWriter()

  const encoder = new TextEncoder()

  const data = await getUserStarsRelationRepos(email)
  const current = data.filter(e => e.isVectorized).length

  // Function to send data to the client
  async function sendData(data: { current: number; total: number; error: boolean; status: number }) {
    const formattedData = `data: ${JSON.stringify(data)}\n\n`
    try {
      await writer.write(encoder.encode(formattedData))
    } catch (error) {
      throw new Error('Error occurred while sending data')
    }
  }

  // Initial Progress
  sendData({ current, total: data.length, error: false, status: 200 })

  async function StartVectorize() {
    for (let i = current; i < data.length && !isClose; i++) {
      try {
        const repo_id = data[i].repo_id
        const res = await fetch(`${process.env.TRANSFORMER_URL}/vectorize`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${process.env.AUTHENTICATION_TOKEN}`,
          },
          body: JSON.stringify({
            name,
            repo_id,
          }),
        })
        if (res.status === 201 || res.status === 200) {
          await updateStarRelation(email as string, repo_id, true)
          await sendData({ current: i, total: data.length, error: false, status: res.status })
        } else {
          await sendData({ current: i, total: data.length, error: true, status: res.status })
        }
      } catch (error) {
        console.log('error occur', error)
        await writer.close()
        break
      }
    }
  }

  StartVectorize()

  // Close if client disconnects
  request.signal.onabort = async () => {
    console.log('closing writer')
    await writer.close()
  }

  return new Response(responseStream.readable, {
    headers: {
      'Content-Type': 'text/event-stream',
      Connection: 'keep-alive',
      'Cache-Control': 'no-cache, no-transform',
    },
  })
}

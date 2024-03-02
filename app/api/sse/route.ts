import { NextRequest } from 'next/server'
import { getServerSession } from 'next-auth'
import { getUserStarsRelationRepos } from '@/actions/neo4j'

export const runtime = 'nodejs'
// This is required to enable streaming
export const dynamic = 'force-dynamic'

export async function GET(request: NextRequest) {
  // check referer is from the same origin
  const origin = request.headers.get('host')
  const referer = request.headers.get('referer')
  const session = await getServerSession()

  if (!session?.user) {
    return new Response('Forbidden', { status: 401 })
  }

  const { email } = session.user
  if (!email) {
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
  function sendData(data: { current: number; total: number }) {
    const formattedData = `data: ${JSON.stringify(data)}\n\n`
    writer.write(encoder.encode(formattedData))
  }

  // Initial Progress
  sendData({ current, total: data.length })

  async function StartVectorize() {
    for (let i = current; i < data.length; i++) {
      await delay(10)
      sendData({ current: i, total: data.length })
    }
    writer.close()
  }

  StartVectorize()

  // Close if client disconnects
  request.signal.onabort = () => {
    console.log('closing writer')
    writer.close()
  }

  return new Response(responseStream.readable, {
    headers: {
      'Content-Type': 'text/event-stream',
      Connection: 'keep-alive',
      'Cache-Control': 'no-cache, no-transform',
    },
  })
}

function delay(ms: number) {
  return new Promise(resolve => setTimeout(resolve, ms))
}

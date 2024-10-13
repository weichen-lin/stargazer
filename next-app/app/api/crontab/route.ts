import { GetUser } from '@/actions'
import CrontabClient from '@/client/crontab'
import { type NextRequest } from 'next/server'

export const dynamic = 'force-dynamic' // defaults to auto
export async function GET() {
  const client = new CrontabClient()

  const { data, status_code } = await client.getCrontab()

  return Response.json(data, { status: status_code })
}

export async function POST() {
  const client = new CrontabClient()

  const { data, status_code } = await client.createCronTab()

  return Response.json(data, { status: status_code })
}

export async function PATCH(req: NextRequest) {
  const client = new CrontabClient()
  const params = req.nextUrl.searchParams

  const triggered_at = params.get('triggered_at') as string

  if (!triggered_at) {
    return new Response('error', {
      status: 400,
    })
  }

  const { data, status_code } = await client.updateCronTab(triggered_at)

  return Response.json(data, { status: status_code })
}

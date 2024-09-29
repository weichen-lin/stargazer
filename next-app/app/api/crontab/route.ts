import { GetUser } from '@/actions'
import CrontabClient from '@/client/crontab-client'
import { type NextRequest } from 'next/server'

export const dynamic = 'force-dynamic' // defaults to auto
export async function GET() {
  const { email } = await GetUser()
  const client = new CrontabClient(email)

  const data = await client.getCrontab()

  try {
    return Response.json(data)
  } catch (error) {
    return new Response('error', {
      status: 400,
    })
  }
}

export async function POST() {
  const { email } = await GetUser()
  const client = new CrontabClient(email)

  const data = await client.createCronTab()

  try {
    return Response.json(data)
  } catch (error) {
    return new Response('error', {
      status: 400,
    })
  }
}

export async function PATCH(req: NextRequest) {
  const { email } = await GetUser()
  const client = new CrontabClient(email)
  const params = req.nextUrl.searchParams

  const hour = params.get('hour') as string

  if (!hour) {
    return new Response('error', {
      status: 400,
    })
  }

  const data = await client.updateCronTab(hour)

  try {
    return Response.json(data)
  } catch (error) {
    return new Response('error', {
      status: 400,
    })
  }
}

import { GetUser } from '@/actions'
import CrontabClient from '@/client/crontab-client'

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

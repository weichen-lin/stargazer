'use server'

import { getUserProviderInfo } from './neo4j'

const PRODUCER_URL = process.env.PRODUCER_URL
const TOKEN = process.env.AUTHENTICATION_TOKEN

interface syncUserStarsParams {
  user_id: string
  user_name: string
  page: number
}

export async function syncUserStars(name: string): Promise<{ status: number; title: string; message: string }> {
  const providerId = await getUserProviderInfo(name)
  const params: syncUserStarsParams = {
    user_id: providerId,
    user_name: name,
    page: 1,
  }

  console.log({ url: PRODUCER_URL })
  const response = await fetch(`${PRODUCER_URL}/get_user_stars`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${TOKEN}`,
    },
    body: JSON.stringify(params),
  })

  const status = response.status
  let message: string
  let title: string

  if (status === 200) {
    title = 'Scheduled: Catch up'
    message = "You will receive a confirmation email once it's completed."
  } else if (status === 409) {
    const json = await response.json()
    title = 'In Progress: Catch up'
    message = `You'll be able to make another request after ${json.expires}`
  } else {
    title = 'An error occurred'
    message = 'Please try again later.'
  }

  return {
    status,
    title,
    message,
  }
}

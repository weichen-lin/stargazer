'use server'

import { getUserProviderInfo } from './neo4j'

const PRODUCER_URL = process.env.PRODUCER_URL
const TOKEN = process.env.AUTHENTICATION_TOKEN

interface syncUserStarsParams {
  user_id: string
  username: string
  page: number
}

export async function syncUserStars(name: string): Promise<{ status: number; message: string }> {
  const providerId = await getUserProviderInfo(name)
  const params: syncUserStarsParams = {
    user_id: providerId,
    username: name,
    page: 1,
  }

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

  if (status === 200) {
    message = "We have received your request, and you will receive a confirmation email once it's completed."
  } else if (status === 409) {
    const json = await response.json()
    message = `Your request is currently being processed. You'll be able to make another request after ${json.expires}`
  } else {
    message = 'An error occurred while processing your request. Please try again later.'
  }

  return {
    status,
    message,
  }
}

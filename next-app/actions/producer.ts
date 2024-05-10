'use server'
import { generateAccessToken } from './util'

const PRODUCER_URL = process.env.PRODUCER_URL
const TOKEN = process.env.AUTHENTICATION_TOKEN

interface syncUserStarsParams {
  user_id: string
  user_name: string
  page: number
}

export async function syncUserStars(email: string): Promise<{ status: number; title: string; message: string }> {
  const jwtToken = await generateAccessToken(email)
  const response = await fetch(`${PRODUCER_URL}/get_user_stars`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      Authorization: jwtToken,
    },
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

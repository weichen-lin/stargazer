'use server'
import { generateAccessToken } from './util'

const PRODUCER_URL = process.env.PRODUCER_URL

export async function syncUserStars(email: string): Promise<{ status: number; title: string; message: string }> {
  const jwtToken = await generateAccessToken(email)
  const response = await fetch(`${process.env.KAFKA_SERVICE}/get_user_stars`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      Authorization: jwtToken,
    },
  })

  const status = response.status

  if (status === 200) {
    return {
      status,
      title: 'Scheduled: Catch up',
      message: "You will receive a confirmation email once it's completed.",
    }
  } else if (status === 409) {
    const json = await response.json()
    return {
      status,
      title: 'In Progress: Catch up',
      message: `You'll be able to make another request after ${json.expires}`,
    }
  } else {
    return {
      status,
      title: 'An error occurred',
      message: 'Please try again later.',
    }
  }
}

export async function updateCrontabHour(email: string, hour: number): Promise<boolean> {
  const jwtToken = await generateAccessToken(email)
  const response = await fetch(`${process.env.KAFKA_SERVICE}/update_cron_tab_setting?hour=${hour}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
      Authorization: jwtToken,
    },
  })

  const status = response.status
  if (status !== 200) {
    return false
  }

  return true
}

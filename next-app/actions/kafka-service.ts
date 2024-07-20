'use server'

import logger from '@/utils/logger'
import { generateAccessToken } from './util'

const fetcher = async <T>(email: string, path: string, data: T | null) => {
  const jwtToken = await generateAccessToken(email)
  const now = Date.now()
  try {
    const res = fetch(`${process.env.KAFKA_SERVICE}/${path}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        Authorization: jwtToken,
      },
      body: data ? JSON.stringify(data) : undefined,
    })
    return res
  } catch (error) {
    logger.error({ error, path, email, payload: data })
  } finally {
    logger.info({ path, email, duration: Date.now() - now, payload: data })
  }
}

export async function syncUserStars(email: string): Promise<{ status: number; title: string; message: string }> {
  const res = await fetcher(email, 'get_user_stars', null)

  if (!res) {
    return {
      status: 500,
      title: 'An error occurred',
      message: 'Please try again later.',
    }
  }

  const status = res.status

  if (status === 200) {
    return {
      status,
      title: 'Scheduled: Catch up',
      message: "You will receive a confirmation email once it's completed.",
    }
  } else if (status === 409) {
    const json = await res.json()
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
  const res = await fetcher(email, `update_cron_tab_setting?hour=${hour}`, null)
  if (!res || res.status !== 200) {
    return false
  }

  return true
}

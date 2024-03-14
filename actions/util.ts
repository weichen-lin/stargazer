'use server'

import { sign } from 'jsonwebtoken'

export async function generateAccessToken(user: string) {
  const expired = new Date()
  expired.setMinutes(expired.getMinutes() + 1)
  const payload = {
    user_name: user,
    create_at: new Date(),
    expire_at: expired,
  }
  const TOKEN = process.env.AUTHENTICATION_TOKEN ?? ''

  const token = sign(payload, TOKEN, { algorithm: 'HS256' })

  return token
}

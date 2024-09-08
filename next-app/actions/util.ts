'use server'

import { sign } from 'jsonwebtoken'

const secret = process.env.JWT_SECRET ?? ''

export async function generateAccessToken(email: string) {
  const expired = new Date()
  expired.setMinutes(expired.getMinutes() + 1)

  const payload = {
    email: email,
    create_at: new Date(),
    expire_at: expired,
  }

  const token = sign(payload, secret, { algorithm: 'HS256' })

  return token
}

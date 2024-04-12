'use server'

import { getServerSession } from 'next-auth'

export const GetSuggesions = async (query: string) => {
  const session = await getServerSession()
  const user = session?.user
  if (!user) {
    return []
  }

  const res = await fetch(`${process.env.TRANSFORMER_URL}/get_suggestions`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${process.env.AUTHENTICATION_TOKEN}`,
    },
    body: JSON.stringify({
      message: query,
      name: user.name,
    }),
  })

  if (res.ok) {
    const data = await res.json()
    return data.items
  }

  return []
}

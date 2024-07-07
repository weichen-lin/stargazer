'use server'

import { GetUser } from './user'
import { ISuggestion } from './neo4j/repos'

export const GetSuggesions = async (
  query: string,
): Promise<{
  status: number
  items: ISuggestion[]
}> => {
  const { email } = await GetUser()

  try {
    const res = await fetch(`${process.env.TRANSFORMER_URL}/get_suggestions`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${process.env.AUTHENTICATION_TOKEN}`,
      },
      body: JSON.stringify({
        query,
        email,
      }),
    })

    if (res.ok) {
      const data = await res.json()
      return {
        status: res.status,
        items: data.items,
      }
    }

    return {
      status: res.status,
      items: [],
    }
  } catch (error) {
    console.error('Error fetching suggestions: ', error)
    return {
      status: 500,
      items: [],
    }
  }
}

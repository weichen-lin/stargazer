'use server'

import { getServerSession } from 'next-auth'

export interface ISuggestion {
  avatar_url: string
  description: string | null
  full_name: string
  html_url: string
  readme_summary: string
}

export const GetSuggesions = async (query: string): Promise<ISuggestion[]> => {
  const session = await getServerSession()
  const user = session?.user
  if (!user) {
    return []
  }

  try {
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
  } catch (error) {
    console.error('Error fetching suggestions: ', error)
    return []
  }
}

export const GetFullTextSearch = async (query: string): Promise<ISuggestion[]> => {
  const session = await getServerSession()
  const user = session?.user
  if (!user) {
    return []
  }

  try {
    const res = await fetch(`${process.env.TRANSFORMER_URL}/full_text_search`, {
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
  } catch (error) {
    console.error('Error fetching suggestions: ', error)
    return []
  }
}

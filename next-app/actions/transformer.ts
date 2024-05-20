'use server'

import { GetUser } from './user'

export interface ISuggestion {
  repo_id: number
  avatar_url: string
  description: string | null
  full_name: string
  html_url: string
}

export const GetSuggesions = async (query: string): Promise<ISuggestion[] | boolean> => {
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
      return data.items
    }

    return []
  } catch (error) {
    console.error('Error fetching suggestions: ', error)
    return false
  }
}

export const GetFullTextSearch = async (query: string): Promise<ISuggestion[]> => {
  const { email } = await GetUser()

  try {
    const res = await fetch(`${process.env.TRANSFORMER_URL}/full_text_search`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${process.env.AUTHENTICATION_TOKEN}`,
      },
      body: JSON.stringify({
        email,
        query,
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

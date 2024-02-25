'use server'

export const GetSuggesions = async (query: string) => {
  const res = await fetch('', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: ``,
    },
    body: JSON.stringify({
      message: query,
    }),
  })

  if (res.ok) {
    const data = await res.json()
    return data.items
  }

  return []
}
